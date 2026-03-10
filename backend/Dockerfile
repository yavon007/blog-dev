# 构建阶段
FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bin/api ./cmd/api/main.go

# 运行阶段（最小镜像，非 root 用户）
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -g 1001 appgroup && \
    adduser -u 1001 -G appgroup -s /bin/sh -D appuser

WORKDIR /app

# 创建 uploads 目录并授权给 appuser
RUN mkdir -p /app/uploads && chown -R appuser:appgroup /app/uploads

COPY --from=builder /app/bin/api .

# 以非 root 用户运行
USER appuser

EXPOSE 8080

CMD ["./api"]
