'use client'

import Image from '@/components/image'
import Header1 from '@/components/typography/header1'
import config from '@/config'

// ----------------------------------------------------------------------

const PaymentLogo = () => {
  return (
    <div className="flex flex-col items-center">
      <Image
        src={config.logo.sckPaymentGateway}
        width={100}
        height={100}
        alt="K Payment Gateway"
      />

      <Header1 className="text-green-600 mt-2">SCK Payment Gateway</Header1>
    </div>
  )
}

export default PaymentLogo
