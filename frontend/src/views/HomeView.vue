<template>
  <div class="home-view">
    <div class="welcome-content">
      <h1>Welcome to CRaaS Console</h1>
      <p>Select a repository from the sidebar to view images.</p>
      <p class="hint" v-if="!store.selectedProjectId">Loading projects...</p>
      <p class="hint" v-else>Current Project: {{ store.selectedProjectId }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRegistryStore } from '@/stores/registry'
const store = useRegistryStore()
</script>

<style scoped lang="scss">
@use '@/assets/main.scss' as *;

.home-view {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
  text-align: center;
}

.welcome-content {
    padding: 2rem;
    background: $card-bg;
    border: 1px solid $border-color;
    border-radius: 8px;

    h1 { color: $primary-color; margin-bottom: 1rem; }
    p { color: $text-color; font-size: 1.1rem; }
    .hint { color: $secondary-color; font-size: 0.9rem; margin-top: 1rem; }
}
</style>
