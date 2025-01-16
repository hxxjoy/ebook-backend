package repository

import (
	"book-fiber/model"
	"book-fiber/system/config"
	"book-fiber/system/cache"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"
	"gorm.io/gorm"
)

type BookCategoryRepository struct {
	db     *gorm.DB
	cache  *cache.Cache
	config *config.Config
}

func NewBookCategoryRepository(db *gorm.DB, cache *cache.Cache, config *config.Config) *BookCategoryRepository {
	return &BookCategoryRepository{db: db, cache: cache, config: config}
}

// bookCategories 返回book_categories表的查询构建器
func (r *BookCategoryRepository) bookCategory() *gorm.DB {
	return r.db.Table("book_category bc")
}

func (r *BookCategoryRepository) book() *gorm.DB {
	return r.db.Table("book bc")
}

func (r *BookCategoryRepository) GetOneBySlug(slug string) (model.BookCategorySimple, error) {
	var one model.BookCategorySimple
	err := r.bookCategory().Where("status=1 AND slug=?",slug).Limit(1).First(&one).Error
	return one, err
}

func (r *BookCategoryRepository) GetAllCategories() ([]model.BookCategorySimple, error) {
	var categories []model.BookCategorySimple
	err := r.bookCategory().Where("status=1").Order("sort DESC,id ASC").Find(&categories).Error
	return categories, err
}

// 转换数据结构的辅助函数
func (r *BookCategoryRepository) BuildCategoryTree() ([]model.CategoryResponseSimple, error) {
	cacheKey := "categories"
	cacheData, found := r.cache.Get(cacheKey)
	if found {
		if result, ok := cacheData.([]model.CategoryResponseSimple); ok {
			return result, nil
		}
	}
	// 获取所有分类
	categories, err := r.GetAllCategories()
	if err != nil {
		return nil, err
	}

	// 创建一个map来存储父分类
	categoryMap := make(map[int][]model.BookCategorySimple)

	// 将分类按照parent_id分组
	for _, category := range categories {
		categoryMap[category.ParentID] = append(categoryMap[category.ParentID], category)
	}

	// 构建顶级分类响应
	var result []model.CategoryResponseSimple
	rootCategories := categoryMap[0] // 获取顶级分类（parent_id = 0）

	// 递归构建分类树
	for _, rootCategory := range rootCategories {
		categoryResponse := model.CategoryResponseSimple{
			Title:         rootCategory.Title,
			Slug:          rootCategory.Slug,
			SubCategories: buildSubCategories(categoryMap, rootCategory.ID),
		}
		result = append(result, categoryResponse)
	}
	r.cache.Set("categories", result, 0)
	return result, nil
}

// 递归构建子分类
func buildSubCategories(categoryMap map[int][]model.BookCategorySimple, parentID int) []model.CategoryResponseSimple {
	if subCategories, exists := categoryMap[parentID]; exists {
		var result []model.CategoryResponseSimple
		for _, subCategory := range subCategories {
			categoryResponse := model.CategoryResponseSimple{
				Title:         subCategory.Title,
				Slug:          subCategory.Slug,
				SubCategories: buildSubCategories(categoryMap, subCategory.ID),
			}
			result = append(result, categoryResponse)
		}
		return result
	}
	return nil
}

func (r *BookCategoryRepository) GetSubCategoriesBySlug(slug string) ([]model.BookCategorySimple, error) {
	cacheKey := fmt.Sprintf("subcagegories:slug:%s", slug)
	cacheData, found := r.cache.Get(cacheKey)
	if found {
		if result, ok := cacheData.([]model.BookCategorySimple); ok {
			return result, nil
		}
	}
	var categories []model.BookCategorySimple
	query := `
        SELECT id, title, slug, parent_id
        FROM book_category 
        WHERE parent_id = (SELECT id FROM book_category WHERE slug = ?)
    `

	if err := r.db.Raw(query, slug, slug).Scan(&categories).Error; err != nil {
		return nil, err
	}
	r.cache.Set(cacheKey,categories,0)
	return categories, nil
}

func (r *BookCategoryRepository) GetBooksByCategories(categoryIDs string, page, pageSize int) (map[string]any, error) {
	cacheKey := fmt.Sprintf("books:category:%s:page:%d:size:%d", categoryIDs, page, pageSize)
	cacheData, found := r.cache.Get(cacheKey)
	if found {
		if result, ok := cacheData.(map[string]any); ok {
			return result, nil
		}
	}
	offset := (page - 1) * pageSize
	order := "ORDER BY created_at ASC"
	limit := fmt.Sprintf(" LIMIT %d OFFSET %d", pageSize, offset)

	sql := fmt.Sprintf(`
        SELECT COUNT(*) OVER() as total_count, slug, title, cover, author, description
        FROM book b INNER JOIN ln_book_category lbc ON b.id=lbc.book_id
        WHERE b.status=1 AND lbc.book_category_id IN (%s)
        %s %s`, categoryIDs, order, limit)
	var results []map[string]any
	err := r.db.Raw(sql).Scan(&results).Error

	// 构造响应
	response := map[string]any{
		"items":     results,
		"page":      page,
		"page_size": pageSize,
	}

	// 如果结果集不为空，获取总数
	if len(results) > 0 {
		if totalCount, ok := results[0]["total_count"].(int64); ok {
			response["total_count"] = totalCount
			response["total_pages"] = int(math.Ceil(float64(totalCount) / float64(pageSize)))
		}
	}

	r.cache.Set(cacheKey, response, 0)

	return response, err
}

func (r *BookCategoryRepository) GetList(query *model.BookQuery) (map[string]any, error) {
	/* Part of the code is hidden due to copyright restrictions. */

	if r.config.Cache.Enabled {
		r.cache.Set(cacheKey, response, 5*time.Minute)
	}

	return response, err
}