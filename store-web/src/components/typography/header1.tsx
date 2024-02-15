'use client'

// ----------------------------------------------------------------------

const Header1 = ({
  children,
  id,
  className
}: {
  children: React.ReactNode
  id?: string
  className?: string
}) => {
  return (
    <h1
      id={id}
      className={`mb-5 text-2xl font-bold text-gray-900 ${className}`}
    >
      {children}
    </h1>
  )
}

export default Header1
