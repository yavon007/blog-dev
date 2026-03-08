.PHONY: deps lint test run migrate-up migrate-down compose-up compose-down

# 安装依赖
deps:
	go mod tidy && go mod download

# 代码检查
lint:
	golangci-lint run ./...

# 运行测试
test:
	go test ./... -race -coverprofile=coverage.out

# 启动开发服务器
run:
	go run ./cmd/api/main.go

# 数据库迁移 - 向上
migrate-up:
	@./scripts/migrate.sh up

# 数据库迁移 - 回滚一步
migrate-down:
	@./scripts/migrate.sh down 1

# 启动 docker-compose 服务
compose-up:
	docker compose -f deployments/docker-compose.yml up -d

# 停止 docker-compose 服务
compose-down:
	docker compose -f deployments/docker-compose.yml down

# 查看日志
logs:
	docker compose -f deployments/docker-compose.yml logs -f api

# 构建二进制
build:
	go build -o bin/api ./cmd/api/main.go

# 清理
clean:
	rm -rf bin/ coverage.out
