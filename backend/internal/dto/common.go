package dto

import (
	"time"

	"github.com/devops-command-center/backend/internal/models"
	"github.com/google/uuid"
)

type PaginationQuery struct {
	Page     int    `form:"page,default=1" binding:"min=1"`
	PageSize int    `form:"page_size,default=20" binding:"min=1,max=100"`
	Search   string `form:"search"`
	SortBy   string `form:"sort_by"`
	SortDir  string `form:"sort_dir,default=desc"`
}

func (p PaginationQuery) Offset() int {
	return (p.Page - 1) * p.PageSize
}

type ProjectRequest struct {
	OrganizationID uuid.UUID `json:"organization_id" binding:"required"`
	Name           string    `json:"name" binding:"required,min=2,max=255"`
	Slug           string    `json:"slug" binding:"required,min=2,max=255"`
	Description    string    `json:"description"`
	RepositoryURL  string    `json:"repository_url"`
	Environment    string    `json:"environment" binding:"omitempty,oneof=development staging production"`
}

type DeploymentRequest struct {
	ProjectID   uuid.UUID `json:"project_id" binding:"required"`
	Application string    `json:"application" binding:"required"`
	Environment string    `json:"environment" binding:"required"`
	Version     string    `json:"version" binding:"required"`
	GitCommit   string    `json:"git_commit"`
	Branch      string    `json:"branch"`
	Logs        string    `json:"logs"`
}

type IncidentRequest struct {
	Title       string         `json:"title" binding:"required"`
	Description string         `json:"description"`
	Priority    models.Severity `json:"priority" binding:"required,oneof=critical high medium low"`
	ProjectID   *uuid.UUID     `json:"project_id"`
	AssigneeID  *uuid.UUID     `json:"assignee_id"`
	SLAHours    int            `json:"sla_hours"`
}

type IncidentUpdateRequest struct {
	Status     *models.IncidentStatus `json:"status"`
	Priority   *models.Severity       `json:"priority"`
	AssigneeID *uuid.UUID             `json:"assignee_id"`
	RootCause  *string                `json:"root_cause"`
	Resolution *string                `json:"resolution"`
	Timeline   *string                `json:"timeline"`
}

type IncidentCommentRequest struct {
	Body string `json:"body" binding:"required,min=1"`
}

type AlertMuteRequest struct {
	Minutes int `json:"minutes" binding:"required,min=1"`
}

type ScaleDeploymentRequest struct {
	Replicas int32 `json:"replicas" binding:"required,min=0"`
}

type DashboardStats struct {
	TotalProjects          int64   `json:"total_projects"`
	RunningBuilds          int64   `json:"running_builds"`
	FailedBuilds           int64   `json:"failed_builds"`
	SuccessfulBuilds       int64   `json:"successful_builds"`
	ServersOnline          int64   `json:"servers_online"`
	DockerContainersRunning int64  `json:"docker_containers_running"`
	PodsRunning            int64   `json:"pods_running"`
	CriticalAlerts         int64   `json:"critical_alerts"`
	DeploymentsToday       int64   `json:"deployments_today"`
	OpenIncidents          int64   `json:"open_incidents"`
	CPUUsage               float64 `json:"cpu_usage"`
	MemoryUsage            float64 `json:"memory_usage"`
	DiskUsage              float64 `json:"disk_usage"`
	NetworkTrafficIn       int64   `json:"network_traffic_in"`
	NetworkTrafficOut      int64   `json:"network_traffic_out"`
}

type MetricSeriesPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
}

type MetricSeriesResponse struct {
	Name   string              `json:"name"`
	Unit   string              `json:"unit"`
	Points []MetricSeriesPoint `json:"points"`
}
