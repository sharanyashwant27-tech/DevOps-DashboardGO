package services

import (
	"context"
	"fmt"
	"time"

	"github.com/devops-command-center/backend/config"
	"github.com/devops-command-center/backend/internal/auth"
	"github.com/devops-command-center/backend/internal/dto"
	"github.com/devops-command-center/backend/internal/models"
	"github.com/devops-command-center/backend/internal/repositories"
	"github.com/google/uuid"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
	"go.uber.org/zap"
)

// Services aggregates domain services for DI.
type Services struct {
	Auth       *AuthService
	Dashboard  *DashboardService
	Jenkins    *JenkinsService
	GitHub     *GitHubService
	Docker     *DockerService
	Kubernetes *KubernetesService
	Project    *ProjectService
	Deployment *DeploymentService
	Incident   *IncidentService
	Alert      *AlertService
	Server     *ServerService
	Metrics    *MetricsService
	Notify     *NotificationService
	Audit      *AuditService
}

func New(
	cfg *config.Config,
	repos *repositories.Repositories,
	jwt *auth.JWTManager,
	log *zap.Logger,
) *Services {
	return &Services{
		Auth:       NewAuthService(repos.Users, jwt, cfg.JWT, repos.Audit, log),
		Dashboard:  NewDashboardService(repos, log),
		Jenkins:    NewJenkinsService(cfg.Jenkins, log),
		GitHub:     NewGitHubService(cfg.GitHub, log),
		Docker:     NewDockerService(cfg.Docker, log),
		Kubernetes: NewKubernetesService(cfg.Kubernetes, log),
		Project:    NewProjectService(repos.Projects, repos.Audit, log),
		Deployment: NewDeploymentService(repos.Deployments, repos.Audit, log),
		Incident:   NewIncidentService(repos.Incidents, repos.Audit, log),
		Alert:      NewAlertService(repos.Alerts, repos.Audit, log),
		Server:     NewServerService(repos.Servers, log),
		Metrics:    NewMetricsService(repos.Metrics, log),
		Notify:     NewNotificationService(cfg, log),
		Audit:      NewAuditService(repos.Audit),
	}
}

// --- Project ---

type ProjectService struct {
	repo  repositories.ProjectRepository
	audit repositories.AuditRepository
	log   *zap.Logger
}

func NewProjectService(repo repositories.ProjectRepository, audit repositories.AuditRepository, log *zap.Logger) *ProjectService {
	return &ProjectService{repo: repo, audit: audit, log: log}
}

func (s *ProjectService) Create(ctx context.Context, req dto.ProjectRequest, ownerID uuid.UUID) (*models.Project, error) {
	p := &models.Project{
		OrganizationID: req.OrganizationID,
		Name:           req.Name,
		Slug:           req.Slug,
		Description:    req.Description,
		RepositoryURL:  req.RepositoryURL,
		Environment:    req.Environment,
		Status:         "active",
		OwnerID:        ownerID,
	}
	if p.Environment == "" {
		p.Environment = "production"
	}
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	_ = s.audit.Create(ctx, &models.AuditLog{
		UserID: &ownerID, Action: "create", Resource: "project", ResourceID: p.ID.String(), Status: "success",
	})
	return p, nil
}

func (s *ProjectService) Get(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *ProjectService) List(ctx context.Context, q dto.PaginationQuery) ([]models.Project, int64, error) {
	return s.repo.List(ctx, q.Offset(), q.PageSize, q.Search)
}

func (s *ProjectService) Delete(ctx context.Context, id, actor uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	_ = s.audit.Create(ctx, &models.AuditLog{
		UserID: &actor, Action: "delete", Resource: "project", ResourceID: id.String(), Status: "success",
	})
	return nil
}

// --- Deployment ---

type DeploymentService struct {
	repo  repositories.DeploymentRepository
	audit repositories.AuditRepository
	log   *zap.Logger
}

func NewDeploymentService(repo repositories.DeploymentRepository, audit repositories.AuditRepository, log *zap.Logger) *DeploymentService {
	return &DeploymentService{repo: repo, audit: audit, log: log}
}

