/**
 * Shared TypeScript types for homectl
 * These types should match the Go backend structures
 */

export interface Config {
  version: number;
  settings: Settings;
  groups: Group[];
  icons?: IconsConfig;
}

export interface Settings {
  title?: string;
  theme?: 'dark' | 'light';
  allowHosts?: string[];
  blockPrivateMetaIPs?: boolean;
  requestTimeout?: string;
  cache?: CacheConfig;
  auth?: AuthConfig;
  docker?: DockerConfig;
}

export interface CacheConfig {
  defaultTTL?: number;
  maxEntries?: number;
  widgetTTL?: Record<string, number>;
}

export interface AuthConfig {
  enabled?: boolean;
  provider?: 'local' | 'github' | 'google';
  session?: {
    maxAge?: string;
  };
  github?: {
    clientId?: string;
    clientSecret?: string;
    allowedUsers?: string[];
  };
}

export interface DockerConfig {
  enabled?: boolean;
  socket?: string;
  labelPrefix?: string;
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
  checks?: Check[];
  widgets?: Widget[];
  newTab?: boolean;
  pingEnabled?: boolean;
}

export type Check = HTTPCheck | TCPCheck | PingCheck;

export interface HTTPCheck {
  type: 'http';
  url: string;
  method?: string;
  expectStatus?: number;
  expectBodyContains?: string;
  headers?: Record<string, string>;
  timeout?: string;
  intervalSeconds?: number;
  retries?: number;
}

export interface TCPCheck {
  type: 'tcp';
  host: string;
  port: number;
  timeout?: string;
  intervalSeconds?: number;
}

export interface PingCheck {
  type: 'ping';
  host: string;
  count?: number;
  intervalSeconds?: number;
}

export type Widget = HTTPJSONWidget | HTTPHTMLWidget | TCPPortWidget;

export interface HTTPJSONWidget {
  type: 'httpJson';
  url: string;
  jsonPath: string;
  label?: string;
  format?: 'status' | 'bytes' | 'duration' | 'raw' | 'percent';
  cacheTTL?: number;
}

export interface HTTPHTMLWidget {
  type: 'httpHtml';
  url: string;
  selector: string;
  attribute?: string;
  label?: string;
}

export interface TCPPortWidget {
  type: 'tcpPort';
  host: string;
  port: number;
  label?: string;
  showLatency?: boolean;
}

// API Response types
export interface ServiceWithStatus extends Service {
  status: ServiceStatus;
  widgets?: WidgetResult[];
}

export interface ServiceStatus {
  state: 'up' | 'down' | 'degraded' | 'unknown';
  latency?: number;
  lastCheck?: string;
  error?: string;
}

export interface WidgetResult {
  label?: string;
  value: string | number | boolean;
  formatted?: string;
  state?: 'good' | 'warning' | 'error';
  lastUpdated?: string;
  error?: string;
}

export interface APIHealthResponse {
  status: 'healthy' | 'unhealthy';
  version: string;
  uptime: number;
}

export interface APIConfigResponse {
  title: string;
  theme: 'dark' | 'light';
  groups: Group[];
  icons?: IconsConfig;
  // Note: sensitive data is stripped
}

export interface APIServicesResponse {
  groups: {
    name: string;
    layout?: string;
    collapsed?: boolean;
    services: ServiceWithStatus[];
  }[];
}
