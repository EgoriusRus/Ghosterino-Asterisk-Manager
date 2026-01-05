<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { devicesAPI } from '@/api/client'
import type { Device, DeviceModel } from '@/types/api'
import {
  mdiPlus,
  mdiPencil,
  mdiDelete,
  mdiCellphone
} from '@mdi/js'

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
  deviceModel: null as DeviceModel | null
})

// Form states
const formLoading = ref(false)
const formValid = ref(true)

// Selected device for edit/delete
const selectedDevice = ref<Device | null>(null)

// Table headers
const headers = [
  { title: 'MAC адрес', key: 'mac', sortable: true },
  { title: 'Модель', key: 'deviceModel', sortable: true },
  { title: 'Создано', key: 'createdAt', sortable: true },
  { title: 'Обновлено', key: 'updatedAt', sortable: true },
  { title: 'Действия', key: 'actions', sortable: false, align: 'end' as const }
]

// Validation rules
const rules = {
  required: (v: unknown) => !!v || 'Обязательное поле',
  mac: (v: string) => {
    if (!v) return 'MAC адрес обязателен'
    const macRegex = /^([0-9A-Fa-f]{2}:){5}[0-9A-Fa-f]{2}$/
    return macRegex.test(v) || 'Формат: XX:XX:XX:XX:XX:XX'
  }
}

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
    deviceModel: null
  }
}

// Format MAC address input (auto-add colons)
const formatMACInput = (value: string): string => {
  let cleaned = value.replace(/[^0-9A-Fa-f]/g, '')
  let formatted = ''
  for (let i = 0; i < cleaned.length && i < 12; i++) {
    if (i > 0 && i % 2 === 0) {
      formatted += ':'
    }
    formatted += cleaned[i]
  }
  return formatted.toUpperCase()
}

// Open create modal
const openCreateModal = () => {
  resetForm()
  showCreateModal.value = true
}

// Open edit modal
const openEditModal = (device: Device) => {
  selectedDevice.value = device
  formData.value = {
    mac: device.mac,
    deviceModel: device.deviceModel
  }
  showEditModal.value = true
}

// Open delete dialog
const openDeleteDialog = (device: Device) => {
  selectedDevice.value = device
  showDeleteDialog.value = true
}

