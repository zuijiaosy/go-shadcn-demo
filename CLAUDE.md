## 项目说明

### 项目定位
Go + React + shadcn/ui 全栈应用模板，支持单一二进制部署（前端嵌入到 Go 二进制文件）

### 核心设计原则
- **单一二进制**: 前端通过 Go embed 打包进二进制，部署只需一个文件
- **前后端分离开发**: 开发时 Vite 独立运行（HMR），通过代理调用 Go API
- **GHCR 自动发布**: Docker 镜像由 GitHub Actions 构建并推送到 ghcr.io，不做本地打包
- **极小镜像体积**: Docker 多阶段构建，最终镜像仅 40-60MB

### 目录结构
```
cmd/server/          # 应用入口，主程序
internal/api/        # HTTP API 层，路由和处理器
internal/buildinfo/  # 版本信息注入
frontend/            # React + shadcn/ui 前端
  src/components/ui/ # shadcn/ui 组件（Button、Card 等）
scripts/             # setup-dev.sh（构建前端并嵌入 Go）
.github/workflows/   # CI 与 GHCR 镜像发布
```

---

## 对话语言
一直使用中文进行对话

---

## 开发规范

### API 开发
- 所有 API 路由统一使用 `/api/` 前缀
- 错误响应格式: `{"error": "错误描述"}`
- 成功响应使用标准 HTTP 状态码（200, 201, 204）

### 前端组件开发
- 优先使用 shadcn/ui 组件（已集成 Button、Card）
- 自定义组件放在 `frontend/src/components/` 目录
- 使用 Tailwind CSS 进行样式定制
- 组件文件使用 PascalCase 命名（如 `MyComponent.tsx`）

### 代码注释
- 所有代码注释必须使用中文
- Go 函数注释遵循 godoc 格式
- TypeScript 组件使用 JSDoc 注释

---

## 测试和部署流程

### 开发测试
1. **启动开发环境**
   ```bash
   # 方式一: 前后端分离（推荐）
   make dev              # Terminal 1: 启动 Go API (localhost:8080)
   cd frontend && npm run dev  # Terminal 2: 启动 Vite (localhost:5173)

   # 方式二: Go embed 模式
   ./scripts/setup-dev.sh && make dev
   ```

2. **测试方法**
   - **浏览器测试**: 使用 chrome MCP 工具访问 http://localhost:5173
   - **API 测试**: 使用 `curl` 命令或编写 Python 脚本
   ```bash
   curl http://localhost:8080/api/health
   curl http://localhost:8080/api/version
   ```

3. **热重载验证**
   - 修改 Go 代码 → Air 自动重启
   - 修改 React 代码 → Vite HMR 自动更新

### 构建验证
1. **本地构建测试**
   ```bash
   make build           # 编译生产版本
   ./app                # 运行二进制文件
   # 浏览器访问 http://localhost:8080 验证功能
   ```

2. **Docker 部署验证**（镜像由 GitHub Actions 构建发布到 GHCR）
   ```bash
   make docker-up       # 拉取 GHCR 镜像并启动容器
   make docker-logs     # 查看日志
   ```

3. **验证检查项**
   - [ ] 前端页面正常显示
   - [ ] `/api/health` 返回 200
   - [ ] `/api/version` 返回版本信息
   - [ ] 静态资源加载正常（CSS、JS、图片）
   - [ ] Docker 镜像体积在 60MB 以内

### 测试文件管理
- 编写的临时测试文件（.test.py, .test.go 等）使用后记得删除
- 单元测试文件保留在项目中

---

## MCP 工具使用

### context7
获取依赖框架的最新文档信息
```
示例: 获取 Gin 框架的路由中间件文档
示例: 获取 shadcn/ui Button 组件最新 API
```

### chrome (playwright)
调用本地浏览器调试前端网页
```
示例: 打开 http://localhost:5173 测试前端页面
示例: 截图验证样式是否正确
示例: 测试按钮点击交互
```

### postgres
如果项目集成了数据库，使用此工具查询和调试 SQL

---

## 常见任务

### 添加新的 API 端点
编辑 `internal/api/server.go`，在 `setupAPIRoutes` 中添加路由

### 添加 shadcn/ui 组件
```bash
cd frontend
npx shadcn@latest add [component-name]
```

### 修改构建配置
- Go 构建: 编辑 `Makefile` 的 build 目标
- Docker 构建: 编辑 `Dockerfile` 与 `.github/workflows/docker.yml`
- 前端构建: 编辑 `frontend/vite.config.ts`
