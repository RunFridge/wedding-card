# Wedding Invitation Website

Go (chi) backend + Vue.js 3 frontend wedding invitation with interactive features.

## Tech Stack

- **Backend**: Go + chi router, SQLite (WAL mode), S3/local storage abstraction, in-process moderation job queue, OpenAI moderation (optional), embedded static serving
- **Visitor Frontend** (`frontend/`): Vue 3 + TS + Vite, Tailwind, Cropper.js v1, Phaser 3 pixel-art background, Twemoji SVGs, Vitest
- **Admin Frontend** (`admin-frontend/`): Vue 3 + TS + Vite, shadcn-vue (new-york) + reka-ui, Tailwind, Vitest. Build output: `web/admin/`, served at `/-/admin/`

## Key Directories

- `internal/` — Go backend (config, database, demo, handlers, imaging, middleware, models, moderation, session, storage)
- `frontend/src/` — Visitor Vue app (views, components, composables, game/, services/api.ts, config/wedding.ts, lib/store.ts)
- `admin-frontend/src/` — Admin Vue app (views, components, composables, lib/axios.ts + auth.ts)
- `web/dist/` + `web/admin/` — Embedded build outputs

## Build & Run

```bash
make install          # Install frontend deps
make dev              # All dev servers (backend :8080, visitor :5173, admin :5174)
make build            # Production build (both frontends + Go binary)
make docker-up        # Docker build + compose up + show password
make docker-down      # Compose down -v
```

## Environment Variables

| Name | Default | Description |
|------|---------|-------------|
| `PORT` | `8080` | Server port |
| `DATABASE_PATH` | `./wedding.db` | SQLite database path |
| `DEMO_MODE` | _(unset)_ | Seed demo data, force admin password `demo_1234!`, show DEMO ribbon, block dangerous admin actions |
| `DEMO_RESET_CRON` | `0 2 * * 6` | Full wipe + reseed schedule (demo mode only) |

All other config is done through the admin panel (`/-/admin/settings` for wedding details, `/-/admin/system` for system settings).

## API Endpoints

All prefixed with `/api`. Admin endpoints require `Authorization: Bearer <token>`.

**Public**: `GET /api/config`, `GET /api/health`, `GET /api/guestbook`, `POST /api/guestbook`, `POST /api/guestbook/{id}/verify`, `PUT /api/guestbook/{id}`, `DELETE /api/guestbook/{id}`, `GET /api/game/photos`, `GET /api/game/rankings`, `POST /api/game/rankings`, `POST /api/game/beats`, `GET /api/photos`, `GET /api/photos/status`, `POST /api/photos/upload`, `POST /api/photos/{id}/verify`, `DELETE /api/photos/{id}`, `GET /api/hall-of-fame`, `POST /api/hall-of-fame`

**WebSocket**: `WS /api/ws` (visitor), `WS /api/admin/ws` (admin, auth required)

**Pages**: `GET /simple` — server-rendered plain HTML invitation (Korean-only, `internal/handlers/simple.go` + embedded `simple.html`); 302s to `simple_redirect_url` wedding-config value when set to a valid http(s) URL

**Admin**: `POST /api/admin/verify`, `GET /api/admin/session`, `GET|PATCH|DELETE /api/admin/guestbook[/{id}[/visibility]]`, `GET|DELETE /api/admin/game/rankings[/{id}]`, `POST /api/admin/game/rankings/purge`, `GET|PATCH|DELETE /api/admin/photos[/{id}[/visibility]]`, `GET|POST|PATCH|DELETE /api/admin/asset-photos[/{id}[/game|/main]]`, `GET|DELETE /api/admin/hall-of-fame[/{id}]`, `GET|PUT /api/admin/config`, `GET|PUT /api/admin/system-settings`, `POST /api/admin/system-settings/test-s3`, `POST /api/admin/system-settings/test-moderation`, `GET /api/admin/page-views`, `GET /api/admin/game-beats`, `GET /api/admin/logs` (SSE), `GET /api/admin/moderation/status`, `POST /api/admin/setup/complete`

## Database Tables

- **guestbook_entries**: id, nickname, message, password_hash, ip, hidden, evaluated, secret, created_at
- **game_scores**: id, nickname (3 uppercase letters), time_ms, ip, created_at
- **photo_uploads**: id, name, upload_date, ip_address, hashname, original_hashname, original_filename, hidden (default 1), evaluated, password_hash
- **asset_photos**: id, label, hashname, thumb_hashname, original_filename, use_for_game, is_main_photo, sort_order, created_at
- **wedding_config_overrides**: key (PK; `sys:` prefix for system settings), value, updated_at
- **hall_of_fame**: id, nickname, ip, created_at
- **page_views**: date (PK, YYYY-MM-DD), count
- **game_beats**: date (PK, YYYY-MM-DD), count

## Version Bump

`VERSION` file is single source of truth. `make sync-version` propagates to `frontend/package.json`. Backend reads via `-ldflags`.

```bash
echo "X.Y.Z" > VERSION && make sync-version
git add VERSION frontend/package.json frontend/package-lock.json
git commit -m "chore: bump version to X.Y.Z"
git tag -a vX.Y.Z -m "vX.Y.Z" && git push && git push --tags
```

## Key Conventions

- Admin axios: defaults to `Content-Type: application/json`; FormData uploads must use `{ headers: { 'Content-Type': undefined } }`
- Admin axios has 429 retry interceptor (reads `Retry-After`, retries up to 3x) and 401 redirect to login
- Storage interface: handlers use `Store`/`ReinitStorage` (not S3-specific names); local stores under `<db-dir>/photos/` at `/storage/*`
- CSP `img-src` dynamically includes S3 origin (path-style `s3.{region}.amazonaws.com`); updated at runtime on settings change
- Config override system: `wedding_config_overrides` table, `sys:` prefix for system settings; `SetConfigOverrides` preserves `sys:` keys
- Photos hidden by default, require admin approval; moderation auto-approves clean content
- Verify endpoints return uniform 401 for not-found and wrong-password (prevent enumeration)
- Game anti-cheat: one-time session token from `GET /api/game/photos`, 2min TTL, 3s minimum elapsed
- Go style: `any` not `interface{}`; Korean text literals not unicode escapes; emoji literals
- Prettier: `singleQuote: true, tabWidth: 2` in both frontends
- Rate limiting: read 300/min, write 10/min, auth 3/min+10/hr, admin 120/min; toggleable via `rate_limit_enabled`
- SQLite: WAL mode, `busy_timeout=5000`, `synchronous=NORMAL`
- Demo mode (`internal/demo`): seeds curated Korean mock data + generated placeholder photos, marker key `sys:demo_seeded`; `demoGuard` in server.go 403s password change, system-settings write, S3/moderation tests, restart; reset = full table wipe + storage delete + reseed via robfig/cron
- Gzip handled by Nginx, not Go
- Tests: `npm run test` in each frontend directory (Vitest)

## Rate Limiting

Per-IP via httprate. Toggleable via admin System Settings (`rate_limit_enabled`).

| Tier | Limit | Scope |
|------|-------|-------|
| Read | 300/min | GET endpoints |
| Write | 10/min | POST/PUT/DELETE mutations |
| Upload | 30/min | POST /photos/upload |
| Auth | 3/min + 10/hr | admin/verify |
| Admin | 120/min | All admin/* |
