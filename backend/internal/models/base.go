package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base contains common columns for all models.
type Base struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (b *Base) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

// Role represents RBAC roles.
type Role string

const (
	RoleAdmin    Role = "admin"
	RoleDevOps   Role = "devops"
	RoleDeveloper Role = "developer"
	RoleViewer   Role = "viewer"
)

func (r Role) IsValid() bool {
	switch r {
	case RoleAdmin, RoleDevOps, RoleDeveloper, RoleViewer:
		return true
	default:
		return false
	}
}

// Severity levels for alerts/incidents.
type Severity string

const (
	SeverityCritical Severity = "critical"
	SeverityHigh     Severity = "high"
	SeverityMedium   Severity = "medium"
	SeverityLow      Severity = "low"
)

// Status enums.
type BuildStatus string

const (
	BuildStatusPending  BuildStatus = "pending"
	BuildStatusRunning  BuildStatus = "running"
	BuildStatusSuccess  BuildStatus = "success"
	BuildStatusFailed   BuildStatus = "failed"
	BuildStatusAborted  BuildStatus = "aborted"
)

type DeploymentStatus string

const (
	DeploymentStatusPending   DeploymentStatus = "pending"
	DeploymentStatusRunning   DeploymentStatus = "running"
	DeploymentStatusSuccess   DeploymentStatus = "success"
	DeploymentStatusFailed    DeploymentStatus = "failed"
	DeploymentStatusRolledBack DeploymentStatus = "rolled_back"
)

type IncidentStatus string

const (
	IncidentStatusOpen       IncidentStatus = "open"
	IncidentStatusInvestigating IncidentStatus = "investigating"
	IncidentStatusMitigated  IncidentStatus = "mitigated"
	IncidentStatusResolved   IncidentStatus = "resolved"
	IncidentStatusClosed     IncidentStatus = "closed"
)

type AlertStatus string

const (
	AlertStatusOpen         AlertStatus = "open"
	AlertStatusAcknowledged AlertStatus = "acknowledged"
	AlertStatusResolved     AlertStatus = "resolved"
	AlertStatusMuted        AlertStatus = "muted"
)

type AlertSource string

const (
	AlertSourcePrometheus AlertSource = "prometheus"
	AlertSourceServer     AlertSource = "server"
	AlertSourceDocker     AlertSource = "docker"
	AlertSourceKubernetes AlertSource = "kubernetes"
	AlertSourceJenkins    AlertSource = "jenkins"
	AlertSourceGitHub     AlertSource = "github"
	AlertSourceSystem     AlertSource = "system"
)
