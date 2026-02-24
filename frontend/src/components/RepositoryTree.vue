<template>
  <div class="repository-tree">
     <RepositoryTreeNode
        v-for="node in treeNodes"
        :key="node.fullPath"
        :node="node"
        :depth="0"
        :project-id="projectId"
        :registry-id="registryId"
     />
  </div>
</template>

<script setup lang="ts">
import { computed, provide, reactive } from 'vue'
import type { Repository } from '@/types'
import type { RepoNode } from '@/types/tree'
import RepositoryTreeNode from './RepositoryTreeNode.vue'

const props = defineProps<{
  repositories: Repository[]
  projectId: string | null
  registryId: string
}>()

const expandedKeys = reactive(new Set<string>())

// Toggle function provided to children
const toggle = (path: string) => {
  if (expandedKeys.has(path)) {
    expandedKeys.delete(path)
  } else {
    expandedKeys.add(path)
  }
}

provide('expandedKeys', expandedKeys)
provide('toggle', toggle)

const treeNodes = computed(() => {
  const nodes: RepoNode[] = []
  const map = new Map<string, RepoNode>()

  // Sort repositories by name to ensure consistent tree
  const sortedRepos = [...props.repositories].sort((a, b) => a.name.localeCompare(b.name))

  sortedRepos.forEach(repo => {
    const parts = repo.name.split('/')
    let currentPath = ''

    // Iterate parts to build tree
    parts.forEach((part, index) => {
        const isLast = index === parts.length - 1
        const parentPath = currentPath
        currentPath = currentPath ? `${currentPath}/${part}` : part

        let node = map.get(currentPath)

        if (!node) {
            node = {
                name: part,
                fullPath: currentPath,
                isGroup: !isLast, // Initially assume group if not last part
                children: [],
                repo: undefined
            }
            map.set(currentPath, node)

            if (index === 0) {
                nodes.push(node)
            } else {
                const parent = map.get(parentPath)
                if (parent) {
                    parent.children.push(node)
                    // Ensure parent is marked as group since it has children
                    parent.isGroup = true
                }
            }
        } else {
             // If node exists and we are traversing deeper, it must be a group
             if (!isLast) {
                 node.isGroup = true
             }
        }

        // If this is the last part, it's the repository itself
        if (isLast) {
            node.repo = repo
        }
    })
  })

  return nodes
})
</script>
