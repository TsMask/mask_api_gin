package controller

import (
	"mask_api_gin/src/framework/context"
	"mask_api_gin/src/framework/response"
	"mask_api_gin/src/framework/utils/date"
	"mask_api_gin/src/framework/utils/file"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/modules/monitor/service"
	systemService "mask_api_gin/src/modules/system/service"

	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// NewSysJobLog 实例化控制层
var NewSysJobLog = &SysJobLogController{
	sysJobService:      service.NewSysJob,
	sysJobLogService:   service.NewSysJobLog,
	sysDictTypeService: systemService.NewSysDictType,
}

// SysJobLogController 调度任务日志信息 控制层处理
//
// PATH /monitor/job/log
type SysJobLogController struct {
	sysJobService      *service.SysJob            // 调度任务服务
	sysJobLogService   *service.SysJobLog         // 调度任务日志服务
	sysDictTypeService *systemService.SysDictType // 字典类型服务
}

// List 调度任务日志列表
//
// GET /list
func (s SysJobLogController) List(c *gin.Context) {
	// 查询参数转换map
	query := context.QueryMap(c)
	if jobId := c.Query("jobId"); jobId != "" && jobId != "0" {
		job := s.sysJobService.FindById(jobId)
		query["jobName"] = job.JobName
		query["jobGroup"] = job.JobGroup
	}
	rows, total := s.sysJobLogService.FindByPage(query)
	c.JSON(200, response.OkData(map[string]any{"rows": rows, "total": total}))
}

// Info 调度任务日志信息
//
// GET /:logId
func (s SysJobLogController) Info(c *gin.Context) {
	logId := c.Param("logId")
	if logId == "" {
		c.JSON(400, response.CodeMsg(40010, "bind err: logId is empty"))
		return
	}

	jobLogInfo := s.sysJobLogService.FindById(logId)
	if jobLogInfo.LogId == logId {
		c.JSON(200, response.OkData(jobLogInfo))
		return
	}
	c.JSON(200, response.Err(nil))
}

// Remove 调度任务日志删除
//
// DELETE /:logId
func (s SysJobLogController) Remove(c *gin.Context) {
	logIdStr := c.Param("logId")
	logIds := parse.RemoveDuplicatesToArray(logIdStr, ",")
	if logIdStr == "" || len(logIds) == 0 {
		c.JSON(400, response.CodeMsg(40010, "bind err: logId is empty"))
		return
	}

	rows := s.sysJobLogService.RemoveByIds(logIds)
	if rows > 0 {
		msg := fmt.Sprintf("删除成功：%d", rows)
		c.JSON(200, response.OkMsg(msg))
		return
	}
	c.JSON(200, response.Err(nil))
}

// Clean 调度任务日志清空
//
// DELETE /clean
func (s SysJobLogController) Clean(c *gin.Context) {
	rows := s.sysJobLogService.Clean()
	c.JSON(200, response.OkData(rows))
}

// Export 导出调度任务日志信息
//
// GET /export
func (s SysJobLogController) Export(c *gin.Context) {
	// 查询结果，根据查询条件结果，单页最大值限制
	query := context.QueryMap(c)
	if jobId := c.Query("jobId"); jobId != "" && jobId != "0" {
		job := s.sysJobService.FindById(jobId)
		query["jobName"] = job.JobName
		query["jobGroup"] = job.JobGroup
	}
	rows, total := s.sysJobLogService.FindByPage(query)
	if total == 0 {
		c.JSON(200, response.CodeMsg(40016, "export data record as empty"))
		return
	}

	// 导出文件名称
	fileName := fmt.Sprintf("job_log_export_%d_%d.xlsx", len(rows), time.Now().UnixMilli())
	// 第一行表头标题
	headerCells := map[string]string{
		"A1": "日志序号",
		"B1": "任务名称",
		"C1": "任务组名",
		"D1": "调用目标",
		"E1": "传入参数",
		"F1": "日志信息",
		"G1": "执行状态",
		"H1": "记录时间",
	}
	// 读取任务组名字典数据
	dictSysJobGroup := s.sysDictTypeService.FindDataByType("sys_job_group")
	// 从第二行开始的数据
	dataCells := make([]map[string]any, 0)
	for i, row := range rows {
		idx := strconv.Itoa(i + 2)
		// 任务组名
		sysJobGroup := ""
		for _, v := range dictSysJobGroup {
			if row.JobGroup == v.DataValue {
				sysJobGroup = v.DataLabel
				break
			}
		}
		// 状态
		statusValue := "失败"
		if row.StatusFlag == "1" {
			statusValue = "成功"
		}
		dataCells = append(dataCells, map[string]any{
			"A" + idx: row.LogId,
			"B" + idx: row.JobName,
			"C" + idx: sysJobGroup,
			"D" + idx: row.InvokeTarget,
			"E" + idx: row.TargetParams,
			"F" + idx: row.JobMsg,
			"G" + idx: statusValue,
			"H" + idx: date.ParseDateToStr(row.CreateTime, date.YYYY_MM_DD_HH_MM_SS),
		})
	}

	// 导出数据表格
	saveFilePath, err := file.WriteSheet(headerCells, dataCells, fileName, "")
	if err != nil {
		c.JSON(200, response.ErrMsg(err.Error()))
		return
	}

	c.FileAttachment(saveFilePath, fileName)
}
