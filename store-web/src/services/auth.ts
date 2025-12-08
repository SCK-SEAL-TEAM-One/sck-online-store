import axios, { isAxiosError } from 'axios'

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
  error?: string
  message?: string
}

const authAxiosInstance = axios.create({
  baseURL: process.env.storeServiceURL || 'http://localhost:3000',
  headers: {
    'Accept-Language': 'en'
  },
  withCredentials: true
})

export const Login = async (payload: LoginPayload): Promise<LoginResponse> => {
  try {
    const { data } = await authAxiosInstance.post(`/api/v1/login`, payload)
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
      if (error.response) {
        return {
          status: error.response?.status,
          message: error.response?.data?.message
        }
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
    const response = await authAxiosInstance.get(`/api/v1/refreshToken`)
    const accessToken = response.data.access_token
    const responseData: LoginSuccessResponse = {
      accessToken: accessToken,
      message: response.data.message
    }
    return {
      data: responseData
    }
  } catch (error) {
    if (isAxiosError(error)) {
      if (error.response) {
        return {
          status: error.response?.status,
          message: error.response?.data?.message
        }
      }
    }
    return {
      status: 500,
      message: 'Unknown Error'
    }
  }
}
