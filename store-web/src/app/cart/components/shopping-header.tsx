'use client'

import { Dialog } from '@headlessui/react'
import { XMarkIcon } from '@heroicons/react/24/outline'

// ----------------------------------------------------------------------

type ShoppingHeaderProps = {
  id?: string
  closeShoppingCart: () => void
}

const ShoppingHeader = ({ id, closeShoppingCart }: ShoppingHeaderProps) => {
  return (
    <div className="flex items-start justify-between">
      <Dialog.Title id={id} className="text-lg font-medium text-gray-900">
        Shopping cart
      </Dialog.Title>
      <div className="ml-3 flex h-7 items-center">
        <button
          id="shopping-cart-close-btn"
          type="button"
          className="relative -m-2 p-2 text-gray-400 hover:text-gray-500"
          onClick={closeShoppingCart}
        >
          <span className="absolute -inset-0.5" />
          <span className="sr-only">Close panel</span>
          <XMarkIcon className="h-6 w-6" aria-hidden="true" />
        </button>
      </div>
    </div>
  )
}

export default ShoppingHeader
