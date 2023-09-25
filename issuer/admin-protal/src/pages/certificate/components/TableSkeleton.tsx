import { FC } from "react"
import TableBody from "@mui/material/TableBody/TableBody"
import Skeleton from "@mui/material/Skeleton/Skeleton"
import TableRow from "@mui/material/TableRow/TableRow"
import StyledTableCell from "./StyledTableCell"
import { PAGE_SIZE, tableCellConfig } from "../hooks"

const TableSkeleton: FC = () => {
  return (
    <TableBody>
      {[...new Array(PAGE_SIZE)].map((_, index) => (
        <TableRow key={index} sx={{ textAlign: "center" }}>
          <StyledTableCell align="center" colSpan={tableCellConfig.length}>
            <Skeleton animation="wave" height={30} />
          </StyledTableCell>
        </TableRow>
      ))}
    </TableBody>
  )
}

export default TableSkeleton
