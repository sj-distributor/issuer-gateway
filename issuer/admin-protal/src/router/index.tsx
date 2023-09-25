import { FC } from "react"
import { Navigate, Route, Routes } from "react-router-dom"
import Certificate from "@/pages/certificate"
import Profile from "@/pages/profile"
import Login from "@/pages/login"
import NotFound from "@/pages/not-found/NotFound"
import MainLayout from "@/components/layout/main"
import RequireAuth from "@/components/requreAuth.tsx"

const RenderRoutes: FC = () => {
  return (
    <Routes>
      <Route path="/" element={<Navigate to="/certificate" replace />} />
      <Route path="/login" element={<Login />} />

      <Route element={<RequireAuth />}>
        <Route element={<MainLayout />}>
          <Route path="/certificate" element={<Certificate />} />
          <Route path="/profile" element={<Profile />} />
        </Route>
      </Route>

      <Route path="*" element={<NotFound />} />
    </Routes>
  )
}

export default RenderRoutes
