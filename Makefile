VERSION    := $(shell cat VERSION)
IMAGE_NAME ?= wedding-server
IMAGE_TAG  ?= $(VERSION)

.PHONY: all build build-frontend build-admin build-backend clean dev dev-frontend dev-admin dev-backend install install-admin help docker-build docker-run docker-up docker-down version sync-version test test-backend test-frontend test-admin coverage coverage-go coverage-frontend coverage-admin generate-og

all: build

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  all              Build everything (default)"
	@echo "  install          Install all frontend dependencies"
	@echo "  install-admin    Install admin frontend dependencies"
	@echo "  build            Build frontend, admin, and backend"
	@echo "  build-frontend   Build Vue.js visitor frontend"
	@echo "  build-admin      Build Vue.js admin frontend"
	@echo "  build-backend    Build Go backend (includes both frontends)"
	@echo "  clean            Remove build artifacts"
	@echo "  dev              Run all dev servers in parallel (Ctrl+C to stop)"
	@echo "  dev-frontend     Run visitor frontend dev server (port 5173)"
	@echo "  dev-admin        Run admin frontend dev server (port 5174)"
	@echo "  dev-backend      Run backend dev server"
	@echo "  docker-build     Build Docker image (IMAGE_NAME=$(IMAGE_NAME), IMAGE_TAG=$(IMAGE_TAG))"
	@echo "  docker-run       Run Docker container"
	@echo "  docker-up        Build image, start compose, show admin password"
	@echo "  docker-down      Stop compose and remove volumes"
	@echo "  generate-og      Generate OG image from wedding config"
	@echo "  version          Print current version"
	@echo "  test             Run all tests (backend + frontend)"
	@echo "  test-backend     Run Go backend unit tests"
	@echo "  test-frontend    Run visitor frontend unit tests"
	@echo "  test-admin       Run admin frontend unit tests (if available)"
	@echo "  coverage         Run all coverage reports"
	@echo "  coverage-go      Generate Go HTML coverage report"
	@echo "  coverage-frontend  Show visitor frontend coverage"
	@echo "  coverage-admin   Show admin frontend coverage"
	@echo "  sync-version     Sync VERSION to frontend/package.json"
	@echo "  help             Show this help message"

test: test-backend test-frontend test-admin

test-backend:
	go test ./internal/... -v

test-frontend:
	cd frontend && npm run test

test-admin:
	cd admin-frontend && npm run test

coverage: coverage-go coverage-frontend coverage-admin

coverage-go:
	@echo "── Go Backend Coverage ──"
	go test ./internal/... -coverprofile=coverage.out
	@go tool cover -func=coverage.out | tail -1
	go tool cover -html=coverage.out -o coverage.html
	@echo "HTML report: coverage.html"

coverage-frontend:
	@echo "── Visitor Frontend Coverage ──"
	cd frontend && npx vitest run --coverage 2>&1 || true

coverage-admin:
	@echo "── Admin Frontend Coverage ──"
	cd admin-frontend && npx vitest run --coverage 2>&1 || true

install:
	cd frontend && npm install
	cd admin-frontend && npm install

install-admin:
	cd admin-frontend && npm install

build: build-frontend build-admin build-backend

sync-version:
	cd frontend && npm version $(VERSION) --no-git-tag-version --allow-same-version

generate-og:
	cd frontend && npx tsx scripts/generate-og.ts

build-frontend: sync-version generate-og
	cd frontend && npm run build

build-admin:
	cd admin-frontend && npm run build

build-backend: build-frontend build-admin
	CGO_ENABLED=0 go build -tags timetzdata -ldflags "-s -w -X main.Version=$(VERSION)" -o wedding-server .

clean:
	rm -f wedding-server
	rm -rf web/dist/*
	touch web/dist/.gitkeep
	rm -rf web/admin/*
	touch web/admin/.gitkeep

dev-frontend:
	cd frontend && npm run dev -- --host

dev-admin:
	cd admin-frontend && npm run dev -- --host

dev-backend:
	go run main.go

dev:
	@trap 'kill 0' INT TERM; \
	(cd frontend && npm run dev) & \
	(cd admin-frontend && npm run dev) & \
	go run main.go & \
	wait

version:
	@echo $(VERSION)

docker-build:
	docker build --build-arg VERSION=$(VERSION) -t $(IMAGE_NAME):$(IMAGE_TAG) -t $(IMAGE_NAME):latest .

docker-run: docker-build
	docker run -d -p 8080:8080 -v wedding-data:/data $(IMAGE_NAME):$(IMAGE_TAG)

docker-up: docker-build
	WEDDING_IMAGE=$(IMAGE_NAME):latest docker compose up -d
	@sleep 2
	@docker compose logs wedding-server 2>/dev/null | grep -A 1 "AUTO-GENERATED ADMIN PASSWORD" || echo "Password already changed or not available"

docker-down:
	docker compose down -v
