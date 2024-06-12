import { ref } from "vue"

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

export interface User {
  username: string
  token: string
}

const user = ref<User>()
const netstats = ref<Netstat[]>([])

export { user, netstats }
