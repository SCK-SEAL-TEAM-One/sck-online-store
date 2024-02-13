const Header1 = ({
  children,
  className
}: {
  children: React.ReactNode
  className?: string
}) => {
  return (
    <h1 className={`mb-5 text-2xl font-bold text-gray-900 ${className}`}>
      {children}
    </h1>
  )
}

export default Header1
