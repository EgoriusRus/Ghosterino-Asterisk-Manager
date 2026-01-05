<script setup lang="ts">
import { ref } from 'vue'

interface FormData {
  siteName: string
  siteEmail: string
  notifications: {
    email: boolean
    push: boolean
    sms: boolean
  }
  maintenance: boolean
  timezone: string
}

const formData = ref<FormData>({
  siteName: 'Admin Panel',
  siteEmail: 'admin@example.com',
  notifications: {
    email: true,
    push: false,
    sms: false
  },
  maintenance: false,
  timezone: 'UTC'
})

const handleSave = () => {
  console.log('Saving settings:', formData.value)
  alert('Settings saved successfully!')
}
</script>

<template>
  <div>
    <!-- Page header -->
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-gray-800">Settings</h1>
      <p class="text-gray-600 mt-1">Manage your application settings and preferences</p>
    </div>

    <!-- Settings form -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Sidebar navigation -->
      <div class="lg:col-span-1">
        <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-4">
          <nav class="space-y-1">
            <button
              class="w-full text-left px-3 py-2 rounded-lg bg-indigo-50 text-indigo-600 font-medium"
            >
              General
            </button>
            <button
              class="w-full text-left px-3 py-2 rounded-lg text-gray-700 hover:bg-gray-50"
            >
              Notifications
            </button>
            <button
              class="w-full text-left px-3 py-2 rounded-lg text-gray-700 hover:bg-gray-50"
            >
              Security
            </button>
            <button
              class="w-full text-left px-3 py-2 rounded-lg text-gray-700 hover:bg-gray-50"
            >
              API Keys
            </button>
            <button
              class="w-full text-left px-3 py-2 rounded-lg text-gray-700 hover:bg-gray-50"
            >
              Advanced
            </button>
          </nav>
        </div>
      </div>

      <!-- Settings content -->
      <div class="lg:col-span-2 space-y-6">
        <!-- General settings -->
        <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <h2 class="text-lg font-semibold text-gray-800 mb-4">General Settings</h2>

          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Site Name
              </label>
              <input
                v-model="formData.siteName"
                type="text"
                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent outline-none"
              />
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Site Email
              </label>
              <input
                v-model="formData.siteEmail"
                type="email"
                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent outline-none"
              />
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Timezone
              </label>
              <select
                v-model="formData.timezone"
                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent outline-none"
              >
                <option value="UTC">UTC</option>
                <option value="America/New_York">Eastern Time</option>
                <option value="America/Chicago">Central Time</option>
                <option value="America/Los_Angeles">Pacific Time</option>
                <option value="Europe/London">London</option>
              </select>
            </div>
          </div>
        </div>

        <!-- Notification settings -->
        <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <h2 class="text-lg font-semibold text-gray-800 mb-4">Notification Preferences</h2>

          <div class="space-y-4">
            <div class="flex items-center justify-between">
              <div>
                <p class="text-sm font-medium text-gray-800">Email Notifications</p>
                <p class="text-sm text-gray-500">Receive email updates about activity</p>
              </div>
              <label class="relative inline-flex items-center cursor-pointer">
                <input
                  v-model="formData.notifications.email"
                  type="checkbox"
                  class="sr-only peer"
                />
                <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-indigo-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-indigo-600"></div>
              </label>
            </div>

            <div class="flex items-center justify-between">
              <div>
                <p class="text-sm font-medium text-gray-800">Push Notifications</p>
                <p class="text-sm text-gray-500">Receive push notifications in browser</p>
              </div>
              <label class="relative inline-flex items-center cursor-pointer">
                <input
                  v-model="formData.notifications.push"
                  type="checkbox"
                  class="sr-only peer"
                />
                <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-indigo-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-indigo-600"></div>
              </label>
            </div>

            <div class="flex items-center justify-between">
              <div>
                <p class="text-sm font-medium text-gray-800">SMS Notifications</p>
                <p class="text-sm text-gray-500">Receive SMS alerts for critical events</p>
              </div>
              <label class="relative inline-flex items-center cursor-pointer">
                <input
                  v-model="formData.notifications.sms"
                  type="checkbox"
                  class="sr-only peer"
                />
                <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-indigo-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-indigo-600"></div>
              </label>
            </div>
          </div>
        </div>

        <!-- Maintenance mode -->
        <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <h2 class="text-lg font-semibold text-gray-800 mb-4">Maintenance Mode</h2>

          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm font-medium text-gray-800">Enable Maintenance Mode</p>
              <p class="text-sm text-gray-500">Temporarily disable public access to the site</p>
            </div>
            <label class="relative inline-flex items-center cursor-pointer">
              <input
                v-model="formData.maintenance"
                type="checkbox"
                class="sr-only peer"
              />
              <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-red-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-red-600"></div>
            </label>
          </div>
        </div>

        <!-- Save button -->
        <div class="flex justify-end">
          <button
            @click="handleSave"
            class="px-6 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors font-medium"
          >
            Save Changes
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
