import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import client from '@/api/client'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<string | null>(localStorage.getItem('auth_user'))

  const isAuthenticated = computed(() => !!user.value)

  const checkAuth = async () => {
    try {
      const res = await client.get<{authenticated: boolean, user: string}>('/auth/check')
      if (res.data.authenticated) {
        user.value = res.data.user
        localStorage.setItem('auth_user', res.data.user)
      } else {
        throw new Error('Not authenticated')
      }
    } catch (e) {
      user.value = null
      localStorage.removeItem('auth_user')
      throw e
    }
  }

  const login = async (creds: {login: string, password: string}) => {
    const res = await client.post<{user: string}>('/login', creds)

    // Set user state
    user.value = res.data.user
    localStorage.setItem('auth_user', res.data.user)

    // Verify session
    await checkAuth()
  }

  const logout = async () => {
    try {
      await client.post('/logout')
    } catch {
      // Ignore errors during logout
    } finally {
      user.value = null
      localStorage.removeItem('auth_user')
    }
  }

  return {
    user,
    isAuthenticated,
    login,
    logout,
    checkAuth
  }
})
