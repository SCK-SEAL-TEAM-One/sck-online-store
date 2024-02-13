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
  const orderId = search.get('order_id')
  const paymentDate = search.get('payment_date')
  const shippingMethodId = search.get('shipping_method_id')
  const trackingId = search.get('tracking_id')

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

      <Text size="md" className="text-gray-600 my-10">
        {`Date and time of payment ${dayjs(paymentDate).format('DD/MM/YYYY HH:mm:ss')} Order number `}
        <a
          className="text-sm font-medium text-indigo-600"
          href="#?order_id=102323"
        >
          {orderId}
        </a>
        {shipping ? (
          <div>
            {` You can track the product via ${shipping.name} with number `}
            <a
              className="text-sm font-medium text-indigo-600"
              href={`#?tracking_id=${trackingId}`}
            >
              {trackingId}
            </a>
          </div>
        ) : null}
      </Text>

      {/* <Text size='md' className='text-gray-600 my-10'>
        {'วันเวลาที่ชําระเงิน 1/3/2020 13:30:00 หมายเลขคําสั่งซื้อ '}
        <a
          className='text-sm font-medium text-indigo-600'
          href='/orders/102323'
        >
          102323
        </a>
        {' คุณสามารถติดตามสินค้าผ่านช่องทาง Kerry ด้วยหมายเลข '}
        <a
          className='text-sm font-medium text-indigo-600'
          href='#?tracking_id=51547878755545848512'
        >
          51547878755545848512
        </a>
      </Text> */}
    </div>
  )
}

export default OrderResult
