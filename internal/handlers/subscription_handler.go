package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/gopher-95/go-subscription-api/internal/service"
)

type SubscriptionHandler struct {
	service *service.SubscriptionService
}

func NewSubscriptionHandler(s *service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{
		service: s,
	}
}

func (h *SubscriptionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "ошибка сервера")
		return
	}

	subscription, err := h.service.GetByID(int64(idInt))
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "не удалось получить запись по id")
		return
	}

	jsonResponse(w, http.StatusOK, "запись успешно получена", "id", subscription.ID)

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
