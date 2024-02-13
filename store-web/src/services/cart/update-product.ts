import axiosShoppingMallApi from '@/utils/axios'
import { handleServiceError } from '@/utils/helper'

// ------------------------------------------------

export type UpdateProductInCartServiceResponse = {
  data?: {
    status: string // deleted, updated
  }
  message?: string
}

type UpdateProductInCartServiceRequest = {
  productId: number
  quantity: number
}

const updateProductInCartService = async ({
  productId,
  quantity
}: UpdateProductInCartServiceRequest): Promise<UpdateProductInCartServiceResponse> => {
  try {
    const { data } = await axiosShoppingMallApi.put(`/api/v1/updateCart`, {
      product_id: productId,
      quantity: quantity
    })
    return {
      data: data
    }
  } catch (error) {
    return handleServiceError(error)
  }
}

export default updateProductInCartService
