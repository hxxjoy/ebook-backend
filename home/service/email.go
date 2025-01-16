package service

import (
    "context"
    "encoding/json"
    "fmt"
    "math/rand"
    "net/smtp"
    "time"
    
    "github.com/redis/go-redis/v9"
    "gorm.io/gorm"
	"book-fiber/system/config"

    //"book-fiber/model"
)

type EmailService interface {
    SendVerificationCode(email, ip string) error
    VerifyCode(email, code, ip string) bool
    //GetVerificationLogs(email string) []model.VerificationLog
}

type emailService struct {
    smtpHost     string
    smtpPort     int
    smtpUsername string
    smtpPassword string
    redis        *redis.Client
    db           *gorm.DB
}

type EmailLimitConfig struct {
    MaxAttemptsPerIP     int
    MaxAttemptsPerEmail  int
    MaxVerifyAttempts    int
    CodeExpiration       time.Duration
    IPBlockDuration      time.Duration
    EmailBlockDuration   time.Duration
    RequestCooldown      time.Duration
}

var defaultConfig = EmailLimitConfig{
    MaxAttemptsPerIP:     10,
    MaxAttemptsPerEmail:  5,
    MaxVerifyAttempts:    5,
    CodeExpiration:       10 * time.Minute,
    IPBlockDuration:      24 * time.Hour,
    EmailBlockDuration:   24 * time.Hour,
    RequestCooldown:      60 * time.Second,
}

type codeInfo struct {
    Code      string    `json:"code"`
    CreatedAt time.Time `json:"created_at"`
    Attempts  int       `json:"attempts"`
}

func NewEmailService(config *config.Config, redis *redis.Client, db *gorm.DB) EmailService {
    return &emailService{
        smtpHost:     config.SMTP.Host,
        smtpPort:     config.SMTP.Port,
        smtpUsername: config.SMTP.Username,
        smtpPassword: config.SMTP.Password,
        redis:        redis,
        db:           db,
    }
}

// 检查IP是否被封禁
func (s *emailService) isIPBlocked(ctx context.Context, ip string) bool {
    key := fmt.Sprintf("blocked_ip:%s", ip)
    exists, _ := s.redis.Exists(ctx, key).Result()
    return exists == 1
}

// 检查邮箱是否被封禁
func (s *emailService) isEmailBlocked(ctx context.Context, email string) bool {
    key := fmt.Sprintf("blocked_email:%s", email)
    exists, _ := s.redis.Exists(ctx, key).Result()
    return exists == 1
}

// 检查发送频率
func (s *emailService) checkRateLimit(ctx context.Context, email, ip string) error {
    // 检查IP是否被封禁
    if s.isIPBlocked(ctx, ip) {
        return fmt.Errorf("ip is blocked")
    }

    // 检查邮箱是否被封禁
    if s.isEmailBlocked(ctx, email) {
        return fmt.Errorf("email is blocked")
    }

    // 检查IP最近一次请求时间
    ipLastReqKey := fmt.Sprintf("last_req:ip:%s", ip)
    lastReq, _ := s.redis.Get(ctx, ipLastReqKey).Time()
    if time.Since(lastReq) < defaultConfig.RequestCooldown {
        return fmt.Errorf("please wait before requesting again")
    }

    // 检查IP每小时请求次数
    ipHourlyKey := fmt.Sprintf("hourly:ip:%s", ip)
    ipCount, _ := s.redis.Incr(ctx, ipHourlyKey).Result()
    if ipCount == 1 {
        s.redis.Expire(ctx, ipHourlyKey, time.Hour)
    }
    if ipCount > int64(defaultConfig.MaxAttemptsPerIP) {
        // 封禁IP
        blockKey := fmt.Sprintf("blocked_ip:%s", ip)
        s.redis.Set(ctx, blockKey, "blocked", defaultConfig.IPBlockDuration)
        return fmt.Errorf("too many requests from this IP")
    }

    // 检查邮箱每小时接收次数
    emailHourlyKey := fmt.Sprintf("hourly:email:%s", email)
    emailCount, _ := s.redis.Incr(ctx, emailHourlyKey).Result()
    if emailCount == 1 {
        s.redis.Expire(ctx, emailHourlyKey, time.Hour)
    }
    if emailCount > int64(defaultConfig.MaxAttemptsPerEmail) {
        // 封禁邮箱
        blockKey := fmt.Sprintf("blocked_email:%s", email)
        s.redis.Set(ctx, blockKey, "blocked", defaultConfig.EmailBlockDuration)
        return fmt.Errorf("too many requests for this email")
    }

    // 更新最后请求时间
    s.redis.Set(ctx, ipLastReqKey, time.Now(), time.Hour)
    
    return nil
}

