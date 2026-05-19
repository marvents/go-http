package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Todo struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type ApiResponse struct {
	Message string `json:"message"`
	Data any `json:"data"`
}

var Todos []Todo

func Spliter(p string) []string {
	return strings.Split(p, "/")
}

func getId(p string) string {
	return Spliter(p)[3]
}

func JsonResponse(w http.ResponseWriter, status int, message string, data any) {
	w.Header().Set("content-Type", "application/json");
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ApiResponse{
		Message: message,
		Data: data,
	})
}

func TodosHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		JsonResponse(w, http.StatusOK, "success", Todos)
	case http.MethodPost:
		var todo Todo
		json.NewDecoder(r.Body).Decode(&todo)
		todo.Id = len(Todos) + 1;
		Todos = append(Todos, todo);
		JsonResponse(w, http.StatusCreated, "todo created successfully", todo)
	}
}

func TodoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(getId(r.URL.Path), 10, 32)
	if err != nil {
		fmt.Print(err)
	}

	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		for _, t := range Todos {
			if t.Id == int(id) {
				JsonResponse(w, http.StatusOK, "success", t)
				return
			}
		}

	case http.MethodPut:
		// I already build it before that
	case http.MethodDelete:
		var FiltredTodos []Todo
		for _, t := range Todos {
			if t.Id == int(id) {
				continue
			}
			FiltredTodos = append(FiltredTodos, t)
		}
		Todos = FiltredTodos;
		JsonResponse(w, http.StatusOK, "todo delete successfully", nil)
	}

}

func main() {
	http.HandleFunc("/api/todos", TodosHandler)
	http.HandleFunc("/api/todos/", TodoHandler)

	http.ListenAndServe(":8080", nil)
}
