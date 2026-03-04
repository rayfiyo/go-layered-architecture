package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/rayfiyo/go-layered-architecture/internal/domain"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserEmailConflict = errors.New("email already exists")
)

type UserRepositorySQLX struct {
	db *sqlx.DB
}

func NewUserRepositorySQLX(db *sqlx.DB) *UserRepositorySQLX {
	return &UserRepositorySQLX{db: db}
}

type userRow struct {
	ID        int64  `db:"id"`
	Name      string `db:"name"`
	Email     string `db:"email"`
	CreatedAt string `db:"created_at"`
}

func (r *UserRepositorySQLX) Create(
	ctx context.Context, u domain.User,
) (domain.User, error) {
	const q = `
INSERT INTO users (name, email, created_at)
VALUES (?, ?, ?);
`
	createdAt := u.CreatedAt.UTC().Format(time.RFC3339Nano)
	res, err := r.db.ExecContext(ctx, q, u.Name, u.Email, createdAt)
	if err != nil {
		// SQLite の UNIQUE 制約違反を吸収
		// （driver 依存の詳細コードには踏み込まない）
		if looksLikeUniqueViolation(err) {
			return domain.User{}, ErrUserEmailConflict
		}
		return domain.User{}, fmt.Errorf("insert user: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return domain.User{}, fmt.Errorf("last insert id: %w", err)
	}

	u.ID = id
	return u, nil
}

func (r *UserRepositorySQLX) GetByID(ctx context.Context, id int64) (domain.User, error) {
	const q = `
SELECT id, name, email, created_at
FROM users
WHERE id = ?
LIMIT 1;
`
	var row userRow
	if err := r.db.GetContext(ctx, &row, q, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, ErrUserNotFound
		}
		return domain.User{}, fmt.Errorf("select user: %w", err)
	}

	t, err := time.Parse(time.RFC3339Nano, row.CreatedAt)
	if err != nil {
		return domain.User{}, fmt.Errorf("parse created_at: %w", err)
	}

	return domain.User{
		ID:        row.ID,
		Name:      row.Name,
		Email:     row.Email,
		CreatedAt: t.UTC(),
	}, nil
}

func looksLikeUniqueViolation(err error) bool {
	// modernc.org/sqlite はメッセージに "UNIQUE constraint failed" を含むことが多い
	msg := err.Error()
	return strings.Contains(msg, "UNIQUE constraint failed") ||
		strings.Contains(msg, "constraint failed")
}
