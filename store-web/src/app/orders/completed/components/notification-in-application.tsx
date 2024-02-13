'use client'

// ----------------------------------------------------------------------

type NotificationInApplicationProps = {
  onChange: (e: boolean) => void
}

const NotificationInApplication = ({
  onChange
}: NotificationInApplicationProps) => {
  const handleIsApplicationChange = (e: { target: { checked: boolean } }) => {
    onChange(e.target.checked)
  }

  return (
    <div className="flex items-center mb-6">
      <input
        id="in-applications"
        type="checkbox"
        value="true"
        className="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500 focus:ring-2"
        onChange={handleIsApplicationChange}
      />
      <label
        htmlFor="in-applications"
        className="ms-2 text-sm font-medium text-gray-900"
      >
        In Applications
      </label>
    </div>
  )
}

export default NotificationInApplication
