package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/devops-command-center/backend/config"
	_ "github.com/devops-command-center/backend/docs"
	"github.com/devops-command-center/backend/internal/auth"
	"github.com/devops-command-center/backend/internal/controllers"
	"github.com/devops-command-center/backend/internal/database"
	"github.com/devops-command-center/backend/internal/repositories"
	"github.com/devops-command-center/backend/internal/routes"
	"github.com/devops-command-center/backend/internal/scheduler"
	"github.com/devops-command-center/backend/internal/services"
	"github.com/devops-command-center/backend/internal/websocket"
	"github.com/devops-command-center/backend/pkg/logger"
	redisclient "github.com/devops-command-center/backend/pkg/redis"
	"go.uber.org/zap"
)

// @title DevOps Command Center API
// @version 1.0
// @description Enterprise DevOps Dashboard REST API
// @host localhost:8095
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfgPath := os.Getenv("DCC_CONFIG")
	if cfgPath == "" {
		cfgPath = filepath.Join("config", "config.yaml")
		if _, err := os.Stat(cfgPath); err != nil {
			cfgPath = filepath.Join("backend", "config", "config.yaml")
		}
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		panic(fmt.Errorf("load config: %w", err))
	}

	log := logger.MustNew(cfg.Server.Mode)
	defer logger.Sync(log)
	log = logger.WithPID(log)

	db, err := database.Connect(cfg.Database, log)
	if err != nil {
		log.Fatal("database connection failed", zap.Error(err))
	}
	if err := database.AutoMigrate(db); err != nil {
		log.Fatal("migration failed", zap.Error(err))
	}
	if err := database.Seed(db, cfg.Seed, log); err != nil {
		log.Fatal("seed failed", zap.Error(err))
	}

	rdb, err := redisclient.New(cfg.Redis, log)
	if err != nil {
		log.Warn("redis unavailable, continuing without cache", zap.Error(err))
	} else {
		defer rdb.Close()
	}

	jwtMgr := auth.NewJWTManager(cfg.JWT)
	repos := repositories.New(db)
	svc := services.New(cfg, repos, jwtMgr, log)
	ctrl := controllers.New(svc, log)
	hub := websocket.NewHub(log)
	go hub.Run()

	sched := scheduler.New(cfg.Scheduler, svc, hub, log)
	if err := sched.Start(); err != nil {
		log.Fatal("scheduler failed", zap.Error(err))
	}
	defer sched.Stop()

	router := routes.Setup(cfg, ctrl, jwtMgr, hub, log)
	srv := &http.Server{
		Addr:         cfg.Server.Addr(),
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	go func() {
		log.Info("server starting", zap.String("addr", cfg.Server.Addr()))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	_ = svc.Docker.Close()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("forced shutdown", zap.Error(err))
	}
	log.Info("server stopped")
}
