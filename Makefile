.PHONY: help backend-run frontend-run test tidy compose-up compose-down migrate seed docker-build

help:
	@echo "DevOps Command Center targets:"
	@echo "  make tidy          - go mod tidy"
	@echo "  make test          - run Go tests"
	@echo "  make backend-run   - run API server"
	@echo "  make frontend-run  - run Vite UI"
	@echo "  make compose-up    - start full stack"
	@echo "  make compose-down  - stop stack"
	@echo "  make docker-build  - build images"

tidy:
	cd backend && go mod tidy

test:
	cd backend && go test ./...

backend-run:
	cd backend && go run ./cmd/server

frontend-run:
	cd frontend && npm install && npm run dev

compose-up:
	docker compose -f deployments/docker-compose.yml up -d --build

compose-down:
	docker compose -f deployments/docker-compose.yml down

docker-build:
	docker build -f deployments/Dockerfile.backend -t dcc-backend:latest .
	docker build -f deployments/Dockerfile.frontend -t dcc-frontend:latest .
