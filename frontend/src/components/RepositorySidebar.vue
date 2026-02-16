<template>
  <aside class="sidebar">
    <div class="sidebar-header">
      <h2>Repositories</h2>
      <button @click="refresh" :disabled="store.loading" class="refresh-btn" title="Refresh">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-arrow-clockwise" viewBox="0 0 16 16" :class="{ spinning: store.loading }">
          <path fill-rule="evenodd" d="M8 3a5 5 0 1 0 4.546 2.914.5.5 0 0 1 .908-.417A6 6 0 1 1 8 2v1z"/>
          <path d="M8 4.466V.534a.25.25 0 0 1 .41-.192l2.36 1.966c.12.1.12.284 0 .384L8.41 4.658A.25.25 0 0 1 8 4.466z"/>
        </svg>
      </button>
    </div>

    <div v-if="store.error" class="error-msg">{{ store.error }}</div>

    <div class="repo-list">
      <div v-for="registry in store.registries" :key="registry.id" class="registry-group">
        <div class="registry-header">
            <router-link
                :to="`/projects/${store.selectedProjectId}/registries/${registry.id}`"
                class="registry-link"
                active-class="active"
            >
                <span class="registry-icon">📦</span>
                {{ registry.name }}
            </router-link>
            <span v-if="registry.loadingRepos" class="loading-indicator">...</span>
        </div>

        <div class="repo-items">
            <div v-if="!registry.repositories || registry.repositories.length === 0" class="empty-repos">
                No repositories
            </div>
            <router-link
                v-else
                v-for="repo in registry.repositories"
                :key="repo.name"
                :to="`/projects/${store.selectedProjectId}/registries/${registry.id}/repositories/${encodeURIComponent(repo.name)}`"
                class="repo-link"
                active-class="active"
            >
                <span class="repo-icon">📄</span>
                {{ repo.name }}
            </router-link>
        </div>
      </div>

      <div v-if="store.registries.length === 0 && !store.loading" class="empty-state">
          No registries found.
      </div>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { useRegistryStore } from '@/stores/registry'

const store = useRegistryStore()

const refresh = async () => {
    await store.refreshStructure()
}
</script>

<style scoped lang="scss">
@use '@/assets/main.scss' as *;

.sidebar {
  width: 300px;
  background-color: $card-bg;
  border-right: 1px solid $border-color;
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
}

.sidebar-header {
  padding: 1rem;
  border-bottom: 1px solid $border-color;
  display: flex;
  justify-content: space-between;
  align-items: center;

  h2 {
      font-size: 1.1rem;
      margin: 0;
      color: $text-color;
      font-weight: 600;
  }
}

.refresh-btn {
    background: transparent;
    border: none;
    color: $secondary-color;
    cursor: pointer;
    padding: 0.5rem;
    border-radius: 50%;
    transition: background-color 0.2s, color 0.2s;
    display: flex;
    align-items: center;
    justify-content: center;

    &:hover:not(:disabled) {
        background-color: rgba($text-color, 0.05);
        color: $primary-color;
    }

    &:disabled {
        opacity: 0.5;
        cursor: not-allowed;
    }
}

.spinning {
    animation: spin 1s linear infinite;
}

@keyframes spin {
    from { transform: rotate(0deg); }
    to { transform: rotate(360deg); }
}

.error-msg {
    padding: 1rem;
    color: $danger-color;
    font-size: 0.85rem;
    background-color: rgba($danger-color, 0.1);
}

.repo-list {
    flex: 1;
    overflow-y: auto;
    padding: 1rem 0;
}

.registry-group {
    margin-bottom: 1rem;
}

.registry-header {
    padding: 0.5rem 1rem;
    font-size: 0.95rem;
    font-weight: 600;
    color: $text-color;
    display: flex;
    align-items: center;
    justify-content: space-between;
}

.registry-link {
    color: inherit;
    text-decoration: none;
    display: flex;
    align-items: center;
    gap: 0.5rem;
    flex: 1;

    &:hover {
        color: $primary-color;
    }

    &.active {
        color: $primary-color;
    }
}

.repo-items {
    padding-left: 0.5rem;
}

.repo-link {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem 1rem 0.5rem 2rem; // Indent
    font-size: 0.9rem;
    color: $secondary-color;
    text-decoration: none;
    border-left: 2px solid transparent;
    transition: all 0.2s;

    &:hover {
        background-color: rgba($primary-color, 0.05);
        color: $text-color;
    }

    &.active {
        background-color: rgba($primary-color, 0.1);
        color: $primary-color;
        border-left-color: $primary-color;
        font-weight: 500;
    }
}

.registry-icon, .repo-icon {
    opacity: 0.7;
}

.empty-repos {
    padding: 0.5rem 1rem 0.5rem 2.5rem;
    font-size: 0.85rem;
    color: $secondary-color;
    font-style: italic;
}

.empty-state {
    padding: 2rem;
    text-align: center;
    color: $secondary-color;
}

.loading-indicator {
    font-size: 0.8rem;
    color: $secondary-color;
}
</style>
