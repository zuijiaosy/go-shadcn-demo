package api

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"runtime"
	"slices"
	"strings"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/zuijiaosy/go-shadcn-demo/internal/buildinfo"
)

// 嵌入前端文件系统（编译时）
// all: 前缀确保包含隐藏文件（如 .gitkeep 占位文件），
// 使得在未构建前端时项目也能正常编译
//
//go:embed all:frontend/dist
var frontendFS embed.FS

// Server HTTP 服务器
type Server struct {
	router *gin.Engine
}

// NewServer 创建新的服务器实例
func NewServer() *Server {
	s := &Server{router: gin.Default()}
	s.setupMiddleware()
	s.setupAPIRoutes()
	s.setupStaticFiles()
	return s
}

// Handler 返回 http.Handler（用于外部 http.Server）
func (s *Server) Handler() http.Handler {
	return s.router
}

// setupMiddleware 设置中间件
func (s *Server) setupMiddleware() {
	// Gzip 压缩
	s.router.Use(gzip.Gzip(gzip.DefaultCompression))

	// CORS（开发端口内置，其余通过 CORS_ORIGINS 环境变量追加）
	allowedOrigins := []string{"http://localhost:5173", "http://localhost:3000"}
	if extra := os.Getenv("CORS_ORIGINS"); extra != "" {
		allowedOrigins = append(allowedOrigins, strings.Split(extra, ",")...)
	}
	s.router.Use(func(c *gin.Context) {
		if origin := c.GetHeader("Origin"); slices.Contains(allowedOrigins, origin) {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		}
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	// 安全响应头
	s.router.Use(func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "SAMEORIGIN")
		c.Header("Referrer-Policy", "no-referrer-when-downgrade")
		c.Next()
	})
}

// setupAPIRoutes 设置 API 路由
func (s *Server) setupAPIRoutes() {
	api := s.router.Group("/api")
	{
		api.GET("/health", s.healthCheck)
		api.GET("/version", s.getVersion)
	}
}

// healthCheck 健康检查
func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// getVersion 获取版本信息
func (s *Server) getVersion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version":    buildinfo.Version,
		"git_commit": buildinfo.GitCommit,
		"build_time": buildinfo.BuildTime,
		"go_version": runtime.Version(),
	})
}

// setupStaticFiles 设置静态文件服务（SPA 模式）
func (s *Server) setupStaticFiles() {
	distFS, _ := fs.Sub(frontendFS, "frontend/dist")

	// 未构建前端时目录中只有 .gitkeep 占位文件
	hasFrontend := false
	if _, err := fs.Stat(distFS, "index.html"); err == nil {
		hasFrontend = true
	} else {
		log.Println("[警告] 前端构建产物未找到，静态文件服务已禁用")
		log.Println("[提示] 运行 ./scripts/setup-dev.sh 构建并嵌入前端文件")
	}

	httpFS := http.FS(distFS)
	s.router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		switch {
		// 未匹配的 API 路径返回 JSON 404
		case strings.HasPrefix(path, "/api/"):
			c.JSON(http.StatusNotFound, gin.H{"error": "API endpoint not found"})
		// 前端未构建时给出明确提示
		case !hasFrontend:
			c.String(http.StatusNotFound, "前端未构建，请运行 ./scripts/setup-dev.sh")
		// 静态文件存在则直接返回
		case fileExists(distFS, strings.TrimPrefix(path, "/")):
			c.FileFromFS(strings.TrimPrefix(path, "/"), httpFS)
		// 缺失的静态资源返回 404，避免 SPA 回退掩盖错误
		case strings.HasPrefix(path, "/assets/"):
			c.Status(http.StatusNotFound)
		// 其余路径回退到 index.html（SPA 前端路由）
		default:
			c.FileFromFS("/", httpFS)
		}
	})
}

// fileExists 判断嵌入文件系统中是否存在指定的普通文件
func fileExists(fsys fs.FS, name string) bool {
	if name == "" {
		return false
	}
	info, err := fs.Stat(fsys, name)
	return err == nil && !info.IsDir()
}
