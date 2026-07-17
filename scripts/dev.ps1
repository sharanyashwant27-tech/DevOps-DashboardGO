# Local development helper for Windows PowerShell
$ErrorActionPreference = "Stop"

Write-Host "Starting postgres + redis..."
docker compose -f deployments/docker-compose.yml up -d postgres redis

Write-Host "Backend: cd backend; go run ./cmd/server"
Write-Host "Frontend: cd frontend; npm install; npm run dev"
Write-Host "Default admin: admin@devops.local / Admin@12345"
