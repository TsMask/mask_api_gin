package controller

import (
	"fmt"
	"mask_api_gin/src/framework/utils/ctx"
	"mask_api_gin/src/framework/utils/date"
	"mask_api_gin/src/framework/utils/file"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/framework/vo/result"
	commonService "mask_api_gin/src/modules/common/service"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/service"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 实例化控制层 SysLogLoginController 结构体
var NewSysLogLogin = &SysLogLoginController{
	sysLogLoginService: service.NewSysLogLoginImpl,
	accountService:     commonService.NewAccountImpl,
}

// 系统登录日志信息
//
// PATH /system/log/login
type SysLogLoginController struct {
	// 系统登录日志服务
	sysLogLoginService service.ISysLogLogin
	// 账号身份操作服务
	accountService commonService.IAccount
}

// 系统登录日志列表
//
// GET /list
func (s *SysLogLoginController) List(c *gin.Context) {
	querys := ctx.QueryMap(c)
	data := s.sysLogLoginService.SelectSysLogLoginPage(querys)
	c.JSON(200, result.Ok(data))
}

// 系统登录日志删除
//
// DELETE /:infoIds
func (s *SysLogLoginController) Remove(c *gin.Context) {
	infoIds := c.Param("infoIds")
	if infoIds == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 处理字符转id数组后去重
	ids := strings.Split(infoIds, ",")
	uniqueIDs := parse.RemoveDuplicates(ids)
	if len(uniqueIDs) <= 0 {
		c.JSON(200, result.Err(nil))
		return
	}
	rows := s.sysLogLoginService.DeleteSysLogLoginByIds(uniqueIDs)
	if rows > 0 {
		msg := fmt.Sprintf("删除成功：%d", rows)
		c.JSON(200, result.OkMsg(msg))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 系统登录日志清空
//
// DELETE /clean
func (s *SysLogLoginController) Clean(c *gin.Context) {
	err := s.sysLogLoginService.CleanSysLogLogin()
	if err != nil {
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}
	c.JSON(200, result.Ok(nil))
}

// 系统登录日志账户解锁
//
// PUT /unlock/:userName
func (s *SysLogLoginController) Unlock(c *gin.Context) {
	userName := c.Param("userName")
	if userName == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}
	ok := s.accountService.ClearLoginRecordCache(userName)
	if ok {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 导出系统登录日志信息
//
// POST /export
func (s *SysLogLoginController) Export(c *gin.Context) {
	// 查询结果，根据查询条件结果，单页最大值限制
	querys := ctx.BodyJSONMap(c)
	data := s.sysLogLoginService.SelectSysLogLoginPage(querys)
	if data["total"].(int64) == 0 {
		c.JSON(200, result.ErrMsg("导出数据记录为空"))
		return
	}
	rows := data["rows"].([]model.SysLogLogin)

	// 导出文件名称
	fileName := fmt.Sprintf("sys_log_login_export_%d_%d.xlsx", len(rows), time.Now().UnixMilli())
	// 第一行表头标题
	headerCells := map[string]string{
		"A1": "序号",
		"B1": "用户账号",
		"C1": "登录状态",
		"D1": "登录地址",
		"E1": "登录地点",
		"F1": "浏览器",
		"G1": "操作系统",
		"H1": "提示消息",
		"I1": "访问时间",
	}
	// 从第二行开始的数据
	dataCells := make([]map[string]any, 0)
	for i, row := range rows {
		idx := strconv.Itoa(i + 2)
		// 状态
		statusValue := "失败"
		if row.Status == "1" {
			statusValue = "成功"
		}
		dataCells = append(dataCells, map[string]any{
			"A" + idx: row.LoginID,
			"B" + idx: row.UserName,
			"C" + idx: statusValue,
			"D" + idx: row.IPAddr,
			"E" + idx: row.LoginLocation,
			"F" + idx: row.Browser,
			"G" + idx: row.OS,
			"H" + idx: row.Msg,
			"I" + idx: date.ParseDateToStr(row.LoginTime, date.YYYY_MM_DD_HH_MM_SS),
		})
	}

	// 导出数据表格
	saveFilePath, err := file.WriteSheet(headerCells, dataCells, fileName, "")
	if err != nil {
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}

	c.FileAttachment(saveFilePath, fileName)
}
