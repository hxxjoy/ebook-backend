// internal/model/author.go
package model

import (
	"time"

	"github.com/google/uuid"
)

type Author struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Slug      uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name      string    `json:"name"`
	Biography string    `json:"biography"`
	Avatar    string    `json:"avatar"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 关联
	Books []Book `json:"books" gorm:"foreignKey:AuthorID"`
}

func (m *Author) GetSlugString() string {
	return m.Slug.String()
}