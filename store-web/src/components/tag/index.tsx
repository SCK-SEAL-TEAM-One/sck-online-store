'use client'

import { XMarkIcon } from '@heroicons/react/24/outline'

// ----------------------------------------------------------------------

type TagProps = {
  name: string
  onRemove: () => void
}

const Tag = ({ name, onRemove }: TagProps) => {
  return (
    <span className="inline-flex items-center px-3 py-2 me-2 text-sm font-medium text-gray-600 bg-gray-100 rounded">
      <span>{name}</span>
      <button
        type="button"
        className="inline-flex items-center p-1 ms-2 text-sm text-gray-400 bg-transparent rounded-sm hover:bg-gray-200 "
        aria-label="Remove"
        onClick={onRemove}
      >
        <XMarkIcon width={12} height={12} />
      </button>
    </span>
  )
}

export default Tag
