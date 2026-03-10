<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { settingsApi } from '@/api/settings'

interface SiteSettings {
  site_title: string
  site_description: string
  default_meta_title: string
  default_meta_description: string
  og_image_url: string
}

const form = ref<SiteSettings>({
  site_title: '',
  site_description: '',
  default_meta_title: '',
  default_meta_description: '',
  og_image_url: '',
})
const loading = ref(false)
const saving = ref(false)
const saved = ref(false)

async function fetchSettings() {
  loading.value = true
  try {
    const res = await settingsApi.getSeo() as unknown as { data: SiteSettings }
    form.value = res.data
  } finally {
    loading.value = false
  }
}

async function save() {
  saving.value = true
  try {
    await settingsApi.updateSeo(form.value)
    saved.value = true
    setTimeout(() => { saved.value = false }, 3000)
  } finally {
    saving.value = false
  }
}

onMounted(fetchSettings)
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold">网站设置</h1>
    </div>

    <div v-if="loading" class="card p-6 text-center text-gray-400">加载中...</div>

    <div v-else class="space-y-6">
      <!-- 基本信息 -->
      <div class="card p-6">
        <h2 class="font-semibold mb-4">基本信息</h2>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium mb-1">网站标题</label>
            <input v-model="form.site_title" type="text" class="input-base" placeholder="My Blog" />
          </div>
          <div>
            <label class="block text-sm font-medium mb-1">网站描述</label>
            <input v-model="form.site_description" type="text" class="input-base" placeholder="一个技术博客" />
          </div>
        </div>
      </div>

      <!-- SEO 设置 -->
      <div class="card p-6">
        <h2 class="font-semibold mb-4">默认 SEO 设置</h2>
        <p class="text-sm text-gray-500 dark:text-gray-400 mb-4">
          这些设置将用于首页、RSS 等未单独设置 SEO 的页面。文章页面可在编辑时单独设置。
        </p>
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium mb-1">默认 Meta 标题</label>
            <input v-model="form.default_meta_title" type="text" class="input-base" placeholder="My Blog - 技术分享" />
          </div>
          <div>
            <label class="block text-sm font-medium mb-1">默认 Meta 描述</label>
            <textarea
              v-model="form.default_meta_description"
              rows="2"
              class="input-base resize-none"
              placeholder="分享编程技术、开发经验和个人见解"
            />
          </div>
          <div>
            <label class="block text-sm font-medium mb-1">默认 OG 图片</label>
            <input v-model="form.og_image_url" type="url" class="input-base" placeholder="https://example.com/og.png" />
            <img v-if="form.og_image_url" :src="form.og_image_url" alt="OG 预览" class="w-full h-32 object-cover rounded mt-2" />
          </div>
        </div>
      </div>

      <!-- 保存按钮 -->
      <div class="flex items-center gap-4">
        <button
          class="btn-primary"
          :disabled="saving"
          @click="save"
        >
          {{ saving ? '保存中...' : '保存设置' }}
        </button>
        <span v-if="saved" class="text-sm text-green-600 dark:text-green-400">保存成功</span>
      </div>
    </div>
  </div>
</template>
