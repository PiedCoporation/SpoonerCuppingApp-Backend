package constants

import (
	"encoding/json"
	"time"
)

const (
	EventSendMessage    = "send_message"
	EventNewMessage     = "new_message"
	EventCreateEvent    = "create_event"
	EventJoinEvent       = "join_event"
	EventNewJoinEvent    = "new_join_event"
	EventLeaveEvent      = "leave_event"
	EventNewLeaveEvent   = "new_leave_event"
	EventStartEvent      = "start_event"
	EventNewStartEvent   = "new_start_event"
	EventEndEvent        = "end_event"
	EventNewEndEvent     = "new_end_event"
	EventMarkRound       = "mark_round"
	EventNewMarkRound    = "new_mark_round"
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

type CreateEventEvent struct {
	EventId string `json:"event_id"`
}

type JoinEventEvent struct {
	EventId string `json:"event_id"`
}

type LeaveEventEvent struct {
	EventId string `json:"event_id"`
}

type StartEventEvent struct {
	EventId string `json:"event_id"`
}

type EndEventEvent struct {
	EventId string `json:"event_id"`
}

type NewJoinEventEvent struct {
	EventId string `json:"event_id"`
	UserId string `json:"user_id"`
	UserName string `json:"user_name"`
}

type NewLeaveEventEvent struct {
	EventId string `json:"event_id"`
	UserId string `json:"user_id"`
	UserName string `json:"user_name"`
}

type MarkRoundEvent struct {
	EventId string `json:"event_id"`
	Time string `json:"time"`
}