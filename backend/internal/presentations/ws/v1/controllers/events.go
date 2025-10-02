package controllers

import (
	"backend/internal/presentations/ws/v1/configs"
	"backend/internal/presentations/ws/v1/constants"
	"backend/internal/usecases"
	"fmt"
)

func SendMessageEventController(event constants.Event, client interface{}) error {
    c, ok := client.(*configs.Client)
    if !ok {
        return fmt.Errorf("invalid client type")
    }
    return usecases.SendMessageWS(event, c)
}

func CreateEventEventController(event constants.Event, client interface{}) error {
    c, ok := client.(*configs.Client)
    if !ok {
        return fmt.Errorf("invalid client type")
    }
    return usecases.CreateEventWS(event, c)
}

func JoinEventEventController(event constants.Event, client interface{}) error {
    c, ok := client.(*configs.Client)
    if !ok {
        return fmt.Errorf("invalid client type")
    }
    return usecases.JoinEventWS(event, c)
}

func LeaveEventEventController(event constants.Event, client interface{}) error {
    c, ok := client.(*configs.Client)
    if !ok {
        return fmt.Errorf("invalid client type")
    }
    return usecases.LeaveEventWS(event, c)
}

func StartEventEventController(event constants.Event, client interface{}) error {
    c, ok := client.(*configs.Client)
    if !ok {
        return fmt.Errorf("invalid client type")
    }
    return usecases.StartEventWS(event, c)
}

func EndEventEventController(event constants.Event, client interface{}) error {
    c, ok := client.(*configs.Client)
    if !ok {
        return fmt.Errorf("invalid client type")
    }
    return usecases.EndEventWS(event, c)
}

func MarkRoundEventController(event constants.Event, client interface{}) error {
    c, ok := client.(*configs.Client)
    if !ok {
        return fmt.Errorf("invalid client type")
    }
    return usecases.MarkRoundWS(event, c)
}


