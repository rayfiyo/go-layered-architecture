package handler

import "github.com/gin-gonic/gin"

// NewRouter は HTTP ルーティングを構築する。
// プレゼンテーション層の責務として、エンドポイントとハンドラーの紐付けを行う。
func NewRouter(userHandler *UserHandler) *gin.Engine {
	r := gin.Default()

	// ユーザー作成
	r.POST("/users", userHandler.CreateUser)
	// ユーザー取得
	r.GET("/users", userHandler.GetUser)

	return r
}
