import { mockOrderUpdateStatusResponse } from '@/mock'

// ------------------------------------------------

type OrderUpdateStatusServiceRequest = {
  orderId: number
  otp: number
  otpRef: string
}
export type OrderUpdateStatusServiceResponse = {
  order_id: number
  payment_date: string
  shipping_method_id: number
  tracking_id: string
}

const orderUpdateStatusService = async ({
  orderId,
  otp,
  otpRef
}: OrderUpdateStatusServiceRequest): Promise<OrderUpdateStatusServiceResponse> => {
  //   const result = await axiosShoppingMallApi.put(
  //     `/api/v1/order`,
  //   )

  // Mock Response
  let result = mockOrderUpdateStatusResponse.body

  return result
}

export default orderUpdateStatusService
