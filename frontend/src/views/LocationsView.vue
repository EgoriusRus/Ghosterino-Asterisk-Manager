<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { locationsAPI } from '@/api/client'
import type { Location } from '@/types/api'
import {
  mdiPlus,
  mdiPencil,
  mdiDelete,
  mdiMapMarker
} from '@mdi/js'

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
const formLoading = ref(false)
const formValid = ref(true)

// Selected location for edit/delete
const selectedLocation = ref<Location | null>(null)

// Table headers
const headers = [
  { title: 'Название', key: 'name', sortable: true },
  { title: 'Сервер', key: 'server', sortable: true },
  { title: 'Подсеть', key: 'subnet', sortable: false },
  { title: 'VoIP VLAN', key: 'voipVlan', sortable: true },
  { title: 'VLAN', key: 'vlan', sortable: true },
  { title: 'Создано', key: 'createdAt', sortable: true },
  { title: 'Действия', key: 'actions', sortable: false, align: 'end' as const }
]

// Validation rules
const rules = {
  required: (v: unknown) => !!v || 'Обязательное поле',
  ip: (v: string) => {
    if (!v) return 'IP адрес обязателен'
    const ipRegex = /^(\d{1,3}\.){3}\d{1,3}$/
    return ipRegex.test(v) || 'Некорректный формат IP адреса'
  }
}

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
  showEditModal.value = true
}

// Open delete dialog
const openDeleteDialog = (location: Location) => {
  selectedLocation.value = location
  showDeleteDialog.value = true
}

// Create location
const createLocation = async () => {
  if (!formValid.value) return

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
  if (!formValid.value) return

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
    month: 'short',
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
    <div class="d-flex justify-space-between align-center mb-6">
      <div>
        <h1 class="text-h4 font-weight-bold">Локации</h1>
        <p class="text-body-2 text-medium-emphasis">Управление локациями и их сетевыми настройками</p>
      </div>
      <v-btn
        color="primary"
        :prepend-icon="mdiPlus"
        @click="openCreateModal"
      >
        Добавить локацию
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
            Всего локаций: <strong>{{ locations.length }}</strong>
          </span>
        </div>
      </v-card-text>

      <v-data-table
        :headers="headers"
        :items="locations"
        :loading="loading"
        item-value="id"
      >
        <template #item.name="{ item }">
          <div class="d-flex align-center">
            <v-avatar color="success" size="32" class="mr-3" variant="tonal">
              <v-icon :icon="mdiMapMarker" size="small" />
            </v-avatar>
            <span class="font-weight-medium">{{ item.name }}</span>
          </div>
        </template>

        <template #item.server="{ item }">
          <code>{{ item.server }}</code>
        </template>

        <template #item.subnet="{ item }">
          <code v-if="item.subnet">{{ item.subnet }}</code>
          <span v-else class="text-medium-emphasis">—</span>
        </template>

        <template #item.voipVlan="{ item }">
          <v-chip color="purple" size="small" variant="tonal">
            {{ item.voipVlan }}
          </v-chip>
        </template>

        <template #item.vlan="{ item }">
          <v-chip color="info" size="small" variant="tonal">
            {{ item.vlan }}
          </v-chip>
        </template>

        <template #item.createdAt="{ item }">
          <span class="text-body-2">{{ formatDate(item.createdAt) }}</span>
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
            <p class="text-medium-emphasis">Локации не найдены</p>
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
        <v-card-title class="text-h6">Создать локацию</v-card-title>
        <v-card-text>
          <v-form v-model="formValid" @submit.prevent="createLocation">
            <v-text-field
              v-model="formData.name"
              label="Название"
              :rules="[rules.required]"
              placeholder="Офис Москва"
              class="mb-4"
            />

            <v-text-field
              v-model="formData.server"
              label="IP сервера"
              :rules="[rules.ip]"
              placeholder="192.168.1.1"
              class="mb-4"
            />

            <v-text-field
              v-model="formData.subnet"
              label="Подсеть"
              placeholder="192.168.1.0/24"
              class="mb-4"
            />

            <v-row dense>
              <v-col cols="6">
                <v-text-field
                  v-model.number="formData.voipVlan"
                  label="VoIP VLAN"
                  type="number"
                  placeholder="100"
                />
              </v-col>
              <v-col cols="6">
                <v-text-field
                  v-model.number="formData.vlan"
                  label="VLAN"
                  type="number"
                  placeholder="10"
                />
              </v-col>
            </v-row>
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
            @click="createLocation"
          >
            Создать
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Edit Dialog -->
    <v-dialog v-model="showEditModal" max-width="500" persistent>
      <v-card>
        <v-card-title class="text-h6">Редактировать локацию</v-card-title>
        <v-card-text>
          <v-form v-model="formValid" @submit.prevent="updateLocation">
            <v-text-field
              v-model="formData.name"
              label="Название"
              :rules="[rules.required]"
              placeholder="Офис Москва"
              class="mb-4"
            />

            <v-text-field
              v-model="formData.server"
              label="IP сервера"
              :rules="[rules.ip]"
              placeholder="192.168.1.1"
              class="mb-4"
            />

            <v-text-field
              v-model="formData.subnet"
              label="Подсеть"
              placeholder="192.168.1.0/24"
              class="mb-4"
            />

            <v-row dense>
              <v-col cols="6">
                <v-text-field
                  v-model.number="formData.voipVlan"
                  label="VoIP VLAN"
                  type="number"
                  placeholder="100"
                />
              </v-col>
              <v-col cols="6">
                <v-text-field
                  v-model.number="formData.vlan"
                  label="VLAN"
                  type="number"
                  placeholder="10"
                />
              </v-col>
            </v-row>
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
            @click="updateLocation"
          >
            Сохранить
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Delete Dialog -->
    <v-dialog v-model="showDeleteDialog" max-width="400">
      <v-card>
        <v-card-title class="text-h6">Удалить локацию?</v-card-title>
        <v-card-text>
          Вы уверены, что хотите удалить локацию "{{ selectedLocation?.name }}"?
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
            @click="deleteLocation"
          >
            Удалить
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>
