package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"todo/db"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

type Todo struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type ApiResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type ResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (rw *ResponseWriter) WriteHeader(code int) {
	rw.StatusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func JsonResponse(w http.ResponseWriter, status int, message string, data any) {
	w.Header().Set("content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ApiResponse{
		Message: message,
		Data:    data,
	})
}

func GetAllTodos(w http.ResponseWriter, r *http.Request) {
	var todos []Todo
	err := db.DB.Select(&todos, "SELECT * FROM todos")
	if err != nil {
		fmt.Printf("[DB ERROR] %s\n", err)
		JsonResponse(w, http.StatusInternalServerError, "db error", nil)
		return
	}
	JsonResponse(w, http.StatusOK, "success", todos)
}

func AddNewTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		JsonResponse(w, http.StatusBadRequest, "invalid json", nil)
		return
	}
	// todo.Id = len(Todos) + 1
	// Todos = append(Todos, todo)
	err = db.DB.QueryRow(
		"INSERT INTO todos (title, content) VALUES ($1, $2) RETURNING id",
		todo.Title, todo.Content,
	).Scan(&todo.Id)

	if err != nil {
		fmt.Printf("[DB ERROR] %s\n", err)
		JsonResponse(w, http.StatusInternalServerError, "db error", nil)
		return
	}
	JsonResponse(w, http.StatusCreated, "todo created successfully", todo)
}

func GetTodoById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		JsonResponse(w, http.StatusNotFound, "id not found", nil)
		return
	}
	var todo Todo
	err = db.DB.QueryRow(
		"SELECT * FROM todos WHERE id = $1", id,
	).Scan(&todo.Id, &todo.Title, &todo.Content)

	if err != nil {
		fmt.Printf("[DB ERROR] %s\n", err)
		JsonResponse(w, http.StatusInternalServerError, "db error", nil)
		return
	}

	JsonResponse(w, http.StatusOK, "success", todo)

}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		JsonResponse(w, http.StatusNotFound, "id not found", nil)
		return
	}
	var todo Todo
	json.NewDecoder(r.Body).Decode(&todo)

	_, err = db.DB.Exec(
		"UPDATE todos SET title = $1, content = $2 WHERE id = $3",
		todo.Title, todo.Content, id,
	)
	if err != nil {
		JsonResponse(w, http.StatusNotFound, "todo not found", nil)
		return
	}
	JsonResponse(w, http.StatusOK, "todo Updated", todo)
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		JsonResponse(w, http.StatusNotFound, "id not found", nil)
		return
	}
	_, err = db.DB.Exec(
		"DELETE FROM todos WHERE id = $1",
		id,
	)
	if err != nil {
		JsonResponse(w, http.StatusNotFound, "todo not found", nil)
		return
	}
	JsonResponse(w, http.StatusOK, "todo delete successfully", nil)
}

func Middelware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &ResponseWriter{ResponseWriter: w, StatusCode: 200}
		next.ServeHTTP(rw, r)
		fmt.Printf("[%s] %s - %d\n", r.Method, r.URL.Path, rw.StatusCode)
	})
}

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

	r.Use(Middelware)

	r.Get("/api/todos", GetAllTodos)
	r.Post("/api/todos", AddNewTodo)
	r.Get("/api/todos/{id}", GetTodoById)
	r.Put("/api/todos/{id}", UpdateTodo)
	r.Delete("/api/todos/{id}", DeleteTodo)

	http.ListenAndServe(":8080", r)
}
