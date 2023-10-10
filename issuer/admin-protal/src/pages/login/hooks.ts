import { ChangeEvent, useRef, useState } from "react"
import { useNavigate } from "react-router-dom"
import { loginRequest } from "@/api"
import { StorageKeys } from "@/utils/storage"
import { useAuth } from "@/hooks"
import { storage } from "../../utils/storage/index"

export const useAction = () => {
  const navigate = useNavigate()
  const { onSetAuth } = useAuth()
  const [showHelpText, setShowHelpText] = useState<boolean>(false)
  const loginInfo = useRef<{ account: string; password: string }>({
    account: "",
    password: "",
  })

  const onChangeAccount = (e: ChangeEvent<HTMLInputElement>) => {
    loginInfo.current.account = e.target.value
  }

  const onChangePassword = (e: ChangeEvent<HTMLInputElement>) => {
    loginInfo.current.password = e.target.value
  }

  const onLogin = async () => {
    if (!loginInfo.current.account || !loginInfo.current.password) {
      setShowHelpText(true)
      return
    }
    const { success, data, msg } = await loginRequest({
      name: loginInfo.current.account,
      pass: loginInfo.current.password,
    })
    if (success && data) {
      storage.set(StorageKeys.TOKEN, data.token)
      onSetAuth(true)
      navigate("/")
    } else {
      msg &&
        global.$toast.onOpen({
          text: msg,
          type: "error",
        })
    }
  }

  return { showHelpText, onLogin, onChangeAccount, onChangePassword }
}
