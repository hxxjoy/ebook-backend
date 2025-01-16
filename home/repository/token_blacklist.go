package repository

import (
	"book-fiber/model"

	"gorm.io/gorm"
)

type TokenBlacklistRepository interface {
    Create(blacklist *model.TokenBlacklist) error
    IsBlacklisted(slug string) bool
    DeleteExpired() error
}

type tokenBlacklistRepo struct {
    db *gorm.DB
}

func NewTokenBlacklistRepository(db *gorm.DB) TokenBlacklistRepository {
    return &tokenBlacklistRepo{
        db: db,
    }
}

func (r *tokenBlacklistRepo) Create(blacklist *model.TokenBlacklist) error {
    return r.db.Create(blacklist).Error
}

func (r *tokenBlacklistRepo) IsBlacklisted(slug string) bool {
    var blacklist *model.TokenBlacklist
    result := r.db.Where("slug = ?", slug).First(&blacklist)
    return result.Error == nil
}

func (r *tokenBlacklistRepo) DeleteExpired() error {
    return r.db.Where("expires_at < CURRENT_TIMESTAMP").Delete(&model.TokenBlacklist{}).Error
}