import { useRouter } from 'vue-router'

export function useAuth() {
  const router = useRouter()

  function logout() {
    router.push({ path: '/login' })
  }

  function toHome() {
    router.push({ path: '/' })
  }

  function login() {
    toHome()
  }

  return {
    logout,
    login,
  }
}
