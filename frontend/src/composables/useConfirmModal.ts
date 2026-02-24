import { reactive, computed, type Ref } from 'vue'
import type { Image } from '@/types'

export function useConfirmModal(rname: Ref<string>, selectedImagesCount: Ref<number>) {
  const modalState = reactive({
    isOpen: false,
    type: 'single' as 'single' | 'bulk' | 'repo',
    targetDigest: '',
    targetTags: [] as string[]
  })

  const modalTitle = computed(() => {
    switch (modalState.type) {
      case 'single': return 'Delete Image'
      case 'bulk': return 'Confirm Bulk Deletion'
      case 'repo': return 'Delete Repository'
      default: return 'Confirm Action'
    }
  })

  const modalMessage = computed(() => {
    switch (modalState.type) {
      case 'single':
        return 'Are you sure you want to delete this image? This cannot be undone.'
      case 'bulk':
        return `Are you sure you want to delete ${selectedImagesCount.value} selected images?`
      case 'repo':
        return `Are you sure you want to delete repository '${rname.value}'? All images within it will be permanently lost.`
      default:
        return 'Are you sure you want to proceed?'
    }
  })

  const modalConfirmText = computed(() => {
    return modalState.type === 'repo' ? 'Delete Repository' : 'Delete'
  })

  const modalVerificationValue = computed(() => {
    return modalState.type === 'repo' ? rname.value : undefined
  })

  const openDeleteImageModal = (image: Image) => {
    modalState.type = 'single'
    modalState.targetDigest = image.digest
    modalState.targetTags = image.tags || []
    modalState.isOpen = true
  }

  const openBulkDeleteModal = () => {
    modalState.type = 'bulk'
    modalState.isOpen = true
  }

  const openDeleteRepoModal = () => {
    modalState.type = 'repo'
    modalState.isOpen = true
  }

  const closeModal = () => {
    modalState.isOpen = false
    modalState.targetDigest = ''
    modalState.targetTags = []
  }

  return {
    modalState,
    modalTitle,
    modalMessage,
    modalConfirmText,
    modalVerificationValue,
    openDeleteImageModal,
    openBulkDeleteModal,
    openDeleteRepoModal,
    closeModal
  }
}
