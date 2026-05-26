package response

import "github.com/gin-gonic/gin"

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func JSON(
	c *gin.Context,
	status int,
	message string,
	data interface{},
	err interface{},
) {

	if message == "" {
		message = DefaultMessages[status]
	}

	c.JSON(status, APIResponse{
		Success: status < 400,
		Message: message,
		Data:    data,
		Error:   err,
	})
}