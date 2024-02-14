import axiosShoppingMallApi from '@/utils/axios'
import { handleServiceError } from '@/utils/helper'

// ------------------------------------------------

export type ProductDetailInCart = {
  id: number
  user_id: number
  product_id: number
  quantity: number
  product_name: string
  product_price: number
  product_price_thb: number
  product_price_full_thb: number
  product_image: string
  stock: number
  product_brand: string
}

export type ProductDetailInCartSummary = {
  total_price: number
  total_price_thb: number
  total_price_full_thb: number
  receive_point: number
}

export type GetProductInCartServiceResponse = {
  data?: {
    carts: ProductDetailInCart[]
    summary: ProductDetailInCartSummary
  }
  message?: string
}

const GetProductInCartService = async (
  userId: number
): Promise<GetProductInCartServiceResponse> => {
  try {
    const { data } = await axiosShoppingMallApi.get(`/api/v1/cart`)
    return {
      data: data
    }
  } catch (error) {
    return handleServiceError(error)
  }
}

export default GetProductInCartService
