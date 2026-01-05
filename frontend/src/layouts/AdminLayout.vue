<script setup lang="ts">
import { ref } from 'vue'
import Sidebar from './components/Sidebar.vue'

// Sidebar state management (mobile only)
const isMobileSidebarOpen = ref(false)

const toggleMobileSidebar = () => {
  isMobileSidebarOpen.value = !isMobileSidebarOpen.value
}

const closeMobileSidebar = () => {
  isMobileSidebarOpen.value = false
}
</script>

<template>
  <div class="min-h-screen bg-gray-50 flex">
    <!-- Sidebar -->
    <Sidebar
      :is-open="isMobileSidebarOpen"
      @close="closeMobileSidebar"
    />

    <!-- Main content area -->
    <div class="flex-1 flex flex-col">
      <!-- Mobile menu toggle button (visible only on mobile) -->
      <div class="lg:hidden bg-white border-b border-gray-200 px-4 py-3 flex items-center">
        <button
          @click="toggleMobileSidebar"
          class="p-2 rounded-lg text-gray-600 hover:bg-gray-100 transition-colors"
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
        <h1 class="ml-3 text-lg font-semibold text-gray-800">Admin Panel</h1>
      </div>

      <!-- Main content -->
      <main class="flex-1 p-6 overflow-y-auto">
        <router-view />
      </main>
    </div>
  </div>
</template>
