<template>
  <div class="app-layout">
    <header class="app-header">
      <div class="header-left">
        <button class="mobile-toggle" @click="toggleSidebar">
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <line x1="3" y1="12" x2="21" y2="12"></line>
            <line x1="3" y1="6" x2="21" y2="6"></line>
            <line x1="3" y1="18" x2="21" y2="18"></line>
          </svg>
        </button>
        <div class="logo">
          <router-link to="/">
            <img src="/logo.png" alt="CRaaS Logo" class="logo-img" />
            <span class="logo-text">CRaaS Console</span>
          </router-link>
        </div>
      </div>

      <div class="header-center">
        <!-- Breadcrumbs moved here -->
        <div class="breadcrumbs" v-if="breadcrumbs.length > 0">
           <span v-for="(crumb, index) in breadcrumbs" :key="index">
            <router-link v-if="crumb.to" :to="crumb.to">{{ crumb.label }}</router-link>
            <span v-else>{{ crumb.label }}</span>
            <span v-if="index < breadcrumbs.length - 1" class="separator">/</span>
          </span>
        </div>
        <div id="header-actions"></div>
      </div>

      <div class="header-right">
        <a href="https://github.com/fussraider/selectel-craas-ui" target="_blank" rel="noopener noreferrer" class="github-link" title="GitHub Project">
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M9 19c-5 1.5-5-2.5-7-3m14 6v-3.87a3.37 3.37 0 0 0-.94-2.61c3.14-.35 6.44-1.54 6.44-7A5.44 5.44 0 0 0 20 4.77 5.07 5.07 0 0 0 19.91 1S18.73.65 16 2.48a13.38 13.38 0 0 0-7 0C6.27.65 5.09 1 5.09 1A5.07 5.07 0 0 0 5 4.77a5.44 5.44 0 0 0-1.5 3.78c0 5.42 3.3 6.61 6.44 7A3.37 3.37 0 0 0 9 18.13V22"></path>
          </svg>
        </a>
      </div>
    </header>

    <div class="app-body">
      <aside class="sidebar-container" :class="{ open: sidebarOpen }">
        <RepositorySidebar />
      </aside>

      <main class="app-content" @click="closeSidebarIfMobile">
        <slot />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRoute } from 'vue-router'
import RepositorySidebar from './RepositorySidebar.vue'
import { useRegistryStore } from '@/stores/registry'
import { onMounted } from 'vue'

const store = useRegistryStore()
const route = useRoute()

onMounted(async () => {
    // Ensure project is loaded since selector is removed
    await store.fetchProjects()
    // Logic to select project moved to store/here if needed, but store handles auto-select on fetchProjects
    if (store.projects.length > 0 && !store.selectedProjectId) {
        const firstProject = store.projects[0]
        if (firstProject) {
            store.selectedProjectId = firstProject.id
            await store.loadProjectData(store.selectedProjectId)
        }
    } else if (store.selectedProjectId) {
        // Refresh data if ID exists
        await store.loadProjectData(store.selectedProjectId)
    }
})
const sidebarOpen = ref(false)

const toggleSidebar = () => {
  sidebarOpen.value = !sidebarOpen.value
}

const closeSidebarIfMobile = () => {
    if (window.innerWidth < 768) {
        sidebarOpen.value = false
    }
}

// Close sidebar on route change (mobile)
watch(() => route.fullPath, () => {
    if (window.innerWidth < 768) {
        sidebarOpen.value = false
    }
})

const breadcrumbs = computed(() => {
  const crumbs = []

  const pid = route.params.pid as string
  const rid = route.params.rid as string
  const rname = route.params.rname as string

  if (rid) {
      crumbs.push({ label: `Registry ${rid}`, to: `/projects/${pid}/registries/${rid}` })
  }

  if (rname) {
      crumbs.push({ label: `Repo ${rname}`, to: `/projects/${pid}/registries/${rid}/repositories/${encodeURIComponent(rname)}` })
  }

  return crumbs
})
</script>

<style scoped lang="scss">
@use "sass:color";
@use '@/assets/main.scss' as *;

.app-layout {
  display: flex;
  flex-direction: column;
  height: 100vh;
  overflow: hidden;
}

.app-header {
  background-color: $card-bg;
  border-bottom: 1px solid $border-color;
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 1rem;
  flex-shrink: 0;
  z-index: 50;

  .header-left {
      display: flex;
      align-items: center;
      gap: 1rem;
  }

  .logo a {
    font-size: 1.25rem;
    font-weight: 700;
    text-decoration: none;
    color: $primary-color;
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }

  .logo-img {
    height: 50px;
    width: auto;
    transition: height 0.2s;
  }

  @media (max-width: 768px) {
      padding: 0 0.5rem;

      .logo-text {
          display: none;
      }

      .logo-img {
          height: 32px;
      }
  }
}

.mobile-toggle {
    display: none;
    background: none;
    border: none;
    cursor: pointer;
    color: $text-color;

    @media (max-width: 768px) {
        display: block;
    }
}

.app-body {
    display: flex;
    flex: 1;
    overflow: hidden;
    position: relative;
}

.sidebar-container {
    width: 300px;
    height: 100%;
    border-right: 1px solid $border-color;
    background: $card-bg;
    transition: transform 0.3s ease;

    @media (max-width: 768px) {
        position: absolute;
        top: 0;
        left: 0;
        bottom: 0;
        z-index: 40;
        transform: translateX(-100%);
        box-shadow: 2px 0 8px rgba(0,0,0,0.1);

        &.open {
            transform: translateX(0);
        }
    }
}

.app-content {
    flex: 1;
    overflow-y: auto;
    padding: 2rem;
    background-color: $background-color;

    @media (max-width: 768px) {
        padding: 1rem;
    }
}

.breadcrumbs {
  font-size: 0.9rem;
  color: $secondary-color;
  display: flex;
  align-items: center;
  white-space: nowrap;

  a {
      color: $primary-color;
      text-decoration: none;
      padding: 0.2rem 0.4rem;
      border-radius: 4px;
      transition: background-color 0.2s;

      &:hover {
          background-color: rgba($primary-color, 0.1);
      }
  }

  .separator {
      margin: 0 0.2rem;
      color: $border-color;
  }
}

.header-center {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: space-between;
    overflow: hidden;
    padding: 0 1rem;
    gap: 1rem;
}

#header-actions {
    display: flex;
    align-items: center;
    margin-left: auto;
}

@media (max-width: 768px) {
    .breadcrumbs {
        display: none;
    }
}

.github-link {
    color: $secondary-color;
    transition: color 0.2s;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 0.5rem;
    border-radius: 50%;

    &:hover {
        color: $text-color;
        background-color: rgba($text-color, 0.1);
    }
}
</style>
