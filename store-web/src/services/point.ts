import axiosShoppingMallApi from '@/utils/axios'
import { handleServiceError } from '@/utils/helper'

// ------------------------------------------------

export type GetPointServiceResponse = {
  data?: {
    point: number
  }
  message?: string
}

const getPointService = async (): Promise<GetPointServiceResponse> => {
  try {
    const { data } = await axiosShoppingMallApi.put(`/api/v1/point`)
    return {
      data: data
    }
  } catch (error) {
    return handleServiceError(error)
  }
}

export default getPointService
