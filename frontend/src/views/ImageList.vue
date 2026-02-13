<template>
  <div class="view-container">
    <div class="header">
      <h1>Images ({{ rname }})</h1>
    </div>
    <div v-if="store.loading" class="loading">Loading images...</div>
    <div v-else-if="store.error" class="error">{{ store.error }}</div>
    <div v-else class="list-container">
      <div class="list-controls" v-if="store.images.length > 0">
         <div class="select-all">
            <input type="checkbox" id="selectAll" :checked="allSelected" @change="toggleSelectAll" />
            <label for="selectAll">Select All</label>
         </div>
         <button v-if="selectedImages.size > 0" @click="confirmDelete" class="bulk-delete-btn">
            Delete Selected ({{ selectedImages.size }})
         </button>
      </div>

      <div v-if="store.images.length === 0" class="empty-state">No images found.</div>
      <div v-for="image in store.images" :key="image.digest" class="list-item">
        <div class="checkbox-container">
           <input type="checkbox" :value="image.digest" :checked="selectedImages.has(image.digest)" @change="toggleSelection(image.digest)" />
        </div>
        <div class="item-info">
          <div class="digest-row">
            <span class="digest" :title="image.digest">{{ image.digest }}</span>
            <button class="copy-btn" @click="copyToClipboard(image.digest)" title="Copy Digest">
                <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" fill="currentColor" class="bi bi-clipboard" viewBox="0 0 16 16">
                    <path d="M4 1.5H3a2 2 0 0 0-2 2V14a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2V3.5a2 2 0 0 0-2-2h-1v1h1a1 1 0 0 1 1 1V14a1 1 0 0 1-1 1H3a1 1 0 0 1-1-1V3.5a1 1 0 0 1 1-1h1v-1z"/>
                    <path d="M9.5 1a.5.5 0 0 1 .5.5v1a.5.5 0 0 1-.5.5h-3a.5.5 0 0 1-.5-.5v-1a.5.5 0 0 1 .5-.5h3zm-3-1A1.5 1.5 0 0 0 5 1.5v1A1.5 1.5 0 0 0 6.5 4h3A1.5 1.5 0 0 0 11 2.5v-1A1.5 1.5 0 0 0 9.5 0h-3z"/>
                </svg>
            </button>
          </div>
          <div class="tags">
            <span v-for="tag in image.tags" :key="tag" class="tag">{{ tag }}</span>
            <span v-if="!image.tags || image.tags.length === 0" class="no-tags">No tags</span>
          </div>
          <div class="item-meta">
            <span>Size: {{ (image.size / 1024 / 1024).toFixed(2) }} MB</span>
            <span>Created: {{ new Date(image.createdAt).toLocaleString() }}</span>
          </div>
        </div>
        <button @click="deleteImg(image.digest)" class="delete-btn" title="Delete Image">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-trash" viewBox="0 0 16 16">
                <path d="M5.5 5.5A.5.5 0 0 1 6 6v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5zm2.5 0a.5.5 0 0 1 .5.5v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5zm3 .5a.5.5 0 0 0-1 0v6a.5.5 0 0 0 1 0V6z"/>
                <path fill-rule="evenodd" d="M14.5 3a1 1 0 0 1-1 1H13v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V4h-.5a1 1 0 0 1-1-1V2a1 1 0 0 1 1-1H6a1 1 0 0 1 1-1h2a1 1 0 0 1 1 1h3.5a1 1 0 0 1 1 1v1zM4.118 4 4 4.059V13a1 1 0 0 0 1 1h6a1 1 0 0 0 1-1V4.059L11.882 4H4.118zM2.5 3V2h11v1h-11z"/>
            </svg>
        </button>
      </div>
    </div>
  </div>

  <dialog ref="deleteModal" class="modal">
    <div class="modal-content">
      <h2>Confirm Deletion</h2>
      <p>Are you sure you want to delete {{ selectedImages.size }} selected images?</p>

      <div class="form-group">
        <label>
            <input type="checkbox" v-model="deleteWithGC">
            Run garbage collection?
        </label>
        <p class="help-text">This will free up space immediately.</p>
      </div>

      <div class="modal-actions">
        <button @click="closeModal" class="cancel-btn">Cancel</button>
        <button @click="executeBulkDelete" class="confirm-btn">Delete</button>
      </div>
    </div>
  </dialog>
</template>

<script setup lang="ts">
import { useRegistryStore } from '@/stores/registry'
import { onMounted, computed, ref } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const store = useRegistryStore()
const pid = computed(() => route.params.pid as string)
const rid = computed(() => route.params.rid as string)
const rname = computed(() => route.params.rname as string)

const selectedImages = ref(new Set<string>())
const deleteWithGC = ref(true)
const deleteModal = ref<HTMLDialogElement | null>(null)

const allSelected = computed(() => {
    return store.images.length > 0 && selectedImages.value.size === store.images.length
})

onMounted(() => {
  store.fetchImages(pid.value, rid.value, rname.value)
})

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
        store.images.forEach(img => selectedImages.value.add(img.digest))
    }
}

const confirmDelete = () => {
    deleteModal.value?.showModal()
}

const closeModal = () => {
    deleteModal.value?.close()
}

const executeBulkDelete = async () => {
    try {
        await store.cleanupRepository(
            pid.value,
            rid.value,
            rname.value,
            Array.from(selectedImages.value),
            !deleteWithGC.value // API expects disable_gc
        )
        selectedImages.value.clear()
        closeModal()
    } catch (err) {
        // Error is handled in store, displayed in UI
        closeModal()
    }
}

const deleteImg = async (digest: string) => {
  if (confirm(`Are you sure you want to delete image ${digest}?`)) {
    await store.deleteImage(pid.value, rid.value, rname.value, digest)
  }
}

const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text).catch(err => {
        console.error('Failed to copy text: ', err)
    })
}
</script>

<style scoped lang="scss">
@use "sass:color";
@use '@/assets/main.scss' as *;

.view-container {
  h1 {
    margin-bottom: 2rem;
    color: $primary-color;
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
  border-radius: 6px;
  padding: 1.25rem;
  display: flex;
  justify-content: flex-start; // Changed from space-between
  align-items: flex-start;
  gap: 1rem; // Add gap for spacing
  transition: box-shadow 0.2s;

  &:hover {
    box-shadow: 0 2px 4px rgba(0,0,0,0.2);
  }
}

.item-info {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    flex: 1; // Take up remaining space
    min-width: 0; // Allow shrinking for text overflow
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
    word-break: break-all; // Handle long hashes
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

    &:hover {
        color: $primary-color;
    }
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
  flex: 0 0 auto; // Prevent shrinking

  &:hover {
    background-color: rgba(220, 53, 69, 0.1);
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
    min-height: 48px; // Prevent layout shift when delete button appears
}

.select-all {
    display: flex;
    gap: 0.5rem;
    align-items: center;
    font-weight: bold;
}

.checkbox-container {
    flex: 0 0 auto; // Prevent shrinking or growing
    display: flex;
    align-items: center;
    padding-top: 0.25rem; // Visual alignment
}

.bulk-delete-btn {
    background-color: $danger-color;
    color: white;
    border: none;
    padding: 0.5rem 1rem;
    border-radius: 4px;
    cursor: pointer;
    font-size: 0.9rem;

    &:hover {
        background-color: color.adjust($danger-color, $lightness: -10%);
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
