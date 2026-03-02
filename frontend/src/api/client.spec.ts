import { describe, it, expect, vi, beforeEach } from 'vitest'
import axios from 'axios'
import client from './client'
import { createPinia, setActivePinia } from 'pinia'
import { useNotificationStore } from '@/stores/notifications'

// Mock window.location
const mockLocation = {
  pathname: '/',
  href: ''
}

Object.defineProperty(window, 'location', {
  value: mockLocation,
  writable: true
})

interface MockAxiosError extends Error {
  isAxiosError: boolean
  response?: {
    status: number
    data?: unknown
  }
}

type InterceptorHandlers = { handlers: Array<{ rejected: (error: unknown) => Promise<never> }> }

describe('API Client Interceptor', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockLocation.pathname = '/'
    mockLocation.href = ''
    localStorage.clear()
    vi.clearAllMocks()
  })

  it('redirects to login and clears user data on 401, without adding notification', async () => {
    localStorage.setItem('auth_user', 'testuser')

    // Create a mock error response
    const error = new Error('Unauthorized') as MockAxiosError
    error.isAxiosError = true
    error.response = { status: 401 }

    // Get the interceptor function
    const rejectInterceptor = (client.interceptors.response as unknown as InterceptorHandlers).handlers[0]!.rejected

    // Spy on notification store
    const notificationStore = useNotificationStore()
    const addNotificationSpy = vi.spyOn(notificationStore, 'addNotification')

    // Expect it to reject with the error
    await expect(rejectInterceptor(error)).rejects.toThrow('Unauthorized')

    // Should have cleared storage and redirected
    expect(localStorage.getItem('auth_user')).toBeNull()
    expect(mockLocation.href).toBe('/login')

    // IMPORTANT: Should NOT have added a notification
    expect(addNotificationSpy).not.toHaveBeenCalled()
  })

  it('adds an error notification for non-401 errors', async () => {
    const error = new Error('Server Error') as MockAxiosError
    error.isAxiosError = true
    error.response = { status: 500, data: 'Internal Server Error' }

    const rejectInterceptor = (client.interceptors.response as unknown as InterceptorHandlers).handlers[0]!.rejected

    const notificationStore = useNotificationStore()
    const addNotificationSpy = vi.spyOn(notificationStore, 'addNotification')

    await expect(rejectInterceptor(error)).rejects.toThrow('Server Error')

    // Should have added a notification
    expect(addNotificationSpy).toHaveBeenCalledWith('Internal Server Error', 'error')
  })

  it('does not add notification if error was manually cancelled', async () => {
    const error = new Error('Cancelled') as MockAxiosError
    error.isAxiosError = true
    error.response = { status: 500 }

    // Mock axios.isCancel to return true
    vi.spyOn(axios, 'isCancel').mockReturnValue(true)

    const rejectInterceptor = (client.interceptors.response as unknown as InterceptorHandlers).handlers[0]!.rejected

    const notificationStore = useNotificationStore()
    const addNotificationSpy = vi.spyOn(notificationStore, 'addNotification')

    await expect(rejectInterceptor(error)).rejects.toThrow('Cancelled')

    // Should NOT have added a notification
    expect(addNotificationSpy).not.toHaveBeenCalled()
  })
})
