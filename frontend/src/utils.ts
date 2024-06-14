import moment from "moment"
import { computed, ref } from "vue"

export interface User {
  username: string
  token: string
}

export interface Netstat {
  Time: number
  LocalIP: string
  LocalPort: number
  RemoteIP: string
  RemotePort: number
  Executable: string
  Location: string
  Reason: string
  Risk: number
  TIX: string
  Confidence: number
}

export interface Threat {
  // TODO
}

export const user = ref<User>()
export const netstats = ref<Netstat[]>([])
export const threats = ref<Threat[]>([])

/**
 * column: date, unknown, safe, normal, suspicious, malicious
 * 
 * row: day (1-7), weekly (8), all (9)
 */
export const stats = computed(() => {
  const riskCounts = Array.from({ length: 10 }, () => Array(9).fill(0))
  const now = Date.now()
  const aweek = 604800000
  const aday = 86400000

  riskCounts[0] = ['日期', '未知', '安全', '正常', '可疑', '恶意']
  for (const n of netstats.value) {
    riskCounts[9][n.Risk + 1]++
    if (n.Time > now - aweek) {
      riskCounts[8][n.Risk + 1]++
      const day = 7 - Math.floor((now - n.Time) / aday)
      riskCounts[day][n.Risk + 1]++
    }
  }
  for (let i = 1; i <= 7; i++) {
    riskCounts[i][0] = moment(now - (7 - i) * aday).format('YYYY/MM/DD')
  }
  riskCounts[8][0] = 'Weekly'
  riskCounts[9][0] = 'All'
  return riskCounts
})