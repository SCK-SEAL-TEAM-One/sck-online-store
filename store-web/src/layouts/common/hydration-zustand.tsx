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
    return <p>Loading data ...</p> // #TODO: Full loading component
  }

  return <>{children}</>
}

export default HydrationZustand
