<template>
  <div class="login-container">
    <div class="login-card">
      <div class="login-header">
        <img src="/logo.png" alt="Logo" class="login-logo" />
        <h2>Sign In</h2>
        <p class="subtitle">Access the Container Registry</p>
      </div>

      <form @submit.prevent="handleSubmit">
        <div class="form-group">
          <label for="username">Username</label>
          <input
            type="text"
            id="username"
            v-model="username"
            required
            placeholder="Enter username"
            :disabled="loading"
          />
        </div>

        <div class="form-group">
          <label for="password">Password</label>
          <input
            type="password"
            id="password"
            v-model="password"
            required
            placeholder="Enter password"
            :disabled="loading"
          />
        </div>

        <div v-if="error" class="error-alert">
          {{ error }}
        </div>

        <button type="submit" class="btn-primary" :disabled="loading">
          <span v-if="loading" class="spinner"></span>
          <span v-else>Sign In</span>
        </button>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { formatError } from '@/api/client'

const username = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')
const router = useRouter()
const authStore = useAuthStore()

const handleSubmit = async () => {
  loading.value = true
  error.value = ''
  try {
    await authStore.login({ login: username.value, password: password.value })
    router.push('/')
  } catch (err) {
    error.value = formatError(err)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped lang="scss">
@use '@/assets/main.scss' as *;
@use "sass:color";

.login-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background-color: $background-color;
  padding: 1rem;
}

.login-card {
  background-color: $card-bg;
  border: 1px solid $border-color;
  border-radius: 8px;
  padding: 1.5rem 2.5rem 2.5rem 2.5rem;
  width: 100%;
  max-width: 400px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);

  // Animation for appearance
  animation: fadeIn 0.5s ease-out;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(-20px); }
  to { opacity: 1; transform: translateY(0); }
}

.login-header {
  text-align: center;
  margin-bottom: 2rem;

  .login-logo {
    height: 120px;
    width: auto;
    margin-bottom: 1rem;
  }

  h2 {
    color: $text-color;
    margin: 0;
    font-size: 1.5rem;
    font-weight: 600;
  }

  .subtitle {
    color: $secondary-color;
    margin-top: 0.5rem;
    font-size: 0.9rem;
  }
}

.form-group {
  margin-bottom: 1.5rem;

  label {
    display: block;
    margin-bottom: 0.5rem;
    color: $text-color;
    font-weight: 500;
    font-size: 0.9rem;
  }

  input {
    width: 100%;
    padding: 0.75rem;
    border: 1px solid $border-color;
    background-color: color.adjust($background-color, $lightness: 5%);
    color: $text-color;
    border-radius: 4px;
    font-size: 1rem;
    transition: border-color 0.2s, box-shadow 0.2s;

    &:focus {
      outline: none;
      border-color: $primary-color;
      box-shadow: 0 0 0 2px rgba($primary-color, 0.2);
    }

    &:disabled {
      opacity: 0.6;
      cursor: not-allowed;
    }
  }
}

.error-alert {
  background-color: rgba($danger-color, 0.1);
  color: color.adjust($danger-color, $lightness: 10%);
  padding: 0.75rem;
  border-radius: 4px;
  margin-bottom: 1.5rem;
  font-size: 0.9rem;
  border: 1px solid rgba($danger-color, 0.2);
  display: flex;
  align-items: center;
}

.btn-primary {
  width: 100%;
  padding: 0.75rem;
  background-color: $primary-color;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.2s, transform 0.1s;
  display: flex;
  justify-content: center;
  align-items: center;
  height: 48px; // Fixed height to prevent jump on spinner

  &:hover:not(:disabled) {
    background-color: $primary-hover;
  }

  &:active:not(:disabled) {
    transform: translateY(1px);
  }

  &:disabled {
    background-color: color.adjust($primary-color, $saturation: -20%);
    opacity: 0.7;
    cursor: wait;
  }
}

.spinner {
  display: inline-block;
  width: 20px;
  height: 20px;
  border: 3px solid rgba(255, 255, 255, 0.3);
  border-radius: 50%;
  border-top-color: white;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
