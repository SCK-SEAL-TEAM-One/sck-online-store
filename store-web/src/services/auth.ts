import axiosShoppingMallApi from '@/utils/axios'
import { handleServiceError } from '@/utils/helper'

export interface LoginPayload {
  username: string
  password: string
}

export interface LoginSuccessResponse {
  accessToken: string
  message: string
}

export type LoginResponse = {
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
    // #TODO: Adjust error handle logic
    return handleServiceError(error)
  }
}
