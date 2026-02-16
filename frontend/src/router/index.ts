import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import RegistrySettings from '../views/RegistrySettings.vue'
import ImageList from '../views/ImageList.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
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
    // Keep the old route as redirect or alias if needed, but for now we simplify
    {
      path: '/projects/:pid/registries/:rid/repositories/:rname/images',
      redirect: to => {
        return { name: 'images', params: to.params }
      }
    }
  ],
})

export default router
