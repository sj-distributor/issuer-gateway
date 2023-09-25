import { FC } from "react"
import DialogTitle from "@mui/material/DialogTitle/DialogTitle"
import Dialog from "@mui/material/Dialog/Dialog"
import DialogContent from "@mui/material/DialogContent/DialogContent"
import DialogContentText from "@mui/material/DialogContentText/DialogContentText"
import TextField from "@mui/material/TextField/TextField"
import DialogActions from "@mui/material/DialogActions/DialogActions"
import Button from "@mui/material/Button/Button"
import useTheme from "@mui/material/styles/useTheme"

const DeleteDialog: FC<{
  open: boolean
  title: string
  content: string
  onConfirm: () => void
  onCancel: () => void
}> = ({ open, title, content, onConfirm, onCancel }) => {
  const theme = useTheme()
  // TODO: 等待接口中
  return (
    <Dialog open={open} scroll={"paper"} maxWidth="md" fullWidth>
      <DialogTitle>{title}</DialogTitle>
      <DialogContent>
        <DialogContentText mb={theme.spacing(2)}>{content}</DialogContentText>
        <TextField
          label="请输入域名"
          variant="outlined"
          fullWidth
          required
          onChange={() => {}}
        />
      </DialogContent>
      <DialogActions>
        <Button onClick={onConfirm}>确认</Button>
        <Button onClick={onCancel}>取消</Button>
      </DialogActions>
    </Dialog>
  )
}

export default DeleteDialog
