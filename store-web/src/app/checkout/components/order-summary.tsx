'use client'

import SummaryText from '@/app/checkout/components/summary-text'
import Header3 from '@/components/typography/header3'
import useOrderStore from '@/hooks/use-order-store'

// ----------------------------------------------------------------------

const OrderSummary = () => {
  const { summary, totalPayment, receivePoint, shipping } = useOrderStore(
    (state) => state
  )
  return (
    <div className="mb-6">
      <Header3>Summary</Header3>

      <div className="mb-6 pb-2 border-b border-gray-200 text-gray-800">
        <SummaryText
          id="order-summary-subtotal"
          text="Merchandise Subtotal"
          value={summary.total_price_thb}
        />
        <SummaryText
          id="order-summary-receive-point"
          text="Receive Points"
          format="number"
          className="font-semibold"
          unit="Points"
          value={receivePoint}
        />
        <SummaryText
          id="order-summary-shipping-fee"
          text="Shipping Fee"
          value={shipping.shippingFee}
        />
        {/* <SummaryText
          id="order-summary-point-discount"
          text="Points Discount"
          textBeforeValue="-"
          format="number"
          className="text-red-600 font-semibold"
          value={point.burnPoint}
        /> */}

        {/* Not use for now */}
        {/* <SummaryText text='Discount' value={20.0} />
        <SummaryText text='Tax (7%)' value={3.99} /> */}
      </div>
      <div>
        <SummaryText
          id="order-summary-total-payment"
          className="font-semibold"
          text="Total Payment"
          value={totalPayment}
        />
      </div>
    </div>
  )
}

export default OrderSummary
