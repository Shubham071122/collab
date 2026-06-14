package auth

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"github.com/shubham071122/collab/internal/subscription"
	"github.com/shubham071122/collab/internal/user"
	"github.com/shubham071122/collab/pkg/email"
	"github.com/shubham071122/collab/pkg/jwt"
)

type Service struct {
	authRepo *Repository
	subRepo  *subscription.Repository
}

func NewService(authRepo *Repository, subRepo *subscription.Repository) *Service {
	return &Service{
		authRepo: authRepo,
		subRepo:  subRepo,
	}
}

func generateOTP() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%06d", r.Intn(900000)+100000)
}

func (s *Service) Register(req RegisterRequest) error {

	if req.Name == "" || req.Email == "" || req.Password == "" {
		return errors.New("All fields are required!")
	}

	emailStr := strings.ToLower(strings.TrimSpace(req.Email))

	if _, err := s.authRepo.GetUserByEmail(emailStr); err == nil {
		return errors.New("User already exists!")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	otpCode := generateOTP()
	otpExpires := time.Now().Add(30 * time.Minute)

	userID, err := s.authRepo.CreateUser(
		req.Name,
		emailStr,
		string(hashedPassword),
		otpCode,
		otpExpires,
	)
	if err != nil {
		return err
	}

	err = s.subRepo.CreateDefaultSubscription(userID)
	if err != nil {
		return err
	}

	go func() {
		err := email.SendVerificationEmail(emailStr, req.Name, otpCode)
		if err != nil {
			fmt.Printf("Error sending verification email: %v\n", err)
		}
	}()

	return nil
}

func (s *Service) Login(req LoginRequest) (*AuthResponse, error) {

	if req.Email == "" || req.Password == "" {
		return nil, errors.New("Email and password are required!")
	}

	fmt.Printf("Attempting to find user with email: %s\n", req.Email)

	emailStr := strings.ToLower(strings.TrimSpace(req.Email))

	foundUser, err := s.authRepo.GetUserByEmail(emailStr)
	if err != nil {
		fmt.Printf("Error finding user: %v\n", err)
		return nil, errors.New("Invalid credentials!")
	}
	fmt.Printf("Found user: %+v\n", foundUser)

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(req.Password))

	if err != nil {
		return nil, errors.New("Invalid credentials!")
	}
	token, err := jwt.GenerateJWT(foundUser.ID, foundUser.IsVerified)
	if err != nil {
		return nil, err
	}
	response := &AuthResponse{
		Token: token,
		User: user.UserDTO{
			ID:         foundUser.ID,
			Name:       foundUser.Name,
			Email:      foundUser.Email,
			IsVerified: foundUser.IsVerified,
		},
	}

	return response, nil
}

func (s *Service) VerifyOTP(req VerifyOTPRequest) (*AuthResponse, error) {
	if req.Email == "" || req.Code == "" {
		return nil, errors.New("Email and verification code are required!")
	}

	emailStr := strings.ToLower(strings.TrimSpace(req.Email))
	foundUser, err := s.authRepo.GetUserByEmail(emailStr)
	if err != nil {
		return nil, errors.New("User not found!")
	}

	if foundUser.IsVerified {
		return nil, errors.New("User is already verified!")
	}

	if foundUser.VerificationCode == nil || *foundUser.VerificationCode != req.Code {
		return nil, errors.New("Invalid verification code!")
	}

	if foundUser.VerificationExpires == nil || foundUser.VerificationExpires.Before(time.Now()) {
		return nil, errors.New("Verification code has expired!")
	}

	err = s.authRepo.VerifyUser(foundUser.ID)
	if err != nil {
		return nil, err
	}

	token, err := jwt.GenerateJWT(foundUser.ID, true)
	if err != nil {
		return nil, err
	}

	response := &AuthResponse{
		Token: token,
		User: user.UserDTO{
			ID:         foundUser.ID,
			Name:       foundUser.Name,
			Email:      foundUser.Email,
			IsVerified: true,
		},
	}

	return response, nil
}

func (s *Service) ResendOTP(req ResendOTPRequest) error {
	if req.Email == "" {
		return errors.New("Email is required!")
	}

	emailStr := strings.ToLower(strings.TrimSpace(req.Email))
	foundUser, err := s.authRepo.GetUserByEmail(emailStr)
	if err != nil {
		return errors.New("User not found!")
	}

	if foundUser.IsVerified {
		return errors.New("User is already verified!")
	}

	otpCode := generateOTP()
	otpExpires := time.Now().Add(30 * time.Minute)

	err = s.authRepo.UpdateVerificationCode(foundUser.ID, otpCode, otpExpires)
	if err != nil {
		return err
	}

	go func() {
		err := email.SendVerificationEmail(emailStr, foundUser.Name, otpCode)
		if err != nil {
			fmt.Printf("Error sending verification email: %v\n", err)
		}
	}()

	return nil
}

func (s *Service) Logout() error {
	// Implement logout logic here
	return nil
}
