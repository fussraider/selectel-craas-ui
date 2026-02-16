<template>
  <div class="registry-settings">
    <div class="header">
      <h1>Registry Settings: {{ rid }}</h1>
    </div>

    <div v-if="store.error" class="error-msg">{{ store.error }}</div>
    <div v-if="store.success" class="success-msg">{{ store.success }}</div>

    <div class="card info-card">
        <h3>Registry Information</h3>
        <div class="info-grid">
            <div class="info-item">
                <span class="label">ID:</span>
                <span class="value">{{ rid }}</span>
            </div>
             <!-- We might need to find the registry object in store to get status/size/created -->
             <div v-if="registry" class="info-item">
                <span class="label">Name:</span>
                <span class="value">{{ registry.name }}</span>
            </div>
             <div v-if="registry" class="info-item">
                <span class="label">Status:</span>
                <span class="value status" :class="registry.status.toLowerCase()">{{ registry.status }}</span>
            </div>
             <div v-if="registry" class="info-item">
                <span class="label">Size:</span>
                <span class="value">{{ (registry.size / 1024 / 1024).toFixed(2) }} MB</span>
            </div>
             <div v-if="registry" class="info-item">
                <span class="label">Created:</span>
                <span class="value">{{ new Date(registry.createdAt).toLocaleDateString() }}</span>
            </div>
        </div>
    </div>

    <div class="card gc-card">
        <h3>Garbage Collection</h3>
        <p class="description">
            Garbage collection frees up space by removing unreferenced blobs.
            This operation puts the registry in read-only mode.
        </p>

        <div v-if="store.gcLoading" class="loading">Loading GC info...</div>
        <div v-else-if="store.gcInfo" class="gc-stats">
            <div class="stat">
                <span class="stat-label">Potential Savings</span>
                <span class="stat-value">{{ (store.gcInfo.sizeSummary / 1024 / 1024).toFixed(2) }} MB</span>
            </div>
             <div class="stat">
                <span class="stat-label">Untagged Size</span>
                <span class="stat-value">{{ (store.gcInfo.sizeUntagged / 1024 / 1024).toFixed(2) }} MB</span>
            </div>
        </div>

        <button @click="triggerGC" :disabled="store.gcLoading" class="btn primary">
            Start Garbage Collection
        </button>
    </div>

    <div class="card danger-zone">
        <h3>Danger Zone</h3>
        <div class="action-row">
            <div class="action-info">
                <h4>Delete Registry</h4>
                <p>This action cannot be undone. All repositories and images will be lost.</p>
            </div>
            <button @click="deleteReg" :disabled="store.loading" class="btn danger">
                Delete Registry
            </button>
        </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRegistryStore } from '@/stores/registry'
import { onMounted, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const store = useRegistryStore()

const pid = computed(() => route.params.pid as string)
const rid = computed(() => route.params.rid as string)

const registry = computed(() => store.registries.find(r => r.id === rid.value))

onMounted(() => {
    store.clearNotifications()
    if (pid.value && rid.value) {
        store.fetchGCInfo(pid.value, rid.value)
    }
})

watch(() => rid.value, (newRid) => {
    if (newRid && pid.value) {
        store.fetchGCInfo(pid.value, newRid)
    }
})

const triggerGC = async () => {
    if (confirm("Garbage collection makes the registry read-only until it completes. Are you sure you want to proceed?")) {
        await store.startGC(pid.value, rid.value)
    }
}

const deleteReg = async () => {
  if (confirm('Are you sure you want to delete this registry? This action is irreversible.')) {
    await store.deleteRegistry(pid.value, rid.value)
    // Redirect to project root or dashboard
    router.push('/')
  }
}
</script>

<style scoped lang="scss">
@use '@/assets/main.scss' as *;

.registry-settings {
    max-width: 800px;
    margin: 0 auto;
}

.header {
    margin-bottom: 2rem;
    h1 {
        font-size: 1.75rem;
        color: $text-color;
    }
}

.card {
    background: $card-bg;
    border: 1px solid $border-color;
    border-radius: 8px;
    padding: 1.5rem;
    margin-bottom: 2rem;

    h3 {
        margin-top: 0;
        margin-bottom: 1rem;
        font-size: 1.25rem;
        color: $primary-color;
        border-bottom: 1px solid $border-color;
        padding-bottom: 0.5rem;
    }
}

.info-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 1rem;
}

.info-item {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;

    .label {
        font-size: 0.85rem;
        color: $secondary-color;
        font-weight: 500;
    }

    .value {
        font-size: 1rem;
        color: $text-color;

        &.status {
             display: inline-block;
             padding: 0.1rem 0.4rem;
             border-radius: 4px;
             font-size: 0.8rem;
             font-weight: bold;
             text-transform: uppercase;
             width: fit-content;

            &.active { background-color: #198754; color: #fff; }
            &.creating { background-color: #ffc107; color: #212529; }
            &.error { background-color: #dc3545; color: #fff; }
        }
    }
}

.description {
    color: $secondary-color;
    font-size: 0.95rem;
    margin-bottom: 1rem;
}

.gc-stats {
    display: flex;
    gap: 2rem;
    margin-bottom: 1.5rem;
    background: rgba($background-color, 0.5);
    padding: 1rem;
    border-radius: 6px;

    .stat {
        display: flex;
        flex-direction: column;

        .stat-label {
            font-size: 0.8rem;
            color: $secondary-color;
        }

        .stat-value {
            font-size: 1.25rem;
            font-weight: 600;
            color: $primary-color;
        }
    }
}

.danger-zone {
    border-color: rgba($danger-color, 0.3);

    h3 {
        color: $danger-color;
        border-color: rgba($danger-color, 0.2);
    }
}

.action-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 1rem;

    @media (max-width: 600px) {
        flex-direction: column;
        align-items: flex-start;

        button {
            width: 100%;
        }
    }
}

.action-info {
    h4 { margin: 0 0 0.5rem 0; font-size: 1rem; }
    p { margin: 0; font-size: 0.9rem; color: $secondary-color; }
}

.btn {
    padding: 0.6rem 1.2rem;
    border-radius: 4px;
    border: none;
    font-weight: 500;
    cursor: pointer;
    transition: opacity 0.2s;

    &:hover:not(:disabled) { opacity: 0.9; }
    &:disabled { opacity: 0.5; cursor: not-allowed; }

    &.primary {
        background-color: $primary-color;
        color: white;
    }

    &.danger {
        background-color: $danger-color;
        color: white;
    }
}

.error-msg {
    padding: 1rem;
    margin-bottom: 1rem;
    background-color: rgba($danger-color, 0.1);
    color: $danger-color;
    border-radius: 4px;
}

.success-msg {
    padding: 1rem;
    margin-bottom: 1rem;
    background-color: rgba(#198754, 0.1);
    color: #198754;
    border-radius: 4px;
}
</style>
