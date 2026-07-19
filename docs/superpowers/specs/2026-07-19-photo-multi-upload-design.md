# Photo Multi-Upload + Faster Uploads

## Problem

Production photo uploads are slow and single-file only:

- Raw camera files (5–10MB) upload unresized, but the server immediately downsizes everything to 2048px JPEG (`internal/imaging`) — most of the transferred bytes are thrown away.
- Edited uploads send two full-size files (edited + raw original), doubling the payload.
- `PhotoView.vue` holds exactly one `selectedFile`; guests must repeat the whole form per photo.
- `POST /api/photos/upload` sits in the 10/min write rate tier, which would 429 a multi-photo batch.

## Design

### Client-side compression (frontend/src/lib/upload.ts)

`compressImage(file)`: decode via `createImageBitmap` (`imageOrientation: 'from-image'`), scale to max 2048px, re-encode as JPEG q0.85 via canvas. Skip if already a small JPEG within bounds; fall back to the original file on any error or if compression doesn't shrink it. Matches the server's own output parameters, so no quality is lost versus today. Applied to every uploaded file (edited and original) just before upload.

`runWithConcurrency(items, limit, worker)`: tiny worker-pool used to upload 3 photos at a time.

### Multi-select queue (PhotoView.vue)

- Gallery input gains `multiple`; camera input stays single-shot. New selections append to a queue (`{file, original, preview, status}`); the retake-confirm flow is removed.
- Queue renders as a thumbnail grid: tap a thumbnail to open the existing PhotoEditor for that photo (confirm replaces that item's file/preview); an × button removes an item.
- One shared name + password for the batch. Submit uploads all pending items with concurrency 3, reusing the existing per-photo `POST /api/photos/upload` — per-photo status, error isolation, and moderation flow unchanged.
- Per-item status drives the UI (spinner overlay while uploading, red ring on failure). Upload button shows `{done}/{total}` progress. Successful items leave the queue; failed items stay for one-tap retry. Partial failure shows a translated count message; a 503 keeps the existing "uploads disabled" message.

### Backend

`POST /api/photos/upload` moves out of the 10/min write tier into its own 30/min upload tier in `internal/server/server.go`. Handler, validation, moderation, and storage are unchanged.

### i18n

New keys in both `frontend/src/locales/{ko,en}.json`: `photo.addMore`, `photo.uploadCount`, `photo.uploadingProgress`, `photo.uploadPartialFailed`, `photo.editHint`, `photo.removePhotoAria`. Removed: `photo.reselect`, `photo.edit`, `photo.retakeConfirm`. Key-parity test keeps catalogs in sync.

## Out of scope

- Batch upload endpoint (per-photo requests reuse everything and give free retry granularity).
- Server-side resize algorithm changes (client compression makes the server resize a no-op for the common path).
- Upload progress bars per file (status states suffice).

## Verification

- `npm run test` + `npm run build` in `frontend/`
- `go test ./internal/...`
- Manual: select several photos, edit one via tap, upload; confirm batch appears in admin pending list; confirm >10-photo batch doesn't 429.
