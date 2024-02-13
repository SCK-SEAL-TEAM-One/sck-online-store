'use client'

import Logo from '@/layouts/common/components/logo'
import MenuList from '@/layouts/common/components/menu-list'
import RightMenu from '@/layouts/common/components/right-menu'

// ----------------------------------------------------------------------

export type HeaderProps = {
  setShoppingCartOpen: (set: boolean) => void
}

const Header = ({ setShoppingCartOpen }: HeaderProps) => {
  return (
    <header className="bg-white">
      <nav
        className="mx-auto flex max-w-7xl items-center justify-between p-6 lg:px-8"
        aria-label="Global"
      >
        <Logo />
        <MenuList />
        <RightMenu setShoppingCartOpen={setShoppingCartOpen} />
      </nav>
    </header>
  )
}

export default Header
