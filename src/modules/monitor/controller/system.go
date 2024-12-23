package controller

import (
	"mask_api_gin/src/framework/response"
	"mask_api_gin/src/modules/monitor/service"

	"github.com/gin-gonic/gin"
)

// NewSystem 实例化控制层
var NewSystem = &SystemController{
	systemInfoService: service.NewSystemInfo,
}

// SystemController 服务器监控信息 控制层处理
//
// PATH /monitor/system
type SystemController struct {
	// 服务器系统相关信息服务
	systemInfoService *service.SystemInfo
}

// Info 服务器信息
//
// GET /
func (s SystemController) Info(c *gin.Context) {
	c.JSON(200, response.OkData(map[string]any{
		"project": s.systemInfoService.ProjectInfo(),
		"cpu":     s.systemInfoService.CPUInfo(),
		"memory":  s.systemInfoService.MemoryInfo(),
		"network": s.systemInfoService.NetworkInfo(),
		"time":    s.systemInfoService.TimeInfo(),
		"system":  s.systemInfoService.SystemInfo(),
		"disk":    s.systemInfoService.DiskInfo(),
	}))
}
