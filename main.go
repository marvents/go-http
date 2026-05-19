package main

import (
	"log"
	"net/http"
	"todo/db"
	"todo/middleware"
	"todo/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)


func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("[ERROR] filed to load donenv")
	}

	err = db.Connect()
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()

	r.Use(middleware.Middleware)

	r.Get("/api/todos", handlers.GetAllTodos)
	r.Post("/api/todos", handlers.AddNewTodo)
	r.Get("/api/todos/{id}", handlers.GetTodoById)
	r.Put("/api/todos/{id}", handlers.UpdateTodo)
	r.Delete("/api/todos/{id}", handlers.DeleteTodo)

	http.ListenAndServe(":8080", r)
}
