'use client'

import DiscountPoint from '@/app/checkout/components/discount-point'
import Header3 from '@/components/typography/header3'
// import DiscountForm from '@/app/checkout/components/discount-form'

// ----------------------------------------------------------------------

const Discount = () => {
  return (
    <div className="w-full mx-auto rounded-lg bg-white border border-gray-200 p-3 text-gray-800 font-light mb-6">
      <Header3>Discounts</Header3>

      {/* <DiscountForm /> */}
      <DiscountPoint />
    </div>
  )
}

export default Discount
