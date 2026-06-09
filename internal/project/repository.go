package project

import (
	"database/sql"
	"fmt"

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
		SELECT id, user_id, name, description, canvas, created_at, updated_at
		FROM projects
		WHERE id = $1
	`

	err := r.db.QueryRow(query, id).Scan(
		&project.ID,
		&project.OwnerID,
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
		SET name = $1, description = $2, canvas = CASE WHEN $3 = '' THEN canvas ELSE $3::jsonb END, updated_at = NOW()
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

func (r *Repository) GetProjectMembers(projectID string, requestingUserID string) (*MemberInfo, error) {
	var isMember bool
	checkQuery := `
		SELECT (
			SELECT COUNT(*) > 0 FROM projects WHERE id = $1 AND user_id = $2
		) OR (
			SELECT COUNT(*) > 0 FROM project_collaborators WHERE project_id = $1 AND user_id = $2
		)
	`
	err := r.db.QueryRow(checkQuery, projectID, requestingUserID).Scan(&isMember)
	if err != nil || !isMember {
		return nil, fmt.Errorf("access denied")
	}

	colQuery := `
		SELECT u.id, u.name, u.email, pc.permission
		FROM project_collaborators pc
		JOIN users u ON pc.user_id = u.id
		WHERE pc.project_id = $1
	`
	rows, err := r.db.Query(colQuery, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var collaborators []Collaborator
	for rows.Next() {
		var c Collaborator
		if err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.Permission); err != nil {
			return nil, err
		}
		collaborators = append(collaborators, c)
	}
	if collaborators == nil {
		collaborators = []Collaborator{}
	}

	var owner Collaborator
	owner.Permission = "owner"
	ownerQuery := `SELECT u.id, u.name, u.email FROM projects p JOIN users u ON p.user_id = u.id WHERE p.id = $1`
	_ = r.db.QueryRow(ownerQuery, projectID).Scan(&owner.ID, &owner.Name, &owner.Email)

	return &MemberInfo{
		TotalCount:    len(collaborators) + 1,
		Owner:         owner,
		Collaborators: collaborators,
	}, nil
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

func (r *Repository) IsProjectCollaborator(projectID string, userID string) (bool, error) {
	query := `
		SELECT COUNT(*) > 0
		FROM project_collaborators
		WHERE project_id = $1 AND user_id = $2
	`

	var isCollab bool
	err := r.db.QueryRow(query, projectID, userID).Scan(&isCollab)
	if err != nil {
		return false, err
	}
	return isCollab, nil
}

func (r *Repository) GetCollaboratorPermission(projectID string, userID string) (string, error) {
	var permission string
	query := `
		SELECT permission
		FROM project_collaborators
		WHERE project_id = $1 AND user_id = $2
	`
	err := r.db.QueryRow(query, projectID, userID).Scan(&permission)
	if err != nil {
		return "", err
	}
	return permission, nil
}
