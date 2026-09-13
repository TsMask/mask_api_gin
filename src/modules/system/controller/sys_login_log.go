package controller

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"mask_api_gin/src/internal/web/reqctx"
	"mask_api_gin/src/internal/web/resp"
	"mask_api_gin/src/modules/system/service"
)

// NewSysLoginLog 实例化控制层
var NewSysLoginLog = &SysLoginLogController{
	sysLoginLogService: service.NewSysLoginLog,
}

// SysLoginLogController 系统登录日志信息 控制层处理
//
// PATH /system/logs/login
type SysLoginLogController struct {
	sysLoginLogService *service.SysLoginLog // 系统登录日志服务
}

// List 系统登录日志列表
//
// GET /list
func (s SysLoginLogController) List(c *gin.Context) {
	query := reqctx.QueryMap(c)
	dataScopeSQL := service.DataScopeSQL(c, "sys_user", "sys_user")
	rows, total := s.sysLoginLogService.FindByPage(query, dataScopeSQL)
	c.JSON(200, resp.OkData(map[string]any{"rows": rows, "total": total}))
}

// Clean 系统登录日志清空
//
// DELETE /clean
func (s SysLoginLogController) Clean(c *gin.Context) {
	rows := s.sysLoginLogService.Clean()
	c.JSON(200, resp.OkData(rows))
}

// Export 导出系统登录日志信息
//
// GET /export
func (s SysLoginLogController) Export(c *gin.Context) {
	// 查询结果，根据查询条件结果，单页最大值限制
	query := reqctx.QueryMap(c)
	dataScopeSQL := service.DataScopeSQL(c, "sys_user", "sys_user")
	rows, total := s.sysLoginLogService.FindByPage(query, dataScopeSQL)
	if total == 0 {
		c.JSON(200, resp.ErrMsg("export data record as empty"))
		return
	}

	// 导出文件名称
	fileName := fmt.Sprintf("sys_login_log_export_%d_%d.csv", len(rows), time.Now().UnixMilli())
	// 导出数据表格
	saveFilePath, err := s.sysLoginLogService.ExportData(rows, fileName)
	if err != nil {
		c.JSON(200, resp.ErrMsg(err.Error()))
		return
	}

	c.FileAttachment(saveFilePath, fileName)
}
