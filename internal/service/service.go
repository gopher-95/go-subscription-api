package service

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gopher-95/go-subscription-api/internal/models"
)

type SubscriptionStorage interface {
	Create(sub *models.Subscription) (int, error)
	Get(id int) (*models.Subscription, error)
	Update(id int, sub *models.Subscription) (int, error)
	Delete(id int) (int, error)
	GetAll(limit int, offest int) ([]models.Subscription, error)
	GetSubscriptionsForPeriod(startDate, endDate time.Time, userID, serviceName *string) ([]models.Subscription, error)
}

type Service struct {
	storage SubscriptionStorage
}

func NewService(storage SubscriptionStorage) *Service {
	log.Println("инициализация сервиса подписок")
	return &Service{storage: storage}
}

func (s *Service) Create(sub models.UpdateCreateSubscriptionRequest) (int, error) {
	log.Printf("создание подписки: service=%s, price=%d, user=%s", sub.ServiceName, sub.Price, sub.UserID)
	if sub.ServiceName == "" {
		log.Println("ошибка валидации: название подписки не указано")
		return 0, errors.New("название подписки не указано")
	}

	if sub.Price < 0 {
		log.Printf("ошибка валидации: отрицательная цена %d", sub.Price)
		return 0, errors.New("цена не может быть меньше нуля")
	}

	if sub.UserID == "" {
		log.Println("ошибка валидации: не указан user_id")
		return 0, errors.New("не указан user_id")
	}

	if sub.StartDate == "" {
		log.Println("ошибка валидации: не указана дата начала подписки")
		return 0, errors.New("не указана дата начала подписки")
	}

	startDateParse, err := time.Parse("01-2006", sub.StartDate)
	if err != nil {
		log.Printf("ошибка парсинга даты начала: %s", sub.StartDate)
		return 0, errors.New("ошибка парсинга времени")
	}

	var endDateParse *time.Time
	if sub.EndDate != nil {
		parsed, err := time.Parse("01-2006", *sub.EndDate)
		if err != nil {
			log.Printf("ошибка парсинга даты окончания: %s", *sub.EndDate)
			return 0, errors.New("некорректный формат даты окончания подписки")
		}
		endDateParse = &parsed
	}

	if endDateParse != nil && endDateParse.Before(startDateParse) {
		log.Printf("ошибка: дата окончания %v раньше даты начала %v", endDateParse, startDateParse)
		return 0, errors.New("дата окончания подписки не может быть раньше даты начала подписки")
	}

	if startDateParse.Before(time.Now().Truncate(24 * time.Hour)) {
		log.Printf("ошибка: дата начала %v в прошлом", startDateParse)
		return 0, errors.New("дата начала подписки не может быть в прошлом")
	}

	subscription := &models.Subscription{
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID,
		StartDate:   startDateParse,
		EndDate:     endDateParse,
	}

	id, err := s.storage.Create(subscription)
	if err != nil {
		log.Printf("ошибка создания подписки в бд: %v", err)
		return 0, fmt.Errorf("ошибка отправки данных в бд: %w", err)
	}
	log.Printf("подписка создана: id=%d", id)
	return id, nil
}

func (s *Service) Get(id int) (*models.Subscription, error) {
	log.Printf("получение подписки id=%d", id)
	if id == 0 {
		log.Println("ошибка: передан id=0")
		return nil, errors.New("id не может быть равен нулю")
	}
	subscription, err := s.storage.Get(id)
	if err != nil {
		log.Printf("ошибка получения подписки id=%d из бд: %v", id, err)
		return nil, fmt.Errorf("ошибка получения данных из бд: %w", err)
	}

	if subscription == nil {
		log.Printf("подписка id=%d не найдена", id)
		return nil, errors.New("запись с таким id не найдена")
	}

	log.Printf("подписка id=%d получена", id)
	return subscription, nil
}

func (s *Service) Update(id int, sub models.UpdateCreateSubscriptionRequest) (int, error) {

	log.Printf("обновление подписки id=%d: service=%s, price=%d", id, sub.ServiceName, sub.Price)
	if sub.ServiceName == "" {
		log.Println("ошибка валидации при обновлении: не указан service_name")
		return 0, errors.New("ошибка выполнения update запроса, не указан service_name")
	}

	if sub.Price < 0 {
		log.Printf("ошибка валидации при обновлении: отрицательная цена %d", sub.Price)
		return 0, errors.New("ошибка выполнения update запроса, не указан price")
	}

	if sub.UserID == "" {
		return 0, errors.New("не указан user_id")
	}

	if sub.StartDate == "" {
		log.Println("ошибка валидации при обновлении: не указан user_id")
		return 0, errors.New("ошибка выполнения update запроса, не указан start_date")
	}

	startDateParse, err := time.Parse("01-2006", sub.StartDate)
	if err != nil {
		log.Printf("ошибка парсинга даты начала при обновлении: %s", sub.StartDate)
		return 0, errors.New("не удалось запарсить время")
	}

	var endDateParse *time.Time
	if sub.EndDate != nil {
		parsed, err := time.Parse("01-2006", *sub.EndDate)
		if err != nil {
			log.Printf("ошибка парсинга даты окончания при обновлении: %s", *sub.EndDate)
			return 0, errors.New("ошибка парсинга времени в update запросе: end_date")
		}

		endDateParse = &parsed

	}

	if endDateParse != nil && endDateParse.Before(startDateParse) {
		log.Printf("ошибка: дата окончания %v раньше даты начала %v", endDateParse, startDateParse)
		return 0, errors.New("дата окончания подписки не может быть раньше даты начала подписки")
	}

	if startDateParse.Before(time.Now().Truncate(24 * time.Hour)) {
		log.Printf("ошибка: дата начала %v в прошлом", startDateParse)
		return 0, errors.New("дата начала подписки не может быть в прошлом")
	}

	subscription := models.Subscription{
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID,
		StartDate:   startDateParse,
		EndDate:     endDateParse,
	}

	rowsAffected, err := s.storage.Update(id, &subscription)
	if err != nil {
		log.Printf("ошибка обновления подписки id=%d в бд: %v", id, err)
		return 0, fmt.Errorf("ошибка выполнения update запроса к бд: %w", err)
	}

	log.Printf("подписка id=%d обновлена, затронуто строк: %d", id, rowsAffected)
	return rowsAffected, nil
}

