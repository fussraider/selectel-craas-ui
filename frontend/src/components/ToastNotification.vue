<template>
  <div v-if="message" :class="['toast', type]" role="alert">
    <div class="content">
      <span class="icon">{{ type === 'success' ? '✅' : '⚠️' }}</span>
      <span class="text">{{ message }}</span>
    </div>
    <button class="close-btn" @click="$emit('close')">&times;</button>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  message: string | null
  type: 'success' | 'error' | 'info'
}>()

defineEmits<{
  (e: 'close'): void
}>()
</script>

<style scoped lang="scss">
@use '@/assets/main.scss' as *;

.toast {
  /* Removed fixed positioning, as ToastContainer handles layout */
  background-color: $card-bg;
  color: $text-color;
  padding: 1rem 1.5rem;
  border-radius: 8px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.5), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
  z-index: 50;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  min-width: 300px;
  max-width: 90vw;
  border: 1px solid $border-color;
  pointer-events: auto; /* Ensure clicks work since container ignores them */

  &.success {
    border-left: 4px solid #10b981; // Green-500
  }

  &.error {
    border-left: 4px solid $danger-color;
  }

  &.info {
    border-left: 4px solid $primary-color;
  }
}

.content {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.icon {
  font-size: 1.25rem;
}

.text {
  font-size: 0.95rem;
  font-weight: 500;
}

.close-btn {
  background: transparent;
  border: none;
  color: $secondary-color;
  font-size: 1.5rem;
  line-height: 1;
  padding: 0;
  cursor: pointer;
  transition: color 0.2s;

  &:hover {
    color: $text-color;
  }
}
</style>
