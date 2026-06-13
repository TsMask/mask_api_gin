package controller

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"mask_api_gin/src/framework/reqctx"
	"mask_api_gin/src/framework/resp"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/modules/monitor/model"
	"mask_api_gin/src/modules/monitor/service"
)

// NewSysJob 实例化控制层
var NewSysJob = &SysJobController{
	sysJobService: service.NewSysJob,
}

// SysJobController 调度任务信息 控制层处理
//
// PATH /monitor/job
type SysJobController struct {
	sysJobService *service.SysJob // 调度任务服务
}

// List 调度任务列表
//
// GET /list
func (s SysJobController) List(c *gin.Context) {
	query := reqctx.QueryMap(c)
	rows, total := s.sysJobService.FindByPage(query)
	c.JSON(200, resp.OkData(map[string]any{"rows": rows, "total": total}))
}

// Info 调度任务信息
//
// GET /:jobId
func (s SysJobController) Info(c *gin.Context) {
	jobId := parse.Number(c.Param("jobId"))
	if jobId <= 0 {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: jobId is empty"))
		return
	}

	jobInfo := s.sysJobService.FindById(jobId)
	if jobInfo.JobId == jobId {
		c.JSON(200, resp.OkData(jobInfo))
		return
	}
	c.JSON(200, resp.Err(nil))
}

// Add 调度任务新增
//
// POST /
func (s SysJobController) Add(c *gin.Context) {
	var body model.SysJob
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_PARSER, errMsgs))
		return
	}
	if body.JobId > 0 {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: jobId not is empty"))
		return
	}

	// 检查cron表达式格式
	if parse.CronExpression(body.CronExpression) == 0 {
		msg := fmt.Sprintf("调度任务新增【%s】失败，Cron表达式不正确", body.JobName)
		c.JSON(200, resp.ErrMsg(msg))
		return
	}

	// 检查任务调用传入参数是否json格式
	if body.TargetParams != "" {
		msg := fmt.Sprintf("调度任务新增【%s】失败，任务传入参数json字符串不正确", body.JobName)
		if len(body.TargetParams) < 7 {
			c.JSON(200, resp.ErrMsg(msg))
			return
		}
		if !json.Valid([]byte(body.TargetParams)) {
			c.JSON(200, resp.ErrMsg(msg))
			return
		}
	}

	// 检查属性值唯一
	uniqueJob := s.sysJobService.CheckUniqueByJobName(body.JobName, body.JobGroup, 0)
	if !uniqueJob {
		msg := fmt.Sprintf("调度任务新增【%s】失败，同任务组内有相同任务名称", body.JobName)
		c.JSON(200, resp.ErrMsg(msg))
		return
	}

	body.CreateBy = reqctx.LoginUserToUserName(c)
	insertId := s.sysJobService.Insert(body)
	if insertId > 0 {
		c.JSON(200, resp.OkData(insertId))
		return
	}
	c.JSON(200, resp.Err(nil))
}

// Edit 调度任务修改
//
// PUT /
func (s SysJobController) Edit(c *gin.Context) {
	var body model.SysJob
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_PARSER, errMsgs))
		return
	}
	if body.JobId <= 0 {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: jobId is empty"))
		return
	}

	// 检查cron表达式格式
	if parse.CronExpression(body.CronExpression) == 0 {
		msg := fmt.Sprintf("调度任务修改【%s】失败，Cron表达式不正确", body.JobName)
		c.JSON(200, resp.ErrMsg(msg))
		return
	}

	// 检查任务调用传入参数是否json格式
	if body.TargetParams != "" {
		msg := fmt.Sprintf("调度任务修改【%s】失败，任务传入参数json字符串不正确", body.JobName)
		if len(body.TargetParams) < 7 {
			c.JSON(200, resp.ErrMsg(msg))
			return
		}
		if !json.Valid([]byte(body.TargetParams)) {
			c.JSON(200, resp.ErrMsg(msg))
			return
		}
	}

	// 检查属性值唯一
	uniqueJob := s.sysJobService.CheckUniqueByJobName(body.JobName, body.JobGroup, body.JobId)
	if !uniqueJob {
		msg := fmt.Sprintf("调度任务修改【%s】失败，同任务组内有相同任务名称", body.JobName)
		c.JSON(200, resp.ErrMsg(msg))
		return
	}

	// 检查是否存在
	jobInfo := s.sysJobService.FindById(body.JobId)
	if jobInfo.JobId != body.JobId {
		c.JSON(200, resp.ErrMsg("没有权限访问调度任务数据！"))
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
	jobInfo.UpdateBy = reqctx.LoginUserToUserName(c)
	rows := s.sysJobService.Update(jobInfo)
	if rows > 0 {
		c.JSON(200, resp.Ok(nil))
		return
	}
	c.JSON(200, resp.Err(nil))
}

