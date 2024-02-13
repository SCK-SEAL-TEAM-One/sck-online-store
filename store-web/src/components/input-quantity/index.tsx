'use client'

import { MinusIcon, PlusIcon } from '@heroicons/react/24/outline'

// ----------------------------------------------------------------------

type InputFieldProps = React.InputHTMLAttributes<HTMLInputElement> & {
  label?: string
  id?: string
  isHiddenLable?: boolean
  decrement: () => void
  increment: () => void
}

const InputQuantity = (props: InputFieldProps) => {
  return (
    <div>
      {props.label && !props.isHiddenLable ? (
        <label
          htmlFor={props.id}
          className="block mb-2 text-sm font-medium text-gray-900"
        >
          {props.label}
        </label>
      ) : null}

      <div className="relative flex items-center max-w-[10rem]">
        <button
          type="button"
          id="decrement-button"
          data-input-counter-decrement="quantity-input"
          className="bg-gray-100 hover:bg-gray-200 border border-gray-300 rounded-s-lg p-3 h-11 focus:ring-gray-100 focus:ring-2 focus:outline-none"
          onClick={props.decrement}
        >
          <MinusIcon className="w-3 h-3 text-gray-900" />
        </button>
        <input
          type="text"
          className="bg-gray-100 border border-gray-300 border-x-0 h-11 text-center text-gray-900 text-sm focus:ring-blue-500 focus:border-blue-500 block w-full py-2.5"
          {...props}
        />
        <button
          type="button"
          id="increment-button"
          data-input-counter-increment="quantity-input"
          className="bg-gray-100 hover:bg-gray-200 border border-gray-300 rounded-e-lg p-3 h-11 focus:ring-gray-100"
          onClick={props.increment}
        >
          <PlusIcon className="w-3 h-3 text-gray-900" />
        </button>
      </div>
    </div>
  )
}

export default InputQuantity
