import { useEffect, useState } from 'react'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'

interface VersionInfo {
  version: string
  git_commit: string
  build_time: string
  go_version: string
}

function App() {
  const [versionInfo, setVersionInfo] = useState<VersionInfo | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetch('/api/version')
      .then((res) => res.json())
      .then((data) => {
        setVersionInfo(data)
        setLoading(false)
      })
      .catch((err) => {
        console.error('获取版本信息失败:', err)
        setLoading(false)
      })
  }, [])

  return (
    <div className="min-h-screen bg-background flex items-center justify-center p-4">
      <div className="w-full max-w-4xl space-y-8">
        {/* 欢迎卡片 */}
        <Card>
          <CardHeader>
            <CardTitle>Go + React + shadcn/ui 全栈模板</CardTitle>
            <CardDescription>
              单一二进制部署的现代化全栈应用模板，前端通过 Go embed 嵌入
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <h3 className="font-semibold mb-2">后端技术栈</h3>
                <ul className="text-sm text-muted-foreground space-y-1">
                  <li>• Go 1.24+ / Gin 框架</li>
                  <li>• embed.FS 前端嵌入</li>
                  <li>• Gzip 压缩 + CORS</li>
                  <li>• Docker 多阶段构建</li>
                </ul>
              </div>
              <div>
                <h3 className="font-semibold mb-2">前端技术栈</h3>
                <ul className="text-sm text-muted-foreground space-y-1">
                  <li>• React 19 + TypeScript</li>
                  <li>• Vite 7 构建工具</li>
                  <li>• Tailwind CSS 样式</li>
                  <li>• shadcn/ui 组件库</li>
                </ul>
              </div>
            </div>
          </CardContent>
          <CardFooter className="flex gap-2">
            <Button onClick={() => window.open('/api/health', '_blank')}>
              健康检查
            </Button>
            <Button variant="outline" onClick={() => window.open('/api/version', '_blank')}>
              版本信息 API
            </Button>
          </CardFooter>
        </Card>

        {/* 版本信息卡片 */}
        <Card>
          <CardHeader>
            <CardTitle>版本信息</CardTitle>
            <CardDescription>从 Go 后端获取的构建信息</CardDescription>
          </CardHeader>
          <CardContent>
            {loading ? (
              <p className="text-sm text-muted-foreground">加载中...</p>
            ) : versionInfo ? (
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="space-y-2">
                  <div>
                    <span className="text-sm font-medium">版本号：</span>
                    <span className="text-sm text-muted-foreground">
                      {versionInfo.version}
                    </span>
                  </div>
                  <div>
                    <span className="text-sm font-medium">Git Commit：</span>
                    <span className="text-sm text-muted-foreground font-mono">
                      {versionInfo.git_commit.substring(0, 8)}
                    </span>
                  </div>
                </div>
                <div className="space-y-2">
                  <div>
                    <span className="text-sm font-medium">构建时间：</span>
                    <span className="text-sm text-muted-foreground">
                      {versionInfo.build_time}
                    </span>
                  </div>
                  <div>
                    <span className="text-sm font-medium">Go 版本：</span>
                    <span className="text-sm text-muted-foreground">
                      {versionInfo.go_version}
                    </span>
                  </div>
                </div>
              </div>
            ) : (
              <p className="text-sm text-destructive">
                无法连接到后端 API，请确保 Go 服务器正在运行
              </p>
            )}
          </CardContent>
        </Card>

        {/* 快速开始卡片 */}
        <Card>
          <CardHeader>
            <CardTitle>快速开始</CardTitle>
            <CardDescription>开发和部署指南</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <div>
                <h3 className="font-semibold mb-2">开发模式</h3>
                <pre className="bg-muted p-3 rounded-md text-sm overflow-x-auto">
                  <code>{`# Terminal 1: Go 后端热重载
make dev

# Terminal 2: Vite 前端 HMR
cd frontend && npm run dev`}</code>
                </pre>
              </div>
              <div>
                <h3 className="font-semibold mb-2">生产构建</h3>
                <pre className="bg-muted p-3 rounded-md text-sm overflow-x-auto">
                  <code>{`# Docker 部署
docker-compose up -d

# 本地编译
./scripts/build.sh && ./app`}</code>
                </pre>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}

export default App
