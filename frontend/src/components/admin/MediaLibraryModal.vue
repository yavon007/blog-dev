<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { mediaApi } from '@/api/media'
import ImageDropzone from './ImageDropzone.vue'
import ImageGrid from './ImageGrid.vue'
import type { MediaFile, UploadResult } from '@/types'

const emit = defineEmits<{
  close: []
  select: [url: string, alt: string]
}>()

const images = ref<MediaFile[]>([])
const selected = ref<number | null>(null)
const loading = ref(false)
const page = ref(1)
const total = ref(0)

async function fetchImages() {
  loading.value = true
  try {
    const res = await mediaApi.list(page.value, 20) as unknown as { data: { items: MediaFile[]; total: number } }
    images.value = res.data.items ?? []
    total.value = res.data.total ?? 0
  } finally {
    loading.value = false
  }
}

function handleUpload(_result: UploadResult) {
  fetchImages()
}

function handleSelect(img: MediaFile) {
  selected.value = selected.value === img.id ? null : img.id
}

function confirm() {
  const img = images.value.find(i => i.id === selected.value)
  if (img) {
    emit('select', img.url, img.alt_text || img.original_name)
  }
}

onMounted(fetchImages)
</script>

<template>
  <div class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" @click.self="emit('close')">
    <div class="bg-white dark:bg-gray-800 rounded-lg w-full max-w-4xl max-h-[80vh] overflow-hidden flex flex-col">
      <!-- Header -->
      <div class="p-4 border-b border-gray-200 dark:border-gray-700 flex justify-between items-center">
        <h3 class="font-bold text-lg">媒体库</h3>
        <button
          class="p-1 hover:bg-gray-100 dark:hover:bg-gray-700 rounded"
          @click="emit('close')"
        >
          <span class="i-carbon-close text-xl" />
        </button>
      </div>

      <!-- Content -->
      <div class="flex-1 overflow-auto p-4 space-y-4">
        <ImageDropzone @upload="handleUpload" />
        <ImageGrid
          :images="images"
          :selected="selected"
          :loading="loading"
          @select="handleSelect"
        />
      </div>

      <!-- Footer -->
      <div class="p-4 border-t border-gray-200 dark:border-gray-700 flex justify-between items-center">
        <span class="text-sm text-gray-500 dark:text-gray-400">
          共 {{ total }} 张图片
        </span>
        <div class="flex gap-2">
          <button class="btn-secondary" @click="emit('close')">取消</button>
          <button
            class="btn-primary"
            :disabled="!selected"
            @click="confirm"
          >
            确认选择
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
