package controller

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"mask_api_gin/src/internal/shared/parse"
	"mask_api_gin/src/internal/web/reqctx"
	"mask_api_gin/src/internal/web/resp"
	"mask_api_gin/src/modules/job/service"
)

// NewSysJobLog 实例化控制层
var NewSysJobLog = &SysJobLogController{
	sysJobService:    service.NewSysJob,
	sysJobLogService: service.NewSysJobLog,
}

// SysJobLogController 调度任务日志信息 控制层处理
//
// PATH /api/v1/job/logs
type SysJobLogController struct {
	sysJobService    *service.SysJob    // 调度任务服务
	sysJobLogService *service.SysJobLog // 调度任务日志服务
}

// List 调度任务日志列表
//
// GET /list
func (s SysJobLogController) List(c *gin.Context) {
	query := reqctx.QueryMap(c)
	if idStr := c.Query("id"); idStr != "" {
		if id := parse.Number(idStr); id > 0 {
			job := s.sysJobService.FindById(id)
			if job.ID != id {
				c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: id no found"))
				return
			}
			query["job_name"] = job.JobName
			query["job_group"] = job.JobGroup
		}
	}
	rows, total := s.sysJobLogService.FindByPage(query)
	c.JSON(200, resp.OkData(map[string]any{"rows": rows, "total": total}))
}

// Info 调度任务日志信息
//
// GET /:id
func (s SysJobLogController) Info(c *gin.Context) {
	id := parse.Number(c.Param("id"))
	if id <= 0 {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: id is empty"))
		return
	}

	data := s.sysJobLogService.FindById(id)
	if data.ID == id {
		c.JSON(200, resp.OkData(data))
		return
	}
	c.JSON(200, resp.Err(nil))
}

// Remove 调度任务日志删除
//
// DELETE /:id
func (s SysJobLogController) Remove(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: id is empty"))
		return
	}

	// 处理字符转id数组后去重
	uniqueIDs := parse.RemoveDuplicatesToArray(id, ",")
	// 转换成int64数组类型
	ids := make([]int64, 0)
	for _, v := range uniqueIDs {
		ids = append(ids, parse.Number(v))
	}

	rows := s.sysJobLogService.DeleteByIds(ids)
	if rows > 0 {
		msg := fmt.Sprintf("删除成功：%d", rows)
		c.JSON(200, resp.OkMsg(msg))
		return
	}
	c.JSON(200, resp.Err(nil))
}

// Clean 调度任务日志清空
//
// DELETE /clean
func (s SysJobLogController) Clean(c *gin.Context) {
	rows := s.sysJobLogService.Clean()
	c.JSON(200, resp.OkData(rows))
}

// Export 导出调度任务日志信息
//
// GET /export
func (s SysJobLogController) Export(c *gin.Context) {
	// 查询结果，根据查询条件结果，单页最大值限制
	query := reqctx.QueryMap(c)
	if idStr := c.Query("id"); idStr != "" {
		if id := parse.Number(idStr); id > 0 {
			job := s.sysJobService.FindById(id)
			query["job_name"] = job.JobName
			query["job_group"] = job.JobGroup
		}
	}
	rows, total := s.sysJobLogService.FindByPage(query)
	if total == 0 {
		c.JSON(200, resp.ErrMsg("export data record as empty"))
		return
	}

	// 导出文件名称
	fileName := fmt.Sprintf("job_log_export_%d_%d.csv", len(rows), time.Now().UnixMilli())
	// 导出数据表格
	saveFilePath, err := s.sysJobLogService.ExportData(rows, fileName)
	if err != nil {
		c.JSON(200, resp.ErrMsg(err.Error()))
		return
	}

	c.FileAttachment(saveFilePath, fileName)
}
