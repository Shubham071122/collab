package collaboration

import "encoding/json"

type Hub struct {
	rooms      map[string]map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan Message
}

func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan Message),
	}
}

func (h *Hub) Run() {
	for {
		select {

		case client := <-h.register:
			if _, exists := h.rooms[client.ProjectID]; !exists {
				h.rooms[client.ProjectID] = make(map[*Client]bool)
			}
			h.rooms[client.ProjectID][client] = true

		case client := <-h.unregister:
			if clients, exists := h.rooms[client.ProjectID]; exists {
				delete(clients, client)
				close(client.Send)

				if len(clients) == 0 {
					delete(h.rooms, client.ProjectID)
				}
			}

		case message := <-h.broadcast:
			if clients, exists := h.rooms[message.ProjectID]; exists {
				for client := range clients {
					select {
					case client.Send <- mustMarshal(message):
						
					default:
						close(client.Send)
						delete(clients, client)
					}
				}
			}
		}
	}
}

func mustMarshal(message Message) []byte {
	data, _ := json.Marshal(message)
	return data
}