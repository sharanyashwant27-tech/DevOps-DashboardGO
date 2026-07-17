package routes

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/devops-command-center/backend/config"
	"github.com/devops-command-center/backend/internal/auth"
	"github.com/devops-command-center/backend/internal/controllers"
	"github.com/devops-command-center/backend/internal/middleware"
	"github.com/devops-command-center/backend/internal/models"
	"github.com/devops-command-center/backend/internal/websocket"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

// Setup configures the Gin engine and all routes.
func Setup(cfg *config.Config, ctrl *controllers.Controllers, jwtMgr *auth.JWTManager, hub *websocket.Hub, log *zap.Logger) *gin.Engine {
	gin.SetMode(cfg.Server.Mode)
	r := gin.New()

	r.Use(middleware.Recovery(log))
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger(log))
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.RateLimiter(cfg.RateLimit.RequestsPerMinute, cfg.RateLimit.Burst))
	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORS.AllowedOrigins,
		AllowMethods:     cfg.CORS.AllowedMethods,
		AllowHeaders:     cfg.CORS.AllowedHeaders,
		AllowCredentials: true,
	}))

	r.GET("/health", ctrl.Health)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/ws", hub.HandleWS)

	v1 := r.Group("/api/v1")
	{
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", ctrl.Register)
			authGroup.POST("/login", ctrl.Login)
			authGroup.POST("/refresh", ctrl.Refresh)
			authGroup.POST("/forgot-password", ctrl.ForgotPassword)
			authGroup.GET("/me", middleware.JWTAuth(jwtMgr), ctrl.Me)
		}

		protected := v1.Group("")
		protected.Use(middleware.JWTAuth(jwtMgr))
		{
			protected.GET("/dashboard/stats", ctrl.DashboardStats)

			projects := protected.Group("/projects")
			{
				projects.GET("", ctrl.ListProjects)
				projects.POST("", middleware.RequireRoles(models.RoleAdmin, models.RoleDevOps, models.RoleDeveloper), ctrl.CreateProject)
				projects.GET("/:id", ctrl.GetProject)
				projects.DELETE("/:id", middleware.RequireRoles(models.RoleAdmin, models.RoleDevOps), ctrl.DeleteProject)
			}

			deployments := protected.Group("/deployments")
			{
				deployments.GET("", ctrl.ListDeployments)
				deployments.POST("", middleware.RequireRoles(models.RoleAdmin, models.RoleDevOps, models.RoleDeveloper), ctrl.CreateDeployment)
				deployments.POST("/:id/rollback", middleware.RequireRoles(models.RoleAdmin, models.RoleDevOps), ctrl.RollbackDeployment)
			}

			incidents := protected.Group("/incidents")
			{
				incidents.GET("", ctrl.ListIncidents)
				incidents.POST("", middleware.RequireRoles(models.RoleAdmin, models.RoleDevOps, models.RoleDeveloper), ctrl.CreateIncident)
				incidents.GET("/:id", ctrl.GetIncident)
				incidents.PATCH("/:id", middleware.RequireRoles(models.RoleAdmin, models.RoleDevOps, models.RoleDeveloper), ctrl.UpdateIncident)
				incidents.POST("/:id/comments", ctrl.AddIncidentComment)
			}

			alerts := protected.Group("/alerts")
			{
				alerts.GET("", ctrl.ListAlerts)
				alerts.POST("/:id/acknowledge", ctrl.AcknowledgeAlert)
				alerts.POST("/:id/resolve", middleware.RequireRoles(models.RoleAdmin, models.RoleDevOps), ctrl.ResolveAlert)
				alerts.POST("/:id/mute", middleware.RequireRoles(models.RoleAdmin, models.RoleDevOps), ctrl.MuteAlert)
			}

			jenkins := protected.Group("/jenkins")
			{
				jenkins.GET("/jobs", ctrl.JenkinsJobs)
				jenkins.GET("/jobs/:job/builds", ctrl.JenkinsBuilds)
				jenkins.GET("/queue", ctrl.JenkinsQueue)
				jenkins.GET("/stats", ctrl.JenkinsStats)
				jenkins.POST("/jobs/:job/build", middleware.RequireRoles(models.RoleAdmin, models.RoleDevOps, models.RoleDeveloper), ctrl.JenkinsTrigger)
				jenkins.POST("/jobs/:job/builds/:number/stop", middleware.RequireRoles(models.RoleAdmin, models.RoleDevOps), ctrl.JenkinsStop)
				jenkins.GET("/jobs/:job/builds/:number/console", ctrl.JenkinsConsole)
			}

			github := protected.Group("/github")
			{
				github.GET("/repos", ctrl.GitHubRepos)
				github.GET("/repos/:owner/:repo/branches", ctrl.GitHubBranches)
				github.GET("/repos/:owner/:repo/commits", ctrl.GitHubCommits)
				github.GET("/repos/:owner/:repo/pulls", ctrl.GitHubPRs)
				github.GET("/repos/:owner/:repo/issues", ctrl.GitHubIssues)
				github.GET("/repos/:owner/:repo/releases", ctrl.GitHubReleases)
				github.GET("/repos/:owner/:repo/actions/runs", ctrl.GitHubWorkflows)
				github.GET("/repos/:owner/:repo/contributors", ctrl.GitHubContributors)
				github.GET("/repos/:owner/:repo/health", ctrl.GitHubHealth)
			}

			docker := protected.Group("/docker")
			{
				docker.GET("/containers", ctrl.DockerContainers)
				docker.GET("/images", ctrl.DockerImages)
				docker.GET("/volumes", ctrl.DockerVolumes)
				docker.GET("/networks", ctrl.DockerNetworks)
				docker.GET("/containers/:id/stats", ctrl.DockerStats)
				docker.GET("/containers/:id/logs", ctrl.DockerLogs)
				docker.POST("/containers/:id/start", middleware.RequireRoles(models.RoleAdmin, models.RoleDevOps), ctrl.DockerStart)
				docker.POST("/containers/:id/stop", middleware.RequireRoles(models.RoleAdmin, models.RoleDevOps), ctrl.DockerStop)
				docker.POST("/containers/:id/restart", middleware.RequireRoles(models.RoleAdmin, models.RoleDevOps), ctrl.DockerRestart)
				docker.DELETE("/containers/:id", middleware.RequireRoles(models.RoleAdmin, models.RoleDevOps), ctrl.DockerDelete)
			}

			k8s := protected.Group("/kubernetes")
			{
				k8s.GET("/namespaces", ctrl.K8sNamespaces)
				k8s.GET("/pods", ctrl.K8sPods)
				k8s.GET("/deployments", ctrl.K8sDeployments)
				k8s.GET("/replicasets", ctrl.K8sReplicaSets)
				k8s.GET("/daemonsets", ctrl.K8sDaemonSets)
				k8s.GET("/services", ctrl.K8sServices)
				k8s.GET("/ingresses", ctrl.K8sIngresses)
				k8s.GET("/nodes", ctrl.K8sNodes)
				k8s.GET("/persistentvolumes", ctrl.K8sPVs)
				k8s.GET("/persistentvolumeclaims", ctrl.K8sPVCs)
				k8s.GET("/events", ctrl.K8sEvents)
				k8s.GET("/pods/:pod/logs", ctrl.K8sPodLogs)
				k8s.POST("/deployments/:name/scale", middleware.RequireRoles(models.RoleAdmin, models.RoleDevOps), ctrl.K8sScale)
				k8s.POST("/deployments/:name/restart", middleware.RequireRoles(models.RoleAdmin, models.RoleDevOps), ctrl.K8sRestart)
				k8s.DELETE("/pods/:pod", middleware.RequireRoles(models.RoleAdmin, models.RoleDevOps), ctrl.K8sDeletePod)
			}

			protected.GET("/servers", ctrl.ListServers)
			protected.GET("/servers/local", ctrl.ServerDetails)
			protected.GET("/metrics/:name", ctrl.MetricSeries)
			protected.GET("/audit", middleware.RequireRoles(models.RoleAdmin, models.RoleDevOps), ctrl.ListAudit)
		}
	}

	mountFrontend(r, log)
	return r
}

