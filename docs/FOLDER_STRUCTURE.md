# Folder-by-Folder Explanation

## Root

| Path | Role |
|------|------|
| `backend/` | Go API (Clean Architecture) |
| `frontend/` | React + TypeScript UI |
| `deployments/` | Docker, NGINX, Prometheus, Grafana, K8s |
| `docs/` | Architecture, API, ER, deployment guides |
| `scripts/` | Helper scripts |
| `.github/workflows/` | CI/CD |

## Backend (idiomatic `internal/` layout)

Requested top-level folders map as follows:

| Requested | Implemented |
|-----------|-------------|
| `cmd/` | `backend/cmd/server` |
| `config/` | `backend/config` |
| `controllers/` | `backend/internal/controllers` |
| `models/` | `backend/internal/models` |
| `dto/` | `backend/internal/dto` |
| `repositories/` | `backend/internal/repositories` |
| `services/` | `backend/internal/services` |
| `middleware/` | `backend/internal/middleware` |
| `routes/` | `backend/internal/routes` |
| `websocket/` | `backend/internal/websocket` |
| `scheduler/` | `backend/internal/scheduler` |
| `auth/` | `backend/internal/auth` |
| `migrations/` | `backend/migrations` |
| `tests/` | `backend/tests` |
| `utils/` / helpers | `backend/pkg/*` |

### Layer responsibilities

1. **Controllers** – HTTP binding, status codes, no business rules  
2. **Services** – business logic, integrations (Jenkins/GitHub/Docker/K8s)  
3. **Repositories** – GORM persistence behind interfaces  
4. **Models / DTOs** – domain entities vs transport contracts  
5. **Middleware** – JWT, RBAC, CORS, rate limit, security headers, logging  
6. **WebSocket hub** – live dashboard/alerts/metrics fan-out  
7. **Scheduler** – metrics collection + alert evaluation cron  

## Frontend

| Path | Role |
|------|------|
| `src/pages/` | Route-level screens |
| `src/layouts/` | Shell / navigation |
| `src/hooks/` | WebSocket and shared hooks |
| `src/services/` | Axios API clients |
| `src/context/` | Auth + theme providers |
| `src/components/*` | Module UI building blocks |
| `src/types/` | Shared TypeScript types |

## Deployments

| Path | Role |
|------|------|
| `docker-compose.yml` | Full local/prod-like stack |
| `Dockerfile.*` | Multi-stage builds |
| `nginx/` | Reverse proxy + SPA |
| `prometheus/` | Scrape config |
| `grafana/` | Datasource provisioning |
| `kubernetes/` | Namespace, Deployments, Ingress, RBAC |
