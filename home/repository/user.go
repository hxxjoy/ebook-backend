package repository

import (
	"book-fiber/model"
	"book-fiber/system/cache"
	"book-fiber/system/config"
	"errors"

	"gorm.io/gorm"
)

type UserRepository struct {
	db     *gorm.DB
	cache  *cache.Cache
	config *config.Config
}

func NewUserRepository(db *gorm.DB, cache *cache.Cache, config *config.Config) *UserRepository {
	return &UserRepository{db: db, cache: cache, config: config}
}

func (r *UserRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // 未找到记录返回 nil, nil
		}
		return &model.User{}, result.Error
	}
	return &user, nil
}

func (r *UserRepository) GetBySlug(slug string) (*model.User, error) {
	var user model.User
	result := r.db.Where("slug = ?", slug).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // 未找到记录返回 nil, nil
		}
		return &model.User{}, result.Error
	}
	return &user, nil
}

func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}
