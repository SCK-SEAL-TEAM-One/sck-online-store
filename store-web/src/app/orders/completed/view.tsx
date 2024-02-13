'use client'

import Notification from '@/app/orders/completed/components/notification'
import OrderResult from '@/app/orders/completed/components/order-result'

// ----------------------------------------------------------------------

const SuccessView = () => {
  return (
    <div className="bg-white flex justify-center items-center min-h-[calc(100vh-88px)]">
      <div className="mx-auto max-w-2xl px-4">
        <OrderResult />
        <Notification />
      </div>
    </div>
  )
}

export default SuccessView
