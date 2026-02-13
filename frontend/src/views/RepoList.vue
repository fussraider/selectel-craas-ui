<template>
  <div class="container">
    <div class="header">
      <router-link to="/">Projects</router-link> &gt;
      <router-link :to="`/projects/${pid}/registries`">Registries</router-link> &gt;
      Repositories
      <h1>Repositories in Registry {{ rid }}</h1>
    </div>
    <div v-if="store.loading" class="loading">Loading...</div>
    <div v-else-if="store.error" class="error">{{ store.error }}</div>
    <div v-else class="repo-list">
      <div v-if="store.repositories.length === 0">No repositories found.</div>
      <div v-for="repo in store.repositories" :key="repo.name" class="repo-item">
        <div class="repo-info">
          <router-link :to="`/projects/${pid}/registries/${rid}/repositories/${encodeURIComponent(repo.name)}/images`" class="repo-link">
            {{ repo.name }}
          </router-link>
          <div class="repo-details">Size: {{ repo.size }} bytes | Updated: {{ repo.updatedAt }}</div>
        </div>
        <button @click="deleteRepo(repo.name)" class="delete-btn">Delete</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRegistryStore } from '@/stores/registry'
import { onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const store = useRegistryStore()
const pid = computed(() => route.params.pid as string)
const rid = computed(() => route.params.rid as string)

onMounted(() => {
  store.fetchRepositories(pid.value, rid.value)
})

const deleteRepo = async (rname: string) => {
  if (confirm(`Are you sure you want to delete repository '${rname}'?`)) {
    await store.deleteRepository(pid.value, rid.value, rname)
  }
}
</script>

<style scoped>
.container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}
.repo-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.repo-item {
  border: 1px solid #ddd;
  padding: 15px;
  border-radius: 4px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.repo-link {
  font-weight: bold;
  font-size: 1.1em;
  text-decoration: none;
  color: #007bff;
}
.repo-details {
  font-size: 0.9em;
  color: #666;
}
.delete-btn {
  background-color: #dc3545;
  color: white;
  border: none;
  padding: 5px 10px;
  border-radius: 3px;
  cursor: pointer;
}
.delete-btn:hover {
  background-color: #c82333;
}
.error {
  color: red;
}
</style>
