package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/devops-command-center/backend/internal/dto"
	"github.com/devops-command-center/backend/internal/middleware"
	"github.com/devops-command-center/backend/internal/services"
	"github.com/devops-command-center/backend/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Controllers groups HTTP handlers.
type Controllers struct {
	svc *services.Services
	log *zap.Logger
}

func New(svc *services.Services, log *zap.Logger) *Controllers {
	return &Controllers{svc: svc, log: log}
}

// --- Auth ---

// Register godoc
// @Summary Register user
// @Tags auth
// @Accept json
// @Produce json
// @Param body body dto.RegisterRequest true "register"
// @Success 201 {object} response.APIResponse
// @Router /api/v1/auth/register [post]
func (h *Controllers) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	out, err := h.svc.Auth.Register(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, services.ErrEmailTaken) {
			response.Fail(c, http.StatusConflict, "EMAIL_TAKEN", err.Error(), "")
			return
		}
		response.Internal(c, err.Error())
		return
	}
	response.Created(c, "registered", out)
}

// Login godoc
// @Summary Login
// @Tags auth
// @Accept json
// @Produce json
// @Param body body dto.LoginRequest true "login"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/auth/login [post]
func (h *Controllers) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	out, err := h.svc.Auth.Login(c.Request.Context(), req, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}
	response.OK(c, "login successful", out)
}

func (h *Controllers) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	out, err := h.svc.Auth.Refresh(c.Request.Context(), req.RefreshToken, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}
	response.OK(c, "token refreshed", out)
}

func (h *Controllers) ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	_ = h.svc.Auth.ForgotPassword(c.Request.Context(), req.Email)
	response.OK(c, "if the email exists, a reset link will be sent", nil)
}

func (h *Controllers) Me(c *gin.Context) {
	user, err := h.svc.Auth.Me(c.Request.Context(), middleware.GetUserID(c))
	if err != nil {
		response.NotFound(c, "user not found")
		return
	}
	response.OK(c, "ok", user)
}

// --- Dashboard ---

func (h *Controllers) DashboardStats(c *gin.Context) {
	stats, err := h.svc.Dashboard.Stats(c.Request.Context())
	if err != nil {
		response.Internal(c, err.Error())
		return
	}
	response.OK(c, "dashboard stats", stats)
}

// --- Projects ---

func (h *Controllers) CreateProject(c *gin.Context) {
	var req dto.ProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	p, err := h.svc.Project.Create(c.Request.Context(), req, middleware.GetUserID(c))
	if err != nil {
		response.Internal(c, err.Error())
		return
	}
	response.Created(c, "project created", p)
}

func (h *Controllers) ListProjects(c *gin.Context) {
	var q dto.PaginationQuery
	_ = c.ShouldBindQuery(&q)
	if q.Page == 0 {
		q.Page = 1
	}
	if q.PageSize == 0 {
		q.PageSize = 20
	}
	items, total, err := h.svc.Project.List(c.Request.Context(), q)
	if err != nil {
		response.Internal(c, err.Error())
		return
	}
	response.Paginated(c, "projects", items, q.Page, q.PageSize, total)
}

func (h *Controllers) GetProject(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	p, err := h.svc.Project.Get(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "project not found")
		return
	}
	response.OK(c, "ok", p)
}

func (h *Controllers) DeleteProject(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	if err := h.svc.Project.Delete(c.Request.Context(), id, middleware.GetUserID(c)); err != nil {
		response.Internal(c, err.Error())
		return
	}
	response.OK(c, "deleted", nil)
}

// --- Deployments ---

func (h *Controllers) CreateDeployment(c *gin.Context) {
	var req dto.DeploymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	user, _ := h.svc.Auth.Me(c.Request.Context(), middleware.GetUserID(c))
	name := ""
	if user != nil {
		name = user.Name
	}
	d, err := h.svc.Deployment.Create(c.Request.Context(), req, middleware.GetUserID(c), name)
	if err != nil {
		response.Internal(c, err.Error())
		return
	}
	response.Created(c, "deployment recorded", d)
}

