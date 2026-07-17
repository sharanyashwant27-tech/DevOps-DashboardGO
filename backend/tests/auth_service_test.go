package tests

import (
	"context"
	"testing"
	"time"

	"github.com/devops-command-center/backend/config"
	"github.com/devops-command-center/backend/internal/auth"
	"github.com/devops-command-center/backend/internal/dto"
	"github.com/devops-command-center/backend/internal/models"
	"github.com/devops-command-center/backend/internal/services"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type mockUserRepo struct {
	users  map[string]*models.User
	tokens map[string]*models.RefreshToken
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: map[string]*models.User{}, tokens: map[string]*models.RefreshToken{}}
}

func (m *mockUserRepo) Create(ctx context.Context, user *models.User) error {
	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	m.users[user.Email] = user
	return nil
}
func (m *mockUserRepo) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	for _, u := range m.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}
func (m *mockUserRepo) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	u, ok := m.users[email]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return u, nil
}
func (m *mockUserRepo) Update(ctx context.Context, user *models.User) error {
	m.users[user.Email] = user
	return nil
}
func (m *mockUserRepo) List(ctx context.Context, offset, limit int, search string) ([]models.User, int64, error) {
	return nil, 0, nil
}
func (m *mockUserRepo) SaveRefreshToken(ctx context.Context, token *models.RefreshToken) error {
	token.ID = uuid.New()
	m.tokens[token.TokenHash] = token
	return nil
}
func (m *mockUserRepo) FindRefreshToken(ctx context.Context, hash string) (*models.RefreshToken, error) {
	t, ok := m.tokens[hash]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return t, nil
}
func (m *mockUserRepo) RevokeRefreshToken(ctx context.Context, hash string) error {
	if t, ok := m.tokens[hash]; ok {
		t.Revoked = true
	}
	return nil
}

type mockAuditRepo struct{}

func (m *mockAuditRepo) Create(ctx context.Context, a *models.AuditLog) error { return nil }
func (m *mockAuditRepo) List(ctx context.Context, offset, limit int) ([]models.AuditLog, int64, error) {
	return nil, 0, nil
}

func TestAuthRegisterAndLogin(t *testing.T) {
	jwtCfg := config.JWTConfig{
		AccessSecret:     "test-access-secret-32-characters!!",
		RefreshSecret:    "test-refresh-secret-32-characters!",
		AccessTTLMinutes: 15,
		RefreshTTLHours:  24,
	}
	svc := services.NewAuthService(newMockUserRepo(), auth.NewJWTManager(jwtCfg), jwtCfg, &mockAuditRepo{}, zap.NewNop())

	reg, err := svc.Register(context.Background(), dto.RegisterRequest{
		Email: "dev@example.com", Password: "Password1!", Name: "Dev User", Role: "developer",
	})
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	if reg.AccessToken == "" || reg.RefreshToken == "" {
		t.Fatal("expected tokens")
	}

	login, err := svc.Login(context.Background(), dto.LoginRequest{
		Email: "dev@example.com", Password: "Password1!",
	}, "127.0.0.1", "test")
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	if login.User.Email != "dev@example.com" {
		t.Fatalf("unexpected email: %s", login.User.Email)
	}
}

func TestHashPassword(t *testing.T) {
	hash, err := auth.HashPassword("secret123")
	if err != nil {
		t.Fatal(err)
	}
	if !auth.CheckPassword(hash, "secret123") {
		t.Fatal("password should match")
	}
	if auth.CheckPassword(hash, "wrong") {
		t.Fatal("password should not match")
	}
}
