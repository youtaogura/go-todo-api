package handler

import (
	"encoding/json"
	"go_todo/src/model"
	"go_todo/src/service"
	"go_todo/src/types"
	"go_todo/src/util"
	"net/http"

	"gorm.io/gorm"
)

type AuthHandler struct {
	DB *gorm.DB
	AuthService service.AuthService
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var body types.LoginRequest
	json.NewDecoder(r.Body).Decode(&body)

	res := h.AuthService.Login(body)
	if res == nil {
		http.Error(w, "Invalid username or password", http.StatusNotFound)
		return
	}

	setAuthHeader(w, res.Token)
	util.ReturnJson(w, util.ReturnJsonOptions{
		Content: map[string]interface{}{
			"id": res.User.ID,
			"username": res.User.Username,
		},
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(types.UserKey{}).(*model.User)
	cookie, err := r.Cookie("access_token")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	logoutChan := make(chan bool)
	var ok bool
	util.Wait(
		func() { h.AuthService.Logout(user, cookie.Value, logoutChan) },
		func() { deleteAuthHeader(w) },
		func() { ok = <- logoutChan },
	)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	util.ReturnJson(w, util.ReturnJsonOptions{})
}

func setAuthHeader(w http.ResponseWriter, token string) {
	w.Header().Set("Authorization", "Bearer " + token)
}

func deleteAuthHeader(w http.ResponseWriter) {
	w.Header().Del("Authorization")
}