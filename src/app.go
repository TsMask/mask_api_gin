package src

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"mask_api_gin/src/internal/config"
	"mask_api_gin/src/internal/web/engine"
	"mask_api_gin/src/modules/auth"
	"mask_api_gin/src/modules/common"
	"mask_api_gin/src/modules/demo"
	"mask_api_gin/src/modules/job"
	"mask_api_gin/src/modules/monitor"
	"mask_api_gin/src/modules/system"
)

// RunServer 运行服务程序，并在 ctx 取消后排空在途请求后返回。
func RunServer(ctx context.Context) error {
	app := engine.NewApp()

	// 模块加载顺序即初始化顺序；各模块在 Setup 内部消化自己的启动副作用。
	modulesSetup(app)

	// 读取服务配置
	app.ForwardedByClientIP = config.GetBoolOr("server.proxy", false)
	host := config.GetStringOr("server.host", "")
	port := config.GetIntOr("server.port", 0)
	addr := net.JoinHostPort(host, strconv.Itoa(port))

	server := &http.Server{
		Addr:              addr,
		Handler:           app,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	// 启动服务
	serverErr := make(chan error, 1)
	go func() {
		cert := config.GetStringOr("server.cert", "")
		key := config.GetStringOr("server.key", "")
		if cert != "" && key != "" {
			fmt.Printf("\nListening and serving HTTPS on %s \n\n", addr)
			serverErr <- server.ListenAndServeTLS(cert, key)
			return
		}
		fmt.Printf("\nListening and serving HTTP on %s \n\n", addr)
		serverErr <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		return server.Shutdown(shutdownCtx)
	}
}

// modulesSetup 加载全部业务模块。
func modulesSetup(app *gin.Engine) {
	// 注入操作日志与任务执行日志持久化钩子（须在路由注册、任务初始化之前）
	system.RegisterOnOperationLog()
	job.RegisterOnJobLog()

	demo.Setup(app)    // 演示模块（可移除）
	auth.Setup(app)    // 认证模块
	common.Setup(app)  // 通用模块
	system.Setup(app)  // 系统模块
	job.Setup(app)     // 任务模块
	monitor.Setup(app) // 监控模块
}
