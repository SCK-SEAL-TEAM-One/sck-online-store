import { mockUpdateCartResponse } from '@/mock'

// ------------------------------------------------

export type UpdateProductInCartServiceResponse = {
  status: string // deleted, updated
}

type UpdateProductInCartServiceRequest = {
  product_id: number
  quantity: number
}

const updateProductInCartService = async ({
  product_id,
  quantity
}: UpdateProductInCartServiceRequest): Promise<UpdateProductInCartServiceResponse> => {
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

  let result = mockUpdateCartResponse.body

  return result
}

export default updateProductInCartService
