package repository

import (
	"book-fiber/model"
	"book-fiber/system/cache"
	"book-fiber/system/config"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type BookRepository struct {
	db     *gorm.DB
	cache  *cache.Cache
	config *config.Config
}

func NewBookRepository(db *gorm.DB, cache *cache.Cache, config *config.Config) *BookRepository {
	return &BookRepository{db: db, cache: cache, config: config}
}

// books 返回books表的查询构建器
func (r *BookRepository) book() *gorm.DB {
	return r.db.Table("book")
}

// bookCategories 返回book_categories表的查询构建器
func (r *BookRepository) bookCategory() *gorm.DB {
	return r.db.Table("book_category bc")
}

func (r *BookRepository) bookChapter() *gorm.DB {
	return r.db.Table("book_chapter bch")
}

func (r *BookRepository) GetBookBySlug(slug string) (*model.Book, error) {
	cacheKey := fmt.Sprintf("book:%s", slug)
	cacheData, found := r.cache.Get(cacheKey)
	if found {
		if result, ok := cacheData.(*model.Book); ok {
			return result, nil
		}
	}
	/* Part of the code is hidden due to copyright restrictions. */
	return data, nil
}

func (r *BookRepository) GetBookCategory(book_id int32) ([]*model.Book, error) {
	bookID := strconv.Itoa(int(book_id))
	cacheKey := fmt.Sprintf("book:cagegories:%s", bookID)
	cacheData, found := r.cache.Get(cacheKey)
	if found {
		if result, ok := cacheData.([]map[string]any); ok {
			return result, nil
		}
	}
	categories := /* Part of the code is hidden due to copyright restrictions. */
	r.cache.Set(cacheKey, categories, 0)
	return categories, nil
}

func (r *BookRepository) GetChapters(book_id int32) ([]*model.Book, error) {
	bookID := strconv.Itoa(int(book_id))
	cacheKey := fmt.Sprintf("book:chapters:%s", bookID)
	cacheData, found := r.cache.Get(cacheKey)
	if found {
		if result, ok := cacheData.([]*model.Book); ok {
			return result, nil
		}
	}
	data := /* Part of the code is hidden due to copyright restrictions. */
	r.cache.Set(cacheKey, data, 0)
	return data, nil
}

func (r *BookRepository) GetChapter(slug string) (*model.Book, error) {
	cacheKey := fmt.Sprintf("chapter:%s", slug)
	cacheData, found := r.cache.Get(cacheKey)
	if found {
		if result, ok := cacheData.(map[string]any); ok {
			return result, nil
		}
	}
	data := /* Part of the code is hidden due to copyright restrictions. */

	r.cache.Set(cacheKey, data, 0)
	return data, nil
}

func (r *BookRepository) GetList(query *model.BookQuery) (map[string]any, error) {
	cacheKey := r.generateCacheKey(query)
	cacheData, found := r.cache.Get(cacheKey)
	if found {
		if result, ok := cacheData.(map[string]any); ok {
			return result, nil
		}
	}
	var results 
	/* Part of the code is hidden due to copyright restrictions. */

	r.cache.Set(cacheKey, response, 0)

	return response, err
}


func (r *BookRepository) generateCacheKey(query *model.BookQuery) string {
	// 将查询参数序列化为缓存键
	params := map[string]interface{}{
		"title":       query.Title,
		"recommended": query.IsRecommended,
		"category_id": query.CategoryID,
		"author":      query.Author,
		"language":    query.Language,
		"search":      query.Search,
		"page":        query.Page,
		"page_size":   query.PageSize,
		"order_by":    query.OrderBy,
	}

	bytes, _ := json.Marshal(params)
	return fmt.Sprintf("book:list:%x", md5.Sum(bytes))
}
