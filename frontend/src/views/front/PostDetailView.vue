<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { postApi } from '@/api/post'
import { commentApi } from '@/api/comment'
import type { Post, Comment } from '@/types'

const route = useRoute()
const router = useRouter()

const post = ref<Post | null>(null)
const comments = ref<Comment[]>([])
const loading = ref(true)

const commentForm = ref({ author_name: '', author_email: '', body: '' })
const submitting = ref(false)
const submitSuccess = ref(false)

async function fetchPost() {
  loading.value = true
  try {
    const res = await postApi.getBySlug(route.params.slug as string) as unknown as { data: Post }
    post.value = res.data
    fetchComments()
  } catch {
    router.push('/404')
  } finally {
    loading.value = false
  }
}

async function fetchComments() {
  try {
    const res = await commentApi.list(route.params.slug as string) as unknown as { data: { items: Comment[] } }
    comments.value = res.data.items ?? []
  } catch {
    comments.value = []
  }
}

async function submitComment() {
  submitting.value = true
  try {
    await commentApi.create(route.params.slug as string, commentForm.value)
    commentForm.value = { author_name: '', author_email: '', body: '' }
    submitSuccess.value = true
    setTimeout(() => { submitSuccess.value = false }, 3000)
  } catch {
    // handle error
  } finally {
    submitting.value = false
  }
}

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' })
}

onMounted(fetchPost)
</script>

<template>
  <div>
    <!-- 加载骨架屏 -->
    <div v-if="loading" class="animate-pulse space-y-4">
      <div class="h-8 bg-gray-200 dark:bg-gray-700 rounded w-3/4" />
      <div class="h-4 bg-gray-200 dark:bg-gray-700 rounded w-1/4" />
      <div class="space-y-2 mt-8">
        <div v-for="i in 5" :key="i" class="h-4 bg-gray-200 dark:bg-gray-700 rounded" :style="`width: ${60 + Math.random() * 40}%`" />
      </div>
    </div>

    <!-- 文章内容 -->
    <article v-else-if="post">
      <header class="mb-8">
        <!-- 封面 -->
        <img
          v-if="post.cover_url"
          :src="post.cover_url"
          :alt="post.title"
          class="w-full h-64 object-cover rounded-xl mb-6"
        />

        <div class="flex items-center gap-2 text-sm text-gray-500 dark:text-gray-400 mb-3">
          <span v-if="post.category_name" class="bg-primary-100 dark:bg-primary-900/30 text-primary-700 dark:text-primary-300 px-2.5 py-0.5 rounded-full text-xs">
            {{ post.category_name }}
          </span>
          <time :datetime="post.published_at ?? post.created_at">
            {{ formatDate(post.published_at ?? post.created_at) }}
          </time>
        </div>

        <h1 class="text-3xl font-bold text-gray-900 dark:text-gray-100 mb-4">
          {{ post.title }}
        </h1>

        <div v-if="post.tags.length > 0" class="flex flex-wrap gap-2">
          <span
            v-for="tag in post.tags"
            :key="tag.id"
            class="text-xs px-2 py-0.5 bg-gray-100 dark:bg-gray-700 rounded-full text-gray-600 dark:text-gray-300"
          >
            #{{ tag.name }}
          </span>
        </div>
      </header>

      <!-- Markdown 渲染内容 -->
      <!-- eslint-disable-next-line vue/no-v-html -->
      <div class="prose dark:prose-invert max-w-none" v-html="post.content_html_cached || post.content_md" />

      <!-- 评论区 -->
      <section class="mt-12 pt-8 border-t border-gray-200 dark:border-gray-700">
        <h2 class="text-xl font-bold mb-6">评论 ({{ comments.length }})</h2>

        <!-- 评论列表 -->
        <div v-if="comments.length > 0" class="space-y-4 mb-8">
          <div
            v-for="comment in comments"
            :key="comment.id"
            class="card p-4"
          >
            <div class="flex items-center gap-2 mb-2">
              <span class="font-medium text-sm">{{ comment.author_name }}</span>
              <time class="text-xs text-gray-500 dark:text-gray-400" :datetime="comment.created_at">
                {{ formatDate(comment.created_at) }}
              </time>
            </div>
            <p class="text-sm text-gray-700 dark:text-gray-300 whitespace-pre-line">{{ comment.body }}</p>
          </div>
        </div>

        <!-- 评论表单 -->
        <div class="card p-6">
          <h3 class="font-semibold mb-4">发表评论</h3>
          <div v-if="submitSuccess" class="mb-4 p-3 bg-green-50 dark:bg-green-900/20 text-green-700 dark:text-green-300 rounded-lg text-sm">
            评论已提交，等待审核后显示。
          </div>
          <form class="space-y-4" @submit.prevent="submitComment">
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium mb-1">昵称 *</label>
                <input v-model="commentForm.author_name" required class="input-base" placeholder="您的昵称" />
              </div>
              <div>
                <label class="block text-sm font-medium mb-1">邮箱 *</label>
                <input v-model="commentForm.author_email" type="email" required class="input-base" placeholder="your@email.com" />
              </div>
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">评论内容 *</label>
              <textarea
                v-model="commentForm.body"
                required
                rows="4"
                class="input-base resize-none"
                placeholder="写下您的想法..."
              />
            </div>
            <button
              type="submit"
              :disabled="submitting"
              class="btn-primary disabled:opacity-60"
            >
              {{ submitting ? '提交中...' : '提交评论' }}
            </button>
          </form>
        </div>
      </section>
    </article>
  </div>
</template>
