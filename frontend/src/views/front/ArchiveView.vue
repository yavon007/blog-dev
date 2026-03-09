<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { postApi, type ArchiveItem } from '@/api/post'

const router = useRouter()
const archiveItems = ref<ArchiveItem[]>([])
const loading = ref(true)

// Group by year
const groupedArchive = computed(() => {
  const groups: Record<number, { year: number; months: { month: number; count: number }[] }> = {}
  for (const item of archiveItems.value) {
    if (!groups[item.year]) {
      groups[item.year] = { year: item.year, months: [] }
    }
    groups[item.year]!.months.push({ month: item.month, count: item.count })
  }
  return Object.values(groups).sort((a, b) => b.year - a.year)
})

onMounted(async () => {
  try {
    const res = await postApi.getArchive() as unknown as { data: ArchiveItem[] }
    archiveItems.value = res.data ?? []
  } finally {
    loading.value = false
  }
})

function goToMonth(year: number, month: number) {
  router.push(`/archive/${year}/${month}`)
}
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold mb-6">文章归档</h1>

    <div v-if="loading" class="space-y-4">
      <div v-for="i in 3" :key="i" class="h-20 bg-gray-200 dark:bg-gray-700 rounded animate-pulse" />
    </div>

    <div v-else-if="groupedArchive.length === 0" class="text-center py-16 text-gray-500 dark:text-gray-400">
      <span class="i-carbon-document-blank text-6xl block mx-auto mb-4" />
      <p>暂无文章</p>
    </div>

    <div v-else class="space-y-8">
      <div v-for="group in groupedArchive" :key="group.year" class="relative pl-6 border-l-2 border-gray-200 dark:border-gray-700">
        <div class="absolute -left-3 top-0 w-6 h-6 rounded-full bg-primary-500 text-white text-xs flex items-center justify-center font-bold">
          {{ group.year.toString().slice(-2) }}
        </div>
        <h2 class="text-xl font-bold mb-4 text-gray-900 dark:text-gray-100">{{ group.year }}年</h2>
        <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-3">
          <button
            v-for="m in group.months"
            :key="m.month"
            class="p-3 bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 hover:border-primary-500 dark:hover:border-primary-500 hover:shadow-sm transition-all text-left"
            @click="goToMonth(group.year, m.month)"
          >
            <span class="font-medium text-gray-900 dark:text-gray-100">{{ m.month }}月</span>
            <span class="text-sm text-gray-500 dark:text-gray-400 ml-2">{{ m.count }}篇</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
