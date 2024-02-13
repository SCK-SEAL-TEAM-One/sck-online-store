'use client'

import MenuItem from '@/layouts/common/components/menu-item'
import { Popover } from '@headlessui/react'

// ---------------------------------------------------

const MenuList = () => {
  return (
    <Popover.Group className="hidden lg:flex lg:gap-x-12">
      <MenuItem id='header-menu-home' link="/product" name="Home" />
      <MenuItem id='header-menu-for-kids' link="/product" name="For Kids" />
      <MenuItem id='header-menu-categories' link="/product" name="Categories" />
    </Popover.Group>
  )
}

export default MenuList
