# Makefile

.PHONY: all backend frontend

all: backend frontend

backend:
	cd backend && go build -o server ./cmd/server

frontend:
	cd frontend && npm install && npm run build

run-backend:
	cd backend && go run ./cmd/server

run-frontend:
	cd frontend && npm run dev

test-backend:
	cd backend && go test ./...

test-frontend:
	@echo "No frontend tests"
