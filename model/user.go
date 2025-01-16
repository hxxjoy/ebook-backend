package model

import (
	"time"

	"github.com/google/uuid"
)

type UserQuery struct {
	ID       int64     `json:"id"`
	Slug     uuid.UUID `json:"slug"`
	Username *string   `db:"username" json:"username"`
	Email    *string   `db:"email" json:"email"`
}

type User struct {
	ID           int64      `json:"id"`
	Slug         uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4()"`
	Username     *string    `db:"username" json:"username"`
	PasswordHash *string    `db:"password_hash" json:"password_hash"`
	Nickname     *string    `db:"nickname" json:"nickname"`
	Email        *string    `db:"email" json:"email"`
	Avatar       *string    `db:"avatar" json:"avatart"`
	Age          *int16     `db:"age" json:"age"`
	Role         *string    `db:"role" json:"role"`
	Status       *int16     `json:"status,omitempty"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
}

// TableName 指定表名
func (User) TableName() string {
	return "user" // 返回实际的表名
}

func (u *User) GetSlugString() string {
	return u.Slug.String()
}
