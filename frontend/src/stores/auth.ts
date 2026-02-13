import { defineStore } from 'pinia'
import axios from 'axios'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    status: 'checking' as 'checking' | 'authenticated' | 'error',
    errorMessage: '',
  }),
  actions: {
    async checkAuth() {
      try {
        await axios.get('/api/auth/status')
        this.status = 'authenticated'
      } catch (err) {
        this.status = 'error'
        if (axios.isAxiosError(err)) {
          this.errorMessage = err.response?.data || err.message
        } else {
          this.errorMessage = String(err)
        }
      }
    },
  },
})
