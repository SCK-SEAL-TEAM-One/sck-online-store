import { mockOrderCheckoutResponse } from '@/mock'

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
  order_id: number
}

const orderCheckoutService = async (
  order: OrderRequest
): Promise<OrderCheckoutServiceResponse> => {
  //   const result = await axiosShoppingMallApi.put(
  //     `/api/v1/order`,
  //   )

  // Mock Response
  let result = mockOrderCheckoutResponse.body

  return result
}

export default orderCheckoutService
