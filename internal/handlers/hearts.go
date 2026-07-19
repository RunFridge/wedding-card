package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"

	"github.com/RunFridge/wedding-card/internal/models"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type wsIncoming struct {
	Type    string `json:"type"`
	PhotoID int64  `json:"photo_id"`
	Count   int64  `json:"count"`
}

type heartUpdate struct {
	PhotoID int64 `json:"photo_id"`
	Hearts  int64 `json:"hearts"`
}

type wsOutgoing struct {
	Type    string        `json:"type"`
	Updates []heartUpdate `json:"updates"`
}

type HeartsConfigFn func() (flushIntervalMs int, flushBatchSize int)

type clientRole int

const (
	roleVisitor clientRole = iota
	roleAdmin
)

type client struct {
	hub  *Hub
	role clientRole
	conn *websocket.Conn
	send chan []byte
}

type Hub struct {
	mu         sync.Mutex
	visitors   map[*client]bool
	admins     map[*client]bool
	buffer     map[int64]int64
	bufferSize int64
	configFn   HeartsConfigFn
	done       chan struct{}
}

var HeartsHub *Hub

func NewHub(configFn HeartsConfigFn) *Hub {
	h := &Hub{
		visitors: make(map[*client]bool),
		admins:   make(map[*client]bool),
		buffer:   make(map[int64]int64),
		configFn: configFn,
		done:     make(chan struct{}),
	}
	go h.run()
	return h
}

func (h *Hub) Stop() {
	close(h.done)
}

func (h *Hub) ConnectionCount() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.visitors)
}

func (h *Hub) register(c *client) {
	h.mu.Lock()
	switch c.role {
	case roleAdmin:
		h.admins[c] = true
	default:
		h.visitors[c] = true
	}
	h.mu.Unlock()

	if c.role == roleVisitor {
		h.broadcastPresence()
	}
}

func (h *Hub) unregister(c *client) {
	h.mu.Lock()
	switch c.role {
	case roleAdmin:
		if _, ok := h.admins[c]; ok {
			delete(h.admins, c)
			close(c.send)
		}
	default:
		if _, ok := h.visitors[c]; ok {
			delete(h.visitors, c)
			close(c.send)
		}
	}
	h.mu.Unlock()

	if c.role == roleVisitor {
		h.broadcastPresence()
	}
}

func (h *Hub) broadcastHeartIncrement(sender *client, photoID, count int64) {
	msg, err := json.Marshal(map[string]any{
		"type":     "heart_increment",
		"photo_id": photoID,
		"count":    count,
	})
	if err != nil {
		return
	}

	h.mu.Lock()
	for c := range h.visitors {
		if c == sender {
			continue
		}
		select {
		case c.send <- msg:
		default:
			delete(h.visitors, c)
			close(c.send)
		}
	}
	for c := range h.admins {
		select {
		case c.send <- msg:
		default:
			delete(h.admins, c)
			close(c.send)
		}
	}
	h.mu.Unlock()
}

func (h *Hub) BroadcastContentUpdate(contentType string) {
	msg, err := json.Marshal(map[string]string{
		"type":         "content_update",
		"content_type": contentType,
	})
	if err != nil {
		return
	}

	h.mu.Lock()
	for c := range h.visitors {
		select {
		case c.send <- msg:
		default:
			delete(h.visitors, c)
			close(c.send)
		}
	}
	h.mu.Unlock()
}

func (h *Hub) BroadcastRankingUpdate(nickname string, timeMs int) {
	msg, err := json.Marshal(map[string]any{
		"type":    "ranking_update",
		"nickname": nickname,
		"time_ms":  timeMs,
	})
	if err != nil {
		return
	}

	h.mu.Lock()
	for c := range h.visitors {
		select {
		case c.send <- msg:
		default:
			delete(h.visitors, c)
			close(c.send)
		}
	}
	h.mu.Unlock()
}

func (h *Hub) BroadcastGameReset() {
	msg, err := json.Marshal(map[string]string{"type": "game_reset"})
	if err != nil {
		return
	}

	h.mu.Lock()
	for c := range h.visitors {
		select {
		case c.send <- msg:
		default:
			delete(h.visitors, c)
			close(c.send)
		}
	}
	h.mu.Unlock()
}

func (h *Hub) broadcastPresence() {
	h.mu.Lock()
	count := len(h.visitors)
	msg, err := json.Marshal(map[string]any{
		"type":          "presence",
		"visitor_count": count,
	})
	if err != nil {
		h.mu.Unlock()
		return
	}
	for c := range h.admins {
		select {
		case c.send <- msg:
		default:
			delete(h.admins, c)
			close(c.send)
		}
	}
	h.mu.Unlock()
}

