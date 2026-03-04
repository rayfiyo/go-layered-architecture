package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rayfiyo/go-layered-architecture/internal/application/service"
	"github.com/rayfiyo/go-layered-architecture/internal/domain"
)

// UserHandler は HTTP の入出力（プレゼンテーション層）を担当する。
// gin.Context を受け取り、アプリケーション層（UserService）を呼び出してレスポンスを返す。
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler は UserHandler を生成する。
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

type createUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// CreateUser はユーザー作成 API（POST /users）。
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "リクエストボディの形式が不正です"})
		return
	}

	user, err := h.userService.CreateUser(c.Request.Context(), req.Name, req.Email)
	if err != nil {
		// ドメインの検証エラーは 400 とする。
		switch err {
		case domain.ErrInvalidUserName, domain.ErrInvalidUserEmail:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError,
				gin.H{"error": "サーバー内部でエラーが発生しました"})
			return
		}
	}

	c.JSON(http.StatusCreated, user)
}

// GetUser はユーザー取得 API（GET /users?id=1）。
func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "クエリパラメータ id は必須です"})
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id は正の整数で指定してください"})
		return
	}

	user, err := h.userService.GetUser(c.Request.Context(), id)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError,
				gin.H{"error": "サーバー内部でエラーが発生しました"})
			return
		}
	}

	c.JSON(http.StatusOK, user)
}
