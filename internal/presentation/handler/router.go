package handler

import (
	"net/http"

	"github.com/rayfiyo/go-layered-architecture/internal/application/service"
)

func NewRouter(userSvc *service.UserService) http.Handler {
	mux := http.NewServeMux()

	uh := NewUserHandler(userSvc)
	mux.HandleFunc("POST /users", uh.HandleCreateUser)
	mux.HandleFunc("GET /users", uh.HandleGetUser)

	// ヘルスチェック（任意だが省略しない）
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok\n"))
	})

	return mux
}
