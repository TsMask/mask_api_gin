package controller

import (
	"fmt"
	"mask_api_gin/src/framework/utils/ctx"
	"mask_api_gin/src/framework/utils/date"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/framework/vo/result"
	"mask_api_gin/src/modules/monitor/model"
	"mask_api_gin/src/modules/monitor/service"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

// 操作日志记录信息
//
// PATH /monitor/operlog
var SysOperLog = &sysOperLog{
	sysOperLogService: service.SysOperLogImpl,
}

type sysOperLog struct {
	sysOperLogService service.ISysOperLog
}

// 操作日志列表
//
// GET /list
func (s *sysOperLog) List(c *gin.Context) {
	querys := ctx.QueryMapString(c)
	data := s.sysOperLogService.SelectOperLogPage(querys)
	c.JSON(200, result.Ok(data))
}

// 操作日志删除
//
// DELETE /:operIds
func (s *sysOperLog) Remove(c *gin.Context) {
	operIds := c.Param("operIds")
	if operIds == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 处理字符转id数组
	ids := strings.Split(operIds, ",")
	if len(ids) <= 0 {
		c.JSON(200, result.Err(nil))
		return
	}
	// 去重id
	uniqueIDs := parse.RemoveDuplicates(ids)
	rows := s.sysOperLogService.DeleteOperLogByIds(uniqueIDs)
	if rows > 0 {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 操作日志清空
//
// DELETE /clean
func (s *sysOperLog) Clean(c *gin.Context) {
	err := s.sysOperLogService.CleanOperLog()
	if err != nil {
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}
	c.JSON(200, result.Ok(nil))
}

// 导出操作日志
//
// POST /export
func (s *sysOperLog) Export(c *gin.Context) {
	// 查询结果，根据查询条件结果，单页最大值限制
	querys := ctx.QueryMapString(c)
	data := s.sysOperLogService.SelectOperLogPage(querys)

	// 导出数据组装
	fileName := fmt.Sprintf("operlog_export_%d_%d.xlsx", data["total"], date.NowTimestamp())
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
	file.SetCellValue(sheet, "A1", "操作序号")
	file.SetCellValue(sheet, "B1", "操作模块")
	file.SetCellValue(sheet, "C1", "业务类型")
	file.SetCellValue(sheet, "D1", "请求方法")
	file.SetCellValue(sheet, "E1", "请求方式")
	file.SetCellValue(sheet, "F1", "操作类别")
	file.SetCellValue(sheet, "G1", "操作人员")
	file.SetCellValue(sheet, "H1", "部门名称")
	file.SetCellValue(sheet, "I1", "请求地址")
	file.SetCellValue(sheet, "J1", "操作地点")
	file.SetCellValue(sheet, "K1", "请求参数")
	file.SetCellValue(sheet, "L1", "操作消息")
	file.SetCellValue(sheet, "M1", "状态")
	file.SetCellValue(sheet, "N1", "操作时间")

	for i, row := range data["rows"].([]model.SysOperLog) {
		idx := i + 2
		file.SetCellValue(sheet, "A"+strconv.Itoa(idx), row.OperID)
		file.SetCellValue(sheet, "B"+strconv.Itoa(idx), row.Title)
		file.SetCellValue(sheet, "C"+strconv.Itoa(idx), row.BusinessType)
		file.SetCellValue(sheet, "D"+strconv.Itoa(idx), row.Method)
		file.SetCellValue(sheet, "E"+strconv.Itoa(idx), row.RequestMethod)
		file.SetCellValue(sheet, "F"+strconv.Itoa(idx), row.OperatorType)
		file.SetCellValue(sheet, "G"+strconv.Itoa(idx), row.OperName)
		file.SetCellValue(sheet, "H"+strconv.Itoa(idx), row.DeptName)
		file.SetCellValue(sheet, "I"+strconv.Itoa(idx), row.OperURL)
		file.SetCellValue(sheet, "J"+strconv.Itoa(idx), row.OperIP)
		file.SetCellValue(sheet, "K"+strconv.Itoa(idx), row.OperLocation)
		file.SetCellValue(sheet, "L"+strconv.Itoa(idx), row.OperParam)
		if row.Status == "0" {
			file.SetCellValue(sheet, "M"+strconv.Itoa(idx), "失败")
		} else {
			file.SetCellValue(sheet, "M"+strconv.Itoa(idx), "成功")
		}
		file.SetCellValue(sheet, "N"+strconv.Itoa(idx), row.OperTime)
	}

	// 根据指定路径保存文件
	if err := file.SaveAs(fileName); err != nil {
		fmt.Println(err)
	}
	// 导出数据表格
	c.FileAttachment(fileName, fileName)
}
