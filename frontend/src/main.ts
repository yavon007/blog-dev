import { createApp } from 'vue'
import { createPinia } from 'pinia'
import '@unocss/reset/tailwind.css'
import 'virtual:uno.css'
import App from './App.vue'
import router from './router'
import { useAppStore } from './store/app'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)

// 初始化主题
const appStore = useAppStore()
appStore.initTheme()

app.mount('#app')
