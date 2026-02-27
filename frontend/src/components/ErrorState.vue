<template>
  <div class="error-state">
    <p>{{ title }}</p>
    <button v-if="retry" @click="retry" class="btn small-btn">Retry</button>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  title?: string
  retry?: () => void
}>()
</script>

<style scoped lang="scss">
@use "sass:color";
@use '@/assets/main.scss' as *;

.error-state {
    padding: 1rem;
    color: $danger-color;
    border: 1px solid rgba($danger-color, 0.3);
    background-color: rgba($danger-color, 0.05);
    border-radius: 6px;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1rem;
    text-align: center;
    margin: 1rem 0;

    p {
        margin: 0;
        font-weight: bold;
        font-size: 1rem;
    }

    button {
        background-color: $card-bg;
        border: 1px solid $border-color;
        color: $text-color;
        padding: 0.3rem 0.8rem;
        border-radius: 4px;
        cursor: pointer;
        font-weight: 500;
        transition: background-color 0.2s;

        &:hover {
            background-color: color.adjust($card-bg, $lightness: 5%);
        }
    }
}
</style>
