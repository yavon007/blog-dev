[根目录](../CLAUDE.md) > **frontend**

# Frontend 模块文档

前端应用：Vue 3 + TypeScript + Vite + UnoCSS

## 模块职责

博客系统的用户界面，包括：
- 前台：文章展示、标签浏览、评论互动
- 后台：文章管理、分类/标签管理、评论审核

## 入口与启动

**入口文件**: `src/main.ts`

```bash
# 开发环境
cd frontend && pnpm dev

# 类型检查
cd frontend && pnpm type-check

# 生产构建
cd frontend && pnpm build
```

**启动流程**:
1. 创建 Vue 应用
2. 注册 Pinia 状态管理
3. 注册 Vue Router
4. 初始化主题（从 localStorage 恢复）
5. 挂载到 `#app`

## 对外接口

### 路由结构

#### 前台路由 (`FrontLayout`)

| 路径 | 组件 | 说明 |
|------|------|------|
| `/` | `HomeView.vue` | 首页/文章列表 |
| `/post/:slug` | `PostDetailView.vue` | 文章详情 |
| `/tags` | `TagsView.vue` | 标签页 |

#### 后台路由 (`AdminLayout`, 需认证)

| 路径 | 组件 | 说明 |
|------|------|------|
| `/admin` | `DashboardView.vue` | 仪表盘 |
| `/admin/posts` | `PostListView.vue` | 文章列表 |
| `/admin/posts/new` | `EditorView.vue` | 新建文章 |
| `/admin/posts/:id/edit` | `EditorView.vue` | 编辑文章 |
| `/admin/categories` | `CategoryView.vue` | 分类管理 |
| `/admin/tags` | `TagView.vue` | 标签管理 |
| `/admin/comments` | `CommentView.vue` | 评论管理 |

#### 公开路由

| 路径 | 组件 | 说明 |
|------|------|------|
| `/login` | `LoginView.vue` | 登录页 |
| `/:pathMatch(.*)*` | `NotFoundView.vue` | 404 页 |

### API 模块

| 模块 | 文件 | 说明 |
|------|------|------|
| `authApi` | `src/api/auth.ts` | 登录/刷新/登出 |
| `postApi` | `src/api/post.ts` | 文章 CRUD |
| `categoryApi` | `src/api/taxonomy.ts` | 分类管理 |
| `tagApi` | `src/api/taxonomy.ts` | 标签管理 |
| `commentApi` | `src/api/comment.ts` | 评论管理 |

### 类型定义

核心类型位于 `src/types/index.ts`：

- `ApiResponse<T>` - 统一响应格式
- `PagedData<T>` - 分页数据
- `Post` - 文章
- `Category` - 分类
- `Tag` - 标签
- `Comment` - 评论
- `TokenPair` - JWT 令牌对

## 关键依赖与配置

### 主要依赖

- **框架**: `vue@3.5`
- **路由**: `vue-router@4`
- **状态管理**: `pinia@2`
- **HTTP 客户端**: `axios@1`
- **CSS 框架**: `unocss@66`
- **Markdown 编辑器**: `md-editor-v3@4`
- **工具库**: `@vueuse/core`

### Vite 配置

```typescript
// vite.config.ts
{
  server: {
    port: 5173,
    proxy: {
      '/api': 'http://localhost:8080'
    }
  }
}
```

### UnoCSS 配置

- 预设: `presetUno`, `presetTypography`, `presetIcons`
- 主题色: `primary` (蓝色系)
- 快捷方式: `btn`, `btn-primary`, `btn-secondary`, `card`, `input-base`

## 数据模型

### 状态管理

| Store | 文件 | 说明 |
|-------|------|------|
| `useUserStore` | `src/store/user.ts` | 用户认证状态、Token 管理 |
| `useAppStore` | `src/store/app.ts` | 应用主题、全局状态 |

### Token 存储

- Access Token: `localStorage.blog_access_token`
- Refresh Token: `localStorage.blog_refresh_token`

## 测试与质量

### 当前状态

- **单元测试**: 暂无 `*.spec.ts` / `*.test.ts` 文件
- **类型检查**: TypeScript strict 模式启用

```bash
cd frontend && pnpm type-check  # 类型检查
```

### 质量工具建议

- 添加 `src/api/*.spec.ts`
- 添加 `src/store/*.spec.ts`
- 配置 Vitest 进行单元测试

## 常见问题 (FAQ)

**Q: Markdown 编辑器加载慢？**

A: `md-editor-v3` 仅在 `/admin` 路由下使用，建议使用 `defineAsyncComponent` 懒加载。

**Q: 如何切换深色模式？**

A: 调用 `useAppStore().toggleDark()`，状态自动持久化到 localStorage。

**Q: API 请求 401 怎么办？**

A: 响应拦截器会自动处理 401，清除 Token 并跳转登录页。

**Q: 如何添加新的 API？**

A:
1. 在 `src/types/index.ts` 定义类型
2. 在 `src/api/` 下创建或修改模块
3. 使用 `request.get/post/put/delete` 调用

## 相关文件清单

```
frontend/
├── src/
│   ├── main.ts              # 入口
│   ├── App.vue              # 根组件
│   ├── router/index.ts      # 路由配置
│   ├── store/
│   │   ├── user.ts          # 用户状态
│   │   └── app.ts           # 应用状态
│   ├── api/
│   │   ├── auth.ts          # 认证 API
│   │   ├── post.ts          # 文章 API
│   │   ├── taxonomy.ts      # 分类/标签 API
│   │   └── comment.ts       # 评论 API
│   ├── utils/request.ts     # Axios 封装
│   ├── types/index.ts       # 类型定义
│   ├── layouts/
│   │   ├── FrontLayout.vue  # 前台布局
│   │   └── AdminLayout.vue  # 后台布局
│   └── views/
│       ├── front/           # 前台页面
│       └── admin/           # 后台页面
├── vite.config.ts           # Vite 配置
├── uno.config.ts            # UnoCSS 配置
├── tsconfig.json            # TypeScript 配置
├── package.json             # 依赖管理
├── Dockerfile               # Docker 构建
└── nginx.conf               # Nginx 配置
```

## 变更记录 (Changelog)

| 日期 | 变更 |
|------|------|
| 2026-03-09 | 初始化模块文档 |
