package service

import (
	"book-fiber/home/repository"
	"book-fiber/model"
	"errors"
)

type BookService struct {
	repo *repository.BookRepository
}

func NewBookService(repo *repository.BookRepository) *BookService {
	return &BookService{repo: repo}
}

func (s *BookService) GetList(query *model.BookQuery) (map[string]any, error) {
	return s.repo.GetList(query)
}
/*
func (s *BookService) GetSearchList(word string, page int) (map[string]any, error) {
	return s.repo.GetSearchList(word, page)
}*/

func (s *BookService) GetOne(slug string) (map[string]any, error) {
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

func (s *BookService) GetChapters(slug string) (map[string]any, error) {
	book, err := s.repo.GetBookBySlug(slug)
	if err != nil {
		return nil, errors.New("book not found")
	}
	id, ok := book["id"].(int32)
	if !ok {
		return nil, errors.New("invalid book id type")
	}
	chapters, err := s.repo.GetChapters(id)
	if err != nil {
		return nil, errors.New("chapters not found")
	}
	response := map[string]any{"book": book, "chapters": chapters}
	return response, nil
}

func (s *BookService) GetChapter(slug string) (map[string]any, error) {
	return s.repo.GetChapter(slug)
}
