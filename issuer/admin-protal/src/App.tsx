import { BrowserRouter } from "react-router-dom"
import { AuthProvider } from "./context/authProvider"
import RenderRoutes from "./router"
import GlobalToast from "./components/toast/index"
import { ResetCss } from "./styles/ResetCss"

function App() {
  return (
    <>
      <ResetCss />
      <GlobalToast />
      <BrowserRouter>
        <AuthProvider>
          <RenderRoutes />
        </AuthProvider>
      </BrowserRouter>
    </>
  )
}

export default App
