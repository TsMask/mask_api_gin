package controller

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"mask_api_gin/src/framework/reqctx"
	"mask_api_gin/src/framework/resp"
	authService "mask_api_gin/src/modules/auth/service"
	"mask_api_gin/src/modules/system/service"
)

// NewSysLogLogin 实例化控制层
var NewSysLogLogin = &SysLogLoginController{
	sysLogLoginService: service.NewSysLogLogin,
	accountService:     authService.NewAccount,
}

// SysLogLoginController 系统登录日志信息 控制层处理
//
// PATH /system/log/login
type SysLogLoginController struct {
	sysLogLoginService *service.SysLogLogin // 系统登录日志服务
	accountService     *authService.Account // 账号身份操作服务
}

// List 系统登录日志列表
//
// GET /list
func (s SysLogLoginController) List(c *gin.Context) {
	query := reqctx.QueryMap(c)
	dataScopeSQL := reqctx.LoginUserToDataScopeSQL(c, "sys_user", "sys_user")
	rows, total := s.sysLogLoginService.FindByPage(query, dataScopeSQL)
	c.JSON(200, resp.OkData(map[string]any{"rows": rows, "total": total}))
}

// Clean 系统登录日志清空
//
// DELETE /clean
func (s SysLogLoginController) Clean(c *gin.Context) {
	rows := s.sysLogLoginService.Clean()
	c.JSON(200, resp.OkData(rows))
}

// Unlock 系统登录日志账户解锁
//
// PUT /unlock/:userName
func (s SysLogLoginController) Unlock(c *gin.Context) {
	userName := c.Param("userName")
	if userName == "" {
		c.JSON(422, resp.CodeMsg(422002, "bind err: userName is empty"))
		return
	}

	ok := s.accountService.CleanLoginRecordCache(userName)
	if ok {
		c.JSON(200, resp.Ok(nil))
		return
	}
	c.JSON(200, resp.Err(nil))
}

// Export 导出系统登录日志信息
//
// GET /export
func (s SysLogLoginController) Export(c *gin.Context) {
	// 查询结果，根据查询条件结果，单页最大值限制
	query := reqctx.QueryMap(c)
	dataScopeSQL := reqctx.LoginUserToDataScopeSQL(c, "sys_user", "sys_user")
	rows, total := s.sysLogLoginService.FindByPage(query, dataScopeSQL)
	if total == 0 {
		c.JSON(200, resp.ErrMsg("export data record as empty"))
		return
	}

	// 导出文件名称
	fileName := fmt.Sprintf("sys_log_login_export_%d_%d.xlsx", len(rows), time.Now().UnixMilli())
	// 导出数据表格
	saveFilePath, err := s.sysLogLoginService.ExportData(rows, fileName)
	if err != nil {
		c.JSON(200, resp.ErrMsg(err.Error()))
		return
	}

	c.FileAttachment(saveFilePath, fileName)
}