func (h *Controllers) ListDeployments(c *gin.Context) {
	var q dto.PaginationQuery
	_ = c.ShouldBindQuery(&q)
	if q.Page == 0 {
		q.Page = 1
	}
	if q.PageSize == 0 {
		q.PageSize = 20
	}
	items, total, err := h.svc.Deployment.List(c.Request.Context(), q, c.Query("environment"))
	if err != nil {
		response.Internal(c, err.Error())
		return
	}
	response.Paginated(c, "deployments", items, q.Page, q.PageSize, total)
}

func (h *Controllers) RollbackDeployment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	d, err := h.svc.Deployment.Rollback(c.Request.Context(), id, middleware.GetUserID(c))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OK(c, "rollback initiated", d)
}

// --- Incidents ---

func (h *Controllers) CreateIncident(c *gin.Context) {
	var req dto.IncidentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	inc, err := h.svc.Incident.Create(c.Request.Context(), req, middleware.GetUserID(c))
	if err != nil {
		response.Internal(c, err.Error())
		return
	}
	response.Created(c, "incident created", inc)
}

func (h *Controllers) ListIncidents(c *gin.Context) {
	var q dto.PaginationQuery
	_ = c.ShouldBindQuery(&q)
	if q.Page == 0 {
		q.Page = 1
	}
	if q.PageSize == 0 {
		q.PageSize = 20
	}
	items, total, err := h.svc.Incident.List(c.Request.Context(), q, c.Query("status"), c.Query("priority"))
	if err != nil {
		response.Internal(c, err.Error())
		return
	}
	response.Paginated(c, "incidents", items, q.Page, q.PageSize, total)
}

func (h *Controllers) GetIncident(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	inc, err := h.svc.Incident.Get(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "incident not found")
		return
	}
	response.OK(c, "ok", inc)
}

func (h *Controllers) UpdateIncident(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	var req dto.IncidentUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	inc, err := h.svc.Incident.Update(c.Request.Context(), id, req)
	if err != nil {
		response.Internal(c, err.Error())
		return
	}
	response.OK(c, "updated", inc)
}

func (h *Controllers) AddIncidentComment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	var req dto.IncidentCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	comment, err := h.svc.Incident.AddComment(c.Request.Context(), id, middleware.GetUserID(c), req.Body)
	if err != nil {
		response.Internal(c, err.Error())
		return
	}
	response.Created(c, "comment added", comment)
}

// --- Alerts ---

func (h *Controllers) ListAlerts(c *gin.Context) {
	var q dto.PaginationQuery
	_ = c.ShouldBindQuery(&q)
	if q.Page == 0 {
		q.Page = 1
	}
	if q.PageSize == 0 {
		q.PageSize = 20
	}
	items, total, err := h.svc.Alert.List(c.Request.Context(), q, c.Query("severity"), c.Query("status"), c.Query("source"))
	if err != nil {
		response.Internal(c, err.Error())
		return
	}
	response.Paginated(c, "alerts", items, q.Page, q.PageSize, total)
}

func (h *Controllers) AcknowledgeAlert(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	a, err := h.svc.Alert.Acknowledge(c.Request.Context(), id, middleware.GetUserID(c))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.OK(c, "acknowledged", a)
}

func (h *Controllers) ResolveAlert(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	a, err := h.svc.Alert.Resolve(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.OK(c, "resolved", a)
}

func (h *Controllers) MuteAlert(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	var req dto.AlertMuteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	a, err := h.svc.Alert.Mute(c.Request.Context(), id, req.Minutes)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.OK(c, "muted", a)
}

// --- Jenkins ---

func (h *Controllers) JenkinsJobs(c *gin.Context) {
	jobs, err := h.svc.Jenkins.ListJobs(c.Request.Context(), c.Query("search"))
	if err != nil {
		response.Fail(c, http.StatusServiceUnavailable, "JENKINS_UNAVAILABLE", err.Error(), "")
		return
	}
	response.OK(c, "jenkins jobs", jobs)
}

func (h *Controllers) JenkinsBuilds(c *gin.Context) {
	builds, err := h.svc.Jenkins.GetJobBuilds(c.Request.Context(), c.Param("job"))
	if err != nil {
		response.Fail(c, http.StatusServiceUnavailable, "JENKINS_UNAVAILABLE", err.Error(), "")
		return
	}
	response.OK(c, "builds", builds)
}

func (h *Controllers) JenkinsQueue(c *gin.Context) {
	queue, err := h.svc.Jenkins.GetQueue(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusServiceUnavailable, "JENKINS_UNAVAILABLE", err.Error(), "")
		return
	}
	response.OK(c, "queue", queue)
}

