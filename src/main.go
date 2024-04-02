package main

import (
	"context"
	"fmt"
	"go_todo/src/database"
	"go_todo/src/handler"
	"go_todo/src/service"
	"go_todo/src/types"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	database.Setup()
	StartWebServer()
}

func StartWebServer() error {
	fmt.Println("Rest API with Mux Routers")
	router := mux.NewRouter().StrictSlash(true)
	router.Use(logMW)
	setupRoutes(router)

	return http.ListenAndServe(fmt.Sprintf(":%d", 8080), router)
}

func setupRoutes(router *mux.Router) {
	db := database.GetConnection()

	r := func (path ...string) string {
		pathStr := "/api/" + os.Getenv("API_VERSION")
		for _, s := range path {
			pathStr += "/" + s
		}
		return pathStr
	}

	/** Auth routes */
	authHandler := handler.AuthHandler{
		DB: db,
		AuthService: service.AuthService{
			DB: db,
			SessionService: service.SessionService{DB: db},
		},
	}
	router.HandleFunc(r("login"), authHandler.Login).Methods("POST")
	router.HandleFunc(r("logout"), authMW(authHandler.Logout)).Methods("POST")

	/** User routes */
	userHandler := handler.UserHandler{
		UserService: service.UserService{
			DB: db,
			AuthService: service.AuthService{DB: db},
		},
	}
	router.HandleFunc(r("users"), userHandler.ListUsers).Methods("GET")
	router.HandleFunc(r("users", "register"), userHandler.RegisterUser).Methods("POST")

	/** Todo routes */
	todoHandler := handler.TodoHandler{
		TodoService: service.TodoService{DB: db},
	}
	router.HandleFunc(r("todos"), authMW(todoHandler.NewTodo)).Methods("POST")
	router.HandleFunc(r("todos"), authMW(todoHandler.ListTodos)).Methods("GET")
	router.HandleFunc(r("todos", "{id}"), authMW(todoHandler.UpdateTodo)).Methods("PUT")
	router.HandleFunc(r("todos", "{id}"), authMW(todoHandler.DeleteTodo)).Methods("DELETE")
}

func logMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s - %s (%s)", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

func authMW(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ss := service.SessionService{DB: database.GetConnection()}
		cookie, err := r.Cookie("access_token")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}
		user := ss.SessionUser(cookie.Value)
		if user == nil || user.ID == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}
		ctx := context.WithValue(r.Context(), types.UserKey{}, user)
		r = r.WithContext(ctx)
		next(w, r)
	}
}