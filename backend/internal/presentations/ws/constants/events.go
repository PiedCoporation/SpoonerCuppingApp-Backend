package constants

import (
	"encoding/json"
	"time"
)

const (
	EventSendMessage    = "send_message"
	EventNewMessage     = "new_message"
	EventChangeChatRoom = "change_room"
)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// EventHandler is kept decoupled from the websocket client type to avoid import cycles.
// The handler receives the event and an opaque client reference.
type EventHandler func(event Event, client interface{}) error

type SendMessageEvent struct {
	Message string `json:"message"`
	From    string `json:"from"`
}

type NewMessageEvent struct {
	SendMessageEvent
	Sent time.Time `json:"sent"`
}

type ChangeChatRoomEvent struct {
	Name string `json:"room"`
}