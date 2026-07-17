package models

import (
	"time"

	"github.com/google/uuid"
)

// Deployment stores every deployment event for rollback and audit.
type Deployment struct {
	Base
	ProjectID       uuid.UUID        `gorm:"type:uuid;index;not null" json:"project_id"`
	Application     string           `gorm:"size:255;not null;index" json:"application"`
	Environment     string           `gorm:"size:50;not null;index" json:"environment"`
	Version         string           `gorm:"size:100;not null" json:"version"`
	GitCommit       string           `gorm:"size:64;index" json:"git_commit,omitempty"`
	Branch          string           `gorm:"size:255" json:"branch,omitempty"`
	TriggeredByID   *uuid.UUID       `gorm:"type:uuid;index" json:"triggered_by_id,omitempty"`
	TriggeredByName string           `gorm:"size:255" json:"triggered_by,omitempty"`
	Status          DeploymentStatus `gorm:"size:50;index;not null" json:"status"`
	RollbackVersion string           `gorm:"size:100" json:"rollback_version,omitempty"`
	Logs            string           `gorm:"type:text" json:"logs,omitempty"`
	DeployedAt      time.Time        `gorm:"index;not null" json:"deployed_at"`
	FinishedAt      *time.Time       `json:"finished_at,omitempty"`
	Project         Project          `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	TriggeredBy     *User            `gorm:"foreignKey:TriggeredByID" json:"triggered_by_user,omitempty"`
}

func (Deployment) TableName() string { return "deployments" }
