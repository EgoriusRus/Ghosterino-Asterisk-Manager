<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import type { MenuItem } from '@/types/navigation'

interface Props {
  isOpen: boolean
}

const props = defineProps<Props>()
const emit = defineEmits<{
  close: []
}>()

const route = useRoute()

const menuItems: MenuItem[] = [
  {
    id: 'profiles',
    label: 'Сотрудники',
    icon: 'M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z',
    path: '/admin/profiles'
  },
  {
    id: 'devices',
    label: 'Устройства',
    icon: 'M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z',
    path: '/admin/devices'
  },
  {
    id: 'locations',
    label: 'Локации',
    icon: 'M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z M15 11a3 3 0 11-6 0 3 3 0 016 0z',
    path: '/admin/locations'
  }
]

const isActiveRoute = (path: string): boolean => {
  return route.path === path
}

const sidebarClasses = computed(() => {
  return [
    'fixed lg:static inset-y-0 left-0 z-30',
    'flex flex-col w-64 h-screen px-4 py-8 overflow-y-auto',
    'bg-white border-r rtl:border-r-0 rtl:border-l dark:bg-gray-900 dark:border-gray-700',
    'transition-transform duration-300 ease-in-out',
    props.isOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0'
  ]
})
</script>

<template>
  <aside :class="sidebarClasses">
    <!-- Logo/Brand -->
    <router-link to="/admin/profiles">
      <img class="w-auto h-6 sm:h-7" src="https://merakiui.com/images/logo.svg" alt="Logo">
    </router-link>

    <!-- Navigation -->
    <div class="flex flex-col justify-between flex-1 mt-6">
      <nav>
        <router-link
          v-for="(item, index) in menuItems"
          :key="item.id"
          :to="item.path"
          :class="[
            'flex items-center px-4 py-2 rounded-md transition-colors duration-300 transform',
            index > 0 ? 'mt-5' : '',
            isActiveRoute(item.path)
              ? 'text-gray-700 bg-gray-100 dark:bg-gray-800 dark:text-gray-200'
              : 'text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800 dark:hover:text-gray-200 hover:text-gray-700'
          ]"
          @click="emit('close')"
        >
          <svg
            class="w-5 h-5"
            viewBox="0 0 24 24"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              :d="item.icon"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            />
          </svg>
          <span class="mx-4 font-medium">{{ item.label }}</span>
        </router-link>
      </nav>

      <!-- User Profile -->
      <a href="#" class="flex items-center px-4 -mx-2">
        <img
          class="object-cover mx-2 rounded-full h-9 w-9"
          src="https://images.unsplash.com/photo-1531427186611-ecfd6d936c79?ixid=MXwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHw%3D&ixlib=rb-1.2.1&auto=format&fit=crop&w=634&q=80"
          alt="avatar"
        />
        <span class="mx-2 font-medium text-gray-800 dark:text-gray-200">Администратор</span>
      </a>
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
