const Header2 = ({
  children,
  className = ''
}: {
  children: React.ReactNode
  className?: string
}) => {
  return (
    <h2 className={`mb-5 text-xl font-medium text-gray-600 ${className}`}>
      {children}
    </h2>
  )
}

export default Header2
