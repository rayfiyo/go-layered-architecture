package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rayfiyo/go-layered-architecture/internal/domain"
	"github.com/rayfiyo/go-layered-architecture/internal/infrastructure/repository"
)

var (
	ErrBadRequest = errors.New("bad request")
	ErrNotFound   = errors.New("not found")
	ErrConflict   = errors.New("conflict")
)

type UserService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) *UserService {
	return &UserService{repo: repo}
}

type CreateUserInput struct {
	Name  string
	Email string
}

type UserOutput struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *UserService) CreateUser(
	ctx context.Context, in CreateUserInput,
) (UserOutput, error) {
	u, err := domain.NewUser(in.Name, in.Email, time.Now())
	if err != nil {
		if errors.Is(err, domain.ErrInvalidUserName) ||
			errors.Is(err, domain.ErrInvalidUserEmail) {
			return UserOutput{}, fmt.Errorf("%w: %v", ErrBadRequest, err)
		}
		return UserOutput{}, err
	}

	created, err := s.repo.Create(ctx, u)
	if err != nil {
		if errors.Is(err, repository.ErrUserEmailConflict) {
			return UserOutput{}, fmt.Errorf("%w: %v", ErrConflict, err)
		}
		return UserOutput{}, err
	}

	return UserOutput{
		ID:        created.ID,
		Name:      created.Name,
		Email:     created.Email,
		CreatedAt: created.CreatedAt,
	}, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id int64) (UserOutput, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return UserOutput{}, fmt.Errorf("%w: %v", ErrNotFound, err)
		}
		return UserOutput{}, err
	}

	return UserOutput{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}, nil
}
