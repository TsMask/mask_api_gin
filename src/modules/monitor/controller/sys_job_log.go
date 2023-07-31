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

// 实例化控制层 SysJobLogController 结构体
var NewSysJobLog = &SysJobLogController{
	sysJobLogService: service.NewSysJobLogImpl,
}

// 调度任务日志信息
//
// PATH /monitor/jobLog
type SysJobLogController struct {
	sysJobLogService service.ISysJobLog
}

// 调度任务日志列表
//
// GET /list
func (s *SysJobLogController) List(c *gin.Context) {
	// 查询参数转换map
	querys := ctx.QueryMapString(c)
	list := s.sysJobLogService.SelectJobLogPage(querys)
	c.JSON(200, result.Ok(list))
}

// 调度任务日志信息
//
// GET /:jobLogId
func (s *SysJobLogController) Info(c *gin.Context) {
	jobLogId := c.Param("jobLogId")
	if jobLogId == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}
	data := s.sysJobLogService.SelectJobLogById(jobLogId)
	if data.JobLogID == jobLogId {
		c.JSON(200, result.OkData(data))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 调度任务日志删除
//
// DELETE /:jobLogIds
func (s *SysJobLogController) Remove(c *gin.Context) {
	jobLogIds := c.Param("jobLogIds")
	if jobLogIds == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 处理字符转id数组后去重
	ids := strings.Split(jobLogIds, ",")
	uniqueIDs := parse.RemoveDuplicates(ids)
	if len(uniqueIDs) <= 0 {
		c.JSON(200, result.Err(nil))
		return
	}
	rows := s.sysJobLogService.DeleteJobLogByIds(uniqueIDs)
	if rows > 0 {
		msg := fmt.Sprintf("删除成功：%d", rows)
		c.JSON(200, result.OkMsg(msg))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 调度任务日志清空
//
// DELETE /clean
func (s *SysJobLogController) Clean(c *gin.Context) {
	err := s.sysJobLogService.CleanJobLog()
	if err != nil {
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}
	c.JSON(200, result.Ok(nil))
}

// 导出调度任务日志信息
//
// POST /export
func (s *SysJobLogController) Export(c *gin.Context) {
	// 查询结果，根据查询条件结果，单页最大值限制
	querys := ctx.QueryMapString(c)
	data := s.sysJobLogService.SelectJobLogPage(querys)

	// 导出数据组装
	fileName := fmt.Sprintf("jobLog_export_%d_%d.xlsx", data["total"], date.NowTimestamp())
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
	file.SetCellValue(sheet, "A1", "日志序号")
	file.SetCellValue(sheet, "B1", "任务名称")
	file.SetCellValue(sheet, "C1", "任务组名")
	file.SetCellValue(sheet, "D1", "调用目标")
	file.SetCellValue(sheet, "E1", "传入参数")
	file.SetCellValue(sheet, "F1", "日志信息")
	file.SetCellValue(sheet, "G1", "执行状态")
	file.SetCellValue(sheet, "H1", "记录时间")

	for i, row := range data["rows"].([]model.SysJobLog) {
		idx := i + 2
		file.SetCellValue(sheet, "A"+strconv.Itoa(idx), row.JobLogID)
		file.SetCellValue(sheet, "B"+strconv.Itoa(idx), row.JobName)
		file.SetCellValue(sheet, "C"+strconv.Itoa(idx), row.JobGroup)
		file.SetCellValue(sheet, "D"+strconv.Itoa(idx), row.InvokeTarget)
		file.SetCellValue(sheet, "E"+strconv.Itoa(idx), row.TargetParams)
		file.SetCellValue(sheet, "F"+strconv.Itoa(idx), row.JobMsg)
		file.SetCellValue(sheet, "G"+strconv.Itoa(idx), row.Status)
		file.SetCellValue(sheet, "H"+strconv.Itoa(idx), row.CreateTime)

	}

	// 根据指定路径保存文件
	if err := file.SaveAs(fileName); err != nil {
		fmt.Println(err)
	}
	// 导出数据表格
	c.FileAttachment(fileName, fileName)
}
