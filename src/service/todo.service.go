package service

import (
	"go_todo/src/model"
	"go_todo/src/types"

	"gorm.io/gorm"
)

type TodoService struct {
	DB *gorm.DB
}

func (h *TodoService) ListTodos(user *model.User) []model.Todo {
	var todos []model.Todo
	h.DB.Where("user_id = ?", user.ID).Find(&todos)
	return todos
}

func (h *TodoService) NewTodo(user *model.User, body types.NewTodoRequest) *model.Todo {
	todo := model.Todo{
		UserID: user.ID,
		Title: body.Title,
		Description: body.Description,
	}
	if todo.Title == "" {
		todo.Title = "Untitled"
	}
	h.DB.Create(&todo)

	return &todo
}

func (h *TodoService) UpdateTodo(user *model.User, todoId int64, body types.UpdateTodoRequest) *model.Todo {
	todo := model.Todo{}
	h.DB.Where("id = ?", todoId).First(&todo)

	if todo.ID == 0 {
		return nil
	}

	if todo.UserID != user.ID {
		return nil
	}

	if body.Title != "" {
		todo.Title = body.Title
	}
	if body.Description != "" {
		todo.Description = body.Description
	}
	todo.Completed = body.Completed
	h.DB.Save(&todo)

	return &todo
}

func (h *TodoService) DeleteTodo(user *model.User, todoId int) bool {
	todo := model.Todo{}
	h.DB.Where("id = ?", todoId).First(&todo)

	if todo.ID == 0 {
		return false
	}

	if todo.UserID != user.ID {
		return false
	}
	h.DB.Delete(&todo)

	return true
}

