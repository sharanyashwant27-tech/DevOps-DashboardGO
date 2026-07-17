package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/devops-command-center/backend/config"
	"go.uber.org/zap"
)

type GitHubService struct {
	cfg    config.GitHubConfig
	client *http.Client
	log    *zap.Logger
}

func NewGitHubService(cfg config.GitHubConfig, log *zap.Logger) *GitHubService {
	svc := &GitHubService{
		cfg:    cfg,
		client: &http.Client{Timeout: 30 * time.Second},
		log:    log,
	}
	switch {
	case cfg.Live():
		log.Info("github live mode (authenticated)", zap.String("user", cfg.Username))
	case cfg.PublicUser() != "":
		log.Info("github public mode (no token)", zap.String("user", cfg.PublicUser()))
	case cfg.DemoMode:
		log.Info("github running in demo mode")
	}
	return svc
}

func (s *GitHubService) Enabled() bool { return s.cfg.Enabled() }

func (s *GitHubService) Mode() string {
	if s.cfg.Live() {
		return "live"
	}
	if s.cfg.PublicUser() != "" {
		return "public"
	}
	if s.cfg.DemoMode {
		return "demo"
	}
	return "disabled"
}

func (s *GitHubService) useDemo() bool {
	return !s.cfg.Live() && s.cfg.PublicUser() == "" && s.cfg.DemoMode
}

