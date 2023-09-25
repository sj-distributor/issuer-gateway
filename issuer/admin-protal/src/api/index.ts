import axios from "./axios"
import { getResponseError, tranResponse } from "./lib"
import { ApiResult, Certs } from "@/entity/types"

export const loginRequest = async ({
  name,
  pass,
}: {
  name: string
  pass: string
}): Promise<ApiResult<{ token: string }>> => {
  try {
    const result = await axios.request({
      url: "/api/user/login",
      method: "POST",
      data: {
        name,
        pass,
      },
    })
    return tranResponse(result)
  } catch (error) {
    return getResponseError(error)
  }
}

export const getCertsListRequest = async ({
  page,
  size,
}: {
  page: number
  size: number
}): Promise<ApiResult<{ certs: Certs[]; total: number }>> => {
  try {
    const result = await axios.request({
      url: "/api/certs",
      method: "GET",
      params: {
        page,
        size,
      },
    })
    return tranResponse(result)
  } catch (error) {
    return getResponseError(error)
  }
}

export const addDomainRequest = async ({
  domain,
  email,
  target,
}: {
  domain: string
  email: string
  target: string
}): Promise<ApiResult<void>> => {
  try {
    const result = await axios.request({
      url: "/api/domain",
      method: "POST",
      data: {
        domain,
        email,
        target,
      },
    })
    return tranResponse(result)
  } catch (error) {
    return getResponseError(error)
  }
}

// TODO: 待对接
export const applyCertRequest = async (
  id: number
): Promise<ApiResult<void>> => {
  try {
    const result = await axios.request({
      url: "/api/cert",
      method: "POST",
      data: {
        id,
      },
    })
    return tranResponse(result)
  } catch (error) {
    return getResponseError(error)
  }
}

// TODO: 待对接
export const uploadCertRequest = async ({
  id,
  certificate,
  privateKey,
  issuerCertificate,
}: {
  id: number
  certificate: string
  privateKey: string
  issuerCertificate: string
}): Promise<ApiResult<void>> => {
  try {
    const result = await axios.request({
      url: "/api/cert/upload",
      method: "POST",
      data: { id, certificate, privateKey, issuerCertificate },
    })
    return tranResponse(result)
  } catch (error) {
    return getResponseError(error)
  }
}

// TODO: 待对接
export const renewCertRequest = async (
  id: number
): Promise<ApiResult<void>> => {
  try {
    const result = await axios.request({
      url: "/api/cert/upload",
      method: "PUT",
      data: { id },
    })
    return tranResponse(result)
  } catch (error) {
    return getResponseError(error)
  }
}
