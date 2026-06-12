package collaboration

import (
	"database/sql"
	"net/http"
	"strings"

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

func getProjectPermission(db *sql.DB, projectID string, userID string) (string, error) {
	var ownerID string
	ownerQuery := `SELECT user_id FROM projects WHERE id = $1`
	err := db.QueryRow(ownerQuery, projectID).Scan(&ownerID)
	if err != nil {
		return "", err
	}
	if ownerID == userID {
		return "owner", nil
	}

	var permission string
	collabQuery := `SELECT permission FROM project_collaborators WHERE project_id = $1 AND user_id = $2`
	err = db.QueryRow(collabQuery, projectID, userID).Scan(&permission)
	if err == sql.ErrNoRows {
		return "", nil
	} else if err != nil {
		return "", err
	}

	return permission, nil
}

func (h *Handler) Connect(c *gin.Context) {

	projectID := c.Param("id")

	userID := c.GetString("user_id")

	permission, err := getProjectPermission(h.Hub.db, projectID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify project permissions"})
		return
	}
	if permission == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: you do not have access to this project"})
		return
	}

	protocol := c.Request.Header.Get("Sec-WebSocket-Protocol")
	var responseHeader http.Header
	if protocol != "" {
		parts := strings.Split(protocol, ",")
		for _, part := range parts {
			trimmed := strings.TrimSpace(part)
			if trimmed != "" && trimmed != "Bearer" {
				responseHeader = make(http.Header)
				responseHeader.Set("Sec-WebSocket-Protocol", trimmed)
				break
			}
		}
	}

	conn, err := upgrader.Upgrade(
		c.Writer,
		c.Request,
		responseHeader,
	)

	if err != nil {
		return
	}

	client := &Client{
		UserID:     userID,
		ProjectID:  projectID,
		Permission: permission,
		Conn:       conn,
		Send:       make(chan []byte, 256),
		Hub:        h.Hub,
	}

	h.Hub.register <- client

	go client.writePump()
	go client.readPump()
}