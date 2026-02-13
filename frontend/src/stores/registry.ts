import { defineStore } from 'pinia'
import axios from 'axios'

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

export const useRegistryStore = defineStore('registry', {
  state: () => ({
    projects: [] as Project[],
    registries: [] as Registry[],
    repositories: [] as Repository[],
    images: [] as Image[],
    loading: false,
    error: null as string | null,
  }),
  actions: {
    handleError(err: unknown) {
      if (axios.isAxiosError(err)) {
        this.error = String(err.response?.data || err.message)
      } else {
        this.error = String(err)
      }
    },
    async fetchProjects() {
      this.loading = true
      this.error = null
      try {
        const res = await axios.get('/api/projects')
        this.projects = res.data
      } catch (err) {
        this.handleError(err)
      } finally {
        this.loading = false
      }
    },
    async fetchRegistries(pid: string) {
      this.loading = true
      this.error = null
      try {
        const res = await axios.get(`/api/projects/${pid}/registries`)
        this.registries = res.data
      } catch (err) {
        this.handleError(err)
      } finally {
        this.loading = false
      }
    },
    async deleteRegistry(pid: string, rid: string) {
        try {
            await axios.delete(`/api/projects/${pid}/registries/${rid}`)
            this.registries = this.registries.filter(r => r.id !== rid)
        } catch (err) {
            this.handleError(err)
            throw err
        }
    },
    async fetchRepositories(pid: string, rid: string) {
        this.loading = true
        this.error = null
        try {
            const res = await axios.get(`/api/projects/${pid}/registries/${rid}/repositories`)
            this.repositories = res.data
        } catch (err) {
            this.handleError(err)
        } finally {
            this.loading = false
        }
    },
    async deleteRepository(pid: string, rid: string, rname: string) {
        try {
            await axios.delete(`/api/projects/${pid}/registries/${rid}/repository`, { params: { name: rname } })
            this.repositories = this.repositories.filter(r => r.name !== rname)
        } catch (err) {
            this.handleError(err)
            throw err
        }
    },
    async fetchImages(pid: string, rid: string, rname: string) {
        this.loading = true
        this.error = null
        try {
            const res = await axios.get(`/api/projects/${pid}/registries/${rid}/images`, { params: { repository: rname } })
            this.images = res.data
        } catch (err) {
            this.handleError(err)
        } finally {
            this.loading = false
        }
    },
    async deleteImage(pid: string, rid: string, rname: string, digest: string) {
        try {
            await axios.delete(`/api/projects/${pid}/registries/${rid}/images/${digest}`, { params: { repository: rname } })
            this.images = this.images.filter(i => i.digest !== digest)
        } catch (err) {
            this.handleError(err)
            throw err
        }
    },
    async cleanupRepository(pid: string, rid: string, rname: string, digests: string[], disableGC: boolean = false) {
        try {
            await axios.post(`/api/projects/${pid}/registries/${rid}/cleanup`, {
                digests: digests,
                disable_gc: disableGC
            }, {
                params: { repository: rname }
            })
            // Remove deleted images from the store
            this.images = this.images.filter(i => !digests.includes(i.digest))
        } catch (err) {
            this.handleError(err)
            throw err
        }
    }
  }
})
