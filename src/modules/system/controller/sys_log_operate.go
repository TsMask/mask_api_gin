package controller

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"mask_api_gin/src/framework/reqctx"
	"mask_api_gin/src/framework/resp"
	"mask_api_gin/src/modules/system/service"
)

// NewSysLogOperate 实例化控制层
var NewSysLogOperate = &SysLogOperateController{
	sysLogOperateService: service.NewSysLogOperate,
}

// SysLogOperateController 操作日志记录信息 控制层处理
//
// PATH /system/log/operate
type SysLogOperateController struct {
	sysLogOperateService *service.SysLogOperate // 操作日志服务
}

// List 操作日志列表
//
// GET /list
func (s SysLogOperateController) List(c *gin.Context) {
	query := reqctx.QueryMap(c)
	dataScopeSQL := reqctx.LoginUserToDataScopeSQL(c, "sys_user", "sys_user")
	rows, total := s.sysLogOperateService.FindByPage(query, dataScopeSQL)
	c.JSON(200, resp.OkData(map[string]any{"rows": rows, "total": total}))
}

// Clean 操作日志清空
//
// DELETE /clean
func (s SysLogOperateController) Clean(c *gin.Context) {
	rows := s.sysLogOperateService.Clean()
	c.JSON(200, resp.OkData(rows))
}

// Export 导出操作日志
//
// GET /export
func (s SysLogOperateController) Export(c *gin.Context) {
	// 查询结果，根据查询条件结果，单页最大值限制
	query := reqctx.QueryMap(c)
	dataScopeSQL := reqctx.LoginUserToDataScopeSQL(c, "sys_user", "sys_user")
	rows, total := s.sysLogOperateService.FindByPage(query, dataScopeSQL)
	if total == 0 {
		c.JSON(200, resp.CodeMsg(40016, "export data record as empty"))
		return
	}

	// 导出文件名称
	fileName := fmt.Sprintf("sys_log_operate_export_%d_%d.xlsx", len(rows), time.Now().UnixMilli())
	// 导出数据表格
	saveFilePath, err := s.sysLogOperateService.ExportData(rows, fileName)
	if err != nil {
		c.JSON(200, resp.ErrMsg(err.Error()))
		return
	}

	c.FileAttachment(saveFilePath, fileName)
}
