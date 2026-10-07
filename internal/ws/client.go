package ws

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/SanjayM-2002/chat-filter/internal/filter"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second  // max time to write one message
	pongWait       = 60 * time.Second  // max time without hearing from the client
	pingPeriod     = pongWait * 9 / 10 // how often we ping (must be < pongWait)
	maxMessageSize = 1024              // max bytes in one incoming message
	mailboxSize    = 256               // buffered messages per client
)

// incoming is what the browser sends.
type incoming struct {
	Text string `json:"text"`
}

// outgoing is what we broadcast to the room.
type outgoing struct {
	User string `json:"user"`
	Text string `json:"text"`
	TS   int64  `json:"ts"`
}

// Client is one connected user (one open WebSocket).
type Client struct {
	hub    *Hub
	filter *filter.AbuseFilter
	conn   *websocket.Conn
	send   chan []byte // mailbox: messages waiting to go to this user
	room   string
	user   string
}

func NewClient(hub *Hub, f *filter.AbuseFilter, conn *websocket.Conn, room, user string) *Client {
	return &Client{
		hub:    hub,
		filter: f,
		conn:   conn,
		send:   make(chan []byte, mailboxSize),
		room:   room,
		user:   user,
	}
}

// Run joins the room and starts the reader and writer.
func (c *Client) Run() {
	c.hub.Join(c)
	go c.writePump()
	go c.readPump()
}

// readPump: messages FROM this user -> clean -> broadcast to the room.
func (c *Client) readPump() {
	defer func() {
		c.hub.Leave(c)
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("read error room=%s user=%s: %v", c.room, c.user, err)
			}
			return
		}

		var in incoming
		if err := json.Unmarshal(data, &in); err != nil {
			continue // ignore junk
		}
		text := strings.TrimSpace(in.Text)
		if text == "" {
			continue
		}

		clean := c.filter.Mask(text) // filter ONCE, before fan-out
		if clean != text {
			log.Printf("masked room=%s user=%s: %q -> %q", c.room, c.user, text, clean)
		}

		out, _ := json.Marshal(outgoing{User: c.user, Text: clean, TS: time.Now().Unix()})
		c.hub.Broadcast(c.room, out)
	}
}

// writePump: messages FOR this user (from the mailbox) -> down the socket.
// Also sends a ping every pingPeriod so we notice dead connections.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// hub closed the mailbox: say goodbye and stop
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
