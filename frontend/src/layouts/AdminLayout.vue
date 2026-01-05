<script setup lang="ts">
import { ref } from 'vue'
import Sidebar from './components/Sidebar.vue'
import Header from './components/Header.vue'
import type { SidebarState } from '@/types/navigation'

// Sidebar state management
const sidebarState = ref<SidebarState>({
  isOpen: false, // For mobile overlay
  isCollapsed: false // For desktop collapse
})

const toggleSidebar = () => {
  sidebarState.value.isCollapsed = !sidebarState.value.isCollapsed
}

const toggleMobileSidebar = () => {
  sidebarState.value.isOpen = !sidebarState.value.isOpen
}

const closeMobileSidebar = () => {
  sidebarState.value.isOpen = false
}

const handleLogout = () => {
  // In production, this would call auth store logout
  console.log('Logout clicked')
  // Redirect to login page (placeholder)
  // router.push('/login')
  alert('Logout functionality - implement auth integration')
}
</script>

<template>
  <div class="min-h-screen bg-gray-50 flex">
    <!-- Sidebar -->
    <Sidebar
      :is-open="sidebarState.isOpen"
      :is-collapsed="sidebarState.isCollapsed"
      @close="closeMobileSidebar"
      @logout="handleLogout"
    />

    <!-- Main content area -->
    <div
      :class="[
        'flex-1 flex flex-col transition-all duration-300',
        'lg:ml-0',
        sidebarState.isCollapsed ? 'lg:ml-16' : 'lg:ml-64'
      ]"
    >
      <!-- Header -->
      <Header
        :is-collapsed="sidebarState.isCollapsed"
        @toggle-sidebar="toggleSidebar"
        @toggle-mobile-sidebar="toggleMobileSidebar"
      />

      <!-- Main content -->
      <main class="flex-1 p-6 overflow-y-auto">
        <router-view />
      </main>

      <!-- Footer (optional) -->
      <footer class="bg-white border-t border-gray-200 py-4 px-6">
        <div class="flex flex-col sm:flex-row justify-between items-center text-sm text-gray-600">
          <p>&copy; 2026 Admin Panel. All rights reserved.</p>
          <div class="flex space-x-4 mt-2 sm:mt-0">
            <a href="#" class="hover:text-indigo-600 transition-colors">Privacy</a>
            <a href="#" class="hover:text-indigo-600 transition-colors">Terms</a>
            <a href="#" class="hover:text-indigo-600 transition-colors">Support</a>
          </div>
        </div>
      </footer>
    </div>
  </div>
</template>
