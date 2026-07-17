package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds all application configuration.
type Config struct {
	Server     ServerConfig     `mapstructure:"server"`
	Database   DatabaseConfig   `mapstructure:"database"`
	Redis      RedisConfig      `mapstructure:"redis"`
	JWT        JWTConfig        `mapstructure:"jwt"`
	CORS       CORSConfig       `mapstructure:"cors"`
	RateLimit  RateLimitConfig  `mapstructure:"rate_limit"`
	Jenkins    JenkinsConfig    `mapstructure:"jenkins"`
	GitHub     GitHubConfig     `mapstructure:"github"`
	Docker     DockerConfig     `mapstructure:"docker"`
	Kubernetes KubernetesConfig `mapstructure:"kubernetes"`
	SMTP       SMTPConfig       `mapstructure:"smtp"`
	Slack      SlackConfig      `mapstructure:"slack"`
	Teams      TeamsConfig      `mapstructure:"teams"`
	Scheduler  SchedulerConfig  `mapstructure:"scheduler"`
	Seed       SeedConfig       `mapstructure:"seed"`
}

type ServerConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Mode         string `mapstructure:"mode"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
}

func (s ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

type DatabaseConfig struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	Name            string `mapstructure:"name"`
	SSLMode         string `mapstructure:"sslmode"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}

func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode,
	)
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

func (r RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

type JWTConfig struct {
	AccessSecret     string `mapstructure:"access_secret"`
	RefreshSecret    string `mapstructure:"refresh_secret"`
	AccessTTLMinutes int    `mapstructure:"access_ttl_minutes"`
	RefreshTTLHours  int    `mapstructure:"refresh_ttl_hours"`
}

func (j JWTConfig) AccessTTL() time.Duration {
	return time.Duration(j.AccessTTLMinutes) * time.Minute
}

func (j JWTConfig) RefreshTTL() time.Duration {
	return time.Duration(j.RefreshTTLHours) * time.Hour
}

type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
	AllowedMethods []string `mapstructure:"allowed_methods"`
	AllowedHeaders []string `mapstructure:"allowed_headers"`
}

type RateLimitConfig struct {
	RequestsPerMinute int `mapstructure:"requests_per_minute"`
	Burst             int `mapstructure:"burst"`
}

type JenkinsConfig struct {
	URL      string `mapstructure:"url"`
	Username string `mapstructure:"username"`
	Token    string `mapstructure:"token"`
	DemoMode bool   `mapstructure:"demo_mode"`
}

func (j JenkinsConfig) Enabled() bool {
	return strings.TrimSpace(j.URL) != "" || j.DemoMode
}

func (j JenkinsConfig) Live() bool {
	return strings.TrimSpace(j.URL) != ""
}

type GitHubConfig struct {
	Token    string `mapstructure:"token"`
	Username string `mapstructure:"username"`
	BaseURL  string `mapstructure:"base_url"`
	DemoMode bool   `mapstructure:"demo_mode"`
}

func (g GitHubConfig) Enabled() bool {
	return strings.TrimSpace(g.Token) != "" || strings.TrimSpace(g.Username) != "" || g.DemoMode
}

func (g GitHubConfig) Live() bool {
	return strings.TrimSpace(g.Token) != ""
}

// PublicUser returns the configured GitHub username for unauthenticated public API access.
func (g GitHubConfig) PublicUser() string {
	return strings.TrimSpace(g.Username)
}

type DockerConfig struct {
	Host    string `mapstructure:"host"`
	Enabled bool   `mapstructure:"enabled"`
}

type KubernetesConfig struct {
	Kubeconfig string `mapstructure:"kubeconfig"`
	InCluster  bool   `mapstructure:"in_cluster"`
	Enabled    bool   `mapstructure:"enabled"`
}

type SMTPConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	From     string `mapstructure:"from"`
}

func (s SMTPConfig) Enabled() bool {
	return strings.TrimSpace(s.Host) != ""
}

type SlackConfig struct {
	WebhookURL string `mapstructure:"webhook_url"`
}

func (s SlackConfig) Enabled() bool {
	return strings.TrimSpace(s.WebhookURL) != ""
}

type TeamsConfig struct {
	WebhookURL string `mapstructure:"webhook_url"`
}

func (t TeamsConfig) Enabled() bool {
	return strings.TrimSpace(t.WebhookURL) != ""
}

type SchedulerConfig struct {
	MetricsInterval  string `mapstructure:"metrics_interval"`
	AlertInterval    string `mapstructure:"alert_interval"`
	CleanupInterval  string `mapstructure:"cleanup_interval"`
}

type SeedConfig struct {
	Enabled       bool   `mapstructure:"enabled"`
	AdminEmail    string `mapstructure:"admin_email"`
	AdminPassword string `mapstructure:"admin_password"`
	AdminName     string `mapstructure:"admin_name"`
}

// Load reads configuration from file and environment variables.
func Load(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	v.SetEnvPrefix("DCC")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	// Environment overrides for secrets
	bindEnv(v, "database.host", "DCC_DATABASE_HOST")
	bindEnv(v, "database.port", "DCC_DATABASE_PORT")
	bindEnv(v, "database.user", "DCC_DATABASE_USER")
	bindEnv(v, "database.password", "DCC_DATABASE_PASSWORD")
	bindEnv(v, "database.name", "DCC_DATABASE_NAME")
	bindEnv(v, "redis.host", "DCC_REDIS_HOST")
	bindEnv(v, "redis.port", "DCC_REDIS_PORT")
	bindEnv(v, "redis.password", "DCC_REDIS_PASSWORD")
	bindEnv(v, "server.port", "DCC_SERVER_PORT")
	bindEnv(v, "jenkins.url", "DCC_JENKINS_URL")
	bindEnv(v, "jenkins.demo_mode", "DCC_JENKINS_DEMO_MODE")
	bindEnv(v, "jwt.access_secret", "DCC_JWT_ACCESS_SECRET")
	bindEnv(v, "jwt.refresh_secret", "DCC_JWT_REFRESH_SECRET")
	bindEnv(v, "jenkins.url", "DCC_JENKINS_URL")
	bindEnv(v, "jenkins.username", "DCC_JENKINS_USERNAME")
	bindEnv(v, "jenkins.token", "DCC_JENKINS_TOKEN")
	bindEnv(v, "github.token", "DCC_GITHUB_TOKEN")
	bindEnv(v, "github.username", "DCC_GITHUB_USERNAME")
	bindEnv(v, "github.demo_mode", "DCC_GITHUB_DEMO_MODE")
	bindEnv(v, "slack.webhook_url", "DCC_SLACK_WEBHOOK_URL")
	bindEnv(v, "teams.webhook_url", "DCC_TEAMS_WEBHOOK_URL")

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}
	return &cfg, nil
}

func bindEnv(v *viper.Viper, key, env string) {
	_ = v.BindEnv(key, env)
}
