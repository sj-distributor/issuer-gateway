import { FC } from "react"
import LoadingButton from "@mui/lab/LoadingButton/LoadingButton"
import PlaylistAddSharpIcon from "@mui/icons-material/PlaylistAddSharp"
import FileUpload from "@mui/icons-material/FileUpload"
import DeleteIcon from "@mui/icons-material/Delete"
import Refresh from "@mui/icons-material/Refresh"
import { Certs } from "@/entity/types"
import StyledTableCell from "./StyledTableCell"
import { OnDeleteParams, OPERATION_ROW_WIDTH } from "../hooks"

const OperationCell: FC<{
  cert: Certs
  onApplyCert: (id: number) => Promise<void>
  onDeleteItem: (params: OnDeleteParams) => void
}> = ({ cert, onApplyCert, onDeleteItem }) => {
  return (
    <StyledTableCell
      component="th"
      scope="row"
      align="center"
      width={OPERATION_ROW_WIDTH}
      sx={{
        display: "flex",
        gap: 2,
        justifyContent: "center",
        alignItems: "center",
      }}
    >
      {cert.expire <= 0 && (
        <LoadingButton
          size="small"
          onClick={() => onApplyCert(cert.id)}
          endIcon={<PlaylistAddSharpIcon />}
          loading={false}
          loadingPosition="end"
          variant="contained"
        >
          申请证书
        </LoadingButton>
      )}
      {cert.expire > 0 && (
        <LoadingButton
          size="small"
          onClick={() => {}}
          endIcon={<Refresh />}
          loading={false}
          loadingPosition="end"
          variant="contained"
        >
          重新申请
        </LoadingButton>
      )}
      <LoadingButton
        size="small"
        onClick={() => {}}
        endIcon={<FileUpload />}
        loading={false}
        loadingPosition="end"
        variant="contained"
      >
        上传证书
      </LoadingButton>
      <LoadingButton
        color="error"
        size="small"
        onClick={() =>
          onDeleteItem({
            id: cert.id,
            title: "",
            content: `确认删除 ${cert.domain}`,
          })
        }
        endIcon={<DeleteIcon />}
        loading={false}
        loadingPosition="end"
        variant="contained"
      >
        <p>删除</p>
      </LoadingButton>
    </StyledTableCell>
  )
}

export default OperationCell