func (h *Controllers) JenkinsTrigger(c *gin.Context) {
	if err := h.svc.Jenkins.TriggerBuild(c.Request.Context(), c.Param("job")); err != nil {
		response.Fail(c, http.StatusBadRequest, "TRIGGER_FAILED", err.Error(), "")
		return
	}
	response.OK(c, "build triggered", nil)
}

func (h *Controllers) JenkinsStop(c *gin.Context) {
	num, _ := strconv.Atoi(c.Param("number"))
	if err := h.svc.Jenkins.StopBuild(c.Request.Context(), c.Param("job"), num); err != nil {
		response.Fail(c, http.StatusBadRequest, "STOP_FAILED", err.Error(), "")
		return
	}
	response.OK(c, "build stopped", nil)
}

func (h *Controllers) JenkinsConsole(c *gin.Context) {
	num, _ := strconv.Atoi(c.Param("number"))
	log, err := h.svc.Jenkins.ConsoleLog(c.Request.Context(), c.Param("job"), num)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "LOG_FAILED", err.Error(), "")
		return
	}
	response.OK(c, "console", gin.H{"log": log})
}

func (h *Controllers) JenkinsStats(c *gin.Context) {
	stats, err := h.svc.Jenkins.Stats(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusServiceUnavailable, "JENKINS_UNAVAILABLE", err.Error(), "")
		return
	}
	response.OK(c, "stats", stats)
}

// --- GitHub ---

func (h *Controllers) GitHubRepos(c *gin.Context) {
	repos, err := h.svc.GitHub.ListRepos(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusServiceUnavailable, "GITHUB_UNAVAILABLE", err.Error(), "")
		return
	}
	response.OK(c, "repositories", repos)
}

func (h *Controllers) GitHubBranches(c *gin.Context) {
	items, err := h.svc.GitHub.ListBranches(c.Request.Context(), c.Param("owner"), c.Param("repo"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "GITHUB_ERROR", err.Error(), "")
		return
	}
	response.OK(c, "branches", items)
}

func (h *Controllers) GitHubCommits(c *gin.Context) {
	items, err := h.svc.GitHub.ListCommits(c.Request.Context(), c.Param("owner"), c.Param("repo"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "GITHUB_ERROR", err.Error(), "")
		return
	}
	response.OK(c, "commits", items)
}

func (h *Controllers) GitHubPRs(c *gin.Context) {
	items, err := h.svc.GitHub.ListPullRequests(c.Request.Context(), c.Param("owner"), c.Param("repo"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "GITHUB_ERROR", err.Error(), "")
		return
	}
	response.OK(c, "pull_requests", items)
}

func (h *Controllers) GitHubIssues(c *gin.Context) {
	items, err := h.svc.GitHub.ListIssues(c.Request.Context(), c.Param("owner"), c.Param("repo"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "GITHUB_ERROR", err.Error(), "")
		return
	}
	response.OK(c, "issues", items)
}

func (h *Controllers) GitHubReleases(c *gin.Context) {
	items, err := h.svc.GitHub.ListReleases(c.Request.Context(), c.Param("owner"), c.Param("repo"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "GITHUB_ERROR", err.Error(), "")
		return
	}
	response.OK(c, "releases", items)
}

