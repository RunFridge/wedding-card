# Configurable Charter Bus Notice

Date: 2026-07-09
Status: Approved

## Problem

`frontend/src/views/MapView.vue` (lines 75–90) hardcodes a "📢 전세버스 탑승 안내" notice
stating the bride's charter bus is expected to be full (만석) and asking guests to report
attendance changes. It renders whenever any charter bus entry exists. The text is
wedding-specific and cannot be changed or hidden without editing code — a problem for
open-sourcing the project and for any deployment that doesn't want this exact wording.

## Decision

Replace the hardcoded text with a configurable multi-line text setting,
`charter_bus_notice`. Empty text hides the notice entirely — the emptiness is the toggle.
Default is empty, so the open-sourced repo ships with no wedding-specific text.

Rejected alternatives:

- Boolean toggle with hardcoded text: keeps deployment-specific Korean text in the code.
- Boolean toggle plus editable text: two settings for one small notice; empty-means-hidden
  covers both needs.

## Changes

### Backend — `internal/config/config.go`

- Add `CharterBusNotice string` with JSON key `charter_bus_notice` to `WeddingConfig`.
- Default: `""` in `loadWeddingDefaults()`.
- Wire into `ApplyOverrides` via `str("charter_bus_notice", &c.CharterBusNotice)`.
- Wire into `DiffToOverrideMap` via `strDiff("charter_bus_notice", ...)`.
- The field flows through the existing `wedding_config_overrides` persistence and is
  exposed by `GET /api/config` like every other `WeddingConfig` field. No handler changes.

### Visitor frontend

- `frontend/src/config/wedding.ts`: `export let CHARTER_BUS_NOTICE = '';` and hydrate in
  `loadConfig()` with `if (c.charter_bus_notice) CHARTER_BUS_NOTICE = c.charter_bus_notice;`.
- `frontend/src/views/MapView.vue`: replace the three hardcoded paragraphs with a single
  paragraph rendering `CHARTER_BUS_NOTICE` using `whitespace-pre-line break-keep`. The
  notice panel gets `v-if="CHARTER_BUS_NOTICE"`. The "📢 전세버스 탑승 안내" header stays
  hardcoded inside the panel.

### Admin frontend

- `admin-frontend/src/types/admin.ts`: add `charter_bus_notice: string` to the wedding
  config type.
- `admin-frontend/src/views/SettingsView.vue`: add a `Textarea` bound to
  `form.charter_bus_notice` under the charter bus list, following the `main_greet_text`
  pattern. Label indicates that empty hides the notice (e.g. "전세버스 탑승 안내 문구 —
  비워두면 표시되지 않습니다").

## Behavior

- Notice panel shows only when `CHARTER_BUS_NOTICE` is non-empty.
- The charter bus list itself continues to show/hide independently based on entries.
- Newlines in the configured text are preserved (`whitespace-pre-line`); inline emphasis
  (the current `<strong>` on "좌석이 만석으로 예상") is not supported — plain text only.

## Testing

- Go: extend `internal/config/config_test.go` to cover apply/diff of `charter_bus_notice`.
- Frontends: run both Vitest suites (`npm run test` in `frontend/` and `admin-frontend/`).
