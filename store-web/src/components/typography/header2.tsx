'use client'

// ----------------------------------------------------------------------

const Header2 = ({
  children,
  id,
  className = ''
}: {
  children: React.ReactNode
  id?: string
  className?: string
}) => {
  return (
    <h2
      id={id}
      className={`mb-5 text-xl font-medium text-gray-600 ${className}`}
    >
      {children}
    </h2>
  )
}

export default Header2
