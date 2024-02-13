'use client'

import Text from '@/components/typography/text'

// ----------------------------------------------------------------------

type PaymentTextType = {
  label: string
  text: string
}

const PaymentText = ({ label, text }: PaymentTextType) => {
  return (
    <div className="flex items-center gap-2">
      <Text size="md" className="font-semibold text-right w-28">
        {`${label}:`}
      </Text>
      <Text>{text}</Text>
    </div>
  )
}

export default PaymentText
