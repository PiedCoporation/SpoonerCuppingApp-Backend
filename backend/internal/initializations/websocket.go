package initializations

import (
	"backend/internal/presentations/ws/configs"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

var (
	webSocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
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
	if origin == "" {
		// Some clients (e.g., Postman) may omit Origin; allow in dev/local
		return true
	}

	u, err := url.Parse(origin)
	if err != nil {
		return false
	}

	// Allow localhost and 127.0.0.1 on any scheme/port
	switch u.Hostname() {
	case "localhost", "127.0.0.1":
		return true
	}

	return false
}
