<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { tagApi } from '@/api/taxonomy'
import type { Tag } from '@/types'
import { useRouter } from 'vue-router'

const tags = ref<Tag[]>([])
const router = useRouter()

onMounted(async () => {
  const res = await tagApi.list() as unknown as { data: Tag[] }
  tags.value = res.data ?? []
})

function searchByTag(slug: string) {
  router.push({ path: '/', query: { tag: slug } })
}
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold mb-6">所有标签</h1>
    <div class="flex flex-wrap gap-3">
      <button
        v-for="tag in tags"
        :key="tag.id"
        class="px-4 py-2 bg-gray-100 dark:bg-gray-700 rounded-full text-sm hover:bg-primary-100 dark:hover:bg-primary-900/30 hover:text-primary-700 dark:hover:text-primary-300 transition-colors"
        @click="searchByTag(tag.slug)"
      >
        #{{ tag.name }}
        <span v-if="tag.post_count" class="ml-1 text-xs opacity-60">{{ tag.post_count }}</span>
      </button>
    </div>
    <div v-if="tags.length === 0" class="text-center py-16 text-gray-500">
      暂无标签
    </div>
  </div>
</template>
