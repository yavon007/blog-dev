<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { categoryApi } from '@/api/taxonomy'
import type { Category, CreateCategoryPayload } from '@/types'

const categories = ref<Category[]>([])
const showForm = ref(false)
const editing = ref<Category | null>(null)
const form = ref<CreateCategoryPayload>({ name: '', slug: '', description: '' })

function generateSlug(name: string) {
  return name.toLowerCase().replace(/[\s]+/g, '-').replace(/[^a-z0-9-]/g, '')
}

function openCreate() {
  editing.value = null
  form.value = { name: '', slug: '', description: '' }
  showForm.value = true
}

function openEdit(cat: Category) {
  editing.value = cat
  form.value = { name: cat.name, slug: cat.slug, description: cat.description }
  showForm.value = true
}

async function save() {
  if (editing.value) {
    await categoryApi.update(editing.value.id, form.value)
    await fetchCategories()
  } else {
    await categoryApi.create(form.value)
    await fetchCategories()
  }
  showForm.value = false
}

async function deleteCategory(id: number) {
  if (!confirm('确定删除？')) return
  await categoryApi.delete(id)
  categories.value = categories.value.filter(c => c.id !== id)
}

async function fetchCategories() {
  const res = await categoryApi.list() as unknown as { data: Category[] }
  categories.value = res.data ?? []
}

onMounted(fetchCategories)
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold">分类管理</h1>
      <button class="btn-primary" @click="openCreate">
        <span class="i-carbon-add mr-1" />新建分类
      </button>
    </div>

    <!-- 表单 -->
    <div v-if="showForm" class="card p-6 mb-6">
      <h2 class="font-semibold mb-4">{{ editing ? '编辑分类' : '新建分类' }}</h2>
      <div class="space-y-3">
        <input
          v-model="form.name"
          class="input-base"
          placeholder="分类名称"
          @input="() => { if (!editing) form.slug = generateSlug(form.name) }"
        />
        <input v-model="form.slug" class="input-base" placeholder="slug" />
        <textarea v-model="form.description" class="input-base resize-none" rows="2" placeholder="描述（可选）" />
        <div class="flex gap-2">
          <button class="btn-primary" @click="save">保存</button>
          <button class="btn-secondary" @click="showForm = false">取消</button>
        </div>
      </div>
    </div>

    <!-- 列表 -->
    <div class="card overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-gray-50 dark:bg-gray-700/50">
          <tr>
            <th class="px-4 py-3 text-left font-medium text-gray-600 dark:text-gray-300">名称</th>
            <th class="px-4 py-3 text-left font-medium text-gray-600 dark:text-gray-300">Slug</th>
            <th class="px-4 py-3 text-left font-medium text-gray-600 dark:text-gray-300">文章数</th>
            <th class="px-4 py-3 text-right font-medium text-gray-600 dark:text-gray-300">操作</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-100 dark:divide-gray-700">
          <tr v-for="cat in categories" :key="cat.id" class="hover:bg-gray-50 dark:hover:bg-gray-700/30">
            <td class="px-4 py-3 font-medium">{{ cat.name }}</td>
            <td class="px-4 py-3 text-gray-500 dark:text-gray-400 font-mono text-xs">{{ cat.slug }}</td>
            <td class="px-4 py-3 text-gray-500 dark:text-gray-400">{{ cat.post_count ?? 0 }}</td>
            <td class="px-4 py-3 text-right">
              <button class="p-1.5 hover:bg-gray-100 dark:hover:bg-gray-700 rounded mr-1" @click="openEdit(cat)">
                <span class="i-carbon-edit text-base" />
              </button>
              <button class="p-1.5 hover:bg-red-50 dark:hover:bg-red-900/20 text-red-500 rounded" @click="deleteCategory(cat.id)">
                <span class="i-carbon-trash-can text-base" />
              </button>
            </td>
          </tr>
          <tr v-if="categories.length === 0">
            <td colspan="4" class="px-4 py-12 text-center text-gray-400">暂无分类</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
