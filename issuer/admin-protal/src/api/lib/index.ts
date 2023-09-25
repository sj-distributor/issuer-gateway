import request, { AxiosResponse } from "axios"
import { ResponseCode } from "@/entity/enum"
import { ApiResponse, ApiResult } from "@/entity/types"

export const tranResponse = (
  result: AxiosResponse<ApiResponse<unknown>, any>
): ApiResult<any> => {
  return {
    status: result?.status,
    code: result?.data.code,
    data: result?.data.data,
    msg: result?.data?.msg,
    success:
      result?.status === 200 && result.data?.code === ResponseCode.Success,
  }
}

export const getResponseError = (error: unknown) => {
  if (request.isAxiosError(error) && error.response?.data) {
    return {
      status: error.status,
      success: false,
    }
  } else {
    return {
      status: undefined,
      success: false,
    }
  }
}
