import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import RegistrySettings from '../views/RegistrySettings.vue'
import ImageList from '../views/ImageList.vue'
import LoginView from '../views/LoginView.vue'
import { useConfigStore } from '@/stores/config'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: LoginView,
      meta: { hideLayout: true }
    },
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/projects/:pid/registries/:rid',
      name: 'registry-settings',
      component: RegistrySettings,
    },
    {
      path: '/projects/:pid/registries/:rid/repositories/:rname',
      name: 'images',
      component: ImageList,
    },
    {
      path: '/projects/:pid/registries/:rid/repositories/:rname/images',
      redirect: to => {
        return { name: 'images', params: to.params }
      }
    }
  ],
})

router.beforeEach((to, from, next) => {
  const configStore = useConfigStore()
  const authStore = useAuthStore()

  if (configStore.authEnabled) {
    if (!authStore.isAuthenticated && to.name !== 'login') {
      next({ name: 'login' })
      return
    }
    if (authStore.isAuthenticated && to.name === 'login') {
      next({ name: 'home' })
      return
    }
  }
  next()
})

export default router
