# 使用官方Go镜像作为构建环境
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的包
RUN apk add --no-cache git ca-certificates tzdata

# 复制go.mod和go.sum
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 使用alpine作为运行环境
FROM alpine:latest

# 安装ca-certificates和tzdata
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

WORKDIR /app

# 从builder阶段复制二进制文件
COPY --from=builder /app/main .

# 创建uploads目录
RUN mkdir -p uploads

# 创建非root用户
RUN adduser -D -s /bin/sh appuser
RUN chown -R appuser:appuser /app
USER appuser

# 暴露端口
EXPOSE 8080

# 启动应用
CMD ["./main"]
