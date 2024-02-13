'use client'

import Text from '@/components/typography/text'
import { convertCurrency } from '@/utils/format'

// ----------------------------------------------------------------------

type SubTotalProps = {
  total: number
}

const SubTotal = ({ total }: SubTotalProps) => {
  return (
    <div className="flex justify-between text-base font-medium text-gray-900">
      <Text id="shopping-cart-subtotal-label">Subtotal</Text>
      <Text id="shopping-cart-subtotal-price">
        {convertCurrency(total, 'THB')}
      </Text>
    </div>
  )
}

export default SubTotal
