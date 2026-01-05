import { ref, computed } from 'vue'
import type { User, LoginRequest, LoginResponse } from '@/types/api'

const TOKEN_KEY = 'auth_token'
const USER_KEY = 'auth_user'

// State
const token = ref<string | null>(localStorage.getItem(TOKEN_KEY))
const user = ref<User | null>(null)
const loading = ref(false)
const error = ref<string | null>(null)

// Initialize user from localStorage
const storedUser = localStorage.getItem(USER_KEY)
if (storedUser) {
  try {
    user.value = JSON.parse(storedUser)
  } catch {
    localStorage.removeItem(USER_KEY)
  }
}

// Computed
const isAuthenticated = computed(() => !!token.value && !!user.value)

// API base URL
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api'

// Actions
async function login(credentials: LoginRequest): Promise<boolean> {
  loading.value = true
  error.value = null

  try {
    const response = await fetch(`${API_BASE_URL}/auth/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(credentials),
    })

    if (!response.ok) {
      const data = await response.json()
      throw new Error(data.error || 'Ошибка авторизации')
    }

    const data: LoginResponse = await response.json()

    // Save to state and localStorage
    token.value = data.token
    user.value = data.user
    localStorage.setItem(TOKEN_KEY, data.token)
    localStorage.setItem(USER_KEY, JSON.stringify(data.user))

    return true
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Ошибка авторизации'
    return false
  } finally {
    loading.value = false
  }
}

function logout() {
  token.value = null
  user.value = null
  localStorage.removeItem(TOKEN_KEY)
  localStorage.removeItem(USER_KEY)
}

async function checkAuth(): Promise<boolean> {
  if (!token.value) {
    return false
  }

  try {
    const response = await fetch(`${API_BASE_URL}/auth/me`, {
      headers: {
        'Authorization': `Bearer ${token.value}`,
      },
    })

    if (!response.ok) {
      logout()
      return false
    }

    const userData: User = await response.json()
    user.value = userData
    localStorage.setItem(USER_KEY, JSON.stringify(userData))
    return true
  } catch {
    logout()
    return false
  }
}

function getToken(): string | null {
  return token.value
}

// Export composable
export function useAuth() {
  return {
    // State
    token,
    user,
    loading,
    error,
    // Computed
    isAuthenticated,
    // Actions
    login,
    logout,
    checkAuth,
    getToken,
  }
}
