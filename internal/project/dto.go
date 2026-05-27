package project

type CreateProjectRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	UserID      string `json:"-"`
}

type UpdateProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
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