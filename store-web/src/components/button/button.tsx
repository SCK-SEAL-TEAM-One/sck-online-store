'use client'

// ----------------------------------------------------------------------

type InputFieldProps = React.ButtonHTMLAttributes<HTMLButtonElement> & {
  children: React.ReactNode
  color?: string
  isblock?: string
  size?: string
}

const Button = (props: InputFieldProps) => {
  const { isblock = false, color = 'primary', size } = props
  let customClassName = ''

  if (isblock) {
    customClassName += 'w-full'
  }

  if (size === 'sm') {
    customClassName += ' text-sm font-normal px-6 py-2'
  } else {
    customClassName += ' text-base font-medium px-8 py-3'
  }

  if (color === 'primary') {
    customClassName += ' bg-indigo-600 text-white hover:bg-indigo-700 focus:ring-indigo-500'
  } else {
    customClassName += ' bg-gray-300 text-black hover:bg-gray-400 focus:ring-gray-500'
  }

  return (
    <button
      {...props}
      className={`${customClassName} ${props.className} flex items-center justify-center rounded-md border border-transparent focus:outline-none focus:ring-2 focus:ring-offset-2`}
    >
      {props.children}
    </button>
  )
}

export default Button
