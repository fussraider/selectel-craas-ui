<template>
  <div class="view-container">
    <div class="header">
      <h1>Repositories ({{ rid }})</h1>
    </div>

    <div class="gc-section" v-if="store.gcInfo">
      <div class="gc-header">
          <h3>Garbage Collection</h3>
          <span class="gc-size">Potential Savings: {{ (store.gcInfo.sizeSummary / 1024 / 1024).toFixed(2) }} MB</span>
      </div>
      <button @click="triggerGC" :disabled="store.gcLoading" class="gc-btn">Start Garbage Collection</button>
    </div>

    <div v-if="store.loading" class="loading">Loading repositories...</div>
    <div v-else-if="store.error" class="error">{{ store.error }}</div>
    <div v-else class="list-container">
      <div v-if="store.repositories.length === 0" class="empty-state">No repositories found.</div>
      <div v-for="repo in store.repositories" :key="repo.name" class="list-item">
        <div class="item-info">
          <router-link :to="`/projects/${pid}/registries/${rid}/repositories/${encodeURIComponent(repo.name)}/images`" class="item-link">
            {{ repo.name }}
          </router-link>
          <div class="item-meta">
            <span>Size: {{ (repo.size / 1024 / 1024).toFixed(2) }} MB</span>
            <span>Updated: {{ new Date(repo.updatedAt).toLocaleDateString() }}</span>
          </div>
        </div>
        <button @click="deleteRepo(repo.name)" class="delete-btn" title="Delete Repository">
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
const rid = computed(() => route.params.rid as string)

onMounted(() => {
  store.fetchRepositories(pid.value, rid.value)
  store.fetchGCInfo(pid.value, rid.value)
})

const triggerGC = async () => {
    if (confirm("Garbage collection makes the registry read-only until it completes. Are you sure you want to proceed?")) {
        await store.startGC(pid.value, rid.value)
    }
}

const deleteRepo = async (rname: string) => {
  if (confirm(`Are you sure you want to delete repository '${rname}'?`)) {
    await store.deleteRepository(pid.value, rid.value, rname)
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

.gc-section {
  background: $card-bg;
  border: 1px solid $border-color;
  border-radius: 6px;
  padding: 1.25rem;
  margin-bottom: 2rem;
  display: flex;
  justify-content: space-between;
  align-items: center;

  .gc-header {
      display: flex;
      flex-direction: column;
      gap: 0.5rem;

      h3 {
          font-size: 1.1rem;
          margin: 0;
          color: $primary-color;
      }

      .gc-size {
          font-size: 0.9rem;
          color: $secondary-color;
      }
  }
}

.gc-btn {
    background-color: $primary-color;
    color: white;
    border: none;
    padding: 0.6rem 1.2rem;
    border-radius: 4px;
    cursor: pointer;
    font-weight: 500;
    transition: opacity 0.2s;

    &:hover:not(:disabled) {
        opacity: 0.9;
    }

    &:disabled {
        background-color: #ccc;
        cursor: not-allowed;
    }
}

.list-container {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.list-item {
  background: $card-bg;
  border: 1px solid $border-color;
  border-radius: 6px;
  padding: 1.25rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
  transition: box-shadow 0.2s;

  &:hover {
    box-shadow: 0 2px 4px rgba(0,0,0,0.2);
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