func (h *Controllers) GitHubWorkflows(c *gin.Context) {
	items, err := h.svc.GitHub.ListWorkflowRuns(c.Request.Context(), c.Param("owner"), c.Param("repo"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "GITHUB_ERROR", err.Error(), "")
		return
	}
	response.OK(c, "workflow_runs", items)
}

func (h *Controllers) GitHubContributors(c *gin.Context) {
	items, err := h.svc.GitHub.ListContributors(c.Request.Context(), c.Param("owner"), c.Param("repo"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "GITHUB_ERROR", err.Error(), "")
		return
	}
	response.OK(c, "contributors", items)
}

func (h *Controllers) GitHubHealth(c *gin.Context) {
	items, err := h.svc.GitHub.RepositoryHealth(c.Request.Context(), c.Param("owner"), c.Param("repo"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "GITHUB_ERROR", err.Error(), "")
		return
	}
	response.OK(c, "health", items)
}

// --- Docker ---

func (h *Controllers) DockerContainers(c *gin.Context) {
	items, err := h.svc.Docker.ListContainers(c.Request.Context(), c.Query("search"))
	if err != nil {
		response.Fail(c, http.StatusServiceUnavailable, "DOCKER_UNAVAILABLE", err.Error(), "")
		return
	}
	response.OK(c, "containers", items)
}

func (h *Controllers) DockerImages(c *gin.Context) {
	items, err := h.svc.Docker.ListImages(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusServiceUnavailable, "DOCKER_UNAVAILABLE", err.Error(), "")
		return
	}
	response.OK(c, "images", items)
}

func (h *Controllers) DockerVolumes(c *gin.Context) {
	items, err := h.svc.Docker.ListVolumes(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusServiceUnavailable, "DOCKER_UNAVAILABLE", err.Error(), "")
		return
	}
	response.OK(c, "volumes", items)
}

func (h *Controllers) DockerNetworks(c *gin.Context) {
	items, err := h.svc.Docker.ListNetworks(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusServiceUnavailable, "DOCKER_UNAVAILABLE", err.Error(), "")
		return
	}
	response.OK(c, "networks", items)
}

func (h *Controllers) DockerStats(c *gin.Context) {
	stats, err := h.svc.Docker.ContainerStats(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "DOCKER_ERROR", err.Error(), "")
		return
	}
	response.OK(c, "stats", stats)
}

func (h *Controllers) DockerLogs(c *gin.Context) {
	tail := c.DefaultQuery("tail", "200")
	logs, err := h.svc.Docker.Logs(c.Request.Context(), c.Param("id"), tail)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "DOCKER_ERROR", err.Error(), "")
		return
	}
	response.OK(c, "logs", gin.H{"log": logs})
}

func (h *Controllers) DockerStart(c *gin.Context) {
	if err := h.svc.Docker.Start(c.Request.Context(), c.Param("id")); err != nil {
		response.Fail(c, http.StatusBadRequest, "DOCKER_ERROR", err.Error(), "")
		return
	}
	response.OK(c, "started", nil)
}

func (h *Controllers) DockerStop(c *gin.Context) {
	if err := h.svc.Docker.Stop(c.Request.Context(), c.Param("id")); err != nil {
		response.Fail(c, http.StatusBadRequest, "DOCKER_ERROR", err.Error(), "")
		return
	}
	response.OK(c, "stopped", nil)
}

func (h *Controllers) DockerRestart(c *gin.Context) {
	if err := h.svc.Docker.Restart(c.Request.Context(), c.Param("id")); err != nil {
		response.Fail(c, http.StatusBadRequest, "DOCKER_ERROR", err.Error(), "")
		return
	}
	response.OK(c, "restarted", nil)
}

func (h *Controllers) DockerDelete(c *gin.Context) {
	force := c.Query("force") == "true"
	if err := h.svc.Docker.Delete(c.Request.Context(), c.Param("id"), force); err != nil {
		response.Fail(c, http.StatusBadRequest, "DOCKER_ERROR", err.Error(), "")
		return
	}
	response.OK(c, "deleted", nil)
}

