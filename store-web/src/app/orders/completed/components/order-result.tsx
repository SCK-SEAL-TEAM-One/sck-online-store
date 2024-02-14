'use client'

import Header1 from '@/components/typography/header1'
import Header3 from '@/components/typography/header3'
import Text from '@/components/typography/text'
import { getShippingMethodById } from '@/utils/shipping'
import dayjs from 'dayjs'
import { useSearchParams } from 'next/navigation'
import { useEffect, useState } from 'react'

// ----------------------------------------------------------------------

type ShippingType = {
  id: number
  name: string
  shippingTime: string
  condition: string
  price: number
}

const OrderResult = () => {
  const search = useSearchParams()
  const orderId = search.get('order')
  const paymentDate = search.get('payment')
  const shippingMethodId = search.get('shipping')
  const trackingNumber = search.get('tracking')

  const [shipping, setShipping] = useState<ShippingType | null>(null)

  useEffect(() => {
    if (shippingMethodId) {
      const shipping = getShippingMethodById(Number(shippingMethodId))
      if (shipping) {
        setShipping(shipping)
      } else {
        setShipping(null)
      }
    }
  }, [shippingMethodId])

  return (
    <div className="text-gray-800 border-b border-gray-200 mb-4">
      <Header3 className="text-green-600">Payment successful</Header3>
      <Header1 className="text-6xl font-extrabold my-5">
        Thank you for your order
      </Header1>

      <Text id="order-success-text" size="md" className="text-gray-600 my-10">
        {`Date and time of payment `}
        <span id="order-success-order-payment-date" className="font-semibold">
          {dayjs(paymentDate).format('DD/MM/YYYY HH:mm:ss')}
        </span>
        {` Order number `}
        <a
          id="order-success-order-id"
          className="text-sm font-medium text-indigo-600"
          href="#?order_id=102323"
        >
          {orderId}
        </a>
        {shipping ? (
          <div>
            {` You can track the product via `}
            <span id="order-success-shipping-method" className="font-semibold">
              {shipping.name}
            </span>
            {` with number `}
            <a
              id="order-success-tracking-id"
              className="text-sm font-medium text-indigo-600"
              target="_blank"
              href={`https://th.kerryexpress.com/th/track/?track=${trackingNumber}`}
            >
              {trackingNumber}
            </a>
          </div>
        ) : null}
      </Text>
    </div>
  )
}

export default OrderResult
