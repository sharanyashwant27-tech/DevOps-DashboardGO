package database

import (
	"context"
	"time"

	"github.com/devops-command-center/backend/config"
	"github.com/devops-command-center/backend/internal/auth"
	"github.com/devops-command-center/backend/internal/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Seed inserts bootstrap admin, organization, sample project data.
func Seed(db *gorm.DB, cfg config.SeedConfig, log *zap.Logger) error {
	if !cfg.Enabled {
		return nil
	}
	ctx := context.Background()

	var org models.Organization
	err := db.WithContext(ctx).Where("slug = ?", "acme-corp").First(&org).Error
	if err == gorm.ErrRecordNotFound {
		org = models.Organization{
			Name:        "Acme Corp",
			Slug:        "acme-corp",
			Description: "Default organization for DevOps Command Center",
		}
		if err := db.Create(&org).Error; err != nil {
			return err
		}
	}

	var admin models.User
	err = db.WithContext(ctx).Where("email = ?", cfg.AdminEmail).First(&admin).Error
	if err == gorm.ErrRecordNotFound {
		hash, err := auth.HashPassword(cfg.AdminPassword)
		if err != nil {
			return err
		}
		admin = models.User{
			Email:          cfg.AdminEmail,
			PasswordHash:   hash,
			Name:           cfg.AdminName,
			Role:           models.RoleAdmin,
			IsActive:       true,
			OrganizationID: &org.ID,
		}
		if err := db.Create(&admin).Error; err != nil {
			return err
		}
		log.Info("seeded admin user", zap.String("email", cfg.AdminEmail))
	}

	var projectCount int64
	db.Model(&models.Project{}).Count(&projectCount)
	if projectCount == 0 {
		project := models.Project{
			OrganizationID: org.ID,
			Name:           "Platform API",
			Slug:           "platform-api",
			Description:    "Core platform services",
			RepositoryURL:  "https://github.com/example/platform-api",
			Environment:    "production",
			Status:         "active",
			OwnerID:        admin.ID,
		}
		if err := db.Create(&project).Error; err != nil {
			return err
		}

		pipeline := models.Pipeline{
			ProjectID:  project.ID,
			Name:       "main-ci",
			Provider:   "jenkins",
			ExternalID: "platform-api",
			IsActive:   true,
		}
		_ = db.Create(&pipeline).Error

		now := time.Now()
		builds := []models.Build{
			{PipelineID: pipeline.ID, ProjectID: project.ID, BuildNumber: 101, Status: models.BuildStatusSuccess, Branch: "main", DurationMs: 120000, StartedAt: &now},
			{PipelineID: pipeline.ID, ProjectID: project.ID, BuildNumber: 102, Status: models.BuildStatusFailed, Branch: "main", DurationMs: 80000, StartedAt: &now},
			{PipelineID: pipeline.ID, ProjectID: project.ID, BuildNumber: 103, Status: models.BuildStatusRunning, Branch: "feature/x", DurationMs: 0, StartedAt: &now},
		}
		for i := range builds {
			_ = db.Create(&builds[i]).Error
		}

		dep := models.Deployment{
			ProjectID:       project.ID,
			Application:     "platform-api",
			Environment:     "production",
			Version:         "v1.2.0",
			GitCommit:       "abc1234",
			Branch:          "main",
			TriggeredByID:   &admin.ID,
			TriggeredByName: admin.Name,
			Status:          models.DeploymentStatusSuccess,
			RollbackVersion: "v1.1.0",
			Logs:            "Deployment completed successfully",
			DeployedAt:      time.Now(),
		}
		_ = db.Create(&dep).Error

		alerts := []models.Alert{
			{Title: "Pod CrashLoopBackOff", Description: "api-pod restarting", Severity: models.SeverityCritical, Status: models.AlertStatusOpen, Source: models.AlertSourceKubernetes},
			{Title: "High disk usage", Description: "disk > 85%", Severity: models.SeverityMedium, Status: models.AlertStatusOpen, Source: models.AlertSourceServer},
		}
		for i := range alerts {
			_ = db.Create(&alerts[i]).Error
		}

		incident := models.Incident{
			Title:       "API latency spike",
			Description: "p99 latency above SLO",
			Priority:    models.SeverityHigh,
			Status:      models.IncidentStatusInvestigating,
			AssigneeID:  &admin.ID,
			ReporterID:  admin.ID,
			ProjectID:   &project.ID,
			Timeline:    "[" + time.Now().Format(time.RFC3339) + "] Incident opened\n",
		}
		deadline := time.Now().Add(4 * time.Hour)
		incident.SLADeadline = &deadline
		_ = db.Create(&incident).Error

		log.Info("seeded sample project data")
	}

	return nil
}
