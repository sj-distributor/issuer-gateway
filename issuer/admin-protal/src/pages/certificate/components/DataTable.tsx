import { FC } from "react"
import { Certs } from "@/entity/types"
import Paper from "@mui/material/Paper/Paper"
import Table from "@mui/material/Table/Table"
import TableBody from "@mui/material/TableBody/TableBody"
import TableContainer from "@mui/material/TableContainer/TableContainer"
import TableHead from "@mui/material/TableHead/TableHead"
import TableRow from "@mui/material/TableRow/TableRow"
import Tooltip from "@mui/material/Tooltip/Tooltip"
import { formatDateTime } from "@/utils/time"
import { DOMAIN_ROW_MAX_WIDTH, OnDeleteParams, tableCellConfig } from "../hooks"
import StyledTableCell from "./StyledTableCell"
import StyledTableRow from "./StyledTableRow"
import TableSkeleton from "./TableSkeleton"
import OperationCell from "./OperationCell"
import * as styles from "../styles"

const DataTableBody: FC<{
  finishInit: boolean
  certsList: Certs[]
  total: number
  onApplyCert: (id: number) => Promise<void>
  onDeleteItem: (params: OnDeleteParams) => void
}> = ({ finishInit, certsList, total, onApplyCert, onDeleteItem }) => {
  if (!finishInit) {
    return <TableSkeleton />
  }
  return (
    <TableBody>
      {total ? (
        certsList.map((cert) => (
          <StyledTableRow key={cert.id}>
            <StyledTableCell
              component="th"
              scope="row"
              sx={{ maxWidth: DOMAIN_ROW_MAX_WIDTH, overflow: "hidden" }}
            >
              <Tooltip title={cert.domain} placement="bottom-start">
                <p css={styles.domainCell}>{cert.domain}</p>
              </Tooltip>
            </StyledTableCell>
            <StyledTableCell component="th" scope="row" align="left">
              {cert.target}
            </StyledTableCell>
            <StyledTableCell component="th" scope="row" align="center">
              {cert.expire}
            </StyledTableCell>
            <StyledTableCell component="th" scope="row" align="center">
              {`${formatDateTime(cert.created_at * 1000, "yyyy-MM-dd")}`}
            </StyledTableCell>
            <OperationCell
              cert={cert}
              onApplyCert={onApplyCert}
              onDeleteItem={onDeleteItem}
            />
          </StyledTableRow>
        ))
      ) : (
        <TableRow>
          <StyledTableCell
            component="th"
            scope="row"
            align="center"
            colSpan={tableCellConfig.length}
          >
            No Data...
          </StyledTableCell>
        </TableRow>
      )}
    </TableBody>
  )
}

const DataTable: FC<{
  finishInit: boolean
  certsList: Certs[]
  total: number
  onApplyCert: (id: number) => Promise<void>
  onDeleteItem: (params: OnDeleteParams) => void
}> = ({ finishInit, certsList, total, onApplyCert, onDeleteItem }) => {
  return (
    <TableContainer component={Paper}>
      <Table sx={{ minWidth: 700 }}>
        <TableHead>
          <TableRow>
            {tableCellConfig.map((cell) => (
              <StyledTableCell
                key={cell.name}
                align={cell.alignPosition}
                sx={cell.styles}
              >
                {cell.name}
              </StyledTableCell>
            ))}
          </TableRow>
        </TableHead>
        <DataTableBody
          finishInit={finishInit}
          certsList={certsList}
          total={total}
          onApplyCert={onApplyCert}
          onDeleteItem={onDeleteItem}
        />
      </Table>
    </TableContainer>
  )
}

export default DataTable
