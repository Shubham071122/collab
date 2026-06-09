package auth

import (
	"database/sql"

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
) error {

	query := `
		INSERT INTO users (
			id,
			name,
			email,
			password_hash
		)
		VALUES ($1, $2, lower($3), $4)
	`

	_, err := r.db.Exec(
		query,
		uuid.New().String(),
		name,
		email,
		passwordHash,
	)

	return err
}

func (r *Repository) GetUserByEmail(email string) (*user.User, error) {

	query := `
		SELECT
			id,
			name,
			email,
			password_hash,
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
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &u, nil
}