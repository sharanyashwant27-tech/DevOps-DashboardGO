package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/devops-command-center/backend/config"
	"go.uber.org/zap"
)

type JenkinsService struct {
	cfg    config.JenkinsConfig
	client *http.Client
	log    *zap.Logger

	mu      sync.RWMutex
	demo    *demoJenkins
}

func NewJenkinsService(cfg config.JenkinsConfig, log *zap.Logger) *JenkinsService {
	svc := &JenkinsService{
		cfg:    cfg,
		client: &http.Client{Timeout: 30 * time.Second},
		log:    log,
	}
	if !cfg.Live() && cfg.DemoMode {
		svc.demo = newDemoJenkins()
		log.Info("jenkins running in demo mode (set jenkins.url to connect a real server)")
	}
	return svc
}

func (s *JenkinsService) Enabled() bool { return s.cfg.Enabled() }

func (s *JenkinsService) Mode() string {
	if s.cfg.Live() {
		return "live"
	}
	if s.demo != nil {
		return "demo"
	}
	return "disabled"
}

type JenkinsJob struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	Color string `json:"color"`
	Class string `json:"_class"`
}

type JenkinsJobsResponse struct {
	Jobs []JenkinsJob `json:"jobs"`
}

type JenkinsBuild struct {
	Number    int    `json:"number"`
	URL       string `json:"url"`
	Result    string `json:"result"`
	Building  bool   `json:"building"`
	Duration  int64  `json:"duration"`
	Timestamp int64  `json:"timestamp"`
}

type JenkinsQueueItem struct {
	ID    int    `json:"id"`
	Stuck bool   `json:"stuck"`
	Task  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"task"`
}

func (s *JenkinsService) useDemo() bool {
	return !s.cfg.Live() && s.demo != nil
}

func (s *JenkinsService) do(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	if !s.cfg.Live() {
		return nil, fmt.Errorf("jenkins is not configured")
	}
	u := strings.TrimRight(s.cfg.URL, "/") + path
	req, err := http.NewRequestWithContext(ctx, method, u, body)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(s.cfg.Username, s.cfg.Token)
	req.Header.Set("Accept", "application/json")
	return s.client.Do(req)
}

func (s *JenkinsService) ListJobs(ctx context.Context, search string) ([]JenkinsJob, error) {
	if s.useDemo() {
		return s.demo.listJobs(search), nil
	}
	resp, err := s.do(ctx, http.MethodGet, "/api/json?tree=jobs[name,url,color,_class]", nil)
	if err != nil {
		if s.cfg.DemoMode {
			s.ensureDemo()
			s.log.Warn("jenkins live call failed, falling back to demo", zap.Error(err))
			return s.demo.listJobs(search), nil
		}
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		err := fmt.Errorf("jenkins list jobs: %s", string(b))
		if s.cfg.DemoMode {
			s.ensureDemo()
			return s.demo.listJobs(search), nil
		}
		return nil, err
	}
	var out JenkinsJobsResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	if search == "" {
		return out.Jobs, nil
	}
	q := strings.ToLower(search)
	filtered := make([]JenkinsJob, 0)
	for _, j := range out.Jobs {
		if strings.Contains(strings.ToLower(j.Name), q) {
			filtered = append(filtered, j)
		}
	}
	return filtered, nil
}

func (s *JenkinsService) GetJobBuilds(ctx context.Context, jobName string) ([]JenkinsBuild, error) {
	if s.useDemo() {
		return s.demo.getBuilds(jobName), nil
	}
	path := fmt.Sprintf("/job/%s/api/json?tree=builds[number,url,result,building,duration,timestamp]", url.PathEscape(jobName))
	resp, err := s.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		if s.cfg.DemoMode {
			s.ensureDemo()
			return s.demo.getBuilds(jobName), nil
		}
		return nil, err
	}
	defer resp.Body.Close()
	var payload struct {
		Builds []JenkinsBuild `json:"builds"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}
	return payload.Builds, nil
}

