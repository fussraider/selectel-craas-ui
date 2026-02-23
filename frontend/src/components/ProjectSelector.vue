<template>
  <div class="project-selector">
    <label for="project-select" class="selector-label">Project:</label>
    <select
      id="project-select"
      v-model="selectedProject"
      @change="onProjectChange"
      :disabled="store.loading"
    >
      <option v-for="project in store.projects" :key="project.id" :value="project.id">
        {{ project.name }} ({{ project.id }})
      </option>
    </select>
  </div>
</template>

<script setup lang="ts">
import { useRegistryStore } from '@/stores/registry'
import { onMounted, ref, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'

const store = useRegistryStore()
const router = useRouter()
const route = useRoute()
const selectedProject = ref<string>("")

onMounted(async () => {
  await store.fetchProjects()
  if (store.projects.length > 0) {
      const pidFromUrl = route.params.pid as string

      // Prioritize URL param if valid
      if (pidFromUrl && store.projects.some(p => p.id === pidFromUrl)) {
          store.selectedProjectId = pidFromUrl
          selectedProject.value = pidFromUrl
          await store.loadProjectData(pidFromUrl)
      } else if (store.selectedProjectId) {
          selectedProject.value = store.selectedProjectId
          await store.loadProjectData(selectedProject.value)
      } else if (store.projects.length > 0 && store.projects[0]) {
          // Fallback to first project
          store.selectedProjectId = store.projects[0].id
          selectedProject.value = store.projects[0].id
          await store.loadProjectData(selectedProject.value)
      }
  }
})

// Sync URL changes to store (e.g. back/forward navigation)
watch(() => route.params.pid, async (newPid) => {
    const pid = newPid as string
    if (pid && pid !== store.selectedProjectId && store.projects.some(p => p.id === pid)) {
        store.selectedProjectId = pid
        selectedProject.value = pid
        await store.loadProjectData(pid)
    }
})

// Watch store for changes (in case changed elsewhere)
watch(() => store.selectedProjectId, (newVal) => {
    if (newVal) selectedProject.value = newVal
})

const onProjectChange = async () => {
    if (selectedProject.value) {
        store.selectedProjectId = selectedProject.value
        await store.loadProjectData(selectedProject.value)
        router.push('/') // Reset view to home/dashboard or keep current if valid?
        // Actually, if we change project, the current repo view is invalid.
        // We should probably go to the project root or just stay on "/" until a repo is selected.
    }
}
</script>

<style scoped lang="scss">
@use '@/assets/main.scss' as *;

.project-selector {
  display: flex;
  align-items: center;
  gap: 0.5rem;

  label {
      font-weight: 600;
      color: $text-color;
  }

  select {
      background-color: $card-bg;
      color: $text-color;
      border: 1px solid $border-color;
      padding: 0.5rem;
      border-radius: 4px;
      font-size: 0.95rem;
      min-width: 200px;
      cursor: pointer;

      &:focus {
          outline: 2px solid $primary-color;
          border-color: transparent;
      }
  }

  @media (max-width: 768px) {
      .selector-label {
          display: none;
      }

      select {
          min-width: 120px;
          max-width: 150px;
          font-size: 0.85rem;
      }
  }
}
</style>