func (s *GitHubService) do(ctx context.Context, path string) ([]byte, error) {
	base := strings.TrimRight(s.cfg.BaseURL, "/")
	if base == "" {
		base = "https://api.github.com"
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, base+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("User-Agent", "DevOps-Command-Center")
	if s.cfg.Live() {
		req.Header.Set("Authorization", "Bearer "+s.cfg.Token)
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("github api %s: %s", path, string(body))
	}
	return body, nil
}

func (s *GitHubService) ListRepos(ctx context.Context) ([]map[string]interface{}, error) {
	if s.useDemo() {
		return demoRepos(), nil
	}
	var path string
	if s.cfg.Live() {
		path = "/user/repos?per_page=100&sort=updated"
	} else {
		path = fmt.Sprintf("/users/%s/repos?per_page=100&sort=updated&type=all", s.cfg.PublicUser())
	}
	body, err := s.do(ctx, path)
	if err != nil {
		if s.cfg.DemoMode {
			s.log.Warn("github live/public call failed, falling back to demo", zap.Error(err))
			return demoRepos(), nil
		}
		return nil, err
	}
	var repos []map[string]interface{}
	return repos, json.Unmarshal(body, &repos)
}

func (s *GitHubService) ListBranches(ctx context.Context, owner, repo string) ([]map[string]interface{}, error) {
	if s.useDemo() {
		return demoBranches(owner, repo), nil
	}
	body, err := s.do(ctx, fmt.Sprintf("/repos/%s/%s/branches?per_page=100", owner, repo))
	if err != nil {
		if s.cfg.DemoMode {
			return demoBranches(owner, repo), nil
		}
		return nil, err
	}
	var items []map[string]interface{}
	return items, json.Unmarshal(body, &items)
}

func (s *GitHubService) ListCommits(ctx context.Context, owner, repo string) ([]map[string]interface{}, error) {
	if s.useDemo() {
		return demoCommits(owner, repo), nil
	}
	body, err := s.do(ctx, fmt.Sprintf("/repos/%s/%s/commits?per_page=30", owner, repo))
	if err != nil {
		if s.cfg.DemoMode {
			return demoCommits(owner, repo), nil
		}
		return nil, err
	}
	var items []map[string]interface{}
	return items, json.Unmarshal(body, &items)
}

func (s *GitHubService) ListPullRequests(ctx context.Context, owner, repo string) ([]map[string]interface{}, error) {
	if s.useDemo() {
		return demoPulls(owner, repo), nil
	}
	body, err := s.do(ctx, fmt.Sprintf("/repos/%s/%s/pulls?state=all&per_page=30", owner, repo))
	if err != nil {
		if s.cfg.DemoMode {
			return demoPulls(owner, repo), nil
		}
		return nil, err
	}
	var items []map[string]interface{}
	return items, json.Unmarshal(body, &items)
}

func (s *GitHubService) ListIssues(ctx context.Context, owner, repo string) ([]map[string]interface{}, error) {
	if s.useDemo() {
		return demoIssues(owner, repo), nil
	}
	body, err := s.do(ctx, fmt.Sprintf("/repos/%s/%s/issues?state=all&per_page=30", owner, repo))
	if err != nil {
		if s.cfg.DemoMode {
			return demoIssues(owner, repo), nil
		}
		return nil, err
	}
	var items []map[string]interface{}
	return items, json.Unmarshal(body, &items)
}

func (s *GitHubService) ListReleases(ctx context.Context, owner, repo string) ([]map[string]interface{}, error) {
	if s.useDemo() {
		return demoReleases(owner, repo), nil
	}
	body, err := s.do(ctx, fmt.Sprintf("/repos/%s/%s/releases?per_page=20", owner, repo))
	if err != nil {
		if s.cfg.DemoMode {
			return demoReleases(owner, repo), nil
		}
		return nil, err
	}
	var items []map[string]interface{}
	return items, json.Unmarshal(body, &items)
}

func (s *GitHubService) ListWorkflowRuns(ctx context.Context, owner, repo string) (map[string]interface{}, error) {
	if s.useDemo() {
		return demoWorkflows(owner, repo), nil
	}
	body, err := s.do(ctx, fmt.Sprintf("/repos/%s/%s/actions/runs?per_page=20", owner, repo))
	if err != nil {
		if s.cfg.DemoMode {
			return demoWorkflows(owner, repo), nil
		}
		return nil, err
	}
	var out map[string]interface{}
	return out, json.Unmarshal(body, &out)
}

func (s *GitHubService) ListContributors(ctx context.Context, owner, repo string) ([]map[string]interface{}, error) {
	if s.useDemo() {
		return demoContributors(owner, repo), nil
	}
	body, err := s.do(ctx, fmt.Sprintf("/repos/%s/%s/contributors?per_page=30", owner, repo))
	if err != nil {
		if s.cfg.DemoMode {
			return demoContributors(owner, repo), nil
		}
		return nil, err
	}
	var items []map[string]interface{}
	return items, json.Unmarshal(body, &items)
}

func (s *GitHubService) RepositoryHealth(ctx context.Context, owner, repo string) (map[string]interface{}, error) {
	if s.useDemo() {
		return demoHealth(owner, repo), nil
	}
	body, err := s.do(ctx, fmt.Sprintf("/repos/%s/%s", owner, repo))
	if err != nil {
		if s.cfg.DemoMode {
			return demoHealth(owner, repo), nil
		}
		return nil, err
	}
	var repoInfo map[string]interface{}
	if err := json.Unmarshal(body, &repoInfo); err != nil {
		return nil, err
	}
	prs, _ := s.ListPullRequests(ctx, owner, repo)
	issues, _ := s.ListIssues(ctx, owner, repo)
	runs, _ := s.ListWorkflowRuns(ctx, owner, repo)
	health := map[string]interface{}{
		"name":             repoInfo["full_name"],
		"open_issues":      repoInfo["open_issues_count"],
		"stars":            repoInfo["stargazers_count"],
		"forks":            repoInfo["forks_count"],
		"default_branch":   repoInfo["default_branch"],
		"pushed_at":        repoInfo["pushed_at"],
		"description":      repoInfo["description"],
		"html_url":         repoInfo["html_url"],
		"language":         repoInfo["language"],
		"pull_requests":    len(prs),
		"issues_listed":    len(issues),
		"workflow_payload": runs,
		"score":            "healthy",
		"mode":             s.Mode(),
		"profile":          "https://github.com/" + owner,
	}
	if open, ok := repoInfo["open_issues_count"].(float64); ok && open > 50 {
		health["score"] = "needs_attention"
	}
	return health, nil
}

// --- Demo data (fallback only) ---

func demoRepos() []map[string]interface{} {
	now := time.Now().UTC().Format(time.RFC3339)
	return []map[string]interface{}{
		{
			"full_name": "acme-corp/platform-api", "name": "platform-api", "private": true,
			"language": "Go", "stargazers_count": 42, "forks_count": 8, "open_issues_count": 3,
			"default_branch": "main", "pushed_at": now, "html_url": "https://github.com/acme-corp/platform-api",
			"description": "Core platform services API",
		},
	}
}

func demoBranches(owner, repo string) []map[string]interface{} {
	return []map[string]interface{}{
		{"name": "main", "commit": map[string]interface{}{"sha": "a1b2c3d4e5f67890"}},
		{"name": "develop", "commit": map[string]interface{}{"sha": "b2c3d4e5f6789012"}},
	}
}

func demoCommits(owner, repo string) []map[string]interface{} {
	now := time.Now().UTC()
	return []map[string]interface{}{
		{
			"sha": "a1b2c3d4e5f6789012345678abcdef01",
			"commit": map[string]interface{}{
				"message": "feat: update " + repo,
				"author":  map[string]interface{}{"name": "demo", "date": now.Add(-2 * time.Hour).Format(time.RFC3339)},
			},
		},
	}
}

func demoPulls(owner, repo string) []map[string]interface{} {
	return []map[string]interface{}{
		{"number": 1, "title": "Sample PR", "state": "open", "user": map[string]interface{}{"login": owner}},
	}
}

func demoIssues(owner, repo string) []map[string]interface{} {
	return []map[string]interface{}{
		{"number": 1, "title": "Sample issue", "state": "open", "user": map[string]interface{}{"login": owner}},
	}
}

func demoReleases(owner, repo string) []map[string]interface{} {
	return []map[string]interface{}{
		{"tag_name": "v1.0.0", "name": "v1.0.0", "draft": false, "prerelease": false},
	}
}

func demoWorkflows(owner, repo string) map[string]interface{} {
	return map[string]interface{}{
		"total_count":   0,
		"workflow_runs": []map[string]interface{}{},
	}
}

func demoContributors(owner, repo string) []map[string]interface{} {
	return []map[string]interface{}{
		{"login": owner, "contributions": 1},
	}
}

func demoHealth(owner, repo string) map[string]interface{} {
	return map[string]interface{}{
		"name": owner + "/" + repo, "open_issues": 0, "stars": 0, "forks": 0,
		"default_branch": "main", "score": "healthy", "mode": "demo",
	}
}
