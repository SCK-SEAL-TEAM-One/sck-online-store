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

// Not used now
export type NotificationServiceResponse = {
  data?: {
    status: string
  }
  message?: string
}

const notificationService = async (
  notificationInfo: NotificationServiceRequest
): Promise<NotificationServiceResponse> => {
  try {
    const { status } = await axiosShoppingMallApi.post(
      `/api/v1/notification`,
      {
        user_id: notificationInfo.userId,
        order_id: notificationInfo.orderId,
        in_applicaition: notificationInfo.isApplication,
        email: notificationInfo.email,
        mobile: notificationInfo.mobile
      }
    )

    if (status === 200) {
      return {
        data: {
          status: 'success'
        }
      }
    }
    
    return {
      data: {
        status: 'error'
      }
    }
  } catch (error) {
    return handleServiceError(error)
  }
}

export default notificationService
