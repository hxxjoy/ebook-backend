// 这是目前最好的配置使用方式，性能比用map高几十倍
package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

var (
	C    *Config
	once sync.Once
	mu   sync.RWMutex
)

// Config 配置结构体
type Config struct {
	App      AppConfig      `yaml:"app"`
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	Log      LogConfig      `yaml:"log"`
	Cache    CacheConfig    `yaml:"cache"`
	SMTP     SMTPConfig     `yaml:"smtp"`
	JWT      JWTConfig      `yaml:"jwt"`
}

// MustLoad 加载配置文件，如果出错则panic
func MustLoad(filepath string) {
	once.Do(func() {
		if err := load(filepath); err != nil {
			panic(fmt.Sprintf("load config failed: %v", err))
		}
	})
}

// Load 加载配置文件
func Load(filepath string) error {
	mu.Lock()
	defer mu.Unlock()
	return load(filepath)
}
func load(configPath string) error {
	// 检查文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("config file not found: %s", configPath)
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("read config file failed: %w", err)
	}

	// 创建新的配置实例
	config := new(Config)
	if err := yaml.Unmarshal(data, config); err != nil {
		return fmt.Errorf("unmarshal config failed: %w", err)
	}

	// 加载 .env 文件
	if err := loadEnv(config); err != nil {
		return fmt.Errorf("load env failed: %w", err)
	}

	// 验证配置
	if err := config.validate(); err != nil {
		return fmt.Errorf("validate config failed: %w", err)
	}

	// 设置默认值
	config.setDefaults()

	// 更新全局配置
	C = config
	return nil
}

// loadEnv 加载环境变量
func loadEnv(c *Config) error {
	// 尝试加载 .env 文件，如果文件不存在也不报错
	_ = godotenv.Load()

	// App 配置
	if env := os.Getenv("APP_NAME"); env != "" {
		c.App.Name = env
	}
	if env := os.Getenv("APP_MODE"); env != "" {
		c.App.Mode = env
	}
	if env := os.Getenv("APP_PAGE_SIZE"); env != "" {
		if v, err := strconv.Atoi(env); err == nil {
			c.App.PageSize = v
		}
	}
	if env := os.Getenv("APP_MAX_FILE_SIZE"); env != "" {
		if v, err := strconv.Atoi(env); err == nil {
			c.App.MaxFileSize = v
		}
	}
	if env := os.Getenv("APP_JWT_TIMEOUT"); env != "" {
		if v, err := time.ParseDuration(env); err == nil {
			c.App.JWTTimeout = v
		}
	}

	// Server 配置
	if env := os.Getenv("SERVER_HOST"); env != "" {
		c.Server.Host = env
	}
	if env := os.Getenv("SERVER_PORT"); env != "" {
		if v, err := strconv.Atoi(env); err == nil {
			c.Server.Port = v
		}
	}
	if env := os.Getenv("SERVER_READ_TIMEOUT"); env != "" {
		if v, err := time.ParseDuration(env); err == nil {
			c.Server.ReadTimeout = v
		}
	}
	if env := os.Getenv("SERVER_WRITE_TIMEOUT"); env != "" {
		if v, err := time.ParseDuration(env); err == nil {
			c.Server.WriteTimeout = v
		}
	}
	// Cache 配置
	if env := os.Getenv("CACHE_ENABLED"); env != "" {
		c.Cache.Enabled = env == "true"
	}
	if env := os.Getenv("CACHE_TTL"); env != "" {
		if v, err := strconv.Atoi(env); err == nil {
			c.Cache.TTL = v
		}
	}
	if env := os.Getenv("CACHE_CLEARUP"); env != "" {
		if v, err := strconv.Atoi(env); err == nil {
			c.Cache.Clearup = v
		}
	}

	// Database 配置
	if env := os.Getenv("DB_DRIVER"); env != "" {
		c.Database.Driver = env
	}
	if env := os.Getenv("DB_HOST"); env != "" {
		c.Database.Host = env
	}
	if env := os.Getenv("DB_MAX_OPEN_CONNS"); env != "" {
		if v, err := strconv.Atoi(env); err == nil {
			c.Database.MaxOpenConns = v
		}
	}
	if env := os.Getenv("DB_MAX_IDLE_CONNS"); env != "" {
		if v, err := strconv.Atoi(env); err == nil {
			c.Database.MaxIdleConns = v
		}
	}
	if env := os.Getenv("DB_CONN_MAX_LIFETIME"); env != "" {
		if v, err := time.ParseDuration(env); err == nil {
			c.Database.ConnMaxLifetime = v
		}
	}
	if env := os.Getenv("DB_CONN_MAX_IDLE_TIME"); env != "" {
		if v, err := time.ParseDuration(env); err == nil {
			c.Database.ConnMaxIdleTime = v
		}
	}

	// Redis 配置
	if env := os.Getenv("REDIS_HOST"); env != "" {
		c.Redis.Host = env
	}
	if env := os.Getenv("REDIS_PASSWORD"); env != "" {
		c.Redis.Password = env
	}
	if env := os.Getenv("REDIS_DB"); env != "" {
		if v, err := strconv.Atoi(env); err == nil {
			c.Redis.DB = v
		}
	}
	if env := os.Getenv("REDIS_POOL_SIZE"); env != "" {
		if v, err := strconv.Atoi(env); err == nil {
			c.Redis.PoolSize = v
		}
	}
	if env := os.Getenv("REDIS_MAX_CONN_AGE"); env != "" {
		if v, err := time.ParseDuration(env); err == nil {
			c.Redis.MaxConnAge = v
		}
	}
	//smtp
	if env := os.Getenv("SMTP_HOST"); env != "" {
		c.SMTP.Host = env
	}
	if env := os.Getenv("SMTP_USERNAME"); env != "" {
		c.SMTP.Username = env
	}
	if env := os.Getenv("SMTP_PORT"); env != "" {
		if v, err := strconv.Atoi(env); err == nil {
			c.SMTP.Port = v
		}
	}
	if env := os.Getenv("SMTP_PASSWORD"); env != "" {
		c.SMTP.Password = env
	}
	//jwt
	if env := os.Getenv("JWT_ACCESS_TOKEN_SECRET"); env != "" {
		c.JWT.AccessTokenSecret = env
	}
	if env := os.Getenv("JWT_REFRESH_TOKEN_SECRET"); env != "" {
		c.JWT.RefreshTokenSecret = env
	}
	if env := os.Getenv("JWT_ACCESS_TOKEN_EXPIRE"); env != "" {
		if v, err := time.ParseDuration(env); err == nil {
			c.JWT.AccessTokenExpire = v
		}
	}
	if env := os.Getenv("JWT_REFRESH_TOKEN_EXPIRE"); env != "" {
		if v, err := time.ParseDuration(env); err == nil {
			c.JWT.RefreshTokenExpire = v
		}
	}
	// Log 配置
	if env := os.Getenv("LOG_LEVEL"); env != "" {
		c.Log.Level = env
	}
	if env := os.Getenv("LOG_FILENAME"); env != "" {
		c.Log.Filename = env
	}
	if env := os.Getenv("LOG_MAX_SIZE"); env != "" {
		if v, err := strconv.Atoi(env); err == nil {
			c.Log.MaxSize = v
		}
	}
	if env := os.Getenv("LOG_MAX_BACKUPS"); env != "" {
		if v, err := strconv.Atoi(env); err == nil {
			c.Log.MaxBackups = v
		}
	}
	if env := os.Getenv("LOG_MAX_AGE"); env != "" {
		if v, err := strconv.Atoi(env); err == nil {
			c.Log.MaxAge = v
		}
	}

	return nil
}

