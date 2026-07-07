FROM docker.io/oven/bun:1.3.10-alpine AS assets
WORKDIR /app

COPY package.json bun.lock ./
RUN --mount=type=cache,target=/root/.bun/install/cache \
    bun install --frozen-lockfile

COPY assets ./assets
COPY view ./view
RUN bunx @tailwindcss/cli -i assets/app.css -o app.css

FROM golang:1.26.1-alpine3.23 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download && go mod verify

COPY . .
COPY --from=assets /app/app.css ./internal/webui/static/app.css

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go tool sqlc generate && \
    go tool templ generate && \
    go tool go-licenses check ./cmd/api --disallowed_types=forbidden,restricted,unknown && \
    CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-w -s" -o /app/app ./cmd/api && \
    go tool go-licenses report ./cmd/api > /app/third-party-licenses.csv

FROM alpine:3.21

RUN --mount=type=cache,target=/var/cache/apk,sharing=locked \
    apk add --update \
    ca-certificates \
    wireguard-tools \
    iproute2 \
    iptables \
    kmod

COPY --from=builder /app/app /app
COPY --from=builder /app/third-party-licenses.csv /licenses/third-party-licenses.csv
COPY LICENSE-MIT LICENSE-APACHE /licenses/

EXPOSE 8080/tcp 51820/udp

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget -q -O- http://localhost:8080/api/health

ENTRYPOINT ["/app"]
