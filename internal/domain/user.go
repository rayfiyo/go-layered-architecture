// ドメイン層：ビジネスモデル・エンティティを定義
// ユーザーのエンティティを定義

package domain

// ユーザーのエンティティを表現
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
