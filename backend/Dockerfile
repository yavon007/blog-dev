# 构建阶段
FROM golang:1.23-alpine AS builder

WORKDIR /app

# 安装依赖
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# 编译
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bin/api ./cmd/api/main.go

# 运行阶段（最小镜像）
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/bin/api .

EXPOSE 8080

CMD ["./api"]
