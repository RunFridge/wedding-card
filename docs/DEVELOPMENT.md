# Development Guide

Technical documentation for developing and building the wedding invitation website. For deploying and using the app, see the [README](../README.md).

## Tech Stack

| Layer | Technology |
|-------|------------|
| Backend | Go 1.25, [chi](https://github.com/go-chi/chi) router |
| Database | SQLite (WAL mode, busy_timeout, synchronous=NORMAL) via [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite) (pure-Go, no CGO) |
| Frontend | Vue 3, TypeScript, Vite 6 |
| Styling | Tailwind CSS 3 |
| Animation | Phaser 3 |
| Image Editor | [Cropper.js](https://github.com/fengyuanchen/cropperjs) v1 |
| Auth | bcrypt password hashing ([golang.org/x/crypto](https://pkg.go.dev/golang.org/x/crypto)), session tokens |
| Rate Limiting | [httprate](https://github.com/go-chi/httprate) (per-IP) |
| Storage | S3-compatible (AWS SDK v2) or local filesystem |
| Queue | In-process buffered channel (moderation jobs) |
| Moderation | OpenAI Moderation API (`omni-moderation-latest`) |
| Config | Admin panel (DB-backed, no .env) |

## Project Structure

```
.
├── main.go                    # Entry point, chi router setup, embedded static serving
├── VERSION                    # Single source of truth for project version (semver)
├── Dockerfile                 # Multi-stage build (node → go (CGO-free, cross-compiled) → scratch)
├── docker-compose.yml         # Sample Compose file with default settings
├── Makefile                   # Build/dev shortcuts
├── internal/
│   ├── config/config.go       # App configuration (defaults + DB overrides)
│   ├── database/
│   │   ├── database.go        # SQLite connection (WAL mode)
│   │   └── migrations.go      # Auto-run schema migrations
│   ├── handlers/
│   │   ├── admin.go           # Admin panel endpoints
│   │   ├── guestbook.go       # CRUD + password verify endpoints
│   │   ├── game.go            # Photo selection + leaderboard endpoints
│   │   ├── hall_of_fame.go    # Hall of Fame endpoints
│   │   ├── photo.go           # Photo upload + admin photo management
│   │   ├── media.go           # Unified image proxy (streams from Store with immutable caching)
│   │   └── health.go          # Health check
│   ├── imaging/
│   │   └── imaging.go         # Server-side image processing (resize, thumbnails)
│   ├── middleware/
│   │   └── admin.go           # Admin auth middleware (Bearer + session tokens)
│   ├── models/
│   │   ├── guestbook.go       # GuestbookEntry model + queries
│   │   ├── game.go            # GameScore model + queries
│   │   ├── hall_of_fame.go    # HallOfFame model + queries
│   │   └── photo.go           # PhotoUpload model + queries
│   ├── moderation/
│   │   ├── queue.go           # In-process job queue (buffered channel)
│   │   ├── moderation.go      # OpenAI moderation API client
│   │   └── worker.go          # Background moderation worker
│   ├── session/
│   │   └── session.go         # In-memory session token store (24h TTL)
│   └── storage/
│       ├── storage.go         # Storage interface
│       ├── s3.go              # S3-compatible storage client
│       └── local.go           # Local filesystem storage
├── admin-frontend/
│   ├── index.html             # Standalone admin panel (served at /-/admin)
│   └── src/views/
│       └── HallOfFameView.vue # Admin hall of fame management
├── frontend/
│   ├── vite.config.ts         # Vite config with dev proxy + Phaser chunk splitting
│   └── src/
│       ├── main.ts
│       ├── App.vue            # Root layout with Phaser background + bottom nav
│       ├── router/index.ts    # Vue Router (/, /guestbook, /game, /map, /photo)
│       ├── services/api.ts    # Typed fetch wrappers for all API endpoints
│       ├── config/wedding.ts  # All wedding details (names, venue, dates, etc.)
│       ├── types/             # TypeScript type definitions
│       ├── lib/
│       │   └── store.ts       # Reactive persistent store
│       ├── composables/       # Vue composables (useHearts, useAchievements, useTheme)
│       ├── game/              # Phaser scene (BackgroundScene)
│       ├── views/
│       │   ├── HomeView.vue           # Main invitation page
│       │   ├── GuestbookView.vue      # Bulletin board
│       │   ├── GameView.vue           # Card game + photo gallery lightbox
│       │   ├── AchievementsView.vue   # Unlockable achievements
│       │   ├── HallOfFameView.vue     # Hall of Fame for completionists
│       │   ├── MapView.vue            # Venue directions
│       │   ├── PhotoView.vue          # Guest photo upload & gallery
│       │   └── NotFoundView.vue       # 404 page
│       ├── components/
│       │   ├── BoringAvatar.vue   # Avatar generator for guestbook
│       │   ├── LightboxViewer.vue # Image lightbox/gallery viewer
│       │   ├── PhotoEditor.vue    # Client-side image editor (crop, filters)
│       │   └── StardewDialog.vue  # Pixel-art dialog component
│       └── assets/
│           ├── main_photo.jpg # Main wedding photo
│           ├── photos/        # Bundled wedding photos (detail + thumbnail)
│           ├── characters/    # Pixel-art character sprites
│           ├── icons/         # Map provider icons
│           └── cursor/        # Custom pixel cursors
└── web/
    └── dist/                  # Vite build output (embedded into Go binary)
```

## Development

Prerequisites: Go 1.25+, Node.js 22+

```bash
# Install frontend dependencies
make install

# Run all dev servers in parallel (Ctrl+C to stop all)
make dev
```

This starts the Go backend, visitor frontend (Vite), and admin frontend (Vite) in a single terminal. The Vite dev servers proxy all `/api` requests to the Go backend at `localhost:8080`.

## Production Build

```bash
# Build frontend then Go binary (single command)
make build

# Run the server
./wedding-server
```

The `make build` command:
1. Syncs the `VERSION` file to `frontend/package.json`
2. Builds the Vue frontend via Vite → outputs to `web/dist/`
3. Compiles the Go binary with `web/dist/` embedded via `//go:embed` and version injected via ldflags

The result is a single `wedding-server` binary that serves both the API and the SPA.

## Environment Variables

Almost everything is configured in the admin panel. Only these live outside it:

| Environment variable | Default | Description |
|----------------------|---------|-------------|
| `PORT` | `8080` | Server listen port |
| `DATABASE_PATH` | `./wedding.db` | Database file path (set to `/data/wedding.db` in the Docker image) |
| `TZ` | system | Container timezone (e.g. `Asia/Seoul`) |

## API Endpoints

All endpoints are prefixed with `/api`.

### Guestbook

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/guestbook` | List visible messages (supports `cursor` and `limit` query params for pagination) |
| POST | `/api/guestbook` | Create message (nickname, message, password, optional `secret` flag for private messages) |
| POST | `/api/guestbook/{id}/verify` | Verify entry password |
| PUT | `/api/guestbook/{id}` | Update message (requires password) |
| DELETE | `/api/guestbook/{id}` | Delete message (requires password) |

### Match Card Game

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/game/photos` | Get 6 random photo IDs for a game round |
| GET | `/api/game/rankings` | Get top 10 leaderboard |
| POST | `/api/game/rankings` | Submit score (3-letter nickname + time in ms) |

### Guest Photos

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/photos` | List visible photo uploads with URLs (supports `offset` and `limit` query params for pagination) |
| GET | `/api/photos/status` | S3 storage availability check (HeadBucket) |
| GET | `/api/photos/{kind}/{hashname}` | Stream a photo (`kind` = `photo` for guest uploads, `asset` for curated assets). Returned with `Cache-Control: immutable` + hashname-based `ETag`; honors `If-None-Match` for `304` responses. |
| POST | `/api/photos/upload` | Upload photo (multipart: name, image, optional original) |

### Hall of Fame

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/hall-of-fame` | List all hall of fame entries |
| POST | `/api/hall-of-fame` | Submit hall of fame entry (nickname) |

### Health

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/health` | Health check (returns `{"status":"ok","version":"..."}`) |

### Admin (all require `Authorization: Bearer <token>` header)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/admin/verify` | Verify admin password, returns session token |
| GET | `/api/admin/guestbook` | All entries with IP + hidden status |
| PATCH | `/api/admin/guestbook/{id}/visibility` | Toggle hidden flag `{"hidden": bool}` |
| GET | `/api/admin/game/rankings` | All scores with IP |
| DELETE | `/api/admin/game/rankings/{id}` | Delete a score |
| POST | `/api/admin/game/rankings/purge` | Purge all rankings |
| GET | `/api/admin/photos` | All uploads with IP, hashname, hidden, preview URLs |
| PATCH | `/api/admin/photos/{id}/visibility` | Toggle hidden `{"hidden": bool}` |
| DELETE | `/api/admin/photos/{id}` | Delete from storage + DB |
| GET | `/api/admin/hall-of-fame` | All hall of fame entries with IP |
| DELETE | `/api/admin/hall-of-fame/{id}` | Delete a hall of fame entry |

## Database Schema

SQLite with WAL journaling. Migrations run automatically on startup.

### guestbook_entries

| Column | Type | Description |
|--------|------|-------------|
| id | INTEGER | Primary key, auto-increment |
| nickname | TEXT | Poster's display name |
| message | TEXT | Message content |
| password_hash | TEXT | bcrypt hash for edit/delete auth |
| ip | TEXT | Poster's IP address |
| hidden | INTEGER | 0 = visible, 1 = hidden |
| secret | INTEGER | 0 = public, 1 = secret (only visible to the couple) |
| evaluated | INTEGER | 0 = not evaluated, 1 = evaluated by moderation |
| created_at | DATETIME | Timestamp |

### game_scores

| Column | Type | Description |
|--------|------|-------------|
| id | INTEGER | Primary key, auto-increment |
| nickname | TEXT | 3 uppercase letters (arcade style) |
| time_ms | INTEGER | Completion time in milliseconds |
| ip | TEXT | Player's IP address |
| created_at | DATETIME | Timestamp |

### photo_uploads

| Column | Type | Description |
|--------|------|-------------|
| id | INTEGER | Primary key, auto-increment |
| name | TEXT | Uploader's name (required) |
| upload_date | DATETIME | Timestamp |
| ip_address | TEXT | Uploader's IP address |
| hashname | TEXT | UUID-based filename in S3 |
| original_hashname | TEXT | UUID-based filename for original (unedited) backup |
| original_filename | TEXT | Original file name |
| hidden | INTEGER | 1 = hidden (default), 0 = visible |
| evaluated | INTEGER | 0 = not evaluated, 1 = evaluated by moderation |

### hall_of_fame

| Column | Type | Description |
|--------|------|-------------|
| id | INTEGER | Primary key, auto-increment |
| nickname | TEXT | Player's display name |
| ip | TEXT | Player's IP address |
| created_at | DATETIME | Timestamp |

## Rate Limiting

Per-IP rate limiting via [httprate](https://github.com/go-chi/httprate). Returns `429 Too Many Requests` with `Retry-After` header when exceeded.

| Tier | Limit | Endpoints |
|------|-------|-----------|
| Read | 300 req/min | GET guestbook, photos, photos/status, game/photos, game/rankings, hall-of-fame, health |
| Write | 10 req/min | POST guestbook, PUT/DELETE guestbook/{id}, POST game/rankings, POST photos/upload, POST hall-of-fame |
| Auth | 3 req/min + 10 req/hour | POST admin/verify |
| Admin | 120 req/min | All authenticated admin/* endpoints |

> **Note:** Rate limiting can be toggled on/off via the admin panel.

## Version Management

The project version is defined in the `VERSION` file at the project root (semver, e.g. `0.1.0`). This single source of truth propagates to:

- **Go binary** — injected via `-ldflags "-X main.Version=..."` and logged on startup
- **Docker image tag** — `make docker-build` tags the image as `wedding-server:<version>`
- **Health endpoint** — `GET /api/health` returns `{"status":"ok","version":"..."}`
- **Frontend** — `make build` syncs the version to `frontend/package.json`

```bash
# Print current version
make version
```

## Docker Image Internals

The Dockerfile uses a multi-stage build:

1. **node:22-alpine** — installs deps and builds the visitor + admin frontends (pinned to `$BUILDPLATFORM` for native speed on the host runner)
2. **golang:1.25-alpine** — cross-compiles the Go binary with `CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH` and the `timetzdata` build tag so timezone data is baked in (no `/usr/share/zoneinfo` needed). Pinned to `$BUILDPLATFORM` so arm64 builds don't run the compiler under QEMU.
3. **alpine:3.21** (certs stage) — source of `ca-certificates.crt` only; not shipped in the final image.
4. **scratch** — the final image contains only the statically-linked binary (~20MB) plus the CA bundle. No shell, no libc, no package manager.

The container's `HEALTHCHECK` invokes the binary itself (`/wedding-server -healthcheck`) — a `-healthcheck` subcommand is recognized at startup and exits 0/1 based on an in-process GET to `/api/health`. This replaces the previous `wget`-based check since scratch has no coreutils. HTTP server timeouts are configured: ReadTimeout 15s, WriteTimeout 60s, IdleTimeout 120s.

```bash
# Build image (tags with version from VERSION file)
make docker-build

# Or build manually
docker build --build-arg VERSION=$(cat VERSION) -t wedding-server:$(cat VERSION) .
```

Docker log rotation is configured at 50MB max size with 3 rotated files (~150MB cap, sufficient for 3+ months of logs).

## Content Moderation Internals

When moderation is enabled, a background worker processes uploaded photos and guestbook entries through the OpenAI Moderation API via an in-process job queue (a buffered Go channel):

1. A photo is uploaded or guestbook entry is created/updated
2. A moderation job is enqueued to the in-process queue
3. The background worker picks up the job (5s poll timeout)
4. The worker calls the OpenAI Moderation API (with 3 retries and exponential backoff)
5. Based on the result, the `evaluated` and `hidden` flags are updated in the database

| Content | Passes moderation | Fails moderation | API error (all retries exhausted) |
|---------|-------------------|------------------|-----------------------------------|
| Photo | `evaluated=1, hidden=0` (auto-approve) | `evaluated=1, hidden=1` (stays hidden) | `evaluated=0` (manual review) |
| Guestbook | `evaluated=1, hidden=0` (stays visible) | `evaluated=1, hidden=1` (auto-hide) | `evaluated=0` (manual review) |

Admin can always override via the show/hide toggle in the admin panel. When a guestbook entry is updated, its `evaluated` flag resets to 0 and a new moderation job is enqueued.

The queue is not persisted: SQLite is the durable state. On startup the server re-enqueues everything still pending — hidden photos and non-secret guestbook entries with `evaluated=0` — so jobs lost to a restart or a full queue are picked up automatically.

## Game Rules

- 20 wedding photos available, each game randomly picks 6 (12 cards = 6 pairs)
- 3x4 grid layout
- 30-second time limit (configurable via admin panel)
- Timer displays centiseconds (2 decimal places) for precision scoring
- Score = completion time (faster is better)
- Nickname: exactly 3 uppercase English letters (arcade style)
- Confetti celebration for top 3 rankings
- NPC typing hint for first-time players (configurable via admin panel)
- After a game, browse all 20 photos in a lightbox gallery
- 11 unlockable achievements including game medals (gold/silver/bronze), photo interactions, and a secret wedding day achievement
- Hall of Fame for players who complete all achievements, with confetti celebration

## Adding a Language

Both apps ship Korean and English locales; contributions for more are welcome:

1. Copy `frontend/src/locales/ko.json` to `frontend/src/locales/<code>.json` and translate the values (same for `admin-frontend/src/locales/`)
2. Register it in the app's `src/i18n.ts`: add an `import` for the new file and add it to the `messages` object
3. Run `npm run test` in that app — a key-parity test verifies your file covers every key

The language switchers pick up new locales automatically.

## Testing

```bash
# Go backend tests
go test ./...

# Frontend tests (Vitest)
cd frontend && npm run test
cd admin-frontend && npm run test
```
