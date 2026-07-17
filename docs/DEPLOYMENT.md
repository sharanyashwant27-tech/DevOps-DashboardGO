# Deployment Guide

## Prerequisites

- Docker & Docker Compose
- Go 1.22+ (local backend)
- Node.js 20+ (local frontend)
- Optional: Kubernetes cluster, kubectl
- Optional: Jenkins URL + token, GitHub PAT

## Quick Start (Docker Compose)

```bash
cp .env.example .env
docker compose -f deployments/docker-compose.yml up -d --build
```

Services:

| Service | URL |
|---------|-----|
| UI (frontend) | http://localhost:3000 |
| API | http://localhost:8080 |
| NGINX gateway | http://localhost:80 |
| Swagger | http://localhost:8080/swagger/index.html |
| Prometheus | http://localhost:9090 |
| Grafana | http://localhost:3001 (admin/admin) |
| Loki | http://localhost:3100 |

Default admin (seeded):

- Email: `admin@devops.local`
- Password: `Admin@12345`

## Local Development

```bash
# Infra only
docker compose -f deployments/docker-compose.yml up -d postgres redis

# Backend
cd backend
go mod tidy
go run ./cmd/server

# Frontend
cd frontend
npm install
npm run dev
```

## Kubernetes

```bash
kubectl apply -f deployments/kubernetes/namespace.yaml
kubectl apply -f deployments/kubernetes/secrets.yaml
kubectl apply -f deployments/kubernetes/configmap.yaml
kubectl apply -f deployments/kubernetes/backend-deployment.yaml
kubectl apply -f deployments/kubernetes/frontend-deployment.yaml
```

Update image names in manifests to your registry, then configure Ingress DNS (`devops.local` by default).

## HTTPS

Terminate TLS at NGINX Ingress / external load balancer. The API is HTTPS-ready behind reverse proxy headers (`X-Forwarded-Proto`).

## Secrets Management

Prefer environment variables / Kubernetes Secrets for:

- JWT secrets
- Database password
- GitHub PAT
- Jenkins token
- Slack / Teams webhooks

Never commit `.env` or real secrets.

## CI/CD

GitHub Actions workflow (`.github/workflows/ci.yml`) runs:

1. Go vet + tests + build
2. Frontend build
3. Docker image build/push (main)
4. Kubernetes deploy hook (configure `KUBECONFIG`)

## Observability

- App exposes Prometheus metrics at `/metrics`
- Grafana is pre-provisioned with Prometheus + Loki datasources
- Zap structured logs to stdout (collect with Loki/Promtail in production)
