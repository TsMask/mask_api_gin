package controller

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"mask_api_gin/src/internal/web/reqctx"
	"mask_api_gin/src/internal/web/resp"
	"mask_api_gin/src/modules/system/service"
)

// NewSysOperationLog 实例化控制层
var NewSysOperationLog = &SysOperationLogController{
	sysOperationLogService: service.NewSysOperationLog,
}

// SysOperationLogController 操作日志记录信息 控制层处理
//
// PATH /system/logs/operation
type SysOperationLogController struct {
	sysOperationLogService *service.SysOperationLog // 操作日志服务
}

// List 操作日志列表
//
// GET /list
func (s SysOperationLogController) List(c *gin.Context) {
	query := reqctx.QueryMap(c)
	dataScopeSQL := service.DataScopeSQL(c, "sys_user", "sys_user")
	rows, total := s.sysOperationLogService.FindByPage(query, dataScopeSQL)
	c.JSON(200, resp.OkData(map[string]any{"rows": rows, "total": total}))
}

// Clean 操作日志清空
//
// DELETE /clean
func (s SysOperationLogController) Clean(c *gin.Context) {
	rows := s.sysOperationLogService.Clean()
	c.JSON(200, resp.OkData(rows))
}

// Export 导出操作日志
//
// GET /export
func (s SysOperationLogController) Export(c *gin.Context) {
	// 查询结果，根据查询条件结果，单页最大值限制
	query := reqctx.QueryMap(c)
	dataScopeSQL := service.DataScopeSQL(c, "sys_user", "sys_user")
	rows, total := s.sysOperationLogService.FindByPage(query, dataScopeSQL)
	if total == 0 {
		c.JSON(200, resp.ErrMsg("export data record as empty"))
		return
	}

	// 导出文件名称
	fileName := fmt.Sprintf("sys_operation_log_export_%d_%d.csv", len(rows), time.Now().UnixMilli())
	// 导出数据表格
	saveFilePath, err := s.sysOperationLogService.ExportData(rows, fileName)
	if err != nil {
		c.JSON(200, resp.ErrMsg(err.Error()))
		return
	}

	c.FileAttachment(saveFilePath, fileName)
}
