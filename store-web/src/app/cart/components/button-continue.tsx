'use client'

import { ArrowRightIcon } from '@heroicons/react/24/outline'

// ----------------------------------------------------------------------

type ButtonContiueShoppingProps = {
  onClick?: () => void
}

const ButtonContiueShopping = ({ onClick }: ButtonContiueShoppingProps) => {
  return (
    <button
      type="button"
      className="font-medium text-indigo-600 hover:text-indigo-500 flex gap-1 items-center"
      onClick={onClick}
    >
      Continue Shopping
      <ArrowRightIcon width={16} />
    </button>
  )
}

export default ButtonContiueShopping
