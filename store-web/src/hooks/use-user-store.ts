import { create } from 'zustand'
import { devtools, persist } from 'zustand/middleware'

export interface UserInfo {
  userId: number
  firstName: string
  lastName: string
  username: string
}

type UserState = {
  user: UserInfo | null
  setUser: (user: UserInfo) => void
  clearUser: () => void
}

export const useUserStore = create<UserState>()(
  persist(
    devtools((set) => ({
      user: null,
      setUser: (user: UserInfo) => set({ user }),
      clearUser: () => set({ user: null })
    })),
    {
      name: 'user'
    }
  )
)
