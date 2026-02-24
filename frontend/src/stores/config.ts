import { defineStore } from 'pinia'
import { ref } from 'vue'
import client, { formatError } from '@/api/client'

interface Config {
  enableDeleteRegistry: boolean
  enableDeleteRepository: boolean
  enableDeleteImage: boolean
  protectedTags?: string[]
  authEnabled: boolean
}

export const useConfigStore = defineStore('config', () => {
  const enableDeleteRegistry = ref(false)
  const enableDeleteRepository = ref(false)
  const enableDeleteImage = ref(false)
  const protectedTags = ref<string[]>([])
  const authEnabled = ref(false)

  const loading = ref(false)
  const error = ref<string | null>(null)

  const fetchConfig = async () => {
    loading.value = true
    try {
      const res = await client.get<Config>('/config')
      enableDeleteRegistry.value = res.data.enableDeleteRegistry
      enableDeleteRepository.value = res.data.enableDeleteRepository
      enableDeleteImage.value = res.data.enableDeleteImage
      protectedTags.value = res.data.protectedTags || []
      authEnabled.value = res.data.authEnabled
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
    protectedTags,
    authEnabled,
    fetchConfig,
    loading,
    error
  }
})
