import { mockNotificationResponse } from '@/mock'

// ------------------------------------------------

type NotificationServiceRequest = {
  userId: number
  orderId: number
  isApplication: boolean
  email: string
  mobile: string
}

export type OrderCheckoutServiceResponse = {
  status: string
}

const notificationService = async (
  data: NotificationServiceRequest
): Promise<OrderCheckoutServiceResponse> => {
  //   const result = await axiosShoppingMallApi.put(
  //     `/api/v1/order`,
  //   )

  // Mock Response
  let result = mockNotificationResponse.status

  return {
    status: result === 200 ? 'success' : 'fail'
  }
}

export default notificationService
