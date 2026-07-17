-- DevOps Command Center schema (reference / optional init)
-- GORM AutoMigrate is the primary migration path; this file documents the ER model.

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Core identity
-- users, organizations, refresh_tokens
-- Domain
-- projects, pipelines, builds, deployments
-- Infrastructure
-- docker_hosts, containers, clusters, pods, nodes, servers
-- Operations
-- alerts, incidents, incident_comments, incident_attachments
-- Observability
-- metrics, logs, audit_logs, notifications

-- Indexes commonly used by the application are created by GORM tags.
-- Seed data is applied by backend/internal/database/seed.go when seed.enabled=true.
