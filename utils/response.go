package utils

import (
	"encoding/json"
	"net/http"
)

type ApiResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func JsonResponse(w http.ResponseWriter, status int, message string, data any) {
	w.Header().Set("content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ApiResponse{
		Message: message,
		Data:    data,
	})
}
