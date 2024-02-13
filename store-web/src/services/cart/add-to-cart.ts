import { mockAddToCartResponse } from '@/mock'

// ------------------------------------------------

export type AddToCartServiceResponse = {
  status: string // added, updated
}

type AddToCartServiceRequest = {
  product_id: number
  quantity: number
}

const addToCartService = async ({
  product_id,
  quantity
}: AddToCartServiceRequest): Promise<AddToCartServiceResponse> => {
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

  let result = mockAddToCartResponse.body

  return result
}

export default addToCartService
