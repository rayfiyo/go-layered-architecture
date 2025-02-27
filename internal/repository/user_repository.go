// リポジトリ層：データアクセスを担当
// ユーザーのデータ保存／取得の実装（ここでは InMemory の例）

package repository

import (
	"errors"
	"sync"

	"github.com/rayfiyo/layered/internal/domain"
)

/*
// こんな感じで DB の設定を書く仕様の場合もあるだろう

// SQLite を利用したユーザーデータの実装
type SQLiteUserRepository struct {
	db *sql.DB
}

// DB ファイルで SQLiteUserRepository を初期化
func NewSQLiteUserRepository(dbFile string) (*SQLiteUserRepository, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	// テーブルが存在しない場合は作成
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE
	);
	`
	_, err = db.Exec(query)
	if err != nil {
		return nil, err
	}

	return &SQLiteUserRepository{db: db}, nil
}
*/

// ユーザーデータの取得・保存のためのインターフェース
type UserRepository interface {
	GetByID(id int) (*domain.User, error)
	Create(user *domain.User) error
}

// メモリ上にユーザーデータを保持するリポジトリの実装
type InMemoryUserRepository struct {
	mu     sync.RWMutex
	users  map[int]*domain.User
	nextID int
}

// InMemoryUserRepository のインスタンスを返す
func NewInMemoryUserRepository() UserRepository {
	return &InMemoryUserRepository{
		users:  make(map[int]*domain.User),
		nextID: 1,
	}
}

// 指定した ID のユーザーを取得
func (r *InMemoryUserRepository) GetByID(id int) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, exists := r.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// 新しいユーザーを作成
func (r *InMemoryUserRepository) Create(user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	user.ID = r.nextID
	r.users[r.nextID] = user
	r.nextID++
	return nil
}
