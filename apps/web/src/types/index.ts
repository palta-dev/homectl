export interface HealthResponse {
  status: 'healthy' | 'unhealthy';
  version: string;
  uptime: number;
}

export interface ConfigResponse {
  title: string;
  theme: 'dark' | 'light';
  background?: string;
  backgroundOpacity?: number;
  allowHosts?: string[];
  requestTimeout?: string;
  password?: string;
  docker?: DockerConfig;
  groups: Group[];
  icons?: IconsConfig;
  passwordProtected: boolean;
}

export interface DockerConfig {
  enabled: boolean;
  socket?: string;
  labelPrefix?: string;
  hosts?: HostConfig[];
  subnet?: string;
  ignore?: string[];
}

export interface HostConfig {
  name?: string;
  address: string;
  ports?: number[];
  tags?: string[];
}

export interface IconsConfig {
  sources: IconSource[];
}

export type IconSource =
  | { type: 'local'; path: string }
  | { type: 'simpleicons'; cache?: boolean; cacheTTL?: number }
  | { type: 'url'; baseUrl: string; pathTemplate?: string };

export interface Group {
  name: string;
  layout?: 'grid' | 'list' | 'compact';
  collapsed?: boolean;
  services: Service[];
}

export interface Service {
  name: string;
  url: string;
  icon?: string;
  description?: string;
  tags?: string[];
  newTab?: boolean;
  pingEnabled?: boolean;
  widgets?: Widget[];
}

export interface Widget {
  type: 'httpJson' | 'httpHtml' | 'tcpPort' | 'httpStatus' | 'system';
  url?: string;
  host?: string;
  port?: number;
  jsonPath?: string;
  selector?: string;
  attribute?: string;
  label?: string;
  format?: 'status' | 'bytes' | 'duration' | 'raw' | 'percent';
  cacheTTL?: number;
  options?: Record<string, string>;
}

export interface ServicesResponse {
  groups: GroupWithStatus[];
}

export interface GroupWithStatus {
  name: string;
  layout?: string;
  collapsed?: boolean;
  services: ServiceWithStatus[];
}

export interface ServiceWithStatus {
  name: string;
  url: string;
  icon?: string;
  favicon?: string;
  description?: string;
  tags?: string[];
  newTab?: boolean;
  pingEnabled?: boolean;
  status: ServiceStatus;
  widgets?: WidgetResult[];
  isDiscovered?: boolean;
}

export interface ServiceStatus {
  state: 'up' | 'down' | 'degraded' | 'unknown';
  latency?: number;
  lastCheck?: string;
  error?: string;
}

export interface WidgetResult {
  label?: string;
  value: string | number | boolean | null;
  formatted?: string;
  state?: 'good' | 'warning' | 'error';
  lastUpdated?: string;
  error?: string;
}
