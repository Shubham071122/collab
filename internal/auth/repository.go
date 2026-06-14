package auth

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/shubham071122/collab/internal/user"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateUser(
	name string,
	email string,
	passwordHash string,
	verificationCode string,
	verificationExpires time.Time,
) (string, error) {

	query := `
		INSERT INTO users (
			id,
			name,
			email,
			password_hash,
			is_verified,
			verification_code,
			verification_expires
		)
		VALUES ($1, $2, lower($3), $4, FALSE, $5, $6)
	`

	userID := uuid.New().String()
	_, err := r.db.Exec(
		query,
		userID,
		name,
		email,
		passwordHash,
		verificationCode,
		verificationExpires,
	)

	if err != nil {
		return "", err
	}

	return userID, nil
}

func (r *Repository) GetUserByEmail(email string) (*user.User, error) {

	query := `
		SELECT
			id,
			name,
			email,
			password_hash,
			is_verified,
			verification_code,
			verification_expires,
			created_at,
			updated_at
		FROM users
		WHERE lower(email) = lower($1)
	`

	var u user.User

	err := r.db.QueryRow(query, email).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.PasswordHash,
		&u.IsVerified,
		&u.VerificationCode,
		&u.VerificationExpires,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *Repository) UpdateVerificationCode(userID string, code string, expires time.Time) error {
	query := `
		UPDATE users
		SET verification_code = $1, verification_expires = $2, updated_at = NOW()
		WHERE id = $3
	`
	_, err := r.db.Exec(query, code, expires, userID)
	return err
}

func (r *Repository) VerifyUser(userID string) error {
	query := `
		UPDATE users
		SET is_verified = TRUE, verification_code = NULL, verification_expires = NULL, updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.Exec(query, userID)
	return err
}