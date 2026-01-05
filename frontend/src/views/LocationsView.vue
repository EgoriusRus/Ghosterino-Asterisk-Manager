<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { locationsAPI } from '@/api/client'
import type { Location } from '@/types/api'
import BaseModal from '@/components/BaseModal.vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'

// State
const locations = ref<Location[]>([])
const loading = ref(false)
const error = ref<string | null>(null)

// Modal states
const showCreateModal = ref(false)
const showEditModal = ref(false)
const showDeleteDialog = ref(false)

// Form data
const formData = ref({
  id: 0,
  name: '',
  server: '',
  subnet: '',
  voipVlan: null as number | null,
  vlan: null as number | null
})

// Form states
const formErrors = ref<Record<string, string>>({})
const formLoading = ref(false)

// Selected location for edit/delete
const selectedLocation = ref<Location | null>(null)

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

// Reset form
const resetForm = () => {
  formData.value = {
    id: 0,
    name: '',
    server: '',
    subnet: '',
    voipVlan: null,
    vlan: null
  }
  formErrors.value = {}
}

// Validate form
const validateForm = (): boolean => {
  formErrors.value = {}

  if (!formData.value.name.trim()) {
    formErrors.value.name = 'Название обязательно'
  }

  if (!formData.value.server.trim()) {
    formErrors.value.server = 'IP сервера обязателен'
  } else {
    // Basic IP validation
    const ipRegex = /^(\d{1,3}\.){3}\d{1,3}$/
    if (!ipRegex.test(formData.value.server)) {
      formErrors.value.server = 'Некорректный формат IP адреса'
    }
  }

  return Object.keys(formErrors.value).length === 0
}

// Open create modal
const openCreateModal = () => {
  resetForm()
  showCreateModal.value = true
}

// Open edit modal
const openEditModal = (location: Location) => {
  selectedLocation.value = location
  formData.value = {
    id: location.id,
    name: location.name,
    server: location.server,
    subnet: location.subnet,
    voipVlan: location.voipVlan,
    vlan: location.vlan
  }
  formErrors.value = {}
  showEditModal.value = true
}

// Open delete dialog
const openDeleteDialog = (location: Location) => {
  selectedLocation.value = location
  showDeleteDialog.value = true
}

// Create location
const createLocation = async () => {
  if (!validateForm()) {
    return
  }

  formLoading.value = true
  error.value = null

  try {
    await locationsAPI.create({
      name: formData.value.name,
      server: formData.value.server,
      subnet: formData.value.subnet,
      voipVlan: formData.value.voipVlan || 0,
      vlan: formData.value.vlan || 0
    })

    showCreateModal.value = false
    resetForm()
    await loadLocations()
  } catch (err) {
    error.value = 'Не удалось создать локацию'
    console.error('Failed to create location:', err)
  } finally {
    formLoading.value = false
  }
}

// Update location
const updateLocation = async () => {
  if (!validateForm()) {
    return
  }

  formLoading.value = true
  error.value = null

  try {
    await locationsAPI.update(formData.value.id, {
      name: formData.value.name,
      server: formData.value.server,
      subnet: formData.value.subnet,
      voipVlan: formData.value.voipVlan || 0,
      vlan: formData.value.vlan || 0
    })

    showEditModal.value = false
    resetForm()
    await loadLocations()
  } catch (err) {
    error.value = 'Не удалось обновить локацию'
    console.error('Failed to update location:', err)
  } finally {
    formLoading.value = false
  }
}

