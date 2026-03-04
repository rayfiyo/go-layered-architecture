package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/rayfiyo/go-layered-architecture/internal/application/service"
)

type UserHandler struct {
	userSvc *service.UserService
}

func NewUserHandler(userSvc *service.UserService) *UserHandler {
	return &UserHandler{userSvc: userSvc}
}

type createUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func (h *UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	out, err := h.userSvc.CreateUser(r.Context(), service.CreateUserInput{
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, out)
}

func (h *UserHandler) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeJSON(w, http.StatusBadRequest,
			errorResponse{Error: "missing query parameter: id"})
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid id"})
		return
	}

	out, err := h.userSvc.GetUserByID(r.Context(), id)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, out)
}

func writeServiceError(w http.ResponseWriter, err error) {
	// アプリケーション層が wrap した sentinel を文字列で判定せず errors.Is で分類する
	switch {
	case errors.Is(err, service.ErrBadRequest):
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
	case errors.Is(err, service.ErrNotFound):
		writeJSON(w, http.StatusNotFound, errorResponse{Error: err.Error()})
	case errors.Is(err, service.ErrConflict):
		writeJSON(w, http.StatusConflict, errorResponse{Error: err.Error()})
	default:
		writeJSON(w, http.StatusInternalServerError,
			errorResponse{Error: "internal server error"})
	}
}

func decodeJSON(r *http.Request, dst any) error {
	if r.Body == nil {
		return errors.New("empty body")
	}
	defer func() {
		_ = r.Body.Close()
	}()

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		return err
	}
	// 余分な JSON が続いていないかチェック
	if dec.More() {
		return errors.New("invalid json: multiple values")
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(true)

	_ = enc.Encode(v)
}
