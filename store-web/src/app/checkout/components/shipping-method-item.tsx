'use client'

import { convertCurrency } from '@/utils/format'

// ----------------------------------------------------------------------

type ShippingMethodItemProps = {
  id: number
  name: string
  shippingTime: string
  condition: string
  price: number
  onChange: (event: React.ChangeEvent<HTMLInputElement>) => void
  shippingMethodSelected: number
}

const ShippingMethodItem = ({
  id,
  name,
  shippingTime,
  condition,
  price,
  onChange,
  shippingMethodSelected
}: ShippingMethodItemProps) => {
  return (
    <li id={`shipping-method-${id}-card`}>
      <input
        type="radio"
        id={`shipping-method-${id}-input`}
        name="shipping-method"
        value={id}
        className="hidden peer"
        onChange={onChange}
        data-fee={price}
        checked={shippingMethodSelected === id}
        required
      />
      <label
        htmlFor={`shipping-method-${id}-input`}
        className="inline-flex items-center justify-between w-full p-5 text-gray-700 bg-white border border-gray-200 rounded-lg cursor-pointer peer-checked:border-blue-600 peer-checked:text-blue-600 hover:text-gray-600 hover:bg-gray-100"
      >
        <div className="block">
          <div
            id={`shipping-method-${id}-name`}
            className="w-full text-lg font-semibold first-letter:uppercase"
          >
            {name}
          </div>
          <div
            id={`shipping-method-${id}-time`}
            className="w-full text-gray-600 text-sm py-2"
          >
            {shippingTime}
          </div>
          {condition ? (
            <div
              id={`shipping-method-${id}-note`}
              className="w-full text-red-400 text-xs"
            >
              {condition ?? condition}
            </div>
          ) : (
            <div className="mt-6"></div>
          )}

          <div
            id={`shipping-method-${id}-fee`}
            className="w-full mt-2 font-semibold"
          >
            {convertCurrency(price, 'THB')}
          </div>
        </div>
      </label>
    </li>
  )
}

export default ShippingMethodItem
