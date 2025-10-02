package usecases

import (
	"backend/internal/presentations/ws/v1/configs"
	"backend/internal/presentations/ws/v1/constants"
	"encoding/json"
	"fmt"
	"time"
)

func SendMessageWS(event constants.Event, client *configs.Client) error {
    var chatevent constants.SendMessageEvent
    if err := json.Unmarshal(event.Payload, &chatevent); err != nil {
        return fmt.Errorf("could not unmarshal payload: %w", err)
    }

    var broadcastEvent constants.NewMessageEvent
    broadcastEvent.Message = chatevent.Message
    broadcastEvent.From = chatevent.From
    broadcastEvent.Sent = time.Now()

    data, err := json.Marshal(broadcastEvent)
    if err != nil {
        return fmt.Errorf("could not marshal payload: %w", err)
    }

    outgoingEvent := constants.Event{
        Type:    constants.EventNewMessage,
        Payload: data,
    }

    client.Hub().Broadcast(client.Room(), outgoingEvent)
    return nil
}

func CreateEventWS(event constants.Event, client *configs.Client) error {
    var createEventEvent constants.CreateEventEvent
    if err := json.Unmarshal(event.Payload, &createEventEvent); err != nil {
        return fmt.Errorf("could not unmarshal payload: %w", err)
    }
    client.SetRoom(createEventEvent.EventId)
    client.Hub().CreateNewRoom(client)
    client.SetIsLeader(true)
    client.SetLeaderId(client.Id)
    return nil
}

func JoinEventWS(event constants.Event, client *configs.Client) error {
    var chatRoomEvent constants.JoinEventEvent
    if err := json.Unmarshal(event.Payload, &chatRoomEvent); err != nil {
        return fmt.Errorf("could not unmarshal payload: %w", err)
    }

    leaderId := client.Hub().GetLeaderIdByRoom(chatRoomEvent.EventId)
    if leaderId == "" {
        return fmt.Errorf("could not get leader id")
    }

    client.SetLeaderId(leaderId)
    client.SetRoom(chatRoomEvent.EventId)
    client.SetIsLeader(false)
    client.Hub().JoinNewRoom(client)

    broadcastEvent := constants.NewJoinEventEvent{
        EventId: chatRoomEvent.EventId,
        UserId:  client.Id,
        UserName: "test",
    }

    data, err := json.Marshal(broadcastEvent)
    if err != nil {
        return fmt.Errorf("could not marshal payload: %w", err)
    }

    outgoingEvent := constants.Event{
        Type:    constants.EventNewJoinEvent,
        Payload: data,
    }

    client.Hub().Broadcast(chatRoomEvent.EventId, outgoingEvent)
    return nil
}

func LeaveEventWS(event constants.Event, client *configs.Client) error {
    var chatRoomEvent constants.LeaveEventEvent
    if err := json.Unmarshal(event.Payload, &chatRoomEvent); err != nil {
        return fmt.Errorf("could not unmarshal payload: %w", err)
    }

    client.Hub().LeaveNewRoom(client)
    client.SetIsLeader(false)
    client.SetLeaderId("")
    client.SetRoom("")

    broadcastEvent := constants.NewLeaveEventEvent{
        EventId: chatRoomEvent.EventId,
        UserId:  client.Id,
        UserName: "test",
    }

    data, err := json.Marshal(broadcastEvent)
    if err != nil {
        return fmt.Errorf("could not marshal payload: %w", err)
    }

    outgoingEvent := constants.Event{
        Type:    constants.EventNewLeaveEvent,
        Payload: data,
    }

    client.Hub().Broadcast(chatRoomEvent.EventId, outgoingEvent)
    return nil
}

func StartEventWS(event constants.Event, client *configs.Client) error {
    var chatRoomEvent constants.StartEventEvent
    if err := json.Unmarshal(event.Payload, &chatRoomEvent); err != nil {
        return fmt.Errorf("could not unmarshal payload: %w", err)
    }

    broadcastEvent := constants.StartEventEvent{
        EventId: chatRoomEvent.EventId,
    }

    data, err := json.Marshal(broadcastEvent)
    if err != nil {
        return fmt.Errorf("could not marshal payload: %w", err)
    }

    outgoingEvent := constants.Event{
        Type:    constants.EventNewStartEvent,
        Payload: data,
    }

    client.Hub().Broadcast(chatRoomEvent.EventId, outgoingEvent)
    return nil
}

func EndEventWS(event constants.Event, client *configs.Client) error {
    var chatRoomEvent constants.EndEventEvent
    if err := json.Unmarshal(event.Payload, &chatRoomEvent); err != nil {
        return fmt.Errorf("could not unmarshal payload: %w", err)
    }

    broadcastEvent := constants.EndEventEvent{
        EventId: chatRoomEvent.EventId,
    }

    data, err := json.Marshal(broadcastEvent)
    if err != nil {
        return fmt.Errorf("could not marshal payload: %w", err)
    }

    outgoingEvent := constants.Event{
        Type:    constants.EventNewEndEvent,
        Payload: data,
    }

    client.Hub().Broadcast(chatRoomEvent.EventId, outgoingEvent)
    return nil
}

func MarkRoundWS(event constants.Event, client *configs.Client) error {
    var chatRoomEvent constants.MarkRoundEvent
    if err := json.Unmarshal(event.Payload, &chatRoomEvent); err != nil {
        return fmt.Errorf("could not unmarshal payload: %w", err)
    }

    broadcastEvent := constants.MarkRoundEvent{
        EventId: chatRoomEvent.EventId,
        Time:    chatRoomEvent.Time,
    }

    data, err := json.Marshal(broadcastEvent)
    if err != nil {
        return fmt.Errorf("could not marshal payload: %w", err)
    }

    outgoingEvent := constants.Event{
        Type:    constants.EventNewMarkRound,
        Payload: data,
    }

    client.Hub().Broadcast(chatRoomEvent.EventId, outgoingEvent)
    return nil
}


