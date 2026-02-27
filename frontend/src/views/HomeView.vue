<template>
  <div class="home-view">
    <div class="welcome-content">
      <h1>Welcome to CRaaS Console</h1>
      <p>Select a repository from the sidebar to view images.</p>

      <ErrorState
        v-if="store.error"
        title="Failed to load projects."
        :retry="store.fetchProjects"
      />

      <p class="hint" v-else-if="store.loading">Loading projects...</p>
      <p class="hint" v-else-if="store.selectedProjectId">Current Project: {{ store.selectedProjectId }}</p>

      <!-- Toast Notifications for HomeView -->
      <ToastNotification
        v-if="store.error"
        type="error"
        :message="store.error"
        @close="store.clearNotifications"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRegistryStore } from '@/stores/registry'
import ToastNotification from '@/components/ToastNotification.vue'
import ErrorState from '@/components/ErrorState.vue'
const store = useRegistryStore()
</script>

<style scoped lang="scss">
@use "sass:color";
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
