'use client'

import ShippingMethodItem from '@/app/checkout/components/shipping-method-item'
import Header3 from '@/components/typography/header3'

import SHIPPING_METHOD from '@/assets/data/shipping_method.json'
import useOrderStore from '@/hooks/use-order-store'
import { useEffect } from 'react'

// ----------------------------------------------------------------------

const ShippingMethod = () => {
  const { shipping, setShippingMethod } = useOrderStore((state) => state)
  const { provinceId } = shipping.shippingInformation
  const shippingMethodId = shipping.shippingMethod
  const BMPProvinceIDs = [1, 2, 3, 4, 58, 59] // กรุงเทพมหานครและปริมณฑล (Bangkok Metropolitan Region: BMR) ประกอบด้วย 6 จังหวัด คือ กรุงเทพมหานคร, นนทบุรี, ปทุมธานี, สมุทรปราการ, นครปฐม, และสมุทรสาคร

  useEffect(() => {
    if (shippingMethodId === 3 && !BMPProvinceIDs.includes(provinceId)) {
      const defaultShippingMethod = SHIPPING_METHOD[0]
      setShippingMethod(defaultShippingMethod.id, defaultShippingMethod.price)
      alert(
        `Out of Delivery Area of ${SHIPPING_METHOD[2].name.toLocaleUpperCase()}. Please recheck again.`
      )
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [provinceId, shippingMethodId])

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
        {SHIPPING_METHOD.map((shipp) => {
          const inDeliveryArea = BMPProvinceIDs.includes(
            shipping.shippingInformation.provinceId
          )
          return (
            <ShippingMethodItem
              {...shipp}
              key={`shipping-${shipp.id}`}
              onChange={handleShippingMethodChange}
              shippingMethodSelected={shipping.shippingMethod}
              disabled={shipp.id === 3 && !inDeliveryArea}
            />
          )
        })}
      </ul>
    </div>
  )
}

export default ShippingMethod
