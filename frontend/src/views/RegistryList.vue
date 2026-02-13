<template>
  <div class="view-container">
    <div class="header">
      <h1>Registries (Project {{ pid }})</h1>
    </div>
    <div v-if="store.loading" class="loading">Loading registries...</div>
    <div v-else-if="store.error" class="error">{{ store.error }}</div>
    <div v-else class="list-container">
      <div v-if="store.registries.length === 0" class="empty-state">No registries found.</div>
      <div v-for="registry in store.registries" :key="registry.id" class="list-item">
        <div class="item-info">
          <router-link :to="`/projects/${pid}/registries/${registry.id}/repositories`" class="item-link">
            {{ registry.name }}
          </router-link>
          <div class="item-meta">
            <span class="status" :class="registry.status.toLowerCase()">{{ registry.status }}</span>
            <span class="created">Created: {{ new Date(registry.createdAt).toLocaleDateString() }}</span>
          </div>
        </div>
        <button @click="deleteReg(registry.id)" class="delete-btn" title="Delete Registry">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-trash" viewBox="0 0 16 16">
                <path d="M5.5 5.5A.5.5 0 0 1 6 6v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5zm2.5 0a.5.5 0 0 1 .5.5v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5zm3 .5a.5.5 0 0 0-1 0v6a.5.5 0 0 0 1 0V6z"/>
                <path fill-rule="evenodd" d="M14.5 3a1 1 0 0 1-1 1H13v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V4h-.5a1 1 0 0 1-1-1V2a1 1 0 0 1 1-1H6a1 1 0 0 1 1-1h2a1 1 0 0 1 1 1h3.5a1 1 0 0 1 1 1v1zM4.118 4 4 4.059V13a1 1 0 0 0 1 1h6a1 1 0 0 0 1-1V4.059L11.882 4H4.118zM2.5 3V2h11v1h-11z"/>
            </svg>
        </button>
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

<style scoped lang="scss">
@use '@/assets/main.scss' as *;

.view-container {
  h1 {
    margin-bottom: 2rem;
    color: $primary-color;
  }
}

.list-container {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.list-item {
  background: white;
  border: 1px solid $border-color;
  border-radius: 6px;
  padding: 1.25rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
  transition: box-shadow 0.2s;

  &:hover {
    box-shadow: 0 2px 4px rgba(0,0,0,0.05);
  }
}

.item-info {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
}

.item-link {
  font-weight: 600;
  font-size: 1.1rem;
}

.item-meta {
  font-size: 0.85rem;
  color: $secondary-color;
  display: flex;
  gap: 1rem;
  align-items: center;
}

.status {
    padding: 0.2rem 0.5rem;
    border-radius: 4px;
    font-size: 0.75rem;
    text-transform: uppercase;
    font-weight: bold;

    &.active {
        background-color: #d4edda;
        color: #155724;
    }

    &.creating {
        background-color: #fff3cd;
        color: #856404;
    }

    &.error {
        background-color: #f8d7da;
        color: #721c24;
    }
}

.delete-btn {
  background-color: transparent;
  color: $danger-color;
  border: 1px solid transparent;
  padding: 0.5rem;
  border-radius: 4px;
  transition: background-color 0.2s;

  &:hover {
    background-color: rgba(220, 53, 69, 0.1);
  }
}

.empty-state {
    text-align: center;
    padding: 3rem;
    color: $secondary-color;
    border: 1px dashed $border-color;
    border-radius: 8px;
}
</style>
