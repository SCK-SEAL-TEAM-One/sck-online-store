// ------------------------------------------------

import axiosShoppingMallApi from '@/utils/axios'
import { handleServiceError } from '@/utils/helper'

export type ProductDetailType = {
  id: number
  product_name: string
  product_price: number
  product_price_thb: number
  product_price_full_thb: number
  product_image: string
  stock: number
  product_brand: string
}

export type GetProductDetailServiceResponse = {
  data?: ProductDetailType
  message?: string
}

const getProductDetailService = async (
  id: string
): Promise<GetProductDetailServiceResponse> => {
  try {
    const { data } = await axiosShoppingMallApi.get(`/api/v1/product/${id}`)
    return {
      data: data
    }
  } catch (error) {
    return handleServiceError(error)
  }
}

export default getProductDetailService
