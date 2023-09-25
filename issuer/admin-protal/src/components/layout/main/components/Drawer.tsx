import MuiDrawer from "@mui/material/Drawer"
import { styled } from "@mui/material/styles"
import { DRAWER_WIDTH } from ".."

const Drawer = styled(MuiDrawer)(() => ({
  width: DRAWER_WIDTH,
  flexShrink: 0,
  whiteSpace: "nowrap",
  boxSizing: "border-box",
  overflowX: "hidden",
  "& .MuiDrawer-paper": {
    width: DRAWER_WIDTH,
  },
}))

export default Drawer
