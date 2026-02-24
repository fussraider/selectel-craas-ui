<script setup lang="ts">
import { RouterView, useRoute } from 'vue-router'
import AppLayout from '@/components/AppLayout.vue'

const route = useRoute()
</script>

<template>
  <component :is="route.meta.hideLayout ? 'div' : AppLayout">
    <RouterView :key="$route.fullPath" />
  </component>
</template>

<style lang="scss">
@use '@/assets/main.scss';
</style>
