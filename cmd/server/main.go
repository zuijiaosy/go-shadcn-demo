package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/zuijiaosy/go-shadcn-demo/internal/api"
	"github.com/zuijiaosy/go-shadcn-demo/internal/buildinfo"
)

func main() {
	log.Printf("🚀 启动应用 [版本 %s | commit %s | 构建于 %s | %s]",
		buildinfo.Version, buildinfo.GitCommit, buildinfo.BuildTime, runtime.Version())

	addr := os.Getenv("SERVER_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	httpServer := &http.Server{
		Addr:         addr,
		Handler:      api.NewServer().Handler(),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// 监听中断信号，收到后触发优雅关闭
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("🌐 HTTP 服务器启动在 http://localhost%s", addr)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("❌ 服务器启动失败: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("⏳ 正在关闭服务器...")

	// 优雅关闭，最多等待 30 秒
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("❌ 服务器关闭失败: %v", err)
		return
	}
	log.Println("✅ 服务器已安全关闭")
}
