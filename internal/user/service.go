package user

import (
	"errors"
)

type Service struct {
	userRepo *Repository
}

func NewService(userRepo *Repository) *Service {
	return &Service{
		userRepo: userRepo,
	}
}

func (s *Service) CurrentUser(userID string) (*User, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	return s.userRepo.GetUserByID(userID)
}
