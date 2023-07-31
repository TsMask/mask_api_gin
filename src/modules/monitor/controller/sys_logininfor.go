package controller

import (
	"fmt"
	"mask_api_gin/src/framework/utils/ctx"
	"mask_api_gin/src/framework/utils/date"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/framework/vo/result"
	commonService "mask_api_gin/src/modules/common/service"
	"mask_api_gin/src/modules/monitor/model"
	"mask_api_gin/src/modules/monitor/service"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

// 实例化控制层 SysLogininforController 结构体
var NewSysLogininfor = &SysLogininforController{
	sysLogininforService: service.NewSysLogininforImpl,
	accountService:       commonService.NewAccountImpl,
}

// 登录访问信息
//
// PATH /monitor/logininfor
type SysLogininforController struct {
	// 系统登录访问服务
	sysLogininforService service.ISysLogininfor
	// 账号身份操作服务
	accountService commonService.IAccount
}

// 登录访问列表
//
// GET /list
func (s *SysLogininforController) List(c *gin.Context) {
	querys := ctx.QueryMapString(c)
	data := s.sysLogininforService.SelectLogininforPage(querys)
	c.JSON(200, result.Ok(data))
}

// 登录访问删除
//
// DELETE /:infoIds
func (s *SysLogininforController) Remove(c *gin.Context) {
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
	rows := s.sysLogininforService.DeleteLogininforByIds(uniqueIDs)
	if rows > 0 {
		msg := fmt.Sprintf("删除成功：%d", rows)
		c.JSON(200, result.OkMsg(msg))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 登录访问清空
//
// DELETE /clean
func (s *SysLogininforController) Clean(c *gin.Context) {
	err := s.sysLogininforService.CleanLogininfor()
	if err != nil {
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}
	c.JSON(200, result.Ok(nil))
}

// 登录访问账户解锁
//
// PUT /unlock/:userName
func (s *SysLogininforController) Unlock(c *gin.Context) {
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

// 导出登录访问信息
//
// POST /export
func (s *SysLogininforController) Export(c *gin.Context) {
	// 查询结果，根据查询条件结果，单页最大值限制
	querys := ctx.QueryMapString(c)
	data := s.sysLogininforService.SelectLogininforPage(querys)

	// 导出数据组装
	fileName := fmt.Sprintf("logininfor_export_%d_%d.xlsx", data["total"], date.NowTimestamp())
	file := excelize.NewFile()
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// 创建一个工作表
	sheet := "Sheet1"
	index, err := file.NewSheet(sheet)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 设置工作簿的默认工作表
	file.SetActiveSheet(index)
	// 设置名为 Sheet1 工作表上 A 到 H 列的宽度为 20
	file.SetColWidth("Sheet1", "A", "H", 20)
	// 设置单元格的值
	file.SetCellValue(sheet, "A1", "序号")
	file.SetCellValue(sheet, "B1", "用户账号")
	file.SetCellValue(sheet, "C1", "登录状态")
	file.SetCellValue(sheet, "D1", "登录地址")
	file.SetCellValue(sheet, "E1", "登录地点")
	file.SetCellValue(sheet, "F1", "浏览器")
	file.SetCellValue(sheet, "G1", "操作系统")
	file.SetCellValue(sheet, "H1", "提示消息")
	file.SetCellValue(sheet, "I1", "访问时间")

	for i, row := range data["rows"].([]model.SysLogininfor) {
		idx := i + 2
		file.SetCellValue(sheet, "A"+strconv.Itoa(idx), row.InfoID)
		file.SetCellValue(sheet, "B"+strconv.Itoa(idx), row.UserName)
		file.SetCellValue(sheet, "C"+strconv.Itoa(idx), row.Status)
		file.SetCellValue(sheet, "D"+strconv.Itoa(idx), row.IPAddr)
		file.SetCellValue(sheet, "E"+strconv.Itoa(idx), row.LoginLocation)
		file.SetCellValue(sheet, "F"+strconv.Itoa(idx), row.Browser)
		file.SetCellValue(sheet, "G"+strconv.Itoa(idx), row.OS)
		file.SetCellValue(sheet, "H"+strconv.Itoa(idx), row.Msg)
		file.SetCellValue(sheet, "I"+strconv.Itoa(idx), row.LoginTime)
	}

	// 根据指定路径保存文件
	if err := file.SaveAs(fileName); err != nil {
		fmt.Println(err)
	}
	// 导出数据表格
	c.FileAttachment(fileName, fileName)
}
