'use client'

// ---------------------------------------------------

type MenuItemProps = {
  id?: string
  link: string
  name: string
}

const MenuItem = ({ id, link, name }: MenuItemProps) => {
  return (
    <a
      id={id}
      href={link}
      className="text-sm font-semibold leading-6 text-gray-900"
    >
      {name}
    </a>
  )
}

export default MenuItem