func (s *emailService) SendEmail(to, subject, body string) error {
    addr := fmt.Sprintf("%s:%d", s.smtpHost, s.smtpPort)
    auth := smtp.PlainAuth("", s.smtpUsername, s.smtpPassword, s.smtpHost)
    
    msg := []byte(fmt.Sprintf("To: %s\r\n"+
        "Subject: %s\r\n"+
        "Content-Type: text/plain; charset=UTF-8\r\n"+
        "\r\n"+
        "%s\r\n", to, subject, body))
        
    return smtp.SendMail(addr, auth, s.smtpUsername, []string{to}, msg)
}

func (s *emailService) SendVerificationCode(email, ip string) error {
    ctx := context.Background()
    
    // 检查频率限制
    if err := s.checkRateLimit(ctx, email, ip); err != nil {
        return err
    }
    
    // 生成验证码
    code := fmt.Sprintf("%06d", rand.Intn(1000000))
    
    // 存储验证码
    info := codeInfo{
        Code:      code,
        CreatedAt: time.Now(),
        Attempts:  0,
    }
    
    codeJSON, _ := json.Marshal(info)
    key := fmt.Sprintf("vcode:%s", email)
    result, err := s.redis.Set(ctx, key, codeJSON, defaultConfig.CodeExpiration).Result()
    if err != nil {
        return fmt.Errorf("failed to save verification code in Redis: %w", err)
    }
    if result != "OK" {
        return fmt.Errorf("unexpected Redis response: %s", result)
    }
    
    // 发送邮件
    return s.SendEmail(email, "Verification Code", 
        fmt.Sprintf("Your verification code is: %s\nThis verification code is valid for %d minutes.", 
            code, defaultConfig.CodeExpiration/time.Minute))
}

func (s *emailService) VerifyCode(email, code, ip string) bool {
    ctx := context.Background()

    // 检查IP是否被封禁
    if s.isIPBlocked(ctx, ip) {
        return false
    }

    key := fmt.Sprintf("vcode:%s", email)
    data, err := s.redis.Get(ctx, key).Result()
    if err != nil {
        return false
    }

    var codeInfo struct {
        Code      string    `json:"code"`
        CreatedAt time.Time `json:"created_at"`
        Attempts  int       `json:"attempts"`
    }

    if err := json.Unmarshal([]byte(data), &codeInfo); err != nil {
        return false
    }

    // 增加尝试次数
    codeInfo.Attempts++
    
    // 检查尝试次数
    if codeInfo.Attempts > defaultConfig.MaxVerifyAttempts {
        // 封禁IP
        blockKey := fmt.Sprintf("blocked_ip:%s", ip)
        s.redis.Set(ctx, blockKey, "blocked", defaultConfig.IPBlockDuration)
        s.redis.Del(ctx, key) // 删除验证码
        return false
    }

    // 更新尝试次数
    codeJSON, _ := json.Marshal(codeInfo)
    s.redis.Set(ctx, key, codeJSON, defaultConfig.CodeExpiration)

    // 验证码匹配检查
    if codeInfo.Code != code {
        return false
    }

    // 验证成功，删除验证码
    s.redis.Del(ctx, key)
    return true
}