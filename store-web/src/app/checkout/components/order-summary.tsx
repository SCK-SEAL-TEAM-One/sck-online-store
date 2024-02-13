'use client'

import SummaryText from '@/app/checkout/components/summary-text'
import useOrderStore from '@/hooks/use-order-store'
import Header3 from '@/components/typography/header3'

// ----------------------------------------------------------------------

const OrderSummary = () => {
  const { subTotal, totalPayment, receivePoint, point, shipping } =
    useOrderStore((state) => state)
  return (
    <div className="mb-6">
      <Header3>Summary</Header3>

      <div className="mb-6 pb-2 border-b border-gray-200 text-gray-800">
        <SummaryText text="Merchandise Subtotal" value={subTotal} />
        <SummaryText text="Shipping Fee" value={shipping.shippingFee} />
        <SummaryText
          text="Points Discount"
          textBeforeValue="-"
          format="number"
          className="text-red-600 font-semibold"
          value={point.burnPoint}
        />

        {/* Not use for now */}
        {/* <SummaryText text='Discount' value={20.0} />
        <SummaryText text='Tax (7%)' value={3.99} /> */}
      </div>
      <div>
        <SummaryText
          text="Receive Point"
          format="number"
          className="font-semibold"
          unit="Points"
          value={receivePoint}
        />
        <SummaryText
          className="font-semibold"
          text="Total Payment"
          value={totalPayment}
        />
      </div>
    </div>
  )
}

export default OrderSummary
