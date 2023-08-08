package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(
		cors.Handler(cors.Options{
			AllowedOrigins:   []string{"https://", "http://"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			AllowCredentials: false,
			MaxAge:           300,
		}))

	mux.Get("/api/test", app.TestHandler)
	mux.Get("/api/todos", app.GetToDos)
	mux.Get("/api/todos/{id}", app.GetToDoByID)
	mux.Post("/api/todos", app.CreateTodo)
	mux.Put("/api/todos", app.EditToDo)
	mux.Delete("/api/todos/{id}", app.DeleteToDo)

	mux.Get("/api/users", app.GetUsers)
	mux.Get("/api/users/{id}", app.GetUserByID)
	mux.Post("/api/users", app.CreateUser)
	mux.Put("/api/users", app.EditUser)
	mux.Delete("/api/users/{id}", app.DeleteUser)
	return mux
}
