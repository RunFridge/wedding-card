# Configurable Charter Bus Notice Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace the hardcoded 전세버스 탑승 안내 notice in the visitor Map view with a `charter_bus_notice` wedding config text setting where empty text hides the notice.

**Architecture:** Add one string field to the Go `WeddingConfig` struct (persisted via the existing `wedding_config_overrides` mechanism, exposed automatically through `GET /api/config`), hydrate it in the visitor frontend's runtime config loader, render it conditionally in `MapView.vue`, and expose a textarea in the admin Settings view.

**Tech Stack:** Go (chi), Vue 3 + TypeScript, Vitest, Go testing.

**Spec:** `docs/superpowers/specs/2026-07-09-charter-bus-notice-design.md`

> **Note on commits:** This repository currently has NO commits (fresh `git init` for open-sourcing). Confirm with the user before running any commit step — they may want to create their initial commit first. If commits are skipped, still complete every other step.

---

### Task 1: Backend `charter_bus_notice` field

**Files:**
- Modify: `internal/config/config.go` (struct ~line 102, `ApplyOverrides` ~line 312, `DiffToOverrideMap` ~line 361)
- Test: `internal/config/config_test.go`

- [ ] **Step 1: Write the failing test**

Append to `internal/config/config_test.go`:

```go
func TestCharterBusNoticeOverride(t *testing.T) {
	base := loadWeddingDefaults()

	if base.CharterBusNotice != "" {
		t.Errorf("default CharterBusNotice = %q, want empty", base.CharterBusNotice)
	}

	notice := "좌석이 만석으로 예상됩니다.\n탑승 여부에 변동이 있으신 분은 미리 연락 부탁드립니다."
	applied := ApplyOverrides(base, map[string]string{"charter_bus_notice": notice})
	if applied.CharterBusNotice != notice {
		t.Errorf("CharterBusNotice = %q, want %q", applied.CharterBusNotice, notice)
	}

	diff := DiffToOverrideMap(base, applied)
	if diff["charter_bus_notice"] != notice {
		t.Errorf("diff[charter_bus_notice] = %q, want %q", diff["charter_bus_notice"], notice)
	}

	if _, ok := DiffToOverrideMap(base, base)["charter_bus_notice"]; ok {
		t.Error("diff of unchanged config should not contain charter_bus_notice")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/config/ -run TestCharterBusNoticeOverride`
Expected: FAIL to compile with `base.CharterBusNotice undefined`

- [ ] **Step 3: Implement the field**

In `internal/config/config.go`, add to the `WeddingConfig` struct directly after the `CharterBus` field (~line 102):

```go
	CharterBusNotice       string             `json:"charter_bus_notice"`
```

(The Go zero value covers the empty default in `loadWeddingDefaults()` — no line needed there.)

In `ApplyOverrides`, directly after the `charter_bus` unmarshal block (the `if v, ok := overrides["charter_bus"]; ok { ... }` block, ~line 317), add:

```go
	str("charter_bus_notice", &c.CharterBusNotice)
```

In `DiffToOverrideMap`, next to the other `strDiff` calls (after `strDiff("car_info", ...)`, ~line 361), add:

```go
	strDiff("charter_bus_notice", base.CharterBusNotice, current.CharterBusNotice)
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/config/ ./internal/handlers/ ./internal/models/`
Expected: `ok` for all three packages

- [ ] **Step 5: Commit** (see commit note in header — confirm with user first)

```bash
git add internal/config/config.go internal/config/config_test.go
git commit -m "feat: add configurable charter bus notice to wedding config"
```

---

### Task 2: Visitor frontend — hydrate and render the notice

**Files:**
- Modify: `frontend/src/config/wedding.ts` (~line 90 and ~line 179)
- Modify: `frontend/src/views/MapView.vue` (lines 75–90 template, ~line 162 imports)

- [ ] **Step 1: Add the config export**

In `frontend/src/config/wedding.ts`, directly after the `CHARTER_BUS` declaration (ends ~line 90), add:

```ts
export let CHARTER_BUS_NOTICE = '';
```

- [ ] **Step 2: Hydrate in loadConfig**

In the same file's `loadConfig()`, directly after `if (c.charter_bus) CHARTER_BUS = c.charter_bus;` (~line 179), add:

```ts
    if (c.charter_bus_notice) CHARTER_BUS_NOTICE = c.charter_bus_notice;
```

- [ ] **Step 3: Replace the hardcoded notice in MapView.vue**

In `frontend/src/views/MapView.vue`, replace the entire hardcoded notice block (lines 75–90):

