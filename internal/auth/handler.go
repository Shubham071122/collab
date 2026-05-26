package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/shubham071122/collab/internal/response"
)

type Handler struct {
	authService *Service
}

func NewHandler(authService *Service) *Handler {
	return &Handler{
		authService: authService,
	}
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, response.StatusBadRequest, "", nil, err.Error())
		return
	}

	if err := h.authService.Register(req); err != nil {
		response.JSON(c, response.StatusInternalServerError, "", nil, err.Error())
		return
	}

	response.JSON(c, response.StatusOK, "User registered successful", nil, nil)
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, response.StatusBadRequest, "", nil, err.Error())
		return
	}

	resp, err := h.authService.Login(req)
	if err != nil {
		response.JSON(c, response.StatusUnauthorized, "", nil, err.Error())
		return
	}

	response.JSON(c, response.StatusOK, "Login successful", resp, nil)
}

func (h *Handler) Logout(c *gin.Context) {
	if err := h.authService.Logout(); err != nil {
		response.JSON(c, response.StatusInternalServerError, "", nil, err.Error())
		return
	}

	response.JSON(c, response.StatusOK, "Logout successful", nil, nil)
}
