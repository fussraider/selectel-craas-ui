import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import client from '@/api/client'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem('auth_token'))
  const user = ref<string | null>(localStorage.getItem('auth_user'))

  const isAuthenticated = computed(() => !!token.value)

  const login = async (creds: {login: string, password: string}) => {
    const res = await client.post<{token: string}>('/login', creds)
    token.value = res.data.token
    localStorage.setItem('auth_token', res.data.token)

    user.value = creds.login
    localStorage.setItem('auth_user', creds.login)
  }

  const logout = () => {
    token.value = null
    user.value = null
    localStorage.removeItem('auth_token')
    localStorage.removeItem('auth_user')
    // We might need to redirect here or let the caller handle it.
    // Often it's cleaner to let the router guard handle redirection after state change.
  }

  return {
    token,
    user,
    isAuthenticated,
    login,
    logout
  }
})
