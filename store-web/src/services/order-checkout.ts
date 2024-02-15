import { mockOrderCheckoutResponse } from '@/mock'
import axiosShoppingMallApi from '@/utils/axios'
import { handleServiceError } from '@/utils/helper'

// ------------------------------------------------

type Cart = {
  product_id: number
  quantity: number
}

type PaymentInformation = {
  card_name: string
  card_number: string
  expire_date: string
  cvv: string
}

type OrderRequest = {
  cart: Cart[]
  shipping_method_id: number
  shipping_address: string
  shipping_sub_district: string
  shipping_district: string
  shipping_province: string
  shipping_zip_code: string
  recipient_first_name: string
  recipient_last_name: string
  recipient_phone_number: string
  payment_method_id: number
  burn_point: number
  sub_total_price: number
  discount_price: number
  total_price: number
  payment_information: PaymentInformation
}

export type OrderCheckoutServiceResponse = {
  data?: {
    order_id: number
  }
  message?: string
}

const orderCheckoutService = async (
  order: OrderRequest
): Promise<OrderCheckoutServiceResponse> => {
  try {
    const { data } = await axiosShoppingMallApi.post(`/api/v1/order`, order)
    return {
      data: data
    }
  } catch (error) {
    return handleServiceError(error)
  }
}

export default orderCheckoutService
