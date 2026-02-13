<template>
  <div class="container">
    <h1>Select Project</h1>
    <div v-if="store.loading" class="loading">Loading projects...</div>
    <div v-else-if="store.error" class="error">{{ store.error }}</div>
    <div v-else class="project-list">
      <div v-if="store.projects.length === 0">No projects found.</div>
      <div v-for="project in store.projects" :key="project.id" class="project-item">
        <router-link :to="`/projects/${project.id}/registries`" class="project-link">
          <div class="project-name">{{ project.name }}</div>
          <div class="project-id">{{ project.id }}</div>
        </router-link>
      </div>
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

<style scoped>
.container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}
.project-list {
  display: grid;
  gap: 10px;
}
.project-item {
  border: 1px solid #ddd;
  padding: 15px;
  border-radius: 4px;
}
.project-link {
  text-decoration: none;
  color: inherit;
  display: block;
}
.project-link:hover {
  background-color: #f9f9f9;
}
.project-name {
  font-weight: bold;
  font-size: 1.1em;
}
.project-id {
  color: #666;
  font-size: 0.9em;
}
.error {
  color: red;
}
</style>
