.PHONY: help dev prod prod-force down logs ps build clean init-env

COMPOSE = docker compose
ENV_FILE = .env

help: ## 显示帮助信息
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

init-env: ## 初始化环境变量文件（首次使用）
	@if [ ! -f $(ENV_FILE) ]; then \
		cp .env.example $(ENV_FILE); \
		echo "已创建 .env 文件，请编辑填写必要配置（POSTGRES_PASSWORD、JWT_SECRET）"; \
	else \
		echo ".env 文件已存在"; \
	fi

dev: ## 启动开发环境（含数据库和缓存）
	$(COMPOSE) up -d postgres redis
	@echo "数据库和缓存已启动，运行 'cd backend && make run' 和 'cd frontend && pnpm dev'"

prod: ## 启动生产环境（所有服务，强制重建容器）
	@[ -f $(ENV_FILE) ] || { echo "错误：请先运行 'make init-env' 并配置 .env"; exit 1; }
	$(COMPOSE) up -d --build --force-recreate

prod-app: ## 仅重建应用服务（backend/frontend），不动数据库
	@[ -f $(ENV_FILE) ] || { echo "错误：请先运行 'make init-env' 并配置 .env"; exit 1; }
	$(COMPOSE) up -d --build --force-recreate backend frontend

prod-force: ## 强制无缓存重建（确保代码更新）
	@[ -f $(ENV_FILE) ] || { echo "错误：请先运行 'make init-env' 并配置 .env"; exit 1; }
	$(COMPOSE) build --no-cache
	$(COMPOSE) up -d --force-recreate

down: ## 停止所有服务
	$(COMPOSE) down

down-volumes: ## 停止所有服务并删除数据卷（危险！会清空数据库）
	$(COMPOSE) down -v

logs: ## 查看所有服务日志
	$(COMPOSE) logs -f

logs-backend: ## 查看后端日志
	$(COMPOSE) logs -f backend

logs-frontend: ## 查看前端日志
	$(COMPOSE) logs -f frontend

ps: ## 查看服务状态
	$(COMPOSE) ps

build: ## 重新构建所有镜像
	$(COMPOSE) build --no-cache

migrate-up: ## 执行数据库迁移（需要服务运行）
	$(COMPOSE) run --rm migrate

shell-db: ## 进入数据库 shell
	$(COMPOSE) exec postgres psql -U $${POSTGRES_USER:-blog} -d $${POSTGRES_DB:-blog}

shell-backend: ## 进入后端容器 shell
	$(COMPOSE) exec backend sh

clean: ## 清理停止的容器和未使用镜像
	docker system prune -f
