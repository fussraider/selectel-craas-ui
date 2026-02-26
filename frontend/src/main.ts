import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import { useConfigStore } from '@/stores/config'
import { useAuthStore } from '@/stores/auth'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)

const configStore = useConfigStore()
await configStore.fetchConfig()

if (configStore.authEnabled) {
  const authStore = useAuthStore()
  try {
    await authStore.checkAuth()
  } catch {
    // ignore
  }
}

app.use(router)

app.mount('#app')
