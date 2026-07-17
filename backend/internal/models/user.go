package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents an authenticated platform user.
type User struct {
	Base
	Email        string     `gorm:"uniqueIndex;size:255;not null" json:"email"`
	PasswordHash string     `gorm:"size:255;not null" json:"-"`
	Name         string     `gorm:"size:255;not null" json:"name"`
	Role         Role       `gorm:"size:50;not null;default:viewer;index" json:"role"`
	AvatarURL    string     `gorm:"size:512" json:"avatar_url,omitempty"`
	IsActive     bool       `gorm:"default:true;index" json:"is_active"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	OrganizationID *uuid.UUID `gorm:"type:uuid;index" json:"organization_id,omitempty"`
	Organization   *Organization `gorm:"foreignKey:OrganizationID" json:"organization,omitempty"`
}

func (User) TableName() string { return "users" }

// RefreshToken stores refresh tokens for JWT rotation.
type RefreshToken struct {
	Base
	UserID    uuid.UUID `gorm:"type:uuid;index;not null" json:"user_id"`
	TokenHash string    `gorm:"size:255;uniqueIndex;not null" json:"-"`
	ExpiresAt time.Time `gorm:"index;not null" json:"expires_at"`
	Revoked   bool      `gorm:"default:false" json:"revoked"`
	UserAgent string    `gorm:"size:512" json:"user_agent,omitempty"`
	IPAddress string    `gorm:"size:64" json:"ip_address,omitempty"`
	User      User      `gorm:"foreignKey:UserID" json:"-"`
}

func (RefreshToken) TableName() string { return "refresh_tokens" }

// Organization groups projects and users.
type Organization struct {
	Base
	Name        string `gorm:"uniqueIndex;size:255;not null" json:"name"`
	Slug        string `gorm:"uniqueIndex;size:255;not null" json:"slug"`
	Description string `gorm:"type:text" json:"description,omitempty"`
}

func (Organization) TableName() string { return "organizations" }
