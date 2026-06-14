package project

import "time"

type Project struct {
	ID          string    `json:"id"`
	OwnerID     string    `json:"owner_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Canvas      string    `json:"canvas"`
	IsLocked    bool      `json:"is_locked"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type MemberInfo struct {
	TotalCount    int            `json:"total_count"`
	Owner         Collaborator   `json:"owner"`
	Collaborators []Collaborator `json:"collaborators"`
}