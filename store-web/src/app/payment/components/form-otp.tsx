'use client'

import Button from '@/components/button/button'
import InputField from '@/components/input-field'
import { isNumber } from '@/utils/format'

// ----------------------------------------------------------------------

type FormOtpProps = {
  otpRef: string
  otp: string
  onChange: (e: { target: { value: string } }) => void
}

const FormOtp = ({ otpRef, otp, onChange }: FormOtpProps) => {
  return (
    <div className="-mx-2 flex items-end justify-between">
      <div className="flex-grow px-2 lg:max-w-sm">
        <InputField
          id="otp"
          type="text"
          label={`OTP (Ref: ${otpRef})`}
          placeholder="XXXXXX"
          maxLength={6}
          onChange={onChange}
          value={otp !== '' ? otp : ''}
          required
        />
      </div>
      <div className="px-2 pb-2">
        <Button
          size="sm"
          className="bg-gray-300 text-black hover:bg-gray-400"
          color="default"
        >
          Request OTP
        </Button>
      </div>
    </div>
  )
}

export default FormOtp
