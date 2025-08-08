SHELL := /bin/bash

.PHONY: run test lint fmt ci docker-build docker-run dev-up dev-down migrate-up migrate-down

run:
	go run ./cmd/server

test:
	go test ./... -race -v

lint:
	go vet ./...
	@gofmt -l . | (! grep .) || (echo "Run 'gofmt -w .' to format the files above" && exit 1)

fmt:
	gofmt -s -w .

ci: lint test

docker-build:
	docker build -t go-todo-api:local .

docker-run:
	docker run --rm -p 8080:8080 -e PORT=8080 go-todo-api:local

dev-up:
	docker compose up --build

dev-down:
	docker compose down -v

DB_DSN ?= host=localhost port=5432 user=${POSTGRES_USER} password=${POSTGRES_PASSWORD} dbname=${POSTGRES_DB} sslmode=disable

migrate-up:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING='${DB_DSN}' go run github.com/pressly/goose/v3/cmd/goose@v3.24.3 -dir ./migrations up

migrate-down:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING='${DB_DSN}' go run github.com/pressly/goose/v3/cmd/goose@v3.24.3 -dir ./migrations down