// Delete location
const deleteLocation = async () => {
  if (!selectedLocation.value) return

  formLoading.value = true
  error.value = null

  try {
    await locationsAPI.delete(selectedLocation.value.id)
    showDeleteDialog.value = false
    selectedLocation.value = null
    await loadLocations()
  } catch (err) {
    error.value = 'Не удалось удалить локацию'
    console.error('Failed to delete location:', err)
  } finally {
    formLoading.value = false
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
      <div class="mt-4 sm:mt-0">
        <button
          @click="openCreateModal"
          class="inline-flex items-center px-4 py-2 bg-indigo-600 text-white text-sm font-medium rounded-lg hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-colors"
        >
          <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          Добавить локацию
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
              <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                Действия
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
              <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                <button
                  @click="openEditModal(location)"
                  class="text-indigo-600 hover:text-indigo-900 mr-3"
                  title="Редактировать"
                >
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                  </svg>
                </button>
                <button
                  @click="openDeleteDialog(location)"
                  class="text-red-600 hover:text-red-900"
                  title="Удалить"
                >
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  </svg>
                </button>
              </td>
            </tr>
            <tr v-if="locations.length === 0">
              <td colspan="7" class="px-6 py-8 text-center text-gray-500">
                Локации не найдены
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Create Modal -->
    <BaseModal
      :show="showCreateModal"
      title="Создать локацию"
      @close="showCreateModal = false"
    >
      <form @submit.prevent="createLocation" class="space-y-4">
        <!-- Name -->
        <div>
          <label for="create-name" class="block text-sm font-medium text-gray-700 mb-1">
            Название <span class="text-red-500">*</span>
          </label>
          <input
            id="create-name"
            v-model="formData.name"
            type="text"
            :class="[
              'w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500',
              formErrors.name ? 'border-red-300' : 'border-gray-300'
            ]"
            placeholder="Например: Офис Москва"
          />
          <p v-if="formErrors.name" class="mt-1 text-sm text-red-600">{{ formErrors.name }}</p>
        </div>

        <!-- Server -->
        <div>
          <label for="create-server" class="block text-sm font-medium text-gray-700 mb-1">
            IP сервера <span class="text-red-500">*</span>
          </label>
          <input
            id="create-server"
            v-model="formData.server"
            type="text"
            :class="[
              'w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500',
              formErrors.server ? 'border-red-300' : 'border-gray-300'
            ]"
            placeholder="192.168.1.1"
          />
          <p v-if="formErrors.server" class="mt-1 text-sm text-red-600">{{ formErrors.server }}</p>
        </div>

        <!-- Subnet -->
        <div>
          <label for="create-subnet" class="block text-sm font-medium text-gray-700 mb-1">
            Подсеть
          </label>
          <input
            id="create-subnet"
            v-model="formData.subnet"
            type="text"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500"
            placeholder="192.168.1.0/24"
          />
        </div>

        <!-- VoIP VLAN -->
        <div>
          <label for="create-voip-vlan" class="block text-sm font-medium text-gray-700 mb-1">
            VoIP VLAN
          </label>
          <input
            id="create-voip-vlan"
            v-model.number="formData.voipVlan"
            type="number"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500"
            placeholder="100"
          />
        </div>

        <!-- VLAN -->
        <div>
          <label for="create-vlan" class="block text-sm font-medium text-gray-700 mb-1">
            VLAN
          </label>
          <input
            id="create-vlan"
            v-model.number="formData.vlan"
            type="number"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500"
            placeholder="10"
          />
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
            @click="createLocation"
            :disabled="formLoading"
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
      title="Редактировать локацию"
      @close="showEditModal = false"
    >
      <form @submit.prevent="updateLocation" class="space-y-4">
        <!-- Name -->
        <div>
          <label for="edit-name" class="block text-sm font-medium text-gray-700 mb-1">
            Название <span class="text-red-500">*</span>
          </label>
          <input
            id="edit-name"
            v-model="formData.name"
            type="text"
            :class="[
              'w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500',
              formErrors.name ? 'border-red-300' : 'border-gray-300'
            ]"
            placeholder="Например: Офис Москва"
          />
          <p v-if="formErrors.name" class="mt-1 text-sm text-red-600">{{ formErrors.name }}</p>
        </div>

        <!-- Server -->
        <div>
          <label for="edit-server" class="block text-sm font-medium text-gray-700 mb-1">
            IP сервера <span class="text-red-500">*</span>
          </label>
          <input
            id="edit-server"
            v-model="formData.server"
            type="text"
            :class="[
              'w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500',
              formErrors.server ? 'border-red-300' : 'border-gray-300'
            ]"
            placeholder="192.168.1.1"
          />
          <p v-if="formErrors.server" class="mt-1 text-sm text-red-600">{{ formErrors.server }}</p>
        </div>

        <!-- Subnet -->
        <div>
          <label for="edit-subnet" class="block text-sm font-medium text-gray-700 mb-1">
            Подсеть
          </label>
          <input
            id="edit-subnet"
            v-model="formData.subnet"
            type="text"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500"
            placeholder="192.168.1.0/24"
          />
        </div>

        <!-- VoIP VLAN -->
        <div>
          <label for="edit-voip-vlan" class="block text-sm font-medium text-gray-700 mb-1">
            VoIP VLAN
          </label>
          <input
            id="edit-voip-vlan"
            v-model.number="formData.voipVlan"
            type="number"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500"
            placeholder="100"
          />
        </div>

        <!-- VLAN -->
        <div>
          <label for="edit-vlan" class="block text-sm font-medium text-gray-700 mb-1">
            VLAN
          </label>
          <input
            id="edit-vlan"
            v-model.number="formData.vlan"
            type="number"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500"
            placeholder="10"
          />
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
            @click="updateLocation"
            :disabled="formLoading"
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
      title="Удалить локацию?"
      :message="`Вы уверены, что хотите удалить локацию &quot;${selectedLocation?.name}&quot;? Это действие нельзя отменить.`"
      confirm-text="Удалить"
      cancel-text="Отмена"
      variant="danger"
      :loading="formLoading"
      @confirm="deleteLocation"
      @cancel="showDeleteDialog = false"
    />
  </div>
</template>
