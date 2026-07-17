package services

import (
	"context"
	"time"

	"github.com/devops-command-center/backend/internal/dto"
	"github.com/devops-command-center/backend/internal/models"
	"github.com/devops-command-center/backend/internal/repositories"
	"go.uber.org/zap"
)

type DashboardService struct {
	repos *repositories.Repositories
	log   *zap.Logger
}

func NewDashboardService(repos *repositories.Repositories, log *zap.Logger) *DashboardService {
	return &DashboardService{repos: repos, log: log}
}

func (s *DashboardService) Stats(ctx context.Context) (*dto.DashboardStats, error) {
	stats := &dto.DashboardStats{}

	var err error
	if stats.TotalProjects, err = s.repos.Projects.Count(ctx); err != nil {
		return nil, err
	}
	if stats.RunningBuilds, err = s.repos.Builds.CountByStatus(ctx, models.BuildStatusRunning); err != nil {
		return nil, err
	}
	if stats.FailedBuilds, err = s.repos.Builds.CountByStatus(ctx, models.BuildStatusFailed); err != nil {
		return nil, err
	}
	if stats.SuccessfulBuilds, err = s.repos.Builds.CountByStatus(ctx, models.BuildStatusSuccess); err != nil {
		return nil, err
	}
	if stats.ServersOnline, err = s.repos.Servers.CountOnline(ctx); err != nil {
		return nil, err
	}
	if stats.DockerContainersRunning, err = s.repos.Containers.CountRunning(ctx); err != nil {
		return nil, err
	}
	if stats.PodsRunning, err = s.repos.Clusters.CountPodsRunning(ctx); err != nil {
		return nil, err
	}
	if stats.CriticalAlerts, err = s.repos.Alerts.CountCriticalOpen(ctx); err != nil {
		return nil, err
	}
	startOfDay := time.Now().Truncate(24 * time.Hour)
	if stats.DeploymentsToday, err = s.repos.Deployments.CountSince(ctx, startOfDay); err != nil {
		return nil, err
	}
	if stats.OpenIncidents, err = s.repos.Incidents.CountOpen(ctx); err != nil {
		return nil, err
	}

	cpu, mem, disk, netIn, netOut, err := s.repos.Servers.AvgMetrics(ctx)
	if err != nil {
		s.log.Warn("avg metrics failed", zap.Error(err))
	} else {
		stats.CPUUsage = cpu
		stats.MemoryUsage = mem
		stats.DiskUsage = disk
		stats.NetworkTrafficIn = netIn
		stats.NetworkTrafficOut = netOut
	}

	return stats, nil
}
