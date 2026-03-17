package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/gopher-95/go-subscription-api/internal/models"
	"github.com/gopher-95/go-subscription-api/internal/service"
)

type SubscriptionHandler struct {
	service *service.Service
}

func NewSubscriptionHandler(service *service.Service) *SubscriptionHandler {
	return &SubscriptionHandler{service: service}
}

func (handler *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	subscriptionRequest := models.UpdateCreateSubscriptionRequest{}
	err := json.NewDecoder(r.Body).Decode(&subscriptionRequest)
	if err != nil {
		jsonError(w, http.StatusBadRequest, "не удалось выполнить create запрос")
		return
	}

	id, err := handler.service.Create(subscriptionRequest)
	if err != nil {
		switch err.Error() {
		case "название подписки не указано",
			"цена не может быть меньше нуля",
			"не указан user_id",
			"не указана дата начала подписки",
			"ошибка парсинга времени",
			"некорректный формат даты окончания подписки",
			"дата окончания подписки не может быть раньше даты начала подписки",
			"дата начала подписки не может быть в прошлом":
			jsonError(w, http.StatusBadRequest, err.Error())
		default:
			// Неизвестная ошибка - проблема с БД или сервером
			jsonError(w, http.StatusInternalServerError, "внутренняя ошибка сервера")

		}

		return
	}

	jsonResponse(w, http.StatusCreated, map[string]interface{}{
		"message": "подписка успешно создана",
		"id":      id,
	})
}

func (handler *SubscriptionHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := readIDParam(r)
	if err != nil {
		jsonError(w, http.StatusBadRequest, "неправильно указан id")
		return
	}

	subscription, err := handler.service.Get(id)
	if err != nil {
		if err.Error() == "запись с таким id не найдена" {
			jsonError(w, http.StatusNotFound, "подписка не найдена ")
			return
		}
		jsonError(w, http.StatusInternalServerError, "не удалось выполнить запрос к бд")
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"message":      "подписка успешно получена",
		"subscription": subscription,
	})
}

func (handler *SubscriptionHandler) Update(w http.ResponseWriter, r *http.Request) {

	id, err := readIDParam(r)
	if err != nil {
		jsonError(w, http.StatusBadRequest, "неверно указан id записи, которую хотите изменить")
		return
	}

	var updateRequest models.UpdateCreateSubscriptionRequest

	err = json.NewDecoder(r.Body).Decode(&updateRequest)
	if err != nil {
		jsonError(w, http.StatusBadRequest, "некорректный JSON")
		return
	}

	rowsAffected, err := handler.service.Update(id, updateRequest)
	if err != nil {
		switch err.Error() {
		case "ошибка выполнения update запроса, не указан service_name":
			jsonError(w, http.StatusBadRequest, "название сервиса не может быть пустым")

		case "ошибка выполнения update запроса, не указан price":
			jsonError(w, http.StatusBadRequest, "цена не указана")

		case "не указан user_id":
			jsonError(w, http.StatusBadRequest, "не указан user_id")

		case "ошибка выполнения update запроса, не указан start_date":
			jsonError(w, http.StatusBadRequest, "не указана дата начала")

		case "не удалось запарсить время":
			jsonError(w, http.StatusBadRequest, "некорректный формат даты начала. Используйте MM-YYYY")

		case "ошибка парсинга времени в update запросе: end_date":
			jsonError(w, http.StatusBadRequest, "некорректный формат даты окончания. Используйте MM-YYYY")

		case "дата окончания подписки не может быть раньше даты начала подписки":
			jsonError(w, http.StatusBadRequest, "дата окончания не может быть раньше даты начала")

		case "дата начала подписки не может быть в прошлом":
			jsonError(w, http.StatusBadRequest, "дата начала не может быть в прошлом")

		default:
			if err.Error() == "ошибка выполнения update запроса к бд: sql: no rows in result set" {
				jsonError(w, http.StatusNotFound, "подписка с указанным id не найдена")
				return
			}
			jsonError(w, http.StatusInternalServerError, "внутренняя ошибка сервера")
		}
		return
	}

	if rowsAffected == 0 {
		jsonError(w, http.StatusNotFound, "подписка с указанным id не найдена")
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"message":       "подписка успешно обновлена",
		"id":            id,
		"rows_affected": rowsAffected,
	})

}

func (handler *SubscriptionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := readIDParam(r)
	if err != nil {
		jsonError(w, http.StatusBadRequest, "неверно указан id")
		return
	}

	deletedRows, err := handler.service.Delete(id)
	if err != nil {
		switch err.Error() {
		case "ошибка выполнения delete запроса: не указан id":
			jsonError(w, http.StatusBadRequest, "не указан id в delete запросе")
			return
		default:
			jsonError(w, http.StatusInternalServerError, "ошибка сервера")
			return
		}
	}

	if deletedRows == 0 {
		jsonError(w, http.StatusNotFound, "подписка не найдена")
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"message":       "подписка успешно удалена",
		"deleted_count": deletedRows,
	})
}

func (handler *SubscriptionHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10
	offset := 0

	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil || parsedLimit < 1 {
			jsonError(w, http.StatusBadRequest, "limit должен быть положительным числом")
			return
		}
		limit = parsedLimit
	}

	if offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err != nil || parsedOffset < 0 {
			jsonError(w, http.StatusBadRequest, "offset должен быть неотрицательным числом")
			return
		}
		offset = parsedOffset
	}

	subscriptions, err := handler.service.GetAll(limit, offset)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "не удалось получить список подписок")
		return
	}

	if subscriptions == nil {
		subscriptions = []models.Subscription{}
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"subscriptions": subscriptions,
		"pagination": map[string]interface{}{
			"limit":  limit,
			"offset": offset,
			"count":  len(subscriptions),
		},
	})
}

func (handler *SubscriptionHandler) GetTotalCost(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	startDate := query.Get("start_date")
	endDate := query.Get("end_date")
	userID := query.Get("user_id")
	serviceName := query.Get("service_name")

	if startDate == "" {
		jsonError(w, http.StatusBadRequest, "не указан параметр start_date")
		return
	}
	if endDate == "" {
		jsonError(w, http.StatusBadRequest, "не указан параметр end_date")
		return
	}

	req := models.TotalCostRequest{
		StartDate: startDate,
		EndDate:   endDate,
	}

	if userID != "" {
		req.UserID = &userID
	}
	if serviceName != "" {
		req.ServiceName = &serviceName
	}

	result, err := handler.service.CalculateTotalCost(req)
	if err != nil {
		switch err.Error() {
		case "не указана дата начала периода",
			"не указана дата окончания периода",
			"некорректный формат даты начала периода. Используйте MM-YYYY",
			"некорректный формат даты окончания периода. Используйте MM-YYYY",
			"дата начала периода не может быть позже даты окончания":
			jsonError(w, http.StatusBadRequest, err.Error())
		default:
			jsonError(w, http.StatusInternalServerError, "внутренняя ошибка сервера")
		}
		return
	}

	// 5. Успешный ответ
	jsonResponse(w, http.StatusOK, result)
}

func readIDParam(r *http.Request) (int, error) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if id < 1 || err != nil {
		return 0, errors.New("некорректный id")
	}

	return id, nil
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
func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
