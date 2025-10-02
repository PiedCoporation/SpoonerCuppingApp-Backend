package configs

import (
	"backend/internal/presentations/ws/v1/constants"
	"fmt"
	"sync"
)

type Hub struct {
    rooms map[string]ClientLists
	clients ClientLists

	register   chan *Client
	unregister chan *Client
	createNewRoom chan *Client
	joinNewRoom chan *Client
	leaveNewRoom chan *Client

	mu     sync.RWMutex
	closed bool

	eventHandlers map[string]constants.EventHandler
}

type Message struct {
	Room string
	Data []byte
}

func NewHub() *Hub {
	h := &Hub{
        rooms:      make(map[string]ClientLists),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		createNewRoom: make(chan *Client),
		joinNewRoom: make(chan *Client),
		leaveNewRoom: make(chan *Client),
		eventHandlers: make(map[string]constants.EventHandler),
	}
	return h
}

func (h *Hub) routeEvent(event constants.Event, c *Client) error {
	if handler, ok := h.eventHandlers[event.Type]; ok {
        if err := handler(event, c); err != nil {
			return err
		}
		return nil
	} else {
		fmt.Println("No handler for event type:", event.Type)
	}
	return fmt.Errorf("no handler for event type: %s", event.Type)
}

// RegisterHandler registers a handler for a specific event type.
func (h *Hub) RegisterHandler(eventType string, handler constants.EventHandler) {
    h.eventHandlers[eventType] = handler
}

func (h *Hub) InitialHub() {
	for {
		select {

		case c := <-h.register:
			h.mu.Lock()
			h.clients[c] = true
			h.mu.Unlock()

		case c := <-h.unregister:
			h.mu.Lock()
			delete(h.clients, c)
			h.mu.Unlock()
		
		case c := <-h.createNewRoom:
			h.mu.Lock()
			h.rooms[c.room] = make(map[*Client]bool)
			h.rooms[c.room][c] = true
			h.mu.Unlock()

		case c := <-h.joinNewRoom:
			h.mu.Lock()
			if _, ok := h.rooms[c.room]; ok {
				h.rooms[c.room][c] = true
			}
			h.mu.Unlock()

		case c := <-h.leaveNewRoom:
			h.mu.Lock()
			if _, ok := h.rooms[c.room]; ok {
				delete(h.rooms[c.room], c)
			}
			h.mu.Unlock()

		}
	}
}

func (h *Hub) Broadcast(room string, event constants.Event) {
    h.mu.RLock()
    defer h.mu.RUnlock()
    if set, ok := h.rooms[room]; ok {
        for client := range set {
            if client.room == room {
                client.egress <- event
            }
        }
    }
}

func (h *Hub) Register(c *Client) {
    if h.closed {
        return
    }
    h.register <- c
}

func (h *Hub) CreateNewRoom(c *Client) {
    if h.closed {
        return
    }
    h.createNewRoom <- c
}

func (h *Hub) JoinNewRoom(c *Client) {
    if h.closed {
        return
    }
    h.joinNewRoom <- c
}

func (h *Hub) LeaveNewRoom(c *Client) {
    if h.closed {
        return
    }
    h.leaveNewRoom <- c
}

func (h *Hub) Unregister(c *Client) {
    if h.closed {
        return
    }
    h.unregister <- c
}

func (h *Hub) GetLeaderIdByRoom(room string) string {
    if _, ok := h.rooms[room]; ok {
        for client := range h.rooms[room] {
            return client.LeaderID()
        }
    }
    return ""
}

func (h *Hub) Close() {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.closed {
		return
	}
	h.closed = true
    for _, set := range h.rooms {
		for c := range set {
			c.Close()
		}
	}
}