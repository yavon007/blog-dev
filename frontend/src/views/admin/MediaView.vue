<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { mediaApi } from '@/api/media'
import ImageDropzone from '@/components/admin/ImageDropzone.vue'
import ImageGrid from '@/components/admin/ImageGrid.vue'
import type { MediaFile } from '@/types'

const images = ref<MediaFile[]>([])
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

onMounted(fetchImages)
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold">媒体库</h1>
      <span class="text-sm text-gray-500 dark:text-gray-400">共 {{ total }} 张图片</span>
    </div>

    <div class="card p-6 mb-6">
      <h2 class="font-semibold mb-4">上传图片</h2>
      <ImageDropzone @upload="fetchImages" />
    </div>

    <div class="card p-6">
      <h2 class="font-semibold mb-4">图片列表</h2>
      <ImageGrid :images="images" :loading="loading" :selected="null" />

      <!-- 分页 -->
      <div v-if="total > 20" class="flex justify-center gap-2 mt-6">
        <button
          :disabled="page <= 1"
          class="btn-secondary disabled:opacity-40"
          @click="page--; fetchImages()"
        >
          上一页
        </button>
        <span class="px-4 py-2 text-sm text-gray-600 dark:text-gray-400">
          第 {{ page }} 页 / 共 {{ Math.ceil(total / 20) }} 页
        </span>
        <button
          :disabled="page >= Math.ceil(total / 20)"
          class="btn-secondary disabled:opacity-40"
          @click="page++; fetchImages()"
        >
          下一页
        </button>
      </div>
    </div>
  </div>
</template>
