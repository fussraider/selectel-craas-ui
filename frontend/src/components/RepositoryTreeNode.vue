<template>
  <div class="repo-tree-node">
    <!-- Group Header -->
    <div
        v-if="node.isGroup"
        class="group-header"
        @click="toggle(node.fullPath)"
        :style="{ paddingLeft: `${paddingLeft}rem` }"
    >
      <div class="group-info">
        <span class="folder-icon">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
              <path d="M9.828 3h3.982a2 2 0 0 1 1.992 2.181l-.637 7A2 2 0 0 1 13.174 14H2.826a2 2 0 0 1-1.991-1.819l-.637-7a1.99 1.99 0 0 1 .442-1.511l.617-2.09A2.01 2.01 0 0 1 3.196 1H7.828a2 2 0 0 1 1.999 1.586l.288 1.414zM2.19 3l-.617 2.09A.995.995 0 0 0 2 5.5v.07l.637 7a.995.995 0 0 0 .991.819h10.354a.995.995 0 0 0 .99-.819l.638-7a1 1 0 0 0-.916-1.077h-.07l-3.328-1H3.196a1.01 1.01 0 0 0-.999.793l-.127.625z"/>
            </svg>
        </span>
        <span class="group-name">{{ node.name }}</span>
      </div>
      <span class="chevron" :class="{ rotated: isExpanded(node.fullPath) }">
        <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" fill="currentColor" viewBox="0 0 16 16">
          <path fill-rule="evenodd" d="M4.646 1.646a.5.5 0 0 1 .708 0l6 6a.5.5 0 0 1 0 .708l-6 6a.5.5 0 0 1-.708-.708L10.293 8 4.646 2.354a.5.5 0 0 1 0-.708z"/>
        </svg>
      </span>
    </div>

    <!-- Repo Item -->
    <RepositoryItem
      v-if="node.repo && (!node.isGroup || isExpanded(node.fullPath))"
      :repo="node.repo"
      :project-id="projectId"
      :registry-id="registryId"
      :display-name="node.isGroup ? node.name : node.name"
      :depth="node.isGroup ? depth + 1 : depth"
    />

    <!-- Children -->
    <div v-if="node.isGroup && isExpanded(node.fullPath)" class="node-children">
      <RepositoryTreeNode
        v-for="child in node.children"
        :key="child.fullPath"
        :node="child"
        :depth="depth + 1"
        :project-id="projectId"
        :registry-id="registryId"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { inject, computed } from 'vue'
import type { RepoNode } from '@/types/tree'
import RepositoryItem from './RepositoryItem.vue'
import RepositoryTreeNode from './RepositoryTreeNode.vue'

defineOptions({
  name: 'RepositoryTreeNode'
})

const props = defineProps<{
  node: RepoNode
  depth: number
  projectId: string | null
  registryId: string
}>()

const expandedKeys = inject('expandedKeys') as Set<string>
const toggle = inject('toggle') as (path: string) => void

const isExpanded = (path: string) => expandedKeys.has(path)

const paddingLeft = computed(() => {
  return 1 + props.depth * 1.5
})
</script>

<style scoped lang="scss">
@use '@/assets/main.scss' as *;

.group-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.5rem 1rem 0.5rem 1rem;
  cursor: pointer;
  color: $secondary-color;
  transition: all 0.2s;
  user-select: none;
  font-size: 0.95rem;

  &:hover {
    color: $text-color;
    background-color: rgba($text-color, 0.05);
  }
}

.group-info {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-weight: 500;
}

.folder-icon {
  opacity: 0.7;
  display: flex;
  color: $primary-color;
}

.chevron {
  transition: transform 0.2s;
  opacity: 0.5;
  display: flex;

  &.rotated {
    transform: rotate(90deg);
  }
}
</style>
