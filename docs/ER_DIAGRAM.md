# Database ER Diagram

```mermaid
erDiagram
    ORGANIZATIONS ||--o{ USERS : has
    ORGANIZATIONS ||--o{ PROJECTS : owns
    USERS ||--o{ PROJECTS : owns
    USERS ||--o{ REFRESH_TOKENS : has
    USERS ||--o{ AUDIT_LOGS : writes
    USERS ||--o{ NOTIFICATIONS : receives
    PROJECTS ||--o{ PIPELINES : contains
    PROJECTS ||--o{ BUILDS : tracks
    PROJECTS ||--o{ DEPLOYMENTS : ships
    PROJECTS ||--o{ ALERTS : raises
    PROJECTS ||--o{ INCIDENTS : opens
    PIPELINES ||--o{ BUILDS : produces
    USERS ||--o{ INCIDENTS : reports
    USERS ||--o{ INCIDENTS : assigned
    INCIDENTS ||--o{ INCIDENT_COMMENTS : has
    INCIDENTS ||--o{ INCIDENT_ATTACHMENTS : has
    DOCKER_HOSTS ||--o{ CONTAINERS : hosts
    CLUSTERS ||--o{ PODS : runs
    CLUSTERS ||--o{ NODES : contains

    ORGANIZATIONS {
        uuid id PK
        string name
        string slug
    }
    USERS {
        uuid id PK
        string email
        string role
        uuid organization_id FK
    }
    PROJECTS {
        uuid id PK
        uuid organization_id FK
        uuid owner_id FK
        string name
        string environment
    }
    PIPELINES {
        uuid id PK
        uuid project_id FK
        string provider
    }
    BUILDS {
        uuid id PK
        uuid pipeline_id FK
        uuid project_id FK
        string status
    }
    DEPLOYMENTS {
        uuid id PK
        uuid project_id FK
        string version
        string status
    }
    ALERTS {
        uuid id PK
        string severity
        string source
        string status
    }
    INCIDENTS {
        uuid id PK
        string priority
        string status
        timestamp sla_deadline
    }
    SERVERS {
        uuid id PK
        string hostname
        float cpu_percent
        float memory_percent
    }
    METRICS {
        uuid id PK
        string name
        float value
        timestamp recorded_at
    }
    AUDIT_LOGS {
        uuid id PK
        uuid user_id FK
        string action
        string resource
    }
```

## Table Inventory

| Table | Purpose |
|-------|---------|
| users | Auth identities + RBAC roles |
| organizations | Multi-tenant org boundary |
| refresh_tokens | JWT refresh rotation |
| projects | Managed applications |
| pipelines | CI definitions |
| builds | Build history |
| deployments | Release history + rollback |
| docker_hosts / containers | Docker inventory cache |
| clusters / pods / nodes | Kubernetes inventory cache |
| servers | Host monitoring |
| alerts | Multi-source alerts |
| incidents (+ comments/attachments) | Incident management |
| metrics / logs | Telemetry |
| audit_logs | Security audit trail |
| notifications | In-app notifications |
