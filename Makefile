SHELL := /bin/sh

APP_OUT := app
CSS_IN  := assets/app.css
CSS_OUT := internal/webui/static/app.css

.PHONY: all generate build run css css-watch clean

all: generate css build

generate:
	go tool sqlc generate
	go tool templ generate

css:
	bunx @tailwindcss/cli -i $(CSS_IN) -o $(CSS_OUT)

css-watch:
	bunx @tailwindcss/cli -i $(CSS_IN) -o $(CSS_OUT) --watch

build:
	go build -o $(APP_OUT) ./cmd/api

run: build
	sudo ./$(APP_OUT)

clean:
	rm -f $(APP_OUT)
	rm -f $(CSS_OUT)
	rm -f internal/store/models.go internal/store/db.go internal/store/*.sql.go
	find view -name '*_templ.go' -delete
