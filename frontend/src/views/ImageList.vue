<template>
  <div class="view-container">
    <div class="header">
      <h1>Images ({{ rname }})</h1>
      <span :title="!configStore.enableDeleteRepository ? 'Disabled by environment configuration' : ''" class="tooltip-wrapper">
          <button
              @click="openDeleteRepoModal"
              class="btn danger-outline"
              :disabled="!configStore.enableDeleteRepository"
          >
              Delete Repository
          </button>
      </span>
    </div>

    <div v-if="store.imagesLoading && store.images.length === 0" class="loading">Loading images...</div>

    <!-- Toast Notifications -->
    <ToastNotification
      v-if="store.error"
      type="error"
      :message="store.error"
      @close="store.clearNotifications"
    />
    <ToastNotification
      v-if="store.success"
      type="success"
      :message="store.success"
      @close="store.clearNotifications"
    />

    <div v-if="!store.imagesLoading || store.images.length > 0" class="list-container">
      <div class="list-controls" v-if="store.images.length > 0">
         <div class="select-all">
            <input type="checkbox" id="selectAll" :checked="allSelected" @change="toggleSelectAll" />
            <label for="selectAll">Select All</label>
         </div>

         <div class="search-box">
             <input type="text" v-model="searchQuery" placeholder="Search by tag..." class="search-input" />
         </div>

         <span :title="!configStore.enableDeleteImage ? 'Disabled by environment configuration' : ''" class="tooltip-wrapper" :class="{ 'hidden-wrapper': selectedImages.size === 0 }">
             <button
                 @click="openBulkDeleteModal"
                 class="bulk-delete-btn"
                 :disabled="selectedImages.size === 0 || !configStore.enableDeleteImage"
             >
                Delete Selected ({{ selectedImages.size }})
             </button>
         </span>
      </div>

      <div v-if="filteredImages.length === 0 && !store.imagesLoading" class="empty-state">No images found.</div>
      <div
        v-for="image in filteredImages"
        :key="image.digest"
        class="list-item"
        :class="{ 'protected-item': isProtected(image), 'deleting-item': store.deletionLoading.has(image.digest) }"
      >
        <div class="checkbox-container">
           <input
             type="checkbox"
             :value="image.digest"
             :checked="selectedImages.has(image.digest)"
             :disabled="isProtected(image) || store.deletionLoading.has(image.digest)"
             @change="toggleSelection(image.digest)"
           />
        </div>
        <div class="item-info">
          <div class="digest-row">
            <span v-if="isProtected(image)" class="protected-icon" title="This image is protected and cannot be deleted">
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-shield-lock-fill" viewBox="0 0 16 16">
                  <path fill-rule="evenodd" d="M8 0c-.69 0-1.843.265-2.928.56-1.11.3-2.229.655-2.887.87a1.54 1.54 0 0 0-1.044 1.262c-.596 4.477.787 7.795 2.465 9.99a11.777 11.777 0 0 0 2.517 2.453c.386.273.744.482 1.048.625.28.132.581.24.829.24s.548-.108.829-.24a7.159 7.159 0 0 0 1.048-.625 11.775 11.775 0 0 0 2.517-2.453c1.678-2.195 3.061-5.513 2.465-9.99a1.541 1.541 0 0 0-1.044-1.263 62.467 62.467 0 0 0-2.887-.87C9.843.266 8.69 0 8 0zm0 5a1.5 1.5 0 0 1 .5 2.915l.385 1.99a.5.5 0 0 1-.491.595h-.788a.5.5 0 0 1-.49-.595l.384-1.99A1.5 1.5 0 0 1 8 5z"/>
                </svg>
            </span>
            <span class="digest" :title="image.digest">{{ image.digest }}</span>
            <button class="copy-btn" @click="copyToClipboard(image.digest, image.digest)" title="Copy Digest">
                <span v-if="copiedState[image.digest]" class="success-icon">
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-check-lg" viewBox="0 0 16 16">
                        <path d="M12.736 3.97a.733.733 0 0 1 1.047 0c.286.289.29.756.01 1.05L7.88 12.01a.733.733 0 0 1-1.065.02L3.217 8.384a.757.757 0 0 1 0-1.06.733.733 0 0 1 1.047 0l3.052 3.093 5.4-6.425a.247.247 0 0 1 .02-.022Z"/>
                    </svg>
                </span>
                <svg v-else xmlns="http://www.w3.org/2000/svg" width="14" height="14" fill="currentColor" class="bi bi-clipboard" viewBox="0 0 16 16">
                    <path d="M4 1.5H3a2 2 0 0 0-2 2V14a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2V3.5a2 2 0 0 0-2-2h-1v1h1a1 1 0 0 1 1 1V14a1 1 0 0 1-1 1H3a1 1 0 0 1-1-1V3.5a1 1 0 0 1 1-1h1v-1z"/>
                    <path d="M9.5 1a.5.5 0 0 1 .5.5v1a.5.5 0 0 1-.5.5h-3a.5.5 0 0 1-.5-.5v-1a.5.5 0 0 1 .5-.5h3zm-3-1A1.5 1.5 0 0 0 5 1.5v1A1.5 1.5 0 0 0 6.5 4h3A1.5 1.5 0 0 0 11 2.5v-1A1.5 1.5 0 0 0 9.5 0h-3z"/>
                </svg>
            </button>
          </div>
          <div class="tags">
            <span
                v-for="tag in image.tags"
                :key="tag"
                class="tag"
                :class="{ 'copied': copiedState[image.digest + tag] }"
                @click="copyToClipboard(tag, image.digest + tag)"
                title="Click to copy tag"
            >
                {{ tag }}
                <span v-if="copiedState[image.digest + tag]" class="tag-check">✓</span>
            </span>
            <span v-if="!image.tags || image.tags.length === 0" class="no-tags">No tags</span>
          </div>
          <div class="item-meta">
            <span>Size: {{ (image.size / 1024 / 1024).toFixed(2) }} MB</span>
            <span>Created: {{ new Date(image.createdAt).toLocaleString() }}</span>
          </div>
        </div>
        <span class="tooltip-wrapper" :title="!configStore.enableDeleteImage ? 'Disabled by environment configuration' : (isProtected(image) ? 'Protected Image' : 'Delete Image')">
            <button
                v-if="!store.deletionLoading.has(image.digest)"
                @click="openDeleteImageModal(image)"
                class="delete-btn"
                :disabled="!configStore.enableDeleteImage || isProtected(image)"
            >
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-trash" viewBox="0 0 16 16">
                    <path d="M5.5 5.5A.5.5 0 0 1 6 6v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5zm2.5 0a.5.5 0 0 1 .5.5v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5zm3 .5a.5.5 0 0 0-1 0v6a.5.5 0 0 0 1 0V6z"/>
                    <path fill-rule="evenodd" d="M14.5 3a1 1 0 0 1-1 1H13v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V4h-.5a1 1 0 0 1-1-1V2a1 1 0 0 1 1-1H6a1 1 0 0 1 1-1h2a1 1 0 0 1 1 1h3.5a1 1 0 0 1 1 1v1zM4.118 4 4 4.059V13a1 1 0 0 0 1 1h6a1 1 0 0 0 1-1V4.059L11.882 4H4.118zM2.5 3V2h11v1h-11z"/>
                </svg>
            </button>
            <div v-else class="spinner-small" title="Deleting..."></div>
        </span>
      </div>
    </div>
  </div>

  <!-- Reusable Confirmation Modal -->
  <ConfirmModal
    v-if="modalState.isOpen"
    :is-open="modalState.isOpen"
    :title="modalTitle"
    :message="modalMessage"
    :is-danger="true"
    :confirm-text="modalConfirmText"
    :verification-value="modalVerificationValue"
    @update:is-open="modalState.isOpen = $event"
    @confirm="handleModalConfirm"
    @cancel="closeModal"
  >
    <!-- Single Deletion Details -->
    <div v-if="modalState.type === 'single'" class="modal-detail-container">
        <div class="modal-detail">
            <label>Digest:</label>
            <div class="modal-digest" :title="modalState.targetDigest">
                {{ modalState.targetDigest }}
            </div>
        </div>
        <div class="modal-detail" v-if="modalState.targetTags && modalState.targetTags.length > 0">
            <label>Tags:</label>
            <div class="tags">
                <span v-for="tag in modalState.targetTags" :key="tag" class="tag">{{ tag }}</span>
            </div>
        </div>
    </div>

    <!-- Checkbox for both Single and Bulk deletion -->
    <div v-if="modalState.type === 'bulk' || modalState.type === 'single'" class="form-group">
      <label>
          <input type="checkbox" v-model="deleteWithGC">
          Run garbage collection?
      </label>
      <p class="help-text">This will free up space immediately.</p>
    </div>
  </ConfirmModal>

