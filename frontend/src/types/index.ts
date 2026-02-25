// Auth Types
export interface UserInfo {
  id: string;
  username: string;
  email?: string;
}

export interface LoginCredentials {
  username: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  expiresAt: string;
  user: UserInfo;
}

// Router Types
export interface Router {
  id: string | number;
  name: string;
  host: string;
  port: number;
  useSsl: boolean;
  timeout: number;
  isActive: boolean;
  description?: string;
  lastConnected?: string;
  createdAt: string;
  updatedAt: string;
}

// Hotspot Types
export interface HotspotUser {
  id: string;
  server?: string;
  name: string;
  password?: string;
  profile?: string | { name: string };
  macAddress?: string;
  ipAddress?: string;
  uptime?: string;
  bytesIn: number;
  bytesOut: number;
  limitUptime?: string;
  limitBytesTotal: number;
  comment?: string;
  disabled: boolean;
}

export interface UserProfile {
  id: string;
  name: string;
  addressPool?: string;
  sharedUsers: number;
  rateLimit?: string;
  parentQueue?: string;
  onLogin?: string;
  onLogout?: string;
  expireMode?: string;
  validity?: string;
  price: number;
  sellingPrice: number;
  lockUser?: string;
  lockServer?: string;
}

export interface HotspotActive {
  id: string;
  server?: string;
  user: string;
  address: string;
  macAddress: string;
  loginBy: string;
  uptime: string;
  sessionTimeLeft?: string;
  idleTime?: string;
  bytesIn: number;
  bytesOut: number;
  packetsIn: number;
  packetsOut: number;
  radius: boolean;
  comment?: string;
}

export interface HotspotHost {
  id: string;
  macAddress: string;
  address?: string;
  toAddress?: string;
  server?: string;
  authorized: boolean;
  bypassed: boolean;
  blocked: boolean;
  foundBy?: string;
  comment?: string;
}

// Voucher Types
export interface Voucher {
  id?: string;
  username: string;
  password?: string;
  profile: string | { name: string };
  server?: string;
  timeLimit?: string;
  dataLimit?: string;
  comment?: string;
  name?: string;
}

export interface GenerateVoucherRequest {
  quantity: number;
  server?: string;
  mode: 'vc' | 'up';
  nameLength: number;
  prefix?: string;
  characterSet: string;
  profile: string;
  timeLimit?: string;
  dataLimit?: string;
  comment?: string;
}

// Dashboard Types
export interface SystemIdentity {
  name: string;
}

export interface SystemResource {
  uptime?: string;
  version?: string;
  buildTime?: string;
  freeMemory?: number;
  totalMemory?: number;
  freeHddSpace?: number;
  totalHddSpace?: number;
  writeSectSinceReboot?: number;
  writeSectTotal?: number;
  badBlocks?: number;
  architectureName?: string;
  boardName?: string;
  cpu?: string;
  cpuCount?: number;
  cpuFrequency?: number;
  cpuLoad?: number;
}

export interface DashboardData {
  identity?: SystemIdentity;
  routerName?: string;
  resource?: SystemResource;
  stats?: {
    activeUsers: number;
    totalUsers: number;
  };
  activeUsers: number;
  totalUsers: number;
  monthlyIncome: number;
  dailyIncome: number;
}

export interface SystemResources {
  cpuLoad?: number;
  memoryUsed?: number;
  freeMemory?: number;
  memoryTotal?: number;
  totalMemory?: number;
  hddUsed?: number;
  freeHddSpace?: number;
  hddTotal?: number;
  totalHddSpace?: number;
  voltage?: number | string;
  temperature?: number | string;
  uptime?: string;
  version?: string;
  boardName?: string;
}

export interface SystemInfo {
  uptime: string;
  boardName: string;
  model?: string;
  version: string;
}

// Report Types
export interface SalesReport {
  date: string;
  time: string;
  username: string;
  price: number;
  ipAddress: string;
  macAddress: string;
  validity: string;
  profile: string;
  comment?: string;
}

// API Response
export interface ApiResponse<T> {
  success: boolean;
  message?: string;
  data?: T;
  error?: string;
}

// Filter Types
export interface UserFilter {
  profile?: string;
  comment?: string;
  search?: string;
}

export interface ReportFilter {
  date?: string;
  month?: string;
  year?: string;
}
