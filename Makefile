.PHONY: all backend frontend

VERSION := $(shell git describe --tags --exact-match 2>/dev/null || git rev-parse --short HEAD 2>/dev/null || echo "dev")

all: backend frontend

backend:
	cd backend && go build -ldflags "-X main.Version=$(VERSION)" -o server ./cmd/server

frontend:
	cd frontend && npm install && VITE_APP_VERSION=$(VERSION) npm run build

run-backend:
	cd backend && go run ./cmd/server

run-frontend:
	cd frontend && npm run dev

test-backend:
	cd backend && go test ./...

test-frontend:
	@echo "No frontend tests"

lint-frontend:
	cd frontend && npm run lint
