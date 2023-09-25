import { styled } from "@mui/material/styles"
import TableCell from "@mui/material/TableCell/TableCell"
import tableCellClasses from "@mui/material/TableCell/tableCellClasses"

const StyledTableCell = styled(TableCell)(({ theme }) => ({
  [`&.${tableCellClasses.head}`]: {
    backgroundColor: theme.palette.primary.main,
    color: theme.palette.common.white,
  },
  [`&.${tableCellClasses.body}`]: {
    fontSize: 14,
  },
}))

export default StyledTableCell
