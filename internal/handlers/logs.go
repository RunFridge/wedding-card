package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

type LogBuffer struct {
	mu    sync.Mutex
	lines []string
	size  int
	subs  map[chan string]struct{}
}

func NewLogBuffer(size int) *LogBuffer {
	return &LogBuffer{
		lines: make([]string, 0, size),
		size:  size,
		subs:  make(map[chan string]struct{}),
	}
}

func (b *LogBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	text := strings.TrimRight(string(p), "\n")
	for _, line := range strings.Split(text, "\n") {
		if line == "" {
			continue
		}
		if len(b.lines) >= b.size {
			b.lines = b.lines[1:]
		}
		b.lines = append(b.lines, line)

		for ch := range b.subs {
			select {
			case ch <- line:
			default:
			}
		}
	}

	return len(p), nil
}

func (b *LogBuffer) subscribe() (chan string, []string) {
	ch := make(chan string, 64)
	b.mu.Lock()
	history := make([]string, len(b.lines))
	copy(history, b.lines)
	b.subs[ch] = struct{}{}
	b.mu.Unlock()
	return ch, history
}

func (b *LogBuffer) unsubscribe(ch chan string) {
	b.mu.Lock()
	delete(b.subs, ch)
	b.mu.Unlock()
	close(ch)
}

var LogBuf *LogBuffer

func AdminStreamLogs(w http.ResponseWriter, r *http.Request) {
	if LogBuf == nil {
		http.Error(w, "Logs not available", http.StatusServiceUnavailable)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ch, history := LogBuf.subscribe()
	defer LogBuf.unsubscribe(ch)

	for _, line := range history {
		fmt.Fprintf(w, "data: %s\n\n", line)
	}
	flusher.Flush()

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case line := <-ch:
			fmt.Fprintf(w, "data: %s\n\n", line)
			flusher.Flush()
		case <-ticker.C:
			fmt.Fprintf(w, ": keepalive\n\n")
			flusher.Flush()
		}
	}
}
