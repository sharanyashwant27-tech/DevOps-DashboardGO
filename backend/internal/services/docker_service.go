package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/devops-command-center/backend/config"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"go.uber.org/zap"
)

type DockerService struct {
	cli     *client.Client
	enabled bool
	log     *zap.Logger
}

func NewDockerService(cfg config.DockerConfig, log *zap.Logger) *DockerService {
	svc := &DockerService{enabled: cfg.Enabled, log: log}
	if !cfg.Enabled {
		return svc
	}
	opts := []client.Opt{client.FromEnv, client.WithAPIVersionNegotiation()}
	if cfg.Host != "" {
		opts = append(opts, client.WithHost(cfg.Host))
	}
	cli, err := client.NewClientWithOpts(opts...)
	if err != nil {
		log.Warn("docker client unavailable", zap.Error(err))
		svc.enabled = false
		return svc
	}
	svc.cli = cli
	return svc
}

func (s *DockerService) Enabled() bool { return s.enabled && s.cli != nil }

func (s *DockerService) ensure() error {
	if !s.Enabled() {
		return fmt.Errorf("docker is not available")
	}
	return nil
}

func (s *DockerService) ListContainers(ctx context.Context, search string) ([]map[string]interface{}, error) {
	if err := s.ensure(); err != nil {
		return nil, err
	}
	items, err := s.cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return nil, err
	}
	out := make([]map[string]interface{}, 0, len(items))
	q := strings.ToLower(search)
	for _, c := range items {
		name := ""
		if len(c.Names) > 0 {
			name = strings.TrimPrefix(c.Names[0], "/")
		}
		if q != "" && !strings.Contains(strings.ToLower(name), q) && !strings.Contains(strings.ToLower(c.Image), q) {
			continue
		}
		out = append(out, map[string]interface{}{
			"Id":      c.ID,
			"Names":   c.Names,
			"Image":   c.Image,
			"State":   c.State,
			"Status":  c.Status,
			"Ports":   c.Ports,
			"Created": c.Created,
		})
	}
	return out, nil
}

func (s *DockerService) ListImages(ctx context.Context) ([]image.Summary, error) {
	if err := s.ensure(); err != nil {
		return nil, err
	}
	return s.cli.ImageList(ctx, image.ListOptions{All: true})
}

func (s *DockerService) ListVolumes(ctx context.Context) ([]*volume.Volume, error) {
	if err := s.ensure(); err != nil {
		return nil, err
	}
	resp, err := s.cli.VolumeList(ctx, volume.ListOptions{Filters: filters.NewArgs()})
	if err != nil {
		return nil, err
	}
	return resp.Volumes, nil
}

func (s *DockerService) ListNetworks(ctx context.Context) ([]network.Summary, error) {
	if err := s.ensure(); err != nil {
		return nil, err
	}
	return s.cli.NetworkList(ctx, network.ListOptions{})
}

func (s *DockerService) ContainerStats(ctx context.Context, id string) (map[string]float64, error) {
	if err := s.ensure(); err != nil {
		return nil, err
	}
	stats, err := s.cli.ContainerStats(ctx, id, false)
	if err != nil {
		return nil, err
	}
	defer stats.Body.Close()
	var v types.StatsJSON
	if err := json.NewDecoder(stats.Body).Decode(&v); err != nil {
		return nil, err
	}
	cpu := calculateCPUPercent(&v)
	mem := 0.0
	if v.MemoryStats.Limit > 0 {
		mem = float64(v.MemoryStats.Usage) / float64(v.MemoryStats.Limit) * 100.0
	}
	return map[string]float64{"cpu_percent": cpu, "memory_percent": mem}, nil
}

func calculateCPUPercent(v *types.StatsJSON) float64 {
	cpuDelta := float64(v.CPUStats.CPUUsage.TotalUsage - v.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(v.CPUStats.SystemUsage - v.PreCPUStats.SystemUsage)
	if systemDelta > 0 && cpuDelta > 0 {
		online := float64(v.CPUStats.OnlineCPUs)
		if online == 0 {
			online = float64(len(v.CPUStats.CPUUsage.PercpuUsage))
		}
		return (cpuDelta / systemDelta) * online * 100.0
	}
	return 0
}

func (s *DockerService) Logs(ctx context.Context, id string, tail string) (string, error) {
	if err := s.ensure(); err != nil {
		return "", err
	}
	reader, err := s.cli.ContainerLogs(ctx, id, container.LogsOptions{
		ShowStdout: true, ShowStderr: true, Tail: tail, Timestamps: true,
	})
	if err != nil {
		return "", err
	}
	defer reader.Close()
	b, err := io.ReadAll(reader)
	return string(b), err
}

func (s *DockerService) Start(ctx context.Context, id string) error {
	if err := s.ensure(); err != nil {
		return err
	}
	return s.cli.ContainerStart(ctx, id, container.StartOptions{})
}

func (s *DockerService) Stop(ctx context.Context, id string) error {
	if err := s.ensure(); err != nil {
		return err
	}
	timeout := 10
	return s.cli.ContainerStop(ctx, id, container.StopOptions{Timeout: &timeout})
}

func (s *DockerService) Restart(ctx context.Context, id string) error {
	if err := s.ensure(); err != nil {
		return err
	}
	timeout := 10
	return s.cli.ContainerRestart(ctx, id, container.StopOptions{Timeout: &timeout})
}

func (s *DockerService) Delete(ctx context.Context, id string, force bool) error {
	if err := s.ensure(); err != nil {
		return err
	}
	return s.cli.ContainerRemove(ctx, id, container.RemoveOptions{Force: force})
}

func (s *DockerService) Close() error {
	if s.cli != nil {
		return s.cli.Close()
	}
	return nil
}

func (s *DockerService) Ping(ctx context.Context) error {
	if err := s.ensure(); err != nil {
		return err
	}
	_, err := s.cli.Ping(ctx)
	return err
}

func (s *DockerService) Health(ctx context.Context) map[string]interface{} {
	if !s.Enabled() {
		return map[string]interface{}{"enabled": false}
	}
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	err := s.Ping(ctx)
	msg := ""
	if err != nil {
		msg = err.Error()
	}
	return map[string]interface{}{"enabled": true, "reachable": err == nil, "error": msg}
}
