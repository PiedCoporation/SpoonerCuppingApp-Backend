package configs

import (
	"backend/internal/presentations/ws/constants"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

var (
	pongWait     = 10 * time.Second
	pingInterval = (pongWait * 9) / 10
)

type Client struct {
	hub  *Hub
	conn *websocket.Conn

	room string
	egress chan constants.Event
}

// NewClient creates a new websocket client with initialized fields.
func NewClient(hub *Hub, conn *websocket.Conn) *Client {
	return &Client{
		hub:  hub,
		conn: conn,
		egress:     make(chan constants.Event),
	}
}

func (c *Client) ReadMessages() {
	defer func() {
		// clean up connection
		c.hub.Unregister(c)
	}()

	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	c.conn.SetReadLimit(512)

	c.conn.SetPongHandler(c.pongHandler)

	for {
		_, payload, err := c.conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("Error: %v\n", err)
			}
			break
		}

		var request constants.Event

		if err := json.Unmarshal(payload, &request); err != nil {
			fmt.Printf("Error unmarshaling message: %v\n", err)
			continue
		}

		if err := c.hub.routeEvent(request, c); err != nil {
			fmt.Printf("Error routing event: %v\n", err)
		}
	}
}

func (c *Client) WriteMessages() {
	defer func() {
		// clean up connection
		c.hub.Unregister(c)
	}()

	ticker := time.NewTicker(pingInterval)

	for {
		select {

		case message, ok := <-c.egress:
			if !ok {
				if err := c.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					fmt.Printf("Error: %v\n", err)
				}
				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				fmt.Printf("Error marshaling message: %v\n", err)
				continue
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
			fmt.Printf("Message sent")

		case <-ticker.C:
			fmt.Printf("ping \n")

			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
		}

	}
}

func (c *Client) pongHandler(pongMsg string) error {
	fmt.Printf("pong \n")
	return c.conn.SetReadDeadline(time.Now().Add(pongWait))
}


func (c *Client) Close() {
	// trigger writer to exit
	select {
	case <-c.egress:
	default:
	}
	close(c.egress)
	_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
}

