package model

import (
	"time"

	"github.com/google/uuid"
)

type BookQuery struct {
	Page          int    `query:"page"  binding:"omitempty,min=1" default:"1"`
	PageSize      int    `query:"page_size" binding:"omitempty,min=1,max=100" default:"20"`
	Title         string `query:"title"`
	Author        string `query:"author"`
	IsRecommended int    `query:"is_recommended"`
	CategoryID    int    `query:"category_id"`
	Language      string `query:"language"`
	OrderBy       string `query:"order_by"`
	Search        string `query:"search"`
}

// BookListResponse 列表响应结构
type BookListResponse struct {
	Total int64                 `json:"total"`
	Items []*BookResponseSimple `json:"items"`
}

type BookResponseSimple struct {
	Slug     uuid.UUID `json:"slug"`
	Title    string    `json:"title"`
	Subtitle *string   `json:"subtitle,omitempty"`
	Cover    *string   `json:"cover,omitempty"`
	Language *string   `json:"language,omitempty"`
	Author   *string   `json:"author,omitempty"`
}

type BookResponse struct {
	Slug           uuid.UUID          `json:"slug"`
	Title          string             `json:"title"`
	Subtitle       *string            `json:"subtitle,omitempty"`
	AuthorID       *int64             `json:"author_id,omitempty"`
	BookCategoryID *int64             `json:"book_category_id,omitempty"`
	Cover          *string            `json:"cover,omitempty"`
	Description    *string            `json:"description,omitempty"`
	Publisher      *string            `json:"publisher,omitempty"`
	PublishedDate  *time.Time         `json:"published_date,omitempty"`
	Language       *string            `json:"language,omitempty"`
	PageCount      *int32             `json:"page_count,omitempty"`
	Rating         *float64           `json:"rating,omitempty"`
	Tags           []string           `json:"tags,omitempty"`
	IsFree         *bool              `json:"is_free,omitempty"`
	DownloadCount  *int32             `json:"download_count,omitempty"`
	ViewCount      *int32             `json:"view_count,omitempty"`
	IsRecommended  *bool              `json:"is_recommended,omitempty"`
	OtherBookName  *string            `json:"other_book_name,omitempty"`
	LowerBookName  *string            `json:"lower_book_name,omitempty"`
	Author         *string            `json:"author,omitempty"`
	ChapterCount   *int16             `json:"chapter_count,omitempty"`
	Categories     []CategoryResponse `json:"categories"`
}

type Book struct {
	ID              int64              `json:"id"`
	Slug         uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4()"`
	Title           string             `json:"title"`
	Subtitle        *string            `json:"subtitle,omitempty"`
	AuthorID        *int64             `json:"author_id,omitempty"`
	BookCategoryID  *int64             `json:"book_category_id,omitempty"`
	Cover           *string            `json:"cover,omitempty"`
	Description     *string            `json:"description,omitempty"`
	Publisher       *string            `json:"publisher,omitempty"`
	PublishedDate   *time.Time         `json:"published_date,omitempty"`
	Language        *string            `json:"language,omitempty"`
	PageCount       *int32             `json:"page_count,omitempty"`
	FileSize        *int32             `json:"file_size,omitempty"`
	FileFormat      *string            `json:"file_format,omitempty"`
	FilePath        *string            `json:"file_path,omitempty"`
	ReadingProgress *int32             `json:"reading_progress,omitempty"`
	Rating          *float64           `json:"rating,omitempty"`
	Tags            []string           `json:"tags,omitempty"`
	IsFree          *bool              `json:"is_free,omitempty"`
	Price           *float64           `json:"price,omitempty"`
	DiscountPrice   *float64           `json:"discount_price,omitempty"`
	DownloadCount   *int32             `json:"download_count,omitempty"`
	ViewCount       *int32             `json:"view_count,omitempty"`
	Status          *int16             `json:"status,omitempty"`
	CreatedAt       *time.Time         `json:"created_at,omitempty"`
	UpdatedAt       *time.Time         `json:"updated_at,omitempty"`
	IsRecommended   *bool              `json:"is_recommended,omitempty"`
	OtherBookName   *string            `json:"other_book_name,omitempty"`
	LowerBookName   *string            `json:"lower_book_name,omitempty"`
	FromURL         *string            `json:"from_url,omitempty"`
	FromChapterURL  *string            `json:"from_chapter_url,omitempty"`
	OriCover        *string            `json:"ori_cover,omitempty"`
	Author          *string            `json:"author,omitempty"`
	IsCategoried    *int16             `json:"is_categoried,omitempty"`
	ChapterCount    *int16             `json:"chapter_count,omitempty"`
	Categories      []CategoryResponse `json:"categories" gorm:"-"`
	BookCategories  []BookCategory     `json:"-" gorm:"many2many:ln_book_category;"`
}

// TableName 指定表名
func (Book) TableName() string {
	return "book" // 返回实际的表名
}

func (m *Book) GetSlugString() string {
	return m.Slug.String()
}