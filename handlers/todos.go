package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"todo/db"
	"todo/models"
	"todo/utils"
	"github.com/go-chi/chi/v5"
)


func GetAllTodos(w http.ResponseWriter, r *http.Request) {
	var todos []models.Todo
	err := db.DB.Select(&todos, "SELECT * FROM todos")
	if err != nil {
		fmt.Printf("[DB ERROR] %s\n", err)
		utils.JsonResponse(w, http.StatusInternalServerError, "db error", nil)
		return
	}
	utils.JsonResponse(w, http.StatusOK, "success", todos)
}


func AddNewTodo(w http.ResponseWriter, r *http.Request) {
	var todo models.Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		utils.JsonResponse(w, http.StatusBadRequest, "invalid json", nil)
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
		utils.JsonResponse(w, http.StatusInternalServerError, "db error", nil)
		return
	}
	utils.JsonResponse(w, http.StatusCreated, "todo created successfully", todo)
}

func GetTodoById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.JsonResponse(w, http.StatusNotFound, "id not found", nil)
		return
	}
	var todo models.Todo
	err = db.DB.QueryRow(
		"SELECT * FROM todos WHERE id = $1", id,
	).Scan(&todo.Id, &todo.Title, &todo.Content)

	if err != nil {
		fmt.Printf("[DB ERROR] %s\n", err)
		utils.JsonResponse(w, http.StatusInternalServerError, "db error", nil)
		return
	}

	utils.JsonResponse(w, http.StatusOK, "success", todo)

}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.JsonResponse(w, http.StatusNotFound, "id not found", nil)
		return
	}
	var todo models.Todo
	json.NewDecoder(r.Body).Decode(&todo)

	_, err = db.DB.Exec(
		"UPDATE todos SET title = $1, content = $2 WHERE id = $3",
		todo.Title, todo.Content, id,
	)
	if err != nil {
		utils.JsonResponse(w, http.StatusNotFound, "todo not found", nil)
		return
	}
	utils.JsonResponse(w, http.StatusOK, "todo Updated", todo)
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.JsonResponse(w, http.StatusNotFound, "id not found", nil)
		return
	}
	_, err = db.DB.Exec(
		"DELETE FROM todos WHERE id = $1",
		id,
	)
	if err != nil {
		utils.JsonResponse(w, http.StatusNotFound, "todo not found", nil)
		return
	}
	utils.JsonResponse(w, http.StatusOK, "todo delete successfully", nil)
}