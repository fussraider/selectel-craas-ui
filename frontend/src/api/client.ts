import axios from 'axios'
import { useNotificationStore } from '@/stores/notifications'

const client = axios.create({
  baseURL: window.config?.apiBaseUrl || '/api',
  headers: {
    'Content-Type': 'application/json'
  }
})

client.interceptors.response.use(
  response => response,
  error => {
    if (axios.isAxiosError(error) && error.response?.status === 401) {
      // Clear auth data
      localStorage.removeItem('auth_user')

      // Redirect to login if not already there to prevent loop
      if (!window.location.pathname.startsWith('/login')) {
        window.location.href = '/login'
      }
    }

    // We only want to show notification if not cancelled manually by axios
    if (!axios.isCancel(error)) {
      try {
        const store = useNotificationStore()
        store.addNotification(formatError(error), 'error')
      } catch (e) {
        // Fallback if Pinia is somehow not ready
        console.error("Failed to add notification:", e)
      }
    }

    return Promise.reject(error)
  }
)

// Helper to format error messages consistently
export const formatError = (err: unknown): string => {
  if (axios.isAxiosError(err)) {
    const data = err.response?.data
    if (typeof data === 'string') {
        // Avoid returning huge HTML pages from proxies or bad gateways
        if (data.trim().startsWith('<')) {
            return `Server error: ${err.message}`
        }
        return data
    }
    if (data && typeof data === 'object' && 'error' in data) {
      return String((data as { error: unknown }).error)
    }
    return err.message
  } else if (err instanceof Error) {
    return err.message
  }
  return String(err)
}

export default client
