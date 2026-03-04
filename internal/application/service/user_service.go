package service

import (
	"context"
	"errors"

	"github.com/rayfiyo/go-layered-architecture/internal/domain"
)

// ErrUserNotFound はユーザーが存在しない場合のアプリケーション層エラー。
var ErrUserNotFound = errors.New("ユーザーが見つかりません")

// UserService はユースケース（アプリケーション層）を提供する。
// ここでは「ユーザー作成」「ユーザー取得」の2つを実装する。
// アプリケーション層はドメイン層（エンティティ、リポジトリ抽象）に依存する。
type UserService struct {
	repo domain.UserRepository
}

// NewUserService は UserService を生成する。
func NewUserService(repo domain.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// CreateUser はユーザーを作成するユースケース。
func (s *UserService) CreateUser(
	ctx context.Context, name, email string,
) (*domain.User, error) {
	user, err := domain.NewUser(name, email)
	if err != nil {
		// ドメインの検証エラーはそのまま返す。
		return nil, err
	}

	if err := s.repo.Save(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

// GetUser はユーザーを取得するユースケース。
func (s *UserService) GetUser(ctx context.Context, id int64) (*domain.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}