func (s *Service) Delete(id int) (int, error) {
	log.Printf("удаление подписки id=%d", id)
	if id == 0 {
		log.Println("ошибка: попытка удаления с id=0")
		return 0, errors.New("ошибка выполнения delete запроса: не указан id")
	}

	rowsAffected, err := s.storage.Delete(id)
	if err != nil {
		log.Printf("ошибка удаления подписки id=%d из бд: %v", id, err)
		return 0, fmt.Errorf("ошибка выполнения delete запроса к бд: %w", err)
	}

	if rowsAffected == 0 {
		log.Printf("подписка id=%d не найдена для удаления", id)
	} else {
		log.Printf("подписка id=%d удалена, затронуто строк: %d", id, rowsAffected)
	}

	return rowsAffected, nil
}

func (s *Service) GetAll(limit, offset int) ([]models.Subscription, error) {
	log.Printf("получение списка подписок: limit=%d, offset=%d", limit, offset)
	if limit <= 0 {
		limit = 10
		log.Printf("limit скорректирован на %d", limit)
	}

	if offset < 0 {
		offset = 0
		log.Printf("offset скорректирован на %d", offset)
	}

	subscriptions, err := s.storage.GetAll(limit, offset)
	if err != nil {
		log.Printf("ошибка получения списка подписок: %v", err)
		return subscriptions, fmt.Errorf("ошибка выполнения getall запроса к бд: %w", err)
	}

	log.Printf("получено подписок: %d", len(subscriptions))
	return subscriptions, nil

}

func (s *Service) CalculateTotalCost(req models.TotalCostRequest) (*models.TotalCostResponse, error) {

	log.Printf("расчет стоимости за период: %s - %s", req.StartDate, req.EndDate)

	if req.UserID != nil {
		log.Printf("фильтр по user_id: %s", *req.UserID)
	}
	if req.ServiceName != nil {
		log.Printf("фильтр по service_name: %s", *req.ServiceName)
	}

	if req.StartDate == "" {
		log.Println("ошибка: не указана дата начала периода")
		return nil, errors.New("не указана дата начала периода")
	}

	if req.EndDate == "" {
		log.Println("ошибка: не указана дата окончания периода")
		return nil, errors.New("не указана дата окончания периода")
	}

	startPeriod, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		log.Printf("ошибка парсинга даты начала периода: %s", req.StartDate)
		return nil, errors.New("некорректный формат даты начала периода. Используйте MM-YYYY")
	}

	endPeriod, err := time.Parse("01-2006", req.EndDate)
	if err != nil {
		log.Printf("ошибка парсинга даты окончания периода: %s", req.EndDate)
		return nil, errors.New("некорректный формат даты окончания периода. Используйте MM-YYYY")
	}

	if startPeriod.After(endPeriod) {
		log.Printf("ошибка: начало периода %v позже конца %v", startPeriod, endPeriod)
		return nil, errors.New("дата начала периода не может быть позже даты окончания")
	}

	subscriptions, err := s.storage.GetSubscriptionsForPeriod(
		startPeriod,
		endPeriod,
		req.UserID,
		req.ServiceName,
	)
	if err != nil {
		log.Printf("ошибка получения подписок за период: %v", err)
		return nil, fmt.Errorf("ошибка получения подписок за период: %w", err)
	}

	log.Printf("получено подписок для расчета: %d", len(subscriptions))

	totalCost := 0
	for _, sub := range subscriptions {

		subStart := sub.StartDate
		subEnd := time.Now()
		if sub.EndDate != nil {
			subEnd = *sub.EndDate
		}

		calcStart := subStart
		if startPeriod.After(subStart) {
			calcStart = startPeriod
		}

		calcEnd := subEnd
		if endPeriod.Before(subEnd) {
			calcEnd = endPeriod
		}

		if calcStart.After(calcEnd) {
			continue
		}

		months := 0
		for y := calcStart.Year(); y <= calcEnd.Year(); y++ {
			startMonth := 1
			if y == calcStart.Year() {
				startMonth = int(calcStart.Month())
			}
			endMonth := 12
			if y == calcEnd.Year() {
				endMonth = int(calcEnd.Month())
			}
			months += (endMonth - startMonth + 1)
		}
		subCost := sub.Price * months
		totalCost += subCost

		log.Printf("подписка id=%d (%s): %d руб x %d мес = %d руб",
			sub.ID, sub.ServiceName, sub.Price, months, subCost)
	}

	response := &models.TotalCostResponse{
		TotalCost: totalCost,
		Period:    req.StartDate + " - " + req.EndDate,
	}

	if req.UserID != nil {
		response.UserID = *req.UserID
	}
	if req.ServiceName != nil {
		response.Service = *req.ServiceName
	}

	log.Printf("итоговая стоимость за период: %d руб", totalCost)
	return response, nil
}
