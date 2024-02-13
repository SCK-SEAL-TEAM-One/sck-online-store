'use client'

import ShippingMethodItem from '@/app/checkout/components/shipping-method-item'
import Header3 from '@/components/typography/header3'

import SHIPPING_METHOD from '@/assets/data/shipping_method.json'
import useOrderStore from '@/hooks/use-order-store'

// ----------------------------------------------------------------------

const ShippingMethod = () => {
  const { shipping, setShippingMethod } = useOrderStore((state) => state)

  const handleShippingMethodChange = (
    e: React.ChangeEvent<HTMLInputElement>
  ) => {
    const fee = e.target.getAttribute('data-fee')
    setShippingMethod(Number(e.target.value), Number(fee))
  }

  return (
    <div className="mb-6 border-b border-gray-200 pb-6">
      <Header3>Delivery method</Header3>

      <ul
        id="shipping-method-list"
        className="grid w-full gap-2 md:grid-cols-3"
      >
        {SHIPPING_METHOD.map((shipp) => (
          <ShippingMethodItem
            {...shipp}
            key={`shipping-${shipp.id}`}
            onChange={handleShippingMethodChange}
            shippingMethodSelected={shipping.shippingMethod}
          />
        ))}
      </ul>
    </div>
  )
}

export default ShippingMethod
