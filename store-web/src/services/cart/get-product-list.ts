import { mockCartListResponse } from '@/mock'

// ------------------------------------------------

export type ProductDetailInCart = {
  id: number
  user_id: number
  product_id: number
  quantity: number
  product_name: string
  product_price: number
  product_image: string
  stock: number
  product_brand: string
}

export type GetProductInCartServiceResponse = ProductDetailInCart[]

const GetProductInCartService = async (
  userId: number
): Promise<GetProductInCartServiceResponse> => {
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

  let result = mockCartListResponse.body

  return result
}

export default GetProductInCartService
