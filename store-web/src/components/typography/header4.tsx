'use client'

// ----------------------------------------------------------------------

const Header4 = ({
  id,
  children,
  className
}: {
  id?: string
  children: React.ReactNode
  className?: string
}) => {
  return (
    <h3 id={id} className={`text-md font-medium text-gray-90 ${className}`}>
      {children}
    </h3>
  )
}

export default Header4
