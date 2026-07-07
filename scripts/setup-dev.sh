#!/usr/bin/env bash
# 开发环境设置脚本：构建前端并复制到 Go embed 目录
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
cd "${PROJECT_ROOT}"

echo "🔧 开发环境设置..."

# 检查是否在项目根目录
if [ ! -f "go.mod" ]; then
    echo "❌ 错误: 请在项目根目录运行此脚本"
    exit 1
fi

# 始终先构建前端, 再复制到 Go embed 目录
echo ""
echo "📦 构建前端..."
cd frontend

if [ ! -d "node_modules" ]; then
    echo "  📥 安装前端依赖..."
    npm install
fi

echo "  🔨 构建生产版本..."
npm run build
cd ..

# 复制 dist 目录到 Go embed 目录
echo ""
echo "📋 复制前端构建产物到 Go embed 目录..."
rm -rf internal/api/frontend/dist
mkdir -p internal/api/frontend/dist
cp -r frontend/dist/. internal/api/frontend/dist/
# 恢复占位文件，保持 git 工作区干净（该文件已提交，保证未构建时也能编译）
touch internal/api/frontend/dist/.gitkeep
echo "  ✅ 已复制 frontend/dist/ → internal/api/frontend/dist/"

# 验证文件
if [ -f "internal/api/frontend/dist/index.html" ]; then
    echo "  ✅ index.html 已验证"
else
    echo "  ⚠️  警告: index.html 未找到"
fi

echo ""
echo "🎉 开发环境设置完成！"
echo ""
echo "📖 开发方式:"
echo "   方式一：前后端分离开发（推荐）"
echo "     Terminal 1: make dev                    (Go API http://localhost:8080)"
echo "     Terminal 2: cd frontend && npm run dev  (Vite http://localhost:5173)"
echo ""
echo "   方式二：Go embed 模式（模拟生产）"
echo "     make dev  (访问 http://localhost:8080)"
