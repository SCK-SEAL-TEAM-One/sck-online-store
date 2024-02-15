'use client'

// ----------------------------------------------------------------------

const Header3 = ({
  children,
  id,
  className
}: {
  children: React.ReactNode
  id?: string
  className?: string
}) => {
  return (
    <h3 id={id} className={`mb-5 text-lg font-medium text-gray-90 ${className}`}>
      {children}
    </h3>
  )
}

export default Header3
