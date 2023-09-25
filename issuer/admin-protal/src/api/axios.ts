import axios from "axios"
import history from "@/utils/router/history"
import { storage, StorageKeys } from "@/utils/storage"

const instance = axios.create({
  baseURL: "http://192.167.167.167:9090",
  timeout: 10000,
  headers: {
    "Content-Type": "application/json; charset=utf-8",
  },
})

instance.interceptors.request.use(
  (config) => {
    const token = storage.get<string>(StorageKeys.TOKEN)
    config.headers!.Authorization = token ? `Bearer ${token}` : ""
    return config
  },
  (error) => Promise.reject(error)
)

instance.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error?.response?.status === 401) {
      storage.remove(StorageKeys.TOKEN)
      history.replace("/login")
    }
    global.$toast.onOpen({
      text: "网络错误",
      type: "error",
    })
    return Promise.reject(error)
  }
)

export default instance
