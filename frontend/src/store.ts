import { ref } from "vue"

export interface IPRecord {
  Time: number
  LocalAddr: string
  RemoteAddr: string
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
const records = ref<IPRecord[]>([])

export { user, records }