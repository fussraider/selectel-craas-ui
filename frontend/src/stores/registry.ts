import { defineStore } from 'pinia'
import { ref } from 'vue'
import axios, { AxiosError } from 'axios'

export interface Project {
  id: string
  name: string
}

export interface Registry {
  id: string
  name: string
  status: string
  createdAt: string
  size: number
}

export interface Repository {
  name: string
  size: number
  updatedAt: string
}

export interface Image {
  digest: string
  tags: string[]
  size: number
  createdAt: string
}

export interface CleanupResult {
    deleted: any[]
    failed: any[]
}

export const useRegistryStore = defineStore('registry', () => {
  const projects = ref<Project[]>([])
  const registries = ref<Registry[]>([])
  const repositories = ref<Repository[]>([])
  const images = ref<Image[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)
  const success = ref<string | null>(null) // Added for success messages

  const handleError = (err: unknown) => {
    console.error(err)
    if (axios.isAxiosError(err)) {
        // Try to get message from response body
        const msg = err.response?.data?.error || err.response?.data || err.message
        error.value = typeof msg === 'string' ? msg : JSON.stringify(msg)
    } else if (err instanceof Error) {
        error.value = err.message
    } else {
        error.value = String(err)
    }
  }

  const clearNotifications = () => {
      error.value = null
      success.value = null
  }

  const fetchProjects = async () => {
    loading.value = true
    clearNotifications()
    try {
      const res = await axios.get('/api/projects')
      projects.value = res.data
    } catch (err) {
      handleError(err)
    } finally {
      loading.value = false
    }
  }

  const fetchRegistries = async (pid: string) => {
    loading.value = true
    clearNotifications()
    try {
      const res = await axios.get(`/api/projects/${pid}/registries`)
      registries.value = res.data
    } catch (err) {
      handleError(err)
    } finally {
      loading.value = false
    }
  }

  const deleteRegistry = async (pid: string, rid: string) => {
      loading.value = true
      clearNotifications()
      try {
          await axios.delete(`/api/projects/${pid}/registries/${rid}`)
          registries.value = registries.value.filter(r => r.id !== rid)
          success.value = "Registry deleted successfully"
      } catch (err) {
          handleError(err)
          throw err
      } finally {
          loading.value = false
      }
  }

  const fetchRepositories = async (pid: string, rid: string) => {
      loading.value = true
      clearNotifications()
      try {
          const res = await axios.get(`/api/projects/${pid}/registries/${rid}/repositories`)
          repositories.value = res.data
      } catch (err) {
          handleError(err)
      } finally {
          loading.value = false
      }
  }

  const deleteRepository = async (pid: string, rid: string, rname: string) => {
      loading.value = true
      clearNotifications()
      try {
          await axios.delete(`/api/projects/${pid}/registries/${rid}/repository`, { params: { name: rname } })
          repositories.value = repositories.value.filter(r => r.name !== rname)
          success.value = `Repository ${rname} deleted successfully`
      } catch (err) {
          handleError(err)
          throw err
      } finally {
          loading.value = false
      }
  }

  const fetchImages = async (pid: string, rid: string, rname: string) => {
      loading.value = true
      clearNotifications()
      try {
          const res = await axios.get(`/api/projects/${pid}/registries/${rid}/images`, { params: { repository: rname } })
          images.value = res.data
      } catch (err) {
          handleError(err)
      } finally {
          loading.value = false
      }
  }

  const deleteImage = async (pid: string, rid: string, rname: string, digest: string) => {
      loading.value = true
      clearNotifications()
      try {
          await axios.delete(`/api/projects/${pid}/registries/${rid}/images/${digest}`, { params: { repository: rname } })
          images.value = images.value.filter(i => i.digest !== digest)
          success.value = "Image deleted successfully"
      } catch (err) {
          handleError(err)
          throw err
      } finally {
          loading.value = false
      }
  }

  const cleanupRepository = async (pid: string, rid: string, rname: string, digests: string[], disableGC: boolean = false) => {
      loading.value = true
      clearNotifications()
      try {
          const res = await axios.post<CleanupResult>(`/api/projects/${pid}/registries/${rid}/cleanup`, {
              digests: digests,
              disable_gc: disableGC
          }, {
              params: { repository: rname }
          })

          // Remove deleted images from the store
          // const deletedDigests = new Set(res.data.deleted.map((d: any) => d.digest))

          // Filter out the requested digests
          images.value = images.value.filter(i => !digests.includes(i.digest))

          success.value = `Cleanup successful: ${res.data.deleted.length} images deleted.`
      } catch (err) {
          handleError(err)
          throw err
      } finally {
          loading.value = false
      }
  }

  return {
      projects,
      registries,
      repositories,
      images,
      loading,
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
      clearNotifications
  }
})
