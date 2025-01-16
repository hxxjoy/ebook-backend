package service

import (
	"book-fiber/home/repository"
	"book-fiber/model"
	"strconv"
)

type BookCategoryService struct {
	repo *repository.BookCategoryRepository
}

func NewBookCategoryService(repo *repository.BookCategoryRepository) *BookCategoryService {
	return &BookCategoryService{repo: repo}
}

func (s *BookCategoryService) GetList(query *model.BookQuery) (map[string]any, error) {
	return s.repo.GetList(query)
}

func (s *BookCategoryService) BuildCategoryTree() ([]model.CategoryResponseSimple, error) {
	return s.repo.BuildCategoryTree()
}

func (s *BookCategoryService) GetBooksByCategory(slug string, page int, pageSize int) (map[string]any, error) {
	// 获取子分类和 ID 字符串
	/*subCategories, err := s.repo.GetSubCategoriesBySlug(slug)
		if err != nil {
			return nil, err
		}
		var ids strings.Builder
	    isFirst := true

	    for _, category := range subCategories {
	        if !isFirst {
	            ids.WriteString(",")
	        } else {
	            isFirst = false
	        }
	        ids.WriteString(strconv.FormatUint(uint64(category.ID), 10))
	    }

	    // 如果没有找到分类，返回空结果
	    if ids.Len() == 0 {
	        return map[string]any{}, nil
	    }*/
	category, err := s.repo.GetOneBySlug(slug)
	if err != nil {
		return nil, err
	}
	// 获取图书列表
	return s.repo.GetBooksByCategories(strconv.Itoa(category.ID), page, pageSize)
}

/*
func (s *BookCategoryService) GetOne(slug string) (map[string]any, error) {
	book, err := s.repo.GetBookBySlug(slug)
	if err != nil {
		return nil, errors.New("book not found")
	}
	id, ok := book["id"].(int32)
	if !ok {
		return nil, errors.New("invalid book id type")
	}
	categories, err := s.repo.GetBookCategory(id)
	if err != nil {
		return nil, err
	}
	chapters, err := s.repo.GetChapters(id)
	if err != nil {
		return nil, errors.New("chapters not found")
	}
	response := map[string]any{"book": book, "categories": categories, "chapters": chapters}
	return response, nil
}
*/
