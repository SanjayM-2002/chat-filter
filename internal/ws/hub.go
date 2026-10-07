package ws

import "sync"

// Hub keeps track of which clients are in which room.
//
//	stream123 -> {A, B, C}
//	stream456 -> {D, E}
type Hub struct {
	mu    sync.RWMutex
	rooms map[string]map[*Client]struct{}
}

func NewHub() *Hub {
	return &Hub{rooms: make(map[string]map[*Client]struct{})}
}

// Join adds a client to its room.
func (h *Hub) Join(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.rooms[c.room] == nil {
		h.rooms[c.room] = make(map[*Client]struct{})
	}
	h.rooms[c.room][c] = struct{}{}
}

// Leave removes a client from its room and closes its mailbox,
// which tells the client's writer to stop.
func (h *Hub) Leave(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	clients := h.rooms[c.room]
	if _, ok := clients[c]; !ok {
		return // already left
	}
	delete(clients, c)
	close(c.send)
	if len(clients) == 0 {
		delete(h.rooms, c.room) // last one out, drop the room
	}
}

// Broadcast drops msg into the mailbox of every client in the room.
// A client whose mailbox is full is too slow; we disconnect it so it
// can't hold up everyone else.
func (h *Hub) Broadcast(room string, msg []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for c := range h.rooms[room] {
		select {
		case c.send <- msg:
		default:
			c.conn.Close() // reader will error out and call Leave
		}
	}
}

// RoomSize returns how many clients are in a room.
func (h *Hub) RoomSize(room string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.rooms[room])
}