func (s *DeploymentService) Create(ctx context.Context, req dto.DeploymentRequest, actorID uuid.UUID, actorName string) (*models.Deployment, error) {
	d := &models.Deployment{
		ProjectID:       req.ProjectID,
		Application:     req.Application,
		Environment:     req.Environment,
		Version:         req.Version,
		GitCommit:       req.GitCommit,
		Branch:          req.Branch,
		TriggeredByID:   &actorID,
		TriggeredByName: actorName,
		Status:          models.DeploymentStatusSuccess,
		Logs:            req.Logs,
		DeployedAt:      time.Now(),
	}
	if err := s.repo.Create(ctx, d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *DeploymentService) List(ctx context.Context, q dto.PaginationQuery, env string) ([]models.Deployment, int64, error) {
	return s.repo.List(ctx, q.Offset(), q.PageSize, env)
}

func (s *DeploymentService) Get(ctx context.Context, id uuid.UUID) (*models.Deployment, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *DeploymentService) Rollback(ctx context.Context, id, actorID uuid.UUID) (*models.Deployment, error) {
	current, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if current.RollbackVersion == "" {
		return nil, fmt.Errorf("no rollback version available")
	}
	rollback := &models.Deployment{
		ProjectID:       current.ProjectID,
		Application:     current.Application,
		Environment:     current.Environment,
		Version:         current.RollbackVersion,
		GitCommit:       current.GitCommit,
		Branch:          current.Branch,
		TriggeredByID:   &actorID,
		TriggeredByName: "rollback",
		Status:          models.DeploymentStatusRolledBack,
		RollbackVersion: current.Version,
		Logs:            fmt.Sprintf("Rolled back from %s to %s", current.Version, current.RollbackVersion),
		DeployedAt:      time.Now(),
	}
	if err := s.repo.Create(ctx, rollback); err != nil {
		return nil, err
	}
	current.Status = models.DeploymentStatusRolledBack
	_ = s.repo.Update(ctx, current)
	_ = s.audit.Create(ctx, &models.AuditLog{
		UserID: &actorID, Action: "rollback", Resource: "deployment", ResourceID: id.String(), Status: "success",
	})
	return rollback, nil
}

// --- Incident ---

type IncidentService struct {
	repo  repositories.IncidentRepository
	audit repositories.AuditRepository
	log   *zap.Logger
}

func NewIncidentService(repo repositories.IncidentRepository, audit repositories.AuditRepository, log *zap.Logger) *IncidentService {
	return &IncidentService{repo: repo, audit: audit, log: log}
}

func (s *IncidentService) Create(ctx context.Context, req dto.IncidentRequest, reporterID uuid.UUID) (*models.Incident, error) {
	inc := &models.Incident{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		Status:      models.IncidentStatusOpen,
		AssigneeID:  req.AssigneeID,
		ReporterID:  reporterID,
		ProjectID:   req.ProjectID,
		Timeline:    fmt.Sprintf("[%s] Incident created\n", time.Now().Format(time.RFC3339)),
	}
	if req.SLAHours > 0 {
		deadline := time.Now().Add(time.Duration(req.SLAHours) * time.Hour)
		inc.SLADeadline = &deadline
	}
	if err := s.repo.Create(ctx, inc); err != nil {
		return nil, err
	}
	return inc, nil
}

func (s *IncidentService) Get(ctx context.Context, id uuid.UUID) (*models.Incident, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *IncidentService) List(ctx context.Context, q dto.PaginationQuery, status, priority string) ([]models.Incident, int64, error) {
	return s.repo.List(ctx, q.Offset(), q.PageSize, status, priority, q.Search)
}

func (s *IncidentService) Update(ctx context.Context, id uuid.UUID, req dto.IncidentUpdateRequest) (*models.Incident, error) {
	inc, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Status != nil {
		inc.Status = *req.Status
		inc.Timeline += fmt.Sprintf("[%s] Status -> %s\n", time.Now().Format(time.RFC3339), *req.Status)
		if *req.Status == models.IncidentStatusResolved || *req.Status == models.IncidentStatusClosed {
			now := time.Now()
			inc.ResolvedAt = &now
		}
	}
	if req.Priority != nil {
		inc.Priority = *req.Priority
	}
	if req.AssigneeID != nil {
		inc.AssigneeID = req.AssigneeID
	}
	if req.RootCause != nil {
		inc.RootCause = *req.RootCause
	}
	if req.Resolution != nil {
		inc.Resolution = *req.Resolution
	}
	if req.Timeline != nil {
		inc.Timeline = *req.Timeline
	}
	if err := s.repo.Update(ctx, inc); err != nil {
		return nil, err
	}
	return inc, nil
}

func (s *IncidentService) AddComment(ctx context.Context, incidentID, authorID uuid.UUID, body string) (*models.IncidentComment, error) {
	c := &models.IncidentComment{IncidentID: incidentID, AuthorID: authorID, Body: body}
	if err := s.repo.AddComment(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

// --- Alert ---

type AlertService struct {
	repo  repositories.AlertRepository
	audit repositories.AuditRepository
	log   *zap.Logger
}

func NewAlertService(repo repositories.AlertRepository, audit repositories.AuditRepository, log *zap.Logger) *AlertService {
	return &AlertService{repo: repo, audit: audit, log: log}
}

func (s *AlertService) Create(ctx context.Context, a *models.Alert) error {
	return s.repo.Create(ctx, a)
}

func (s *AlertService) List(ctx context.Context, q dto.PaginationQuery, severity, status, source string) ([]models.Alert, int64, error) {
	return s.repo.List(ctx, q.Offset(), q.PageSize, severity, status, source, q.Search)
}

func (s *AlertService) Acknowledge(ctx context.Context, id, userID uuid.UUID) (*models.Alert, error) {
	a, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	a.Status = models.AlertStatusAcknowledged
	a.AckedByID = &userID
	a.AckedAt = &now
	return a, s.repo.Update(ctx, a)
}

func (s *AlertService) Resolve(ctx context.Context, id uuid.UUID) (*models.Alert, error) {
	a, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	a.Status = models.AlertStatusResolved
	a.ResolvedAt = &now
	return a, s.repo.Update(ctx, a)
}

func (s *AlertService) Mute(ctx context.Context, id uuid.UUID, minutes int) (*models.Alert, error) {
	a, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	until := time.Now().Add(time.Duration(minutes) * time.Minute)
	a.Status = models.AlertStatusMuted
	a.MutedUntil = &until
	return a, s.repo.Update(ctx, a)
}

// --- Server ---

type ServerService struct {
	repo repositories.ServerRepository
	log  *zap.Logger
}

func NewServerService(repo repositories.ServerRepository, log *zap.Logger) *ServerService {
	return &ServerService{repo: repo, log: log}
}

func (s *ServerService) List(ctx context.Context) ([]models.Server, error) {
	return s.repo.List(ctx)
}

func (s *ServerService) Get(ctx context.Context, id uuid.UUID) (*models.Server, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *ServerService) CollectLocal(ctx context.Context) (*models.Server, map[string]interface{}, error) {
	hostInfo, err := host.InfoWithContext(ctx)
	if err != nil {
		return nil, nil, err
	}
	cpuPercent, _ := cpu.PercentWithContext(ctx, 0, false)
	vm, _ := mem.VirtualMemoryWithContext(ctx)
	du, _ := disk.UsageWithContext(ctx, "/")
	avg, _ := load.AvgWithContext(ctx)
	uptime, _ := host.UptimeWithContext(ctx)
	ioCounters, _ := net.IOCountersWithContext(ctx, false)

	cpuVal := 0.0
	if len(cpuPercent) > 0 {
		cpuVal = cpuPercent[0]
	}
	memVal := 0.0
	diskVal := 0.0
	if vm != nil {
		memVal = vm.UsedPercent
	}
	if du != nil {
		diskVal = du.UsedPercent
	}
	var load1, load5, load15 float64
	if avg != nil {
		load1, load5, load15 = avg.Load1, avg.Load5, avg.Load15
	}
	var netIn, netOut int64
	if len(ioCounters) > 0 {
		netIn = int64(ioCounters[0].BytesRecv)
		netOut = int64(ioCounters[0].BytesSent)
	}

	osType := "linux"
	if hostInfo.OS == "windows" {
		osType = "windows"
	}

	server := &models.Server{
		Hostname:      hostInfo.Hostname,
		IPAddress:     "",
		OS:            osType,
		Status:        "online",
		CPUPercent:    cpuVal,
		MemoryPercent: memVal,
		DiskPercent:   diskVal,
		LoadAvg1:      load1,
		LoadAvg5:      load5,
		LoadAvg15:     load15,
		UptimeSeconds: int64(uptime),
		NetworkInBps:  netIn,
		NetworkOutBps: netOut,
		LastSeenAt:    time.Now(),
		AgentVersion:  "1.0.0",
	}
	if err := s.repo.Upsert(ctx, server); err != nil {
		s.log.Warn("upsert server failed", zap.Error(err))
	}

	procs, _ := process.ProcessesWithContext(ctx)
	type procInfo struct {
		PID     int32   `json:"pid"`
		Name    string  `json:"name"`
		CPU     float64 `json:"cpu"`
		Memory  float32 `json:"memory"`
	}
	topCPU := make([]procInfo, 0)
	topMem := make([]procInfo, 0)
	for i, p := range procs {
		if i > 40 {
			break
		}
		name, _ := p.NameWithContext(ctx)
		c, _ := p.CPUPercentWithContext(ctx)
		m, _ := p.MemoryPercentWithContext(ctx)
		info := procInfo{PID: p.Pid, Name: name, CPU: c, Memory: m}
		topCPU = append(topCPU, info)
		topMem = append(topMem, info)
	}

	details := map[string]interface{}{
		"server":            server,
		"top_cpu_processes": topCPU,
		"top_mem_processes": topMem,
		"filesystem":        du,
		"platform":          hostInfo.Platform,
	}
	return server, details, nil
}

// --- Metrics ---

type MetricsService struct {
	repo repositories.MetricRepository
	log  *zap.Logger
}

func NewMetricsService(repo repositories.MetricRepository, log *zap.Logger) *MetricsService {
	return &MetricsService{repo: repo, log: log}
}

func (s *MetricsService) Record(ctx context.Context, name string, value float64, unit, source string) error {
	return s.repo.Create(ctx, &models.Metric{
		Name: name, Value: value, Unit: unit, Source: source, RecordedAt: time.Now(),
	})
}

func (s *MetricsService) Series(ctx context.Context, name string, hours int) (*dto.MetricSeriesResponse, error) {
	since := time.Now().Add(-time.Duration(hours) * time.Hour)
	items, err := s.repo.Series(ctx, name, since)
	if err != nil {
		return nil, err
	}
	points := make([]dto.MetricSeriesPoint, 0, len(items))
	unit := ""
	for _, m := range items {
		unit = m.Unit
		points = append(points, dto.MetricSeriesPoint{Timestamp: m.RecordedAt, Value: m.Value})
	}
	return &dto.MetricSeriesResponse{Name: name, Unit: unit, Points: points}, nil
}

// --- Notification ---

type NotificationService struct {
	cfg *config.Config
	log *zap.Logger
}

func NewNotificationService(cfg *config.Config, log *zap.Logger) *NotificationService {
	return &NotificationService{cfg: cfg, log: log}
}

func (s *NotificationService) SendSlack(message string) error {
	if !s.cfg.Slack.Enabled() {
		return fmt.Errorf("slack not configured")
	}
	s.log.Info("slack notification (stub send)", zap.String("webhook", mask(s.cfg.Slack.WebhookURL)), zap.String("message", message))
	return nil
}

func (s *NotificationService) SendTeams(message string) error {
	if !s.cfg.Teams.Enabled() {
		return fmt.Errorf("teams not configured")
	}
	s.log.Info("teams notification (stub send)", zap.String("message", message))
	return nil
}

func (s *NotificationService) SendEmail(to, subject, body string) error {
	if !s.cfg.SMTP.Enabled() {
		return fmt.Errorf("smtp not configured")
	}
	s.log.Info("email notification (stub send)", zap.String("to", to), zap.String("subject", subject))
	return nil
}

func mask(s string) string {
	if len(s) < 8 {
		return "***"
	}
	return s[:4] + "****"
}

// --- Audit ---

type AuditService struct {
	repo repositories.AuditRepository
}

func NewAuditService(repo repositories.AuditRepository) *AuditService {
	return &AuditService{repo: repo}
}

func (s *AuditService) List(ctx context.Context, q dto.PaginationQuery) ([]models.AuditLog, int64, error) {
	return s.repo.List(ctx, q.Offset(), q.PageSize)
}

func (s *AuditService) Log(ctx context.Context, a *models.AuditLog) error {
	return s.repo.Create(ctx, a)
}
