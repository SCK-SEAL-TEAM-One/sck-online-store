'use client'

import { PhoneIcon } from '@heroicons/react/24/outline'

// ----------------------------------------------------------------------

type NotificationInputMobileProps = {
  value: string
  onChange: (e: string) => void
}

const NotificationInputMobile = ({
  value,
  onChange
}: NotificationInputMobileProps) => {
  const handleMobileChange = (e: { target: { value: string } }) => {
    onChange(e.target.value)
  }

  return (
    <div className="relative mb-6">
      <div className="absolute inset-y-0 start-0 flex items-center ps-3.5 pointer-events-none">
        <PhoneIcon className="text-gray-400" width={24} height={24} />
      </div>
      <input
        id="notification-mobile"
        type="text"
        className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full ps-11 p-2.5"
        placeholder="0923456789"
        value={value}
        onChange={handleMobileChange}
      />
    </div>
  )
}

export default NotificationInputMobile
