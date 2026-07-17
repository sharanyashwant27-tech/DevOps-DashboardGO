#!/usr/bin/env sh
# Seed is applied automatically by the backend on startup when seed.enabled=true.
# This script is a convenience wrapper to restart the API and re-check health.

set -e
echo "Restarting backend to apply seed (idempotent)..."
docker compose -f deployments/docker-compose.yml restart backend
sleep 3
curl -sf http://localhost:8080/health | sed 's/.*/Health: &/'
echo "Admin: admin@devops.local / Admin@12345"
