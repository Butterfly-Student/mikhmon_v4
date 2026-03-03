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
  packetsIn?: number;
  packetsOut?: number;
  limitUptime?: string;
  limitBytesIn?: number;
  limitBytesOut?: number;
  limitBytesTotal: number;
  comment?: string;
  disabled: boolean;
  email?: string;
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
  gencode?: string;
  nameLength: number;
  prefix?: string;
  characterSet: string;
  profile: string;
  timeLimit?: string;
  dataLimit?: string;
  comment?: string;
}

export interface VoucherBatchResult {
  count: number;
  comment: string;
  vouchers: Voucher[];
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

/** System health dari /system/health/print — voltage dan temperature sebagai string dari RouterOS */
export interface SystemHealth {
  voltage?: string;
  temperature?: string;
  fanSpeed?: string;
  fanSpeed2?: string;
  fanSpeed3?: string;
}

/** Info routerboard dari /system/routerboard/print */
export interface RouterBoardInfo {
  routerboard?: string;
  model?: string;
  serialNumber?: string;
  firmwareType?: string;
  factoryFirmware?: string;
  currentFirmware?: string;
  upgradeFirmware?: string;
}

/** Log entry dari /log/print */
export interface LogEntry {
  '.id'?: string;
  id?: string;
  time?: string;
  topics?: string;
  message?: string;
}

/** Log entry dengan sequence ID dari WebSocket */
export interface LogEntryWithSeq {
  seq: number;
  time_ms: number;
  entry: LogEntry;
}

/** Metadata untuk init message */
export interface LogMetaInit {
  count: number;
  topics: string;
  maxSize: number;
  routerID: number;
}

/** Metadata untuk update message */
export interface LogMetaUpdate {
  batchSize: number;
  totalSeq: number;
}

/** WebSocket message type: init */
export interface LogMessageInit {
  type: 'init';
  data: LogEntryWithSeq[];
  meta: LogMetaInit;
}

/** WebSocket message type: update */
export interface LogMessageUpdate {
  type: 'update';
  data: LogEntryWithSeq[];
  meta: LogMetaUpdate;
}

/** WebSocket message type: error */
export interface LogMessageError {
  type: 'error';
  message: string;
}

/** WebSocket message type: status */
export interface LogMessageStatus {
  type: 'status';
  status: string;
}

/** Union type untuk semua log message dari WebSocket */
export type LogMessage = LogMessageInit | LogMessageUpdate | LogMessageError | LogMessageStatus;

/** Network interface dari /interface/print */
export interface NetworkInterface {
  id?: string;
  name?: string;
  type?: string;
  mtu?: number;
  macAddress?: string;
  running?: boolean;
  disabled?: boolean;
  comment?: string;
  rxByte?: number;
  txByte?: number;
  rxPacket?: number;
  txPacket?: number;
  rxDrop?: number;
  txDrop?: number;
  rxError?: number;
  txError?: number;
}

export interface DashboardData {
  routerId?: number;
  routerName?: string;
  identity?: SystemIdentity;
  resource?: SystemResource;
  health?: SystemHealth;
  routerBoard?: RouterBoardInfo;
  stats?: {
    activeUsers: number;
    totalUsers: number;
  };
  interfaces?: NetworkInterface[];
  hotspotLogs?: LogEntry[];
  connectionError?: string;
  // Legacy flat fields (backward compat dengan dashboard page lama)
  activeUsers?: number;
  totalUsers?: number;
  monthlyIncome?: number;
  dailyIncome?: number;
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
  /** Voltage dalam bentuk string dari RouterOS (mis: "24.3") */
  voltage?: string;
  /** Temperature dalam bentuk string dari RouterOS (mis: "42") */
  temperature?: string;
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

// PPP Types
export interface PPPSecret {
  id: string;
  name: string;
  password?: string;
  profile?: string;
  service?: string;
  disabled: boolean;
  callerID?: string;
  localAddress?: string;
  remoteAddress?: string;
  routes?: string;
  comment?: string;
  limitBytesIn?: number;
  limitBytesOut?: number;
  lastLoggedOut?: string;
  lastCallerID?: string;
  lastDisconnectReason?: string;
}

export interface PPPProfile {
  id: string;
  name: string;
  localAddress?: string;
  remoteAddress?: string;
  dnsServer?: string;
  sessionTimeout?: string;
  idleTimeout?: string;
  onlyOne?: boolean;
  comment?: string;
  rateLimit?: string;
  parentQueue?: string;
  queueType?: string;
  useCompression?: boolean;
  useEncryption?: boolean;
}

export interface PPPActive {
  id: string;
  name: string;
  service?: string;
  callerID?: string;
  address?: string;
  uptime?: string;
  sessionID?: string;
  encoding?: string;
  bytesIn?: number;
  bytesOut?: number;
  packetsIn?: number;
  packetsOut?: number;
  limitBytesIn?: number;
  limitBytesOut?: number;
}

// PPPInactive is a PPP Secret that is not currently online
export type PPPInactive = PPPSecret;

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