// Remove 调度任务删除
//
// DELETE /:jobId
func (s SysJobController) Remove(c *gin.Context) {
	jobId := c.Param("jobId")
	if jobId == "" {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: jobId is empty"))
		return
	}

	// 处理字符转id数组后去重
	uniqueIDs := parse.RemoveDuplicatesToArray(jobId, ",")
	// 转换成int64数组类型
	ids := make([]int64, 0)
	for _, v := range uniqueIDs {
		ids = append(ids, parse.Number(v))
	}

	rows, err := s.sysJobService.DeleteByIds(ids)
	if err != nil {
		c.JSON(200, resp.ErrMsg(err.Error()))
		return
	}
	msg := fmt.Sprintf("删除成功：%d", rows)
	c.JSON(200, resp.OkMsg(msg))
}

// Status 调度任务修改状态
//
// PUT /status
func (s SysJobController) Status(c *gin.Context) {
	var body struct {
		JobId      int64  `json:"jobId" binding:"required"`
		StatusFlag string `json:"statusFlag" binding:"required,oneof=0 1"`
	}
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_PARSER, errMsgs))
		return
	}

	// 检查是否存在
	jobInfo := s.sysJobService.FindById(body.JobId)
	if jobInfo.JobId != body.JobId {
		c.JSON(200, resp.ErrMsg("没有权限访问调度任务数据！"))
		return
	}

	// 与旧值相等不变更
	if jobInfo.StatusFlag == body.StatusFlag {
		c.JSON(200, resp.ErrMsg("变更状态与旧值相等！"))
		return
	}

	// 更新状态
	jobInfo.StatusFlag = body.StatusFlag
	jobInfo.UpdateBy = reqctx.LoginUserToUserName(c)
	rows := s.sysJobService.Update(jobInfo)
	if rows > 0 {
		c.JSON(200, resp.Ok(nil))
		return
	}
	c.JSON(200, resp.Err(nil))
}

// Run 调度任务立即执行一次
//
// PUT /run/:jobId
func (s SysJobController) Run(c *gin.Context) {
	jobId := parse.Number(c.Param("jobId"))
	if jobId <= 0 {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: jobId is empty"))
		return
	}

	// 检查是否存在
	jobInfo := s.sysJobService.FindById(jobId)
	if jobInfo.JobId != jobId {
		c.JSON(200, resp.ErrMsg("没有权限访问调度任务数据！"))
		return
	}
	// 执行一次调度任务
	ok := s.sysJobService.Run(jobInfo)
	if ok {
		c.JSON(200, resp.Ok(nil))
		return
	}
	c.JSON(200, resp.Err(nil))
}

// Reset 调度任务重置刷新队列
//
// PUT /reset
func (s SysJobController) Reset(c *gin.Context) {
	s.sysJobService.Reset()
	c.JSON(200, resp.Ok(nil))
}

// Export 导出调度任务信息
//
// GET /export
func (s SysJobController) Export(c *gin.Context) {
	// 查询结果，根据查询条件结果，单页最大值限制
	query := reqctx.QueryMap(c)
	rows, total := s.sysJobService.FindByPage(query)
	if total == 0 {
		c.JSON(200, resp.ErrMsg("export data record as empty"))
		return
	}

	// 导出文件名称
	fileName := fmt.Sprintf("job_export_%d_%d.csv", len(rows), time.Now().UnixMilli())
	// 导出数据表格
	saveFilePath, err := s.sysJobService.ExportData(rows, fileName)
	if err != nil {
		c.JSON(200, resp.ErrMsg(err.Error()))
		return
	}

	c.FileAttachment(saveFilePath, fileName)
}
