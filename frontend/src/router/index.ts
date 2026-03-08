import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/store/user'

const routes: RouteRecordRaw[] = [
  // 前台路由
  {
    path: '/',
    component: () => import('@/layouts/FrontLayout.vue'),
    children: [
      {
        path: '',
        name: 'Home',
        component: () => import('@/views/front/HomeView.vue'),
      },
      {
        path: 'post/:slug',
        name: 'PostDetail',
        component: () => import('@/views/front/PostDetailView.vue'),
      },
      {
        path: 'tags',
        name: 'Tags',
        component: () => import('@/views/front/TagsView.vue'),
      },
    ],
  },

  // 管理后台路由
  {
    path: '/admin',
    component: () => import('@/layouts/AdminLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'Dashboard',
        component: () => import('@/views/admin/DashboardView.vue'),
      },
      {
        path: 'posts',
        name: 'PostList',
        component: () => import('@/views/admin/PostListView.vue'),
      },
      {
        path: 'posts/new',
        name: 'PostCreate',
        component: () => import('@/views/admin/EditorView.vue'),
      },
      {
        path: 'posts/:id/edit',
        name: 'PostEdit',
        component: () => import('@/views/admin/EditorView.vue'),
      },
      {
        path: 'categories',
        name: 'Categories',
        component: () => import('@/views/admin/CategoryView.vue'),
      },
      {
        path: 'tags',
        name: 'AdminTags',
        component: () => import('@/views/admin/TagView.vue'),
      },
      {
        path: 'comments',
        name: 'Comments',
        component: () => import('@/views/admin/CommentView.vue'),
      },
    ],
  },

  // 登录页（无需鉴权）
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/admin/LoginView.vue'),
  },

  // 404
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/front/NotFoundView.vue'),
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(_to, _from, savedPosition) {
    if (savedPosition) return savedPosition
    return { top: 0 }
  },
})

// 路由守卫
router.beforeEach((to, _from, next) => {
  const userStore = useUserStore()

  if (to.meta.requiresAuth && !userStore.isAuthenticated) {
    next({ name: 'Login', query: { redirect: to.fullPath } })
  } else {
    next()
  }
})

export default router
