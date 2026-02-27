<template>
  <div class="home-view">
    <div class="welcome-content">
      <h1>Welcome to CRaaS Console</h1>
      <p>Select a repository from the sidebar to view images.</p>

      <div v-if="store.error" class="error-container">
          <p class="error-text">Failed to load projects.</p>
          <button @click="store.fetchProjects" class="btn small-btn">Retry</button>
      </div>

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

.error-container {
    margin-top: 1rem;
    padding: 1rem;
    background-color: rgba($danger-color, 0.1);
    border: 1px solid rgba($danger-color, 0.3);
    border-radius: 6px;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.5rem;

    .error-text {
        color: $danger-color;
        font-weight: bold;
        font-size: 1rem;
    }

    .btn {
        background-color: $card-bg;
        border: 1px solid $border-color;
        color: $text-color;
        padding: 0.3rem 0.8rem;
        border-radius: 4px;
        cursor: pointer;

        &:hover {
            background-color: color.adjust($card-bg, $lightness: 5%);
        }
    }
}
</style>
