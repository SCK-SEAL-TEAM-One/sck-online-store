import ShoppingCartView from '@/app/cart/shopping-cart'
import { UserInfo, useUserStore } from '@/hooks/use-user-store'
import Header from '@/layouts/common/header'
import { decodeJWT } from '@/utils/jwt'
import { useRouter } from 'next/navigation'
import { useEffect, useState } from 'react'

// ----------------------------------------------------------------------------

const CommonLayout = ({ children }: { children: React.ReactNode }) => {
  const [isShoppingCartOpen, setShoppingCartOpen] = useState(false)
  const { user, setUser } = useUserStore()
  const accessToken = localStorage.getItem('accessToken')
  const route = useRouter()

  useEffect(() => {
    if (!accessToken) {
      route.push('/auth/login')
      return
    }

    if (!user) {
      const payload = decodeJWT(accessToken)
      const user: UserInfo = {
        userId: payload.user_id,
        firstName: payload.first_name,
        lastName: payload.last_name,
        username: payload.username
      }
      setUser(user)
    }
  }, [user, accessToken, route, setUser])

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
