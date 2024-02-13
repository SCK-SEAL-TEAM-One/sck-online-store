const Text = ({
  children,
  size,
  className
}: {
  children: React.ReactNode
  size?: string
  className?: string
}) => {
  let customClassName = ''
  if (size === 'sm') {
    customClassName = 'text-sm'
  } else if (size === 'md') {
    customClassName = 'text-md'
  } else if (size === 'lg') {
    customClassName = 'text-lg'
  } else if (size === 'xl') {
    customClassName = 'text-xl'
  } else {
    customClassName = 'text-base'
  }

  return (
    <p className={`text-gray-600 ${customClassName} ${className}`}>
      {children}
    </p>
  )
}

export default Text
