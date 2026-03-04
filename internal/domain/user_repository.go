package domain

import (
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, u User) (User, error)
	GetByID(ctx context.Context, id int64) (User, error)
}
