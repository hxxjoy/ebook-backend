package model

import (
	"time"

	"github.com/google/uuid"
)


type UserVerification struct {
	ID              int64              `json:"id"`
	Slug            uuid.UUID          `gorm:"type:uuid;default:uuid_generate_v4()"`
	Username        *string    `db:"username" json:"username"`
	PasswordHash    *string    `db:"password_hash" json:"password_hash"`
	Nickname        *string    `db:"nickname" json:"nickname"`
	Email           *string    `db:"email" json:"email"`
	Avatar          *string    `db:"avatar" json:"avatart"` 
	Age             *int16     `db:"age" json:"age"`
	Role            *string    `db:"role" json:"role"`
	Status          *int16             `json:"status,omitempty"`
	CreatedAt       *time.Time         `json:"created_at,omitempty"`
	UpdatedAt       *time.Time         `json:"updated_at,omitempty"`
}

// TableName 指定表名
func (UserVerification) TableName() string {
	return "user_verification" // 返回实际的表名
}
