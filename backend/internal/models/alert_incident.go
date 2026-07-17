package models

import (
	"time"

	"github.com/google/uuid"
)

// Alert represents a monitored alert from any source.
type Alert struct {
	Base
	Title       string      `gorm:"size:512;not null" json:"title"`
	Description string      `gorm:"type:text" json:"description,omitempty"`
	Severity    Severity    `gorm:"size:50;index;not null" json:"severity"`
	Status      AlertStatus `gorm:"size:50;index;not null;default:open" json:"status"`
	Source      AlertSource `gorm:"size:50;index;not null" json:"source"`
	SourceRef   string      `gorm:"size:255;index" json:"source_ref,omitempty"`
	ProjectID   *uuid.UUID  `gorm:"type:uuid;index" json:"project_id,omitempty"`
	AckedByID   *uuid.UUID  `gorm:"type:uuid" json:"acked_by_id,omitempty"`
	AckedAt     *time.Time  `json:"acked_at,omitempty"`
	ResolvedAt  *time.Time  `json:"resolved_at,omitempty"`
	MutedUntil  *time.Time  `json:"muted_until,omitempty"`
	Labels      string      `gorm:"type:text" json:"labels,omitempty"`
	Project     *Project    `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
}

func (Alert) TableName() string { return "alerts" }

// Incident represents an operational incident.
type Incident struct {
	Base
	Title          string         `gorm:"size:512;not null" json:"title"`
	Description    string         `gorm:"type:text" json:"description,omitempty"`
	Priority       Severity       `gorm:"size:50;index;not null" json:"priority"`
	Status         IncidentStatus `gorm:"size:50;index;not null;default:open" json:"status"`
	AssigneeID     *uuid.UUID     `gorm:"type:uuid;index" json:"assignee_id,omitempty"`
	ReporterID     uuid.UUID      `gorm:"type:uuid;index;not null" json:"reporter_id"`
	ProjectID      *uuid.UUID     `gorm:"type:uuid;index" json:"project_id,omitempty"`
	RootCause      string         `gorm:"type:text" json:"root_cause,omitempty"`
	Resolution     string         `gorm:"type:text" json:"resolution,omitempty"`
	Timeline       string         `gorm:"type:text" json:"timeline,omitempty"`
	SLADeadline    *time.Time     `gorm:"index" json:"sla_deadline,omitempty"`
	ResolvedAt     *time.Time     `json:"resolved_at,omitempty"`
	Assignee       *User          `gorm:"foreignKey:AssigneeID" json:"assignee,omitempty"`
	Reporter       User           `gorm:"foreignKey:ReporterID" json:"reporter,omitempty"`
	Project        *Project       `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	Comments       []IncidentComment    `gorm:"foreignKey:IncidentID" json:"comments,omitempty"`
	Attachments    []IncidentAttachment `gorm:"foreignKey:IncidentID" json:"attachments,omitempty"`
}

func (Incident) TableName() string { return "incidents" }

// IncidentComment is a comment on an incident.
type IncidentComment struct {
	Base
	IncidentID uuid.UUID `gorm:"type:uuid;index;not null" json:"incident_id"`
	AuthorID   uuid.UUID `gorm:"type:uuid;index;not null" json:"author_id"`
	Body       string    `gorm:"type:text;not null" json:"body"`
	Author     User      `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
}

func (IncidentComment) TableName() string { return "incident_comments" }

// IncidentAttachment stores file references for incidents.
type IncidentAttachment struct {
	Base
	IncidentID uuid.UUID `gorm:"type:uuid;index;not null" json:"incident_id"`
	Filename   string    `gorm:"size:255;not null" json:"filename"`
	URL        string    `gorm:"size:1024;not null" json:"url"`
	MimeType   string    `gorm:"size:128" json:"mime_type,omitempty"`
	SizeBytes  int64     `json:"size_bytes"`
	UploadedBy uuid.UUID `gorm:"type:uuid" json:"uploaded_by"`
}

func (IncidentAttachment) TableName() string { return "incident_attachments" }
