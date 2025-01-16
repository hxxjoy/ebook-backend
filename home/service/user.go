package service

import (
	"book-fiber/home/repository"
	"book-fiber/model"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

/*
func (s *UserService) GetSearchList(word string, page int) (map[string]any, error) {
	return s.repo.GetSearchList(word, page)
}*/

func (s *UserService) GetByEmail(email string) (*model.User, error) {
	one, err := s.repo.GetByEmail(email)
	if err != nil {
		return &model.User{}, errors.New("record not found")
	}
	return one, nil
}

func (s *UserService) GetBySlug(slug string) (*model.User, error) {
	one, err := s.repo.GetBySlug(slug)
	if err != nil {
		return &model.User{}, errors.New("record not found")
	}
	return one, nil
}

func (s *UserService) CheckCredentials(email string, password string) (*model.User, error) {
	one, err := s.repo.GetByEmail(email)
	if err != nil || !(one.ID > 0){
		return &model.User{}, errors.New("incorrect username or password")
	}
	
	err = bcrypt.CompareHashAndPassword([]byte(*one.PasswordHash), []byte(password))
    if err != nil {
        return nil, errors.New("incorrect username or password")
    }
    // 3. 如果验证通过，返回用户信息
    one.PasswordHash = nil // 不返回密码 hash
    return one, nil
}

func (s *UserService) Create(user *model.User) (*model.User, error) {
	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	hash := string(hashedPassword)
	user.PasswordHash = &hash
	// 保存到数据库
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	// 清除密码后返回
	user.PasswordHash = nil
	return user, nil
}
