package routers

import (
	"backend/internal/presentations/ws/v1/configs"
	"backend/internal/presentations/ws/v1/constants"
	wsControllers "backend/internal/presentations/ws/v1/controllers"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WSRouter struct{}

func (w *WSRouter) InitWS(router *gin.RouterGroup, hub *configs.Hub) {

    wsGroup := router.Group("/ws")
    wsGroup.GET("/", func(c *gin.Context) {
        ServeWS(hub, c.Writer, c.Request)
    })
    hub.RegisterHandler(constants.EventSendMessage, wsControllers.SendMessageEventController)
	hub.RegisterHandler(constants.EventCreateEvent, wsControllers.CreateEventEventController)
    hub.RegisterHandler(constants.EventJoinEvent, wsControllers.JoinEventEventController)
	hub.RegisterHandler(constants.EventLeaveEvent, wsControllers.LeaveEventEventController)
    hub.RegisterHandler(constants.EventStartEvent, wsControllers.StartEventEventController)
    hub.RegisterHandler(constants.EventEndEvent, wsControllers.EndEventEventController)
    hub.RegisterHandler(constants.EventMarkRound, wsControllers.MarkRoundEventController)

    go hub.InitialHub()
}

var (
	webSocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,  // Changed from 1024 to 1024 to match the code
		CheckOrigin:     checkOrigin,
	}
)

func ServeWS(hub *configs.Hub, w http.ResponseWriter, r *http.Request) {

	fmt.Println("New WebSocket connection")

	conn, err := webSocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to upgrade to WebSocket:", err)
		return
	}
	//defer conn.Close()

    client := configs.NewClient(hub, conn)
    hub.Register(client)

    // Client Start Read and Write
    go client.ReadMessages()
	go client.WriteMessages()
}

func checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	
	// Allow empty origin (Postman, curl, etc.)
	if origin == "" {
		return true
	}

	u, err := url.Parse(origin)
	if err != nil {
		return false
	}

	hostname := u.Hostname()
	
	// Allow localhost variants
	if hostname == "localhost" || hostname == "127.0.0.1" {
		return true
	}
	
	// Allow local network IPs (192.168.x.x, 10.x.x.x, etc.)
	if strings.HasPrefix(hostname, "192.168.") ||
	   strings.HasPrefix(hostname, "10.") ||
	   strings.HasPrefix(hostname, "172.16.") {
		return true
	}
	
	// For development: log and accept all origins
	log.Printf("⚠️  Accepting origin: %s (development mode)", origin)
	return true
}

