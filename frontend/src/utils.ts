import { ref } from "vue"

export interface IPRecord {
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
const records = ref<IPRecord[]>([])

export { user, records }

export const fetchGeoLocation = async (ip: any) => {
  try {
    const response = await fetch(`http://ip-api.com/json/${ip}`)
    const data = await response.json()
    return { lat: data.lat, lon: data.lon }
  } catch (error) {
    console.error('Error fetching geolocation:', error)
    return { lat: 0, lon: 0 } // Return default coordinates in case of error
  }
}