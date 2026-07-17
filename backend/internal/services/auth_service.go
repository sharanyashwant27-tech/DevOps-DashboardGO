package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/devops-command-center/backend/config"
	"github.com/devops-command-center/backend/internal/auth"
	"github.com/devops-command-center/backend/internal/dto"
	"github.com/devops-command-center/backend/internal/models"
	"github.com/devops-command-center/backend/internal/repositories"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailTaken         = errors.New("email already registered")
	ErrUserInactive       = errors.New("user account is inactive")
	ErrInvalidRefresh     = errors.New("invalid refresh token")
)

type AuthService struct {
	users  repositories.UserRepository
	jwt    *auth.JWTManager
	cfg    config.JWTConfig
	log    *zap.Logger
	audit  repositories.AuditRepository
}

func NewAuthService(users repositories.UserRepository, jwt *auth.JWTManager, cfg config.JWTConfig, audit repositories.AuditRepository, log *zap.Logger) *AuthService {
	return &AuthService{users: users, jwt: jwt, cfg: cfg, audit: audit, log: log}
}

func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
	if _, err := s.users.FindByEmail(ctx, req.Email); err == nil {
		return nil, ErrEmailTaken
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	role := models.RoleViewer
	if req.Role != "" && models.Role(req.Role).IsValid() {
		role = models.Role(req.Role)
	}

	user := &models.User{
		Email:        req.Email,
		PasswordHash: hash,
		Name:         req.Name,
		Role:         role,
		IsActive:     true,
	}
	if err := s.users.Create(ctx, user); err != nil {
		return nil, err
	}

	return s.issueAuth(ctx, user, "", "")
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest, ip, ua string) (*dto.AuthResponse, error) {
	user, err := s.users.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}
	if !user.IsActive {
		return nil, ErrUserInactive
	}
	if !auth.CheckPassword(user.PasswordHash, req.Password) {
		return nil, ErrInvalidCredentials
	}

	now := time.Now()
	user.LastLoginAt = &now
	_ = s.users.Update(ctx, user)

	_ = s.audit.Create(ctx, &models.AuditLog{
		UserID: &user.ID, Action: "login", Resource: "auth",
		IPAddress: ip, UserAgent: ua, Status: "success",
	})

	return s.issueAuth(ctx, user, ip, ua)
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken, ip, ua string) (*dto.AuthResponse, error) {
	claims, err := s.jwt.ParseRefresh(refreshToken)
	if err != nil {
		return nil, ErrInvalidRefresh
	}
	hash := auth.HashToken(refreshToken)
	stored, err := s.users.FindRefreshToken(ctx, hash)
	if err != nil || stored.Revoked || stored.ExpiresAt.Before(time.Now()) {
		return nil, ErrInvalidRefresh
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, ErrInvalidRefresh
	}
	user, err := s.users.FindByID(ctx, userID)
	if err != nil || !user.IsActive {
		return nil, ErrInvalidRefresh
	}

	_ = s.users.RevokeRefreshToken(ctx, hash)
	return s.issueAuth(ctx, user, ip, ua)
}

func (s *AuthService) ForgotPassword(ctx context.Context, email string) error {
	// Production: generate reset token, store hash, send email via SMTP.
	user, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		// Do not reveal whether email exists.
		return nil
	}
	s.log.Info("password reset requested", zap.String("user_id", user.ID.String()))
	_ = s.audit.Create(ctx, &models.AuditLog{
		UserID: &user.ID, Action: "forgot_password", Resource: "auth", Status: "requested",
	})
	return nil
}

func (s *AuthService) Me(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	resp := dto.ToUserResponse(user)
	return &resp, nil
}

func (s *AuthService) issueAuth(ctx context.Context, user *models.User, ip, ua string) (*dto.AuthResponse, error) {
	pair, err := s.jwt.GeneratePair(user)
	if err != nil {
		return nil, err
	}
	rt := &models.RefreshToken{
		UserID:    user.ID,
		TokenHash: auth.HashToken(pair.RefreshToken),
		ExpiresAt: time.Now().Add(s.cfg.RefreshTTL()),
		UserAgent: ua,
		IPAddress: ip,
	}
	if err := s.users.SaveRefreshToken(ctx, rt); err != nil {
		return nil, fmt.Errorf("save refresh token: %w", err)
	}
	return &dto.AuthResponse{
		User:         dto.ToUserResponse(user),
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
		ExpiresAt:    pair.ExpiresAt,
		TokenType:    pair.TokenType,
	}, nil
}