func (h *Hub) addHearts(sender *client, photoID, count int64) {
	h.broadcastHeartIncrement(sender, photoID, count)

	h.mu.Lock()
	h.buffer[photoID] += count
	h.bufferSize += count
	_, batchSize := h.configFn()
	shouldFlush := h.bufferSize >= int64(batchSize)
	h.mu.Unlock()

	if shouldFlush {
		h.flush()
	}
}

func (h *Hub) run() {
	flushMs, _ := h.configFn()
	ticker := time.NewTicker(time.Duration(flushMs) * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			h.flush()
			newMs, _ := h.configFn()
			if newMs != flushMs {
				flushMs = newMs
				ticker.Reset(time.Duration(flushMs) * time.Millisecond)
			}
		case <-h.done:
			h.flush()
			return
		}
	}
}

func (h *Hub) flush() {
	h.mu.Lock()
	if len(h.buffer) == 0 {
		h.mu.Unlock()
		return
	}
	pending := h.buffer
	h.buffer = make(map[int64]int64)
	h.bufferSize = 0
	h.mu.Unlock()

	if err := models.IncrementPhotoHearts(pending); err != nil {
		log.Printf("hearts flush DB error: %v", err)
		h.mu.Lock()
		for id, count := range pending {
			h.buffer[id] += count
			h.bufferSize += count
		}
		h.mu.Unlock()
		return
	}

	ids := make([]int64, 0, len(pending))
	for id := range pending {
		ids = append(ids, id)
	}

	totals, err := models.GetPhotoHearts(ids)
	if err != nil {
		log.Printf("hearts flush read-back error: %v", err)
		return
	}

	updates := make([]heartUpdate, 0, len(totals))
	for id, hearts := range totals {
		updates = append(updates, heartUpdate{PhotoID: id, Hearts: hearts})
	}

	msg, err := json.Marshal(wsOutgoing{Type: "hearts_update", Updates: updates})
	if err != nil {
		return
	}

	h.mu.Lock()
	for c := range h.visitors {
		select {
		case c.send <- msg:
		default:
			delete(h.visitors, c)
			close(c.send)
		}
	}
	for c := range h.admins {
		select {
		case c.send <- msg:
		default:
			delete(h.admins, c)
			close(c.send)
		}
	}
	h.mu.Unlock()
}

func (c *client) readPump() {
	defer func() {
		c.hub.unregister(c)
		c.conn.Close()
	}()
	c.conn.SetReadLimit(4096)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		if c.role == roleAdmin {
			continue
		}
		var incoming wsIncoming
		if json.Unmarshal(message, &incoming) != nil {
			continue
		}
		if incoming.Type == "heart" && incoming.PhotoID > 0 {
			count := incoming.Count
			if count <= 0 {
				count = 1
			}
			if count > 100 {
				count = 100
			}
			c.hub.addHearts(c, incoming.PhotoID, count)
		}
	}
}

func (c *client) writePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, nil)
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func HandleHeartsWS(w http.ResponseWriter, r *http.Request) {
	if HeartsHub == nil {
		http.Error(w, "Hearts not initialized", http.StatusServiceUnavailable)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	c := &client{
		hub:  HeartsHub,
		role: roleVisitor,
		conn: conn,
		send: make(chan []byte, 64),
	}
	HeartsHub.register(c)
	go c.writePump()
	go c.readPump()
}

func HandleAdminWS(w http.ResponseWriter, r *http.Request) {
	if HeartsHub == nil {
		http.Error(w, "Hub not initialized", http.StatusServiceUnavailable)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	c := &client{
		hub:  HeartsHub,
		role: roleAdmin,
		conn: conn,
		send: make(chan []byte, 64),
	}
	HeartsHub.register(c)

	HeartsHub.mu.Lock()
	count := len(HeartsHub.visitors)
	HeartsHub.mu.Unlock()

	initial, _ := json.Marshal(map[string]any{
		"type":          "presence",
		"visitor_count": count,
	})
	select {
	case c.send <- initial:
	default:
	}

	go c.writePump()
	go c.readPump()
}

func AdminResetPhotoHearts(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		errorJSON(w, http.StatusBadRequest, "invalid_id", "Invalid photo ID")
		return
	}

	if err := models.ResetPhotoHearts(id); err != nil {
		if err == sql.ErrNoRows {
			errorJSON(w, http.StatusNotFound, "not_found", "Photo not found")
			return
		}
		http.Error(w, "Failed to reset hearts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
