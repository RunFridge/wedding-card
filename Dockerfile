# Stage 1: Build visitor frontend
FROM --platform=$BUILDPLATFORM node:22-alpine AS frontend
WORKDIR /app/frontend
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# Stage 2: Build admin frontend
FROM --platform=$BUILDPLATFORM node:22-alpine AS admin
WORKDIR /app/admin-frontend
COPY admin-frontend/package.json admin-frontend/package-lock.json ./
RUN npm ci
COPY admin-frontend/ ./
RUN npm run build

# Stage 3: Build Go binary (cross-compiled on the native builder)
FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS backend
ARG VERSION=dev
ARG TARGETOS
ARG TARGETARCH
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend /app/web/dist ./web/dist
COPY --from=admin /app/web/admin ./web/admin
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -tags timetzdata -ldflags="-s -w -X main.Version=${VERSION}" -o /wedding-server .

# Stage 4: Grab CA certificates from a tiny base
FROM --platform=$BUILDPLATFORM alpine:3.21 AS certs
RUN apk add --no-cache ca-certificates

# Stage 5: Minimal runtime
FROM scratch
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=backend /wedding-server /wedding-server

ENV DATABASE_PATH=/data/wedding.db

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD ["/wedding-server", "-healthcheck"]

VOLUME ["/data"]
EXPOSE 8080

ENTRYPOINT ["/wedding-server"]
