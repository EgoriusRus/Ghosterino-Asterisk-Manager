<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { profilesAPI } from '@/api/client'
import type { ProfileWithLocation, PaginationResponse } from '@/types/api'

// State
const profiles = ref<ProfileWithLocation[]>([])
const pagination = ref<PaginationResponse>({
  total: 0,
  page: 1,
  perPage: 10,
  pages: 0
})
const loading = ref(false)
const error = ref<string | null>(null)

// Load profiles with pagination
const loadProfiles = async () => {
  loading.value = true
  error.value = null

  try {
    const result = await profilesAPI.getAll({
      page: pagination.value.page,
      perPage: pagination.value.perPage
    })

    profiles.value = result.data
    pagination.value = result.pagination
  } catch (err) {
    error.value = 'Не удалось загрузить список сотрудников'
    console.error('Failed to load profiles:', err)
  } finally {
    loading.value = false
  }
}

// Pagination controls
const goToPage = (page: number) => {
  pagination.value.page = page
  loadProfiles()
}

const nextPage = () => {
  if (pagination.value.page < pagination.value.pages) {
    goToPage(pagination.value.page + 1)
  }
}

const previousPage = () => {
  if (pagination.value.page > 1) {
    goToPage(pagination.value.page - 1)
  }
}

// Computed
const hasNextPage = computed(() => pagination.value.page < pagination.value.pages)
const hasPreviousPage = computed(() => pagination.value.page > 1)
const pageNumbers = computed(() => {
  const pages: number[] = []
  const totalPages = pagination.value.pages
  const currentPage = pagination.value.page

  // Show up to 5 page numbers
  let start = Math.max(1, currentPage - 2)
  let end = Math.min(totalPages, start + 4)

  if (end - start < 4) {
    start = Math.max(1, end - 4)
  }

  for (let i = start; i <= end; i++) {
    pages.push(i)
  }

  return pages
})

// Load on mount
onMounted(() => {
  loadProfiles()
})
</script>

<template>
  <div>
    <!-- Page header -->
    <div class="mb-6 flex flex-col sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-800">Сотрудники</h1>
        <p class="text-gray-600 mt-1">Управление профилями сотрудников и их SIP-настройками</p>
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

    <!-- Profiles table -->
    <div v-else class="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead class="bg-gray-50 border-b border-gray-200">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Имя
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Email
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Внутренний номер
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Локация
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Устройство
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Статус
              </th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr
              v-for="profile in profiles"
              :key="profile.id"
              class="hover:bg-gray-50 transition-colors"
            >
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm font-medium text-gray-900">{{ profile.name }}</div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm text-gray-600">{{ profile.email }}</div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm text-gray-900">{{ profile.internalNumber }}</div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm text-gray-900">
                  {{ profile.locationName || '—' }}
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm text-gray-600 font-mono">
                  {{ profile.device || '—' }}
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span
                  :class="[
                    'px-2 py-1 text-xs font-medium rounded-full',
                    profile.isActive
                      ? 'bg-green-100 text-green-800'
                      : 'bg-gray-100 text-gray-800'
                  ]"
                >
                  {{ profile.isActive ? 'Активен' : 'Неактивен' }}
                </span>
              </td>
            </tr>
            <tr v-if="profiles.length === 0">
              <td colspan="6" class="px-6 py-8 text-center text-gray-500">
                Сотрудники не найдены
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div v-if="pagination.pages > 1" class="bg-gray-50 px-6 py-4 border-t border-gray-200">
        <div class="flex items-center justify-between">
          <div class="text-sm text-gray-700">
            Показано <span class="font-medium">{{ (pagination.page - 1) * pagination.perPage + 1 }}</span>
            –
            <span class="font-medium">{{ Math.min(pagination.page * pagination.perPage, pagination.total) }}</span>
            из
            <span class="font-medium">{{ pagination.total }}</span>
            записей
          </div>

          <div class="flex items-center space-x-2">
            <!-- Previous button -->
            <button
              @click="previousPage"
              :disabled="!hasPreviousPage"
              :class="[
                'px-3 py-2 text-sm rounded-lg transition-colors',
                hasPreviousPage
                  ? 'bg-white border border-gray-300 text-gray-700 hover:bg-gray-50'
                  : 'bg-gray-100 text-gray-400 cursor-not-allowed'
              ]"
            >
              Назад
            </button>

            <!-- Page numbers -->
            <button
              v-for="page in pageNumbers"
              :key="page"
              @click="goToPage(page)"
              :class="[
                'px-3 py-2 text-sm rounded-lg transition-colors',
                pagination.page === page
                  ? 'bg-indigo-600 text-white'
                  : 'bg-white border border-gray-300 text-gray-700 hover:bg-gray-50'
              ]"
            >
              {{ page }}
            </button>

            <!-- Next button -->
            <button
              @click="nextPage"
              :disabled="!hasNextPage"
              :class="[
                'px-3 py-2 text-sm rounded-lg transition-colors',
                hasNextPage
                  ? 'bg-white border border-gray-300 text-gray-700 hover:bg-gray-50'
                  : 'bg-gray-100 text-gray-400 cursor-not-allowed'
              ]"
            >
              Вперёд
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
