<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { profilesAPI, devicesAPI, locationsAPI } from '@/api/client'
import type { ProfileWithLocation, PaginationResponse, Device, Location } from '@/types/api'
import BaseModal from '@/components/BaseModal.vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'

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

// Reference data for dropdowns
const locations = ref<Location[]>([])
const devices = ref<Device[]>([])
const loadingReferences = ref(false)

// Modal states
const showCreateModal = ref(false)
const showEditModal = ref(false)
const showDeleteDialog = ref(false)

// Form data
const formData = ref({
  id: 0,
  name: '',
  internalNumber: null as number | null,
  externalNumber: '',
  device: null as string | null,
  locationId: null as number | null,
  ringGroup: null as number | null,
  pickupGroup: null as number | null,
  isActive: true
})

// Form states
const formErrors = ref<Record<string, string>>({})
const formLoading = ref(false)

// Selected profile for edit/delete
const selectedProfile = ref<ProfileWithLocation | null>(null)

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

// Load reference data (locations and devices)
const loadReferenceData = async () => {
  loadingReferences.value = true

  try {
    const [locationsData, devicesData] = await Promise.all([
      locationsAPI.getAll(),
      devicesAPI.getAll()
    ])

    locations.value = locationsData
    devices.value = devicesData
  } catch (err) {
    error.value = 'Не удалось загрузить справочные данные'
    console.error('Failed to load reference data:', err)
  } finally {
    loadingReferences.value = false
  }
}

// Reset form
const resetForm = () => {
  formData.value = {
    id: 0,
    name: '',
    internalNumber: null,
    externalNumber: '',
    device: null,
    locationId: null,
    ringGroup: null,
    pickupGroup: null,
    isActive: true
  }
  formErrors.value = {}
}

// Validate form
const validateForm = (): boolean => {
  formErrors.value = {}

  if (!formData.value.name.trim()) {
    formErrors.value.name = 'ФИО обязательно'
  }

  if (!formData.value.internalNumber) {
    formErrors.value.internalNumber = 'Внутренний номер обязателен'
  } else if (formData.value.internalNumber < 0) {
    formErrors.value.internalNumber = 'Внутренний номер должен быть положительным'
  }

  if (!formData.value.locationId) {
    formErrors.value.locationId = 'Локация обязательна'
  }

  return Object.keys(formErrors.value).length === 0
}

// Open create modal
const openCreateModal = async () => {
  resetForm()
  await loadReferenceData()
  showCreateModal.value = true
}

// Open edit modal
const openEditModal = async (profile: ProfileWithLocation) => {
  selectedProfile.value = profile
  formData.value = {
    id: profile.id,
    name: profile.name,
    internalNumber: profile.internalNumber,
    externalNumber: profile.externalNumber,
    device: profile.device,
    locationId: profile.locationId,
    ringGroup: profile.ringGroup,
    pickupGroup: profile.pickupGroup,
    isActive: profile.isActive
  }
  formErrors.value = {}
  await loadReferenceData()
  showEditModal.value = true
}

// Open delete dialog
const openDeleteDialog = (profile: ProfileWithLocation) => {
  selectedProfile.value = profile
  showDeleteDialog.value = true
}

// Create profile
const createProfile = async () => {
  if (!validateForm()) {
    return
  }

  formLoading.value = true
  error.value = null

  try {
    await profilesAPI.create({
      name: formData.value.name,
      email: '',
      internalNumber: formData.value.internalNumber!,
      externalNumber: formData.value.externalNumber,
      device: formData.value.device,
      locationId: formData.value.locationId,
      ringGroup: formData.value.ringGroup,
      pickupGroup: formData.value.pickupGroup,
      isActive: formData.value.isActive
    })

    showCreateModal.value = false
    resetForm()
    await loadProfiles()
  } catch (err) {
    error.value = 'Не удалось создать сотрудника'
    console.error('Failed to create profile:', err)
  } finally {
    formLoading.value = false
  }
}

