// ----------------------------------------------------------------------------

import config from '@/config'
import { useUserStore } from '@/hooks/use-user-store'
import { decodeJWT } from '@/utils/jwt'
import Image from 'next/image'
import { useRouter } from 'next/navigation'
import { useEffect } from 'react'

const AuthLayout = ({ children }: { children: React.ReactNode }) => {
  const { user, setUser } = useUserStore()
  const accessToken = localStorage.getItem('accessToken')
  const route = useRouter()

  useEffect(() => {
    if (!accessToken) return

    if (!user) {
      const payload = decodeJWT(accessToken)
      setUser({
        userId: payload.user_id,
        firstName: payload.first_name,
        lastName: payload.last_name,
        username: payload.username
      })
    }

    return route.push('/product/list')
  }, [user, accessToken, route, setUser])

  return (
    <div className="flex h-screen">
      <div className="basis-[51%] flex justify-center items-center">
        {children}
      </div>
      <div className="basis-[49%]">
        <Image
          id="auth-page-background-image"
          src={`${config.logo.loginPage}`}
          alt="Siam Chamnankit Oneline Store"
          width={800}
          height={1042}
          className="h-screen w-full object-cover rounded-ss-[45px] rounded-es-[45px]"
        />
      </div>
    </div>
  )
}

export default AuthLayout
