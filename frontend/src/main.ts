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

// We must catch errors, but we also want the app to mount immediately so it can show them.
// By not awaiting these sequentially before mount, the UI can render its loading states
// and toast containers can catch the errors fired.

app.use(router)
app.mount('#app')

// Now run the fetching so that any CustomEvent dispatched by axios is caught by ToastContainer
configStore.fetchConfig().finally(() => {
  if (configStore.authEnabled) {
    const authStore = useAuthStore()
    authStore.checkAuth().catch(() => {
      // ignore
    })
  }
})