</template>

<script setup lang="ts">
import { useRegistryStore } from '@/stores/registry'
import { useConfigStore } from '@/stores/config'
import type { Image } from '@/types'
import { onMounted, computed, ref, watch, reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import ToastNotification from '@/components/ToastNotification.vue'
import ConfirmModal from '@/components/ConfirmModal.vue'

const router = useRouter()
const route = useRoute()
const store = useRegistryStore()
const configStore = useConfigStore()
const pid = computed(() => route.params.pid as string)
const rid = computed(() => route.params.rid as string)
const rname = computed(() => route.params.rname as string)

const isProtected = (image: Image): boolean => {
    if (!configStore.protectedTags || configStore.protectedTags.length === 0) return false
    return !!(image.tags && image.tags.some(tag => configStore.protectedTags && configStore.protectedTags.includes(tag)))
}

const selectedImages = ref(new Set<string>())
const deleteWithGC = ref(true)
const searchQuery = ref("")
const copiedState = ref<Record<string, boolean>>({})

// Modal State Management
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
            return `Are you sure you want to delete ${selectedImages.value.size} selected images?`
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

const filteredImages = computed(() => {
    const images = store.images.slice().sort((a, b) => {
        return new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime()
    })

    if (!searchQuery.value) return images
    const query = searchQuery.value.toLowerCase()
    return images.filter(img =>
        img.tags && img.tags.some(tag => tag.toLowerCase().includes(query))
    )
})

const allSelected = computed(() => {
    return filteredImages.value.length > 0 && selectedImages.value.size === filteredImages.value.length
})

const fetchData = async () => {
    selectedImages.value.clear()
    searchQuery.value = ""
    await store.fetchImages(pid.value, rid.value, rname.value)
}

onMounted(() => {
  fetchData()
})

watch(() => route.fullPath, () => {
    fetchData()
})

watch(() => store.images, (newImages) => {
    // Filter selection to only include images that still exist
    const currentDigests = new Set(newImages.map(img => img.digest))
    const toRemove = []
    for (const digest of selectedImages.value) {
        if (!currentDigests.has(digest)) {
            toRemove.push(digest)
        }
    }
    toRemove.forEach(d => selectedImages.value.delete(d))
}, { deep: true })

const toggleSelection = (digest: string) => {
    if (selectedImages.value.has(digest)) {
        selectedImages.value.delete(digest)
    } else {
        selectedImages.value.add(digest)
    }
}

const toggleSelectAll = () => {
    if (allSelected.value) {
        selectedImages.value.clear()
    } else {
        filteredImages.value.forEach(img => {
            if (!isProtected(img)) {
                selectedImages.value.add(img.digest)
            }
        })
    }
}

// Modal Actions
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

const handleModalConfirm = async () => {
    modalState.isOpen = false
    try {
        if (modalState.type === 'single') {
            // Use cleanupRepository instead of deleteImage to support disable_gc option
            await store.cleanupRepository(
                pid.value,
                rid.value,
                rname.value,
                [modalState.targetDigest],
                !deleteWithGC.value
            )
        } else if (modalState.type === 'bulk') {
            await store.cleanupRepository(
                pid.value,
                rid.value,
                rname.value,
                Array.from(selectedImages.value),
                !deleteWithGC.value
            )
            selectedImages.value.clear()
        } else if (modalState.type === 'repo') {
            await store.deleteRepository(pid.value, rid.value, rname.value)
            router.push('/')
        }
    } catch (err) {
        console.error("Action failed", err)
    } finally {
        // Modal is already closed
    }
}


const copyToClipboard = (text: string, id: string) => {
    navigator.clipboard.writeText(text).then(() => {
        copiedState.value[id] = true
        setTimeout(() => {
            copiedState.value[id] = false
        }, 2000)
    }).catch(err => {
        console.error('Failed to copy text: ', err)
    })
}
</script>

<style scoped lang="scss">
@use "sass:color";
@use '@/assets/main.scss' as *;

.view-container {
  .header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 2rem;

      h1 {
        margin: 0;
        color: $primary-color;
      }
  }
}

