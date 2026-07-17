package models

import (
	"time"

	"github.com/google/uuid"
)

// Project represents a DevOps-managed application/project.
type Project struct {
	Base
	OrganizationID uuid.UUID    `gorm:"type:uuid;index;not null" json:"organization_id"`
	Name           string       `gorm:"size:255;not null;index" json:"name"`
	Slug           string       `gorm:"size:255;not null;index" json:"slug"`
	Description    string       `gorm:"type:text" json:"description,omitempty"`
	RepositoryURL  string       `gorm:"size:512" json:"repository_url,omitempty"`
	Environment    string       `gorm:"size:50;default:production;index" json:"environment"`
	Status         string       `gorm:"size:50;default:active;index" json:"status"`
	OwnerID        uuid.UUID    `gorm:"type:uuid;index" json:"owner_id"`
	Organization   Organization `gorm:"foreignKey:OrganizationID" json:"organization,omitempty"`
	Owner          User         `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
}

func (Project) TableName() string { return "projects" }

// Pipeline tracks CI/CD pipeline definitions.
type Pipeline struct {
	Base
	ProjectID  uuid.UUID `gorm:"type:uuid;index;not null" json:"project_id"`
	Name       string    `gorm:"size:255;not null" json:"name"`
	Provider   string    `gorm:"size:50;not null;index" json:"provider"` // jenkins | github_actions
	ExternalID string    `gorm:"size:255;index" json:"external_id,omitempty"`
	URL        string    `gorm:"size:512" json:"url,omitempty"`
	IsActive   bool      `gorm:"default:true" json:"is_active"`
	Project    Project   `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
}

func (Pipeline) TableName() string { return "pipelines" }

// Build records individual CI builds.
type Build struct {
	Base
	PipelineID    uuid.UUID   `gorm:"type:uuid;index;not null" json:"pipeline_id"`
	ProjectID     uuid.UUID   `gorm:"type:uuid;index;not null" json:"project_id"`
	BuildNumber   int         `gorm:"index" json:"build_number"`
	Status        BuildStatus `gorm:"size:50;index;not null" json:"status"`
	Branch        string      `gorm:"size:255" json:"branch,omitempty"`
	CommitSHA     string      `gorm:"size:64;index" json:"commit_sha,omitempty"`
	TriggeredBy   string      `gorm:"size:255" json:"triggered_by,omitempty"`
	DurationMs    int64       `json:"duration_ms"`
	ConsoleLogURL string      `gorm:"size:512" json:"console_log_url,omitempty"`
	StartedAt     *time.Time  `json:"started_at,omitempty"`
	FinishedAt    *time.Time  `json:"finished_at,omitempty"`
	Pipeline      Pipeline    `gorm:"foreignKey:PipelineID" json:"pipeline,omitempty"`
}

func (Build) TableName() string { return "builds" }