// validate 验证配置
func (c *Config) validate() error {
	if c.App.Name == "" {
		return errors.New("app name is required")
	}

	if c.App.Mode != "dev" && c.App.Mode != "test" && c.App.Mode != "prod" {
		return errors.New("invalid app mode")
	}

	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return errors.New("invalid server port")
	}

	if c.Database.Driver == "" {
		return errors.New("database driver is required")
	}

	if c.Database.MaxOpenConns <= 0 {
		return errors.New("invalid database max open conns")
	}

	return nil
}

// setDefaults 设置默认值
func (c *Config) setDefaults() {
	// App 默认值
	if c.App.PageSize <= 0 {
		c.App.PageSize = 10
	}
	if c.App.MaxFileSize <= 0 {
		c.App.MaxFileSize = 50 // 50MB
	}
	if c.App.JWTTimeout <= 0 {
		c.App.JWTTimeout = 24 * time.Hour
	}

	// Server 默认值
	if c.Server.Host == "" {
		c.Server.Host = "0.0.0.0"
	}
	if c.Server.ReadTimeout <= 0 {
		c.Server.ReadTimeout = 10 * time.Second
	}
	if c.Server.WriteTimeout <= 0 {
		c.Server.WriteTimeout = 10 * time.Second
	}

	// Database 默认值
	if c.Database.MaxIdleConns <= 0 {
		c.Database.MaxIdleConns = 10
	}
	if c.Database.ConnMaxLifetime <= 0 {
		c.Database.ConnMaxLifetime = time.Hour
	}
	if c.Database.ConnMaxIdleTime <= 0 {
		c.Database.ConnMaxIdleTime = 30 * time.Minute
	}

	// Redis 默认值
	if c.Redis.PoolSize <= 0 {
		c.Redis.PoolSize = 10
	}
	if c.Redis.MaxConnAge <= 0 {
		c.Redis.MaxConnAge = time.Hour
	}

	// Log 默认值
	if c.Log.Level == "" {
		c.Log.Level = "info"
	}
	if c.Log.MaxSize <= 0 {
		c.Log.MaxSize = 100
	}
	if c.Log.MaxBackups <= 0 {
		c.Log.MaxBackups = 10
	}
	if c.Log.MaxAge <= 0 {
		c.Log.MaxAge = 30
	}

}

// IsDev 是否为开发环境
func (c *Config) IsDev() bool {
	return c.App.Mode == "dev"
}

// IsTest 是否为测试环境
func (c *Config) IsTest() bool {
	return c.App.Mode == "test"
}

// IsProd 是否为生产环境
func (c *Config) IsProd() bool {
	return c.App.Mode == "prod"
}
