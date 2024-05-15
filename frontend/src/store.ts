import { ref } from "vue"

export interface Record {
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

const user = ref()
const records = ref<Record[]>([])

export { user, records }