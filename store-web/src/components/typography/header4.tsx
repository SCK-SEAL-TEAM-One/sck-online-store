const Header4 = ({
  children,
  className
}: {
  children: React.ReactNode
  className?: string
}) => {
  return (
    <h3 className={`text-md font-medium text-gray-90 ${className}`}>
      {children}
    </h3>
  )
}

export default Header4
