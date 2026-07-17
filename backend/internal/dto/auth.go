package dto

import (
	"time"

	"github.com/devops-command-center/backend/internal/models"
	"github.com/google/uuid"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=72"`
	Name     string `json:"name" binding:"required,min=2,max=255"`
	Role     string `json:"role" binding:"omitempty,oneof=admin devops developer viewer"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type AuthResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresAt    time.Time    `json:"expires_at"`
	TokenType    string       `json:"token_type"`
}

type UserResponse struct {
	ID             uuid.UUID   `json:"id"`
	Email          string      `json:"email"`
	Name           string      `json:"name"`
	Role           models.Role `json:"role"`
	AvatarURL      string      `json:"avatar_url,omitempty"`
	IsActive       bool        `json:"is_active"`
	OrganizationID *uuid.UUID  `json:"organization_id,omitempty"`
	LastLoginAt    *time.Time  `json:"last_login_at,omitempty"`
	CreatedAt      time.Time   `json:"created_at"`
}

func ToUserResponse(u *models.User) UserResponse {
	return UserResponse{
		ID:             u.ID,
		Email:          u.Email,
		Name:           u.Name,
		Role:           u.Role,
		AvatarURL:      u.AvatarURL,
		IsActive:       u.IsActive,
		OrganizationID: u.OrganizationID,
		LastLoginAt:    u.LastLoginAt,
		CreatedAt:      u.CreatedAt,
	}
}
