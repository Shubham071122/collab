package user

import "time"

type User struct {
	ID                  string     `json:"id"`
	Name                string     `json:"name"`
	Email               string     `json:"email"`
	PasswordHash        string     `json:"-"`
	IsVerified          bool       `json:"is_verified"`
	VerificationCode    *string    `json:"verification_code"`
	VerificationExpires *time.Time `json:"verification_expires"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}