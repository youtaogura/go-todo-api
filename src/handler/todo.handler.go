package handler

import (
	"encoding/json"
	"go_todo/src/model"
	"go_todo/src/service"
	"go_todo/src/types"
	request_util "go_todo/src/util"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type TodoHandler struct {
	TodoService service.TodoService
}

func (h *TodoHandler) ListTodos(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(types.UserKey{}).(*model.User)
	todos := h.TodoService.ListTodos(user)
	request_util.ReturnJson(w, request_util.ReturnJsonOptions{
		Content: todos,
	})
}

func (h *TodoHandler) NewTodo(w http.ResponseWriter, r *http.Request) {
	user := request_util.RequestUser(r)
	var body types.NewTodoRequest
	json.NewDecoder(r.Body).Decode(&body)

	if body.Title == "" && body.Description == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	todo := h.TodoService.NewTodo(user, body)
	request_util.ReturnJson(w, request_util.ReturnJsonOptions{
		Content: todo,
		Status: http.StatusCreated,
	})
}

func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(types.UserKey{}).(*model.User)
	var body types.UpdateTodoRequest
	json.NewDecoder(r.Body).Decode(&body)

	todoId, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	todo := h.TodoService.UpdateTodo(user, todoId, body)
	if todo == nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	request_util.ReturnJson(w, request_util.ReturnJsonOptions{
		Content: todo,
	})
}

func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(types.UserKey{}).(*model.User)
	todoId, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ok := h.TodoService.DeleteTodo(user, int(todoId))
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

