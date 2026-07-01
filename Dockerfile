FROM docker.io/oven/bun:1.3.10-alpine AS css
WORKDIR /app

COPY package.json bun.lock ./
RUN bun install --frozen-lockfile

COPY assets ./assets
RUN bunx @tailwindcss/cli -i assets/app.css -o app.css

FROM golang:1.26.1-alpine3.23 AS backend
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
COPY --from=css /app/app.css ./internal/webui/static/app.css

RUN go tool sqlc generate
RUN go tool templ generate
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-w -s" -o /app/app ./cmd/api

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
