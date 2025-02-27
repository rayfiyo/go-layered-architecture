// アプリケーションのエントリーポイント。各レイヤーを組み合わせてサーバーを起動する

package main

import (
	"log"
	"net/http"

	"github.com/rayfiyo/layered/internal/handler"
	"github.com/rayfiyo/layered/internal/repository"
	"github.com/rayfiyo/layered/internal/service"
)

func main() {
	// リポジトリ層：ユーザーのデータ操作を担当するリポジトリの初期化
	userRepo := repository.NewInMemoryUserRepository()

	// サービス層：リポジトリを利用してビジネスロジックを提供するサービスの初期化
	userService := service.NewUserService(userRepo)

	// ハンドラー層：サービス層を利用して HTTP リクエストを処理するハンドラーの初期化
	userHandler := handler.NewUserHandler(userService)

	// HTTP サーバーのルーティング設定
	mux := http.NewServeMux()
	mux.HandleFunc("/users", userHandler.HandleUsers)

	log.Println("Server is running on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
