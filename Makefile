SHELL := /bin/sh

APP_OUT    := app
CSS_IN     := assets/app.css
STATIC_DIR := internal/webui/static
CSS_OUT    := $(STATIC_DIR)/app.css

.PHONY: build install generate assets run clean

build: generate assets
	CGO_ENABLED=0 go build -trimpath -ldflags="-w -s" -o $(APP_OUT) ./cmd/api

install:
	bun install

generate:
	go tool sqlc generate
	go tool templ generate

assets: install
	bunx @tailwindcss/cli -i $(CSS_IN) -o $(CSS_OUT)

run: build
	sudo ./$(APP_OUT)

clean:
	rm -f $(APP_OUT)
	find $(STATIC_DIR) -mindepth 1 -delete
	rm -f internal/store/models.go internal/store/db.go internal/store/*.sql.go
	find view -name '*_templ.go' -delete
