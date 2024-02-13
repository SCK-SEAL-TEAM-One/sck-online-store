'use client'

import ButtonIcon from '@/components/button/button-icon'
import { TrashIcon } from '@heroicons/react/24/outline'

// ----------------------------------------------------------------------

type RemoveItemProps = {
  onClick: () => void
}

const RemoveItem = ({ onClick }: RemoveItemProps) => {
  return (
    <div className="flex pb-3">
      <ButtonIcon
        type="button"
        onClick={onClick}
        className="font-medium text-red-600 hover:text-red-500"
      >
        <TrashIcon width={24} height={24} />
      </ButtonIcon>
    </div>
  )
}

export default RemoveItem
