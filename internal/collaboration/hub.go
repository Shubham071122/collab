package collaboration

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type CanvasSnapshot struct {
	Store  map[string]interface{} `json:"store"`
	Schema interface{}            `json:"schema"`
}

type permissionUpdateTask struct {
	ProjectID  string
	UserID     string
	Permission string
}

type Hub struct {
	rooms            map[string]map[*Client]bool
	register         chan *Client
	unregister       chan *Client
	broadcast        chan Message
	db               *sql.DB
	roomsCanvas      map[string]*CanvasSnapshot
	saveTimers       map[string]*time.Timer
	saveQueue        chan string
	permissionUpdate chan permissionUpdateTask
}

func NewHub(db *sql.DB) *Hub {
	return &Hub{
		rooms:            make(map[string]map[*Client]bool),
		register:         make(chan *Client),
		unregister:       make(chan *Client),
		broadcast:        make(chan Message),
		db:               db,
		roomsCanvas:      make(map[string]*CanvasSnapshot),
		saveTimers:       make(map[string]*time.Timer),
		saveQueue:        make(chan string, 100),
		permissionUpdate: make(chan permissionUpdateTask, 100),
	}
}

func (h *Hub) Run() {
	for {
		select {

		case client := <-h.register:
			if _, exists := h.rooms[client.ProjectID]; !exists {
				h.rooms[client.ProjectID] = make(map[*Client]bool)

				snapshot, err := h.loadCanvasFromDB(client.ProjectID)
				if err != nil {
					log.Printf("Error loading initial canvas from DB for project %s: %v", client.ProjectID, err)
					snapshot = &CanvasSnapshot{Store: make(map[string]interface{})}
				}
				h.roomsCanvas[client.ProjectID] = snapshot
			}
			h.rooms[client.ProjectID][client] = true

		case client := <-h.unregister:
			if clients, exists := h.rooms[client.ProjectID]; exists {
				delete(clients, client)
				close(client.Send)

				if len(clients) == 0 {
					delete(h.rooms, client.ProjectID)

					if timer, exists := h.saveTimers[client.ProjectID]; exists {
						timer.Stop()
						delete(h.saveTimers, client.ProjectID)
					}
					if snapshot, exists := h.roomsCanvas[client.ProjectID]; exists && snapshot != nil {
						delete(h.roomsCanvas, client.ProjectID)
						go h.saveCanvasToDB(client.ProjectID, snapshot)
					}
				}
			}

		case message := <-h.broadcast:
			if clients, exists := h.rooms[message.ProjectID]; exists {
				marshaledMsg := mustMarshal(message)
				for client := range clients {
					select {
					case client.Send <- marshaledMsg:
					default:
						close(client.Send)
						delete(clients, client)
					}
				}
			}

			if message.Type == "canvas_change" {
				projectId := message.ProjectID
				if snapshot, exists := h.roomsCanvas[projectId]; exists && snapshot != nil {
					h.applyCanvasDiff(snapshot, message.Payload)

					if timer, exists := h.saveTimers[projectId]; exists {
						timer.Stop()
					}
					h.saveTimers[projectId] = time.AfterFunc(10*time.Second, func() {
						h.saveQueue <- projectId
					})
				}
			}

		case projectId := <-h.saveQueue:
			snapshot, exists := h.roomsCanvas[projectId]
			if exists && snapshot != nil {
				go h.saveCanvasToDB(projectId, snapshot)
			}

		case task := <-h.permissionUpdate:
			log.Printf("Hub permissionUpdate received: projectID=%s, userID=%s, permission=%s", task.ProjectID, task.UserID, task.Permission)
			if clients, exists := h.rooms[task.ProjectID]; exists {
				log.Printf("Found room %s with %d clients", task.ProjectID, len(clients))
				for client := range clients {
					log.Printf("Client: userID=%s, permission=%s", client.UserID, client.Permission)
					if client.UserID == task.UserID {
						log.Printf("Matching client found! Updating to %s", task.Permission)
						if task.Permission == "" {
							msg := Message{
								Type:      "access_revoked",
								ProjectID: task.ProjectID,
								UserID:    task.UserID,
							}
							client.Send <- mustMarshal(msg)
							go func(conn *websocket.Conn) {
								time.Sleep(150 * time.Millisecond)
								conn.Close()
							}(client.Conn)
						} else {
							// Permission updated!
							client.Permission = task.Permission
							msg := Message{
								Type:      "permission_changed",
								ProjectID: task.ProjectID,
								UserID:    task.UserID,
								Payload:   json.RawMessage(`{"permission":"` + task.Permission + `"}`),
							}
							client.Send <- mustMarshal(msg)
						}
					}
				}
			} else {
				log.Printf("Room %s not found in active rooms", task.ProjectID)
			}
		}
	}
}

func (h *Hub) loadCanvasFromDB(projectId string) (*CanvasSnapshot, error) {
	var canvasJSON string
	query := `SELECT canvas FROM projects WHERE id = $1`
	err := h.db.QueryRow(query, projectId).Scan(&canvasJSON)
	if err != nil {
		return nil, err
	}

	var snapshot CanvasSnapshot
	err = json.Unmarshal([]byte(canvasJSON), &snapshot)
	if err != nil {
		snapshot.Store = make(map[string]interface{})
	}
	if snapshot.Store == nil {
		snapshot.Store = make(map[string]interface{})
	}
	return &snapshot, nil
}

func (h *Hub) applyCanvasDiff(snapshot *CanvasSnapshot, payload []byte) {
	var changes struct {
		Added   map[string]interface{}    `json:"added"`
		Updated map[string][2]interface{} `json:"updated"`
		Removed map[string]interface{}    `json:"removed"`
	}
	err := json.Unmarshal(payload, &changes)
	if err != nil {
		log.Printf("Error unmarshaling canvas changes: %v", err)
		return
	}

	for id, record := range changes.Added {
		snapshot.Store[id] = record
	}
	for id, change := range changes.Updated {
		snapshot.Store[id] = change[1]
	}
	for id := range changes.Removed {
		delete(snapshot.Store, id)
	}
}

func (h *Hub) saveCanvasToDB(projectId string, snapshot *CanvasSnapshot) {
	canvasData, err := json.Marshal(snapshot)
	if err != nil {
		log.Printf("Error marshaling canvas snapshot for project %s: %v", projectId, err)
		return
	}

	query := `
		UPDATE projects
		SET canvas = $1::jsonb, updated_at = NOW()
		WHERE id = $2
	`
	_, err = h.db.Exec(query, string(canvasData), projectId)
	if err != nil {
		log.Printf("Error auto-saving canvas to database for project %s: %v", projectId, err)
	} else {
		log.Printf("Successfully auto-saved canvas to database for project %s", projectId)
		select {
		case h.broadcast <- Message{
			Type:      "canvas_saved",
			ProjectID: projectId,
		}:
		case <-time.After(1 * time.Second):
			log.Printf("Timeout broadcasting canvas_saved for project %s", projectId)
		}
	}
}

func mustMarshal(message Message) []byte {
	data, _ := json.Marshal(message)
	return data
}

func (h *Hub) UpdateUserPermission(projectID string, userID string, permission string) {
	h.permissionUpdate <- permissionUpdateTask{
		ProjectID:  projectID,
		UserID:     userID,
		Permission: permission,
	}
}