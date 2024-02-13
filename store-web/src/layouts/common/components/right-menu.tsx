'use client'

import Cart from '@/layouts/common/components/cart'
import Login from '@/layouts/common/components/login'
import { HeaderProps } from '@/layouts/common/header'

// ---------------------------------------------------

const RightMenu = ({ setShoppingCartOpen }: HeaderProps) => {
  return (
    <div className="flex flex-1 gap-x-12 justify-end">
      <Cart setShoppingCartOpen={setShoppingCartOpen} />
      <Login />
    </div>
  )
}

export default RightMenu
