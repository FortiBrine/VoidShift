FROM docker.io/oven/bun:1.3.10-alpine AS frontend
WORKDIR /app/frontend

COPY frontend/package.json frontend/bun.lock ./
RUN bun install --frozen-lockfile

COPY frontend/ ./
RUN bun run generate

FROM golang:1.26.1-alpine3.23 AS backend
WORKDIR /app/backend

COPY backend/go.mod backend/go.sum ./
RUN go mod download && go mod verify

COPY backend/ ./

COPY --from=frontend /app/frontend/.output/public ./internal/embed/webui

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/app ./cmd/api

FROM alpine:3.21

RUN apk add --no-cache \
    ca-certificates \
    wireguard-tools \
    iproute2 \
    iptables \
    kmod

COPY --from=backend /app/app /app

EXPOSE 8080/tcp
EXPOSE 51820/udp

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget -q -O- http://localhost:8080/api/health

ENTRYPOINT ["/app"]
