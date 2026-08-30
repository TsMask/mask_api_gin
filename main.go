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

	// 首个 SIGINT/SIGTERM 触发优雅停机；停机期间再次收到信号按默认行为强制退出。
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	go func() {
		<-ctx.Done()
		stop()
	}()

	if err := src.RunServer(ctx); err != nil {
		// 启动即失败（如端口占用）：先落盘再退出
		src.ConfigurationClose()
		log.Fatalf("run server error: %v", err)
	}

	// HTTP 请求与调度任务排空后，再释放底层连接。
	src.ConfigurationClose()
	log.Println("server stopped")
}
