import { ref } from "vue"

export interface User {
  username: string
  token: string
}

export interface Netstat {
  id: number
  time: number
  localIP: string
  localPort: number
  remoteIP: string
  remotePort: number
  executable: string
  location: string
}

export interface Threat {
  id: number
  time: number
  ip: string
  tic: string
  reason: string
  risk: RiskLevel
  credibility: CredibilityLevel
}

export enum RiskLevel {
  Safe, Normal, Suspicious, Malicious, Unknown
}

export enum CredibilityLevel {
  Low, Medium, High,
}

export const riskLevel = ['安全', '正常', '可疑', '恶意', '未知']
export const credibilityLevel = ['低', '中', '高']

export const user = ref<User>()
export const netstats = ref<Netstat[]>([])
export const threats = ref<Threat[]>([])