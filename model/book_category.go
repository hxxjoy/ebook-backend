package model

import (
	"time"

	"github.com/google/uuid"
)

type CategoryResponse struct {
	Slug         uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4()"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
}

// BookCategory 图书分类表
type BookCategory struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Description *string    `json:"description,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	Slug         uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4()"`
	BookCount   *int32     `json:"book_count,omitempty"`
	Status      *int16     `json:"status,omitempty"`
	Sort        *int16     `json:"sort,omitempty"`
	ParentID    *int64     `json:"parent_id,omitempty"`

	// 关联字段
	/*Parent   *BookCategory   `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children []*BookCategory `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	Books    []*Book         `json:"books,omitempty" gorm:"many2many:ln_book_category;"`*/
}

// LNBookCategory 学习笔记-图书分类关联表
type LnBookCategory struct {
	ID             int64 `json:"id"`
	BookID         int64 `json:"book_id"`
	BookCategoryID int64 `json:"book_category_id"`

	// 关联字段
	Book         *Book         `json:"book,omitempty" gorm:"foreignKey:BookID"`
	BookCategory *BookCategory `json:"book_category,omitempty" gorm:"foreignKey:BookCategoryID"`
}

type BookCategorySimple struct {
	ID       int    `json:"id" gorm:"primary_key"`
	Title    string `json:"title"`
	Slug         uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4()"`
	ParentID int    `json:"parent_id"`
	Sort     int    `json:"sort"`
}

type CategoryResponseSimple struct {
	Title         string                   `json:"title"`
	Slug         uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4()"`
	SubCategories []CategoryResponseSimple `json:"sub_categories"`
}
