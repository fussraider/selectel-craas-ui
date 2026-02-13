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
  background: white;
  border: 1px solid $border-color;
  border-radius: 8px;
  padding: 1.5rem;
  transition: transform 0.2s, box-shadow 0.2s, border-color 0.2s;
  text-decoration: none;
  color: inherit;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: flex-start;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 6px rgba(0,0,0,0.1);
    border-color: $primary-color;
    background-color: #fff;
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
    background: #f1f1f1;
    padding: 0.2rem 0.4rem;
    border-radius: 4px;
  }
}
</style>
