<script setup lang="ts">
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import {
  mdiAccountGroup,
  mdiCellphone,
  mdiMapMarker,
  mdiMenu
} from '@mdi/js'

const route = useRoute()
const drawer = ref(true)
const rail = ref(false)

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
      v-model="drawer"
      :rail="rail"
      permanent
    >
      <v-list-item
        title="Asterisk Manager"
        nav
      >
        <template #append>
          <v-btn
            :icon="mdiMenu"
            variant="text"
            size="small"
            @click="rail = !rail"
          />
        </template>
      </v-list-item>

      <v-divider />

      <v-list nav density="comfortable">
        <v-list-item
          v-for="item in menuItems"
          :key="item.id"
          :to="item.path"
          :prepend-icon="item.icon"
          :title="item.label"
          :active="isActiveRoute(item.path)"
          color="primary"
        />
      </v-list>
    </v-navigation-drawer>

    <v-main>
      <v-container fluid class="pa-6">
        <router-view />
      </v-container>
    </v-main>
  </v-app>
</template>
