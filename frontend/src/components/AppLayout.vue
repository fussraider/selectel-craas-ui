<template>
  <div class="app-layout">
    <header class="app-header">
      <div class="logo">
        <router-link to="/">CRaaS Console</router-link>
      </div>
    </header>
    <main class="app-content">
      <div class="breadcrumbs" v-if="breadcrumbs.length">
        <span v-for="(crumb, index) in breadcrumbs" :key="index">
          <router-link v-if="crumb.to" :to="crumb.to">{{ crumb.label }}</router-link>
          <span v-else>{{ crumb.label }}</span>
          <span v-if="index < breadcrumbs.length - 1" class="separator">/</span>
        </span>
      </div>
      <slot />
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

const breadcrumbs = computed(() => {
  const crumbs = []
  crumbs.push({ label: 'Projects', to: '/' })

  const pid = route.params.pid as string
  const rid = route.params.rid as string
  const rname = route.params.rname as string

  if (pid) {
    crumbs.push({ label: `Project ${pid}`, to: `/projects/${pid}/registries` })
  }
  if (rid) {
     crumbs.push({ label: `Registry ${rid}`, to: `/projects/${pid}/registries/${rid}/repositories` })
  }
  if (rname) {
      const encodedRname = encodeURIComponent(rname)
      crumbs.push({ label: `Repo ${rname}`, to: `/projects/${pid}/registries/${rid}/repositories/${encodedRname}/images` })
  }

  // If we are on the current page, we might want to disable the link, but Vue Router handles active links well.
  // We can remove the 'to' from the last crumb if desired.
  if (crumbs.length > 0) {
      const last = crumbs[crumbs.length - 1]
      // Check if last crumb matches current path?
      // Simplified: just let it be a link to self.
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
  min-height: 100vh;
}

.app-header {
  background-color: rgba($card-bg, 0.95);
  backdrop-filter: blur(10px);
  color: $text-color;
  padding: 1rem 2rem;
  display: flex;
  align-items: center;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
  border-bottom: 1px solid $border-color;
  position: sticky;
  top: 0;
  z-index: 50;

  .logo a {
    color: $primary-color;
    font-size: 1.5rem;
    font-weight: 700;
    text-decoration: none;
    background: linear-gradient(135deg, $primary-color, color.adjust($primary-color, $lightness: 20%));
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;

    &:hover {
        opacity: 0.8;
    }
  }
}

.app-content {
  flex: 1;
  padding: 2rem;
  max-width: 1200px;
  width: 100%;
  margin: 0 auto;
}

.breadcrumbs {
  margin-bottom: 1.5rem;
  font-size: 0.95rem;
  color: $secondary-color;
  background: $card-bg;
  padding: 0.75rem 1rem;
  border-radius: 4px;
  border: 1px solid $border-color;

  a {
    color: #6ea8fe;
    font-weight: 500;
  }

  .separator {
    margin: 0 0.5rem;
    color: $secondary-color;
  }
}
</style>
