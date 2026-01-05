<script setup lang="ts">
import { useRoute } from 'vue-router'
import {
  mdiAccountGroup,
  mdiCellphone,
  mdiMapMarker,
  mdiPhoneVoip
} from '@mdi/js'

const route = useRoute()

const menuItems = [
  {
    id: 'profiles',
    label: 'Сотрудники',
    icon: mdiAccountGroup,
    path: '/admin/profiles'
  },
  {
    id: 'devices',
    label: 'Устройства',
    icon: mdiCellphone,
    path: '/admin/devices'
  },
  {
    id: 'locations',
    label: 'Локации',
    icon: mdiMapMarker,
    path: '/admin/locations'
  }
]

const isActiveRoute = (path: string): boolean => {
  return route.path === path
}
</script>

<template>
  <v-app>
    <v-navigation-drawer
      permanent
      color="primary"
    >
      <v-list-item class="py-4">
        <template #prepend>
          <v-icon :icon="mdiPhoneVoip" size="28" />
        </template>
        <v-list-item-title class="text-h6 font-weight-bold">
          Asterisk
        </v-list-item-title>
        <v-list-item-subtitle class="text-caption">
          Manager
        </v-list-item-subtitle>
      </v-list-item>

      <v-divider class="my-2" />

      <v-list nav density="comfortable" class="px-2">
        <v-list-item
          v-for="item in menuItems"
          :key="item.id"
          :to="item.path"
          :prepend-icon="item.icon"
          :title="item.label"
          :active="isActiveRoute(item.path)"
          rounded="lg"
          class="mb-1"
        />
      </v-list>
    </v-navigation-drawer>

    <v-main class="bg-background">
      <v-container fluid class="pa-6">
        <router-view />
      </v-container>
    </v-main>
  </v-app>
</template>

<style scoped>
.v-navigation-drawer :deep(.v-list-item--active) {
  background: rgba(255, 255, 255, 0.15) !important;
}
.v-navigation-drawer :deep(.v-list-item:hover:not(.v-list-item--active)) {
  background: rgba(255, 255, 255, 0.08) !important;
}
</style>
