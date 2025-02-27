// ハンドラー層：HTTP リクエストの受付とレスポンス処理を担当
// ユーザーに関する API エンドポイントの実装

package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/rayfiyo/layered/internal/domain"
	"github.com/rayfiyo/layered/internal/service"
)

// HTTP リクエストを処理するためのハンドラー
type UserHandler struct {
	service service.UserService
}

// UserHandler のインスタンスを生成
func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// /users エンドポイントへのリクエストを処理
// GET: クエリパラメータ "id" によるユーザー取得
// POST: リクエストボディからユーザー作成
func (h *UserHandler) HandleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getUser(w, r)
	case http.MethodPost:
		h.createUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GET メソッドでユーザーを取得
func (h *UserHandler) getUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
		return
	}
	user, err := h.service.GetUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// POST メソッドで新しいユーザーを作成
func (h *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var user domain.User
	if err := json.Unmarshal(body, &user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if err := h.service.CreateUser(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
