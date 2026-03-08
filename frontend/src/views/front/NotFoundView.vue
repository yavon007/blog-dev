<script setup lang="ts">
import { RouterLink } from 'vue-router'
</script>

<template>
  <div class="text-center py-20">
    <span class="i-carbon-search text-7xl text-gray-300 dark:text-gray-600 block mx-auto mb-4" />
    <h1 class="text-4xl font-bold text-gray-300 dark:text-gray-600 mb-2">404</h1>
    <p class="text-gray-500 dark:text-gray-400 mb-6">页面不存在</p>
    <RouterLink to="/" class="btn-primary">返回首页</RouterLink>
  </div>
</template>
