import { ChangeEvent, FC } from "react"
import TablePagination from "@mui/material/TablePagination/TablePagination"
import { PAGE_SIZE } from "../hooks"

const Pagination: FC<{
  page: number
  total: number
  onPageChange: (_: ChangeEvent<unknown> | null, value: number) => void
}> = ({ page, total, onPageChange }) => {
  return (
    <TablePagination
      sx={{ mt: "auto" }}
      component="div"
      rowsPerPageOptions={[PAGE_SIZE]}
      count={total}
      rowsPerPage={PAGE_SIZE}
      page={page - 1}
      onPageChange={onPageChange}
    />
  )
}

export default Pagination
