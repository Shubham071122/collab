package collaboration

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Handler struct {
	Hub *Hub
}

func NewHandler(hub *Hub) *Handler {
	return &Handler{
		Hub: hub,
	}
}

func (h *Handler) Connect(c *gin.Context) {

	projectID := c.Param("id")

	userID := c.GetString("user_id")

	conn, err := upgrader.Upgrade(
		c.Writer,
		c.Request,
		nil,
	)

	if err != nil {
		return
	}

	client := &Client{
		UserID:    userID,
		ProjectID: projectID,
		Conn:      conn,
		Send:      make(chan []byte, 256),
		Hub:       h.Hub,
	}

	h.Hub.register <- client

	go client.writePump()
	go client.readPump()
}