import { useContext } from "react"
import AuthContext, { IAuthContext } from "@/context/authProvider"

export const useAuth = () => useContext<IAuthContext>(AuthContext)
