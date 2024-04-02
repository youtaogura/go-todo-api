package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"go_todo/src/model"
	"go_todo/src/types"
	"os"

	"gorm.io/gorm"
)

type AuthService struct {
	DB *gorm.DB
	SessionService SessionService
}

func (h *AuthService) Login(body types.LoginRequest) *types.LoginResult {
	var user model.User
	h.DB.Select("id", "username", "password").Where("username = ?", body.Username).First(&user)
	if user.Password != h.Encrypt(body.Password) {
		return nil
	}
	token := h.SessionService.NewSession(&user)
	return &types.LoginResult{
		User: &user,
		Token: token,
	}
}

func (h *AuthService) Logout(user *model.User, token string) bool {
	ok := h.SessionService.DeleteSession(user, token)
	return ok
}

func (h *AuthService) Encrypt(s string) string {
	var key = []byte(os.Getenv("PASSWORD_SALT"))
	hmac := hmac.New(sha256.New, key)
	hmac.Write([]byte(s))
	hash := hmac.Sum(nil)
	return base64.StdEncoding.EncodeToString(hash)
}
