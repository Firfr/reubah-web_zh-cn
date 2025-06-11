# 构建阶段
FROM golang:1.22-alpine AS builder

# 替换 Alpine 镜像源为阿里云
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk update

WORKDIR /app

# 安装构建依赖
RUN apk add --no-cache \
    gcc \
    g++ \
    pkgconfig \
    libwebp-dev \
    musl-dev \
    nodejs \
    npm \
    libheif-dev \
    x265-dev \
    libde265-dev

# 设置 GOPROXY 环境变量以加速 Go 模块下载
ENV GOPROXY=https://goproxy.cn,direct

# 首先复制 go.mod 和 go.sum 以利用 Docker 缓存
COPY go.mod go.sum ./
RUN go mod download

# 复制应用程序代码并进行构建
COPY . ./
WORKDIR /app/templates
RUN npm install && npm run build

# 返回应用目录并构建应用程序
WORKDIR /app
RUN CGO_ENABLED=1 go build -o reubah ./cmd/server

# 运行阶段
FROM alpine:3.19

# 替换 Alpine 镜像源为阿里云
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk update

WORKDIR /app

# 安装运行时依赖
RUN apk add --no-cache \
    libreoffice \
    ttf-liberation \
    libwebp \
    openjdk11-jre \
    curl \
    libheif-dev \
    x265-dev \
    libde265-dev

# 创建 LibreOffice 所需的目录
RUN mkdir -p /tmp/.cache /tmp/.config /tmp/.local

# 复制二进制文件和静态文件
COPY --from=builder /app/reubah /app/reubah
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static

# 创建非 root 用户
RUN addgroup -g 1000 appgroup && \
    adduser -u 1000 -G appgroup -D appuser && \
    chown -R appuser:appgroup /app /tmp/.cache /tmp/.config /tmp/.local

USER appuser

# 设置 LibreOffice 的环境变量
ENV HOME=/tmp

EXPOSE 8081

# 添加健康检查
HEALTHCHECK --interval=30s --timeout=3s \
  CMD curl -f http://localhost:8081/ || exit 1

CMD ["/app/reubah"]

## docker build -t firfe/reubah:2025.06.11 .