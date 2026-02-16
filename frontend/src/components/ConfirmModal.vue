<template>
  <dialog ref="dialogRef" class="modal" @click="handleBackdropClick">
    <div class="modal-content">
      <h2>{{ title }}</h2>
      <p>{{ message }}</p>

      <slot></slot>

      <div v-if="verificationValue" class="verification-group">
        <label>
            Please type <strong>{{ verificationValue }}</strong> to confirm.
        </label>
        <input
            type="text"
            v-model="inputValue"
            class="verification-input"
            :placeholder="verificationValue"
            @keyup.enter="isConfirmed ? confirm() : null"
        />
      </div>

      <div class="modal-actions">
        <button @click="cancel" class="cancel-btn">{{ cancelText }}</button>
        <button
            @click="confirm"
            class="confirm-btn"
            :class="{ 'danger': isDanger }"
            :disabled="!isConfirmed"
        >{{ confirmText }}</button>
      </div>
    </div>
  </dialog>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, computed } from 'vue'

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
  },
  verificationValue: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:isOpen', 'confirm', 'cancel'])
const dialogRef = ref<HTMLDialogElement | null>(null)
const inputValue = ref('')

const isConfirmed = computed(() => {
    if (!props.verificationValue) return true
    return inputValue.value === props.verificationValue
})

watch(() => props.isOpen, (val) => {
  if (val) {
    inputValue.value = '' // Reset input on open
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
  if (isConfirmed.value) {
      emit('confirm')
  }
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

.verification-group {
    margin: 1.5rem 0;
    display: flex;
    flex-direction: column;
    gap: 0.5rem;

    label {
        font-size: 0.9rem;
    }

    strong {
        user-select: all; // Easy copy-paste
    }

    .verification-input {
        padding: 0.5rem;
        border-radius: 4px;
        border: 1px solid $border-color;
        background: $background-color;
        color: $text-color;

        &:focus {
            outline: 2px solid $primary-color;
            border-color: transparent;
        }
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
        transition: opacity 0.2s, background-color 0.2s;

        &:disabled {
            opacity: 0.5;
            cursor: not-allowed;
        }
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

        &:hover:not(:disabled) {
            background: color.adjust($primary-color, $lightness: -10%);
        }

        &.danger {
            background: $danger-color;
            &:hover:not(:disabled) {
                background: color.adjust($danger-color, $lightness: -10%);
            }
        }
    }
}
</style>
