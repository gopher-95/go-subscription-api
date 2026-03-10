package handlers

import (
	"encoding/json"
	"net/http"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
	w.WriteHeader(http.StatusOK)
}

// функция для формирования json ответа при ошибках
func jsonError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

// функция для формирования json ответа при успешном запросе
func jsonResponse(w http.ResponseWriter, status int, message string, entity string, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": message,
		entity:    response,
	})
}
