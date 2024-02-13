'use client'

import MenuItem from '@/layouts/common/components/menu-item'
import { Popover } from '@headlessui/react'

// ---------------------------------------------------

const MenuList = () => {
  return (
    <Popover.Group className="hidden lg:flex lg:gap-x-12">
      <MenuItem link="/product" name="Home" />
      <MenuItem link="/product" name="For Kids" />
      <MenuItem link="/product" name="Categories" />
    </Popover.Group>
  )
}

export default MenuList
