package configs

import (
	"backend/internal/presentations/ws/v1/constants"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

var (
	pongWait     = 10 * time.Second
	pingInterval = (pongWait * 9) / 10
)

type ClientLists map[*Client]bool

type Client struct {
	Id string
	hub  *Hub
	conn *websocket.Conn

	room string
	isLeader bool
	leaderID string

	egress chan constants.Event
}

// NewClient creates a new websocket client with initialized fields.
func NewClient(hub *Hub, conn *websocket.Conn, user_id string) *Client {
	return &Client{
		Id: user_id,
		hub:  hub,
		conn: conn,
		room: "lobby", // Default room for new connections
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
			fmt.Printf("ping client_id=%s time=%s\n", c.Id, time.Now().Format(time.RFC3339))

			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
		}

	}
}


func (c *Client) pongHandler(pongMsg string) error {
	fmt.Printf("pong client_id=%s time=%s\n", c.Id, time.Now().Format(time.RFC3339))
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

func (c *Client) Hub() *Hub {
    return c.hub
}

func (c *Client) Room() string {
    return c.room
}

func (c *Client) SetRoom(room string) {
    c.room = room
}

func (c *Client) Egress() chan constants.Event {
    return c.egress
}

func (c *Client) IsLeader() bool {
    return c.isLeader
}

func (c *Client) LeaderID() string {
    return c.leaderID
}

func (c *Client) SetLeaderId(leaderID string) {
    c.leaderID = leaderID
}

func (c *Client) SetIsLeader(isLeader bool) {
    c.isLeader = isLeader
}
