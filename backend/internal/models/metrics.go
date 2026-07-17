package models

import (
	"time"

	"github.com/google/uuid"
)

// Metric stores time-series style metric samples.
type Metric struct {
	Base
	Name       string     `gorm:"size:255;index;not null" json:"name"`
	Value      float64    `json:"value"`
	Unit       string     `gorm:"size:50" json:"unit,omitempty"`
	Source     string     `gorm:"size:100;index" json:"source,omitempty"`
	SourceID   *uuid.UUID `gorm:"type:uuid;index" json:"source_id,omitempty"`
	Labels     string     `gorm:"type:text" json:"labels,omitempty"`
	RecordedAt time.Time  `gorm:"index;not null" json:"recorded_at"`
}

func (Metric) TableName() string { return "metrics" }

// LogEntry stores application/operational logs.
type LogEntry struct {
	Base
	Level     string     `gorm:"size:20;index;not null" json:"level"`
	Message   string     `gorm:"type:text;not null" json:"message"`
	Source    string     `gorm:"size:100;index" json:"source,omitempty"`
	SourceID  *uuid.UUID `gorm:"type:uuid;index" json:"source_id,omitempty"`
	Metadata  string     `gorm:"type:text" json:"metadata,omitempty"`
	LoggedAt  time.Time  `gorm:"index;not null" json:"logged_at"`
}

func (LogEntry) TableName() string { return "logs" }

// AuditLog records security-sensitive actions.
type AuditLog struct {
	Base
	UserID     *uuid.UUID `gorm:"type:uuid;index" json:"user_id,omitempty"`
	Action     string     `gorm:"size:100;index;not null" json:"action"`
	Resource   string     `gorm:"size:100;index" json:"resource,omitempty"`
	ResourceID string     `gorm:"size:100;index" json:"resource_id,omitempty"`
	IPAddress  string     `gorm:"size:64" json:"ip_address,omitempty"`
	UserAgent  string     `gorm:"size:512" json:"user_agent,omitempty"`
	Details    string     `gorm:"type:text" json:"details,omitempty"`
	Status     string     `gorm:"size:50" json:"status,omitempty"`
}

func (AuditLog) TableName() string { return "audit_logs" }

// Notification is an in-app notification.
type Notification struct {
	Base
	UserID  uuid.UUID `gorm:"type:uuid;index;not null" json:"user_id"`
	Title   string    `gorm:"size:255;not null" json:"title"`
	Body    string    `gorm:"type:text" json:"body,omitempty"`
	Type    string    `gorm:"size:50;index" json:"type"`
	IsRead  bool      `gorm:"default:false;index" json:"is_read"`
	Link    string    `gorm:"size:512" json:"link,omitempty"`
	User    User      `gorm:"foreignKey:UserID" json:"-"`
}

func (Notification) TableName() string { return "notifications" }
