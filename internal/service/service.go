package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/gopher-95/go-subscription-api/internal/models"
)

type SubscriptionStorage interface {
	Create(sub *models.Subscription) (int, error)
	Get(id int) (*models.Subscription, error)
	Update(id int, sub *models.Subscription) (int, error)
	Delete(id int) (int, error)
	GetAll(limit int, offest int) ([]models.Subscription, error)
}

type Service struct {
	storage SubscriptionStorage
}

func NewService(storage SubscriptionStorage) *Service {
	return &Service{storage: storage}
}

func (s *Service) Create(sub models.UpdateCreateSubscriptionRequest) (int, error) {
	if sub.ServiceName == "" {
		return 0, errors.New("название подписки не указано")
	}

	if sub.Price < 0 {
		return 0, errors.New("цена не может быть меньше нуля")
	}

	if sub.UserID == "" {
		return 0, errors.New("не указан user_id")
	}

	if sub.StartDate == "" {
		return 0, errors.New("не указана дата начала подписки")
	}

	startDateParse, err := time.Parse("01-2006", sub.StartDate)
	if err != nil {
		return 0, errors.New("ошибка парсинга времени")
	}

	var endDateParse *time.Time
	if sub.EndDate != nil {
		parsed, err := time.Parse("01-2006", *sub.EndDate)
		if err != nil {
			return 0, errors.New("некорректный формат даты окончания подписки")
		}
		endDateParse = &parsed
	}

	if endDateParse != nil && endDateParse.Before(startDateParse) {
		return 0, errors.New("дата окончания подписки не может быть раньше даты начала подписки")
	}

	if startDateParse.Before(time.Now().Truncate(24 * time.Hour)) {
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
		return 0, fmt.Errorf("ошибка отправки данных в бд: %w", err)
	}

	return id, nil
}

func (s *Service) Get(id int) (*models.Subscription, error) {
	if id == 0 {
		return nil, errors.New("id не может быть равен нулю")
	}
	subscription, err := s.storage.Get(id)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения данных из бд: %w", err)
	}

	if subscription == nil {
		return nil, errors.New("запись с таким id не найдена")
	}

	return subscription, nil
}

func (s *Service) Update(id int, sub models.UpdateCreateSubscriptionRequest) (int, error) {
	//начало проверки входных данных из хэндлера
	if sub.ServiceName == "" {
		return 0, errors.New("ошибка выполнения update запроса, не указан service_name")
	}

	if sub.Price < 0 {
		return 0, errors.New("ошибка выполнения update запроса, не указан price")
	}

	if sub.UserID == "" {
		return 0, errors.New("не указан user_id")
	}

	if sub.StartDate == "" {
		return 0, errors.New("ошибка выполнения update запроса, не указан start_date")
	}

	startDateParse, err := time.Parse("01-2006", sub.StartDate)
	if err != nil {
		return 0, errors.New("не удалось запарсить время")
	}

	var endDateParse *time.Time
	if sub.EndDate != nil {
		parsed, err := time.Parse("01-2006", *sub.EndDate)
		if err != nil {
			return 0, errors.New("ошибка парсинга времени в update запросе: end_date")
		}

		endDateParse = &parsed

	}

	if endDateParse != nil && endDateParse.Before(startDateParse) {
		return 0, errors.New("дата окончания подписки не может быть раньше даты начала подписки")
	}

	if startDateParse.Before(time.Now().Truncate(24 * time.Hour)) {
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
		return 0, fmt.Errorf("ошибка выполнения update запроса к бд: %w", err)
	}

	return rowsAffected, nil
}

func (s *Service) Delete(id int) (int, error) {
	if id == 0 {
		return 0, errors.New("ошибка выполнения delete запроса: не указан id")
	}

	rowsAffected, err := s.storage.Delete(id)
	if err != nil {
		return 0, fmt.Errorf("ошибка выполнения delete запроса к бд: %w", err)
	}

	return rowsAffected, nil
}

func (s *Service) GetAll(limit, offset int) ([]models.Subscription, error) {
	if limit <= 0 {
		limit = 10
	}

	if offset < 0 {
		offset = 0
	}

	subscriptions, err := s.storage.GetAll(limit, offset)
	if err != nil {
		return subscriptions, fmt.Errorf("ошибка выполнения getall запроса к бд: %w", err)
	}

	return subscriptions, nil

}