.btn {
    padding: 0.5rem 1rem;
    border-radius: 4px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;

    &.danger-outline {
        background: transparent;
        border: 1px solid $danger-color;
        color: $danger-color;

        &:hover:not(:disabled) {
            background: rgba($danger-color, 0.1);
        }

        &:disabled {
            opacity: 0.5;
            cursor: not-allowed;
            border-color: rgba($danger-color, 0.3);
            color: rgba($danger-color, 0.5);
        }
    }
}

.list-container {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.list-item {
  background: $card-bg;
  border: 1px solid $border-color;
  border-radius: 8px;
  padding: 1.5rem;
  display: flex;
  justify-content: flex-start;
  align-items: flex-start;
  gap: 1.25rem;
  transition: transform 0.2s, box-shadow 0.2s, border-color 0.2s;

  &:hover {
    box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
    border-color: color.adjust($border-color, $lightness: 10%);
    transform: translateY(-1px);
  }
}

.item-info {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    flex: 1;
    min-width: 0;
}

.digest-row {
    display: flex;
    align-items: center;
    gap: 0.5rem;
}

.digest {
    font-family: monospace;
    font-weight: bold;
    color: $text-color;
    word-break: break-all;
    font-size: 0.9rem;
}

.copy-btn {
    background: none;
    border: none;
    padding: 2px;
    cursor: pointer;
    color: $secondary-color;
    display: flex;
    align-items: center;
    margin-left: 0.5rem;

    &:hover {
        color: $primary-color;
    }

    .success-icon {
        color: #10b981;
        animation: scaleIn 0.2s ease-out;
    }
}

@keyframes scaleIn {
    from { transform: scale(0); }
    to { transform: scale(1); }
}

.tags {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin: 0.2rem 0;
}

.tag {
  background-color: $muted-bg;
  color: $text-color;
  padding: 0.2rem 0.6rem;
  border-radius: 1rem;
  font-size: 0.85rem;
  font-family: monospace;
  cursor: pointer;
  transition: all 0.2s;
  display: inline-flex;
  align-items: center;
  gap: 0.3rem;

  &:hover {
      background-color: color.adjust($muted-bg, $lightness: 10%);
  }

  &.copied {
      background-color: #10b981; // Green
      color: white;
  }
}

.tag-check {
    font-size: 0.7rem;
    font-weight: bold;
}

.no-tags {
    color: $secondary-color;
    font-style: italic;
    font-size: 0.85rem;
}

.item-meta {
  font-size: 0.85rem;
  color: $secondary-color;
  display: flex;
  gap: 1rem;
}

.delete-btn {
  background-color: transparent;
  color: $danger-color;
  border: 1px solid transparent;
  padding: 0.5rem;
  border-radius: 4px;
  transition: background-color 0.2s;
  margin-top: -0.25rem;
  flex: 0 0 auto;

  &:hover:not(:disabled) {
    background-color: rgba(220, 53, 69, 0.1);
  }

  &:disabled {
      color: rgba($danger-color, 0.5);
      cursor: not-allowed;
  }
}

.empty-state {
    text-align: center;
    padding: 3rem;
    color: $secondary-color;
    border: 1px dashed $border-color;
    border-radius: 8px;
}

.list-controls {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.5rem 1rem;
    background: $control-bg;
    border: 1px solid $border-color;
    border-radius: 6px;
    margin-bottom: 0.5rem;
    min-height: 48px;
    gap: 1rem;
    flex-wrap: wrap;
}

.select-all {
    display: flex;
    gap: 0.5rem;
    align-items: center;
    font-weight: bold;
    white-space: nowrap;
}

.search-box {
    flex: 1;
    display: flex;
    justify-content: center;
    max-width: 400px;

    .search-input {
        width: 100%;
        padding: 0.4rem 0.8rem;
        border-radius: 4px;
        border: 1px solid $border-color;
        background-color: $card-bg;
        color: $text-color;

        &:focus {
            outline: 2px solid $primary-color;
            border-color: transparent;
        }
    }
}

.checkbox-container {
    flex: 0 0 auto;
    display: flex;
    align-items: center;
    padding-top: 0.25rem;
}

.bulk-delete-btn {
    background-color: $danger-color;
    color: white;
    border: none;
    padding: 0.5rem 1rem;
    border-radius: 4px;
    cursor: pointer;
    font-size: 0.9rem;
    transition: opacity 0.2s, background-color 0.2s;

    &:hover:not(:disabled) {
        background-color: color.adjust($danger-color, $lightness: -10%);
    }

    &.hidden-btn {
        opacity: 0;
        pointer-events: none;
    }

    &:disabled {
        opacity: 0.5;
        cursor: not-allowed;
    }
}

.modal {
    border: none;
    border-radius: 8px;
    padding: 0;
    box-shadow: 0 4px 12px rgba(0,0,0,0.5);
    max-width: 500px;
    width: 90%;
    background: $modal-bg;
    color: $text-color;

    &::backdrop {
        background: rgba(0, 0, 0, 0.7);
    }
}

.modal-content {
    padding: 2rem;

    h2 {
        margin-top: 0;
        margin-bottom: 1rem;
        color: $text-color;
    }
}

.form-group {
    margin: 1.5rem 0;

    label {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        font-weight: bold;
    }

    .help-text {
        margin-left: 1.8rem;
        font-size: 0.85rem;
        color: $secondary-color;
        margin-top: 0.25rem;
    }
}

.modal-actions {
    display: flex;
    justify-content: flex-end;
    gap: 1rem;
    margin-top: 2rem;

    button {
        padding: 0.6rem 1.2rem;
        border-radius: 4px;
        border: none;
        cursor: pointer;
        font-weight: bold;
    }

    .cancel-btn {
        background: $muted-bg;
        color: $text-color;

        &:hover {
            background: color.adjust($muted-bg, $lightness: -10%);
        }
    }

    .confirm-btn {
        background: $danger-color;
        color: white;

        &:hover {
            background: color.adjust($danger-color, $lightness: -10%);
        }
    }
}

.tooltip-wrapper {
    display: inline-block;
    cursor: help;

    &.hidden-wrapper {
        display: none;
    }
}

.btn:disabled, .bulk-delete-btn:disabled, .delete-btn:disabled {
    pointer-events: none;
}

.modal-detail-container {
    background: $background-color;
    border: 1px solid $border-color;
    border-radius: 6px;
    padding: 1rem;
    margin-bottom: 1.5rem;
    display: flex;
    flex-direction: column;
    gap: 1rem;
}

.modal-detail {
    label {
        font-weight: bold;
        display: block;
        margin-bottom: 0.25rem;
        font-size: 0.85rem;
        color: $secondary-color;
    }
}

.modal-digest {
    font-family: monospace;
    background: $card-bg;
    padding: 0.5rem;
    border-radius: 4px;
    border: 1px solid $border-color;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    cursor: help;
    font-size: 0.85rem;
    color: $text-color;
}

.protected-item {
    background-color: rgba($primary-color, 0.05);
    border-color: rgba($primary-color, 0.2);
}

.protected-icon {
    color: $primary-color;
    display: flex;
    align-items: center;
    margin-right: 0.25rem;
    cursor: help;
}

.deleting-item {
    opacity: 0.5;
    pointer-events: none;
    background-color: rgba($muted-bg, 0.3);
}

.spinner-small {
    width: 16px;
    height: 16px;
    border: 2px solid rgba($danger-color, 0.3);
    border-radius: 50%;
    border-top-color: $danger-color;
    animation: spin 1s ease-in-out infinite;
    margin: 0.5rem;
}

@keyframes spin {
    to { transform: rotate(360deg); }
}
</style>