// Update profile
const updateProfile = async () => {
  if (!validateForm()) {
    return
  }

  formLoading.value = true
  error.value = null

  try {
    await profilesAPI.update(formData.value.id, {
      name: formData.value.name,
      internalNumber: formData.value.internalNumber!,
      externalNumber: formData.value.externalNumber,
      device: formData.value.device,
      locationId: formData.value.locationId,
      ringGroup: formData.value.ringGroup,
      pickupGroup: formData.value.pickupGroup,
      isActive: formData.value.isActive
    })

    showEditModal.value = false
    resetForm()
    await loadProfiles()
  } catch (err) {
    error.value = 'Не удалось обновить сотрудника'
    console.error('Failed to update profile:', err)
  } finally {
    formLoading.value = false
  }
}

// Delete profile
const deleteProfile = async () => {
  if (!selectedProfile.value) return

  formLoading.value = true
  error.value = null

  try {
    await profilesAPI.delete(selectedProfile.value.id)
    showDeleteDialog.value = false
    selectedProfile.value = null
    await loadProfiles()
  } catch (err) {
    error.value = 'Не удалось удалить сотрудника'
    console.error('Failed to delete profile:', err)
  } finally {
    formLoading.value = false
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
      <div class="mt-4 sm:mt-0">
        <button
          @click="openCreateModal"
          class="inline-flex items-center px-4 py-2 bg-indigo-600 text-white text-sm font-medium rounded-lg hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-colors"
        >
          <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          Добавить сотрудника
        </button>
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
              <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                Действия
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
              <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                <button
                  @click="openEditModal(profile)"
                  class="text-indigo-600 hover:text-indigo-900 mr-3"
                  title="Редактировать"
                >
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                  </svg>
                </button>
                <button
                  @click="openDeleteDialog(profile)"
                  class="text-red-600 hover:text-red-900"
                  title="Удалить"
                >
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  </svg>
                </button>
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

    <!-- Create Modal -->
    <BaseModal
      :show="showCreateModal"
      title="Создать сотрудника"
      size="lg"
      @close="showCreateModal = false"
    >
      <form @submit.prevent="createProfile" class="space-y-4">
        <div v-if="loadingReferences" class="text-center py-4">
          <div class="inline-block animate-spin rounded-full h-6 w-6 border-b-2 border-indigo-600"></div>
          <p class="mt-2 text-sm text-gray-600">Загрузка данных...</p>
        </div>

        <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <!-- Name -->
          <div class="md:col-span-2">
            <label for="create-name" class="block text-sm font-medium text-gray-700 mb-1">
              ФИО <span class="text-red-500">*</span>
            </label>
            <input
              id="create-name"
              v-model="formData.name"
              type="text"
              :class="[
                'w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500',
                formErrors.name ? 'border-red-300' : 'border-gray-300'
              ]"
              placeholder="Иванов Иван Иванович"
            />
            <p v-if="formErrors.name" class="mt-1 text-sm text-red-600">{{ formErrors.name }}</p>
          </div>

          <!-- Internal Number -->
          <div>
            <label for="create-internal" class="block text-sm font-medium text-gray-700 mb-1">
              Внутренний номер <span class="text-red-500">*</span>
            </label>
            <input
              id="create-internal"
              v-model.number="formData.internalNumber"
              type="number"
              :class="[
                'w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500',
                formErrors.internalNumber ? 'border-red-300' : 'border-gray-300'
              ]"
              placeholder="1001"
            />
            <p v-if="formErrors.internalNumber" class="mt-1 text-sm text-red-600">{{ formErrors.internalNumber }}</p>
          </div>

          <!-- External Number -->
          <div>
            <label for="create-external" class="block text-sm font-medium text-gray-700 mb-1">
              Внешний номер
            </label>
            <input
              id="create-external"
              v-model="formData.externalNumber"
              type="text"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500"
              placeholder="+7 (123) 456-78-90"
            />
          </div>

          <!-- Location -->
          <div>
            <label for="create-location" class="block text-sm font-medium text-gray-700 mb-1">
              Локация <span class="text-red-500">*</span>
            </label>
            <select
              id="create-location"
              v-model="formData.locationId"
              :class="[
                'w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500',
                formErrors.locationId ? 'border-red-300' : 'border-gray-300'
              ]"
            >
              <option :value="null">Выберите локацию</option>
              <option v-for="location in locations" :key="location.id" :value="location.id">
                {{ location.name }}
              </option>
            </select>
            <p v-if="formErrors.locationId" class="mt-1 text-sm text-red-600">{{ formErrors.locationId }}</p>
          </div>

          <!-- Device -->
          <div>
            <label for="create-device" class="block text-sm font-medium text-gray-700 mb-1">
              Устройство (MAC)
            </label>
            <select
              id="create-device"
              v-model="formData.device"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 font-mono text-sm"
            >
              <option :value="null">Не назначено</option>
              <option v-for="device in devices" :key="device.mac" :value="device.mac">
                {{ device.mac }} ({{ device.deviceModel }})
              </option>
            </select>
          </div>

          <!-- Ring Group -->
          <div>
            <label for="create-ring" class="block text-sm font-medium text-gray-700 mb-1">
              Ring Group
            </label>
            <input
              id="create-ring"
              v-model.number="formData.ringGroup"
              type="number"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500"
              placeholder="1"
            />
          </div>

          <!-- Pickup Group -->
          <div>
            <label for="create-pickup" class="block text-sm font-medium text-gray-700 mb-1">
              Pickup Group
            </label>
            <input
              id="create-pickup"
              v-model.number="formData.pickupGroup"
              type="number"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500"
              placeholder="1"
            />
          </div>

          <!-- Is Active -->
          <div class="md:col-span-2">
            <label class="flex items-center">
              <input
                v-model="formData.isActive"
                type="checkbox"
                class="w-4 h-4 text-indigo-600 border-gray-300 rounded focus:ring-indigo-500"
              />
              <span class="ml-2 text-sm text-gray-700">Активен</span>
            </label>
          </div>
        </div>
      </form>

      <template #footer>
        <div class="flex justify-end gap-3">
          <button
            @click="showCreateModal = false"
            :disabled="formLoading"
            class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500 disabled:opacity-50 disabled:cursor-not-allowed"
            type="button"
          >
            Отмена
          </button>
          <button
            @click="createProfile"
            :disabled="formLoading || loadingReferences"
            class="px-4 py-2 text-sm font-medium text-white bg-indigo-600 rounded-lg hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed"
            type="button"
          >
            <span v-if="formLoading" class="flex items-center">
              <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              Создание...
            </span>
            <span v-else>Создать</span>
          </button>
        </div>
      </template>
    </BaseModal>

    <!-- Edit Modal -->
    <BaseModal
      :show="showEditModal"
      title="Редактировать сотрудника"
      size="lg"
      @close="showEditModal = false"
    >
      <form @submit.prevent="updateProfile" class="space-y-4">
        <div v-if="loadingReferences" class="text-center py-4">
          <div class="inline-block animate-spin rounded-full h-6 w-6 border-b-2 border-indigo-600"></div>
          <p class="mt-2 text-sm text-gray-600">Загрузка данных...</p>
        </div>

        <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <!-- Name -->
          <div class="md:col-span-2">
            <label for="edit-name" class="block text-sm font-medium text-gray-700 mb-1">
              ФИО <span class="text-red-500">*</span>
            </label>
            <input
              id="edit-name"
              v-model="formData.name"
              type="text"
              :class="[
                'w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500',
                formErrors.name ? 'border-red-300' : 'border-gray-300'
              ]"
              placeholder="Иванов Иван Иванович"
            />
            <p v-if="formErrors.name" class="mt-1 text-sm text-red-600">{{ formErrors.name }}</p>
          </div>

          <!-- Internal Number -->
          <div>
            <label for="edit-internal" class="block text-sm font-medium text-gray-700 mb-1">
              Внутренний номер <span class="text-red-500">*</span>
            </label>
            <input
              id="edit-internal"
              v-model.number="formData.internalNumber"
              type="number"
              :class="[
                'w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500',
                formErrors.internalNumber ? 'border-red-300' : 'border-gray-300'
              ]"
              placeholder="1001"
            />
            <p v-if="formErrors.internalNumber" class="mt-1 text-sm text-red-600">{{ formErrors.internalNumber }}</p>
          </div>

          <!-- External Number -->
          <div>
            <label for="edit-external" class="block text-sm font-medium text-gray-700 mb-1">
              Внешний номер
            </label>
            <input
              id="edit-external"
              v-model="formData.externalNumber"
              type="text"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500"
              placeholder="+7 (123) 456-78-90"
            />
          </div>

          <!-- Location -->
          <div>
            <label for="edit-location" class="block text-sm font-medium text-gray-700 mb-1">
              Локация <span class="text-red-500">*</span>
            </label>
            <select
              id="edit-location"
              v-model="formData.locationId"
              :class="[
                'w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500',
                formErrors.locationId ? 'border-red-300' : 'border-gray-300'
              ]"
            >
              <option :value="null">Выберите локацию</option>
              <option v-for="location in locations" :key="location.id" :value="location.id">
                {{ location.name }}
              </option>
            </select>
            <p v-if="formErrors.locationId" class="mt-1 text-sm text-red-600">{{ formErrors.locationId }}</p>
          </div>

          <!-- Device -->
          <div>
            <label for="edit-device" class="block text-sm font-medium text-gray-700 mb-1">
              Устройство (MAC)
            </label>
            <select
              id="edit-device"
              v-model="formData.device"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 font-mono text-sm"
            >
              <option :value="null">Не назначено</option>
              <option v-for="device in devices" :key="device.mac" :value="device.mac">
                {{ device.mac }} ({{ device.deviceModel }})
              </option>
            </select>
          </div>

          <!-- Ring Group -->
          <div>
            <label for="edit-ring" class="block text-sm font-medium text-gray-700 mb-1">
              Ring Group
            </label>
            <input
              id="edit-ring"
              v-model.number="formData.ringGroup"
              type="number"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500"
              placeholder="1"
            />
          </div>

          <!-- Pickup Group -->
          <div>
            <label for="edit-pickup" class="block text-sm font-medium text-gray-700 mb-1">
              Pickup Group
            </label>
            <input
              id="edit-pickup"
              v-model.number="formData.pickupGroup"
              type="number"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500"
              placeholder="1"
            />
          </div>

          <!-- Is Active -->
          <div class="md:col-span-2">
            <label class="flex items-center">
              <input
                v-model="formData.isActive"
                type="checkbox"
                class="w-4 h-4 text-indigo-600 border-gray-300 rounded focus:ring-indigo-500"
              />
              <span class="ml-2 text-sm text-gray-700">Активен</span>
            </label>
          </div>
        </div>
      </form>

      <template #footer>
        <div class="flex justify-end gap-3">
          <button
            @click="showEditModal = false"
            :disabled="formLoading"
            class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500 disabled:opacity-50 disabled:cursor-not-allowed"
            type="button"
          >
            Отмена
          </button>
          <button
            @click="updateProfile"
            :disabled="formLoading || loadingReferences"
            class="px-4 py-2 text-sm font-medium text-white bg-indigo-600 rounded-lg hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed"
            type="button"
          >
            <span v-if="formLoading" class="flex items-center">
              <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              Сохранение...
            </span>
            <span v-else>Сохранить</span>
          </button>
        </div>
      </template>
    </BaseModal>

    <!-- Delete Confirmation Dialog -->
    <ConfirmDialog
      :show="showDeleteDialog"
      title="Удалить сотрудника?"
      :message="`Вы уверены, что хотите удалить сотрудника &quot;${selectedProfile?.name}&quot;? Это действие нельзя отменить.`"
      confirm-text="Удалить"
      cancel-text="Отмена"
      variant="danger"
      :loading="formLoading"
      @confirm="deleteProfile"
      @cancel="showDeleteDialog = false"
    />
  </div>
</template>
