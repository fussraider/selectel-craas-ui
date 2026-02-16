import { defineStore } from 'pinia'
import { ref } from 'vue'
import client, { formatError } from '@/api/client'
import type { Project, Registry, Repository, Image, GCInfo, CleanupResult } from '@/types'

export const useRegistryStore = defineStore('registry', () => {
  const projects = ref<Project[]>([])
  const registries = ref<Registry[]>([])
  const images = ref<Image[]>([])
  const gcInfo = ref<GCInfo | null>(null)

  const selectedProjectId = ref<string | null>(null)

  // Loading states
  const loading = ref(false) // General/Structure loading
  const imagesLoading = ref(false) // Image list loading
  const gcLoading = ref(false) // GC info/action loading

  const error = ref<string | null>(null)
  const success = ref<string | null>(null)

  const handleError = (err: unknown) => {
    console.error(err)
    error.value = formatError(err)
  }

  const clearNotifications = () => {
      error.value = null
      success.value = null
  }

  const fetchProjects = async () => {
    loading.value = true
    clearNotifications()
    try {
      const res = await client.get<Project[]>('/projects')
      projects.value = res.data

      // Auto-select first project if none selected and list not empty
      if (!selectedProjectId.value && projects.value.length > 0 && projects.value[0]) {
          selectedProjectId.value = projects.value[0].id
      }
    } catch (err) {
      handleError(err)
    } finally {
      loading.value = false
    }
  }

  const fetchRegistries = async (pid: string) => {
    // Note: External callers should manage global loading state if chained
    clearNotifications()
    try {
      const res = await client.get<Registry[]>(`/projects/${pid}/registries`)
      // Map to add UI specific fields
      registries.value = res.data.map((r: any) => ({
          ...r,
          repositories: [],
          loadingRepos: false,
          expanded: true // Expand by default as per requirement to show repos
      }))
    } catch (err) {
      handleError(err)
    }
  }

  const fetchRepositories = async (pid: string, rid: string) => {
      const registry = registries.value.find(r => r.id === rid)
      if (registry) {
          registry.loadingRepos = true
      }

      try {
          const res = await client.get<Repository[]>(`/projects/${pid}/registries/${rid}/repositories`)
          if (registry) {
              registry.repositories = res.data
          }
      } catch (err) {
          handleError(err)
      } finally {
          if (registry) {
              registry.loadingRepos = false
          }
      }
  }

  // Orchestrator: Fetch registries then repositories for all
  const loadProjectData = async (pid: string) => {
      loading.value = true
      selectedProjectId.value = pid
      try {
        await fetchRegistries(pid)

        // Parallel fetch repositories for all registries
        const promises = registries.value.map(r => fetchRepositories(pid, r.id))
        await Promise.all(promises)
      } finally {
        loading.value = false
      }
  }

  const refreshStructure = async () => {
      if (selectedProjectId.value) {
          await loadProjectData(selectedProjectId.value)
      }
  }

  const deleteRegistry = async (pid: string, rid: string) => {
      loading.value = true
      clearNotifications()
      try {
          await client.delete(`/projects/${pid}/registries/${rid}`)
          registries.value = registries.value.filter(r => r.id !== rid)
          success.value = "Registry deleted successfully"
      } catch (err) {
          handleError(err)
          throw err
      } finally {
          loading.value = false
      }
  }

  const deleteRepository = async (pid: string, rid: string, rname: string) => {
      loading.value = true
      clearNotifications()
      try {
          await client.delete(`/projects/${pid}/registries/${rid}/repository`, { params: { name: rname } })

          // Update local state
          const registry = registries.value.find(r => r.id === rid)
          if (registry && registry.repositories) {
              registry.repositories = registry.repositories.filter(r => r.name !== rname)
          }

          success.value = `Repository ${rname} deleted successfully`
      } catch (err) {
          handleError(err)
          throw err
      } finally {
          loading.value = false
      }
  }

  const fetchImages = async (pid: string, rid: string, rname: string) => {
      imagesLoading.value = true
      clearNotifications()
      try {
          const res = await client.get<Image[]>(`/projects/${pid}/registries/${rid}/images`, { params: { repository: rname } })
          images.value = res.data
      } catch (err) {
          handleError(err)
      } finally {
          imagesLoading.value = false
      }
  }

  const deleteImage = async (pid: string, rid: string, rname: string, digest: string) => {
      imagesLoading.value = true
      clearNotifications()
      try {
          await client.delete(`/projects/${pid}/registries/${rid}/images/${digest}`, { params: { repository: rname } })
          images.value = images.value.filter(i => i.digest !== digest)
          success.value = "Image deleted successfully"
      } catch (err) {
          handleError(err)
          throw err
      } finally {
          imagesLoading.value = false
      }
  }

  const cleanupRepository = async (pid: string, rid: string, rname: string, digests: string[], disableGC: boolean = false) => {
      imagesLoading.value = true
      clearNotifications()
      try {
          const res = await client.post<CleanupResult>(`/projects/${pid}/registries/${rid}/cleanup`, {
              digests: digests,
              disable_gc: disableGC
          }, {
              params: { repository: rname }
          })

          images.value = images.value.filter(i => !digests.includes(i.digest))
          success.value = `Cleanup successful: ${res.data.deleted.length} images deleted.`
      } catch (err) {
          handleError(err)
          throw err
      } finally {
          imagesLoading.value = false
      }
  }

  const fetchGCInfo = async (pid: string, rid: string) => {
      gcLoading.value = true
      clearNotifications()
      try {
          const res = await client.get<GCInfo>(`/projects/${pid}/registries/${rid}/gc`)
          gcInfo.value = res.data
      } catch (err) {
          handleError(err)
      } finally {
          gcLoading.value = false
      }
  }

  const startGC = async (pid: string, rid: string) => {
      gcLoading.value = true
      clearNotifications()
      try {
          await client.post(`/projects/${pid}/registries/${rid}/gc`)
          success.value = "Garbage collection initiated"
      } catch (err) {
          handleError(err)
          throw err
      } finally {
          gcLoading.value = false
      }
  }

  return {
      projects,
      registries,
      images,
      gcInfo,
      selectedProjectId,
      loading,
      imagesLoading,
      gcLoading,
      error,
      success,
      fetchProjects,
      fetchRegistries,
      deleteRegistry,
      fetchRepositories,
      deleteRepository,
      fetchImages,
      deleteImage,
      cleanupRepository,
      fetchGCInfo,
      startGC,
      clearNotifications,
      loadProjectData,
      refreshStructure
  }
})
