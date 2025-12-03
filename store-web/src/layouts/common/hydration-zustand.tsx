'use client'

import { useUserStore } from '@/hooks/use-user-store'
import { useEffect, useState } from 'react'
// ----------------------------------------------------------------------

const HydrationZustand = ({ children }: { children: React.ReactNode }) => {
  const [isHydrated, setIsHydrated] = useState(false)

  useEffect(() => {
    setIsHydrated(useUserStore.persist.hasHydrated())
  }, [])

  if (!isHydrated) {
    return (
      <div className="w-screen h-screen flex justify-center items-center">
        <span className="loading loading-spinner loading-xl"></span>
      </div>
    )
  }

  return <>{children}</>
}

export default HydrationZustand
