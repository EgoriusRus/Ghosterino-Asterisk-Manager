<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { devicesAPI } from '@/api/client'
import type { Device } from '@/types/api'

// State
const devices = ref<Device[]>([])
const loading = ref(false)
const error = ref<string | null>(null)

// Load all devices
const loadDevices = async () => {
  loading.value = true
  error.value = null

  try {
    devices.value = await devicesAPI.getAll()
  } catch (err) {
    error.value = 'Не удалось загрузить список устройств'
    console.error('Failed to load devices:', err)
  } finally {
    loading.value = false
  }
}

// Format date
const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('ru-RU', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// Format MAC address
const formatMAC = (mac: string) => {
  return mac.toUpperCase()
}

// Load on mount
onMounted(() => {
  loadDevices()
})
</script>

<template>
  <div>
    <!-- Page header -->
    <div class="mb-6 flex flex-col sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-800">Устройства</h1>
        <p class="text-gray-600 mt-1">Управление IP-телефонами и сетевыми устройствами</p>
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

    <!-- Devices table -->
    <div v-else class="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
      <!-- Stats -->
      <div class="px-6 py-4 bg-gray-50 border-b border-gray-200">
        <p class="text-sm text-gray-600">
          Всего устройств: <span class="font-semibold text-gray-900">{{ devices.length }}</span>
        </p>
      </div>

      <div class="overflow-x-auto">
        <table class="w-full">
          <thead class="bg-gray-50 border-b border-gray-200">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                MAC адрес
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Модель устройства
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Дата создания
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Последнее обновление
              </th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr
              v-for="device in devices"
              :key="device.mac"
              class="hover:bg-gray-50 transition-colors"
            >
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex items-center">
                  <div class="w-8 h-8 bg-indigo-100 rounded flex items-center justify-center mr-3">
                    <svg class="w-4 h-4 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z" />
                    </svg>
                  </div>
                  <div class="text-sm font-mono font-medium text-gray-900">
                    {{ formatMAC(device.mac) }}
                  </div>
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span class="px-2 py-1 text-sm font-medium bg-blue-50 text-blue-700 rounded">
                  {{ device.deviceModel }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm text-gray-600">{{ formatDate(device.createdAt) }}</div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm text-gray-600">{{ formatDate(device.updatedAt) }}</div>
              </td>
            </tr>
            <tr v-if="devices.length === 0">
              <td colspan="4" class="px-6 py-8 text-center text-gray-500">
                Устройства не найдены
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>
