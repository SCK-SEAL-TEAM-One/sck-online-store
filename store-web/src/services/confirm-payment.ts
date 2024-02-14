import axiosShoppingMallApi from '@/utils/axios'
import { handleServiceError } from '@/utils/helper'

// ------------------------------------------------

type OrderConfirmPaymentServiceRequest = {
  orderId: number
  otp: number
  otpRef: string
}

export type OrderConfirmPaymentDetailType = {
  order_id: number
  payment_date: string
  shipping_method_id: number
  tracking_number: string
}

export type OrderConfirmPaymentServiceResponse = {
  data?: OrderConfirmPaymentDetailType
  message?: string
}

const orderConfirmPaymentService = async ({
  orderId,
  otp,
  otpRef
}: OrderConfirmPaymentServiceRequest): Promise<OrderConfirmPaymentServiceResponse> => {
  try {
    const { data } = await axiosShoppingMallApi.post(`/api/v1/confirmPayment`, {
      order_id: orderId,
      otp: otp,
      ref_otp: otpRef
    })
    return {
      data: data
    }
  } catch (error) {
    return handleServiceError(error)
  }
}

export default orderConfirmPaymentService
