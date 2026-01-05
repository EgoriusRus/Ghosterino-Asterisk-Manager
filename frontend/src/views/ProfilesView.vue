<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { profilesAPI, devicesAPI, locationsAPI } from '@/api/client'
import type { ProfileWithLocation, PaginationResponse, Device, Location } from '@/types/api'
import {
  mdiPlus,
  mdiPencil,
  mdiDelete,
  mdiAccount
} from '@mdi/js'

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
const formLoading = ref(false)
const formValid = ref(true)

// Selected profile for edit/delete
const selectedProfile = ref<ProfileWithLocation | null>(null)

// Table headers
const headers = [
  { title: 'Имя', key: 'name', sortable: true },
  { title: 'Внутренний', key: 'internalNumber', sortable: true },
  { title: 'Локация', key: 'locationName', sortable: true },
  { title: 'Устройство', key: 'device', sortable: false },
  { title: 'Статус', key: 'isActive', sortable: true },
  { title: 'Действия', key: 'actions', sortable: false, align: 'end' as const }
]

// Validation rules
const rules = {
  required: (v: unknown) => !!v || 'Обязательное поле',
  positiveNumber: (v: number | null) => !v || v > 0 || 'Должно быть положительным числом'
}

// Location items for select
const locationItems = computed(() =>
  locations.value.map(l => ({ title: l.name, value: l.id }))
)

// Device items for select
const deviceItems = computed(() =>
  devices.value.map(d => ({ title: `${d.mac} (${d.deviceModel})`, value: d.mac }))
)

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
  if (!formValid.value) return

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
  if (!formValid.value) return

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

// Handle page change
const onPageChange = (page: number) => {
  pagination.value.page = page
  loadProfiles()
}

// Load on mount
onMounted(() => {
  loadProfiles()
})
</script>

