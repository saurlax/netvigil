import dayjs from "dayjs"
import { computed, ref } from "vue"

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
  Unknown, Safe, Normal, Suspicious, Malicious
}

export enum CredibilityLevel {
  Low, Medium, High,
}

export const riskLevel = ['未知', '安全', '正常', '可疑', '恶意']
export const credibilityLevel = ['低', '中', '高']

export const user = ref<User>()
export const netstats = ref<Netstat[]>([])
export const threats = ref<Threat[]>([])

export const threatsMap = computed(() => {
  const result: { [ip: string]: Threat } = {}
  for (const t of threats.value) {
    const old = result[t.ip]
    if (!old || t.credibility > old.credibility || t.risk > old.risk) {
      result[t.ip] = t
    }
  }
  return result
})

/**
 * column: date, safe, normal, suspicious, malicious, unknown
 * 
 * row: day (1-7), weekly (8)
 */
export const stats = computed(() => {
  const now = Date.now()
  const aday = 86400000
  const result = Array.from({ length: 9 }, () => Array(6).fill(0))
  result[0] = ['日期', ...riskLevel]
  result[8][0] = '近七日'
  for (let i = 1; i < 8; i++) {
    result[i][0] = dayjs(now - 7 * aday + i * aday).format('YYYY-MM-DD')
  }
  for (const n of netstats.value) {
    const day = Math.floor((n.time - now + 8 * aday) / aday)
    if (day < 1) continue
    const risk = threatsMap.value[n.remoteIP]?.risk ?? RiskLevel.Unknown
    result[day][risk + 1]++
    result[8][risk + 1]++
  }
  return result
})