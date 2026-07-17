package repositories

import (
	"context"
	"time"

	"github.com/devops-command-center/backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repositories aggregates all repository interfaces for DI.
type Repositories struct {
	Users       UserRepository
	Projects    ProjectRepository
	Builds      BuildRepository
	Deployments DeploymentRepository
	Alerts      AlertRepository
	Incidents   IncidentRepository
	Servers     ServerRepository
	Metrics     MetricRepository
	Audit       AuditRepository
	Containers  ContainerRepository
	Clusters    ClusterRepository
}

func New(db *gorm.DB) *Repositories {
	return &Repositories{
		Users:       NewUserRepository(db),
		Projects:    NewProjectRepository(db),
		Builds:      NewBuildRepository(db),
		Deployments: NewDeploymentRepository(db),
		Alerts:      NewAlertRepository(db),
		Incidents:   NewIncidentRepository(db),
		Servers:     NewServerRepository(db),
		Metrics:     NewMetricRepository(db),
		Audit:       NewAuditRepository(db),
		Containers:  NewContainerRepository(db),
		Clusters:    NewClusterRepository(db),
	}
}

// --- Project ---

type ProjectRepository interface {
	Create(ctx context.Context, p *models.Project) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Project, error)
	Update(ctx context.Context, p *models.Project) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, offset, limit int, search string) ([]models.Project, int64, error)
	Count(ctx context.Context) (int64, error)
}

type projectRepository struct{ db *gorm.DB }

func NewProjectRepository(db *gorm.DB) ProjectRepository { return &projectRepository{db: db} }

