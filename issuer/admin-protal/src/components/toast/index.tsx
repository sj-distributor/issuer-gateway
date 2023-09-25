import { FC, forwardRef, useImperativeHandle, useRef, useState } from "react"
import { SxProps, Theme } from "@mui/material/styles"
import Alert, { AlertColor } from "@mui/material/Alert/Alert"
import Snackbar from "@mui/material/Snackbar/Snackbar"

export interface ToastOpenProps {
  text: string
  time?: number
  options?: {
    vertical?: "top" | "bottom"
    horizontal?: "center" | "left" | "right"
    sx?: SxProps<Theme> | undefined
  }
  type?: AlertColor
}

export interface ToastHandle {
  onOpen: (params: ToastOpenProps) => void
  onClose: () => void
}

interface ToastConfig {
  message: string
  time: number
  vertical: "top" | "bottom"
  horizontal: "center" | "left" | "right"
  sx: SxProps<Theme> | undefined
  type: AlertColor | undefined
}

const Toast = forwardRef((_, ref) => {
  const [open, setOpen] = useState<boolean>(false)
  const toastConfig = useRef<ToastConfig>({
    message: "",
    time: 3000,
    vertical: "top",
    horizontal: "center",
    sx: null,
    type: undefined,
  })

  const onOpen = ({ text, time, options, type }: ToastOpenProps) => {
    const msgStr = String(text).trim()
    if (!msgStr) return
    toastConfig.current.message = msgStr
    toastConfig.current.type = type || undefined
    if (options) {
      toastConfig.current.vertical = options.vertical || "top"
      toastConfig.current.horizontal = options.horizontal ?? "center"
      toastConfig.current.sx = options.sx ?? null
      toastConfig.current.time = time ?? 3000
    }
    setOpen(true)
  }

  const onClose = () => {
    setOpen(false)
  }

  useImperativeHandle(ref, () => ({
    onOpen,
    onClose,
  }))

  return (
    <Snackbar
      open={open}
      onClose={onClose}
      anchorOrigin={{
        vertical: toastConfig?.current?.vertical,
        horizontal: toastConfig?.current?.horizontal,
      }}
      message={toastConfig?.current?.message}
      autoHideDuration={toastConfig?.current?.time}
      key={toastConfig?.current?.vertical + toastConfig?.current?.horizontal}
      sx={toastConfig.current.sx}
    >
      {toastConfig?.current?.type && (
        <Alert variant="filled" severity={toastConfig?.current?.type}>
          {toastConfig?.current?.message}
        </Alert>
      )}
    </Snackbar>
  )
})

const GlobalToast: FC = () => {
  return (
    <Toast
      ref={(toastRef: globalThis.ToastHandle) => {
        global.$toast = toastRef
      }}
    />
  )
}

export default GlobalToast
