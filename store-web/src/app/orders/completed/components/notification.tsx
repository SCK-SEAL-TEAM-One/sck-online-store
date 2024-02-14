'use client'

import NotificationInApplication from '@/app/orders/completed/components/notification-in-application'
import NotificationInputEmail from '@/app/orders/completed/components/notification-input-email'
import NotificationInputMobile from '@/app/orders/completed/components/notification-input-mobile'
import Button from '@/components/button/button'
import Header2 from '@/components/typography/header2'
import Text from '@/components/typography/text'
import { useSearchParams } from 'next/navigation'
import { useState } from 'react'

// ----------------------------------------------------------------------

const Notification = () => {
  const search = useSearchParams()
  const orderId = search.get('order_id')

  const [email, setEmail] = useState('')
  const [mobile, setMobile] = useState('')
  const [isApplication, setIsApplication] = useState(false)

  const sendNotification = async () => {
    if (
      confirm(
        'Send notification completed.\n\nClick OK for go to Product lists.'
      )
    ) {
      window.location.href = '/product/list'
    }

    // ยังไม่ได้ใช้งาน
    // const result = await notificationService({
    //   userId: 1,
    //   orderId: Number(orderId),
    //   email,
    //   mobile,
    //   isApplication
    // })

    // if (result.data?.status === 'success') {
    //   if (confirm('Send notification completed.\n\nClick OK for go to Product lists.') == true) {
    //     window.location.href = '/product/list'
    //   }
    // }
  }

  return (
    <div>
      <Header2 className="mb-0">Notifications</Header2>
      <Text size="sm" className="mb-5 text-gray-300">
        Please enter your email or mobile number for send notification.
      </Text>

      <NotificationInputEmail
        id="notification-form-email"
        onChange={setEmail}
        value={email}
      />
      <NotificationInputMobile
        id="notification-form-mobile"
        onChange={setMobile}
        value={mobile}
      />
      <NotificationInApplication
        id="notification-form-in-application"
        onChange={setIsApplication}
      />

      <Button id="send-notification-btn" onClick={sendNotification}>
        Send Notification
      </Button>
    </div>
  )
}

export default Notification
