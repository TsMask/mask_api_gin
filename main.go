package main

import (
	"context"
	"embed"
	"log"
	"os/signal"
	"syscall"

	"mask_api_gin/src"
)

//go:embed src/assets/**
var assetsDir embed.FS

//go:embed src/config/*.yaml
var configDir embed.FS

func main() {
	src.ConfigurationInit(&assetsDir, &configDir)

	// 信号监听在 main：首个 SIGINT/SIGTERM 取消 ctx 触发停机；
	// 停机期间再次收到信号，NotifyContext 已恢复默认处置，进程立即终止（强制退出）
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// app.Run 阻塞监听，放入后台；启动失败通过 serveErr 投递
	serveErr := make(chan error, 1)
	go func() {
		serveErr <- src.RunServer()
	}()

	select {
	case err := <-serveErr:
		// 启动即失败（如端口占用）：先落盘再退出
		src.ConfigurationClose()
		log.Fatalf("run server error: %v", err)
	case <-ctx.Done():
		// 停机信号：关闭 cron/cache/db 与日志落盘后退出（不排空在途 HTTP 请求）
		src.ConfigurationClose()
	}
	log.Println("server stopped")
}
