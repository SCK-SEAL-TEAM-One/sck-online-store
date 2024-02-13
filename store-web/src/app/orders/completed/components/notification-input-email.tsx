'use client'

import { EnvelopeIcon } from '@heroicons/react/24/outline'

// ----------------------------------------------------------------------

type NotificationInputEmailProps = {
  id?: string
  value: string
  onChange: (e: string) => void
}

const NotificationInputEmail = ({
  id,
  value,
  onChange
}: NotificationInputEmailProps) => {
  const handleEmailChange = (e: { target: { value: string } }) => {
    onChange(e.target.value)
  }

  return (
    <div id={id} className="relative mb-6">
      <div
        id={`${id}-icon`}
        className="absolute inset-y-0 start-0 flex items-center ps-3.5 pointer-events-none"
      >
        <EnvelopeIcon className="text-gray-400" width={24} height={24} />
      </div>
      <input
        id={`${id}-input`}
        type="email"
        className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full ps-11 p-2.5"
        placeholder="name@scrum123.com"
        onChange={handleEmailChange}
        value={value}
      />
    </div>
  )
}

export default NotificationInputEmail
