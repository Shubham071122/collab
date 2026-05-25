package auth

import "github.com/gin-gonic/gin"

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Register(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "register route",
	})
}

func (h *Handler) Login(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "login route",
	})
}

func (h *Handler) User(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "current user",
	})
}

func (h *Handler) Logout(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "logout route",
	})
}