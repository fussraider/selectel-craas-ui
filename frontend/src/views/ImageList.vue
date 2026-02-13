<template>
  <div class="container">
    <div class="header">
      <router-link to="/">Projects</router-link> &gt;
      <router-link :to="`/projects/${pid}/registries`">Registries</router-link> &gt;
      <router-link :to="`/projects/${pid}/registries/${rid}/repositories`">Repositories</router-link> &gt;
      Images
      <h1>Images in {{ rname }}</h1>
    </div>
    <div v-if="store.loading" class="loading">Loading...</div>
    <div v-else-if="store.error" class="error">{{ store.error }}</div>
    <div v-else class="image-list">
      <div v-if="store.images.length === 0">No images found.</div>
      <div v-for="image in store.images" :key="image.digest" class="image-item">
        <div class="image-info">
          <div class="digest">{{ image.digest }}</div>
          <div class="tags">
            <span v-for="tag in image.tags" :key="tag" class="tag">{{ tag }}</span>
          </div>
          <div class="image-details">Size: {{ image.size }} bytes | Created: {{ image.createdAt }}</div>
        </div>
        <button @click="deleteImg(image.digest)" class="delete-btn">Delete</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRegistryStore } from '@/stores/registry'
import { onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const store = useRegistryStore()
const pid = computed(() => route.params.pid as string)
const rid = computed(() => route.params.rid as string)
const rname = computed(() => route.params.rname as string)

onMounted(() => {
  store.fetchImages(pid.value, rid.value, rname.value)
})

const deleteImg = async (digest: string) => {
  if (confirm(`Are you sure you want to delete image ${digest}?`)) {
    await store.deleteImage(pid.value, rid.value, rname.value, digest)
  }
}
</script>

<style scoped>
.container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}
.image-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.image-item {
  border: 1px solid #ddd;
  padding: 15px;
  border-radius: 4px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.digest {
  font-family: monospace;
  font-weight: bold;
}
.tags {
  display: flex;
  gap: 5px;
  margin: 5px 0;
}
.tag {
  background-color: #e9ecef;
  padding: 2px 5px;
  border-radius: 3px;
  font-size: 0.9em;
}
.image-details {
  font-size: 0.9em;
  color: #666;
}
.delete-btn {
  background-color: #dc3545;
  color: white;
  border: none;
  padding: 5px 10px;
  border-radius: 3px;
  cursor: pointer;
}
.delete-btn:hover {
  background-color: #c82333;
}
.error {
  color: red;
}
</style>
