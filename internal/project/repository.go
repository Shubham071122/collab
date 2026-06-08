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
		SELECT id, name, description, canvas, created_at, updated_at
		FROM projects
		WHERE id = $1
	`

	err := r.db.QueryRow(query, id).Scan(
		&project.ID,
		&project.Name,
		&project.Description,
		&project.Canvas,
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
		RETURNING id, name, description, canvas, created_at, updated_at
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
		&project.Canvas,
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
		SET name = $1, description = $2, canvas = $3, updated_at = NOW()
		WHERE id = $4
	`

	_, err := r.db.Exec(
		query,
		project.Name,
		project.Description,
		project.Canvas,
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

func (r *Repository) ShareProject(projectID string, userID string, permission string) error {

	query := `
		INSERT INTO project_collaborators (id, project_id, user_id, permission)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (project_id, user_id) DO UPDATE SET permission = EXCLUDED.permission
	`

	_, err := r.db.Exec(query, uuid.New().String(), projectID, userID, permission)

	return err
}

func (r *Repository) GetCollaborators(projectID string) ([]Collaborator, error) {
	var collaborators []Collaborator

	query := `
		SELECT u.id, u.name, u.email, pc.permission
		FROM project_collaborators pc
		JOIN users u ON pc.user_id = u.id
		WHERE pc.project_id = $1
	`

	rows, err := r.db.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var collaborator Collaborator
		if err := rows.Scan(&collaborator.ID, &collaborator.Name, &collaborator.Email, &collaborator.Permission); err != nil {
			return nil, err
		}
		collaborators = append(collaborators, collaborator)
	}

	return collaborators, nil
}

func (r *Repository) UpdateCollaboratorPermission(projectID string, userID string, permission string) error {

	query := `
		UPDATE project_collaborators
		SET permission = $1
		WHERE project_id = $2 AND user_id = $3
	`

	_, err := r.db.Exec(query, permission, projectID, userID)

	return err
}

func (r *Repository) RemoveCollaborator(projectID string, userID string) error {

	query := `
		DELETE FROM project_collaborators
		WHERE project_id = $1 AND user_id = $2
	`

	_, err := r.db.Exec(query, projectID, userID)

	return err
}

func (r *Repository) IsProjectOwner(
	projectID string,
	userID string,
) (bool, error) {

	query := `
		SELECT COUNT(*) > 0
		FROM projects
		WHERE id = $1 AND user_id = $2
	`

	var isOwner bool
	err := r.db.QueryRow(query, projectID, userID).Scan(&isOwner)

	if err != nil {
		return false, err
	}

	return isOwner, nil
}