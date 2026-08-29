package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"mask_api_gin/src/internal/infra/cache"
	"mask_api_gin/src/internal/infra/database"
	"mask_api_gin/src/internal/web/resp"
)

// NewHealth 实例化健康检查控制层。
var NewHealth = &HealthController{}

// HealthController 进程存活与依赖就绪检查。
type HealthController struct{}

// Live 仅表示进程可处理请求。
func (s HealthController) Live(c *gin.Context) {
	c.JSON(http.StatusOK, resp.OkData(gin.H{"status": "up"}))
}

// Ready 检查默认数据库与 Redis 是否可用。
func (s HealthController) Ready(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	db := database.DB("")
	if db == nil {
		c.JSON(http.StatusServiceUnavailable, resp.ErrMsg("database unavailable"))
		return
	}
	sqlDB, err := db.DB()
	if err == nil {
		err = sqlDB.PingContext(ctx)
	}
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, resp.ErrMsg("database unavailable"))
		return
	}

	rdb := cache.RDB("")
	if rdb == nil {
		c.JSON(http.StatusServiceUnavailable, resp.ErrMsg("cache unavailable"))
		return
	}
	if err := rdb.Ping(ctx).Err(); err != nil {
		c.JSON(http.StatusServiceUnavailable, resp.ErrMsg("cache unavailable"))
		return
	}
	c.JSON(http.StatusOK, resp.OkData(gin.H{"status": "ready"}))
}
