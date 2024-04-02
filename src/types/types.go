package types

import "go_todo/src/model"

type UserKey struct{}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResult struct {
	User *model.User `json:"user"`
	Token string `json:"token"`
}

type ListUsersResponse struct {
	ID int64 `json:"id"`
	Username string `json:"username"`
}

type NewTodoRequest struct {
	Title string `json:"title"`
	Description string `json:"description"`
}

type UpdateTodoRequest struct {
	Title string `json:"title"`
	Description string `json:"description"`
	Completed bool `json:"completed"`
}

