<template>
  <div class="view-container">
    <!-- Mobile Repository Header -->
    <div class="mobile-repo-header">
      <div class="repo-title">
        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="repo-icon">
          <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"></path>
          <polyline points="3.27 6.96 12 12.01 20.73 6.96"></polyline>
          <line x1="12" y1="22.08" x2="12" y2="12"></line>
        </svg>
        <span>{{ rname }}</span>
      </div>
    </div>

    <!-- Teleport Delete Repo button to header -->
    <Teleport to="#header-actions">
        <span :title="!configStore.enableDeleteRepository ? 'Disabled by environment configuration' : ''" class="tooltip-wrapper header-action-btn">
            <button
                @click="openDeleteRepoModal"
                class="btn danger-outline small-btn"
                :disabled="!configStore.enableDeleteRepository"
            >
                Delete Repo
            </button>
        </span>
    </Teleport>

    <div v-if="store.imagesLoading" class="list-container skeleton-container">
        <div v-for="n in 5" :key="n" class="list-item skeleton-item-row">
            <div class="skeleton-checkbox"></div>
            <div class="item-info">
                <div class="skeleton-line width-40"></div>
                <div class="tags-row">
                    <div class="skeleton-tag"></div>
                    <div class="skeleton-tag"></div>
                    <div class="skeleton-tag"></div>
                </div>
                <div class="skeleton-line width-20"></div>
            </div>
        </div>
    </div>

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
      <div class="list-controls" v-if="store.images.length > 0" ref="controlsRef">
         <div class="select-all">
            <input type="checkbox" id="selectAll" :checked="allSelected" @change="toggleSelectAll" />
            <label for="selectAll">Select All</label>
         </div>

         <div class="search-box">
             <input type="text" v-model="searchQuery" placeholder="Search by tag..." class="search-input" />
         </div>

         <span :title="!configStore.enableDeleteImage ? 'Disabled by environment configuration' : ''" class="tooltip-wrapper bulk-delete-wrapper" :class="{ 'hidden-wrapper': selectedImages.size === 0 }">
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
            <div class="digest-left">
                <span v-if="isProtected(image)" class="protected-icon" title="This image is protected and cannot be deleted">
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-shield-lock-fill" viewBox="0 0 16 16">
                    <path fill-rule="evenodd" d="M8 0c-.69 0-1.843.265-2.928.56-1.11.3-2.229.655-2.887.87a1.54 1.54 0 0 0-1.044 1.262c-.596 4.477.787 7.795 2.465 9.99a11.777 11.777 0 0 0 2.517 2.453c.386.273.744.482 1.048.625.28.132.581.24.829.24s.548-.108.829-.24a7.159 7.159 0 0 0 1.048-.625 11.775 11.775 0 0 0 2.517-2.453c1.678-2.195 3.061-5.513 2.465-9.99a1.541 1.541 0 0 0-1.044-1.263 62.467 62.467 0 0 0-2.887-.87C9.843.266 8.69 0 8 0zm0 5a1.5 1.5 0 0 1 .5 2.915l.385 1.99a.5.5 0 0 1-.491.595h-.788a.5.5 0 0 1-.49-.595l.384-1.99A1.5 1.5 0 0 1 8 5z"/>
                    </svg>
                </span>
                <span class="digest" :title="image.digest">{{ image.digest }}</span>
                <button class="copy-btn" @click="copyToClipboard(`cr.selcloud.ru/${registryName}/${rname}@${image.digest}`, image.digest)" title="Copy Digest">
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
          <div class="tags">
            <span
                v-for="tag in image.tags"
                :key="tag"
                class="tag"
                :class="{ 'copied': copiedState[image.digest + tag] }"
                @click="copyToClipboard(`cr.selcloud.ru/${registryName}/${rname}:${tag}`, image.digest + tag)"
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
      </div>
    </div>
  </div>

  <Transition name="slide-up">
      <div v-if="selectedImages.size > 0 && !isControlsVisible" class="sticky-actions">
          <button
                 @click="openBulkDeleteModal"
                 class="bulk-delete-btn shadow-btn"
                 :disabled="!configStore.enableDeleteImage"
             >
                Delete Selected ({{ selectedImages.size }})
          </button>
      </div>
  </Transition>

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
import { onMounted, onUnmounted, computed, ref, watch, useTemplateRef } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import ToastNotification from '@/components/ToastNotification.vue'
import ConfirmModal from '@/components/ConfirmModal.vue'
import { useConfirmModal } from '@/composables/useConfirmModal'

const router = useRouter()
const route = useRoute()
const store = useRegistryStore()
const configStore = useConfigStore()
const pid = computed(() => route.params.pid as string)
const rid = computed(() => route.params.rid as string)
const rname = computed(() => route.params.rname as string)
const registry = computed(() => store.registries.find(r => r.id === rid.value))
const registryName = computed(() => registry.value?.name || '')

const isProtected = (image: Image): boolean => {
    if (!configStore.protectedTags || configStore.protectedTags.length === 0) return false
    return !!(image.tags && image.tags.some(tag => configStore.protectedTags && configStore.protectedTags.includes(tag)))
}

const selectedImages = ref(new Set<string>())
const selectedImagesCount = computed(() => selectedImages.value.size)
const deleteWithGC = ref(true)
const searchQuery = ref("")
const copiedState = ref<Record<string, boolean>>({})
const isControlsVisible = ref(true)
const controlsRef = useTemplateRef('controlsRef')
let observer: IntersectionObserver | null = null

