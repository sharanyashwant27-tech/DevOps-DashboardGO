# DevOps Command Center

Enterprise DevOps dashboard for CI/CD, GitHub, Docker, Kubernetes, servers, deployments, incidents, and alerts — with multi-organization support.

![Go](https://img.shields.io/badge/Go-1.25-00ADD8) ![React](https://img.shields.io/badge/React-18-61DAFB) ![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-336791) ![Docker](https://img.shields.io/badge/Docker-Compose-2496ED) ![Kubernetes](https://img.shields.io/badge/Kubernetes-ready-326CE5)

**Repository:** [sharanyashwant27-tech/DevOps-DashboardGO](https://github.com/sharanyashwant27-tech/DevOps-DashboardGO)  
**Default app URL:** [http://localhost:8095](http://localhost:8095)  
**Local image:** `devops-dashboard-go:latest` (v1.1.0)

---

## Features

| Area | Capabilities |
|------|----------------|
| **Auth** | Register, login, forgot password, JWT + refresh, RBAC (`admin`, `devops`, `developer`, `viewer`) |
| **Dashboard** | Colorful KPI cards → modules, Chart.js, WebSocket live updates |
| **Jenkins** | Jobs, queue, builds, trigger/stop, console, filters (+ demo mode) |
| **GitHub** | Live repos (username and/or PAT), commits, PRs, Actions, health |
| **Docker** | Live containers/images/volumes/networks via host Docker socket |
| **Kubernetes** | Pods, deployments, services, nodes, events, scale/restart/logs |
| **Servers** | Host CPU/RAM/disk/load/processes (gopsutil) |
| **Deployments** | History + rollback |
| **Incidents / Alerts** | Priority, SLA, comments, ack/resolve/mute |
| **Ops** | Zap logs, Prometheus `/metrics`, Grafana + Loki, audit trail |

---

## Clone

```bash
git clone https://github.com/sharanyashwant27-tech/DevOps-DashboardGO.git
cd DevOps-DashboardGO
```

---

## Quick start (Docker)

```bash
# 1) Configure secrets (GitHub optional but recommended)
cp .env.example .env
# Edit .env:
#   DCC_GITHUB_USERNAME=your-github-username
#   DCC_GITHUB_TOKEN=your_pat   # optional for private repos

# 2) Build & run stack (API + UI on 8095, Postgres, Redis, …)
docker compose -f deployments/docker-compose.yml up -d --build

# 3) Open the app
# http://localhost:8095
```

| Surface | URL |
|---------|-----|
| **App (UI + API)** | http://localhost:8095 |
| Health | http://localhost:8095/health |
| Swagger | http://localhost:8095/swagger/index.html |
| Grafana | http://localhost:3001 (`admin` / `admin`) |
| Prometheus | http://localhost:9090 |
| Optional nginx UI | http://localhost:3000 |

**Seeded admin:** `admin@devops.local` / `Admin@12345`

### Rebuild the app image

```bash
docker compose -f deployments/docker-compose.yml build backend
docker compose -f deployments/docker-compose.yml up -d --force-recreate backend
```

Standalone image build (mount the Docker socket for Operations Console → Docker Monitoring):

```bash
docker build -f deployments/Dockerfile.backend -t devops-dashboard-go:latest .
docker run --rm -p 8095:8095 \
  --env-file .env \
  -e DCC_DATABASE_HOST=host.docker.internal \
  -e DCC_DATABASE_PORT=5434 \
  -e DCC_REDIS_HOST=host.docker.internal \
  -e DCC_REDIS_PORT=6380 \
  -e DCC_DOCKER_ENABLED=true \
  -e DCC_DOCKER_HOST=unix:///var/run/docker.sock \
  -e DOCKER_HOST=unix:///var/run/docker.sock \
  -v /var/run/docker.sock:/var/run/docker.sock \
  --group-add 0 \
  devops-dashboard-go:latest
```

The **backend image** (`devops-dashboard-go:latest`) is multi-stage:

1. Builds the React UI (`node:20`) — aurora ops theme (cyan / amber / coral accents)
2. Compiles the Go 1.25 API
3. Serves UI + API from Alpine on port **8095** (non-root, healthcheck)

### UI look & feel (v1.1.0)

- Mesh gradient background + soft grid overlay
- Color-coded nav, KPI cards, status chips, and page heroes
- Glass panels, gradient primary buttons, auth screens with floating color orbs
- Dark / light modes with a cyan–amber operations palette (Sora + IBM Plex)

CI pushes the same Dockerfile to GHCR on `main`:

- `ghcr.io/sharanyashwant27-tech/devops-command-center-backend:latest`
- `ghcr.io/sharanyashwant27-tech/devops-command-center-frontend:latest`

---

## Docker Monitoring (Operations Console)

Compose mounts the **host Docker engine** into `dcc-backend` so the Docker page lists real containers:

| Setting | Value |
|---------|--------|
| Socket volume | `/var/run/docker.sock:/var/run/docker.sock` |
| `DCC_DOCKER_HOST` / `DOCKER_HOST` | `unix:///var/run/docker.sock` |
| `group_add` | `0` (socket access for non-root `appuser`) |

**Requirements**

- Docker Desktop (or a Linux Docker daemon) must be running
- On Docker Desktop: enable use of the default Docker socket (Advanced settings)
- Health check should show `"docker":{"reachable":true}` at `/health`

Without the socket mount you will see **Docker daemon unavailable**.

---

## Configuration

| Source | Purpose |
|--------|---------|
| `backend/config/config.yaml` | Defaults |
| `.env` / `DCC_*` env vars | Secrets & overrides (loaded by Compose `env_file`) |

Important variables:

```env
DCC_GITHUB_USERNAME=sharanyashwant27-tech
DCC_GITHUB_TOKEN=                 # PAT from https://github.com/settings/tokens
DCC_JENKINS_URL=
DCC_JENKINS_USERNAME=
DCC_JENKINS_TOKEN=
DCC_DOCKER_ENABLED=true
DCC_DOCKER_HOST=unix:///var/run/docker.sock
DCC_JWT_ACCESS_SECRET=change-me-access-secret-min-32-chars!!
DCC_JWT_REFRESH_SECRET=change-me-refresh-secret-min-32-chars!
```

- **GitHub without token:** public repos for `DCC_GITHUB_USERNAME`  
- **GitHub with token:** authenticated (private repos + higher rate limits)  
- **Jenkins without URL:** built-in demo jobs  
- **Docker:** requires host socket mounted (see above)

Never commit `.env` or PATs.

---

## Architecture

```
Browser → :8095 (Gin: static UI + /api/v1 + /ws)
              → PostgreSQL / Redis
              → Jenkins / GitHub / Docker socket / Kubernetes
              → Prometheus ← /metrics
```

Clean Architecture layout lives under `backend/internal/` (controllers → services → repositories).  
See [docs/FOLDER_STRUCTURE.md](docs/FOLDER_STRUCTURE.md) and [docs/ER_DIAGRAM.md](docs/ER_DIAGRAM.md).

---

## Development roadmap

Phased delivery checklist (status tracked in detail):

| Phase | Module |
|------:|--------|
| 1 | Authentication & user management |
| 2 | Dashboard with real-time WebSocket updates |
| 3 | Jenkins integration |
| 4 | GitHub integration |
| 5 | Docker monitoring |
| 6 | Kubernetes monitoring |
| 7 | Server metrics (Prometheus / Node Exporter) |
| 8 | Deployment history & rollback |
| 9 | Incident management |
| 10 | Alert dashboard & notifications |
| 11 | CI/CD, testing, and production deployment |

Full status and next actions: **[docs/ROADMAP.md](docs/ROADMAP.md)**

---

## Local development (without full stack image)

```bash
docker compose -f deployments/docker-compose.yml up -d postgres redis

cd backend && go mod tidy && go run ./cmd/server
# API + (after build) UI on http://localhost:8095
# On Linux/macOS with Docker: DOCKER_HOST=unix:///var/run/docker.sock

cd frontend && npm install && npm run build
# or: npm run dev  → http://localhost:3000 (proxies API to 8095)
```

Windows note: if Application Control blocks `go run`, use the Docker backend image above (with the socket volume).

---

## Docker layout

| File | Role |
|------|------|
| `deployments/Dockerfile.backend` | Multi-stage API + embedded UI → `devops-dashboard-go:latest` on `:8095` |
| `deployments/Dockerfile.frontend` | Optional nginx SPA → `:80` |
| `deployments/docker-compose.yml` | Full stack + Docker socket mount for live monitoring |
| `.dockerignore` | Keeps build context small / excludes secrets |

Compose host ports (avoid local conflicts):

- App **8095**, Postgres **5434**, Redis **6380**

---

## Documentation

| Doc | Description |
|-----|-------------|
| [docs/ROADMAP.md](docs/ROADMAP.md) | Phase status & next work |
| [docs/API.md](docs/API.md) | REST API reference |
| [docs/DEPLOYMENT.md](docs/DEPLOYMENT.md) | Docker / K8s / HTTPS / secrets |
| [docs/ER_DIAGRAM.md](docs/ER_DIAGRAM.md) | Database ER model |
| [docs/FOLDER_STRUCTURE.md](docs/FOLDER_STRUCTURE.md) | Package map |
| [docs/SCREENSHOTS.md](docs/SCREENSHOTS.md) | UI capture guide |

---

## Testing & CI/CD

```bash
cd backend && go test ./...
cd frontend && npm run build
```

GitHub Actions (`.github/workflows/ci.yml`): Go 1.25 test → frontend build → Docker image build/push to GHCR → Kubernetes deploy hook.

---

## License

MIT — adapt for your organization.
