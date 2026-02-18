import axios from 'axios'

const client = axios.create({
  baseURL: '/api',
  headers: {
    'Content-Type': 'application/json'
  }
})

// Helper to format error messages consistently
export const formatError = (err: unknown): string => {
  if (axios.isAxiosError(err)) {
    const data = err.response?.data
    if (typeof data === 'string') return data
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
