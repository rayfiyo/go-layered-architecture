package domain

import "context"

// UserRepository はユーザーの永続化操作を抽象化するインターフェース。
// DDD の考え方では「ドメイン層がインフラ層に依存しない」ため、
// 具体実装（DB、外部API等）はインフラ層に置き、ドメイン層では抽象（interface）のみを定義する。
type UserRepository interface {
	// Save はユーザーを保存する。
	// 採番が必要な場合、保存後に user.ID が設定されることを期待する。
	Save(ctx context.Context, user *User) error

	// FindByID は ID からユーザーを取得する。
	// 見つからない場合は (nil, nil) を返す。
	FindByID(ctx context.Context, id int64) (*User, error)
}
