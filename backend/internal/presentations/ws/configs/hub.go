package configs

import (
	"backend/internal/presentations/ws/constants"
	"fmt"
	"sync"
)

type Hub struct {
    rooms map[string]map[*Client]bool

	register   chan *Client
	unregister chan *Client

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
        rooms:      make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		eventHandlers: make(map[string]constants.EventHandler),
	}

	h.setupEventHandlers()

	return h
}

func (h *Hub) setupEventHandlers() {

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

func (h *Hub) InitialHub() {
	for {
		select {

		case c := <-h.register:
			h.mu.Lock()
			if h.rooms[c.room] == nil {
				h.rooms[c.room] = make(map[*Client]bool)
			}
			h.rooms[c.room][c] = true
			h.mu.Unlock()

		case c := <-h.unregister:
			h.mu.Lock()
			if set, ok := h.rooms[c.room]; ok {
				if _, ok := set[c]; ok {
					delete(set, c)
					close(c.egress)
					if len(set) == 0 {
						delete(h.rooms, c.room)
					}
				}
			}
			h.mu.Unlock()

		}
	}
}

// Register adds a client to the hub via its internal channel.
func (h *Hub) Register(c *Client) {
    if h.closed {
        return
    }
    h.register <- c
}

// Unregister removes a client from the hub via its internal channel.
func (h *Hub) Unregister(c *Client) {
    if h.closed {
        return
    }
    h.unregister <- c
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