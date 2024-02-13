'use client'

// ---------------------------------------------------

type MenuItemProps = {
  link: string
  name: string
}

const MenuItem = ({ link, name }: MenuItemProps) => {
  return (
    <a href={link} className="text-sm font-semibold leading-6 text-gray-900">
      {name}
    </a>
  )
}

export default MenuItem
