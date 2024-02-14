'use client'

import Image from '@/components/image'
import config from '@/config'

// ---------------------------------------------------

const Logo = () => {
  return (
    <div className="flex lg:flex-1">
      <a id="header-logo" href="/product/list" className="-m-1.5 p-1.5">
        <span id="header-logo-text" className="sr-only">
          SCK Shopping Mall
        </span>
        <Image
          id="header-logo-image"
          className="h-10 w-auto"
          src={config.logo.shoppingMall}
          alt="SCK Shopping Mall"
          width={40}
          height={40}
        />
      </a>
    </div>
  )
}

export default Logo
