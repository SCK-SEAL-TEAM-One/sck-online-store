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
      <Text>Subtotal</Text>
      <Text>{convertCurrency(total)}</Text>
    </div>
  )
}

export default SubTotal
