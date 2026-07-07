.PHONY: help setup dev run build test fmt vet clean install-air docker-up docker-down docker-logs

# 默认目标
help:
	@echo "可用命令:"
	@echo "  make setup        - 构建前端并嵌入 Go（首次运行/模拟生产）"
	@echo "  make dev          - 启动开发服务器（Air 热重载）"
	@echo "  make run          - 直接运行应用（无热重载）"
	@echo "  make build        - 本地构建完整二进制（含前端）"
	@echo "  make test         - 运行 Go 测试"
	@echo "  make fmt          - 格式化 Go 代码"
	@echo "  make vet          - 静态检查 Go 代码"
	@echo "  make clean        - 清理构建产物"
	@echo "  make install-air  - 安装 Air 热重载工具"
	@echo "  make docker-up    - 拉取 GHCR 镜像并启动容器"
	@echo "  make docker-down  - 停止并移除容器"
	@echo "  make docker-logs  - 查看容器日志"
	@echo ""
	@echo "生产镜像由 GitHub Actions 自动构建并发布到 ghcr.io"

# 构建前端并复制到 Go embed 目录
setup:
	@./scripts/setup-dev.sh

# 安装 Air 热重载工具
install-air:
	@echo "📦 安装 Air 热重载工具..."
	@go install github.com/air-verse/air@latest
	@echo "✅ Air 安装完成"

# 开发模式（使用 Air 热重载）
dev:
	@if ! command -v air >/dev/null 2>&1; then \
		echo "❌ Air 未安装，正在安装..."; \
		$(MAKE) install-air; \
	fi
	@echo "🚀 启动开发服务器（热重载）..."
	@air

# 直接运行应用（无热重载）
run:
	@echo "🚀 启动应用..."
	@go run ./cmd/server

# 本地构建完整二进制（前端 + 后端，正式版本信息由 CI 注入）
build: setup
	@echo "🔨 编译二进制..."
	@CGO_ENABLED=0 go build -ldflags="-s -w" -o app ./cmd/server
	@echo "✅ 编译完成: ./app"

# 运行测试
test:
	@echo "🧪 运行测试..."
	@go test -v ./...

# 格式化代码
fmt:
	@gofmt -w .
	@echo "✅ 格式化完成"

# 静态检查
vet:
	@go vet ./...
	@echo "✅ 检查通过"

# 清理构建产物
clean:
	@echo "🧹 清理构建产物..."
	@rm -f app
	@rm -rf tmp frontend/dist
	@echo "✅ 清理完成"

# 拉取 GHCR 镜像并启动容器
docker-up:
	@docker-compose pull
	@docker-compose up -d

# 停止并移除容器
docker-down:
	@docker-compose down

# 查看容器日志
docker-logs:
	@docker-compose logs -f
