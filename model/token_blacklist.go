package model

import (
	"time"

	"github.com/google/uuid"
)

type TokenBlacklist struct {
    Slug      uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4()"`
	Token     string    `gorm:"type:text;not null"`
    ExpiresAt time.Time `gorm:"not null;index;comment:'token过期时间'"`
    CreatedAt time.Time `gorm:"not null;autoCreateTime;comment:'创建时间'"`
}

// TableName 设置表名
func (TokenBlacklist) TableName() string {
    return "token_blacklist"
}

func (m *TokenBlacklist) GetSlugString() string {
	return m.Slug.String()
}