<script setup lang="ts">
import { RouterView, useRoute } from 'vue-router'
import AppLayout from '@/components/AppLayout.vue'
import ToastContainer from '@/components/ToastContainer.vue'

const route = useRoute()
</script>

<template>
  <component :is="route.meta.hideLayout ? 'div' : AppLayout">
    <RouterView :key="$route.fullPath" />
  </component>
  <ToastContainer />
</template>

<style lang="scss">
@use '@/assets/main.scss';
</style>
