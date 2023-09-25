import { FC, ReactNode, createContext, useState } from "react"
import { storage, StorageKeys } from "@/utils/storage"

export interface IAuthContext {
  auth: boolean
  onSetAuth: (auth: boolean) => void
}

const token = storage.get<string>(StorageKeys.TOKEN)
const AuthContext = createContext<IAuthContext>({
  auth: !!token,
} as IAuthContext)

export const AuthProvider: FC<{ children: ReactNode }> = ({ children }) => {
  const [auth, setAuth] = useState<boolean>(!!token)
  const onSetAuth = (auth: boolean) => {
    setAuth(auth)
  }

  return (
    <AuthContext.Provider value={{ auth, onSetAuth }}>
      {children}
    </AuthContext.Provider>
  )
}

export default AuthContext
