<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { tagApi } from '@/api/taxonomy'
import type { Tag, CreateTagPayload } from '@/types'

const tags = ref<Tag[]>([])
const showForm = ref(false)
const editing = ref<Tag | null>(null)
const form = ref<CreateTagPayload>({ name: '', slug: '' })

function generateSlug(name: string) {
  return name.toLowerCase().replace(/[\s]+/g, '-').replace(/[^a-z0-9-]/g, '')
}

function openCreate() {
  editing.value = null
  form.value = { name: '', slug: '' }
  showForm.value = true
}

function openEdit(tag: Tag) {
  editing.value = tag
  form.value = { name: tag.name, slug: tag.slug }
  showForm.value = true
}

async function save() {
  if (editing.value) {
    await tagApi.update(editing.value.id, form.value)
  } else {
    await tagApi.create(form.value)
  }
  await fetchTags()
  showForm.value = false
}

async function deleteTag(id: number) {
  if (!confirm('确定删除？')) return
  await tagApi.delete(id)
  tags.value = tags.value.filter(t => t.id !== id)
}

async function fetchTags() {
  const res = await tagApi.list() as unknown as { data: Tag[] }
  tags.value = res.data ?? []
}

onMounted(fetchTags)
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold">标签管理</h1>
      <button class="btn-primary" @click="openCreate">
        <span class="i-carbon-add mr-1" />新建标签
      </button>
    </div>

    <!-- 表单 -->
    <div v-if="showForm" class="card p-6 mb-6">
      <h2 class="font-semibold mb-4">{{ editing ? '编辑标签' : '新建标签' }}</h2>
      <div class="space-y-3">
        <input
          v-model="form.name"
          class="input-base"
          placeholder="标签名称"
          @input="() => { if (!editing) form.slug = generateSlug(form.name) }"
        />
        <input v-model="form.slug" class="input-base" placeholder="slug" />
        <div class="flex gap-2">
          <button class="btn-primary" @click="save">保存</button>
          <button class="btn-secondary" @click="showForm = false">取消</button>
        </div>
      </div>
    </div>

    <!-- 标签云 -->
    <div class="card p-6">
      <div class="flex flex-wrap gap-3">
        <div
          v-for="tag in tags"
          :key="tag.id"
          class="flex items-center gap-2 px-3 py-1.5 bg-gray-100 dark:bg-gray-700 rounded-full"
        >
          <span class="text-sm">#{{ tag.name }}</span>
          <span class="text-xs text-gray-500 dark:text-gray-400">({{ tag.post_count ?? 0 }})</span>
          <button class="text-gray-400 dark:text-gray-500 hover:text-primary-500 dark:hover:text-primary-400" @click="openEdit(tag)">
            <span class="i-carbon-edit text-xs" />
          </button>
          <button class="text-gray-400 dark:text-gray-500 hover:text-red-500 dark:hover:text-red-400" @click="deleteTag(tag.id)">
            <span class="i-carbon-close text-xs" />
          </button>
        </div>
        <div v-if="tags.length === 0" class="text-gray-400 dark:text-gray-500 text-sm">暂无标签</div>
      </div>
    </div>
  </div>
</template>
