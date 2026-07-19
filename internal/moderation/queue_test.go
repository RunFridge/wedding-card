package moderation

import (
	"context"
	"testing"
	"time"
)

func TestQueueRoundtrip(t *testing.T) {
	q := NewQueue()
	ctx := context.Background()

	if err := q.Enqueue(ctx, Job{Type: JobPhoto, ID: 42}); err != nil {
		t.Fatalf("Enqueue() error: %v", err)
	}

	if n, _ := q.Len(ctx); n != 1 {
		t.Errorf("Len() = %d, want 1", n)
	}

	job, err := q.Dequeue(ctx, time.Second)
	if err != nil {
		t.Fatalf("Dequeue() error: %v", err)
	}
	if job == nil || job.Type != JobPhoto || job.ID != 42 {
		t.Errorf("Dequeue() = %+v, want {photo 42}", job)
	}
}

func TestQueueDequeueTimeout(t *testing.T) {
	q := NewQueue()

	job, err := q.Dequeue(context.Background(), 10*time.Millisecond)
	if err != nil {
		t.Fatalf("Dequeue() error: %v", err)
	}
	if job != nil {
		t.Errorf("Dequeue() = %+v, want nil on timeout", job)
	}
}

func TestQueueDropsWhenFull(t *testing.T) {
	q := NewQueue()
	ctx := context.Background()

	for i := 0; i < queueCapacity+10; i++ {
		if err := q.Enqueue(ctx, Job{Type: JobGuestbook, ID: int64(i)}); err != nil {
			t.Fatalf("Enqueue() error: %v", err)
		}
	}

	if n, _ := q.Len(ctx); n != queueCapacity {
		t.Errorf("Len() = %d, want %d", n, queueCapacity)
	}
}
