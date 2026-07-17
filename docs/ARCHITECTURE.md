# System Architecture

## High-level

```mermaid
flowchart LR
  UI[React UI] -->|REST / WS| API[Gin API]
  API --> Auth[JWT + RBAC]
  API --> SVC[Domain Services]
  SVC --> DB[(PostgreSQL)]
  SVC --> Redis[(Redis)]
  SVC --> Jenkins[Jenkins REST]
  SVC --> GitHub[GitHub API]
  SVC --> Docker[Docker SDK]
  SVC --> K8s[client-go]
  API --> Prom[/metrics]
  Prom --> Grafana
  API --> Zap[Zap Logs]
  Zap --> Loki
```

## Principles

- **SOLID** — service interfaces via repository contracts; controllers thin
- **Dependency Injection** — wired in `cmd/server/main.go`
- **Repository Pattern** — GORM behind interfaces for testability
- **Clean Architecture** — transport → application → domain → infrastructure

## Roles

| Role | Capabilities |
|------|--------------|
| admin | Full access + audit |
| devops | Infra actions (docker/k8s/jenkins stop, alert resolve) |
| developer | Trigger builds, create projects/incidents/deployments |
| viewer | Read-only dashboards |
