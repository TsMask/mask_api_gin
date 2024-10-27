package controller

import (
	"fmt"
	"mask_api_gin/src/framework/utils/ctx"
	"mask_api_gin/src/framework/utils/date"
	"mask_api_gin/src/framework/utils/file"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/framework/vo/result"
	"mask_api_gin/src/modules/monitor/service"
	systemService "mask_api_gin/src/modules/system/service"
	"strconv"
	"strings"
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
// PATH /monitor/jobLog
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
	query := ctx.QueryMap(c)
	if v, ok := query["jobId"]; ok && v != "" && v != "0" {
		job := s.sysJobLogService.FindById(v.(string))
		query["jobName"] = job.JobName
		query["jobGroup"] = job.JobGroup
	}
	rows, total := s.sysJobLogService.FindByPage(query)
	c.JSON(200, result.OkData(map[string]any{"rows": rows, "total": total}))
}

// Info 调度任务日志信息
//
// GET /:jobLogId
func (s SysJobLogController) Info(c *gin.Context) {
	jobLogId := c.Param("jobLogId")
	if jobLogId == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}
	data := s.sysJobLogService.FindById(jobLogId)
	if data.JobLogId == jobLogId {
		c.JSON(200, result.OkData(data))
		return
	}
	c.JSON(200, result.Err(nil))
}

// Remove 调度任务日志删除
//
// DELETE /:jobLogIds
func (s SysJobLogController) Remove(c *gin.Context) {
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
	rows := s.sysJobLogService.RemoveByIds(uniqueIDs)
	if rows > 0 {
		msg := fmt.Sprintf("删除成功：%d", rows)
		c.JSON(200, result.OkMsg(msg))
		return
	}
	c.JSON(200, result.Err(nil))
}

// Clean 调度任务日志清空
//
// DELETE /clean
func (s SysJobLogController) Clean(c *gin.Context) {
	err := s.sysJobLogService.Clean()
	if err != nil {
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}
	c.JSON(200, result.Ok(nil))
}

// Export 导出调度任务日志信息
//
// POST /export
func (s SysJobLogController) Export(c *gin.Context) {
	// 查询结果，根据查询条件结果，单页最大值限制
	query := ctx.BodyJSONMap(c)
	rows, total := s.sysJobLogService.FindByPage(query)
	if total == 0 {
		c.JSON(200, result.ErrMsg("导出数据记录为空"))
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
			if row.JobGroup == v.DictValue {
				sysJobGroup = v.DictLabel
				break
			}
		}
		// 状态
		statusValue := "失败"
		if row.Status == "1" {
			statusValue = "成功"
		}
		dataCells = append(dataCells, map[string]any{
			"A" + idx: row.JobLogId,
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
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}

	c.FileAttachment(saveFilePath, fileName)
}
