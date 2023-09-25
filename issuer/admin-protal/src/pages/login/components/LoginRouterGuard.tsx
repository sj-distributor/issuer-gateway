import { FC, ReactNode } from "react"
import { Navigate } from "react-router-dom"
import { useAuth } from "@/hooks"

const LoginRouterGuard: FC<{ children: ReactNode }> = ({ children }) => {
  const { auth } = useAuth()

  return <>{auth ? <Navigate to="/" replace /> : children}</>
}

export default LoginRouterGuard
