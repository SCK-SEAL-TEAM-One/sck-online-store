import axiosShoppingMallApi from '@/utils/axios'
import { handleServiceError } from '@/utils/helper'

// ------------------------------------------------

export type AddToCartServiceResponse = {
  data?: {
    status: string // added, updated
  }
  message?: string
}

type AddToCartServiceRequest = {
  productId: number
  quantity: number
}

const addToCartService = async ({
  productId,
  quantity
}: AddToCartServiceRequest): Promise<AddToCartServiceResponse> => {
  try {
    const { data } = await axiosShoppingMallApi.put(
      `/api/v1/addCart`,
      {
        product_id: productId,
        quantity: quantity
      }
    )
    return {
      data: data
    }
  } catch (error) {
    return handleServiceError(error)
  }
}

export default addToCartService
