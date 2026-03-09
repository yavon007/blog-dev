<script setup lang="ts">
import type { MediaFile } from '@/types'

defineProps<{
  images: MediaFile[]
  selected: number | null
  loading?: boolean
}>()

const emit = defineEmits<{
  select: [image: MediaFile]
}>()

function formatSize(bytes: number) {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
}
</script>

<template>
  <div class="grid grid-cols-3 sm:grid-cols-4 md:grid-cols-5 gap-3">
    <div
      v-for="img in images"
      :key="img.id"
      class="relative aspect-square rounded-lg overflow-hidden cursor-pointer border-2 transition-all"
      :class="selected === img.id
        ? 'border-primary-500 ring-2 ring-primary-500/30'
        : 'border-gray-200 dark:border-gray-700 hover:border-primary-300 dark:hover:border-primary-600'"
      @click="emit('select', img)"
    >
      <img
        :src="img.url"
        :alt="img.alt_text || img.original_name"
        class="w-full h-full object-cover"
      />
      <div class="absolute inset-x-0 bottom-0 bg-black/50 px-2 py-1">
        <p class="text-xs text-white truncate">{{ img.original_name }}</p>
        <p class="text-xs text-gray-300">{{ formatSize(img.size) }}</p>
      </div>
    </div>

    <div v-if="loading" class="contents">
      <div v-for="i in 8" :key="i" class="aspect-square bg-gray-200 dark:bg-gray-700 rounded-lg animate-pulse" />
    </div>

    <div v-if="!loading && images.length === 0" class="col-span-full text-center py-8 text-gray-500 dark:text-gray-400">
      暂无图片
    </div>
  </div>
</template>
