<script setup lang="ts">
import { ref } from 'vue'
import type { UserInfo } from '@/types/navigation'

interface Props {
  isCollapsed: boolean
}

defineProps<Props>()
const emit = defineEmits<{
  toggleSidebar: []
  toggleMobileSidebar: []
}>()

// Mock user data - in production this would come from auth store
const user = ref<UserInfo>({
  id: '1',
  name: 'Admin User',
  email: 'admin@example.com',
  role: 'Administrator'
})
</script>

<template>
  <header class="bg-white shadow-sm border-b border-gray-200 sticky top-0 z-10">
    <div class="flex items-center justify-between h-16 px-4">
      <!-- Left: Toggle buttons -->
      <div class="flex items-center space-x-2">
        <!-- Mobile hamburger -->
        <button
          @click="emit('toggleMobileSidebar')"
          class="lg:hidden p-2 rounded-lg text-gray-600 hover:bg-gray-100 transition-colors"
          aria-label="Toggle mobile menu"
        >
          <svg
            class="w-6 h-6"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M4 6h16M4 12h16M4 18h16"
            />
          </svg>
        </button>

        <!-- Desktop toggle -->
        <button
          @click="emit('toggleSidebar')"
          class="hidden lg:block p-2 rounded-lg text-gray-600 hover:bg-gray-100 transition-colors"
          aria-label="Toggle sidebar"
        >
          <svg
            class="w-6 h-6"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M4 6h16M4 12h16M4 18h16"
            />
          </svg>
        </button>

        <!-- Page title or breadcrumb can go here -->
        <h1 class="hidden md:block text-xl font-semibold text-gray-800">
          Admin Panel
        </h1>
      </div>

      <!-- Right: User info -->
      <div class="flex items-center space-x-4">
        <!-- Notifications (placeholder) -->
        <button
          class="p-2 rounded-lg text-gray-600 hover:bg-gray-100 transition-colors relative"
          aria-label="Notifications"
        >
          <svg
            class="w-6 h-6"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"
            />
          </svg>
          <!-- Notification badge -->
          <span
            class="absolute top-1 right-1 w-2 h-2 bg-red-500 rounded-full"
          ></span>
        </button>

        <!-- User profile -->
        <div class="flex items-center space-x-3">
          <div class="hidden sm:block text-right">
            <p class="text-sm font-medium text-gray-800">{{ user.name }}</p>
            <p class="text-xs text-gray-500">{{ user.role }}</p>
          </div>
          <div
            class="w-10 h-10 rounded-full bg-indigo-600 flex items-center justify-center text-white font-semibold"
          >
            {{ user.name.charAt(0) }}
          </div>
        </div>
      </div>
    </div>
  </header>
</template>
