package controller

import (
	"mask_api_gin/src/framework/context"
	"mask_api_gin/src/framework/response"
	"mask_api_gin/src/framework/utils/file"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/modules/monitor/model"
	"mask_api_gin/src/modules/monitor/service"
	systemService "mask_api_gin/src/modules/system/service"

	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// NewSysJob 实例化控制层
var NewSysJob = &SysJobController{
	sysJobService:      service.NewSysJob,
	sysDictTypeService: systemService.NewSysDictType,
}

// SysJobController 调度任务信息 控制层处理
//
// PATH /monitor/job
type SysJobController struct {
	sysJobService      *service.SysJob            // 调度任务服务
	sysDictTypeService *systemService.SysDictType // 字典类型服务
}

// List 调度任务列表
//
// GET /list
func (s SysJobController) List(c *gin.Context) {
	query := context.QueryMap(c)
	rows, total := s.sysJobService.FindByPage(query)
	c.JSON(200, response.OkData(map[string]any{"rows": rows, "total": total}))
}

// Info 调度任务信息
//
// GET /:jobId
func (s SysJobController) Info(c *gin.Context) {
	jobId := c.Param("jobId")
	if jobId == "" {
		c.JSON(400, response.CodeMsg(40010, "bind err: jobId is empty"))
		return
	}

	jobInfo := s.sysJobService.FindById(jobId)
	if jobInfo.JobId == jobId {
		c.JSON(200, response.OkData(jobInfo))
		return
	}
	c.JSON(200, response.Err(nil))
}

// Add 调度任务新增
//
// POST /
func (s SysJobController) Add(c *gin.Context) {
	var body model.SysJob
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", response.FormatBindError(err))
		c.JSON(400, response.CodeMsg(40010, errMsgs))
		return
	}
	if body.JobId != "" {
		c.JSON(400, response.CodeMsg(40010, "bind err: jobId not is empty"))
		return
	}

	// 检查cron表达式格式
	if parse.CronExpression(body.CronExpression) == 0 {
		msg := fmt.Sprintf("调度任务新增【%s】失败，Cron表达式不正确", body.JobName)
		c.JSON(200, response.ErrMsg(msg))
		return
	}

	// 检查任务调用传入参数是否json格式
	if body.TargetParams != "" {
		msg := fmt.Sprintf("调度任务新增【%s】失败，任务传入参数json字符串不正确", body.JobName)
		if len(body.TargetParams) < 7 {
			c.JSON(200, response.ErrMsg(msg))
			return
		}
		if !json.Valid([]byte(body.TargetParams)) {
			c.JSON(200, response.ErrMsg(msg))
			return
		}
	}

	// 检查属性值唯一
	uniqueJob := s.sysJobService.CheckUniqueByJobName(body.JobName, body.JobGroup, "")
	if !uniqueJob {
		msg := fmt.Sprintf("调度任务新增【%s】失败，同任务组内有相同任务名称", body.JobName)
		c.JSON(200, response.ErrMsg(msg))
		return
	}

	body.CreateBy = context.LoginUserToUserName(c)
	insertId := s.sysJobService.Insert(body)
	if insertId != "" {
		c.JSON(200, response.OkData(insertId))
		return
	}
	c.JSON(200, response.Err(nil))
}

// Edit 调度任务修改
//
// PUT /
func (s SysJobController) Edit(c *gin.Context) {
	var body model.SysJob
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", response.FormatBindError(err))
		c.JSON(400, response.CodeMsg(40010, errMsgs))
		return
	}
	if body.JobId == "" {
		c.JSON(400, response.CodeMsg(40010, "bind err: jobId is empty"))
		return
	}

	// 检查cron表达式格式
	if parse.CronExpression(body.CronExpression) == 0 {
		msg := fmt.Sprintf("调度任务修改【%s】失败，Cron表达式不正确", body.JobName)
		c.JSON(200, response.ErrMsg(msg))
		return
	}

	// 检查任务调用传入参数是否json格式
	if body.TargetParams != "" {
		msg := fmt.Sprintf("调度任务修改【%s】失败，任务传入参数json字符串不正确", body.JobName)
		if len(body.TargetParams) < 7 {
			c.JSON(200, response.ErrMsg(msg))
			return
		}
		if !json.Valid([]byte(body.TargetParams)) {
			c.JSON(200, response.ErrMsg(msg))
			return
		}
	}

	// 检查属性值唯一
	uniqueJob := s.sysJobService.CheckUniqueByJobName(body.JobName, body.JobGroup, body.JobId)
	if !uniqueJob {
		msg := fmt.Sprintf("调度任务修改【%s】失败，同任务组内有相同任务名称", body.JobName)
		c.JSON(200, response.ErrMsg(msg))
		return
	}

	// 检查是否存在
	jobInfo := s.sysJobService.FindById(body.JobId)
	if jobInfo.JobId != body.JobId {
		c.JSON(200, response.ErrMsg("没有权限访问调度任务数据！"))
		return
	}

	jobInfo.JobName = body.JobName
	jobInfo.JobGroup = body.JobGroup
	jobInfo.InvokeTarget = body.InvokeTarget
	jobInfo.TargetParams = body.TargetParams
	jobInfo.CronExpression = body.CronExpression
	jobInfo.MisfirePolicy = body.MisfirePolicy
	jobInfo.Concurrent = body.Concurrent
	jobInfo.StatusFlag = body.StatusFlag
	jobInfo.SaveLog = body.SaveLog
	jobInfo.Remark = body.Remark
	jobInfo.UpdateBy = context.LoginUserToUserName(c)
	rows := s.sysJobService.Update(jobInfo)
	if rows > 0 {
		c.JSON(200, response.Ok(nil))
		return
	}
	c.JSON(200, response.Err(nil))
}

