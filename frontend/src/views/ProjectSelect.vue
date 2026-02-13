<template>
  <div class="view-container">
    <h1>Select Project</h1>
    <div v-if="store.loading" class="loading">Loading projects...</div>
    <div v-else-if="store.error" class="error">{{ store.error }}</div>
    <div v-else class="project-list">
      <div v-if="store.projects.length === 0">No projects found.</div>
      <router-link v-else v-for="project in store.projects" :key="project.id" :to="`/projects/${project.id}/registries`" class="project-card">
        <div class="project-name">{{ project.name }}</div>
        <div class="project-id">{{ project.id }}</div>
      </router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRegistryStore } from '@/stores/registry'
import { onMounted } from 'vue'

const store = useRegistryStore()

onMounted(() => {
  store.fetchProjects()
})
</script>

<style scoped lang="scss">
@use "sass:color";
@use '@/assets/main.scss' as *;

.view-container {
  h1 {
    margin-bottom: 2rem;
    color: $primary-color;
  }
}

.project-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 1.5rem;
}

.project-card {
  background: $card-bg;
  border: 1px solid $border-color;
  border-radius: 12px;
  padding: 2rem;
  transition: all 0.2s ease-in-out;
  text-decoration: none;
  color: inherit;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: flex-start;

  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
    border-color: $primary-color;
    background-color: color.adjust($card-bg, $lightness: 3%);
    text-decoration: none;
  }

  .project-name {
    font-size: 1.25rem;
    font-weight: 600;
    margin-bottom: 0.5rem;
    color: $text-color;
  }

  .project-id {
    font-size: 0.875rem;
    color: $secondary-color;
    font-family: monospace;
    background: $muted-bg;
    color: $text-color;
    padding: 0.2rem 0.4rem;
    border-radius: 4px;
    word-break: break-all;
    max-width: 100%;
  }
}
</style>
