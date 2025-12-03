import axiosShoppingMallApi from '@/utils/axios'
import { isAxiosError } from 'axios'

export interface LoginPayload {
  username: string
  password: string
}

export interface LoginSuccessResponse {
  accessToken: string
  message: string
}

export type LoginResponse = {
  status?: number
  data?: LoginSuccessResponse
  message?: string
}

export const Login = async (payload: LoginPayload): Promise<LoginResponse> => {
  try {
    const { data } = await axiosShoppingMallApi.post(`/api/v1/login`, payload)
    const accessToken = data.access_token
    const responseData: LoginSuccessResponse = {
      accessToken: accessToken,
      message: data.message
    }
    return {
      data: responseData
    }
  } catch (error) {
    if (isAxiosError(error)) {
      return {
        status: error.status,
        message: error.message
      }
    }
    return {
      status: 500,
      message: 'Unknown Error'
    }
  }
}

export const RefreshToken = async (): Promise<LoginResponse> => {
  try {
    const { data } = await axiosShoppingMallApi.get(`/api/v1/refreshToken`)
    const accessToken = data.access_token
    const responseData: LoginSuccessResponse = {
      accessToken: accessToken,
      message: data.message
    }
    return {
      data: responseData
    }
  } catch (error) {
    if (isAxiosError(error)) {
      return {
        status: error.status,
        message: error.message
      }
    }

    return {
      status: 500,
      message: 'Unknown Error'
    }
  }
}
