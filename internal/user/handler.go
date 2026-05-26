package user

import (
	"github.com/gin-gonic/gin"
	"github.com/shubham071122/collab/internal/response"
)

type Handler struct {
	userService *Service
}

func NewHandler(userService *Service) *Handler {
	return &Handler{
		userService: userService,
	}
}

func (h *Handler) CurrentUser(c *gin.Context) {
	userID := c.Param("id")
	user, err := h.userService.CurrentUser(userID)
	if err != nil {
		response.JSON(c, response.StatusNotFound, "", nil, err.Error())
		return
	}
	response.JSON(c, response.StatusOK, "User retrieved successfully", user, nil)
}
