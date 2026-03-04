package repository

import (
	"context"
	"sync"

	"github.com/rayfiyo/go-layered-architecture/internal/domain"
)

// InMemoryUserRepository は UserRepository のインメモリ実装。
// 本番では RDB/NoSQL 等に置き換えるが、ここでは簡素な map + mutex で実装する。
// 依存性逆転により、このインフラ層がドメイン層（interface）に依存する。
type InMemoryUserRepository struct {
	mu     sync.RWMutex
	nextID int64
	store  map[int64]*domain.User
}

// NewInMemoryUserRepository はインメモリリポジトリを生成する。
func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		nextID: 1,
		store:  make(map[int64]*domain.User),
	}
}

// Save はユーザーを保存する。
// user.ID が 0 の場合は新規作成として採番し、0 でない場合は上書き更新する。
func (r *InMemoryUserRepository) Save(ctx context.Context, user *domain.User) error {
	_ = ctx // インメモリ実装なので未使用だが、インターフェース統一のため受け取る。

	r.mu.Lock()
	defer r.mu.Unlock()

	if user.ID == 0 {
		user.ID = r.nextID
		r.nextID++
	}

	// 外部からの参照を安全にするため、保存時にコピーを作る。
	copied := *user
	r.store[user.ID] = &copied
	return nil
}

// FindByID は ID からユーザーを取得する。
func (r *InMemoryUserRepository) FindByID(
	ctx context.Context, id int64,
) (*domain.User, error) {
	_ = ctx

	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.store[id]
	if !ok {
		return nil, nil
	}

	// 呼び出し側が変更しても store が汚染されないようにコピーを返す。
	copied := *user
	return &copied, nil
}
