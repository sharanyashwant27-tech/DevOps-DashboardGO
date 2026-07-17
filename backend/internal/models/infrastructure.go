package models

import (
	"time"

	"github.com/google/uuid"
)

// DockerHost represents a Docker daemon endpoint.
type DockerHost struct {
	Base
	Name        string `gorm:"size:255;not null" json:"name"`
	Host        string `gorm:"size:512;not null" json:"host"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
	Description string `gorm:"type:text" json:"description,omitempty"`
}

func (DockerHost) TableName() string { return "docker_hosts" }

// Container caches Docker container metadata.
type Container struct {
	Base
	DockerHostID uuid.UUID `gorm:"type:uuid;index;not null" json:"docker_host_id"`
	ContainerID  string    `gorm:"size:128;index;not null" json:"container_id"`
	Name         string    `gorm:"size:255;index;not null" json:"name"`
	Image        string    `gorm:"size:512" json:"image"`
	Status       string    `gorm:"size:50;index" json:"status"`
	State        string    `gorm:"size:50" json:"state"`
	CPUPercent   float64   `json:"cpu_percent"`
	MemoryPercent float64  `json:"memory_percent"`
	MemoryUsage  int64     `json:"memory_usage"`
	Ports        string    `gorm:"type:text" json:"ports,omitempty"`
	DockerHost   DockerHost `gorm:"foreignKey:DockerHostID" json:"docker_host,omitempty"`
}

func (Container) TableName() string { return "containers" }

// Cluster represents a Kubernetes cluster.
type Cluster struct {
	Base
	Name        string `gorm:"size:255;not null;uniqueIndex" json:"name"`
	APIServer   string `gorm:"size:512" json:"api_server,omitempty"`
	Context     string `gorm:"size:255" json:"context,omitempty"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
	Version     string `gorm:"size:50" json:"version,omitempty"`
	NodeCount   int    `json:"node_count"`
	Description string `gorm:"type:text" json:"description,omitempty"`
}

func (Cluster) TableName() string { return "clusters" }

// Pod caches Kubernetes pod info.
type Pod struct {
	Base
	ClusterID   uuid.UUID `gorm:"type:uuid;index;not null" json:"cluster_id"`
	Name        string    `gorm:"size:255;index;not null" json:"name"`
	Namespace   string    `gorm:"size:255;index;not null" json:"namespace"`
	Status      string    `gorm:"size:50;index" json:"status"`
	NodeName    string    `gorm:"size:255" json:"node_name,omitempty"`
	Restarts    int32     `json:"restarts"`
	CPURequest  string    `gorm:"size:50" json:"cpu_request,omitempty"`
	MemoryRequest string  `gorm:"size:50" json:"memory_request,omitempty"`
	Cluster     Cluster   `gorm:"foreignKey:ClusterID" json:"cluster,omitempty"`
}

func (Pod) TableName() string { return "pods" }

// Node caches Kubernetes or host node info.
type Node struct {
	Base
	ClusterID   *uuid.UUID `gorm:"type:uuid;index" json:"cluster_id,omitempty"`
	Name        string     `gorm:"size:255;index;not null" json:"name"`
	Status      string     `gorm:"size:50;index" json:"status"`
	Roles       string     `gorm:"size:255" json:"roles,omitempty"`
	CPUCapacity string     `gorm:"size:50" json:"cpu_capacity,omitempty"`
	MemCapacity string     `gorm:"size:50" json:"mem_capacity,omitempty"`
	KubeletVersion string  `gorm:"size:50" json:"kubelet_version,omitempty"`
	Cluster     *Cluster   `gorm:"foreignKey:ClusterID" json:"cluster,omitempty"`
}

func (Node) TableName() string { return "nodes" }

// Server represents a monitored Linux/Windows host.
type Server struct {
	Base
	Hostname     string    `gorm:"size:255;not null;uniqueIndex" json:"hostname"`
	IPAddress    string    `gorm:"size:64;index" json:"ip_address"`
	OS           string    `gorm:"size:50;index" json:"os"` // linux | windows
	Status       string    `gorm:"size:50;default:online;index" json:"status"`
	CPUPercent   float64   `json:"cpu_percent"`
	MemoryPercent float64  `json:"memory_percent"`
	DiskPercent  float64   `json:"disk_percent"`
	LoadAvg1     float64   `json:"load_avg_1"`
	LoadAvg5     float64   `json:"load_avg_5"`
	LoadAvg15    float64   `json:"load_avg_15"`
	UptimeSeconds int64    `json:"uptime_seconds"`
	NetworkInBps  int64    `json:"network_in_bps"`
	NetworkOutBps int64    `json:"network_out_bps"`
	TemperatureC  *float64 `json:"temperature_c,omitempty"`
	LastSeenAt   time.Time `gorm:"index" json:"last_seen_at"`
	AgentVersion string    `gorm:"size:50" json:"agent_version,omitempty"`
	Tags         string    `gorm:"type:text" json:"tags,omitempty"`
}

func (Server) TableName() string { return "servers" }
