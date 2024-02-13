const Header3 = ({
  children,
  className
}: {
  children: React.ReactNode
  className?: string
}) => {
  return (
    <h3 className={`mb-5 text-lg font-medium text-gray-90 ${className}`}>
      {children}
    </h3>
  )
}

export default Header3
