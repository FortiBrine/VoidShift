SHELL := /bin/sh

FRONTEND_DIR := frontend
BACKEND_DIR := backend
FRONTEND_OUT := $(FRONTEND_DIR)/.output/public
EMBED_DIR := $(BACKEND_DIR)/internal/embed/webui
APP_OUT := app

.PHONY: all frontend backend build clean deps

all: build

build: frontend backend

deps:
	cd $(FRONTEND_DIR) && bun install --frozen-lockfile
	cd $(BACKEND_DIR) && go mod download && go mod verify

frontend:
	cd $(FRONTEND_DIR) && bun install --frozen-lockfile
	cd $(FRONTEND_DIR) && bun run generate

backend: frontend
	rm -rf $(EMBED_DIR)
	mkdir -p $(EMBED_DIR)
	cp -R $(FRONTEND_OUT)/. $(EMBED_DIR)/
	cd $(BACKEND_DIR) && CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-w -s" -o ../$(APP_OUT) ./cmd/api

clean:
	rm -rf $(FRONTEND_DIR)/.output
	rm -rf $(EMBED_DIR)
	rm -f $(APP_OUT)
