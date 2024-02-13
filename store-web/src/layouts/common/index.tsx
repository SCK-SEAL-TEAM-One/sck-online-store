import ShoppingCartView from '@/app/cart/shopping-cart'
import Header from '@/layouts/common/header'
import { useState } from 'react'

// ----------------------------------------------------------------------------

const CommonLayout = ({ children }: { children: React.ReactNode }) => {
  const [isShoppingCartOpen, setShoppingCartOpen] = useState(false)

  return (
    <>
      <Header setShoppingCartOpen={setShoppingCartOpen} />
      <ShoppingCartView
        openShoppingCart={isShoppingCartOpen}
        setOpenShoppingCart={setShoppingCartOpen}
      />
      {children}
    </>
  )
}

export default CommonLayout
