import { mockOrderUpdateStatusResponse } from '@/mock'
import axiosShoppingMallApi from '@/utils/axios'
import { handleServiceError } from '@/utils/helper'

// ------------------------------------------------

type OrderUpdateStatusServiceRequest = {
  orderId: number
  otp: number
  otpRef: string
}

export type OrderUpdateStatusDetailType = {
  order_id: number
  payment_date: string
  shipping_method_id: number
  tracking_id: string
}

export type OrderUpdateStatusServiceResponse = {
  data?: OrderUpdateStatusDetailType
  message?: string
}

const orderUpdateStatusService = async ({
  orderId,
  otp,
  otpRef
}: OrderUpdateStatusServiceRequest): Promise<OrderUpdateStatusServiceResponse> => {
  // try {
  //   const { data } = await axiosShoppingMallApi.put(`/api/v1/order`, {
  //     order_id: orderId,
  //     otp: otp,
  //     ref_otp: otpRef
  //   })
  //   return {
  //     data: data
  //   }
  // } catch (error) {
  //   return handleServiceError(error)
  // }

  // Mock Response
  let result = mockOrderUpdateStatusResponse.body

  return {
    data: result
  }
}

export default orderUpdateStatusService