// --- Kubernetes ---

func (h *Controllers) K8sNamespaces(c *gin.Context) {
	items, err := h.svc.Kubernetes.ListNamespaces(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusServiceUnavailable, "K8S_UNAVAILABLE", err.Error(), "")
		return
	}
	response.OK(c, "namespaces", items)
}

func (h *Controllers) K8sPods(c *gin.Context) {
	ns := c.DefaultQuery("namespace", "")
	items, err := h.svc.Kubernetes.ListPods(c.Request.Context(), ns)
	if err != nil {
		response.Fail(c, http.StatusServiceUnavailable, "K8S_UNAVAILABLE", err.Error(), "")
		return
	}
	response.OK(c, "pods", items)
}

func (h *Controllers) K8sDeployments(c *gin.Context) {
	ns := c.DefaultQuery("namespace", "default")
	items, err := h.svc.Kubernetes.ListDeployments(c.Request.Context(), ns)
	if err != nil {
		response.Fail(c, http.StatusServiceUnavailable, "K8S_UNAVAILABLE", err.Error(), "")
		return
	}
	response.OK(c, "deployments", items)
}

func (h *Controllers) K8sReplicaSets(c *gin.Context) {
	ns := c.DefaultQuery("namespace", "default")
	items, err := h.svc.Kubernetes.ListReplicaSets(c.Request.Context(), ns)
	if err != nil {
		response.Fail(c, http.StatusServiceUnavailable, "K8S_UNAVAILABLE", err.Error(), "")
		return
	}
	response.OK(c, "replicasets", items)
}

func (h *Controllers) K8sDaemonSets(c *gin.Context) {
	ns := c.DefaultQuery("namespace", "default")
	items, err := h.svc.Kubernetes.ListDaemonSets(c.Request.Context(), ns)
	if err != nil {
		response.Fail(c, http.StatusServiceUnavailable, "K8S_UNAVAILABLE", err.Error(), "")
		return
	}
	response.OK(c, "daemonsets", items)
}

func (h *Controllers) K8sServices(c *gin.Context) {
	ns := c.DefaultQuery("namespace", "default")
	items, err := h.svc.Kubernetes.ListServices(c.Request.Context(), ns)
	if err != nil {
		response.Fail(c, http.StatusServiceUnavailable, "K8S_UNAVAILABLE", err.Error(), "")
		return
	}
	response.OK(c, "services", items)
}

func (h *Controllers) K8sIngresses(c *gin.Context) {
	ns := c.DefaultQuery("namespace", "default")
	items, err := h.svc.Kubernetes.ListIngresses(c.Request.Context(), ns)
	if err != nil {
		response.Fail(c, http.StatusServiceUnavailable, "K8S_UNAVAILABLE", err.Error(), "")
		return
	}
	response.OK(c, "ingresses", items)
}

func (h *Controllers) K8sNodes(c *gin.Context) {
	items, err := h.svc.Kubernetes.ListNodes(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusServiceUnavailable, "K8S_UNAVAILABLE", err.Error(), "")
		return
	}
	response.OK(c, "nodes", items)
}

func (h *Controllers) K8sPVs(c *gin.Context) {
	items, err := h.svc.Kubernetes.ListPVs(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusServiceUnavailable, "K8S_UNAVAILABLE", err.Error(), "")
		return
	}
	response.OK(c, "persistent_volumes", items)
}

func (h *Controllers) K8sPVCs(c *gin.Context) {
	ns := c.DefaultQuery("namespace", "default")
	items, err := h.svc.Kubernetes.ListPVCs(c.Request.Context(), ns)
	if err != nil {
		response.Fail(c, http.StatusServiceUnavailable, "K8S_UNAVAILABLE", err.Error(), "")
		return
	}
	response.OK(c, "persistent_volume_claims", items)
}

