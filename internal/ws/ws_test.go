package ws

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/SanjayM-2002/chat-filter/internal/filter"

	"github.com/gorilla/websocket"
)

func dial(t *testing.T, srv *httptest.Server, room, user string) *websocket.Conn {
	t.Helper()
	url := "ws" + strings.TrimPrefix(srv.URL, "http") + "?room=" + room + "&user=" + user
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { conn.Close() })
	return conn
}

func waitForRoom(t *testing.T, hub *Hub, room string, n int) {
	t.Helper()
	for i := 0; i < 100; i++ {
		if hub.RoomSize(room) == n {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("room %s never reached %d clients", room, n)
}

func TestBroadcastMasksAndStaysInRoom(t *testing.T) {
	f := filter.New()
	f.AddWords([]string{"shit"})
	hub := NewHub()
	srv := httptest.NewServer(Handler(hub, f))
	defer srv.Close()

	a := dial(t, srv, "stream123", "alice")
	b := dial(t, srv, "stream123", "bob")
	other := dial(t, srv, "stream456", "eve")
	waitForRoom(t, hub, "stream123", 2)
	waitForRoom(t, hub, "stream456", 1)

	if err := a.WriteJSON(map[string]string{"text": "mondays are shit"}); err != nil {
		t.Fatal(err)
	}

	// both alice (sender) and bob get the masked message
	for _, c := range []*websocket.Conn{a, b} {
		c.SetReadDeadline(time.Now().Add(2 * time.Second))
		var got outgoing
		if err := c.ReadJSON(&got); err != nil {
			t.Fatal(err)
		}
		if got.Text != "mondays are ****" || got.User != "alice" {
			t.Errorf("got %+v", got)
		}
	}

	// eve is in another room and must get nothing
	other.SetReadDeadline(time.Now().Add(200 * time.Millisecond))
	if _, data, err := other.ReadMessage(); err == nil {
		t.Errorf("other room received %s", data)
	}
}

func TestLeaveRemovesClient(t *testing.T) {
	hub := NewHub()
	srv := httptest.NewServer(Handler(hub, filter.New()))
	defer srv.Close()

	a := dial(t, srv, "stream123", "alice")
	waitForRoom(t, hub, "stream123", 1)

	a.Close()
	waitForRoom(t, hub, "stream123", 0)
}

func TestOutgoingJSONShape(t *testing.T) {
	data, _ := json.Marshal(outgoing{User: "u", Text: "t", TS: 1})
	if string(data) != `{"user":"u","text":"t","ts":1}` {
		t.Errorf("got %s", data)
	}
}
