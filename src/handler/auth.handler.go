package handler

import (
	"encoding/json"
	"go_todo/src/model"
	"go_todo/src/service"
	"go_todo/src/types"
	request_util "go_todo/src/util"
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

	setCookie(w, res.Token)
	request_util.ReturnJson(w, request_util.ReturnJsonOptions{
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

	ok := h.AuthService.Logout(user, cookie.Value)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	deleteCookie(w)
	request_util.ReturnJson(w, request_util.ReturnJsonOptions{})
}

func setCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:  "access_token",
		Value: token,
		Path:  "/",
	})
}

func deleteCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   "access_token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}