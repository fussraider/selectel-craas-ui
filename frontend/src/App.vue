<script setup lang="ts">
import { onMounted } from 'vue'
import { RouterView } from 'vue-router'
import AppLayout from '@/components/AppLayout.vue'
import { useConfigStore } from '@/stores/config'

const configStore = useConfigStore()

onMounted(() => {
    configStore.fetchConfig()
})
</script>

<template>
  <AppLayout>
    <RouterView :key="$route.fullPath" />
  </AppLayout>
</template>

<style lang="scss">
@use '@/assets/main.scss';
</style>
