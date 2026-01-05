<script setup lang="ts">
defineProps<{
  msg: string
}>()
</script>

<template>
  <div class="bg-white rounded-2xl shadow-2xl p-8 max-w-2xl w-full">
    <h1 class="text-4xl font-bold text-center text-indigo-600 mb-6">
      {{ msg }}
    </h1>

    <div class="space-y-4">
      <p class="text-lg text-gray-700 text-center">
        Hello World!
      </p>

      <div class="flex justify-center gap-4 mt-8">
        <span class="inline-block px-4 py-2 bg-blue-100 text-blue-800 rounded-full text-sm font-semibold">
          Vue 3
        </span>
        <span class="inline-block px-4 py-2 bg-green-100 text-green-800 rounded-full text-sm font-semibold">
          TypeScript
        </span>
        <span class="inline-block px-4 py-2 bg-purple-100 text-purple-800 rounded-full text-sm font-semibold">
          Tailwind CSS
        </span>
      </div>

      <div class="mt-8 text-center text-gray-600">
        <p class="text-sm">
          Edit <code class="bg-gray-100 px-2 py-1 rounded text-sm font-mono">src/components/HelloWorld.vue</code> to get started
        </p>
      </div>
    </div>
  </div>
</template>
