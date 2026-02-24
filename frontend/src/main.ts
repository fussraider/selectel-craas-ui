import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import { useConfigStore } from '@/stores/config'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)

const configStore = useConfigStore()
await configStore.fetchConfig()

app.use(router)

app.mount('#app')
