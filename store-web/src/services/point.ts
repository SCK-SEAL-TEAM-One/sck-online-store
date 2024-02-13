import { mockGetPointResponse } from '@/mock'

// ------------------------------------------------

export type GetPointServiceResponse = {
  point: number
}

const getPointService = async (): Promise<GetPointServiceResponse> => {
  //   const result = await axiosShoppingMallApi.get(
  //     `/api/v1/point`
  //   )

  // Mock Response
  let result = mockGetPointResponse.body

  return result
}

export default getPointService
