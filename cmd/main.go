package main

import (
	"log"

	"github.com/rayfiyo/go-layered-architecture/internal/application/service"
	"github.com/rayfiyo/go-layered-architecture/internal/infrastructure/repository"
	"github.com/rayfiyo/go-layered-architecture/internal/presentation/handler"
)

// main はエントリーポイント。
// ここで各層のコンポーネントを組み立て（DI）、ルーティングを設定してサーバーを起動する。
func main() {
	// インフラ層（リポジトリの具体実装）
	userRepo := repository.NewInMemoryUserRepository()

	// アプリケーション層（ユースケース）
	userService := service.NewUserService(userRepo)

	// プレゼンテーション層（HTTP ハンドラー）
	userHandler := handler.NewUserHandler(userService)
	router := handler.NewRouter(userHandler)

	// サーバー起動
	addr := ":8080"
	log.Printf("listening on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
