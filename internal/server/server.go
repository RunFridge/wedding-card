package server

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"

	"github.com/RunFridge/wedding-card/internal/config"
	"github.com/RunFridge/wedding-card/internal/handlers"
	customMiddleware "github.com/RunFridge/wedding-card/internal/middleware"
	"github.com/RunFridge/wedding-card/internal/models"
)

type Deps struct {
	StaticFS fs.FS
	AdminFS  fs.FS
}

func NewRouter(deps Deps) chi.Router {
	r := chi.NewRouter()

	r.Use(customMiddleware.SecurityHeaders())
	r.Use(middleware.RealIP)
	r.Use(requestLogger)
	r.Use(middleware.Recoverer)

	conditionalRateLimit := func(requestLimit int, windowLength time.Duration) func(http.Handler) http.Handler {
		limiter := httprate.LimitByIP(requestLimit, windowLength)
		return func(next http.Handler) http.Handler {
			limited := limiter(next)
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if config.Cfg.RateLimitEnabled {
					limited.ServeHTTP(w, r)
				} else {
					next.ServeHTTP(w, r)
				}
			})
		}
	}

	registerAPIRoutes(r, conditionalRateLimit)

	registerSPAHandler(r, deps.StaticFS, deps.AdminFS)

	return r
}

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)
		log.Printf("%s %s %s %d %s %d bytes",
			r.RemoteAddr,
			r.Method,
			r.RequestURI,
			ww.Status(),
			fmt.Sprintf("%.1fms", float64(time.Since(start).Microseconds())/1000),
			ww.BytesWritten(),
		)
	})
}

func registerAPIRoutes(r chi.Router, rateLimit func(int, time.Duration) func(http.Handler) http.Handler) {
	r.Route("/api", func(r chi.Router) {
		// Read tier: 300 req/min per IP
		r.Group(func(r chi.Router) {
			r.Use(rateLimit(300, time.Minute))
			r.Get("/config", handlers.GetWeddingConfig)
			r.Get("/health", handlers.HealthCheck)
			r.Get("/guestbook", handlers.GetGuestbook)
			r.Get("/game/photos", handlers.GetGamePhotos)
			r.Get("/game/rankings", handlers.GetGameRankings)
			r.Get("/photos", handlers.GetVisiblePhotos)
			r.Get("/photos/status", handlers.PhotoStorageStatus)
			r.Get("/photos/{kind}/{hashname}", handlers.GetMedia)
			r.Get("/ws", handlers.HandleHeartsWS)
			r.Get("/photos/ws", handlers.HandleHeartsWS)
			r.Get("/hall-of-fame", handlers.GetHallOfFame)
		})

		// Write tier: 10 req/min per IP
		r.Group(func(r chi.Router) {
			r.Use(rateLimit(10, time.Minute))
			r.Post("/guestbook", handlers.CreateGuestbookEntry)
			r.Post("/guestbook/{id}/verify", handlers.VerifyGuestbookPassword)
			r.Put("/guestbook/{id}", handlers.UpdateGuestbookEntry)
			r.Delete("/guestbook/{id}", handlers.DeleteGuestbookEntry)
			r.Post("/game/rankings", handlers.CreateGameScore)
			r.Post("/game/beats", handlers.RecordGameBeat)
			r.Post("/photos/{id}/verify", handlers.VerifyPhotoPassword)
			r.Delete("/photos/{id}", handlers.UserDeletePhoto)
			r.Post("/hall-of-fame", handlers.CreateHallOfFameEntry)
		})

		// Upload tier: 30 req/min per IP — allows multi-photo batches
		r.Group(func(r chi.Router) {
			r.Use(rateLimit(30, time.Minute))
			r.Post("/photos/upload", handlers.UploadPhoto)
		})

		r.Route("/admin", func(r chi.Router) {
			// Auth tier: 3 req/min + 10 req/hour per IP
			r.Group(func(r chi.Router) {
				r.Use(rateLimit(3, time.Minute))
				r.Use(rateLimit(10, time.Hour))
				r.Post("/verify", handlers.AdminVerify)
			})

			// Admin tier: 120 req/min per IP
			r.Group(func(r chi.Router) {
				r.Use(customMiddleware.AdminAuth)
				r.Use(rateLimit(120, time.Minute))
				r.Post("/logout", handlers.AdminLogout)
				r.Get("/guestbook", handlers.AdminGetGuestbook)
				r.Patch("/guestbook/{id}/visibility", handlers.AdminToggleGuestbookVisibility)
				r.Delete("/guestbook/{id}", handlers.AdminDeleteGuestbookEntry)
				r.Get("/game/rankings", handlers.AdminGetGameRankings)
				r.Delete("/game/rankings/{id}", handlers.AdminDeleteGameScore)
				r.Post("/game/rankings/purge", handlers.AdminPurgeGameRankings)
				r.Get("/photos", handlers.AdminGetPhotos)
				r.Patch("/photos/{id}/visibility", handlers.AdminTogglePhotoVisibility)
				r.Delete("/photos/{id}", handlers.AdminDeletePhoto)
				r.Post("/photos/{id}/reset-hearts", handlers.AdminResetPhotoHearts)
				r.Get("/session", handlers.AdminSession)
				r.Put("/password", handlers.AdminChangePassword)
				r.Get("/moderation/status", handlers.AdminGetModerationStatus)
				r.Get("/config", handlers.AdminGetConfig)
				r.Put("/config", handlers.AdminUpdateConfig)
				r.Get("/asset-photos", handlers.AdminGetAssetPhotos)
				r.Post("/asset-photos", handlers.AdminUploadAssetPhoto)
				r.Patch("/asset-photos/{id}/game", handlers.AdminToggleAssetPhotoGame)
				r.Patch("/asset-photos/{id}/main", handlers.AdminSetMainPhoto)
				r.Patch("/asset-photos/{id}", handlers.AdminUpdateAssetPhoto)
				r.Delete("/asset-photos/{id}", handlers.AdminDeleteAssetPhoto)
				r.Get("/system-settings", handlers.AdminGetSystemSettings)
				r.Put("/system-settings", handlers.AdminUpdateSystemSettings)
				r.Post("/system-settings/test-s3", handlers.AdminTestS3Connection)
				r.Post("/system-settings/test-moderation", handlers.AdminTestModeration)
				r.Get("/hall-of-fame", handlers.AdminGetHallOfFame)
				r.Delete("/hall-of-fame/{id}", handlers.AdminDeleteHallOfFameEntry)
				r.Get("/page-views", handlers.AdminGetPageViews)
				r.Get("/game-beats", handlers.AdminGetGameBeats)
				r.Get("/logs", handlers.AdminStreamLogs)
				r.Post("/setup/complete", handlers.AdminCompleteSetup)
				r.Post("/restart", handlers.AdminRestartServer)
				r.Get("/ws", handlers.HandleAdminWS)
			})
		})
	})
}

func withCacheHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.Contains(path, "-") && (strings.HasSuffix(path, ".js") || strings.HasSuffix(path, ".css") || strings.HasSuffix(path, ".woff") || strings.HasSuffix(path, ".woff2") || strings.HasSuffix(path, ".svg") || strings.HasSuffix(path, ".png") || strings.HasSuffix(path, ".jpg")) {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		} else if strings.HasSuffix(path, ".ico") || strings.HasSuffix(path, ".webmanifest") || strings.HasSuffix(path, ".txt") {
			w.Header().Set("Cache-Control", "public, max-age=86400")
		}
		h.ServeHTTP(w, r)
	})
}

func registerSPAHandler(r chi.Router, staticFS, adminFS fs.FS) {
	distFS, err := fs.Sub(staticFS, "web/dist")
	if err != nil {
		log.Fatalf("Failed to create sub filesystem: %v", err)
	}

	adminSubFS, err := fs.Sub(adminFS, "web/admin")
	if err != nil {
		log.Fatalf("Failed to create admin sub filesystem: %v", err)
	}
	adminFileServer := withCacheHeaders(http.FileServer(http.FS(adminSubFS)))
	fileServer := withCacheHeaders(http.FileServer(http.FS(distFS)))

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		// Admin SPA
		if r.URL.Path == "/-/admin" || strings.HasPrefix(r.URL.Path, "/-/admin/") {
			path := strings.TrimPrefix(r.URL.Path, "/-/admin")
			path = strings.TrimPrefix(path, "/")
			if path == "" {
				r.URL.Path = "/"
				adminFileServer.ServeHTTP(w, r)
				return
			}
			if _, err := fs.Stat(adminSubFS, path); err == nil {
				r.URL.Path = "/" + path
				adminFileServer.ServeHTTP(w, r)
				return
			}
			data, _ := fs.ReadFile(adminSubFS, "index.html")
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(data)
			return
		}

		// Redirect to admin setup on fresh install
		if config.Cfg.SetupRequired && !strings.HasPrefix(r.URL.Path, "/-/admin") {
			http.Redirect(w, r, "/-/admin/", http.StatusFound)
			return
		}

		// Block dotfile access (e.g. /.env, /.git)
		if strings.HasPrefix(r.URL.Path, "/.") {
			http.NotFound(w, r)
			return
		}

		// Visitor frontend
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}
		if _, err := fs.Stat(distFS, path); err == nil {
			if path == "index.html" && shouldCountPageView(r) {
				go models.RecordPageView()
			}
			fileServer.ServeHTTP(w, r)
			return
		}
		// SPA fallback — serves index.html for client-side routes
		if shouldCountPageView(r) {
			go models.RecordPageView()
		}
		r.URL.Path = "/"
		fileServer.ServeHTTP(w, r)
	})
}
