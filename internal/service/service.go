package service

import (
	"errors"
	"time"

	"github.com/gopher-95/go-subscription-api/internal/models"
)

type SubscriptionStorage interface {
	Create(sub models.Subscription) (int64, error)
	Get(id int64) (*models.Subscription, error)
	Update(id int64, sub models.UpdateSubscriptionRequest) error
	Delete(id int64) error
	GetAll(limit int64, offest int64) (models.Subscription, error)
}

type Service struct {
	storage SubscriptionStorage
}

func NewService(storage SubscriptionStorage) *Service {
	return &Service{storage: storage}
}

func (s *Service) Create(sub models.CreateSubscriptionRequest) (*models.Subscription, error) {
	if sub.Price < 0 {
		return nil, errors.New("цена не может быть меньше нуля")
	}

	startDateParse, err := time.Parse("02-01-2006", sub.StartDate)
	if err != nil {
		return nil, errors.New("ошибка парсинга времени")
	}

	var endDateParse *time.Time
	if sub.EndDate != nil {
		parsed, err := time.Parse("02-01-2006", *sub.EndDate)
		if err != nil {
			return nil, errors.New("некорректный формат даты окончания подписки")
		}
		endDateParse = &parsed
	}

	if endDateParse != nil && endDateParse.Before(startDateParse) {
		return nil, errors.New("дата окончания подписки не может быть раньше даты начала подписки")
	}

	if startDateParse.Before(time.Now().Truncate(24 * time.Hour)) {
		return nil, errors.New("дата начала подписки не может быть в прошлом")
	}

	subscription := &models.Subscription{
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID,
		StartDate:   startDateParse,
		EndDate:     endDateParse,
	}

	id, err := s.storage.Create(*subscription)
	if err != nil {
		return nil, err
	}

	subscrip, err := s.storage.Get(id)
	if err != nil {
		return nil, err
	}

	return subscrip, nil
}
