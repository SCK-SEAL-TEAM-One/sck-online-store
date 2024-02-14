'use client'

import FormOtp from '@/app/payment/components/form-otp'
import PaymentLogo from '@/app/payment/components/logo'
import PaymentText from '@/app/payment/components/payment-text'
import Button from '@/components/button/button'
import Text from '@/components/typography/text'
import orderConfirmPaymentService from '@/services/confirm-payment'
import { convertCurrency, isNumber } from '@/utils/format'
import dayjs from 'dayjs'
import { useSearchParams } from 'next/navigation'
import { useState } from 'react'

// ----------------------------------------------------------------------

const PaymentView = () => {
  const searchParams = useSearchParams()
  const orderId = searchParams.get('id')
  const total = Number(searchParams.get('total')) || 0
  const cardNumber = searchParams.get('card')

  const today = dayjs().format('DD/MM/YYYY')

  const [otpRef, setOtpRef] = useState('AXYZ')
  const [otp, setOtp] = useState('')

  const handleOtpChange = (e: { target: { value: string } }) => {
    if (isNumber(e.target.value)) {
      setOtp(e.target.value)
    }
  }

  const handlePaymentConfirm = async () => {
    if (otp.length === 6) {
      const result = await orderConfirmPaymentService({
        orderId: Number(orderId),
        otp: Number(otp),
        otpRef: otpRef
      })

      if (result.data) {
        const convertResultToObject = {
          order: result.data.order_id.toString(),
          shipping: result.data.shipping_method_id.toString(),
          payment: result.data.payment_date,
          tracking: result.data.tracking_number
        }

        const query = new URLSearchParams(convertResultToObject).toString()
        window.location.href = `/orders/completed?${query}`
      }
    }
  }

  const handleCancle = () => {
    window.location.href = '/checkout'
  }

  return (
    <div className="bg-white">
      <div className="min-h-[100vh] flex flex-col items-center mx-auto max-w-2xl px-4 py-[105px]">
        <PaymentLogo />
        <Text size="md" className="text-black mt-5">
          Please check the accuracy of your identity verification message. To
          increase security in making this payment transaction.
        </Text>

        <div className="my-10">
          <PaymentText label="Merchant" text="SCK Shopping Mall" />
          <PaymentText label="Amount" text={convertCurrency(total, 'THB')} />
          <PaymentText label="Date" text={today} />
          <PaymentText
            label="Card Number"
            text={`**** **** **** ${cardNumber}`}
          />
        </div>

        <FormOtp otpRef={otpRef} otp={otp} onChange={handleOtpChange} />

        <div className="flex gap-2 mt-12">
          <Button size="sm" color="primary" onClick={handlePaymentConfirm}>
            PAY NOW
          </Button>
          <Button size="sm" color="default" onClick={handleCancle}>
            Cancle
          </Button>
        </div>
      </div>
    </div>
  )
}

export default PaymentView
