<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { devicesAPI } from '@/api/client'
import type { Device, DeviceModel } from '@/types/api'
import BaseModal from '@/components/BaseModal.vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'

// Device models
const DEVICE_MODELS: DeviceModel[] = ['Yealink T27G', 'Yealink T23G', 'Fanvil', 'Cisco']

// State
const devices = ref<Device[]>([])
const loading = ref(false)
const error = ref<string | null>(null)

// Modal states
const showCreateModal = ref(false)
const showEditModal = ref(false)
const showDeleteDialog = ref(false)

// Form data
const formData = ref({
  mac: '',
  deviceModel: '' as DeviceModel | ''
})

// Form states
const formErrors = ref<Record<string, string>>({})
const formLoading = ref(false)

// Selected device for edit/delete
const selectedDevice = ref<Device | null>(null)
const isEditMode = ref(false)

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

// Reset form
const resetForm = () => {
  formData.value = {
    mac: '',
    deviceModel: ''
  }
  formErrors.value = {}
  isEditMode.value = false
}

// Validate MAC address format
const validateMAC = (mac: string): boolean => {
  // Format: XX:XX:XX:XX:XX:XX (with hex digits)
  const macRegex = /^([0-9A-Fa-f]{2}:){5}[0-9A-Fa-f]{2}$/
  return macRegex.test(mac)
}

// Format MAC address input (auto-add colons)
const formatMACInput = (value: string): string => {
  // Remove all non-hex characters except colons
  let cleaned = value.replace(/[^0-9A-Fa-f]/g, '')

  // Add colons every 2 characters
  let formatted = ''
  for (let i = 0; i < cleaned.length && i < 12; i++) {
    if (i > 0 && i % 2 === 0) {
      formatted += ':'
    }
    formatted += cleaned[i]
  }

  return formatted.toUpperCase()
}

// Handle MAC input
const handleMACInput = (event: Event) => {
  const target = event.target as HTMLInputElement
  const cursorPosition = target.selectionStart || 0
  const oldValue = formData.value.mac
  const newValue = formatMACInput(target.value)

  formData.value.mac = newValue

  // Restore cursor position
  if (newValue.length >= oldValue.length) {
    setTimeout(() => {
      target.setSelectionRange(cursorPosition + (newValue.length - oldValue.length), cursorPosition + (newValue.length - oldValue.length))
    }, 0)
  }
}

// Validate form
const validateForm = (): boolean => {
  formErrors.value = {}

  if (!formData.value.mac.trim()) {
    formErrors.value.mac = 'MAC адрес обязателен'
  } else if (!validateMAC(formData.value.mac)) {
    formErrors.value.mac = 'Некорректный формат MAC адреса (XX:XX:XX:XX:XX:XX)'
  }

  if (!formData.value.deviceModel) {
    formErrors.value.deviceModel = 'Модель устройства обязательна'
  }

  return Object.keys(formErrors.value).length === 0
}

// Open create modal
const openCreateModal = () => {
  resetForm()
  showCreateModal.value = true
}

// Open edit modal
const openEditModal = (device: Device) => {
  selectedDevice.value = device
  isEditMode.value = true
  formData.value = {
    mac: device.mac,
    deviceModel: device.deviceModel
  }
  formErrors.value = {}
  showEditModal.value = true
}

// Open delete dialog
const openDeleteDialog = (device: Device) => {
  selectedDevice.value = device
  showDeleteDialog.value = true
}

// Create device
const createDevice = async () => {
  if (!validateForm()) {
    return
  }

  formLoading.value = true
  error.value = null

  try {
    await devicesAPI.create({
      mac: formData.value.mac,
      deviceModel: formData.value.deviceModel as DeviceModel
    })

    showCreateModal.value = false
    resetForm()
    await loadDevices()
  } catch (err) {
    error.value = 'Не удалось создать устройство'
    console.error('Failed to create device:', err)
  } finally {
    formLoading.value = false
  }
}

// Update device
const updateDevice = async () => {
  if (!validateForm() || !selectedDevice.value) {
    return
  }

  formLoading.value = true
  error.value = null

  try {
    await devicesAPI.update(selectedDevice.value.mac, {
      deviceModel: formData.value.deviceModel as DeviceModel
    })

    showEditModal.value = false
    resetForm()
    await loadDevices()
  } catch (err) {
    error.value = 'Не удалось обновить устройство'
    console.error('Failed to update device:', err)
  } finally {
    formLoading.value = false
  }
}

