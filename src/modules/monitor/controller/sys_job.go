package controller

import (
	"encoding/json"
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
	"github.com/gin-gonic/gin/binding"
	"github.com/xuri/excelize/v2"
)

// 调度任务信息
//
// PATH /monitor/job
var SysJob = &sysJob{
	sysJobService: service.SysJobImpl,
}

type sysJob struct {
	// 调度任务服务
	sysJobService service.ISysJob
}

// 调度任务列表
//
// GET /list
func (s *sysJob) List(c *gin.Context) {
	querys := ctx.QueryMapString(c)
	data := s.sysJobService.SelectJobPage(querys)
	c.JSON(200, result.Ok(data))
}

// 调度任务信息
//
// GET /:jobId
func (s *sysJob) Info(c *gin.Context) {
	jobId := c.Param("jobId")
	if jobId == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	data := s.sysJobService.SelectJobById(jobId)
	if data.JobID == jobId {
		c.JSON(200, result.OkData(data))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 调度任务新增
//
// POST /
func (s *sysJob) Add(c *gin.Context) {
	var body model.SysJob
	err := c.ShouldBindBodyWith(&body, binding.JSON)
	if err != nil || body.JobID != "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 检查cron表达式格式
	if parse.CronExpression(body.CronExpression) == 0 {
		msg := fmt.Sprintf("调度任务新增【%s】失败，Cron表达式不正确", body.JobName)
		c.JSON(200, result.ErrMsg(msg))
		return
	}

	// 检查任务调用传入参数是否json格式
	if body.TargetParams != "" {
		msg := fmt.Sprintf("调度任务新增【%s】失败，任务传入参数json字符串不正确", body.JobName)
		if len(body.TargetParams) < 7 {
			c.JSON(200, result.ErrMsg(msg))
			return
		}
		if !json.Valid([]byte(body.TargetParams)) {
			c.JSON(200, result.ErrMsg(msg))
			return
		}
	}

	// 检查属性值唯一
	uniqueJob := s.sysJobService.CheckUniqueJobName(body.JobName, body.JobGroup, "")
	if !uniqueJob {
		msg := fmt.Sprintf("调度任务新增【%s】失败，同任务组内有相同任务名称", body.JobName)
		c.JSON(200, result.ErrMsg(msg))
		return
	}

	body.CreateBy = ctx.LoginUserToUserName(c)
	insertId := s.sysJobService.InsertJob(body)
	if insertId != "" {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 调度任务修改
//
// PUT /
func (s *sysJob) Edit(c *gin.Context) {
	var body model.SysJob
	err := c.ShouldBindBodyWith(&body, binding.JSON)
	if err != nil || body.JobID == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 检查cron表达式格式
	if parse.CronExpression(body.CronExpression) == 0 {
		msg := fmt.Sprintf("调度任务修改【%s】失败，Cron表达式不正确", body.JobName)
		c.JSON(200, result.ErrMsg(msg))
		return
	}

	// 检查任务调用传入参数是否json格式
	if body.TargetParams != "" {
		msg := fmt.Sprintf("调度任务修改【%s】失败，任务传入参数json字符串不正确", body.JobName)
		if len(body.TargetParams) < 7 {
			c.JSON(200, result.ErrMsg(msg))
			return
		}
		if !json.Valid([]byte(body.TargetParams)) {
			c.JSON(200, result.ErrMsg(msg))
			return
		}
	}

	// 检查属性值唯一
	uniqueJob := s.sysJobService.CheckUniqueJobName(body.JobName, body.JobGroup, body.JobID)
	if !uniqueJob {
		msg := fmt.Sprintf("调度任务修改【%s】失败，同任务组内有相同任务名称", body.JobName)
		c.JSON(200, result.ErrMsg(msg))
		return
	}

	body.UpdateBy = ctx.LoginUserToUserName(c)
	rows := s.sysJobService.UpdateJob(body)
	if rows > 0 {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 调度任务删除
//
// DELETE /:jobIds
func (s *sysJob) Remove(c *gin.Context) {
	jobIds := c.Param("jobIds")
	if jobIds == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}
	// 处理字符转id数组后去重
	ids := strings.Split(jobIds, ",")
	uniqueIDs := parse.RemoveDuplicates(ids)
	if len(uniqueIDs) <= 0 {
		c.JSON(200, result.Err(nil))
		return
	}
	rows, err := s.sysJobService.DeleteJobByIds(uniqueIDs)
	if err != nil {
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}
	msg := fmt.Sprintf("删除成功：%d", rows)
	c.JSON(200, result.OkMsg(msg))
}

// 调度任务修改状态
//
// PUT /changeStatus
func (s *sysJob) Status(c *gin.Context) {
	var body struct {
		// 任务ID
		JobId string `json:"jobId" binding:"required"`
		// 状态
		Status string `json:"status" binding:"required"`
	}
	err := c.ShouldBindBodyWith(&body, binding.JSON)
	if err != nil {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 检查是否存在
	job := s.sysJobService.SelectJobById(body.JobId)
	if job.JobID != body.JobId {
		c.JSON(200, result.ErrMsg("没有权限访问调度任务数据！"))
		return
	}

	// 与旧值相等不变更
	if job.Status == body.Status {
		c.JSON(200, result.ErrMsg("变更状态与旧值相等！"))
		return
	}

	// 更新状态
	job.Status = body.Status
	job.UpdateBy = ctx.LoginUserToUserName(c)
	ok := s.sysJobService.ChangeStatus(job)
	if ok {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 调度任务立即执行一次
//
// PUT /run/:jobId
func (s *sysJob) Run(c *gin.Context) {
	jobId := c.Param("jobId")
	if jobId == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 检查是否存在
	job := s.sysJobService.SelectJobById(jobId)
	if job.JobID != jobId {
		c.JSON(200, result.ErrMsg("没有权限访问调度任务数据！"))
		return
	}

	ok := s.sysJobService.RunQueueJob(job)
	if ok {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 调度任务重置刷新队列
//
// PUT /resetQueueJob
func (s *sysJob) ResetQueueJob(c *gin.Context) {
	s.sysJobService.ResetQueueJob()
	c.JSON(200, result.Ok(nil))
}

// 导出调度任务信息
//
// POST /export
func (s *sysJob) Export(c *gin.Context) {
	// 查询结果，根据查询条件结果，单页最大值限制
	querys := ctx.QueryMapString(c)
	data := s.sysJobService.SelectJobPage(querys)

	// 导出数据组装
	fileName := fmt.Sprintf("job_export_%d_%d.xlsx", data["total"], date.NowTimestamp())
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
	file.SetCellValue(sheet, "A1", "任务编号")
	file.SetCellValue(sheet, "B1", "任务名称")
	file.SetCellValue(sheet, "C1", "任务组名")
	file.SetCellValue(sheet, "D1", "调用目标")
	file.SetCellValue(sheet, "E1", "传入参数")
	file.SetCellValue(sheet, "F1", "执行表达式")
	file.SetCellValue(sheet, "G1", "计划策略")
	file.SetCellValue(sheet, "H1", "并发执行")
	file.SetCellValue(sheet, "I1", "任务状态")
	file.SetCellValue(sheet, "J1", "备注说明")

	for i, row := range data["rows"].([]model.SysJob) {
		idx := i + 2
		file.SetCellValue(sheet, "A"+strconv.Itoa(idx), row.JobID)
		file.SetCellValue(sheet, "B"+strconv.Itoa(idx), row.JobName)
		file.SetCellValue(sheet, "C"+strconv.Itoa(idx), row.JobGroup)
		file.SetCellValue(sheet, "D"+strconv.Itoa(idx), row.InvokeTarget)
		file.SetCellValue(sheet, "E"+strconv.Itoa(idx), row.TargetParams)
		file.SetCellValue(sheet, "F"+strconv.Itoa(idx), row.CronExpression)
		file.SetCellValue(sheet, "G"+strconv.Itoa(idx), row.MisfirePolicy)
		file.SetCellValue(sheet, "H"+strconv.Itoa(idx), row.Concurrent)
		file.SetCellValue(sheet, "I"+strconv.Itoa(idx), row.Status)
		file.SetCellValue(sheet, "J"+strconv.Itoa(idx), row.Remark)

	}

	// 根据指定路径保存文件
	if err := file.SaveAs(fileName); err != nil {
		fmt.Println(err)
	}
	// 导出数据表格
	c.FileAttachment(fileName, fileName)
}
