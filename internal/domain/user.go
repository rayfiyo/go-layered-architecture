package domain

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidUserName  = errors.New("invalid user name")
	ErrInvalidUserEmail = errors.New("invalid user email")
)

type User struct {
	ID        int64
	Name      string
	Email     string
	CreatedAt time.Time
}

func NewUser(name, email string, now time.Time) (User, error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)

	if name == "" {
		return User{}, ErrInvalidUserName
	}
	if email == "" || !strings.Contains(email, "@") {
		return User{}, ErrInvalidUserEmail
	}

	return User{
		ID:        0,
		Name:      name,
		Email:     email,
		CreatedAt: now.UTC(),
	}, nil
}
