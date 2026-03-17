package handlers

import (
	"encoding/json"
	"errors"
	"log"
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
	log.Println("инициализация HTTP обработчиков")
	return &SubscriptionHandler{service: service}
}

func (handler *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	log.Println("POST /api/v1/subscriptions - создание подписки")
	subscriptionRequest := models.UpdateCreateSubscriptionRequest{}
	err := json.NewDecoder(r.Body).Decode(&subscriptionRequest)
	if err != nil {
		log.Printf("ошибка декодирования JSON: %v", err)
		jsonError(w, http.StatusBadRequest, "не удалось выполнить create запрос")
		return
	}

	log.Printf("📦 Данные запроса: service=%s, price=%d, user=%s",
		subscriptionRequest.ServiceName, subscriptionRequest.Price, subscriptionRequest.UserID)

	id, err := handler.service.Create(subscriptionRequest)
	if err != nil {
		log.Printf("ошибка создания подписки: %v", err)
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
	log.Printf("✅ Подписка создана: id=%d", id)
	jsonResponse(w, http.StatusCreated, map[string]interface{}{
		"message": "подписка успешно создана",
		"id":      id,
	})
}

func (handler *SubscriptionHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := readIDParam(r)
	if err != nil {
		log.Printf("ошибка получения ID из URL: %v", err)
		jsonError(w, http.StatusBadRequest, "неправильно указан id")
		return
	}

	log.Printf("GET /api/v1/subscriptions/%d", id)

	subscription, err := handler.service.Get(id)
	if err != nil {
		if err.Error() == "запись с таким id не найдена" {
			log.Printf("подписка id=%d не найдена", id)
			jsonError(w, http.StatusNotFound, "подписка не найдена ")
			return
		}
		log.Printf("ошибка получения подписки id=%d: %v", id, err)
		jsonError(w, http.StatusInternalServerError, "не удалось выполнить запрос к бд")
		return
	}

	log.Printf("✅ Подписка id=%d получена", id)
	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"message":      "подписка успешно получена",
		"subscription": subscription,
	})
}

func (handler *SubscriptionHandler) Update(w http.ResponseWriter, r *http.Request) {

	id, err := readIDParam(r)
	if err != nil {
		log.Printf("ошибка получения ID из URL: %v", err)
		jsonError(w, http.StatusBadRequest, "неверно указан id записи, которую хотите изменить")
		return
	}

	log.Printf("PUT /api/v1/subscriptions/%d - обновление", id)

	var updateRequest models.UpdateCreateSubscriptionRequest

	err = json.NewDecoder(r.Body).Decode(&updateRequest)
	if err != nil {
		log.Printf("ошибка декодирования JSON: %v", err)
		jsonError(w, http.StatusBadRequest, "некорректный JSON")
		return
	}

	log.Printf("данные обновления: service=%s, price=%d, start=%s",
		updateRequest.ServiceName, updateRequest.Price, updateRequest.StartDate)

	rowsAffected, err := handler.service.Update(id, updateRequest)
	if err != nil {
		log.Printf("ошибка обновления подписки id=%d: %v", id, err)

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
	log.Printf("ошибка обновления подписки id=%d: %v", id, err)

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"message":       "подписка успешно обновлена",
		"id":            id,
		"rows_affected": rowsAffected,
	})

}

func (handler *SubscriptionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := readIDParam(r)
	if err != nil {
		log.Printf("ошибка получения ID из URL: %v", err)
		jsonError(w, http.StatusBadRequest, "неверно указан id")
		return
	}

	log.Printf("DELETE /api/v1/subscriptions/%d", id)

	deletedRows, err := handler.service.Delete(id)
	if err != nil {
		log.Printf("ошибка удаления подписки id=%d: %v", id, err)

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
		log.Printf("подписка id=%d не найдена для удаления", id)
		jsonError(w, http.StatusNotFound, "подписка не найдена")
		return
	}

	log.Printf("✅ Подписка id=%d удалена", id)
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

	log.Printf("GET /api/v1/subscriptions - список (limit=%s, offset=%s)", limitStr, offsetStr)

	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil || parsedLimit < 1 {
			log.Printf("некорректный limit: %s", limitStr)
			jsonError(w, http.StatusBadRequest, "limit должен быть положительным числом")
			return
		}
		limit = parsedLimit
	}

	if offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err != nil || parsedOffset < 0 {
			log.Printf("некорректный offset: %s", offsetStr)
			jsonError(w, http.StatusBadRequest, "offset должен быть неотрицательным числом")
			return
		}
		offset = parsedOffset
	}

	subscriptions, err := handler.service.GetAll(limit, offset)
	if err != nil {
		log.Printf("ошибка получения списка подписок: %v", err)
		jsonError(w, http.StatusInternalServerError, "не удалось получить список подписок")
		return
	}

	if subscriptions == nil {
		subscriptions = []models.Subscription{}
	}

	log.Printf("✅ Получено подписок: %d (limit=%d, offset=%d)",
		len(subscriptions), limit, offset)

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

	log.Printf("GET /api/v1/subscriptions/total-cost - период: %s - %s, user=%s, service=%s",
		startDate, endDate, userID, serviceName)

	if startDate == "" {
		log.Printf("отсутствует обязательный параметр start_date")
		jsonError(w, http.StatusBadRequest, "не указан параметр start_date")
		return
	}
	if endDate == "" {
		log.Printf("отсутствует обязательный параметр end_date")
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
		log.Printf("ошибка расчета стоимости: %v", err)

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

	log.Printf("расчет стоимости завершен: %d₽ за период %s - %s",
		result.TotalCost, startDate, endDate)

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
	log.Printf("ответ: статус %d, ошибка: %s", status, message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

// функция для формирования json ответа при успешном запросе
func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	log.Printf("ответ: статус %d", status)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
