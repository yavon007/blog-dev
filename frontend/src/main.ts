import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createUnhead } from '@unhead/vue'
import '@unocss/reset/tailwind.css'
import 'virtual:uno.css'
import App from './App.vue'
import router from './router'
import { useAppStore } from './store/app'

const app = createApp(App)
const pinia = createPinia()
const head = createUnhead()

app.use(pinia)
// Provide head instance globally for useHead composable
app.provide('usehead', head)
app.use(router)

// 初始化主题
const appStore = useAppStore()
appStore.initTheme()

app.mount('#app')
