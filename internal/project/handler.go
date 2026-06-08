package project

import (
	"github.com/gin-gonic/gin"
	"github.com/shubham071122/collab/internal/response"
)

type Handler struct {
	projectService *Service
}

func NewHandler(projectService *Service) *Handler {
	return &Handler{
		projectService: projectService,
	}
}

func (h *Handler) GetProject(c *gin.Context) {
	projectID := c.Param("id")
	project, err := h.projectService.GetProjectByID(projectID)
	if err != nil {
		response.JSON(c, response.StatusNotFound, "Project not found", nil, nil)
		return
	}
	response.JSON(c, response.StatusOK, "Success", project, nil)
}

func (h *Handler) CreateProject(c *gin.Context) {
	var req CreateProjectRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, response.StatusBadRequest, "Invalid request", nil, err.Error())
		return
	}

	req.UserID = c.GetString("user_id")

	project, err := h.projectService.CreateProject(req)
	if err != nil {
		response.JSON(c, response.StatusInternalServerError, "Failed to create project", nil, err.Error())
		return
	}

	response.JSON(c, response.StatusOK, "Project created successfully", project, nil)
}

func (h *Handler) UpdateProject(c *gin.Context) {
	projectID := c.Param("id")
	userID := c.GetString("user_id")
	var req UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, response.StatusBadRequest, "Invalid request", nil, err.Error())
		return
	}

	updatedProject, err := h.projectService.UpdateProject(projectID, userID, req)
	if err != nil {
		response.JSON(c, response.StatusInternalServerError, "Failed to update project", nil, err.Error())
		return
	}

	response.JSON(c, response.StatusOK, "Project updated successfully", updatedProject, nil)
}

func (h *Handler) DeleteProject(c *gin.Context) {
	projectID := c.Param("id")
	userID := c.GetString("user_id")
	
	if err := h.projectService.DeleteProject(projectID, userID); err != nil {
		response.JSON(c, response.StatusInternalServerError, "Failed to delete project", nil, err.Error())
		return
	}
	response.JSON(c, response.StatusOK, "Project deleted successfully", nil, nil)
}

func (h *Handler) ShareProject(c *gin.Context) {
	projectID := c.Param("id")
	var req ShareProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, response.StatusBadRequest, "Invalid request", nil, err.Error())
		return
	}

	ownerID := c.GetString("user_id")
	if err := h.projectService.ShareProject(projectID, req, ownerID); err != nil {
		response.JSON(c, response.StatusInternalServerError, "Failed to share project", nil, err.Error())
		return
	}

	response.JSON(c, response.StatusOK, "Project shared successfully", nil, nil)
}

func (h *Handler) GetCollaborators(c *gin.Context) {
	projectID := c.Param("id")
	userID := c.GetString("user_id")
	collaborators, err := h.projectService.GetCollaborators(projectID, userID)
	if err != nil {
		response.JSON(c, response.StatusInternalServerError, "Failed to get collaborators", nil, err.Error())
		return
	}
	response.JSON(c, response.StatusOK, "Success", collaborators, nil)
}

func (h *Handler) UpdateCollaboratorPermission(c *gin.Context) {
	projectID := c.Param("id")
	userID := c.Param("userId")
	ownerID := c.GetString("user_id")

	var req UpdateCollaboratorPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, response.StatusBadRequest, "Invalid request", nil, err.Error())
		return
	}

	if err := h.projectService.UpdateCollaboratorPermission(projectID, userID, ownerID, req.Permission); err != nil {
		response.JSON(c, response.StatusInternalServerError, "Failed to update collaborator permission", nil, err.Error())
		return
	}

	response.JSON(c, response.StatusOK, "Collaborator permission updated successfully", nil, nil)
}

func (h *Handler) RemoveCollaborator(c *gin.Context) {
	projectID := c.Param("id")
	userID := c.Param("userId")
	ownerID := c.GetString("user_id")

	if err := h.projectService.RemoveCollaborator(projectID, userID, ownerID); err != nil {
		response.JSON(c, response.StatusInternalServerError, "Failed to remove collaborator", nil, err.Error())
		return
	}

	response.JSON(c, response.StatusOK, "Collaborator removed successfully", nil, nil)
}
