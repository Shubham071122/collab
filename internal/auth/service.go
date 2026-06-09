package auth

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"github.com/shubham071122/collab/internal/user"
	"github.com/shubham071122/collab/pkg/jwt"
)

type Service struct{
	authRepo *Repository
}

func NewService(authRepo *Repository) *Service {
	return &Service{
		authRepo: authRepo,
	}
}

func (s *Service) Register(req RegisterRequest) error {

	if req.Name == "" || req.Email == "" || req.Password == "" {
		return errors.New("All fields are required!")
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))

	if _, err := s.authRepo.GetUserByEmail(email); err == nil {
		return errors.New("User already exists!")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = s.authRepo.CreateUser(
		req.Name,
		email,
		string(hashedPassword),
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Login(req LoginRequest) (*AuthResponse, error) {

	if req.Email == "" || req.Password == "" {
		return nil, errors.New("Email and password are required!")
	}

	fmt.Printf("Attempting to find user with email: %s\n", req.Email)

	email := strings.ToLower(strings.TrimSpace(req.Email))

	foundUser, err := s.authRepo.GetUserByEmail(email)
	if err != nil {
		fmt.Printf("Error finding user: %v\n", err)
		return nil, errors.New("Invalid credentials!")
	}
	fmt.Printf("Found user: %+v\n", foundUser)

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(req.Password))

	if err != nil {
		return nil, errors.New("Invalid credentials!")
	}
	token, err := jwt.GenerateJWT(foundUser.ID)
	if err != nil {
		return nil, err
	}
	response := &AuthResponse{
		Token: token,
		User: user.UserDTO{
			ID:    foundUser.ID,
			Name:  foundUser.Name,
			Email: foundUser.Email,
		},
	}

	return response, nil
}


func (s *Service) Logout() error {
	// Implement logout logic here
	return nil
}
