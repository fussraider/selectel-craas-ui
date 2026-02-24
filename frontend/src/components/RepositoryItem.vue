<template>
  <router-link
    :to="`/projects/${projectId}/registries/${registryId}/repositories/${encodeURIComponent(repo.name)}`"
    class="repo-link"
    active-class="active"
    :style="{ paddingLeft: `${paddingLeft}rem` }"
  >
    <div class="repo-content">
      <div class="repo-name">
        <span class="repo-icon">📄</span>
        {{ displayName || repo.name }}
      </div>
      <div class="repo-meta">
        <span v-if="repo.size !== undefined">{{ formatSize(repo.size) }}</span>
      </div>
    </div>
  </router-link>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Repository } from '@/types'

const props = defineProps<{
  repo: Repository
  projectId: string | null
  registryId: string
  displayName?: string
  depth?: number
}>()

const paddingLeft = computed(() => {
  // Base padding 2rem (as in sidebar) + depth * 1.5
  return 2 + (props.depth || 0) * 1.5
})

const formatSize = (bytes: number) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}
</script>

<style scoped lang="scss">
@use "sass:color";
@use '@/assets/main.scss' as *;

.repo-link {
  display: block;
  padding: 0.5rem 1rem 0.5rem 2rem;
  text-decoration: none;
  border-left: 2px solid transparent;
  transition: all 0.2s;
  color: $secondary-color;

  &:hover {
    background-color: rgba($primary-color, 0.05);
    color: $text-color;
  }

  &.active {
    background-color: rgba($primary-color, 0.1);
    color: $primary-color;
    border-left-color: $primary-color;

    .repo-name {
      font-weight: 500;
    }
  }
}

.repo-content {
  display: flex;
  flex-direction: column;
  gap: 0.2rem;
}

.repo-name {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.9rem;
  word-break: break-all;
}

.repo-meta {
  font-size: 0.75rem;
  color: color.adjust($secondary-color, $lightness: -20%);
  margin-left: 1.4rem;
  display: flex;
  gap: 0.5rem;
  opacity: 0.8;
}

.repo-icon {
  opacity: 0.7;
}
</style>
