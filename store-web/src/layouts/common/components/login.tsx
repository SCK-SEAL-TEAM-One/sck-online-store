'use client'

import { ArrowRightEndOnRectangleIcon } from '@heroicons/react/24/outline'

// ---------------------------------------------------

const Login = () => {
  return (
    <a
      id='header-menu-login'
      href="#"
      className="text-sm font-semibold leading-6 text-gray-900 flex items-center gap-1"
    >
      Log in
      <ArrowRightEndOnRectangleIcon width={24} />
    </a>
  )
}

export default Login
