package controller

import (
	"encoding/json"
	"fmt"
	"mask_api_gin/src/framework/utils/ctx"
	"mask_api_gin/src/framework/utils/file"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/framework/vo/result"
	"mask_api_gin/src/modules/monitor/model"
	"mask_api_gin/src/modules/monitor/service"
	systemService "mask_api_gin/src/modules/system/service"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// 实例化控制层 SysJobLogController 结构体
var NewSysJob = &SysJobController{
	sysJobService:      service.NewSysJobImpl,
	sysDictDataService: systemService.NewSysDictDataImpl,
}

// 调度任务信息
//
// PATH /monitor/job
type SysJobController struct {
	// 调度任务服务
	sysJobService service.ISysJob
	// 字典数据服务
	sysDictDataService systemService.ISysDictData
}

// 调度任务列表
//
// GET /list
func (s *SysJobController) List(c *gin.Context) {
	querys := ctx.QueryMap(c)
	data := s.sysJobService.SelectJobPage(querys)
	c.JSON(200, result.Ok(data))
}

// 调度任务信息
//
// GET /:jobId
func (s *SysJobController) Info(c *gin.Context) {
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
func (s *SysJobController) Add(c *gin.Context) {
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
func (s *SysJobController) Edit(c *gin.Context) {
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
func (s *SysJobController) Remove(c *gin.Context) {
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
func (s *SysJobController) Status(c *gin.Context) {
	var body struct {
		JobId  string `json:"jobId" binding:"required"`
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
	rows := s.sysJobService.UpdateJob(job)
	if rows > 0 {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 调度任务立即执行一次
//
// PUT /run/:jobId
func (s *SysJobController) Run(c *gin.Context) {
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
func (s *SysJobController) ResetQueueJob(c *gin.Context) {
	s.sysJobService.ResetQueueJob()
	c.JSON(200, result.Ok(nil))
}

// 导出调度任务信息
//
// POST /export
func (s *SysJobController) Export(c *gin.Context) {
	// 查询结果，根据查询条件结果，单页最大值限制
	querys := ctx.BodyJSONMap(c)
	data := s.sysJobService.SelectJobPage(querys)
	if data["total"].(int64) == 0 {
		c.JSON(200, result.ErrMsg("导出数据记录为空"))
		return
	}
	rows := data["rows"].([]model.SysJob)

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
	dictSysJobGroup := s.sysDictDataService.SelectDictDataByType("sys_job_group")
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
		if row.Status == "1" {
			statusValue = "成功"
		}
		dataCells = append(dataCells, map[string]any{
			"A" + idx: row.JobID,
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
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}

	c.FileAttachment(saveFilePath, fileName)
}
