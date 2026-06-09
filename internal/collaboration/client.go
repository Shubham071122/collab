package collaboration

import "github.com/gorilla/websocket"

type Client struct {
	UserID     string
	ProjectID  string
	Permission string

	Conn *websocket.Conn
	Send chan []byte

	Hub *Hub
}

func (c *Client) readPump() {

	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	for {

		var message Message

		err := c.Conn.ReadJSON(&message)

		if err != nil {
			break
		}

		message.ProjectID = c.ProjectID
		message.UserID = c.UserID

		c.Hub.broadcast <- message
	}
}

func (c *Client) writePump() {

	defer c.Conn.Close()

	for {

		message, ok := <-c.Send

		if !ok {

			c.Conn.WriteMessage(
				websocket.CloseMessage,
				[]byte{},
			)

			return
		}

		err := c.Conn.WriteMessage(
			websocket.TextMessage,
			message,
		)

		if err != nil {
			return
		}
	}
}