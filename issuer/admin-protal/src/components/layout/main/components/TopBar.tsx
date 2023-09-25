import { FC, memo, useState, MouseEvent } from "react"
import { useNavigate } from "react-router-dom"
import MuiAppBar, { AppBarProps } from "@mui/material/AppBar"
import Toolbar from "@mui/material/Toolbar"
import AccountCircleIcon from "@mui/icons-material/AccountCircle"
import IconButton from "@mui/material/IconButton"
import Typography from "@mui/material/Typography"
import Avatar from "@mui/material/Avatar/Avatar"
import MenuItem from "@mui/material/MenuItem/MenuItem"
import Menu from "@mui/material/Menu/Menu"
import { styled, useTheme } from "@mui/material/styles"
import { storage, StorageKeys } from "@/utils/storage"
import { useAuth } from "@/hooks"

const TopBarContainer = styled(MuiAppBar)<AppBarProps>(({ theme }) => ({
  zIndex: theme.zIndex.drawer + 1,
  transition: theme.transitions.create(["width", "margin"], {
    easing: theme.transitions.easing.sharp,
    duration: theme.transitions.duration.leavingScreen,
  }),
}))

const UserAvatar: FC = () => {
  const theme = useTheme()
  const { onSetAuth } = useAuth()
  const navigate = useNavigate()
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null)
  const onOPenSetting = (e: MouseEvent<HTMLElement>) =>
    setAnchorEl(e.currentTarget)
  const onCloseSetting = () => setAnchorEl(null)
  const onLogout = () => {
    storage.remove(StorageKeys.TOKEN)
    onSetAuth(false)
    navigate("/login")
  }

  return (
    <>
      <IconButton onClick={onOPenSetting} sx={{ p: 0 }}>
        <Avatar sx={{ bgcolor: theme.palette.primary.dark }}>
          <AccountCircleIcon />
        </Avatar>
      </IconButton>
      <Menu
        sx={{ mt: "45px" }}
        anchorEl={anchorEl}
        anchorOrigin={{
          vertical: "top",
          horizontal: "right",
        }}
        keepMounted
        transformOrigin={{
          vertical: "top",
          horizontal: "right",
        }}
        open={Boolean(anchorEl)}
        onClose={onCloseSetting}
      >
        <MenuItem onClick={onLogout}>
          <Typography textAlign="center">Log out</Typography>
        </MenuItem>
      </Menu>
    </>
  )
}

const TopBar: FC<{}> = memo(() => {
  return (
    <TopBarContainer position="fixed">
      <Toolbar>
        <Typography variant="h6" noWrap component="p" flex={1}>
          Issuer Gateway Admin Protal
        </Typography>
        <UserAvatar />
      </Toolbar>
    </TopBarContainer>
  )
})

export default TopBar
