# ============================================
# Stage 1: 前端构建（Node.js）
# ============================================
FROM node:22-alpine AS frontend-builder

WORKDIR /build

# 利用 Docker 层缓存（先复制依赖清单）
COPY frontend/package*.json ./
RUN npm ci

# 复制前端源代码并构建
COPY frontend/ ./
RUN npm run build

# ============================================
# Stage 2: 后端构建（Go）
# ============================================
FROM golang:1.24-alpine AS backend-builder

ARG TARGETOS=linux
ARG TARGETARCH
ARG VERSION=dev
ARG GIT_COMMIT=unknown
ARG BUILD_TIME=unknown
# Go 模块代理（GitHub Actions 构建，可通过 --build-arg 覆盖）
ARG GOPROXY=https://proxy.golang.org,direct

WORKDIR /build

# 利用 Docker 层缓存（先下载依赖）
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY cmd/ ./cmd/
COPY internal/ ./internal/

# 从前端构建阶段复制产物到 Go embed 目录
RUN rm -rf internal/api/frontend/dist && \
    mkdir -p internal/api/frontend/dist
COPY --from=frontend-builder /build/dist/. ./internal/api/frontend/dist/

# 验证前端产物已正确复制
RUN test -f ./internal/api/frontend/dist/index.html && \
    echo "✅ Frontend assets verified"

# 编译静态二进制（注入版本信息）
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH:-amd64} \
    go build \
    -ldflags="-s -w \
    -X github.com/zuijiaosy/go-shadcn-demo/internal/buildinfo.Version=${VERSION} \
    -X github.com/zuijiaosy/go-shadcn-demo/internal/buildinfo.GitCommit=${GIT_COMMIT} \
    -X 'github.com/zuijiaosy/go-shadcn-demo/internal/buildinfo.BuildTime=${BUILD_TIME}'" \
    -o /build/app ./cmd/server

# ============================================
# Stage 3: 运行时（最小镜像）
# ============================================
FROM alpine:3.21

ARG VERSION=dev
ARG GIT_COMMIT=unknown
ARG BUILD_TIME=unknown

LABEL org.opencontainers.image.title="go-shadcn-demo" \
      org.opencontainers.image.description="Go + React + shadcn/ui 全栈应用模板" \
      org.opencontainers.image.source="https://github.com/zuijiaosy/go-shadcn-demo" \
      org.opencontainers.image.licenses="MIT" \
      org.opencontainers.image.version="${VERSION}" \
      org.opencontainers.image.revision="${GIT_COMMIT}" \
      org.opencontainers.image.created="${BUILD_TIME}"

WORKDIR /app

# 运行时依赖（wget 由 busybox 提供，无需额外安装）
RUN apk add --no-cache ca-certificates tzdata

# 以非 root 用户运行
RUN adduser -D -u 1000 appuser
USER appuser

# 复制二进制文件（前端已嵌入）
COPY --from=backend-builder /build/app /app/app

EXPOSE 8080
ENV TZ=Asia/Shanghai \
    GIN_MODE=release

# 健康检查（busybox wget）
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD wget -q --spider http://localhost:8080/api/health || exit 1

ENTRYPOINT ["/app/app"]
