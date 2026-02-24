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
import { computed, provide, reactive, watch, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import type { Repository } from '@/types'
import type { RepoNode } from '@/types/tree'
import RepositoryTreeNode from './RepositoryTreeNode.vue'

const props = defineProps<{
  repositories: Repository[]
  projectId: string | null
  registryId: string
}>()

const route = useRoute()
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

const sortNodes = (nodes: RepoNode[]) => {
    // Sort groups first, then by name
    nodes.sort((a, b) => {
        if (a.isGroup === b.isGroup) {
            return a.name.localeCompare(b.name)
        }
        return a.isGroup ? -1 : 1
    })

    // Recursively sort children
    nodes.forEach(node => {
        if (node.children.length > 0) {
            sortNodes(node.children)
        }
    })
}

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

  // Apply sorting: Groups first, then alphanumeric
  sortNodes(nodes)

  return nodes
})

// Auto-expand logic based on current route
const expandActivePath = () => {
    // The route parameter is named 'rname' in router/index.ts
    const repoName = route.params.rname

    if (typeof repoName === 'string') {
        const decodedName = decodeURIComponent(repoName)
        const parts = decodedName.split('/')

        let currentPath = ''
        // Iterate up to the second to last part (parents)
        for (let i = 0; i < parts.length - 1; i++) {
            const part = parts[i] ?? ''
            if (!part) continue
            currentPath = currentPath ? `${currentPath}/${part}` : part
            expandedKeys.add(currentPath)
        }
    }
}

watch(
    () => route.params.rname,
    () => {
        expandActivePath()
    },
    { immediate: true }
)

onMounted(() => {
    expandActivePath()
})
</script>
