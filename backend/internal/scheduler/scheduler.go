package scheduler

import (
	"context"
	"time"

	"github.com/devops-command-center/backend/config"
	"github.com/devops-command-center/backend/internal/models"
	"github.com/devops-command-center/backend/internal/services"
	"github.com/devops-command-center/backend/internal/websocket"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// Scheduler runs background collection and alert evaluation jobs.
type Scheduler struct {
	cron *cron.Cron
	svc  *services.Services
	hub  *websocket.Hub
	cfg  config.SchedulerConfig
	log  *zap.Logger
}

func New(cfg config.SchedulerConfig, svc *services.Services, hub *websocket.Hub, log *zap.Logger) *Scheduler {
	return &Scheduler{
		cron: cron.New(cron.WithSeconds()),
		svc:  svc,
		hub:  hub,
		cfg:  cfg,
		log:  log,
	}
}

func (s *Scheduler) Start() error {
	metricsSpec := normalize(s.cfg.MetricsInterval, "@every 30s")
	alertSpec := normalize(s.cfg.AlertInterval, "@every 1m")
	cleanupSpec := normalize(s.cfg.CleanupInterval, "@daily")

	if _, err := s.cron.AddFunc(toCron(metricsSpec), s.collectMetrics); err != nil {
		return err
	}
	if _, err := s.cron.AddFunc(toCron(alertSpec), s.evaluateAlerts); err != nil {
		return err
	}
	if _, err := s.cron.AddFunc(toCron(cleanupSpec), s.cleanup); err != nil {
		return err
	}
	s.cron.Start()
	s.log.Info("scheduler started")
	return nil
}

func (s *Scheduler) Stop() {
	ctx := s.cron.Stop()
	<-ctx.Done()
}

func (s *Scheduler) collectMetrics() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	server, _, err := s.svc.Server.CollectLocal(ctx)
	if err != nil {
		s.log.Warn("collect local metrics failed", zap.Error(err))
		return
	}
	_ = s.svc.Metrics.Record(ctx, "cpu_usage", server.CPUPercent, "percent", "server")
	_ = s.svc.Metrics.Record(ctx, "memory_usage", server.MemoryPercent, "percent", "server")
	_ = s.svc.Metrics.Record(ctx, "disk_usage", server.DiskPercent, "percent", "server")

	stats, err := s.svc.Dashboard.Stats(ctx)
	if err == nil && s.hub != nil {
		s.hub.Publish("update", "dashboard", stats)
		s.hub.Publish("update", "metrics", map[string]float64{
			"cpu": server.CPUPercent, "memory": server.MemoryPercent, "disk": server.DiskPercent,
		})
	}
}

func (s *Scheduler) evaluateAlerts() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	servers, err := s.svc.Server.List(ctx)
	if err != nil {
		return
	}
	for _, srv := range servers {
		if srv.CPUPercent > 90 {
			alert := &models.Alert{
				Title:       "High CPU usage on " + srv.Hostname,
				Description: "CPU exceeded 90%",
				Severity:    models.SeverityHigh,
				Status:      models.AlertStatusOpen,
				Source:      models.AlertSourceServer,
				SourceRef:   srv.ID.String(),
			}
			_ = s.svc.Alert.Create(ctx, alert)
			if s.hub != nil {
				s.hub.Publish("alert", "alerts", alert)
			}
		}
		if srv.MemoryPercent > 90 {
			alert := &models.Alert{
				Title:       "High memory usage on " + srv.Hostname,
				Description: "Memory exceeded 90%",
				Severity:    models.SeverityHigh,
				Status:      models.AlertStatusOpen,
				Source:      models.AlertSourceServer,
				SourceRef:   srv.ID.String(),
			}
			_ = s.svc.Alert.Create(ctx, alert)
			if s.hub != nil {
				s.hub.Publish("alert", "alerts", alert)
			}
		}
	}
}

func (s *Scheduler) cleanup() {
	s.log.Info("cleanup job executed")
}

func normalize(v, def string) string {
	if v == "" {
		return def
	}
	return v
}

// toCron converts @every / @daily style to seconds-based cron when needed.
func toCron(spec string) string {
	// robfig/cron with seconds: use standard descriptors via cron.New parser isn't used;
	// AddFunc supports @every when using default parser without seconds.
	// We created WithSeconds(), so remap descriptors.
	switch spec {
	case "@daily":
		return "0 0 0 * * *"
	case "@hourly":
		return "0 0 * * * *"
	case "@every 30s":
		return "*/30 * * * * *"
	case "@every 1m":
		return "0 */1 * * * *"
	default:
		if len(spec) > 7 && spec[:7] == "@every " {
			// fallback every minute
			return "0 */1 * * * *"
		}
		return spec
	}
}
