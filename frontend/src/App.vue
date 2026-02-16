<script setup lang="ts">
import { RouterView } from 'vue-router'
import AppLayout from '@/components/AppLayout.vue'
</script>

<template>
  <AppLayout>
    <RouterView :key="$route.fullPath" />
  </AppLayout>
</template>

<style lang="scss">
@use '@/assets/main.scss';
</style>