func (r *projectRepository) Create(ctx context.Context, p *models.Project) error {
	return r.db.WithContext(ctx).Create(p).Error
}
func (r *projectRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	var p models.Project
	if err := r.db.WithContext(ctx).Preload("Organization").First(&p, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}
func (r *projectRepository) Update(ctx context.Context, p *models.Project) error {
	return r.db.WithContext(ctx).Save(p).Error
}
func (r *projectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Project{}, "id = ?", id).Error
}
func (r *projectRepository) List(ctx context.Context, offset, limit int, search string) ([]models.Project, int64, error) {
	var items []models.Project
	var total int64
	q := r.db.WithContext(ctx).Model(&models.Project{})
	if search != "" {
		like := "%" + search + "%"
		q = q.Where("name ILIKE ? OR slug ILIKE ?", like, like)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Offset(offset).Limit(limit).Order("created_at DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}
func (r *projectRepository) Count(ctx context.Context) (int64, error) {
	var n int64
	return n, r.db.WithContext(ctx).Model(&models.Project{}).Count(&n).Error
}

// --- Build ---

type BuildRepository interface {
	Create(ctx context.Context, b *models.Build) error
	List(ctx context.Context, offset, limit int, status string) ([]models.Build, int64, error)
	CountByStatus(ctx context.Context, status models.BuildStatus) (int64, error)
}

type buildRepository struct{ db *gorm.DB }

func NewBuildRepository(db *gorm.DB) BuildRepository { return &buildRepository{db: db} }

func (r *buildRepository) Create(ctx context.Context, b *models.Build) error {
	return r.db.WithContext(ctx).Create(b).Error
}
func (r *buildRepository) List(ctx context.Context, offset, limit int, status string) ([]models.Build, int64, error) {
	var items []models.Build
	var total int64
	q := r.db.WithContext(ctx).Model(&models.Build{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Offset(offset).Limit(limit).Order("created_at DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}
func (r *buildRepository) CountByStatus(ctx context.Context, status models.BuildStatus) (int64, error) {
	var n int64
	return n, r.db.WithContext(ctx).Model(&models.Build{}).Where("status = ?", status).Count(&n).Error
}

// --- Deployment ---

type DeploymentRepository interface {
	Create(ctx context.Context, d *models.Deployment) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Deployment, error)
	Update(ctx context.Context, d *models.Deployment) error
	List(ctx context.Context, offset, limit int, env string) ([]models.Deployment, int64, error)
	CountSince(ctx context.Context, since time.Time) (int64, error)
}

type deploymentRepository struct{ db *gorm.DB }

func NewDeploymentRepository(db *gorm.DB) DeploymentRepository {
	return &deploymentRepository{db: db}
}

func (r *deploymentRepository) Create(ctx context.Context, d *models.Deployment) error {
	return r.db.WithContext(ctx).Create(d).Error
}
func (r *deploymentRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Deployment, error) {
	var d models.Deployment
	if err := r.db.WithContext(ctx).First(&d, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &d, nil
}
func (r *deploymentRepository) Update(ctx context.Context, d *models.Deployment) error {
	return r.db.WithContext(ctx).Save(d).Error
}
func (r *deploymentRepository) List(ctx context.Context, offset, limit int, env string) ([]models.Deployment, int64, error) {
	var items []models.Deployment
	var total int64
	q := r.db.WithContext(ctx).Model(&models.Deployment{})
	if env != "" {
		q = q.Where("environment = ?", env)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Offset(offset).Limit(limit).Order("deployed_at DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}
func (r *deploymentRepository) CountSince(ctx context.Context, since time.Time) (int64, error) {
	var n int64
	return n, r.db.WithContext(ctx).Model(&models.Deployment{}).Where("deployed_at >= ?", since).Count(&n).Error
}

// --- Alert ---

type AlertRepository interface {
	Create(ctx context.Context, a *models.Alert) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Alert, error)
	Update(ctx context.Context, a *models.Alert) error
	List(ctx context.Context, offset, limit int, severity, status, source, search string) ([]models.Alert, int64, error)
	CountCriticalOpen(ctx context.Context) (int64, error)
}

type alertRepository struct{ db *gorm.DB }

func NewAlertRepository(db *gorm.DB) AlertRepository { return &alertRepository{db: db} }

func (r *alertRepository) Create(ctx context.Context, a *models.Alert) error {
	return r.db.WithContext(ctx).Create(a).Error
}
func (r *alertRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Alert, error) {
	var a models.Alert
	if err := r.db.WithContext(ctx).First(&a, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &a, nil
}
func (r *alertRepository) Update(ctx context.Context, a *models.Alert) error {
	return r.db.WithContext(ctx).Save(a).Error
}
func (r *alertRepository) List(ctx context.Context, offset, limit int, severity, status, source, search string) ([]models.Alert, int64, error) {
	var items []models.Alert
	var total int64
	q := r.db.WithContext(ctx).Model(&models.Alert{})
	if severity != "" {
		q = q.Where("severity = ?", severity)
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if source != "" {
		q = q.Where("source = ?", source)
	}
	if search != "" {
		like := "%" + search + "%"
		q = q.Where("title ILIKE ? OR description ILIKE ?", like, like)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Offset(offset).Limit(limit).Order("created_at DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}
func (r *alertRepository) CountCriticalOpen(ctx context.Context) (int64, error) {
	var n int64
	return n, r.db.WithContext(ctx).Model(&models.Alert{}).
		Where("severity = ? AND status = ?", models.SeverityCritical, models.AlertStatusOpen).Count(&n).Error
}

// --- Incident ---

type IncidentRepository interface {
	Create(ctx context.Context, i *models.Incident) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Incident, error)
	Update(ctx context.Context, i *models.Incident) error
	List(ctx context.Context, offset, limit int, status, priority, search string) ([]models.Incident, int64, error)
	AddComment(ctx context.Context, c *models.IncidentComment) error
	CountOpen(ctx context.Context) (int64, error)
}

type incidentRepository struct{ db *gorm.DB }

func NewIncidentRepository(db *gorm.DB) IncidentRepository {
	return &incidentRepository{db: db}
}

func (r *incidentRepository) Create(ctx context.Context, i *models.Incident) error {
	return r.db.WithContext(ctx).Create(i).Error
}
func (r *incidentRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Incident, error) {
	var i models.Incident
	err := r.db.WithContext(ctx).
		Preload("Assignee").Preload("Reporter").Preload("Comments.Author").Preload("Attachments").
		First(&i, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &i, nil
}
func (r *incidentRepository) Update(ctx context.Context, i *models.Incident) error {
	return r.db.WithContext(ctx).Save(i).Error
}
func (r *incidentRepository) List(ctx context.Context, offset, limit int, status, priority, search string) ([]models.Incident, int64, error) {
	var items []models.Incident
	var total int64
	q := r.db.WithContext(ctx).Model(&models.Incident{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if priority != "" {
		q = q.Where("priority = ?", priority)
	}
	if search != "" {
		like := "%" + search + "%"
		q = q.Where("title ILIKE ?", like)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Preload("Assignee").Offset(offset).Limit(limit).Order("created_at DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}
func (r *incidentRepository) AddComment(ctx context.Context, c *models.IncidentComment) error {
	return r.db.WithContext(ctx).Create(c).Error
}
func (r *incidentRepository) CountOpen(ctx context.Context) (int64, error) {
	var n int64
	return n, r.db.WithContext(ctx).Model(&models.Incident{}).
		Where("status IN ?", []models.IncidentStatus{
			models.IncidentStatusOpen, models.IncidentStatusInvestigating, models.IncidentStatusMitigated,
		}).Count(&n).Error
}

// --- Server ---

type ServerRepository interface {
	Upsert(ctx context.Context, s *models.Server) error
	List(ctx context.Context) ([]models.Server, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.Server, error)
	CountOnline(ctx context.Context) (int64, error)
	AvgMetrics(ctx context.Context) (cpu, mem, disk float64, netIn, netOut int64, err error)
}

type serverRepository struct{ db *gorm.DB }

func NewServerRepository(db *gorm.DB) ServerRepository { return &serverRepository{db: db} }

func (r *serverRepository) Upsert(ctx context.Context, s *models.Server) error {
	var existing models.Server
	err := r.db.WithContext(ctx).Where("hostname = ?", s.Hostname).First(&existing).Error
	if err == gorm.ErrRecordNotFound {
		return r.db.WithContext(ctx).Create(s).Error
	}
	if err != nil {
		return err
	}
	s.ID = existing.ID
	s.CreatedAt = existing.CreatedAt
	return r.db.WithContext(ctx).Save(s).Error
}
func (r *serverRepository) List(ctx context.Context) ([]models.Server, error) {
	var items []models.Server
	return items, r.db.WithContext(ctx).Order("hostname").Find(&items).Error
}
func (r *serverRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Server, error) {
	var s models.Server
	if err := r.db.WithContext(ctx).First(&s, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &s, nil
}
func (r *serverRepository) CountOnline(ctx context.Context) (int64, error) {
	var n int64
	return n, r.db.WithContext(ctx).Model(&models.Server{}).Where("status = ?", "online").Count(&n).Error
}
func (r *serverRepository) AvgMetrics(ctx context.Context) (cpu, mem, disk float64, netIn, netOut int64, err error) {
	type agg struct {
		CPU   float64
		Mem   float64
		Disk  float64
		NetIn float64
		NetOut float64
	}
	var a agg
	err = r.db.WithContext(ctx).Model(&models.Server{}).
		Select("COALESCE(AVG(cpu_percent),0) as cpu, COALESCE(AVG(memory_percent),0) as mem, COALESCE(AVG(disk_percent),0) as disk, COALESCE(SUM(network_in_bps),0) as net_in, COALESCE(SUM(network_out_bps),0) as net_out").
		Scan(&a).Error
	return a.CPU, a.Mem, a.Disk, int64(a.NetIn), int64(a.NetOut), err
}

// --- Metric ---

type MetricRepository interface {
	Create(ctx context.Context, m *models.Metric) error
	Series(ctx context.Context, name string, since time.Time) ([]models.Metric, error)
}

type metricRepository struct{ db *gorm.DB }

func NewMetricRepository(db *gorm.DB) MetricRepository { return &metricRepository{db: db} }

func (r *metricRepository) Create(ctx context.Context, m *models.Metric) error {
	return r.db.WithContext(ctx).Create(m).Error
}
func (r *metricRepository) Series(ctx context.Context, name string, since time.Time) ([]models.Metric, error) {
	var items []models.Metric
	return items, r.db.WithContext(ctx).Where("name = ? AND recorded_at >= ?", name, since).
		Order("recorded_at ASC").Find(&items).Error
}

// --- Audit ---

type AuditRepository interface {
	Create(ctx context.Context, a *models.AuditLog) error
	List(ctx context.Context, offset, limit int) ([]models.AuditLog, int64, error)
}

type auditRepository struct{ db *gorm.DB }

func NewAuditRepository(db *gorm.DB) AuditRepository { return &auditRepository{db: db} }

func (r *auditRepository) Create(ctx context.Context, a *models.AuditLog) error {
	return r.db.WithContext(ctx).Create(a).Error
}
func (r *auditRepository) List(ctx context.Context, offset, limit int) ([]models.AuditLog, int64, error) {
	var items []models.AuditLog
	var total int64
	q := r.db.WithContext(ctx).Model(&models.AuditLog{})
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Offset(offset).Limit(limit).Order("created_at DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// --- Container ---

type ContainerRepository interface {
	Upsert(ctx context.Context, c *models.Container) error
	List(ctx context.Context, search string) ([]models.Container, error)
	CountRunning(ctx context.Context) (int64, error)
}

type containerRepository struct{ db *gorm.DB }

func NewContainerRepository(db *gorm.DB) ContainerRepository {
	return &containerRepository{db: db}
}

func (r *containerRepository) Upsert(ctx context.Context, c *models.Container) error {
	var existing models.Container
	err := r.db.WithContext(ctx).Where("container_id = ?", c.ContainerID).First(&existing).Error
	if err == gorm.ErrRecordNotFound {
		return r.db.WithContext(ctx).Create(c).Error
	}
	if err != nil {
		return err
	}
	c.ID = existing.ID
	c.CreatedAt = existing.CreatedAt
	return r.db.WithContext(ctx).Save(c).Error
}
func (r *containerRepository) List(ctx context.Context, search string) ([]models.Container, error) {
	var items []models.Container
	q := r.db.WithContext(ctx).Model(&models.Container{})
	if search != "" {
		like := "%" + search + "%"
		q = q.Where("name ILIKE ? OR image ILIKE ?", like, like)
	}
	return items, q.Order("name").Find(&items).Error
}
func (r *containerRepository) CountRunning(ctx context.Context) (int64, error) {
	var n int64
	return n, r.db.WithContext(ctx).Model(&models.Container{}).Where("state = ?", "running").Count(&n).Error
}

// --- Cluster ---

type ClusterRepository interface {
	List(ctx context.Context) ([]models.Cluster, error)
	CountPodsRunning(ctx context.Context) (int64, error)
}

type clusterRepository struct{ db *gorm.DB }

func NewClusterRepository(db *gorm.DB) ClusterRepository {
	return &clusterRepository{db: db}
}

func (r *clusterRepository) List(ctx context.Context) ([]models.Cluster, error) {
	var items []models.Cluster
	return items, r.db.WithContext(ctx).Find(&items).Error
}
func (r *clusterRepository) CountPodsRunning(ctx context.Context) (int64, error) {
	var n int64
	return n, r.db.WithContext(ctx).Model(&models.Pod{}).Where("status = ?", "Running").Count(&n).Error
}
