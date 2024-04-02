package service

import (
	"go_todo/src/model"
	"go_todo/src/types"

	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
	AuthService AuthService
}

func (h *UserService) ListUsers() []types.ListUsersResponse {
 	var users []model.User
  h.DB.Select("id", "username").Find(&users)
	result := make([]types.ListUsersResponse, len(users))

	for i := range result {
		result[i].ID = users[i].ID
		result[i].Username = users[i].Username
	}
	return result
}

func (h *UserService) RegisterUser(body types.LoginRequest) *model.User {
	var duplicated *model.User
	h.DB.Where("username = ?", body.Username).Find(&duplicated)

	if duplicated.ID != 0 {
		return nil
	}

	var user model.User
	user.Username = body.Username
	user.Password = h.AuthService.Encrypt(body.Password)
	h.DB.Create(&user)

	return &user
}
