# 生产部署指南

## 目录

- [服务器要求](#服务器要求)
- [快速部署（Docker Compose）](#快速部署docker-compose)
- [HTTPS 配置（反向代理）](#https-配置反向代理)
- [安全加固](#安全加固)
- [备份与恢复](#备份与恢复)
- [监控与日志](#监控与日志)
- [代码更新与重新部署](#代码更新与重新部署)
- [常见问题](#常见问题)

---

## 服务器要求

| 资源 | 最低配置 | 推荐配置 |
|------|----------|----------|
| CPU | 1 核 | 2 核 |
| 内存 | 1 GB | 2 GB |
| 硬盘 | 20 GB | 40 GB SSD |
| 操作系统 | Ubuntu 22.04 / Debian 12 | Ubuntu 22.04 LTS |
| Docker | 24.0+ | 最新稳定版 |
| Docker Compose | 2.20+ | 最新稳定版 |

---

## 快速部署（Docker Compose）

### 1. 克隆项目

```bash
git clone <your-repo-url> /opt/blog
cd /opt/blog
```

### 2. 配置环境变量

```bash
cp .env.example .env
```

编辑 `.env`，**必须填写**以下项：

```bash
# 生成强随机密码
openssl rand -hex 32   # 用于 POSTGRES_PASSWORD
openssl rand -hex 32   # 用于 JWT_SECRET
openssl rand -hex 16   # 用于 REDIS_PASSWORD（可选）
```

`.env` 示例：

```env
APP_ENV=production
LOG_LEVEL=info

POSTGRES_DB=blog
POSTGRES_USER=blog
POSTGRES_PASSWORD=<强随机密码>

REDIS_PASSWORD=<随机密码>

JWT_SECRET=<强随机密钥>
JWT_ACCESS_TTL=15m
JWT_REFRESH_TTL=720h

PUBLIC_ALLOWED_ORIGINS=https://yourdomain.com
BACKEND_PORT=8080
FRONTEND_PORT=80
```

### 3. 创建管理员账号

首次部署需要手动插入管理员，数据库启动后执行：

```bash
# 启动数据库
docker compose up -d postgres
docker compose up migrate

# 生成 bcrypt 密码哈希（Go 工具）
# 或使用在线工具生成，cost=12
docker compose exec postgres psql -U blog -d blog -c "
  INSERT INTO admin_users (email, password_hash, role)
  VALUES ('admin@yourdomain.com', '<bcrypt-hash>', 'owner');
"
```

也可以用 Python 一行命令生成 bcrypt 哈希（更简便）：

```bash
docker run --rm python:3-alpine sh -c \
  'pip install bcrypt -q && python3 -c "import bcrypt; print(bcrypt.hashpw(b\"yourpassword\", bcrypt.gensalt(12)).decode())"'
```

或使用 Go（需要联网下载依赖）：

```bash
docker run --rm golang:1.26-alpine sh -c '
  cd $(mktemp -d) &&
  go mod init tmp &&
  cat > main.go <<'"'"'GOEOF'"'"'
package main

import (
    "fmt"
    "golang.org/x/crypto/bcrypt"
)

func main() {
    h, _ := bcrypt.GenerateFromPassword([]byte("yourpassword"), 12)
    fmt.Println(string(h))
}
GOEOF
  go get golang.org/x/crypto/bcrypt &&
  go run .'
```

### 4. 启动所有服务

```bash
make prod
# 或
docker compose up -d --build
```

### 5. 验证服务状态

```bash
make ps
docker compose logs -f
```

---

## HTTPS 配置（反向代理）

生产环境推荐在 Docker 前面使用宿主机 Nginx 或 Caddy 作为反向代理，处理 SSL 证书。

### 方案 A：Caddy（最简单，自动 HTTPS）

```bash
apt install -y debian-keyring debian-archive-keyring apt-transport-https
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | tee /etc/apt/sources.list.d/caddy-stable.list
apt update && apt install caddy
```

`/etc/caddy/Caddyfile`：

```caddy
yourdomain.com {
    # 前端静态文件
    reverse_proxy localhost:80

    # 自动 HTTPS（Let's Encrypt）
    tls your@email.com
}
```

```bash
systemctl enable --now caddy
```

### 方案 B：宿主机 Nginx + Certbot

```bash
apt install -y nginx certbot python3-certbot-nginx
```

`/etc/nginx/sites-available/blog`：

```nginx
# HTTP → HTTPS 重定向
server {
    listen 80;
    server_name yourdomain.com;
    return 301 https://$host$request_uri;
}

# HTTPS
server {
    listen 443 ssl http2;
    server_name yourdomain.com;

    ssl_certificate     /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;
    ssl_protocols       TLSv1.2 TLSv1.3;
    ssl_ciphers         ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;
    ssl_session_cache   shared:SSL:10m;
    ssl_session_timeout 1d;

    # HSTS（谨慎：一旦设置不易撤销）
    add_header Strict-Transport-Security "max-age=63072000; includeSubDomains" always;

    # 代理到 Docker 前端容器
    location / {
        proxy_pass http://127.0.0.1:80;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto https;
    }
}
```

```bash
# 申请 SSL 证书
certbot --nginx -d yourdomain.com --non-interactive --agree-tos -m your@email.com

# 启用站点
ln -s /etc/nginx/sites-available/blog /etc/nginx/sites-enabled/
nginx -t && systemctl reload nginx

# 自动续期
echo "0 0,12 * * * root certbot renew --quiet" >> /etc/crontab
```

### 配置前端 API 域名

如果使用独立域名（如 `api.yourdomain.com`），修改 `.env`：

```env
PUBLIC_ALLOWED_ORIGINS=https://yourdomain.com
```

同时在宿主机 Nginx 中添加 API 代理到 8080 端口。

---

## 安全加固

### 1. 防火墙配置

```bash
# 只开放必要端口
ufw allow 22/tcp    # SSH（建议改为非标准端口）
ufw allow 80/tcp    # HTTP（用于重定向）
ufw allow 443/tcp   # HTTPS
ufw enable

# 禁止直接访问容器端口（通过 Docker iptables 规则可能绕过 ufw）
# 将 docker-compose.yml 中 BACKEND_PORT 和 FRONTEND_PORT 绑定到 127.0.0.1
```

修改 `docker-compose.yml` 中端口绑定（仅监听本地）：

```yaml
backend:
  ports:
    - "127.0.0.1:${BACKEND_PORT:-8080}:8080"

frontend:
  ports:
    - "127.0.0.1:${FRONTEND_PORT:-80}:80"
```

### 2. 数据库安全

```bash
# 连接到数据库后执行
docker compose exec postgres psql -U blog -d blog

-- 限制 blog 用户权限（不允许创建数据库）
REVOKE CREATEDB FROM blog;
```

### 3. 限制 SSH 访问

```bash
# /etc/ssh/sshd_config
PasswordAuthentication no
PermitRootLogin no
AllowUsers yourusername
```

### 4. 定期更新

```bash
# 安全更新系统包
apt update && apt upgrade -y

# 更新 Docker 镜像
docker compose pull
docker compose up -d
```

---

## 备份与恢复

### 数据库备份

```bash
#!/bin/bash
# 创建 /opt/backup/backup-db.sh

BACKUP_DIR="/opt/backup/postgres"
DATE=$(date +%Y%m%d_%H%M%S)

mkdir -p "$BACKUP_DIR"

# 导出数据库
docker compose -f /opt/blog/docker-compose.yml exec -T postgres \
  pg_dump -U "${POSTGRES_USER:-blog}" "${POSTGRES_DB:-blog}" \
  | gzip > "$BACKUP_DIR/blog_$DATE.sql.gz"

# 保留最近 30 天备份
find "$BACKUP_DIR" -name "*.sql.gz" -mtime +30 -delete

echo "备份完成: blog_$DATE.sql.gz"
```

```bash
chmod +x /opt/backup/backup-db.sh

# 每天凌晨 3 点自动备份
echo "0 3 * * * root /opt/backup/backup-db.sh" >> /etc/crontab
```

### 数据库恢复

```bash
# 从备份文件恢复
gunzip -c /opt/backup/postgres/blog_YYYYMMDD_HHMMSS.sql.gz | \
  docker compose exec -T postgres psql -U blog blog
```

### 异地备份（可选）

将备份文件同步到对象存储（如阿里云 OSS、AWS S3）：

```bash
# 安装 rclone 并配置远程存储
rclone copy /opt/backup/postgres/ remote:your-bucket/blog-backup/
```

---

## 监控与日志

### 查看实时日志

```bash
# 所有服务
make logs

# 单个服务
make logs-backend
make logs-frontend
docker compose logs -f postgres
```

### 日志持久化

修改 `docker-compose.yml` 为 backend 添加日志限制：

```yaml
backend:
  logging:
    driver: "json-file"
    options:
      max-size: "50m"
      max-file: "5"
```

### 健康检查端点

```bash
# 后端健康检查
curl http://localhost:8080/api/v1/health

# 预期响应
{"status":"ok"}
```

### 简易监控脚本

```bash
#!/bin/bash
# /opt/blog/scripts/health-check.sh

SERVICES=("backend" "postgres" "redis" "frontend")
ALERT_EMAIL="admin@yourdomain.com"

for service in "${SERVICES[@]}"; do
  STATUS=$(docker compose -f /opt/blog/docker-compose.yml ps --status running "$service" | grep -c "$service")
  if [ "$STATUS" -eq 0 ]; then
    echo "警告：$service 服务停止！" | mail -s "Blog 服务告警" "$ALERT_EMAIL"
    # 自动重启
    docker compose -f /opt/blog/docker-compose.yml restart "$service"
  fi
done
```

```bash
# 每 5 分钟检查一次
echo "*/5 * * * * root /opt/blog/scripts/health-check.sh" >> /etc/crontab
```

---

## 代码更新与重新部署

### 1. 拉取最新代码

```bash
cd /opt/blog
git pull origin main
```

### 2. 重新构建并更新所有服务

```bash
make prod
# 等价于：docker compose up -d --build
```

Docker Compose 会自动重建有变更的镜像，未变更的服务保持运行，**无需手动停止**。

### 3. 仅更新后端（前端无改动时）

```bash
docker compose up -d --build backend
```

### 4. 执行数据库迁移（如有新迁移文件）

```bash
make migrate-up
# 等价于：docker compose run --rm migrate
```

> 迁移操作幂等，重复执行不会造成数据损坏。

### 5. 验证更新结果

```bash
make ps
curl http://localhost:8080/api/v1/health
# 预期：{"status":"ok"}
```

### 6. 回滚（更新失败时）

```bash
# 查看提交历史
git log --oneline -10

# 回滚到指定版本
git checkout <commit-hash>
make prod
```

---

## 常见问题

### migrate 服务一直重启

```bash
# 查看 migrate 日志
docker compose logs migrate

# 手动运行迁移
docker compose run --rm migrate
```

### 后端启动失败：数据库连接错误

```bash
# 确认数据库已就绪
docker compose exec postgres pg_isready -U blog

# 检查环境变量
docker compose exec backend env | grep DATABASE
```

### 端口冲突

修改 `.env` 中的端口映射：
```env
BACKEND_PORT=8081
FRONTEND_PORT=8080
```

### 忘记管理员密码

```bash
# 生成新的 bcrypt hash 并更新
docker compose exec postgres psql -U blog -d blog -c \
  "UPDATE admin_users SET password_hash='<new-bcrypt-hash>' WHERE email='admin@yourdomain.com';"
```

### 查看容器资源占用

```bash
docker stats
```
