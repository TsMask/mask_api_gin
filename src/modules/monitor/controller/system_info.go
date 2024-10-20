package controller

import (
	"mask_api_gin/src/framework/vo/result"
	"mask_api_gin/src/modules/monitor/service"

	"github.com/gin-gonic/gin"
)

// 实例化控制层 SystemInfoController 结构体
var NewSystemInfo = &SystemInfoController{
	systemInfogService: service.NewSystemInfoImpl,
}

// 服务器监控信息
//
// PATH /monitor/system-info
type SystemInfoController struct {
	// 服务器系统相关信息服务
	systemInfogService service.ISystemInfo
}

// 服务器信息
//
// GET /
func (s *SystemInfoController) Info(c *gin.Context) {
	c.JSON(200, result.OkData(map[string]any{
		"project": s.systemInfogService.ProjectInfo(),
		"cpu":     s.systemInfogService.CPUInfo(),
		"memory":  s.systemInfogService.MemoryInfo(),
		"network": s.systemInfogService.NetworkInfo(),
		"time":    s.systemInfogService.TimeInfo(),
		"system":  s.systemInfogService.SystemInfo(),
		"disk":    s.systemInfogService.DiskInfo(),
	}))
}
