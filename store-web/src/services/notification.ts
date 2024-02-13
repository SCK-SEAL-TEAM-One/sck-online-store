import { mockNotificationResponse } from '@/mock'
import axiosShoppingMallApi from '@/utils/axios'
import { handleServiceError } from '@/utils/helper'

// ------------------------------------------------

type NotificationServiceRequest = {
  userId: number
  orderId: number
  isApplication: boolean
  email: string
  mobile: string
}

export type OrderCheckoutServiceResponse = {
  data?: {
    status: string
  }
  message?: string
}

const notificationService = async (
  notificationInfo: NotificationServiceRequest
): Promise<OrderCheckoutServiceResponse> => {
  // try {
  //   const { status } = await axiosShoppingMallApi.post(
  //     `/api/v1/notification`,
  //     {
  //       user_id: notificationInfo.userId,
  //       order_id: notificationInfo.orderId,
  //       in_applicaition: notificationInfo.isApplication,
  //       email: notificationInfo.email,
  //       mobile: notificationInfo.mobile
  //     }
  //   )

  //   if (status === 200) {
  //     return {
  //       data: {
  //         status: 'success'
  //       }
  //     }
  //   }
  // } catch (error) {
  //   return handleServiceError(error)
  // }

  // Mock Response
  let result = mockNotificationResponse.status

  return {
    data: {
      status: result === 200 ? 'success' : 'fail'
    }
  }
}

export default notificationService
