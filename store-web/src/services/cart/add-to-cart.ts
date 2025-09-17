import axiosShoppingMallApi from '@/utils/axios'
import { handleServiceError } from '@/utils/helper'
import {
  ProductDetailInCart,
  ProductDetailInCartSummary
} from './get-product-list'

// ------------------------------------------------

export type AddToCartServiceResponse = {
  data?: {
    carts: ProductDetailInCart[] // added, updated
    summary: ProductDetailInCartSummary
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
    const { data } = await axiosShoppingMallApi.put(`/api/v1/addCart`, {
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

export default addToCartService