// Remove 调度任务删除
//
// DELETE /:jobId
func (s SysJobController) Remove(c *gin.Context) {
	jobIdStr := c.Param("jobId")
	jobIds := parse.RemoveDuplicatesToArray(jobIdStr, ",")
	if jobIdStr == "" || len(jobIds) == 0 {
		c.JSON(400, response.CodeMsg(40010, "bind err: jobId is empty"))
		return
	}

	rows, err := s.sysJobService.DeleteByIds(jobIds)
	if err != nil {
		c.JSON(200, response.ErrMsg(err.Error()))
		return
	}
	msg := fmt.Sprintf("删除成功：%d", rows)
	c.JSON(200, response.OkMsg(msg))
}

// Status 调度任务修改状态
//
// POST /status
func (s SysJobController) Status(c *gin.Context) {
	var body struct {
		JobId      string `json:"jobId" binding:"required"`
		StatusFlag string `json:"statusFlag" binding:"required,oneof=0 1 2"`
	}
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", response.FormatBindError(err))
		c.JSON(400, response.CodeMsg(40010, errMsgs))
		return
	}

	// 检查是否存在
	jobInfo := s.sysJobService.FindById(body.JobId)
	if jobInfo.JobId != body.JobId {
		c.JSON(200, response.ErrMsg("没有权限访问调度任务数据！"))
		return
	}

	// 与旧值相等不变更
	if jobInfo.StatusFlag == body.StatusFlag {
		c.JSON(200, response.ErrMsg("变更状态与旧值相等！"))
		return
	}

	// 更新状态
	jobInfo.StatusFlag = body.StatusFlag
	jobInfo.UpdateBy = context.LoginUserToUserName(c)
	rows := s.sysJobService.Update(jobInfo)
	if rows > 0 {
		c.JSON(200, response.Ok(nil))
		return
	}
	c.JSON(200, response.Err(nil))
}

// Run 调度任务立即执行一次
//
// PUT /run/:jobId
func (s SysJobController) Run(c *gin.Context) {
	jobId := c.Param("jobId")
	if jobId == "" {
		c.JSON(400, response.CodeMsg(40010, "bind err: jobId is empty"))
		return
	}

	// 检查是否存在
	jobInfo := s.sysJobService.FindById(jobId)
	if jobInfo.JobId != jobId {
		c.JSON(200, response.ErrMsg("没有权限访问调度任务数据！"))
		return
	}

	ok := s.sysJobService.Run(jobInfo)
	if ok {
		c.JSON(200, response.Ok(nil))
		return
	}
	c.JSON(200, response.Err(nil))
}

// Reset 调度任务重置刷新队列
//
// PUT /reset
func (s SysJobController) Reset(c *gin.Context) {
	s.sysJobService.Reset()
	c.JSON(200, response.Ok(nil))
}

// Export 导出调度任务信息
//
// GET /export
func (s SysJobController) Export(c *gin.Context) {
	// 查询结果，根据查询条件结果，单页最大值限制
	query := context.QueryMap(c)
	rows, total := s.sysJobService.FindByPage(query)
	if total == 0 {
		c.JSON(200, response.CodeMsg(40016, "export data record as empty"))
		return
	}

	// 导出文件名称
	fileName := fmt.Sprintf("job_export_%d_%d.xlsx", len(rows), time.Now().UnixMilli())
	// 第一行表头标题
	headerCells := map[string]string{
		"A1": "任务编号",
		"B1": "任务名称",
		"C1": "任务组名",
		"D1": "调用目标",
		"E1": "传入参数",
		"F1": "执行表达式",
		"G1": "出错策略",
		"H1": "并发执行",
		"I1": "任务状态",
		"J1": "备注说明",
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
		misfirePolicy := "放弃执行"
		if row.MisfirePolicy == "1" {
			misfirePolicy = "立即执行"
		} else if row.MisfirePolicy == "2" {
			misfirePolicy = "执行一次"
		}
		concurrent := "禁止"
		if row.Concurrent == "1" {
			concurrent = "允许"
		}
		// 状态
		statusValue := "失败"
		if row.StatusFlag == "1" {
			statusValue = "成功"
		}
		dataCells = append(dataCells, map[string]any{
			"A" + idx: row.JobId,
			"B" + idx: row.JobName,
			"C" + idx: sysJobGroup,
			"D" + idx: row.InvokeTarget,
			"E" + idx: row.TargetParams,
			"F" + idx: row.CronExpression,
			"G" + idx: misfirePolicy,
			"H" + idx: concurrent,
			"I" + idx: statusValue,
			"J" + idx: row.Remark,
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
