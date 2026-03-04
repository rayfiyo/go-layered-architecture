package domain

import (
	"errors"
	"strings"
)

// User はドメイン層のエンティティ（ビジネス上の「ユーザー」）を表す。
// ドメイン層は他の層に依存しない。
type User struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var (
	// ErrInvalidUserName はユーザー名が不正な場合のドメインエラー。
	ErrInvalidUserName = errors.New("ユーザー名が不正です")
	// ErrInvalidUserEmail はメールアドレスが不正な場合のドメインエラー。
	ErrInvalidUserEmail = errors.New("メールアドレスが不正です")
)

// NewUser はユーザーの生成を行う（ドメインルールの最小限の検証もここで行う）。
// ここでは簡素に「空文字は不可」「前後の空白は除去」を行う。
func NewUser(name, email string) (*User, error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)

	if name == "" {
		return nil, ErrInvalidUserName
	}
	if email == "" {
		return nil, ErrInvalidUserEmail
	}

	return &User{
		// ID は永続化層（リポジトリ）が採番するため、ここでは 0 のまま。
		ID:    0,
		Name:  name,
		Email: email,
	}, nil
}
