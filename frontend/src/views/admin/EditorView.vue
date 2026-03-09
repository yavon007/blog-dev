<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { MdEditor } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'
import { postApi } from '@/api/post'
import { categoryApi, tagApi } from '@/api/taxonomy'
import type { Category, Tag, CreatePostPayload } from '@/types'

const route = useRoute()
const router = useRouter()
const postId = route.params.id ? Number(route.params.id) : null
const isEdit = Boolean(postId)

const form = ref<CreatePostPayload>({
  title: '',
  slug: '',
  summary: '',
  content_md: '',
  cover_url: '',
  status: 'draft',
  category_id: null,
  tag_ids: [],
})
const categories = ref<Category[]>([])
const tags = ref<Tag[]>([])
const saving = ref(false)
const error = ref('')

// 从标题自动生成 slug
function generateSlug(title: string) {
  return title
    .toLowerCase()
    .replace(/[\s\u4e00-\u9fa5]+/g, '-')
    .replace(/[^a-z0-9-]/g, '')
    .replace(/-+/g, '-')
    .replace(/^-|-$/g, '')
}

function onTitleChange() {
  if (!isEdit || !form.value.slug) {
    form.value.slug = generateSlug(form.value.title)
  }
}

async function fetchData() {
  const [catsRes, tagsRes] = await Promise.all([
    categoryApi.list() as unknown as { data: Category[] },
    tagApi.list() as unknown as { data: Tag[] },
  ])
  categories.value = catsRes.data ?? []
  tags.value = tagsRes.data ?? []

  if (isEdit && postId) {
    const postRes = await postApi.adminGetById(postId) as unknown as { data: typeof form.value & { tags: { id: number }[] } }
    const post = postRes.data
    form.value = {
      title: post.title,
      slug: post.slug,
      summary: post.summary,
      content_md: post.content_md,
      cover_url: post.cover_url,
      status: post.status,
      category_id: post.category_id,
      tag_ids: post.tags?.map((t: { id: number }) => t.id) ?? [],
    }
  }
}

async function save(status: 'draft' | 'published' = form.value.status) {
  saving.value = true
  error.value = ''
  try {
    const payload = { ...form.value, status }
    if (isEdit && postId) {
      await postApi.update(postId, payload)
    } else {
      await postApi.create(payload)
    }
    router.push('/admin/posts')
  } catch (e) {
    error.value = e instanceof Error ? e.message : '保存失败'
  } finally {
    saving.value = false
  }
}

onMounted(fetchData)
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold">{{ isEdit ? '编辑文章' : '新建文章' }}</h1>
      <div class="flex gap-2">
        <button class="btn-secondary" :disabled="saving" @click="save('draft')">
          {{ saving ? '保存中...' : '保存草稿' }}
        </button>
        <button class="btn-primary" :disabled="saving" @click="save('published')">
          发布文章
        </button>
      </div>
    </div>

    <div v-if="error" class="mb-4 p-3 bg-red-50 dark:bg-red-900/20 text-red-700 dark:text-red-300 rounded-lg text-sm">
      {{ error }}
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- 主内容 -->
      <div class="lg:col-span-2 space-y-4">
        <div class="card p-4">
          <input
            v-model="form.title"
            type="text"
            placeholder="文章标题"
            class="w-full text-xl font-bold border-none outline-none bg-transparent"
            @input="onTitleChange"
          />
        </div>

        <div class="card p-4">
          <label class="block text-xs text-gray-500 dark:text-gray-400 mb-1">slug（URL路径）</label>
          <input v-model="form.slug" type="text" placeholder="post-slug" class="input-base text-sm" />
        </div>

        <!-- Markdown 编辑器 -->
        <MdEditor
          v-model="form.content_md"
          style="height: 500px"
          preview-theme="github"
        />
      </div>

      <!-- 右侧设置 -->
      <div class="space-y-4">
        <div class="card p-4 space-y-3">
          <h3 class="font-medium text-sm text-gray-600 dark:text-gray-400">文章摘要</h3>
          <textarea
            v-model="form.summary"
            rows="3"
            class="input-base resize-none text-sm"
            placeholder="文章摘要（可选）"
          />
        </div>

        <div class="card p-4 space-y-3">
          <h3 class="font-medium text-sm text-gray-600 dark:text-gray-400">封面图</h3>
          <input v-model="form.cover_url" type="url" class="input-base text-sm" placeholder="https://..." />
          <img v-if="form.cover_url" :src="form.cover_url" alt="封面预览" class="w-full h-32 object-cover rounded" />
        </div>

        <div class="card p-4 space-y-3">
          <h3 class="font-medium text-sm text-gray-600 dark:text-gray-400">分类</h3>
          <select v-model="form.category_id" class="input-base text-sm">
            <option :value="null">无分类</option>
            <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
          </select>
        </div>

        <div class="card p-4 space-y-3">
          <h3 class="font-medium text-sm text-gray-600 dark:text-gray-400">标签</h3>
          <div class="flex flex-wrap gap-2">
            <label
              v-for="tag in tags"
              :key="tag.id"
              class="flex items-center gap-1.5 cursor-pointer"
            >
              <input
                type="checkbox"
                :value="tag.id"
                :checked="form.tag_ids.includes(tag.id)"
                class="rounded"
                @change="(e) => {
                  const checked = (e.target as HTMLInputElement).checked
                  if (checked) form.tag_ids.push(tag.id)
                  else form.tag_ids = form.tag_ids.filter((id: number) => id !== tag.id)
                }"
              />
              <span class="text-sm">{{ tag.name }}</span>
            </label>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
