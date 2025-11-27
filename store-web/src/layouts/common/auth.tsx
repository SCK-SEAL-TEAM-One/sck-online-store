// ----------------------------------------------------------------------------

import config from '@/config'
import Image from 'next/image'

const AuthLayout = ({ children }: { children: React.ReactNode }) => {
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