// mountFrontend serves the React SPA from frontend/dist (or DCC_FRONTEND_DIR).
func mountFrontend(r *gin.Engine, log *zap.Logger) {
	dist := os.Getenv("DCC_FRONTEND_DIR")
	if dist == "" {
		candidates := []string{
			filepath.Join("..", "frontend", "dist"),
			filepath.Join("frontend", "dist"),
			filepath.Join("web"),
		}
		for _, c := range candidates {
			if info, err := os.Stat(filepath.Join(c, "index.html")); err == nil && !info.IsDir() {
				dist = c
				break
			}
		}
	}
	if dist == "" {
		log.Warn("frontend dist not found; GET / will 404 until frontend is built")
		r.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "DevOps Command Center API",
				"hint":    "Build the UI with: cd frontend && npm install && npm run build",
				"swagger": "/swagger/index.html",
				"health":  "/health",
			})
		})
		return
	}

	abs, err := filepath.Abs(dist)
	if err != nil {
		abs = dist
	}
	log.Info("serving frontend", zap.String("dir", abs))
	r.Static("/assets", filepath.Join(abs, "assets"))
	r.StaticFile("/favicon.ico", filepath.Join(abs, "favicon.ico"))

	index := filepath.Join(abs, "index.html")
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/api") ||
			strings.HasPrefix(path, "/swagger") ||
			strings.HasPrefix(path, "/health") ||
			strings.HasPrefix(path, "/metrics") ||
			strings.HasPrefix(path, "/ws") {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "not found"})
			return
		}
		// Prefer static file if it exists (e.g. vite.svg)
		candidate := filepath.Join(abs, filepath.Clean("/"+path))
		if !strings.HasPrefix(candidate, abs) {
			c.File(index)
			return
		}
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			c.File(candidate)
			return
		}
		c.File(index)
	})
}
