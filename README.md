# Blog System

前后端分离博客系统，支持 Markdown 写作、分类标签管理、评论审核等功能。

## 技术栈

| 层级 | 技术 |
|------|------|
| 前端 | Vue 3 + Vite + TypeScript + UnoCSS + Pinia |
| 后端 | Go + Gin |
| 数据库 | PostgreSQL 16 |
| 缓存 | Redis 7 |
| 容器化 | Docker + docker-compose |

## 项目结构

```
blog-dev/
├── frontend/          # Vue 3 前端应用
├── backend/           # Go 后端 API
├── docs/              # 项目文档
│   ├── api.md         # API 文档
│   ├── database.md    # 数据库设计
│   └── architecture.md # 架构说明
├── docker-compose.yml # 一键启动
└── README.md
```

## 快速开始

### 前置要求

- Go 1.23+
- Node.js 20+ / pnpm 9+
- Docker & Docker Compose
- PostgreSQL 16（或使用 docker-compose）

### 使用 docker-compose 启动（推荐）

```bash
# 复制环境变量配置
cp backend/.env.example backend/.env

# 启动所有服务（postgres、redis、backend、frontend）
docker compose up -d

# 查看日志
docker compose logs -f
```

访问：
- 前台博客：http://localhost:5173
- 后台管理：http://localhost:5173/admin
- 后端 API：http://localhost:8080

### 本地开发

**后端：**

```bash
cd backend
cp .env.example .env     # 配置数据库等环境变量
make deps                # 安装依赖
make migrate-up          # 执行数据库迁移
make run                 # 启动开发服务器（:8080）
```

**前端：**

```bash
cd frontend
pnpm install             # 安装依赖
pnpm dev                 # 启动开发服务器（:5173）
```

## 功能特性

### 公开前台
- 文章列表（分页、搜索）
- 文章详情（Markdown 渲染）
- 分类 / 标签筛选
- 评论查看

### 管理后台（/admin）
- 文章管理（新建、编辑、发布/下线）
- Markdown 编辑器（实时预览）
- 分类 & 标签管理
- 评论审核（通过/拒绝/删除）

## 文档

- [API 文档](./docs/api.md)
- [数据库设计](./docs/database.md)
- [系统架构](./docs/architecture.md)

## License

MIT