func (h *Controllers) K8sEvents(c *gin.Context) {
	ns := c.DefaultQuery("namespace", "default")
	items, err := h.svc.Kubernetes.ListEvents(c.Request.Context(), ns)
	if err != nil {
		response.Fail(c, http.StatusServiceUnavailable, "K8S_UNAVAILABLE", err.Error(), "")
		return
	}
	response.OK(c, "events", items)
}

func (h *Controllers) K8sPodLogs(c *gin.Context) {
	ns := c.DefaultQuery("namespace", "default")
	tail, _ := strconv.ParseInt(c.DefaultQuery("tail", "200"), 10, 64)
	logs, err := h.svc.Kubernetes.PodLogs(c.Request.Context(), ns, c.Param("pod"), c.Query("container"), tail)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "K8S_ERROR", err.Error(), "")
		return
	}
	response.OK(c, "logs", gin.H{"log": logs})
}

func (h *Controllers) K8sScale(c *gin.Context) {
	var req dto.ScaleDeploymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	ns := c.DefaultQuery("namespace", "default")
	if err := h.svc.Kubernetes.ScaleDeployment(c.Request.Context(), ns, c.Param("name"), req.Replicas); err != nil {
		response.Fail(c, http.StatusBadRequest, "K8S_ERROR", err.Error(), "")
		return
	}
	response.OK(c, "scaled", nil)
}

func (h *Controllers) K8sRestart(c *gin.Context) {
	ns := c.DefaultQuery("namespace", "default")
	if err := h.svc.Kubernetes.RestartDeployment(c.Request.Context(), ns, c.Param("name")); err != nil {
		response.Fail(c, http.StatusBadRequest, "K8S_ERROR", err.Error(), "")
		return
	}
	response.OK(c, "restarted", nil)
}

func (h *Controllers) K8sDeletePod(c *gin.Context) {
	ns := c.DefaultQuery("namespace", "default")
	if err := h.svc.Kubernetes.DeletePod(c.Request.Context(), ns, c.Param("pod")); err != nil {
		response.Fail(c, http.StatusBadRequest, "K8S_ERROR", err.Error(), "")
		return
	}
	response.OK(c, "deleted", nil)
}

// --- Servers / Metrics / Audit ---

func (h *Controllers) ListServers(c *gin.Context) {
	items, err := h.svc.Server.List(c.Request.Context())
	if err != nil {
		response.Internal(c, err.Error())
		return
	}
	response.OK(c, "servers", items)
}

func (h *Controllers) ServerDetails(c *gin.Context) {
	_, details, err := h.svc.Server.CollectLocal(c.Request.Context())
	if err != nil {
		response.Internal(c, err.Error())
		return
	}
	response.OK(c, "server details", details)
}

func (h *Controllers) MetricSeries(c *gin.Context) {
	hours, _ := strconv.Atoi(c.DefaultQuery("hours", "24"))
	series, err := h.svc.Metrics.Series(c.Request.Context(), c.Param("name"), hours)
	if err != nil {
		response.Internal(c, err.Error())
		return
	}
	response.OK(c, "series", series)
}

func (h *Controllers) ListAudit(c *gin.Context) {
	var q dto.PaginationQuery
	_ = c.ShouldBindQuery(&q)
	if q.Page == 0 {
		q.Page = 1
	}
	if q.PageSize == 0 {
		q.PageSize = 20
	}
	items, total, err := h.svc.Audit.List(c.Request.Context(), q)
	if err != nil {
		response.Internal(c, err.Error())
		return
	}
	response.Paginated(c, "audit_logs", items, q.Page, q.PageSize, total)
}

func (h *Controllers) Health(c *gin.Context) {
	response.OK(c, "healthy", gin.H{
		"service":    "devops-command-center",
		"jenkins":    gin.H{"enabled": h.svc.Jenkins.Enabled(), "mode": h.svc.Jenkins.Mode()},
		"github":     gin.H{"enabled": h.svc.GitHub.Enabled(), "mode": h.svc.GitHub.Mode()},
		"docker":     h.svc.Docker.Health(c.Request.Context()),
		"kubernetes": h.svc.Kubernetes.Enabled(),
	})
}
