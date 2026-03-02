<template>
  <div v-if="error" class="error-state" role="alert">
    <div class="content">
      <span class="icon">⚠️</span>
      <span class="text">{{ error }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  error: string | null
}>()
</script>

<style scoped lang="scss">
@use "sass:color";
@use '@/assets/main.scss' as *;

.error-state {
  background-color: rgba($danger-color, 0.1);
  color: color.adjust($danger-color, $lightness: 10%);
  padding: 0.75rem 1rem;
  border-radius: 4px;
  border: 1px solid rgba($danger-color, 0.2);
  margin-bottom: 1rem;
  display: flex;
  align-items: center;
  justify-content: flex-start;
  gap: 0.75rem;

  .content {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .icon {
    font-size: 1.1rem;
    flex-shrink: 0;
  }

  .text {
    font-size: 0.9rem;
    font-weight: 500;
  }
}
</style>