package controller

import (
	"mask_api_gin/src/framework/vo/result"
	"mask_api_gin/src/modules/monitor/service"

	"github.com/gin-gonic/gin"
)

// 服务器监控信息
//
// PATH /monitor/server
var ServerController = &serverController{
	systemInfogService: service.SystemInfoImpl,
}

type serverController struct {
	systemInfogService service.ISystemInfo
}

// 服务器信息
//
// GET /
func (s *serverController) Info(c *gin.Context) {
	c.JSON(200, result.OkData(map[string]interface{}{
		"project": s.systemInfogService.ProjectInfo(),
		"cpu":     s.systemInfogService.CPUInfo(),
		"memory":  s.systemInfogService.MemoryInfo(),
		"network": s.systemInfogService.NetworkInfo(),
		"time":    s.systemInfogService.TimeInfo(),
		"system":  s.systemInfogService.SystemInfo(),
		"disk":    s.systemInfogService.DiskInfo(),
	}))
}
