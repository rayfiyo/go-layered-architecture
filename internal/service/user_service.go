// サービス層：ビジネスロジックを担当
// ユーザーに関するビジネスロジックの実装

package service

import (
	"github.com/rayfiyo/layered/internal/domain"
	"github.com/rayfiyo/layered/internal/repository"
)

// ユーザーに関するビジネスロジックを提供するためのインターフェース
type UserService interface {
	GetUser(id int) (*domain.User, error)
	CreateUser(user *domain.User) error
}

type userService struct {
	repo repository.UserRepository
}

// UserService のインスタンスを生成
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

// リポジトリからユーザーを取得
func (s *userService) GetUser(id int) (*domain.User, error) {
	return s.repo.GetByID(id)
}

// リポジトリを通じてユーザーを作成
func (s *userService) CreateUser(user *domain.User) error {
	return s.repo.Create(user)
}
