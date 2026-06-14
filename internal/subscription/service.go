package subscription

import (
	"errors"
)

type Service struct {
	subRepo *Repository
}

func NewService(subRepo *Repository) *Service {
	return &Service{
		subRepo: subRepo,
	}
}

func (s *Service) GetSubscriptionByUserID(userID string) (*UserSubscription, error) {
	return s.subRepo.GetSubscriptionByUserID(userID)
}

func (s *Service) CreateDefaultSubscription(userID string) error {
	return s.subRepo.CreateDefaultSubscription(userID)
}

func (s *Service) CheckProjectLimit(userID string, currentCount int) (bool, error) {
	sub, err := s.subRepo.GetSubscriptionByUserID(userID)
	if err != nil {
		return false, err
	}

	config, exists := Plans[sub.Tier]
	if !exists {
		return false, errors.New("unknown subscription tier")
	}

	if config.MaxProjects == -1 {
		return true, nil
	}

	return currentCount < config.MaxProjects, nil
}

func (s *Service) CheckShareLimit(userID string, currentCount int) (bool, error) {
	sub, err := s.subRepo.GetSubscriptionByUserID(userID)
	if err != nil {
		return false, err
	}

	config, exists := Plans[sub.Tier]
	if !exists {
		return false, errors.New("unknown subscription tier")
	}

	if config.MaxShares == -1 {
		return true, nil
	}

	return currentCount < config.MaxShares, nil
}
