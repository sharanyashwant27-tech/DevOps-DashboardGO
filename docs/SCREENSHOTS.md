# Sample Screenshots Guide

Capture these views after `docker compose up` and logging in as the seeded admin.

## Suggested captures

1. **Login** — `frontend/src/pages/LoginPage.tsx` branded auth card  
2. **Dashboard** — animated KPI grid + Chart.js usage panels  
3. **Jenkins** — job table with trigger/console actions  
4. **GitHub** — repository list + health/commits tabs  
5. **Docker** — container table with lifecycle controls  
6. **Kubernetes** — namespace selector + pods/deployments JSON inspector  
7. **Incidents / Alerts** — priority chips and action buttons  
8. **Dark / Light mode** — toggle from the top bar  

Place PNG/WebP files under `docs/screenshots/` (create locally):

```
docs/screenshots/
  01-login.png
  02-dashboard.png
  03-jenkins.png
  04-github.png
  05-docker.png
  06-kubernetes.png
  07-alerts.png
  08-theme.png
```

Embed in README when available:

```markdown
![Dashboard](docs/screenshots/02-dashboard.png)
```
