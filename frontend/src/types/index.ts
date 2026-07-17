export type Role = 'admin' | 'devops' | 'developer' | 'viewer';

export interface User {
  id: string;
  email: string;
  name: string;
  role: Role;
  avatar_url?: string;
  is_active: boolean;
  organization_id?: string;
}

export interface AuthResponse {
  user: User;
  access_token: string;
  refresh_token: string;
  expires_at: string;
  token_type: string;
}

export interface DashboardStats {
  total_projects: number;
  running_builds: number;
  failed_builds: number;
  successful_builds: number;
  servers_online: number;
  docker_containers_running: number;
  pods_running: number;
  critical_alerts: number;
  deployments_today: number;
  open_incidents: number;
  cpu_usage: number;
  memory_usage: number;
  disk_usage: number;
  network_traffic_in: number;
  network_traffic_out: number;
}

export interface ApiResponse<T> {
  success: boolean;
  message?: string;
  data?: T;
  meta?: {
    page: number;
    page_size: number;
    total: number;
    total_pages: number;
  };
  error?: { code: string; details?: string };
}

export interface Alert {
  id: string;
  title: string;
  description?: string;
  severity: 'critical' | 'high' | 'medium' | 'low';
  status: string;
  source: string;
  created_at: string;
}

export interface Incident {
  id: string;
  title: string;
  description?: string;
  priority: string;
  status: string;
  root_cause?: string;
  resolution?: string;
  sla_deadline?: string;
  created_at: string;
}

export interface Deployment {
  id: string;
  application: string;
  environment: string;
  version: string;
  git_commit?: string;
  branch?: string;
  triggered_by?: string;
  status: string;
  rollback_version?: string;
  deployed_at: string;
}

export interface Project {
  id: string;
  name: string;
  slug: string;
  description?: string;
  environment: string;
  status: string;
  repository_url?: string;
}
