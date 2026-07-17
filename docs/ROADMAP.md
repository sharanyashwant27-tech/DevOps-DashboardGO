# Development Roadmap — DevOps Command Center

Phased plan to keep delivery manageable. Status reflects the current codebase.

| Status | Meaning |
|--------|---------|
| Done | Core feature shipped and usable |
| Partial | Scaffold / demo / basic UI — needs production hardening |
| Next | Recommended focus for upcoming work |

---

## Phase overview

| Phase | Module | Status | Notes |
|------:|--------|--------|-------|
| 1 | Authentication & user management | **Done** | JWT + refresh, register/login/forgot-password, RBAC (`admin` / `devops` / `developer` / `viewer`), seed admin |
| 2 | Dashboard + real-time WebSocket | **Partial** | KPI cards, charts, click-through nav, WS hub + live stats push; deepen channel subscriptions (builds/logs) |
| 3 | Jenkins integration | **Partial** | Live REST client + **demo mode**; trigger/stop/console/queue/filters |
| 4 | GitHub integration | **Done** (live public/auth) | Live mode via PAT / public username (`sharanyashwant27-tech`); repos, commits, PRs, Actions, health |
| 5 | Docker monitoring | **Done** | Docker SDK list/lifecycle/logs/stats; Compose mounts host `/var/run/docker.sock` into `dcc-backend` |
| 6 | Kubernetes monitoring | **Partial** | client-go list/scale/restart/logs; enable `kubernetes.enabled` + kubeconfig / in-cluster |
| 7 | Server metrics (Prometheus / Node Exporter) | **Partial** | Local gopsutil + Prometheus scrape + Grafana/Loki stack; add Node Exporter + PromQL dashboards |
| 8 | Deployment history & rollback | **Done** (basic) | CRUD history, rollback API + UI button; wire to real CD pipelines |
| 9 | Incident management | **Done** (basic) | Create/list/update, comments, SLA deadline, status filters |
| 10 | Alert dashboard & notifications | **Partial** | Ack/resolve/mute + sources; Slack/Teams/email are stubs — implement real webhooks |
| 11 | CI/CD, testing, production deploy | **Partial** | GitHub Actions (test/build/push), Compose + K8s manifests; expand tests, harden secrets, TLS |

---

## Phase details & next actions

### 1 — Authentication & user management ✅
**Have:** JWT access/refresh, bcrypt, roles, audit on login, `/auth/*` APIs, login/register UI.  
**Next (optional):** email password-reset, admin user CRUD UI, org membership invites.

### 2 — Dashboard with real-time WebSocket 🟨
**Have:** Aggregated stats, animated cards → modules, Chart.js, `/ws` hub, scheduler publishes dashboard/metrics.  
**Next:** subscribe UI to `builds` / `alerts` channels; live container/pod log streaming; reconnect UX.

### 3 — Jenkins integration 🟨
**Have:** Full REST surface + demo jobs when `jenkins.url` empty.  
**Next:** set `DCC_JENKINS_*`, persist builds into DB for dashboard KPIs, average build time chart.

### 4 — GitHub integration ✅
**Have:** Live repos for configured user/token; health/commits/PRs/Actions.  
**Next:** richer UI tables (not JSON dumps), contributors/releases tabs, webhook ingress for Actions status.

### 5 — Docker monitoring ✅
**Have:** Containers/images/volumes/networks, start/stop/restart/delete, logs/stats; Compose mounts host Docker socket into `dcc-backend`.  
**Next:** CPU/mem cards per container; multi-host `docker_hosts`.

### 6 — Kubernetes monitoring 🟨
**Have:** Namespaces, pods, deployments, services, nodes, events, scale/restart/delete, logs.  
**Next:** turn on in config, ClusterRole already in manifests; nicer tables; namespace-scoped RBAC.

### 7 — Server metrics (Prometheus / Node Exporter) 🟨
**Have:** Host metrics via gopsutil, `/metrics`, Compose Prometheus/Grafana/Loki.  
**Next:** add `node_exporter` service, scrape config, Grafana dashboard JSON with CPU/mem/disk/network panels.

### 8 — Deployment history & rollback ✅/🟨
**Have:** Store deployments, list UI, rollback endpoint.  
**Next:** capture rollback version automatically on deploy; link commit → GitHub; approval workflow.

### 9 — Incident management ✅/🟨
**Have:** Incidents, priority, assignee, timeline, SLA, comments.  
**Next:** attachments upload, SLA countdown UI, assign engineer picker from users API.

### 10 — Alert dashboard & notifications 🟨
**Have:** Multi-source alerts, severity filters, ack/resolve/mute, scheduler threshold alerts.  
**Next:** real Slack/Teams/SMTP sends; Prometheus Alertmanager webhook receiver; in-app notification center.

### 11 — CI/CD, testing, production deployment 🟨
**Have:** `.github/workflows/ci.yml`, Dockerfiles, Compose, K8s Deployments/Ingress/RBAC.  
**Next:** more unit/integration tests, image scanning, TLS ingress, secrets via sealed-secrets/SOPS, staging env.

---

## Suggested sequence from here

1. **Phase 7** — Node Exporter + Grafana panels (high demo value, low risk)  
2. **Phase 6** — Enable kubeconfig / in-cluster so Kubernetes leaves “unavailable”  
3. **Phase 2** — Richer WebSocket live feeds  
4. **Phase 10** — Real Slack/Teams notifications  
5. **Phase 11** — Test coverage + production hardening  

---

## How to track progress

- Keep this file updated when a phase moves Done → ship a short note under that phase.  
- Prefer vertical slices: API + UI + seed/demo for each module before polishing.  
- Never commit secrets (`.env`, PATs); use Compose `env_file` / K8s Secrets only.
