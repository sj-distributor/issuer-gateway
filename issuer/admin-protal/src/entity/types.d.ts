export interface ApiResponse<T> {
  code: number
  data: T
  msg: string
}

export interface ApiResult<T> {
  status?: number
  code?: number
  data?: T
  msg?: string
  success: boolean
}

export interface Certs {
  created_at: number
  domain: string
  email: string
  expire: number
  id: number
  target: string
}