<template>
  <div>
    <!-- Page header -->
    <div class="d-flex justify-space-between align-center mb-6">
      <div>
        <h1 class="text-h4 font-weight-bold">Сотрудники</h1>
        <p class="text-body-2 text-medium-emphasis">Управление профилями сотрудников и их SIP-настройками</p>
      </div>
      <v-btn
        color="primary"
        :prepend-icon="mdiPlus"
        @click="openCreateModal"
      >
        Добавить сотрудника
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
            Всего сотрудников: <strong>{{ pagination.total }}</strong>
          </span>
        </div>
      </v-card-text>

      <v-data-table
        :headers="headers"
        :items="profiles"
        :loading="loading"
        :items-per-page="pagination.perPage"
        hide-default-footer
      >
        <template #item.name="{ item }">
          <div class="d-flex align-center">
            <v-avatar color="primary" size="32" class="mr-3">
              <v-icon :icon="mdiAccount" size="small" />
            </v-avatar>
            <span class="font-weight-medium">{{ item.name }}</span>
          </div>
        </template>

        <template #item.locationName="{ item }">
          {{ item.locationName || '—' }}
        </template>

        <template #item.device="{ item }">
          <code v-if="item.device">{{ item.device }}</code>
          <span v-else class="text-medium-emphasis">—</span>
        </template>

        <template #item.isActive="{ item }">
          <v-chip
            :color="item.isActive ? 'success' : 'default'"
            size="small"
            variant="tonal"
          >
            {{ item.isActive ? 'Активен' : 'Неактивен' }}
          </v-chip>
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
            <p class="text-medium-emphasis">Сотрудники не найдены</p>
          </div>
        </template>

        <template #loading>
          <v-skeleton-loader type="table-row@5" />
        </template>
      </v-data-table>

      <!-- Pagination -->
      <template v-if="pagination.pages > 1">
        <v-divider />
        <div class="pa-4 d-flex justify-space-between align-center">
          <span class="text-body-2 text-medium-emphasis">
            Показано {{ (pagination.page - 1) * pagination.perPage + 1 }}–{{ Math.min(pagination.page * pagination.perPage, pagination.total) }} из {{ pagination.total }}
          </span>
          <v-pagination
            :model-value="pagination.page"
            :length="pagination.pages"
            :total-visible="5"
            density="comfortable"
            @update:model-value="onPageChange"
          />
        </div>
      </template>
    </v-card>

    <!-- Create Dialog -->
    <v-dialog v-model="showCreateModal" max-width="600" persistent>
      <v-card>
        <v-card-title class="text-h6">Создать сотрудника</v-card-title>
        <v-card-text>
          <v-progress-linear v-if="loadingReferences" indeterminate class="mb-4" />
          <v-form v-else v-model="formValid" @submit.prevent="createProfile">
            <v-row dense>
              <v-col cols="12">
                <v-text-field
                  v-model="formData.name"
                  label="ФИО"
                  :rules="[rules.required]"
                  placeholder="Иванов Иван Иванович"
                />
              </v-col>

              <v-col cols="12" md="6">
                <v-text-field
                  v-model.number="formData.internalNumber"
                  label="Внутренний номер"
                  type="number"
                  :rules="[rules.required, rules.positiveNumber]"
                  placeholder="1001"
                />
              </v-col>

              <v-col cols="12" md="6">
                <v-text-field
                  v-model="formData.externalNumber"
                  label="Внешний номер"
                  placeholder="+7 (123) 456-78-90"
                />
              </v-col>

              <v-col cols="12" md="6">
                <v-select
                  v-model="formData.locationId"
                  label="Локация"
                  :items="locationItems"
                  :rules="[rules.required]"
                  clearable
                />
              </v-col>

              <v-col cols="12" md="6">
                <v-select
                  v-model="formData.device"
                  label="Устройство (MAC)"
                  :items="deviceItems"
                  clearable
                />
              </v-col>

              <v-col cols="12" md="6">
                <v-text-field
                  v-model.number="formData.ringGroup"
                  label="Ring Group"
                  type="number"
                  :rules="[rules.positiveNumber]"
                  placeholder="1"
                />
              </v-col>

              <v-col cols="12" md="6">
                <v-text-field
                  v-model.number="formData.pickupGroup"
                  label="Pickup Group"
                  type="number"
                  :rules="[rules.positiveNumber]"
                  placeholder="1"
                />
              </v-col>

              <v-col cols="12">
                <v-checkbox
                  v-model="formData.isActive"
                  label="Активен"
                  hide-details
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
            :disabled="!formValid || loadingReferences"
            @click="createProfile"
          >
            Создать
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Edit Dialog -->
    <v-dialog v-model="showEditModal" max-width="600" persistent>
      <v-card>
        <v-card-title class="text-h6">Редактировать сотрудника</v-card-title>
        <v-card-text>
          <v-progress-linear v-if="loadingReferences" indeterminate class="mb-4" />
          <v-form v-else v-model="formValid" @submit.prevent="updateProfile">
            <v-row dense>
              <v-col cols="12">
                <v-text-field
                  v-model="formData.name"
                  label="ФИО"
                  :rules="[rules.required]"
                  placeholder="Иванов Иван Иванович"
                />
              </v-col>

              <v-col cols="12" md="6">
                <v-text-field
                  v-model.number="formData.internalNumber"
                  label="Внутренний номер"
                  type="number"
                  :rules="[rules.required, rules.positiveNumber]"
                  placeholder="1001"
                />
              </v-col>

              <v-col cols="12" md="6">
                <v-text-field
                  v-model="formData.externalNumber"
                  label="Внешний номер"
                  placeholder="+7 (123) 456-78-90"
                />
              </v-col>

              <v-col cols="12" md="6">
                <v-select
                  v-model="formData.locationId"
                  label="Локация"
                  :items="locationItems"
                  :rules="[rules.required]"
                  clearable
                />
              </v-col>

              <v-col cols="12" md="6">
                <v-select
                  v-model="formData.device"
                  label="Устройство (MAC)"
                  :items="deviceItems"
                  clearable
                />
              </v-col>

              <v-col cols="12" md="6">
                <v-text-field
                  v-model.number="formData.ringGroup"
                  label="Ring Group"
                  type="number"
                  :rules="[rules.positiveNumber]"
                  placeholder="1"
                />
              </v-col>

              <v-col cols="12" md="6">
                <v-text-field
                  v-model.number="formData.pickupGroup"
                  label="Pickup Group"
                  type="number"
                  :rules="[rules.positiveNumber]"
                  placeholder="1"
                />
              </v-col>

              <v-col cols="12">
                <v-checkbox
                  v-model="formData.isActive"
                  label="Активен"
                  hide-details
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
            :disabled="!formValid || loadingReferences"
            @click="updateProfile"
          >
            Сохранить
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Delete Dialog -->
    <v-dialog v-model="showDeleteDialog" max-width="400">
      <v-card>
        <v-card-title class="text-h6">Удалить сотрудника?</v-card-title>
        <v-card-text>
          Вы уверены, что хотите удалить сотрудника "{{ selectedProfile?.name }}"?
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
            @click="deleteProfile"
          >
            Удалить
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>
