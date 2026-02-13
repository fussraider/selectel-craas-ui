<template>
  <div class="container">
    <div class="header">
      <router-link to="/">Projects</router-link> &gt; Registries
      <h1>Registries for Project {{ pid }}</h1>
    </div>
    <div v-if="store.loading" class="loading">Loading...</div>
    <div v-else-if="store.error" class="error">{{ store.error }}</div>
    <div v-else class="registry-list">
      <div v-if="store.registries.length === 0">No registries found.</div>
      <div v-for="registry in store.registries" :key="registry.id" class="registry-item">
        <div class="registry-info">
          <router-link :to="`/projects/${pid}/registries/${registry.id}/repositories`" class="registry-link">
            {{ registry.name }}
          </router-link>
          <div class="registry-details">
            Status: {{ registry.status }} | Created: {{ registry.createdAt }}
          </div>
        </div>
        <button @click="deleteReg(registry.id)" class="delete-btn">Delete</button>
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

onMounted(() => {
  store.fetchRegistries(pid.value)
})

const deleteReg = async (rid: string) => {
  if (confirm('Are you sure you want to delete this registry? This action is irreversible.')) {
    await store.deleteRegistry(pid.value, rid)
  }
}
</script>

<style scoped>
.container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}
.registry-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.registry-item {
  border: 1px solid #ddd;
  padding: 15px;
  border-radius: 4px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.registry-link {
  font-weight: bold;
  font-size: 1.1em;
  text-decoration: none;
  color: #007bff;
}
.registry-details {
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
