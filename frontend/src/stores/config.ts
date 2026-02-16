import { defineStore } from 'pinia'
import { ref } from 'vue'
import client, { formatError } from '@/api/client'

interface Config {
  enableDeleteRegistry: boolean
  enableDeleteRepository: boolean
  enableDeleteImage: boolean
}

export const useConfigStore = defineStore('config', () => {
  const enableDeleteRegistry = ref(false)
  const enableDeleteRepository = ref(false)
  const enableDeleteImage = ref(false)

  const loading = ref(false)
  const error = ref<string | null>(null)

  const fetchConfig = async () => {
    loading.value = true
    try {
      const res = await client.get<Config>('/config')
      enableDeleteRegistry.value = res.data.enableDeleteRegistry
      enableDeleteRepository.value = res.data.enableDeleteRepository
      enableDeleteImage.value = res.data.enableDeleteImage
    } catch (err) {
      console.error("Failed to load config", err)
      error.value = formatError(err)
    } finally {
      loading.value = false
    }
  }

  return {
    enableDeleteRegistry,
    enableDeleteRepository,
    enableDeleteImage,
    fetchConfig,
    loading,
    error
  }
})
