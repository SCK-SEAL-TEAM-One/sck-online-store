import axiosShoppingMallApi from '@/utils/axios'
import { handleServiceError } from '@/utils/helper'

// ------------------------------------------------

type OrderConfirmPaymentServiceRequest = {
  order_number: number
  otp: number
  otpRef: string
}

export type OrderConfirmPaymentDetailType = {
  order_number: number
  payment_date: string
  shipping_method_id: number
  tracking_number: string
}

export type OrderConfirmPaymentServiceResponse = {
  data?: OrderConfirmPaymentDetailType
  message?: string
}

const orderConfirmPaymentService = async ({
  order_number,
  otp,
  otpRef
}: OrderConfirmPaymentServiceRequest): Promise<OrderConfirmPaymentServiceResponse> => {
  try {
    const { data } = await axiosShoppingMallApi.post(`/api/v1/confirmPayment`, {
      order_number: order_number,
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