// Create device
const createDevice = async () => {
  if (!formValid.value) return

  formLoading.value = true
  error.value = null

  try {
    await devicesAPI.create({
      mac: formData.value.mac,
      deviceModel: formData.value.deviceModel!
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
  if (!formValid.value || !selectedDevice.value) return

  formLoading.value = true
  error.value = null

  try {
    await devicesAPI.update(selectedDevice.value.mac, {
      deviceModel: formData.value.deviceModel!
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
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// Load on mount
onMounted(() => {
  loadDevices()
})
</script>

<template>
  <div>
    <!-- Page header -->
    <div class="d-flex justify-space-between align-center mb-6">
      <div>
        <h1 class="text-h4 font-weight-bold">Устройства</h1>
        <p class="text-body-2 text-medium-emphasis">Управление IP-телефонами и сетевыми устройствами</p>
      </div>
      <v-btn
        color="primary"
        :prepend-icon="mdiPlus"
        @click="openCreateModal"
      >
        Добавить устройство
      </v-btn>
    </div>

    <!-- Error message -->
    <v-alert
      v-if="error"
      type="error"
      variant="tonal"
      closable
      class="mb-6"
      @click:close="error = null"
    >
      {{ error }}
    </v-alert>

    <!-- Data table -->
    <v-card>
      <v-card-text class="pa-0">
        <div class="px-4 py-3 bg-grey-lighten-4">
          <span class="text-body-2 text-medium-emphasis">
            Всего устройств: <strong>{{ devices.length }}</strong>
          </span>
        </div>
      </v-card-text>

      <v-data-table
        :headers="headers"
        :items="devices"
        :loading="loading"
        item-value="mac"
      >
        <template #item.mac="{ item }">
          <div class="d-flex align-center">
            <v-avatar color="primary" size="32" class="mr-3" variant="tonal">
              <v-icon :icon="mdiCellphone" size="small" />
            </v-avatar>
            <code class="font-weight-medium">{{ item.mac.toUpperCase() }}</code>
          </div>
        </template>

        <template #item.deviceModel="{ item }">
          <v-chip color="info" size="small" variant="tonal">
            {{ item.deviceModel }}
          </v-chip>
        </template>

        <template #item.createdAt="{ item }">
          <span class="text-body-2">{{ formatDate(item.createdAt) }}</span>
        </template>

        <template #item.updatedAt="{ item }">
          <span class="text-body-2">{{ formatDate(item.updatedAt) }}</span>
        </template>

        <template #item.actions="{ item }">
          <v-btn
            icon
            variant="text"
            size="small"
            color="primary"
            @click="openEditModal(item)"
          >
            <v-icon :icon="mdiPencil" size="small" />
            <v-tooltip activator="parent" location="top">Редактировать</v-tooltip>
          </v-btn>
          <v-btn
            icon
            variant="text"
            size="small"
            color="error"
            @click="openDeleteDialog(item)"
          >
            <v-icon :icon="mdiDelete" size="small" />
            <v-tooltip activator="parent" location="top">Удалить</v-tooltip>
          </v-btn>
        </template>

        <template #no-data>
          <div class="text-center pa-4">
            <p class="text-medium-emphasis">Устройства не найдены</p>
          </div>
        </template>

        <template #loading>
          <v-skeleton-loader type="table-row@5" />
        </template>
      </v-data-table>
    </v-card>

    <!-- Create Dialog -->
    <v-dialog v-model="showCreateModal" max-width="500" persistent>
      <v-card>
        <v-card-title class="text-h6">Создать устройство</v-card-title>
        <v-card-text>
          <v-form v-model="formValid" @submit.prevent="createDevice">
            <v-text-field
              :model-value="formData.mac"
              @update:model-value="(v: string) => formData.mac = formatMACInput(v)"
              label="MAC адрес"
              :rules="[rules.mac]"
              placeholder="00:11:22:33:44:55"
              maxlength="17"
              class="mb-4"
              hint="Формат: XX:XX:XX:XX:XX:XX"
              persistent-hint
            />

            <v-select
              v-model="formData.deviceModel"
              label="Модель устройства"
              :items="DEVICE_MODELS"
              :rules="[rules.required]"
            />
          </v-form>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn
            variant="text"
            :disabled="formLoading"
            @click="showCreateModal = false"
          >
            Отмена
          </v-btn>
          <v-btn
            color="primary"
            :loading="formLoading"
            :disabled="!formValid"
            @click="createDevice"
          >
            Создать
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Edit Dialog -->
    <v-dialog v-model="showEditModal" max-width="500" persistent>
      <v-card>
        <v-card-title class="text-h6">Редактировать устройство</v-card-title>
        <v-card-text>
          <v-form v-model="formValid" @submit.prevent="updateDevice">
            <v-text-field
              v-model="formData.mac"
              label="MAC адрес"
              readonly
              disabled
              class="mb-4"
              hint="MAC адрес нельзя изменить"
              persistent-hint
            />

            <v-select
              v-model="formData.deviceModel"
              label="Модель устройства"
              :items="DEVICE_MODELS"
              :rules="[rules.required]"
            />
          </v-form>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn
            variant="text"
            :disabled="formLoading"
            @click="showEditModal = false"
          >
            Отмена
          </v-btn>
          <v-btn
            color="primary"
            :loading="formLoading"
            :disabled="!formValid"
            @click="updateDevice"
          >
            Сохранить
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Delete Dialog -->
    <v-dialog v-model="showDeleteDialog" max-width="400">
      <v-card>
        <v-card-title class="text-h6">Удалить устройство?</v-card-title>
        <v-card-text>
          Вы уверены, что хотите удалить устройство с MAC адресом "{{ selectedDevice?.mac }}"?
          Это действие нельзя отменить.
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn
            variant="text"
            :disabled="formLoading"
            @click="showDeleteDialog = false"
          >
            Отмена
          </v-btn>
          <v-btn
            color="error"
            :loading="formLoading"
            @click="deleteDevice"
          >
            Удалить
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>
