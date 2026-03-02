<template>
  <div class="toast-container" aria-live="polite">
    <TransitionGroup name="toast-list">
      <ToastNotification
        v-for="notification in store.notifications"
        :key="notification.id"
        :message="notification.message"
        :type="notification.type"
        @close="store.removeNotification(notification.id)"
      />
    </TransitionGroup>
  </div>
</template>

<script setup lang="ts">
import { useNotificationStore } from '@/stores/notifications'
import ToastNotification from './ToastNotification.vue'

const store = useNotificationStore()
</script>

<style scoped lang="scss">
.toast-container {
  position: fixed;
  bottom: 2rem;
  right: 2rem;
  z-index: 1000;
  display: flex;
  flex-direction: column;
  gap: 1rem;
  pointer-events: none; /* Let clicks pass through the container */
  align-items: flex-end;
}

/* Make sure children restore pointer events */
.toast-container > * {
  pointer-events: auto;
}

.toast-list-enter-active,
.toast-list-leave-active {
  transition: all 0.3s ease;
}
.toast-list-enter-from {
  opacity: 0;
  transform: translateX(30px);
}
.toast-list-leave-to {
  opacity: 0;
  transform: translateY(-20px) scale(0.9);
}
</style>
