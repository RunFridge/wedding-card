package main

import (
	"context"
	"embed"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/RunFridge/wedding-card/internal/config"
	"github.com/RunFridge/wedding-card/internal/database"
	"github.com/RunFridge/wedding-card/internal/demo"
	"github.com/RunFridge/wedding-card/internal/handlers"
	customMiddleware "github.com/RunFridge/wedding-card/internal/middleware"
	"github.com/RunFridge/wedding-card/internal/models"
	"github.com/RunFridge/wedding-card/internal/moderation"
	"github.com/RunFridge/wedding-card/internal/server"
	"github.com/RunFridge/wedding-card/internal/session"
	"github.com/RunFridge/wedding-card/internal/storage"
)

var Version = "dev"

//go:embed web/dist/*
var staticFiles embed.FS

//go:embed web/admin/*
var adminFiles embed.FS

func main() {
	if len(os.Args) > 1 && os.Args[1] == "-healthcheck" {
		runHealthcheck()
		return
	}

	handlers.LogBuf = handlers.NewLogBuffer(500)
	log.SetOutput(io.MultiWriter(os.Stderr, handlers.LogBuf))

	cfg := config.Load()
	session.Init()
	session.InitGame()
	handlers.Version = Version

	if err := database.Init(cfg.DatabasePath); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	config.LoadAndApplySystemSettings()

	if err := config.InitIPHashKey(); err != nil {
		log.Fatalf("Failed to initialize IP hash key: %v", err)
	}

	if err := config.InitAdminPassword(); err != nil {
		log.Fatalf("Failed to initialize admin password: %v", err)
	}

	if err := config.InitSetupState(); err != nil {
		log.Fatalf("Failed to initialize setup state: %v", err)
	}

	initStorage(cfg)

	if cfg.DemoMode {
		if err := demo.Init(cfg.DemoResetCron); err != nil {
			log.Fatalf("Failed to initialize demo mode: %v", err)
		}
	}

	initHeartsHub(cfg)
	initModeration(cfg)

	if err := models.CleanupOldPageViews(6); err != nil {
		log.Printf("Failed to cleanup old page views: %v", err)
	}
	if err := models.CleanupOldGameBeats(6); err != nil {
		log.Printf("Failed to cleanup old game beats: %v", err)
	}

	// Initialize map CSP from config
	merged := cfg.Wedding
	if overrides, err := models.GetConfigOverrides(); err == nil {
		merged = config.ApplyOverrides(merged, overrides)
	}
	customMiddleware.SetMapCSP(merged.MapProviders.EmbedProvider)

	router := server.NewRouter(server.Deps{
		StaticFS: staticFiles,
		AdminFS:  adminFiles,
	})

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	handlers.ShutdownFunc = func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		srv.Shutdown(ctx)
		os.Exit(0)
	}

	log.Printf("Server v%s starting on :%s", Version, cfg.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
}

func initStorage(cfg *config.Config) {
	localPhotosDir := filepath.Join(filepath.Dir(cfg.DatabasePath), "photos")
	if cfg.S3Bucket != "" {
		s3Client, err := storage.NewS3Client(
			cfg.S3Endpoint, cfg.S3Region, cfg.S3Bucket,
			cfg.S3AccessKey, cfg.S3SecretKey,
		)
		if err != nil {
			log.Fatalf("Failed to initialize S3 client: %v", err)
		}
		handlers.Store = s3Client
		log.Println("S3 storage configured for photo uploads")
	} else {
		local, err := storage.NewLocalClient(localPhotosDir)
		if err != nil {
			log.Fatalf("Failed to initialize local storage: %v", err)
		}
		handlers.Store = local
		log.Println("Local file storage configured for photo uploads")
	}

	handlers.ReinitStorage = func(endpoint, region, bucket, accessKey, secretKey string) error {
		if bucket == "" {
			local, err := storage.NewLocalClient(localPhotosDir)
			if err != nil {
				return err
			}
			handlers.Store = local
			log.Println("Switched to local file storage via system settings")
			return nil
		}
		s3Client, err := storage.NewS3Client(endpoint, region, bucket, accessKey, secretKey)
		if err != nil {
			return err
		}
		handlers.Store = s3Client
		log.Println("S3 storage re-initialized via system settings")
		return nil
	}
	handlers.UpdateMapCSP = customMiddleware.SetMapCSP
}

func runHealthcheck() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get("http://127.0.0.1:" + port + "/api/health")
	if err != nil {
		os.Exit(1)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		os.Exit(1)
	}
}

func initHeartsHub(cfg *config.Config) {
	handlers.HeartsHub = handlers.NewHub(func() (int, int) {
		merged := cfg.Wedding
		if overrides, err := models.GetConfigOverrides(); err == nil {
			merged = config.ApplyOverrides(merged, overrides)
		}
		return merged.HeartsFlushIntervalMs, merged.HeartsFlushBatchSize
	})
}

func initModeration(cfg *config.Config) {
	if !cfg.UseModeration {
		return
	}

	queue := moderation.NewQueue()

	moderationClient := moderation.NewClient(cfg.OpenAIAPIKey)
	thresholdGetter := func() map[string]float64 {
		merged := cfg.Wedding
		if overrides, err := models.GetConfigOverrides(); err == nil {
			merged = config.ApplyOverrides(merged, overrides)
		}
		return merged.ModerationThresholds
	}
	worker := moderation.NewWorker(queue, moderationClient, handlers.Store, thresholdGetter)

	handlers.OnPhotoUploaded = func(ctx context.Context, photoID int64) {
		if err := queue.Enqueue(ctx, moderation.Job{Type: moderation.JobPhoto, ID: photoID}); err != nil {
			log.Printf("Failed to enqueue photo moderation job: %v", err)
		}
	}
	handlers.OnGuestbookEntryCreated = func(ctx context.Context, entryID int64) {
		if err := queue.Enqueue(ctx, moderation.Job{Type: moderation.JobGuestbook, ID: entryID}); err != nil {
			log.Printf("Failed to enqueue guestbook moderation job: %v", err)
		}
	}
	handlers.ModerationQueueLen = queue.Len

	worker.OnContentApproved = func(contentType string) {
		handlers.HeartsHub.BroadcastContentUpdate(contentType)
	}

	ctx, cancel := context.WithCancel(context.Background())
	_ = cancel // cancel called on process exit
	go worker.Run(ctx)

	go enqueuePendingModeration(ctx, queue)
}

func enqueuePendingModeration(ctx context.Context, queue *moderation.Queue) {
	photoIDs, err := models.GetUnevaluatedPhotoIDs()
	if err != nil {
		log.Printf("Failed to load unevaluated photos: %v", err)
	}
	for _, id := range photoIDs {
		queue.Enqueue(ctx, moderation.Job{Type: moderation.JobPhoto, ID: id})
	}

	entryIDs, err := models.GetUnevaluatedGuestbookIDs()
	if err != nil {
		log.Printf("Failed to load unevaluated guestbook entries: %v", err)
	}
	for _, id := range entryIDs {
		queue.Enqueue(ctx, moderation.Job{Type: moderation.JobGuestbook, ID: id})
	}

	if n := len(photoIDs) + len(entryIDs); n > 0 {
		log.Printf("Re-enqueued %d pending moderation jobs", n)
	}
}
