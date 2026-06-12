package user

import (
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetUserByID(id string) (*User, error) {
	var user User
	err := r.db.QueryRow(`
		SELECT id, name, email, is_verified, verification_code, verification_expires, created_at, updated_at 
		FROM users 
		WHERE id = $1
	`, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.IsVerified,
		&user.VerificationCode,
		&user.VerificationExpires,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserByEmail(email string) (*User, error) {
	var user User
	err := r.db.QueryRow(`
		SELECT id, name, email, is_verified, verification_code, verification_expires, created_at, updated_at 
		FROM users 
		WHERE email = $1
	`, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.IsVerified,
		&user.VerificationCode,
		&user.VerificationExpires,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
