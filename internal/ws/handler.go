package ws

import (
	"log"
	"net/http"

	"github.com/SanjayM-2002/chat-filter/internal/filter"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Prototype: accept connections from any origin.
	// In production, check r.Header.Get("Origin") against allowed domains.
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Handler upgrades GET /ws?room=stream123&user=sanjay into a WebSocket.
func Handler(hub *Hub, f *filter.AbuseFilter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		room := r.URL.Query().Get("room")
		if room == "" {
			room = "lobby"
		}
		user := r.URL.Query().Get("user")
		if user == "" {
			user = "anon"
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("upgrade: %v", err)
			return // Upgrade already wrote the HTTP error
		}

		NewClient(hub, f, conn, room, user).Run()
		log.Printf("joined room=%s user=%s", room, user)
	}
}
