package project

type CreateProjectRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	UserID      string `json:"-"`
}

type UpdateProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Canvas      string `json:"canvas"`
}

type ShareProjectRequest struct {
	Email     string `json:"email" binding:"required"`
	Permission string `json:"permission" binding:"required,oneof=read edit"`
}

type Collaborator struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Permission string `json:"permission"`
}

type UpdateCollaboratorPermissionRequest struct {
	Permission string `json:"permission" binding:"required,oneof=read edit"`
}

// For the canvas, we can define a struct like this:
// id
// type
// x
// y
// width
// height
// rotation
// strokeColor
// fillColor
// strokeWidth
// opacity
// createdBy
// createdAt
// updatedAt