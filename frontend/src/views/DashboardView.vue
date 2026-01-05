<script setup lang="ts">
import { ref } from 'vue'

interface StatCard {
  title: string
  value: string
  change: string
  changeType: 'increase' | 'decrease'
  icon: string
}

const stats = ref<StatCard[]>([
  {
    title: 'Total Users',
    value: '2,543',
    change: '+12.5%',
    changeType: 'increase',
    icon: 'M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z'
  },
  {
    title: 'Active Sessions',
    value: '1,234',
    change: '+5.2%',
    changeType: 'increase',
    icon: 'M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z'
  },
  {
    title: 'Revenue',
    value: '$45,231',
    change: '-2.3%',
    changeType: 'decrease',
    icon: 'M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z'
  },
  {
    title: 'New Orders',
    value: '324',
    change: '+18.7%',
    changeType: 'increase',
    icon: 'M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z'
  }
])
</script>

<template>
  <div>
    <!-- Page header -->
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-gray-800">Dashboard</h1>
      <p class="text-gray-600 mt-1">Welcome back! Here's what's happening today.</p>
    </div>

    <!-- Stats grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
      <div
        v-for="stat in stats"
        :key="stat.title"
        class="bg-white rounded-lg shadow-sm border border-gray-200 p-6 hover:shadow-md transition-shadow"
      >
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-600">{{ stat.title }}</p>
            <p class="text-2xl font-bold text-gray-800 mt-2">{{ stat.value }}</p>
            <div class="flex items-center mt-2">
              <span
                :class="[
                  'text-sm font-medium',
                  stat.changeType === 'increase' ? 'text-green-600' : 'text-red-600'
                ]"
              >
                {{ stat.change }}
              </span>
              <span class="text-sm text-gray-500 ml-2">vs last month</span>
            </div>
          </div>
          <div class="w-12 h-12 bg-indigo-100 rounded-lg flex items-center justify-center">
            <svg
              class="w-6 h-6 text-indigo-600"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                :d="stat.icon"
              />
            </svg>
          </div>
        </div>
      </div>
    </div>

    <!-- Recent activity placeholder -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Recent users -->
      <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
        <h2 class="text-lg font-semibold text-gray-800 mb-4">Recent Users</h2>
        <div class="space-y-4">
          <div
            v-for="i in 5"
            :key="i"
            class="flex items-center justify-between py-3 border-b border-gray-100 last:border-0"
          >
            <div class="flex items-center space-x-3">
              <div class="w-10 h-10 bg-indigo-600 rounded-full flex items-center justify-center text-white font-semibold">
                U{{ i }}
              </div>
              <div>
                <p class="text-sm font-medium text-gray-800">User {{ i }}</p>
                <p class="text-xs text-gray-500">user{{ i }}@example.com</p>
              </div>
            </div>
            <span class="text-xs text-gray-500">{{ i }}h ago</span>
          </div>
        </div>
      </div>

      <!-- Activity log -->
      <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
        <h2 class="text-lg font-semibold text-gray-800 mb-4">Activity Log</h2>
        <div class="space-y-4">
          <div
            v-for="i in 5"
            :key="i"
            class="flex items-start space-x-3 py-3 border-b border-gray-100 last:border-0"
          >
            <div class="w-2 h-2 bg-indigo-600 rounded-full mt-2"></div>
            <div class="flex-1">
              <p class="text-sm text-gray-800">Activity event {{ i }}</p>
              <p class="text-xs text-gray-500 mt-1">{{ i * 10 }} minutes ago</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
