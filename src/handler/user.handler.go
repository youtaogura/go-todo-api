package handler

import (
	"encoding/json"
	"go_todo/src/service"
	"go_todo/src/types"
	request_util "go_todo/src/util"
	"net/http"
)

type UserHandler struct {
	UserService service.UserService
	AuthService service.AuthService
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
 	users := h.UserService.ListUsers()
	request_util.ReturnJson(w, request_util.ReturnJsonOptions{
		Content: users,
	})
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var body types.LoginRequest
	json.NewDecoder(r.Body).Decode(&body)
	user := h.UserService.RegisterUser(body)

	if user == nil {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	request_util.ReturnJson(w, request_util.ReturnJsonOptions{
		Status: http.StatusCreated,
	})
}
