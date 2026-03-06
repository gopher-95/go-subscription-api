package service

import "github.com/gopher-95/go-subscription-api/internal/domain"

type SubscriptionRepo interface {
	GetByID(id int64) (*domain.Subscription, error)
	DeleteByID(int64) error
}

type SubscriptionService struct {
	repo SubscriptionRepo
}

func (s *SubscriptionService) GetByID(id int64) (*domain.Subscription, error) {
	return s.repo.GetByID(id)
}

func (s *SubscriptionService) DeleteByID(id int64) error {
	err := s.repo.DeleteByID(id)
	if err != nil {
		return err
	}

	return nil
}