const {
  modalState,
  modalTitle,
  modalMessage,
  modalConfirmText,
  modalVerificationValue,
  openDeleteImageModal,
  openBulkDeleteModal,
  openDeleteRepoModal,
  closeModal
} = useConfirmModal(rname, selectedImagesCount)

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
    if (filteredImages.value.length === 0) return false
    const selectable = filteredImages.value.filter(img => !isProtected(img))
    return selectable.length > 0 && selectable.every(img => selectedImages.value.has(img.digest))
})

const fetchData = async () => {
    selectedImages.value.clear()
    searchQuery.value = ""
    await store.fetchImages(pid.value, rid.value, rname.value)
}

onMounted(() => {
  fetchData()

  if (window.IntersectionObserver) {
      observer = new IntersectionObserver((entries) => {
          if (entries[0]) {
              isControlsVisible.value = entries[0].isIntersecting
          }
      }, { threshold: 0 })

      // Watch for ref availability
      watch(controlsRef, (el) => {
          if (el && observer) {
              observer.disconnect()
              observer.observe(el)
          }
      }, { immediate: true })
  }
})

onUnmounted(() => {
    if (observer) observer.disconnect()
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
    /* Main container tweaks if needed */
}

.mobile-repo-header {
  display: none; // Hidden on desktop
  padding: 0.75rem 0;
  margin-bottom: 0.5rem;
  border-bottom: 1px solid $border-color;
  align-items: center;

  @media (max-width: 768px) {
    display: flex;
  }

  .repo-title {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 1.1rem;
    font-weight: 600;
    color: $text-color;
    overflow: hidden;

    span {
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
      max-width: 80vw;
    }

    .repo-icon {
      color: $primary-color;
      flex-shrink: 0;
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

    &.small-btn {
        padding: 0.3rem 0.6rem;
        font-size: 0.85rem;
    }
}

.header-action-btn {
    margin-left: auto;
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
    justify-content: space-between;
    gap: 0.5rem;
    width: 100%;
}

.digest-left {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    min-width: 0;
    flex: 1;
}

.digest {
    font-family: monospace;
    font-weight: bold;
    color: $text-color;
    font-size: 0.9rem;

    display: inline-block;
    max-width: 15ch;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    vertical-align: middle;
}

@media (min-width: 768px) {
    .digest {
        max-width: 40ch;
    }
}

@media (max-width: 768px) {
    .list-item {
        padding: 1rem;
    }
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
  flex-wrap: wrap;
  align-items: center;
  justify-content: flex-start;
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
    flex-wrap: wrap;
    align-items: center;
    padding: 0.5rem;
    background: $control-bg;
    border: 1px solid $border-color;
    border-radius: 6px;
    margin-bottom: 0.5rem;
    gap: 0.5rem;

    @media (min-width: 768px) {
        flex-wrap: nowrap;
        padding: 0.5rem 1rem;
    }
}

.search-box {
    width: 100%;
    order: 1;

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

    @media (min-width: 768px) {
        flex: 1;
        order: 2;
        margin: 0 1rem;
    }
}

.select-all {
    display: flex;
    gap: 0.5rem;
    align-items: center;
    font-weight: bold;
    white-space: nowrap;
    order: 2;

    @media (min-width: 768px) {
        order: 1;
    }
}

.bulk-delete-wrapper {
    order: 3;
    margin-left: auto;

    @media (min-width: 768px) {
        margin-left: 0;
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
    vertical-align: middle;
    cursor: help;

    &.hidden-wrapper {
        visibility: hidden;
        pointer-events: none;
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

/* Skeleton Loader */
@keyframes shimmer {
  0% {
    background-position: -200px 0;
  }
  100% {
    background-position: calc(200px + 100%) 0;
  }
}

.skeleton-item-row {
    pointer-events: none;
    gap: 1.25rem;
}

.skeleton-checkbox {
    width: 16px;
    height: 16px;
    background: #374151; /* Match dark theme base */
    border-radius: 4px;
    opacity: 0.5;
}

.skeleton-line {
    height: 1rem;
    background: #374151;
    background-image: linear-gradient(to right, #374151 0%, #4b5563 20%, #374151 40%, #374151 100%);
    background-repeat: no-repeat;
    background-size: 800px 100%;
    animation: shimmer 1.5s infinite linear forwards;
    border-radius: 4px;
    margin-bottom: 0.5rem;
}

.tags-row {
    display: flex;
    gap: 0.5rem;
    margin: 0.5rem 0;
}

.skeleton-tag {
    height: 1.4rem;
    width: 60px;
    background: #374151;
    background-image: linear-gradient(to right, #374151 0%, #4b5563 20%, #374151 40%, #374151 100%);
    background-repeat: no-repeat;
    background-size: 800px 100%;
    animation: shimmer 1.5s infinite linear forwards;
    border-radius: 1rem;
}

.width-60 { width: 60%; }
.width-40 { width: 40%; }
.width-30 { width: 30%; }
.width-20 { width: 20%; }

.sticky-actions {
    position: fixed;
    bottom: 2rem;
    left: 50%;
    transform: translateX(-50%);
    z-index: 100;
    padding: 0.5rem;

    .shadow-btn {
        box-shadow: 0 4px 12px rgba(0,0,0,0.3);
        padding: 0.8rem 1.5rem;
        font-weight: bold;
    }
}

.slide-up-enter-active,
.slide-up-leave-active {
  transition: all 0.3s ease-out;
}

.slide-up-enter-from,
.slide-up-leave-to {
  transform: translate(-50%, 200%); // Start/End below screen
  opacity: 0;
}
</style>