// Delete device
const deleteDevice = async () => {
  if (!selectedDevice.value) return

  formLoading.value = true
  error.value = null

  try {
    await devicesAPI.delete(selectedDevice.value.mac)
    showDeleteDialog.value = false
    selectedDevice.value = null
    await loadDevices()
  } catch (err) {
    error.value = 'Не удалось удалить устройство'
    console.error('Failed to delete device:', err)
  } finally {
    formLoading.value = false
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

// Format MAC address for display
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
      <div class="mt-4 sm:mt-0">
        <button
          @click="openCreateModal"
          class="inline-flex items-center px-4 py-2 bg-indigo-600 text-white text-sm font-medium rounded-lg hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-colors"
        >
          <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          Добавить устройство
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
              <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                Действия
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
              <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                <button
                  @click="openEditModal(device)"
                  class="text-indigo-600 hover:text-indigo-900 mr-3"
                  title="Редактировать"
                >
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                  </svg>
                </button>
                <button
                  @click="openDeleteDialog(device)"
                  class="text-red-600 hover:text-red-900"
                  title="Удалить"
                >
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  </svg>
                </button>
              </td>
            </tr>
            <tr v-if="devices.length === 0">
              <td colspan="5" class="px-6 py-8 text-center text-gray-500">
                Устройства не найдены
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Create Modal -->
    <BaseModal
      :show="showCreateModal"
      title="Создать устройство"
      @close="showCreateModal = false"
    >
      <form @submit.prevent="createDevice" class="space-y-4">
        <!-- MAC Address -->
        <div>
          <label for="create-mac" class="block text-sm font-medium text-gray-700 mb-1">
            MAC адрес <span class="text-red-500">*</span>
          </label>
          <input
            id="create-mac"
            :value="formData.mac"
            @input="handleMACInput"
            type="text"
            maxlength="17"
            :class="[
              'w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 font-mono',
              formErrors.mac ? 'border-red-300' : 'border-gray-300'
            ]"
            placeholder="00:11:22:33:44:55"
          />
          <p v-if="formErrors.mac" class="mt-1 text-sm text-red-600">{{ formErrors.mac }}</p>
          <p class="mt-1 text-xs text-gray-500">Формат: XX:XX:XX:XX:XX:XX</p>
        </div>

        <!-- Device Model -->
        <div>
          <label for="create-model" class="block text-sm font-medium text-gray-700 mb-1">
            Модель устройства <span class="text-red-500">*</span>
          </label>
          <select
            id="create-model"
            v-model="formData.deviceModel"
            :class="[
              'w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500',
              formErrors.deviceModel ? 'border-red-300' : 'border-gray-300'
            ]"
          >
            <option value="">Выберите модель</option>
            <option v-for="model in DEVICE_MODELS" :key="model" :value="model">
              {{ model }}
            </option>
          </select>
          <p v-if="formErrors.deviceModel" class="mt-1 text-sm text-red-600">{{ formErrors.deviceModel }}</p>
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
            @click="createDevice"
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
      title="Редактировать устройство"
      @close="showEditModal = false"
    >
      <form @submit.prevent="updateDevice" class="space-y-4">
        <!-- MAC Address (read-only in edit mode) -->
        <div>
          <label for="edit-mac" class="block text-sm font-medium text-gray-700 mb-1">
            MAC адрес
          </label>
          <input
            id="edit-mac"
            :value="formData.mac"
            type="text"
            readonly
            class="w-full px-3 py-2 border border-gray-300 rounded-lg bg-gray-50 text-gray-600 font-mono cursor-not-allowed"
          />
          <p class="mt-1 text-xs text-gray-500">MAC адрес нельзя изменить</p>
        </div>

        <!-- Device Model -->
        <div>
          <label for="edit-model" class="block text-sm font-medium text-gray-700 mb-1">
            Модель устройства <span class="text-red-500">*</span>
          </label>
          <select
            id="edit-model"
            v-model="formData.deviceModel"
            :class="[
              'w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500',
              formErrors.deviceModel ? 'border-red-300' : 'border-gray-300'
            ]"
          >
            <option value="">Выберите модель</option>
            <option v-for="model in DEVICE_MODELS" :key="model" :value="model">
              {{ model }}
            </option>
          </select>
          <p v-if="formErrors.deviceModel" class="mt-1 text-sm text-red-600">{{ formErrors.deviceModel }}</p>
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
            @click="updateDevice"
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
      title="Удалить устройство?"
      :message="`Вы уверены, что хотите удалить устройство с MAC адресом &quot;${selectedDevice?.mac}&quot;? Это действие нельзя отменить.`"
      confirm-text="Удалить"
      cancel-text="Отмена"
      variant="danger"
      :loading="formLoading"
      @confirm="deleteDevice"
      @cancel="showDeleteDialog = false"
    />
  </div>
</template>