```html
        <div class="parchment-bg bus-notice p-3 mt-4 text-left">
          <p class="text-secondary text-sm font-bold mb-2 text-center">
            <TwEmoji emoji="📢" size="1rem" /> 전세버스 탑승 안내
          </p>
          <p class="text-wood-dark/85 text-xs break-keep leading-relaxed mb-1.5">
            많은 축하와 성원에 감사드립니다.
          </p>
          <p class="text-wood-dark text-xs break-keep leading-relaxed mb-1.5">
            신부 측 전세버스는 사전 확인 결과
            <strong class="text-primary font-bold">좌석이 만석으로 예상</strong>됩니다.
          </p>
          <p class="text-wood-dark/85 text-xs break-keep leading-relaxed">
            정원제 운행으로 당일 현장에서 추가 탑승이 어려울 수 있으니,
            탑승 여부에 변동이 있으신 분은 미리 연락 부탁드립니다.
          </p>
        </div>
```

with:

```html
        <div v-if="CHARTER_BUS_NOTICE" class="parchment-bg bus-notice p-3 mt-4 text-left">
          <p class="text-secondary text-sm font-bold mb-2 text-center">
            <TwEmoji emoji="📢" size="1rem" /> 전세버스 탑승 안내
          </p>
          <p class="text-wood-dark text-xs break-keep leading-relaxed whitespace-pre-line">{{ CHARTER_BUS_NOTICE }}</p>
        </div>
```

In the same file's `<script setup>` import list from `@/config/wedding` (the one containing `CHARTER_BUS,` ~line 162), add `CHARTER_BUS_NOTICE,` on the line after `CHARTER_BUS,`.

- [ ] **Step 4: Run visitor frontend tests**

Run: `cd frontend && npm run test`
Expected: all tests pass (no existing test covers this notice; this guards against regressions)

- [ ] **Step 5: Commit** (see commit note in header)

```bash
git add frontend/src/config/wedding.ts frontend/src/views/MapView.vue
git commit -m "feat: render charter bus notice from config in map view"
```

---

### Task 3: Admin frontend — settings field

**Files:**
- Modify: `admin-frontend/src/types/admin.ts` (~line 86)
- Modify: `admin-frontend/src/views/SettingsView.vue` (form ~line 93, template ~line 849)

- [ ] **Step 1: Extend the WeddingConfig type**

In `admin-frontend/src/types/admin.ts`, directly after `charter_bus: CharterBusEntry[];` (~line 86), add:

```ts
  charter_bus_notice: string;
```

- [ ] **Step 2: Add the form default**

In `admin-frontend/src/views/SettingsView.vue`, in the `reactive` form object directly after the `charter_bus: [] as ...` line (~line 93), add:

```ts
  charter_bus_notice: '',
```

(`populateForm` spreads `cfg` into the form, so no loader change is needed.)

- [ ] **Step 3: Add the textarea to the template**

In the same file, the Charter Bus section ends with the `v-for="(cb, i) in form.charter_bus"` block's closing `</div>` followed by another `</div>` (~line 849), before a `<Separator />`. Insert between that closing `</div>` and the `<Separator />`:

```html
              <div>
                <label class="mb-1 block text-sm font-medium"
                  >Charter Bus Notice</label
                >
                <textarea
                  v-model="form.charter_bus_notice"
                  rows="4"
                  placeholder="전세버스 탑승 안내 문구를 입력하세요"
                  class="flex w-full rounded-md border border-input bg-transparent px-3 py-2 text-sm shadow-sm placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
                />
                <p class="mt-1 text-xs text-muted-foreground">
                  Shown under the charter bus list. Leave empty to hide the notice.
                </p>
              </div>
```

(This mirrors the existing `main_greet_text` raw-`<textarea>` pattern at ~line 497 — the project does not use a `Textarea` component.)

- [ ] **Step 4: Run admin frontend tests**

Run: `cd admin-frontend && npm run test`
Expected: all tests pass

- [ ] **Step 5: Commit** (see commit note in header)

```bash
git add admin-frontend/src/types/admin.ts admin-frontend/src/views/SettingsView.vue
git commit -m "feat: add charter bus notice field to admin settings"
```

---

### Task 4: End-to-end verification

**Files:** none (verification only)

- [ ] **Step 1: Full backend test suite**

Run: `go test ./internal/...`
Expected: `ok` for every package

- [ ] **Step 2: Production build**

Run: `make build`
Expected: builds both frontends and the `wedding-server` binary without errors (this also type-checks the Vue code via `vue-tsc`)

- [ ] **Step 3: Verify the API exposes the field**

```bash
DATABASE_PATH=/tmp/claude-charter-verify.db ./wedding-server &
sleep 1
curl -s localhost:8080/api/config | grep -o '"charter_bus_notice":""'
kill %1
rm -f /tmp/claude-charter-verify.db*
```

Expected: `"charter_bus_notice":""` is printed (field present, empty by default)

- [ ] **Step 4: Manual smoke test (empty vs set)**

Run `make dev`, open the admin settings (`localhost:5174`), confirm the Charter Bus Notice textarea appears in the transport section, save a multi-line notice, then open the visitor map view (`localhost:5173/map`) and confirm the 📢 notice panel renders the text with line breaks. Clear the field, save, reload — the panel must disappear (the charter bus list itself stays).
