import axios from 'axios';
import type { ApiResponse, AuthResponse, DashboardStats, Alert, Incident, Deployment, Project } from '../types';

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || '',
  timeout: 30000,
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

api.interceptors.response.use(
  (res) => res,
  async (error) => {
    const original = error.config;
    if (error.response?.status === 401 && !original._retry) {
      original._retry = true;
      const refresh = localStorage.getItem('refresh_token');
      if (refresh) {
        try {
          const { data } = await axios.post<ApiResponse<AuthResponse>>('/api/v1/auth/refresh', {
            refresh_token: refresh,
          });
          if (data.data) {
            localStorage.setItem('access_token', data.data.access_token);
            localStorage.setItem('refresh_token', data.data.refresh_token);
            original.headers.Authorization = `Bearer ${data.data.access_token}`;
            return api(original);
          }
        } catch {
          localStorage.clear();
          window.location.href = '/login';
        }
      }
    }
    return Promise.reject(error);
  },
);

export const authApi = {
  login: (email: string, password: string) =>
    api.post<ApiResponse<AuthResponse>>('/api/v1/auth/login', { email, password }),
  register: (payload: { email: string; password: string; name: string }) =>
    api.post<ApiResponse<AuthResponse>>('/api/v1/auth/register', payload),
  forgotPassword: (email: string) =>
    api.post('/api/v1/auth/forgot-password', { email }),
  me: () => api.get<ApiResponse<import('../types').User>>('/api/v1/auth/me'),
};

export const dashboardApi = {
  stats: () => api.get<ApiResponse<DashboardStats>>('/api/v1/dashboard/stats'),
};

export const projectApi = {
  list: (search = '') => api.get<ApiResponse<Project[]>>('/api/v1/projects', { params: { search } }),
};

export const alertApi = {
  list: (params?: Record<string, string>) =>
    api.get<ApiResponse<Alert[]>>('/api/v1/alerts', { params }),
  acknowledge: (id: string) => api.post(`/api/v1/alerts/${id}/acknowledge`),
  resolve: (id: string) => api.post(`/api/v1/alerts/${id}/resolve`),
  mute: (id: string, minutes: number) => api.post(`/api/v1/alerts/${id}/mute`, { minutes }),
};

export const incidentApi = {
  list: (params?: Record<string, string>) =>
    api.get<ApiResponse<Incident[]>>('/api/v1/incidents', { params }),
  create: (payload: Partial<Incident> & { title: string; priority: string }) =>
    api.post('/api/v1/incidents', payload),
  update: (id: string, payload: Record<string, unknown>) =>
    api.patch(`/api/v1/incidents/${id}`, payload),
};

export const deploymentApi = {
  list: () => api.get<ApiResponse<Deployment[]>>('/api/v1/deployments'),
  rollback: (id: string) => api.post(`/api/v1/deployments/${id}/rollback`),
};

export const jenkinsApi = {
  jobs: (search = '') => api.get('/api/v1/jenkins/jobs', { params: { search } }),
  builds: (job: string) => api.get(`/api/v1/jenkins/jobs/${encodeURIComponent(job)}/builds`),
  queue: () => api.get('/api/v1/jenkins/queue'),
  stats: () => api.get('/api/v1/jenkins/stats'),
  trigger: (job: string) => api.post(`/api/v1/jenkins/jobs/${encodeURIComponent(job)}/build`),
  console: (job: string, number: number) =>
    api.get(`/api/v1/jenkins/jobs/${encodeURIComponent(job)}/builds/${number}/console`),
};

export const githubApi = {
  repos: () => api.get('/api/v1/github/repos'),
  health: (owner: string, repo: string) => api.get(`/api/v1/github/repos/${owner}/${repo}/health`),
  commits: (owner: string, repo: string) => api.get(`/api/v1/github/repos/${owner}/${repo}/commits`),
  pulls: (owner: string, repo: string) => api.get(`/api/v1/github/repos/${owner}/${repo}/pulls`),
  workflows: (owner: string, repo: string) =>
    api.get(`/api/v1/github/repos/${owner}/${repo}/actions/runs`),
};

export const dockerApi = {
  containers: (search = '') => api.get('/api/v1/docker/containers', { params: { search } }),
  images: () => api.get('/api/v1/docker/images'),
  volumes: () => api.get('/api/v1/docker/volumes'),
  networks: () => api.get('/api/v1/docker/networks'),
  logs: (id: string) => api.get(`/api/v1/docker/containers/${id}/logs`),
  start: (id: string) => api.post(`/api/v1/docker/containers/${id}/start`),
  stop: (id: string) => api.post(`/api/v1/docker/containers/${id}/stop`),
  restart: (id: string) => api.post(`/api/v1/docker/containers/${id}/restart`),
  remove: (id: string) => api.delete(`/api/v1/docker/containers/${id}`),
};

export const k8sApi = {
  namespaces: () => api.get('/api/v1/kubernetes/namespaces'),
  pods: (namespace = '') => api.get('/api/v1/kubernetes/pods', { params: { namespace } }),
  deployments: (namespace = 'default') =>
    api.get('/api/v1/kubernetes/deployments', { params: { namespace } }),
  nodes: () => api.get('/api/v1/kubernetes/nodes'),
  services: (namespace = 'default') =>
    api.get('/api/v1/kubernetes/services', { params: { namespace } }),
  events: (namespace = 'default') =>
    api.get('/api/v1/kubernetes/events', { params: { namespace } }),
  podLogs: (pod: string, namespace = 'default') =>
    api.get(`/api/v1/kubernetes/pods/${pod}/logs`, { params: { namespace } }),
  scale: (name: string, replicas: number, namespace = 'default') =>
    api.post(`/api/v1/kubernetes/deployments/${name}/scale`, { replicas }, { params: { namespace } }),
  restart: (name: string, namespace = 'default') =>
    api.post(`/api/v1/kubernetes/deployments/${name}/restart`, null, { params: { namespace } }),
};

export const serverApi = {
  list: () => api.get('/api/v1/servers'),
  local: () => api.get('/api/v1/servers/local'),
};

export const metricsApi = {
  series: (name: string, hours = 24) =>
    api.get(`/api/v1/metrics/${name}`, { params: { hours } }),
};

export default api;
