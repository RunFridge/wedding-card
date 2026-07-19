package handlers

import (
	"encoding/json"
	"testing"
	"time"
)

func newTestHub() *Hub {
	return &Hub{
		visitors: make(map[*client]bool),
		admins:   make(map[*client]bool),
		buffer:   make(map[int64]int64),
		configFn: func() (int, int) { return 2000, 100 },
		done:     make(chan struct{}),
	}
}

func newTestClient(h *Hub, role clientRole) *client {
	return &client{
		hub:  h,
		role: role,
		send: make(chan []byte, 64),
	}
}

func drainMessages(c *client) []map[string]any {
	var msgs []map[string]any
	for {
		select {
		case raw := <-c.send:
			var m map[string]any
			json.Unmarshal(raw, &m)
			msgs = append(msgs, m)
		default:
			return msgs
		}
	}
}

func TestConnectionCountVisitorsOnly(t *testing.T) {
	h := newTestHub()

	v1 := newTestClient(h, roleVisitor)
	v2 := newTestClient(h, roleVisitor)
	a1 := newTestClient(h, roleAdmin)

	h.mu.Lock()
	h.visitors[v1] = true
	h.visitors[v2] = true
	h.admins[a1] = true
	h.mu.Unlock()

	if got := h.ConnectionCount(); got != 2 {
		t.Errorf("ConnectionCount() = %d, want 2", got)
	}
}

func TestRegisterVisitor(t *testing.T) {
	h := newTestHub()
	// Add an admin to receive presence broadcasts
	admin := newTestClient(h, roleAdmin)
	h.mu.Lock()
	h.admins[admin] = true
	h.mu.Unlock()

	v := newTestClient(h, roleVisitor)
	h.register(v)

	h.mu.Lock()
	_, inVisitors := h.visitors[v]
	_, inAdmins := h.admins[v]
	h.mu.Unlock()

	if !inVisitors {
		t.Error("visitor not in visitors map after register")
	}
	if inAdmins {
		t.Error("visitor should not be in admins map")
	}

	// Admin should receive a presence broadcast
	msgs := drainMessages(admin)
	if len(msgs) == 0 {
		t.Fatal("admin did not receive presence broadcast on visitor register")
	}
	if msgs[0]["type"] != "presence" {
		t.Errorf("expected presence message, got %v", msgs[0]["type"])
	}
	if int(msgs[0]["visitor_count"].(float64)) != 1 {
		t.Errorf("visitor_count = %v, want 1", msgs[0]["visitor_count"])
	}
}

func TestRegisterAdmin(t *testing.T) {
	h := newTestHub()

	a := newTestClient(h, roleAdmin)
	h.register(a)

	h.mu.Lock()
	_, inAdmins := h.admins[a]
	_, inVisitors := h.visitors[a]
	h.mu.Unlock()

	if !inAdmins {
		t.Error("admin not in admins map after register")
	}
	if inVisitors {
		t.Error("admin should not be in visitors map")
	}
}

func TestUnregisterVisitor(t *testing.T) {
	h := newTestHub()
	admin := newTestClient(h, roleAdmin)
	h.mu.Lock()
	h.admins[admin] = true
	h.mu.Unlock()

	v := newTestClient(h, roleVisitor)
	h.register(v)
	drainMessages(admin) // clear the register presence

	h.unregister(v)

	h.mu.Lock()
	_, stillIn := h.visitors[v]
	h.mu.Unlock()

	if stillIn {
		t.Error("visitor still in map after unregister")
	}

	// Admin should get updated presence (count=0)
	msgs := drainMessages(admin)
	if len(msgs) == 0 {
		t.Fatal("admin did not receive presence broadcast on visitor unregister")
	}
	if int(msgs[0]["visitor_count"].(float64)) != 0 {
		t.Errorf("visitor_count = %v, want 0", msgs[0]["visitor_count"])
	}
}

func TestUnregisterAdmin(t *testing.T) {
	h := newTestHub()

	a := newTestClient(h, roleAdmin)
	h.register(a)
	h.unregister(a)

	h.mu.Lock()
	_, stillIn := h.admins[a]
	h.mu.Unlock()

	if stillIn {
		t.Error("admin still in map after unregister")
	}
}

func TestBroadcastHeartIncrementExcludesSender(t *testing.T) {
	h := newTestHub()

	sender := newTestClient(h, roleVisitor)
	other := newTestClient(h, roleVisitor)
	admin := newTestClient(h, roleAdmin)

	h.mu.Lock()
	h.visitors[sender] = true
	h.visitors[other] = true
	h.admins[admin] = true
	h.mu.Unlock()

	h.broadcastHeartIncrement(sender, 42, 3)

	// Sender should NOT receive the message
	senderMsgs := drainMessages(sender)
	if len(senderMsgs) != 0 {
		t.Errorf("sender received %d messages, want 0", len(senderMsgs))
	}

	// Other visitor should receive it
	otherMsgs := drainMessages(other)
	if len(otherMsgs) != 1 {
		t.Fatalf("other visitor received %d messages, want 1", len(otherMsgs))
	}
	if otherMsgs[0]["type"] != "heart_increment" {
		t.Errorf("type = %v, want heart_increment", otherMsgs[0]["type"])
	}
	if int64(otherMsgs[0]["photo_id"].(float64)) != 42 {
		t.Errorf("photo_id = %v, want 42", otherMsgs[0]["photo_id"])
	}
	if int64(otherMsgs[0]["count"].(float64)) != 3 {
		t.Errorf("count = %v, want 3", otherMsgs[0]["count"])
	}

	// Admin should also receive it
	adminMsgs := drainMessages(admin)
	if len(adminMsgs) != 1 {
		t.Fatalf("admin received %d messages, want 1", len(adminMsgs))
	}
	if adminMsgs[0]["type"] != "heart_increment" {
		t.Errorf("admin msg type = %v, want heart_increment", adminMsgs[0]["type"])
	}
}

