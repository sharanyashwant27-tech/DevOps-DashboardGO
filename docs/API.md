# API Documentation

Base URL: `http://localhost:8080/api/v1`  
Swagger UI: `http://localhost:8080/swagger/index.html`  
Auth: `Authorization: Bearer <access_token>`

## Authentication

| Method | Path | Description |
|--------|------|-------------|
| POST | `/auth/register` | Create user |
| POST | `/auth/login` | Login, returns JWT pair |
| POST | `/auth/refresh` | Rotate tokens |
| POST | `/auth/forgot-password` | Request reset email |
| GET | `/auth/me` | Current user |

## Dashboard

| Method | Path | Description |
|--------|------|-------------|
| GET | `/dashboard/stats` | Aggregated KPIs |

## Projects / Deployments / Incidents / Alerts

| Method | Path | Roles |
|--------|------|-------|
| GET/POST | `/projects` | viewer+ / developer+ |
| GET/DELETE | `/projects/:id` | viewer+ / devops+ |
| GET/POST | `/deployments` | viewer+ / developer+ |
| POST | `/deployments/:id/rollback` | devops+ |
| GET/POST/PATCH | `/incidents` | developer+ for writes |
| POST | `/incidents/:id/comments` | authenticated |
| GET | `/alerts` | authenticated |
| POST | `/alerts/:id/acknowledge\|resolve\|mute` | ops roles for resolve/mute |

## Jenkins

| Method | Path |
|--------|------|
| GET | `/jenkins/jobs?search=` |
| GET | `/jenkins/jobs/:job/builds` |
| GET | `/jenkins/queue` |
| GET | `/jenkins/stats` |
| POST | `/jenkins/jobs/:job/build` |
| POST | `/jenkins/jobs/:job/builds/:number/stop` |
| GET | `/jenkins/jobs/:job/builds/:number/console` |

## GitHub

| Method | Path |
|--------|------|
| GET | `/github/repos` |
| GET | `/github/repos/:owner/:repo/branches` |
| GET | `/github/repos/:owner/:repo/commits` |
| GET | `/github/repos/:owner/:repo/pulls` |
| GET | `/github/repos/:owner/:repo/issues` |
| GET | `/github/repos/:owner/:repo/releases` |
| GET | `/github/repos/:owner/:repo/actions/runs` |
| GET | `/github/repos/:owner/:repo/contributors` |
| GET | `/github/repos/:owner/:repo/health` |

## Docker

| Method | Path |
|--------|------|
| GET | `/docker/containers?search=` |
| GET | `/docker/images\|volumes\|networks` |
| GET | `/docker/containers/:id/stats\|logs` |
| POST | `/docker/containers/:id/start\|stop\|restart` |
| DELETE | `/docker/containers/:id` |

## Kubernetes

| Method | Path |
|--------|------|
| GET | `/kubernetes/namespaces\|pods\|deployments\|replicasets\|daemonsets\|services\|ingresses\|nodes\|events` |
| GET | `/kubernetes/persistentvolumes\|persistentvolumeclaims` |
| GET | `/kubernetes/pods/:pod/logs` |
| POST | `/kubernetes/deployments/:name/scale\|restart` |
| DELETE | `/kubernetes/pods/:pod` |

## Servers / Metrics / Audit

| Method | Path |
|--------|------|
| GET | `/servers` |
| GET | `/servers/local` |
| GET | `/metrics/:name?hours=24` |
| GET | `/audit` | admin/devops |

## Realtime

WebSocket: `ws://localhost:8080/ws`  
Channels: `dashboard`, `alerts`, `builds`, `metrics`

## Health / Metrics

| Method | Path |
|--------|------|
| GET | `/health` |
| GET | `/metrics` | Prometheus scrape |
