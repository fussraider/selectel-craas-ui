<template>
  <div class="view-container">
    <div class="header">
      <h1>Images ({{ rname }})</h1>
      <button
          @click="openDeleteRepoModal"
          class="btn danger-outline"
          :disabled="!configStore.enableDeleteRepository"
          :title="!configStore.enableDeleteRepository ? 'Disabled by environment configuration' : ''"
      >
          Delete Repository
      </button>
    </div>

    <div v-if="store.imagesLoading" class="loading">Loading images...</div>

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

    <div v-if="!store.imagesLoading" class="list-container">
      <div class="list-controls" v-if="store.images.length > 0">
         <div class="select-all">
            <input type="checkbox" id="selectAll" :checked="allSelected" @change="toggleSelectAll" />
            <label for="selectAll">Select All</label>
         </div>

         <div class="search-box">
             <input type="text" v-model="searchQuery" placeholder="Search by tag..." class="search-input" />
         </div>

         <button
             :class="{ 'hidden-btn': selectedImages.size === 0 }"
             @click="openBulkDeleteModal"
             class="bulk-delete-btn"
             :disabled="selectedImages.size === 0 || !configStore.enableDeleteImage"
             :title="!configStore.enableDeleteImage ? 'Disabled by environment configuration' : ''"
         >
            Delete Selected ({{ selectedImages.size }})
         </button>
      </div>

      <div v-if="filteredImages.length === 0" class="empty-state">No images found.</div>
      <div v-for="image in filteredImages" :key="image.digest" class="list-item">
        <div class="checkbox-container">
           <input type="checkbox" :value="image.digest" :checked="selectedImages.has(image.digest)" @change="toggleSelection(image.digest)" />
        </div>
        <div class="item-info">
          <div class="digest-row">
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
        <button
            @click="openDeleteImageModal(image.digest)"
            class="delete-btn"
            title="Delete Image"
            :disabled="!configStore.enableDeleteImage"
        >
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-trash" viewBox="0 0 16 16">
                <path d="M5.5 5.5A.5.5 0 0 1 6 6v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5zm2.5 0a.5.5 0 0 1 .5.5v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5zm3 .5a.5.5 0 0 0-1 0v6a.5.5 0 0 0 1 0V6z"/>
                <path fill-rule="evenodd" d="M14.5 3a1 1 0 0 1-1 1H13v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V4h-.5a1 1 0 0 1-1-1V2a1 1 0 0 1 1-1H6a1 1 0 0 1 1-1h2a1 1 0 0 1 1 1h3.5a1 1 0 0 1 1 1v1zM4.118 4 4 4.059V13a1 1 0 0 0 1 1h6a1 1 0 0 0 1-1V4.059L11.882 4H4.118zM2.5 3V2h11v1h-11z"/>
            </svg>
        </button>
      </div>
    </div>
  </div>

  <!-- Reusable Confirmation Modal -->
  <ConfirmModal
    v-if="modalState.isOpen"
    :is-open="modalState.isOpen"
    :title="modalState.title"
    :message="modalState.message"
    :is-danger="true"
    :confirm-text="modalState.type === 'repo' ? 'Delete Repository' : 'Delete'"
    :verification-value="modalState.verificationValue"
    @update:is-open="modalState.isOpen = $event"
    @confirm="handleModalConfirm"
    @cancel="closeModal"
  >
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

const selectedImages = ref(new Set<string>())
const deleteWithGC = ref(true)
const searchQuery = ref("")
const copiedState = ref<Record<string, boolean>>({})

// Modal State Management
const modalState = reactive({
  isOpen: false,
  type: 'single' as 'single' | 'bulk' | 'repo',
  title: '',
  message: '',
  targetDigest: '',
  verificationValue: undefined as string | undefined
})

const filteredImages = computed(() => {
    if (!searchQuery.value) return store.images
    const query = searchQuery.value.toLowerCase()
    return store.images.filter(img =>
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
        filteredImages.value.forEach(img => selectedImages.value.add(img.digest))
    }
}

// Modal Actions
const openDeleteImageModal = (digest: string) => {
    modalState.type = 'single'
    modalState.targetDigest = digest
    modalState.title = 'Delete Image'
    modalState.message = `Are you sure you want to delete image ${digest}? This cannot be undone.`
    modalState.verificationValue = undefined
    modalState.isOpen = true
}

const openBulkDeleteModal = () => {
    modalState.type = 'bulk'
    modalState.title = 'Confirm Bulk Deletion'
    modalState.message = `Are you sure you want to delete ${selectedImages.value.size} selected images?`
    modalState.verificationValue = undefined
    modalState.isOpen = true
}

const openDeleteRepoModal = () => {
    modalState.type = 'repo'
    modalState.title = 'Delete Repository'
    modalState.message = `Are you sure you want to delete repository '${rname.value}'? All images within it will be permanently lost.`
    modalState.verificationValue = rname.value
    modalState.isOpen = true
}

const closeModal = () => {
    modalState.isOpen = false
    modalState.targetDigest = ''
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
</style>
