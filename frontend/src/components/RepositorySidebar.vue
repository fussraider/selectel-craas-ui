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

    <!-- Registry Loading Skeleton -->
    <div v-if="store.loading && store.registries.length === 0" class="repo-list">
        <div v-for="i in 3" :key="i" class="skeleton-group">
            <div class="skeleton-header"></div>
            <div class="skeleton-item" style="width: 80%"></div>
            <div class="skeleton-item" style="width: 60%"></div>
        </div>
    </div>

    <div v-else class="repo-list">
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
            <!-- Spinner for specific registry loading -->
            <div v-if="registry.loadingRepos" class="spinner-small"></div>
        </div>

        <div class="repo-items">
            <!-- Repos Loading Skeleton -->
            <div v-if="registry.loadingRepos && (!registry.repositories || registry.repositories.length === 0)" class="skeleton-repos">
                 <div class="skeleton-item" style="width: 90%"></div>
                 <div class="skeleton-item" style="width: 70%"></div>
                 <div class="skeleton-item" style="width: 85%"></div>
            </div>

            <div v-else-if="!registry.repositories || registry.repositories.length === 0" class="empty-repos">
                No repositories
            </div>

            <RepositoryTree
                v-else
                :repositories="registry.repositories"
                :registry-id="registry.id"
                :project-id="store.selectedProjectId"
            />
        </div>
      </div>

      <div v-if="store.error && store.registries.length === 0" class="error-container-sidebar">
          <ErrorState
            title="Failed to load registries."
            :retry="refresh"
          />
      </div>
      <div v-else-if="store.registries.length === 0 && !store.loading" class="empty-state">
          No registries found.
      </div>
    </div>

    <div v-if="configStore.authEnabled" class="sidebar-footer">
        <button @click="logout" class="logout-btn">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path>
                <polyline points="16 17 21 12 16 7"></polyline>
                <line x1="21" y1="12" x2="9" y2="12"></line>
            </svg>
            Sign Out
        </button>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { useRegistryStore } from '@/stores/registry'
import { useConfigStore } from '@/stores/config'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import RepositoryTree from './RepositoryTree.vue'
import ErrorState from '@/components/ErrorState.vue'

const store = useRegistryStore()
const configStore = useConfigStore()
const authStore = useAuthStore()
const router = useRouter()

const refresh = async () => {
    await store.refreshStructure()
}

const logout = () => {
    authStore.logout()
    router.push('/login')
}
</script>

<style scoped lang="scss">
@use "sass:color";
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

.registry-icon {
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

.error-container-sidebar {
    padding: 1rem;
}

// Loading Skeletons
@keyframes shimmer {
  0% {
    background-position: -200px 0;
  }
  100% {
    background-position: calc(200px + 100%) 0;
  }
}

.skeleton-group {
    margin-bottom: 1.5rem;
    padding: 0 1rem;
}

.skeleton-header {
    height: 1.2rem;
    background: #374151;
    background-image: linear-gradient(to right, #374151 0%, #4b5563 20%, #374151 40%, #374151 100%);
    background-repeat: no-repeat;
    background-size: 800px 100%;
    animation: shimmer 1.5s infinite linear forwards;
    border-radius: 4px;
    margin-bottom: 0.8rem;
    width: 60%;
}

.skeleton-item {
    height: 1rem;
    background: #374151;
    background-image: linear-gradient(to right, #374151 0%, #4b5563 20%, #374151 40%, #374151 100%);
    background-repeat: no-repeat;
    background-size: 800px 100%;
    animation: shimmer 1.5s infinite linear forwards;
    border-radius: 4px;
    margin-bottom: 0.5rem;
    margin-left: 1.5rem;
}

.skeleton-repos {
    padding-top: 0.5rem;
}

.spinner-small {
    width: 14px;
    height: 14px;
    border: 2px solid rgba($text-color, 0.3);
    border-radius: 50%;
    border-top-color: $primary-color;
    animation: spin 1s ease-in-out infinite;
    margin-left: 0.5rem;
}

.sidebar-footer {
    padding: 1rem;
    border-top: 1px solid $border-color;
    margin-top: auto;
}

.logout-btn {
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
    padding: 0.75rem;
    background-color: transparent;
    border: 1px solid $border-color;
    color: $text-color;
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.2s;
    font-size: 0.95rem;

    &:hover {
        background-color: rgba($danger-color, 0.1);
        color: $danger-color;
        border-color: rgba($danger-color, 0.3);
    }

    svg {
        width: 18px;
        height: 18px;
    }
}
</style>