func TestBroadcastContentUpdateVisitorsOnly(t *testing.T) {
	h := newTestHub()

	visitor := newTestClient(h, roleVisitor)
	admin := newTestClient(h, roleAdmin)

	h.mu.Lock()
	h.visitors[visitor] = true
	h.admins[admin] = true
	h.mu.Unlock()

	h.BroadcastContentUpdate("photos")

	// Visitor should receive it
	vMsgs := drainMessages(visitor)
	if len(vMsgs) != 1 {
		t.Fatalf("visitor received %d messages, want 1", len(vMsgs))
	}
	if vMsgs[0]["type"] != "content_update" {
		t.Errorf("type = %v, want content_update", vMsgs[0]["type"])
	}
	if vMsgs[0]["content_type"] != "photos" {
		t.Errorf("content_type = %v, want photos", vMsgs[0]["content_type"])
	}

	// Admin should NOT receive it
	aMsgs := drainMessages(admin)
	if len(aMsgs) != 0 {
		t.Errorf("admin received %d messages, want 0", len(aMsgs))
	}
}

func TestBroadcastPresenceAdminsOnly(t *testing.T) {
	h := newTestHub()

	visitor := newTestClient(h, roleVisitor)
	admin := newTestClient(h, roleAdmin)

	h.mu.Lock()
	h.visitors[visitor] = true
	h.admins[admin] = true
	h.mu.Unlock()

	h.broadcastPresence()

	// Admin should receive presence with visitor_count=1
	aMsgs := drainMessages(admin)
	if len(aMsgs) != 1 {
		t.Fatalf("admin received %d messages, want 1", len(aMsgs))
	}
	if aMsgs[0]["type"] != "presence" {
		t.Errorf("type = %v, want presence", aMsgs[0]["type"])
	}
	if int(aMsgs[0]["visitor_count"].(float64)) != 1 {
		t.Errorf("visitor_count = %v, want 1", aMsgs[0]["visitor_count"])
	}

	// Visitor should NOT receive it
	vMsgs := drainMessages(visitor)
	if len(vMsgs) != 0 {
		t.Errorf("visitor received %d messages, want 0", len(vMsgs))
	}
}

func TestAddHeartsBuffersAndBroadcasts(t *testing.T) {
	h := newTestHub()

	sender := newTestClient(h, roleVisitor)
	other := newTestClient(h, roleVisitor)

	h.mu.Lock()
	h.visitors[sender] = true
	h.visitors[other] = true
	h.mu.Unlock()

	h.addHearts(sender, 1, 5)

	// Other should get immediate heart_increment
	otherMsgs := drainMessages(other)
	if len(otherMsgs) != 1 {
		t.Fatalf("other received %d messages, want 1", len(otherMsgs))
	}
	if otherMsgs[0]["type"] != "heart_increment" {
		t.Errorf("type = %v, want heart_increment", otherMsgs[0]["type"])
	}

	// Buffer should have the hearts
	h.mu.Lock()
	buffered := h.buffer[1]
	h.mu.Unlock()

	if buffered != 5 {
		t.Errorf("buffer[1] = %d, want 5", buffered)
	}
}

func TestFlushEmptyBuffer(t *testing.T) {
	h := newTestHub()
	visitor := newTestClient(h, roleVisitor)

	h.mu.Lock()
	h.visitors[visitor] = true
	h.mu.Unlock()

	// Flushing empty buffer should not send messages
	h.flush()

	msgs := drainMessages(visitor)
	if len(msgs) != 0 {
		t.Errorf("received %d messages from empty flush, want 0", len(msgs))
	}
}

func TestSlowClientDropped(t *testing.T) {
	h := newTestHub()

	// Create a client with a full send buffer (size 1)
	slow := &client{
		hub:  h,
		role: roleVisitor,
		send: make(chan []byte, 1),
	}
	slow.send <- []byte("blocking")

	normal := newTestClient(h, roleVisitor)

	h.mu.Lock()
	h.visitors[slow] = true
	h.visitors[normal] = true
	h.mu.Unlock()

	h.BroadcastContentUpdate("photos")

	// Allow a moment for cleanup
	time.Sleep(10 * time.Millisecond)

	h.mu.Lock()
	_, slowStillIn := h.visitors[slow]
	_, normalStillIn := h.visitors[normal]
	h.mu.Unlock()

	if slowStillIn {
		t.Error("slow client should have been dropped")
	}
	if !normalStillIn {
		t.Error("normal client should still be registered")
	}
}
