<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { locationsAPI } from '@/api/client'
import type { Location } from '@/types/api'

// State
const locations = ref<Location[]>([])
const loading = ref(false)
const error = ref<string | null>(null)

// Load all locations
const loadLocations = async () => {
  loading.value = true
  error.value = null

  try {
    locations.value = await locationsAPI.getAll()
  } catch (err) {
    error.value = 'Не удалось загрузить список локаций'
    console.error('Failed to load locations:', err)
  } finally {
    loading.value = false
  }
}

// Format date
const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('ru-RU', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

// Load on mount
onMounted(() => {
  loadLocations()
})
</script>

<template>
  <div>
    <!-- Page header -->
    <div class="mb-6 flex flex-col sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-800">Локации</h1>
        <p class="text-gray-600 mt-1">Управление локациями и их сетевыми настройками</p>
      </div>
    </div>

    <!-- Error message -->
    <div v-if="error" class="mb-6 bg-red-50 border border-red-200 rounded-lg p-4">
      <p class="text-red-800">{{ error }}</p>
    </div>

    <!-- Loading state -->
    <div v-if="loading" class="bg-white rounded-lg shadow-sm border border-gray-200 p-8 text-center">
      <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-600"></div>
      <p class="mt-2 text-gray-600">Загрузка...</p>
    </div>

    <!-- Locations table -->
    <div v-else class="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
      <!-- Stats -->
      <div class="px-6 py-4 bg-gray-50 border-b border-gray-200">
        <p class="text-sm text-gray-600">
          Всего локаций: <span class="font-semibold text-gray-900">{{ locations.length }}</span>
        </p>
      </div>

      <div class="overflow-x-auto">
        <table class="w-full">
          <thead class="bg-gray-50 border-b border-gray-200">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Название
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Сервер
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Подсеть
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                VoIP VLAN
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                VLAN
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Дата создания
              </th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr
              v-for="location in locations"
              :key="location.id"
              class="hover:bg-gray-50 transition-colors"
            >
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex items-center">
                  <div class="w-8 h-8 bg-green-100 rounded flex items-center justify-center mr-3">
                    <svg class="w-4 h-4 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
                    </svg>
                  </div>
                  <div class="text-sm font-medium text-gray-900">{{ location.name }}</div>
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm font-mono text-gray-900">{{ location.server }}</div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm font-mono text-gray-900">{{ location.subnet }}</div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span class="px-2 py-1 text-sm font-medium bg-purple-50 text-purple-700 rounded">
                  {{ location.voipVlan }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span class="px-2 py-1 text-sm font-medium bg-blue-50 text-blue-700 rounded">
                  {{ location.vlan }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm text-gray-600">{{ formatDate(location.createdAt) }}</div>
              </td>
            </tr>
            <tr v-if="locations.length === 0">
              <td colspan="6" class="px-6 py-8 text-center text-gray-500">
                Локации не найдены
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>
