import { UserInfo } from '@/apis'
// import Cookies from 'js-cookie'
import { create } from 'zustand'

interface AuthState {
  auth: {
    user: UserInfo | null
    setUser: (user: UserInfo | null) => void
    reset: () => void
  }
}

export const useAuthStore = create<AuthState>()((set) => {
  return {
    auth: {
      user: null,
      setUser: (user) =>
        set((state) => ({ ...state, auth: { ...state.auth, user } })),

      reset: () =>
        set((state) => {
          return {
            ...state,
            auth: { ...state.auth, user: null, accessToken: '' },
          }
        }),
    },
  }
})

// export const useAuth = () => useAuthStore((state) => state.auth)
