package database

import (
	"fmt"
	"time"

	"github.com/devops-command-center/backend/config"
	"github.com/devops-command-center/backend/internal/models"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// Connect opens a PostgreSQL connection with GORM.
func Connect(cfg config.DatabaseConfig, log *zap.Logger) (*gorm.DB, error) {
	level := gormlogger.Warn
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: gormlogger.Default.LogMode(level),
	})
	if err != nil {
		return nil, fmt.Errorf("connect database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql db: %w", err)
	}
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}
	log.Info("database connected", zap.String("host", cfg.Host), zap.String("name", cfg.Name))
	return db, nil
}

// AutoMigrate runs GORM migrations for all models.
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Organization{},
		&models.User{},
		&models.RefreshToken{},
		&models.Project{},
		&models.Pipeline{},
		&models.Build{},
		&models.Deployment{},
		&models.DockerHost{},
		&models.Container{},
		&models.Cluster{},
		&models.Pod{},
		&models.Node{},
		&models.Server{},
		&models.Alert{},
		&models.Incident{},
		&models.IncidentComment{},
		&models.IncidentAttachment{},
		&models.Metric{},
		&models.LogEntry{},
		&models.AuditLog{},
		&models.Notification{},
	)
}
