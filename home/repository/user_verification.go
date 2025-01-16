package repository

import (
	"book-fiber/model"
	"book-fiber/system/cache"
	"book-fiber/system/config"
	"errors"

	"gorm.io/gorm"
)

type UserVerificationRepository struct {
	db     *gorm.DB
	cache  *cache.Cache
	config *config.Config
}

func NewUserVerificationRepository(db *gorm.DB, cache *cache.Cache, config *config.Config) *UserVerificationRepository {
	return &UserVerificationRepository{db: db, cache: cache, config: config}
}

func (r *UserVerificationRepository) GetOneByEmail(email string) (*model.UserVerification, error) {
	var user model.UserVerification
	result := r.db.First(&user, "email = ?", email)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            return nil, nil  // 未找到记录返回 nil, nil
        }
		return &model.UserVerification{}, result.Error
	}
	return &user, nil
}

func (r *UserVerificationRepository) Create(user *model.UserVerification) error {
	return r.db.Create(user).Error
}
