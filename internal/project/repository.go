package project

import (
	"database/sql"

	"github.com/google/uuid"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetProjectByID(id string) (*Project, error) {
	var project Project

	query := `
		SELECT id, name, description, created_at, updated_at
		FROM projects
		WHERE id = $1
	`

	err := r.db.QueryRow(query, id).Scan(
		&project.ID,
		&project.Name,
		&project.Description,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &project, nil
}

func (r *Repository) CreateProject(name string, description string, userId string) (*Project, error) {
	var project Project

	query := `
		INSERT INTO projects (
			id,
			name,
			description,
			user_id
		)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, description, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		uuid.New().String(),
		name,
		description,
		userId,
	).Scan(
		&project.ID,
		&project.Name,
		&project.Description,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &project, nil
}

func (r *Repository) UpdateProject(project *Project) (*Project, error) {

	query := `
		UPDATE projects
		SET name = $1, description = $2, updated_at = NOW()
		WHERE id = $3
	`

	_, err := r.db.Exec(
		query,
		project.Name,
		project.Description,
		project.ID,
	)

	if err != nil {
		return nil, err
	}

	return project, nil
}

func (r *Repository) DeleteProject(id string) error {

	query := `
		DELETE FROM projects
		WHERE id = $1
	`

	_, err := r.db.Exec(query, id)

	return err
}
