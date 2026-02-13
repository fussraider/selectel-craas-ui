import { createRouter, createWebHistory } from 'vue-router'
import ProjectSelect from '../views/ProjectSelect.vue'
import RegistryList from '../views/RegistryList.vue'
import RepoList from '../views/RepoList.vue'
import ImageList from '../views/ImageList.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: ProjectSelect,
    },
    {
      path: '/projects/:pid/registries',
      name: 'registries',
      component: RegistryList,
    },
    {
      path: '/projects/:pid/registries/:rid/repositories',
      name: 'repositories',
      component: RepoList,
    },
    {
      path: '/projects/:pid/registries/:rid/repositories/:rname/images',
      name: 'images',
      component: ImageList,
    },
  ],
})

export default router
