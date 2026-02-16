<template>
  <dialog ref="dialogRef" class="modal" @click="handleBackdropClick">
    <div class="modal-content">
      <h2>{{ title }}</h2>
      <p>{{ message }}</p>

      <slot></slot>

      <div class="modal-actions">
        <button @click="cancel" class="cancel-btn">{{ cancelText }}</button>
        <button @click="confirm" class="confirm-btn" :class="{ 'danger': isDanger }">{{ confirmText }}</button>
      </div>
    </div>
  </dialog>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'

const props = defineProps({
  isOpen: {
    type: Boolean,
    required: true
  },
  title: {
    type: String,
    default: 'Confirm Action'
  },
  message: {
    type: String,
    default: 'Are you sure you want to proceed?'
  },
  confirmText: {
    type: String,
    default: 'Confirm'
  },
  cancelText: {
    type: String,
    default: 'Cancel'
  },
  isDanger: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:isOpen', 'confirm', 'cancel'])
const dialogRef = ref<HTMLDialogElement | null>(null)

watch(() => props.isOpen, (val) => {
  if (val) {
    dialogRef.value?.showModal()
  } else {
    dialogRef.value?.close()
  }
}, { immediate: true })

onMounted(() => {
    if (props.isOpen) {
        dialogRef.value?.showModal()
    }
})

const close = () => {
  emit('update:isOpen', false)
  emit('cancel')
}

const cancel = () => {
  close()
}

const confirm = () => {
  emit('confirm')
  // We don't close automatically, let parent control it (e.g. async action)
}

const handleBackdropClick = (event: MouseEvent) => {
  if (event.target === dialogRef.value) {
    close()
  }
}
</script>

<style scoped lang="scss">
@use "sass:color";
@use '@/assets/main.scss' as *;

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
        background: $primary-color;
        color: white;

        &:hover {
            background: color.adjust($primary-color, $lightness: -10%);
        }

        &.danger {
            background: $danger-color;
            &:hover {
                background: color.adjust($danger-color, $lightness: -10%);
            }
        }
    }
}
</style>