func (s *JenkinsService) GetQueue(ctx context.Context) ([]JenkinsQueueItem, error) {
	if s.useDemo() {
		return s.demo.getQueue(), nil
	}
	resp, err := s.do(ctx, http.MethodGet, "/queue/api/json", nil)
	if err != nil {
		if s.cfg.DemoMode {
			s.ensureDemo()
			return s.demo.getQueue(), nil
		}
		return nil, err
	}
	defer resp.Body.Close()
	var payload struct {
		Items []JenkinsQueueItem `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}
	return payload.Items, nil
}

func (s *JenkinsService) TriggerBuild(ctx context.Context, jobName string) error {
	if s.useDemo() {
		return s.demo.trigger(jobName)
	}
	path := fmt.Sprintf("/job/%s/build", url.PathEscape(jobName))
	resp, err := s.do(ctx, http.MethodPost, path, nil)
	if err != nil {
		if s.cfg.DemoMode {
			s.ensureDemo()
			return s.demo.trigger(jobName)
		}
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("trigger build failed: %s", string(b))
	}
	return nil
}

func (s *JenkinsService) StopBuild(ctx context.Context, jobName string, buildNumber int) error {
	if s.useDemo() {
		return s.demo.stop(jobName, buildNumber)
	}
	path := fmt.Sprintf("/job/%s/%d/stop", url.PathEscape(jobName), buildNumber)
	resp, err := s.do(ctx, http.MethodPost, path, nil)
	if err != nil {
		if s.cfg.DemoMode {
			s.ensureDemo()
			return s.demo.stop(jobName, buildNumber)
		}
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("stop build failed: %s", string(b))
	}
	return nil
}

func (s *JenkinsService) ConsoleLog(ctx context.Context, jobName string, buildNumber int) (string, error) {
	if s.useDemo() {
		return s.demo.console(jobName, buildNumber), nil
	}
	path := fmt.Sprintf("/job/%s/%d/consoleText", url.PathEscape(jobName), buildNumber)
	resp, err := s.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		if s.cfg.DemoMode {
			s.ensureDemo()
			return s.demo.console(jobName, buildNumber), nil
		}
		return "", err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (s *JenkinsService) Stats(ctx context.Context) (map[string]interface{}, error) {
	jobs, err := s.ListJobs(ctx, "")
	if err != nil {
		return nil, err
	}
	running, failed, success := 0, 0, 0
	for _, j := range jobs {
		switch {
		case strings.Contains(j.Color, "anime"):
			running++
		case strings.HasPrefix(j.Color, "red"):
			failed++
		case strings.HasPrefix(j.Color, "blue"):
			success++
		}
	}
	return map[string]interface{}{
		"total_jobs":      len(jobs),
		"running_jobs":    running,
		"failed_jobs":     failed,
		"successful_jobs": success,
		"mode":            s.Mode(),
	}, nil
}

func (s *JenkinsService) ensureDemo() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.demo == nil {
		s.demo = newDemoJenkins()
	}
}

// --- Demo store ---

type demoJenkins struct {
	mu     sync.RWMutex
	jobs   map[string]*demoJob
	queue  []JenkinsQueueItem
	nextQ  int
}

type demoJob struct {
	Job    JenkinsJob
	Builds []JenkinsBuild
	Logs   map[int]string
}

func newDemoJenkins() *demoJenkins {
	now := time.Now().UnixMilli()
	d := &demoJenkins{
		jobs:  map[string]*demoJob{},
		nextQ: 1,
	}
	seed := []struct {
		name  string
		color string
	}{
		{"platform-api-ci", "blue"},
		{"platform-api-deploy", "blue_anime"},
		{"frontend-build", "red"},
		{"infra-terraform", "blue"},
		{"nightly-integration", "aborted"},
		{"security-scan", "yellow"},
	}
	for i, s := range seed {
		builds := []JenkinsBuild{
			{Number: 100 + i*3, URL: "#", Result: "SUCCESS", Building: false, Duration: 90000 + int64(i*1000), Timestamp: now - 3600000},
			{Number: 101 + i*3, URL: "#", Result: "FAILURE", Building: false, Duration: 45000, Timestamp: now - 1800000},
			{Number: 102 + i*3, URL: "#", Result: "", Building: strings.Contains(s.color, "anime"), Duration: 12000, Timestamp: now - 60000},
		}
		logs := map[int]string{}
		for _, b := range builds {
			logs[b.Number] = fmt.Sprintf(
				"[%s] Started by user admin\nBuilding on executor #1\nCheckout SCM...\nRunning tests...\nBuild #%d result=%s\nFinished: %s\n",
				s.name, b.Number, b.Result, map[bool]string{true: "RUNNING", false: "DONE"}[b.Building],
			)
		}
		d.jobs[s.name] = &demoJob{
			Job: JenkinsJob{
				Name:  s.name,
				URL:   "http://demo-jenkins.local/job/" + s.name,
				Color: s.color,
				Class: "hudson.model.FreeStyleProject",
			},
			Builds: builds,
			Logs:   logs,
		}
	}
	item := JenkinsQueueItem{ID: 1, Stuck: false}
	item.Task.Name = "security-scan"
	item.Task.URL = "http://demo-jenkins.local/job/security-scan"
	d.queue = []JenkinsQueueItem{item}
	return d
}

func (d *demoJenkins) listJobs(search string) []JenkinsJob {
	d.mu.RLock()
	defer d.mu.RUnlock()
	out := make([]JenkinsJob, 0, len(d.jobs))
	q := strings.ToLower(search)
	for _, j := range d.jobs {
		if q == "" || strings.Contains(strings.ToLower(j.Job.Name), q) {
			out = append(out, j.Job)
		}
	}
	return out
}

func (d *demoJenkins) getBuilds(jobName string) []JenkinsBuild {
	d.mu.RLock()
	defer d.mu.RUnlock()
	j, ok := d.jobs[jobName]
	if !ok {
		return nil
	}
	cp := make([]JenkinsBuild, len(j.Builds))
	copy(cp, j.Builds)
	return cp
}

func (d *demoJenkins) getQueue() []JenkinsQueueItem {
	d.mu.RLock()
	defer d.mu.RUnlock()
	cp := make([]JenkinsQueueItem, len(d.queue))
	copy(cp, d.queue)
	return cp
}

func (d *demoJenkins) trigger(jobName string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	j, ok := d.jobs[jobName]
	if !ok {
		return fmt.Errorf("job not found: %s", jobName)
	}
	next := 1
	if len(j.Builds) > 0 {
		next = j.Builds[0].Number + 1
	}
	b := JenkinsBuild{
		Number: next, URL: "#", Result: "", Building: true,
		Duration: 0, Timestamp: time.Now().UnixMilli(),
	}
	j.Builds = append([]JenkinsBuild{b}, j.Builds...)
	j.Job.Color = "blue_anime"
	j.Logs[next] = fmt.Sprintf("[%s] Triggered build #%d\nCheckout...\nCompiling...\n", jobName, next)
	return nil
}

func (d *demoJenkins) stop(jobName string, buildNumber int) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	j, ok := d.jobs[jobName]
	if !ok {
		return fmt.Errorf("job not found: %s", jobName)
	}
	for i := range j.Builds {
		if j.Builds[i].Number == buildNumber {
			j.Builds[i].Building = false
			j.Builds[i].Result = "ABORTED"
			j.Job.Color = "aborted"
			j.Logs[buildNumber] += "\nBuild was aborted by user.\n"
			return nil
		}
	}
	return fmt.Errorf("build #%d not found", buildNumber)
}

func (d *demoJenkins) console(jobName string, buildNumber int) string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	j, ok := d.jobs[jobName]
	if !ok {
		return "job not found"
	}
	if log, ok := j.Logs[buildNumber]; ok {
		return log
	}
	return fmt.Sprintf("No console output for %s #%d", jobName, buildNumber)
}
