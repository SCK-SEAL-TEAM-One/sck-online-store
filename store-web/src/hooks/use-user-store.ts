import { produce } from 'immer'
import { create } from 'zustand'
import { devtools, persist } from 'zustand/middleware'

type UserState = {
  userId: number | null
  setUserId: (userId: number) => void
}

export const useUserStore = create<UserState>()(
  persist(
    devtools((set) => ({
      userId: null,
      setUserId: (userId) => {
        set(
          produce((state) => {
            state.userId = userId
          })
        )
      }
    })),
    {
      name: 'user-id'
    }
  )
)
