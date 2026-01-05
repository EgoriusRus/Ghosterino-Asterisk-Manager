<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import type { MenuItem } from '@/types/navigation'

interface Props {
  isOpen: boolean
  isCollapsed: boolean
}

const props = defineProps<Props>()
const emit = defineEmits<{
  close: []
  logout: []
}>()

const route = useRoute()

const menuItems: MenuItem[] = [
  {
    id: 'dashboard',
    label: 'Dashboard',
    icon: 'M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6',
    path: '/admin/dashboard'
  },
  {
    id: 'users',
    label: 'Users',
    icon: 'M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z',
    path: '/admin/users'
  },
  {
    id: 'settings',
    label: 'Settings',
    icon: 'M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z M15 12a3 3 0 11-6 0 3 3 0 016 0z',
    path: '/admin/settings'
  }
]

const isActiveRoute = (path: string): boolean => {
  return route.path === path
}

const sidebarClasses = computed(() => {
  return [
    'fixed lg:static inset-y-0 left-0 z-30',
    'bg-gray-800 text-white',
    'transition-all duration-300 ease-in-out',
    props.isCollapsed ? 'w-16' : 'w-64',
    props.isOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0'
  ]
})
</script>

<template>
  <aside :class="sidebarClasses">
    <!-- Logo/Brand -->
    <div class="flex items-center justify-between h-16 px-4 border-b border-gray-700">
      <div class="flex items-center space-x-2">
        <div class="w-8 h-8 bg-indigo-600 rounded-lg flex items-center justify-center">
          <span class="text-white font-bold text-sm">A</span>
        </div>
        <Transition
          enter-active-class="transition-opacity duration-200"
          leave-active-class="transition-opacity duration-200"
          enter-from-class="opacity-0"
          leave-to-class="opacity-0"
        >
          <span v-if="!isCollapsed" class="font-semibold text-lg">Admin</span>
        </Transition>
      </div>
    </div>

    <!-- Navigation -->
    <nav class="flex-1 px-2 py-4 space-y-1 overflow-y-auto">
      <router-link
        v-for="item in menuItems"
        :key="item.id"
        :to="item.path"
        :class="[
          'flex items-center px-3 py-2 rounded-lg transition-colors duration-200',
          isActiveRoute(item.path)
            ? 'bg-indigo-600 text-white'
            : 'text-gray-300 hover:bg-gray-700 hover:text-white'
        ]"
        @click="emit('close')"
      >
        <svg
          class="w-6 h-6 flex-shrink-0"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            :d="item.icon"
          />
        </svg>
        <Transition
          enter-active-class="transition-opacity duration-200"
          leave-active-class="transition-opacity duration-200"
          enter-from-class="opacity-0"
          leave-to-class="opacity-0"
        >
          <span v-if="!isCollapsed" class="ml-3 font-medium">
            {{ item.label }}
          </span>
        </Transition>
      </router-link>
    </nav>

    <!-- Logout Button -->
    <div class="px-2 py-4 border-t border-gray-700">
      <button
        @click="emit('logout')"
        :class="[
          'flex items-center w-full px-3 py-2 rounded-lg transition-colors duration-200',
          'text-gray-300 hover:bg-red-600 hover:text-white'
        ]"
      >
        <svg
          class="w-6 h-6 flex-shrink-0"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
          />
        </svg>
        <Transition
          enter-active-class="transition-opacity duration-200"
          leave-active-class="transition-opacity duration-200"
          enter-from-class="opacity-0"
          leave-to-class="opacity-0"
        >
          <span v-if="!isCollapsed" class="ml-3 font-medium">Logout</span>
        </Transition>
      </button>
    </div>

    <!-- Mobile overlay -->
    <Transition
      enter-active-class="transition-opacity duration-300"
      leave-active-class="transition-opacity duration-300"
      enter-from-class="opacity-0"
      leave-to-class="opacity-0"
    >
      <div
        v-if="isOpen"
        class="fixed inset-0 bg-black bg-opacity-50 lg:hidden z-20"
        @click="emit('close')"
      />
    </Transition>
  </aside>
</template>
