package moderation

import (
	"context"
	"log"
	"time"
)

const queueCapacity = 1000

type JobType string

const (
	JobPhoto     JobType = "photo"
	JobGuestbook JobType = "guestbook"
)

type Job struct {
	Type JobType `json:"type"`
	ID   int64   `json:"id"`
}

type Queue struct {
	jobs chan Job
}

func NewQueue() *Queue {
	return &Queue{jobs: make(chan Job, queueCapacity)}
}

func (q *Queue) Enqueue(ctx context.Context, job Job) error {
	select {
	case q.jobs <- job:
	default:
		log.Printf("Moderation queue full, dropping job %s:%d (will be re-scanned on restart)", job.Type, job.ID)
	}
	return nil
}

func (q *Queue) Dequeue(ctx context.Context, timeout time.Duration) (*Job, error) {
	select {
	case job := <-q.jobs:
		return &job, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(timeout):
		return nil, nil
	}
}

func (q *Queue) Len(ctx context.Context) (int64, error) {
	return int64(len(q.jobs)), nil
}

func (q *Queue) Close() error {
	return nil
}
