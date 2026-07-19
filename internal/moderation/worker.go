package moderation

import (
	"context"
	"encoding/base64"
	"io"
	"log"
	"math/rand/v2"
	"time"

	"github.com/RunFridge/wedding-card/internal/models"
	"github.com/RunFridge/wedding-card/internal/storage"
)

type Worker struct {
	queue              *Queue
	client             *Client
	store              storage.Storage
	ThresholdGetter    func() map[string]float64
	OnContentApproved  func(contentType string)
}

func NewWorker(queue *Queue, client *Client, store storage.Storage, thresholdGetter func() map[string]float64) *Worker {
	return &Worker{queue: queue, client: client, store: store, ThresholdGetter: thresholdGetter}
}

func (w *Worker) readImageAsDataURL(ctx context.Context, key string) (string, error) {
	reader, info, err := w.store.GetReader(ctx, key)
	if err != nil {
		return "", err
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	ct := info.ContentType
	if ct == "" {
		ct = "image/jpeg"
	}
	return "data:" + ct + ";base64," + base64.StdEncoding.EncodeToString(data), nil
}

func (w *Worker) isFlagged(result *ModerationResult) bool {
	if result.Flagged {
		return true
	}
	if w.ThresholdGetter == nil {
		return false
	}
	thresholds := w.ThresholdGetter()
	for category, score := range result.CategoryScores {
		if threshold, ok := thresholds[category]; ok && score >= threshold {
			log.Printf("Moderation: category %q score %.4f >= threshold %.4f", category, score, threshold)
			return true
		}
	}
	return false
}

func (w *Worker) Run(ctx context.Context) {
	log.Println("Moderation worker started")
	for {
		select {
		case <-ctx.Done():
			log.Println("Moderation worker stopped")
			return
		default:
		}

		job, err := w.queue.Dequeue(ctx, 5*time.Second)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("Moderation dequeue error: %v", err)
			continue
		}
		if job == nil {
			continue
		}

		switch job.Type {
		case JobPhoto:
			w.processPhoto(ctx, job.ID)
		case JobGuestbook:
			w.processGuestbook(ctx, job.ID)
		default:
			log.Printf("Moderation worker: unknown job type %q", job.Type)
		}
	}
}

func (w *Worker) processPhoto(ctx context.Context, id int64) {
	hashname, err := models.GetPhotoHashnameByID(id)
	if err != nil {
		log.Printf("Moderation: failed to get photo %d hashname: %v", id, err)
		return
	}

	dataURL, err := w.readImageAsDataURL(ctx, "photos/"+hashname)
	if err != nil {
		log.Printf("Moderation: failed to read photo %d: %v", id, err)
		return
	}

	result, err := w.moderateWithRetry(ctx, func(ctx context.Context) (*ModerationResult, error) {
		return w.client.ModerateImage(ctx, dataURL)
	})
	if err != nil {
		log.Printf("Moderation: failed to moderate photo %d after retries: %v", id, err)
		return
	}

	hidden := w.isFlagged(result)
	if err := models.SetPhotoEvaluated(id, true, hidden); err != nil {
		log.Printf("Moderation: failed to update photo %d: %v", id, err)
		return
	}

	if hidden {
		log.Printf("Moderation: photo %d flagged", id)
	} else {
		log.Printf("Moderation: photo %d approved", id)
		if w.OnContentApproved != nil {
			w.OnContentApproved("photos")
		}
	}
}

func (w *Worker) processGuestbook(ctx context.Context, id int64) {
	nickname, message, err := models.GetGuestbookContentByID(id)
	if err != nil {
		log.Printf("Moderation: failed to get guestbook entry %d: %v", id, err)
		return
	}

	text := nickname + "\n" + message
	result, err := w.moderateWithRetry(ctx, func(ctx context.Context) (*ModerationResult, error) {
		return w.client.ModerateText(ctx, text)
	})
	if err != nil {
		log.Printf("Moderation: failed to moderate guestbook entry %d after retries: %v", id, err)
		return
	}

	hidden := w.isFlagged(result)
	if err := models.SetGuestbookEvaluated(id, true, hidden); err != nil {
		log.Printf("Moderation: failed to update guestbook entry %d: %v", id, err)
		return
	}

	if hidden {
		log.Printf("Moderation: guestbook entry %d flagged", id)
	} else {
		log.Printf("Moderation: guestbook entry %d approved", id)
		if w.OnContentApproved != nil {
			w.OnContentApproved("guestbook")
		}
	}
}

func (w *Worker) moderateWithRetry(ctx context.Context, fn func(ctx context.Context) (*ModerationResult, error)) (*ModerationResult, error) {
	var lastErr error
	for attempt := range 3 {
		result, err := fn(ctx)
		if err == nil {
			return result, nil
		}
		lastErr = err
		backoff := time.Duration(1<<uint(attempt)) * time.Second
		jitter := time.Duration(rand.Int64N(int64(500 * time.Millisecond)))
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(backoff + jitter):
		}
	}
	return nil, lastErr
}
