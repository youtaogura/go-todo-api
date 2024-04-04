package handler

import (
	"encoding/json"
	"go_todo/src/service"
	"go_todo/src/types"
	"go_todo/src/util"
	"net/http"
)

type UserHandler struct {
	UserService service.UserService
	AuthService service.AuthService
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
 	users := h.UserService.ListUsers()
	util.ReturnJson(w, util.ReturnJsonOptions{
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

	res := h.AuthService.Login(body)
	w.Header().Set("Authorization", "Bearer " + res.Token)

	util.ReturnJson(w, util.ReturnJsonOptions{
		Status: http.StatusCreated,
	})
}
