import { ChangeEvent, useEffect, useRef, useState } from "react"
import { SelectChangeEvent } from "@mui/material/Select/SelectInput"
import { SxProps, Theme } from "@mui/material/styles"
import { TableCellProps } from "@mui/material/TableCell/TableCell"
import { addDomainRequest, applyCertRequest, getCertsListRequest } from "@/api"
import { Certs } from "@/entity/types"

export interface OnDeleteParams {
  id: number
  title: string
  content: string
}

export interface ValidStatus {
  domainError: boolean
  emailError: boolean
  targetError: boolean
}

export enum ProtocolType {
  http = "http://",
  https = "https://",
}

export const PAGE_SIZE = 10
export const DOMAIN_ROW_MAX_WIDTH = 500
export const OPERATION_ROW_WIDTH = 500
export const tableCellConfig: {
  name: string
  alignPosition: TableCellProps["align"]
  styles?: SxProps<Theme>
}[] = [
  {
    name: "Domain",
    alignPosition: "left",
    styles: { maxWidth: DOMAIN_ROW_MAX_WIDTH },
  },
  {
    name: "Server",
    alignPosition: "center",
  },
  {
    name: "Status",
    alignPosition: "center",
  },
  {
    name: "Create Time",
    alignPosition: "center",
  },
  {
    name: "Operation",
    alignPosition: "center",
    styles: { width: OPERATION_ROW_WIDTH },
  },
]

export const useAction = () => {
  const currentPage = useRef<number>(1)
  const [finishInit, setFinishInit] = useState<boolean>(false)
  const [certsData, setCertsData] = useState<{
    certsList: Certs[]
    total: number
  }>({
    certsList: [],
    total: 0,
  })

  const [openDeleteDialog, setOpenDeleteDialog] = useState<boolean>(false)
  const deleteDialogData = useRef<{
    id: number
    title: string
    content: string
  }>({
    id: 0,
    title: "",
    content: "",
  })

  // TODO: 等待接口
  const onDeleteItem = (params: OnDeleteParams) => {
    deleteDialogData.current = {
      id: params.id,
      title: params.title,
      content: params.content,
    }
    setOpenDeleteDialog(true)
  }

  const onConfirmDelete = () => {
    setOpenDeleteDialog(false)
  }

  const onCancelDelete = () => {
    setOpenDeleteDialog(false)
  }

  const onPageChange = (_: ChangeEvent<unknown> | null, value: number) => {
    // <TablePagination /> page第一页下标是0
    currentPage.current = value + 1
    getCertsList()
  }

  const getCertsList = async () => {
    const { success, data, msg } = await getCertsListRequest({
      page: currentPage.current,
      size: PAGE_SIZE,
    })
    if (success && data) {
      setCertsData({
        certsList: data?.certs,
        total: data?.total,
      })
    } else {
      msg &&
        global.$toast.onOpen({
          text: msg,
          type: "error",
        })
    }
  }

  // TODO: 待对接
  const onApplyCert = async (id: number) => {
    const res = await applyCertRequest(id)
    console.log({ res })
  }

  const init = async () => {
    await getCertsList()
    setFinishInit(true)
  }

  useEffect(() => {
    init()
  }, [])

  return {
    currentPage,
    certsData,
    finishInit,
    openDeleteDialog,
    deleteDialogData,
    getCertsList,
    onDeleteItem,
    onConfirmDelete,
    onCancelDelete,
    onPageChange,
    onApplyCert,
  }
}

export const useBindDomain = ({
  getCertsList,
}: {
  getCertsList: () => Promise<void>
}) => {
  const bindingData = useRef<{ domain: string; email: string; target: string }>(
    {
      domain: "",
      email: "",
      target: "",
    }
  )
  const [targetCurrentProtocol, setTargetCurrentProtocol] =
    useState<ProtocolType>(ProtocolType.https)
  const [openBindingDialog, setOpenBindingDialog] = useState<boolean>(false)
  const [validStatus, setValidStatus] = useState<ValidStatus>({
    domainError: false,
    emailError: false,
    targetError: false,
  })

  const onOpenBindingDialog = () => {
    setOpenBindingDialog(true)
  }

  const handleEmail = (e: ChangeEvent<HTMLInputElement>) => {
    bindingData.current.email = e.target.value
  }

  const handleDomain = (e: ChangeEvent<HTMLInputElement>) => {
    bindingData.current.domain = e.target.value
  }

  const handleTarget = (e: ChangeEvent<HTMLInputElement>) => {
    bindingData.current.target = e.target.value
  }

  const onChangeProtocol = (e: SelectChangeEvent<ProtocolType>) => {
    setTargetCurrentProtocol(e.target.value as ProtocolType)
  }

  const onSubmit = async () => {
    const domainRegex = /^(?:[a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$/
    const emailRegex = /^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$/
    const targetRegex = /^www\.[a-zA-Z0-9-]+\.[a-zA-Z]{2,}$/
    const domainPass = domainRegex.test(bindingData.current.domain)
    const targetPass = targetRegex.test(bindingData.current.target)
    const emailPass = emailRegex.test(bindingData.current.email)
    setValidStatus({
      domainError: !domainPass,
      emailError: !emailPass,
      targetError: !targetPass,
    })
    if (domainPass && targetPass && emailPass) {
      const { success, msg } = await addDomainRequest({
        domain: bindingData.current.domain,
        target: `${targetCurrentProtocol}${bindingData.current.target}`,
        email: bindingData.current.email,
      })
      if (success) {
        getCertsList()
        setOpenBindingDialog(false)
      } else {
        msg &&
          global.$toast.onOpen({
            text: msg,
            type: "error",
          })
      }
    }
  }

  const onClose = () => {
    setOpenBindingDialog(false)
    setValidStatus({
      domainError: false,
      emailError: false,
      targetError: false,
    })
  }

  return {
    openBindingDialog,
    validStatus,
    targetCurrentProtocol,
    onOpenBindingDialog,
    handleEmail,
    handleDomain,
    handleTarget,
    onChangeProtocol,
    onSubmit,
    onClose,
  }
}
