import { mockProductDetailResponse } from '@/mock'

// ------------------------------------------------

export type GetProductDetailServiceResponse = {
  id: number
  product_name: string
  product_price: number
  product_image: string
  stock: number
  product_brand: string
}

const getProductDetailService = async (
  id: string
): Promise<GetProductDetailServiceResponse> => {
  // const queryString =
  //     '?' +
  //     new URLSearchParams({
  //       q: keyword,
  //       offset: offset.toString(),
  //       limit: limit.toString()
  //     }).toString()

  //   const result = await axiosShoppingMallApi.get(
  //     `/api/v1/product${queryString}`
  //   )

  let result = mockProductDetailResponse(id).body

  return result
}

export default getProductDetailService
