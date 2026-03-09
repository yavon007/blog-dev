<script setup lang="ts">
import { ref } from 'vue'
import { mediaApi } from '@/api/media'
import type { UploadResult } from '@/types'

const emit = defineEmits<{
  upload: [result: UploadResult]
}>()

const isDragging = ref(false)
const uploading = ref(false)
const error = ref('')

async function handleDrop(e: DragEvent) {
  isDragging.value = false
  const files = e.dataTransfer?.files
  if (files?.length && files[0]) {
    await uploadFile(files[0])
  }
}

async function handleChange(e: Event) {
  const target = e.target as HTMLInputElement
  const files = target.files
  if (files?.length && files[0]) {
    await uploadFile(files[0])
  }
  target.value = ''
}

async function uploadFile(file: File) {
  if (!file.type.startsWith('image/')) {
    error.value = '请上传图片文件'
    return
  }
  if (file.size > 5 * 1024 * 1024) {
    error.value = '文件大小不能超过 5MB'
    return
  }

  uploading.value = true
  error.value = ''
  try {
    const res = await mediaApi.upload(file) as unknown as { data: UploadResult }
    emit('upload', res.data)
  } catch (e) {
    error.value = e instanceof Error ? e.message : '上传失败'
  } finally {
    uploading.value = false
  }
}
</script>

<template>
  <div>
    <div
      class="border-2 border-dashed border-gray-300 dark:border-gray-600 rounded-lg p-8 text-center cursor-pointer transition-colors"
      :class="{
        'border-primary-500 bg-primary-50 dark:bg-primary-900/20': isDragging,
        'hover:border-primary-500 hover:bg-gray-50 dark:hover:bg-gray-800': !isDragging,
      }"
      @dragover.prevent="isDragging = true"
      @dragleave="isDragging = false"
      @drop.prevent="handleDrop"
      @click="($refs.input as HTMLInputElement).click()"
    >
      <input
        ref="input"
        type="file"
        accept="image/*"
        class="hidden"
        @change="handleChange"
      />
      <span v-if="uploading" class="i-carbon-renew text-4xl text-primary-500 animate-spin" />
      <span v-else class="i-carbon-cloud-upload text-4xl text-gray-400" />
      <p class="mt-2 text-sm text-gray-500 dark:text-gray-400">
        {{ uploading ? '上传中...' : '拖拽图片到此处或点击上传' }}
      </p>
      <p class="mt-1 text-xs text-gray-400 dark:text-gray-500">支持 JPG、PNG、GIF、WebP，最大 5MB</p>
    </div>
    <p v-if="error" class="mt-2 text-sm text-red-500">{{ error }}</p>
  </div>
</template>
