package service

import (
	"book-fiber/home/repository"
	"book-fiber/model"
	"book-fiber/system/helper"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo          *repository.UserRepository
	blacklistRepo repository.TokenBlacklistRepository
}

// 实现令牌黑名单
type TokenBlacklist struct {
	Slug      string `gorm:"primarykey"`
	ExpiresAt time.Time
}

func NewAuthService(repo *repository.UserRepository, blacklistRepo repository.TokenBlacklistRepository) *AuthService {
	return &AuthService{repo: repo, blacklistRepo: blacklistRepo}
}

func (s *AuthService) IsTokenBlacklisted(slug string) bool {
	return s.blacklistRepo.IsBlacklisted(slug)
}

func (s *AuthService) RevokeToken(token string) error {
	claims, err := helper.ParseAccessToken(token)
	if err != nil {
		return err
	}
	slug,err := uuid.Parse(claims.Slug)
	if err != nil {
		return err
	}
	blacklist := &model.TokenBlacklist{
		Slug:      slug,
		ExpiresAt: claims.ExpiresAt.Time,
	}

	return s.blacklistRepo.Create(blacklist)
}

func (s *AuthService) CheckCredentials(email string, password string) (*model.User, error) {
	one, err := s.repo.GetByEmail(email)
	if err != nil || !(one.ID > 0) {
		return &model.User{}, errors.New("incorrect username or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(*one.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("incorrect username or password")
	}
	// 3. 如果验证通过，返回用户信息
	one.PasswordHash = nil // 不返回密码 hash
	return one, nil
}

func (s *AuthService) Create(user *model.User) (*model.User, error) {
	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	hash := string(hashedPassword)
	user.PasswordHash = &hash
	// 保存到数据库
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	// 清除密码后返回
	user.PasswordHash = nil
	return user, nil
}

func (s *AuthService) Login(email, password string, rememberMe bool) (*model.User, *helper.TokenPair, error) {
	// 1. 验证用户
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, nil, errors.New("invalid credentials")
	}

	// 2. 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(password)); err != nil {
		return nil, nil, errors.New("invalid credentials")
	}

	// 3. 生成token，使用slug而不是ID
	tokens, err := helper.GenerateTokenPair(user.GetSlugString(), rememberMe)
	if err != nil {
		return nil, nil, err
	}

	user.PasswordHash = nil // 清除敏感信息
	return user, tokens, nil
}

func (s *AuthService) RefreshTokens(refreshToken string) (*helper.TokenPair, error) {
	// 1. 验证刷新令牌
	claims, err := helper.ParseRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// 2. 检查用户是否存在，通过slug查询
	user, err := s.repo.GetBySlug(claims.Slug)
	if err != nil {
		return nil, err
	}

	// 3. 生成新的令牌对，继续使用slug
	return helper.GenerateTokenPair(user.GetSlugString(), false)
}

func (s *AuthService) Logout(token string) error {
	claims, err := helper.ParseAccessToken(token)
	if err != nil {
		return err
	}
	slug,err := uuid.Parse(claims.Slug)
	if err != nil {
		return err
	}
	blacklist := &model.TokenBlacklist{
		Slug:      slug,
		Token:     token,
		ExpiresAt: claims.ExpiresAt.Time,
	}

	return s.blacklistRepo.Create(blacklist)
}
