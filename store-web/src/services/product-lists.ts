import axiosShoppingMallApi from '@/utils/axios'
import { handleServiceError } from '@/utils/helper'

// ------------------------------------------------

type GetProductListServiceRequest = {
  keyword?: string
  offset?: number
  limit?: number
}

export type ProductDetailType = {
  id: number
  product_name: string
  product_price: number
  product_price_thb: number
  product_price_full_thb: number
  product_image: string
}

export type ProductListDataType = {
  total: number
  products: ProductDetailType[]
}

export type GetProductListServiceResponse = {
  data?: ProductListDataType
  message?: string
}

const getProductListService = async ({
  keyword = '',
  offset = 0,
  limit = 50
}: GetProductListServiceRequest): Promise<GetProductListServiceResponse> => {
  try {
    const queryString =
      '?' +
      new URLSearchParams({
        q: keyword,
        offset: offset.toString(),
        limit: limit.toString()
      }).toString()

    const { data } = await axiosShoppingMallApi.get(
      `/api/v1/product${queryString}`
    )
    return {
      data: data
    }
  } catch (error) {
    return handleServiceError(error)
  }
}

export default getProductListService
