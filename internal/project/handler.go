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
	var req UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, response.StatusBadRequest, "Invalid request", nil, err.Error())
		return
	}

	updatedProject, err := h.projectService.UpdateProject(projectID, req)
	if err != nil {
		response.JSON(c, response.StatusInternalServerError, "Failed to update project", nil, err.Error())
		return
	}

	response.JSON(c, response.StatusOK, "Project updated successfully", updatedProject, nil)
}

func (h *Handler) DeleteProject(c *gin.Context) {
	projectID := c.Param("id")
	if err := h.projectService.DeleteProject(projectID); err != nil {
		response.JSON(c, response.StatusInternalServerError, "Failed to delete project", nil, err.Error())
		return
	}
	response.JSON(c, response.StatusOK, "Project deleted successfully", nil, nil)
}
