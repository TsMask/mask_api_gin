package controller

import (
	"mask_api_gin/src/framework/model/result"
	ctxUtils "mask_api_gin/src/framework/utils/ctx"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/modules/monitor/service"
	"strings"

	"github.com/gin-gonic/gin"
)

// 调度任务日志信息
//
// PATH /monitor/jobLog
var SysJobLog = &sysJobLog{
	sysJobLogService: service.SysJobLogImpl,
}

type sysJobLog struct {
	sysJobLogService service.ISysJobLog
}

// 导出调度任务日志信息
//
// POST /export
func (s *sysJobLog) Export(c *gin.Context) {
	c.JSON(200, result.OkMsg("export"))
}

// 调度任务日志列表
//
// GET /list
func (s *sysJobLog) List(c *gin.Context) {
	// 查询参数转换map
	querys := ctxUtils.QueryMapString(c)
	list := s.sysJobLogService.SelectJobLogPage(querys)
	c.JSON(200, result.Ok(list))
}

// 调度任务日志信息
//
// GET /:jobLogId
func (s *sysJobLog) Info(c *gin.Context) {
	jobLogId := c.Param("jobLogId")
	if jobLogId == "" {
		c.JSON(200, result.Err(nil))
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
func (s *sysJobLog) Remove(c *gin.Context) {
	jobLogIds := c.Param("jobLogIds")
	if jobLogIds == "" {
		c.JSON(200, result.Err(nil))
		return
	}
	// 处理字符转id数组
	ids := strings.Split(jobLogIds, ",")
	if len(ids) <= 0 {
		c.JSON(200, result.Err(nil))
		return
	}
	// 去重id
	uniqueIDs := parse.RemoveDuplicates(ids)
	rows := s.sysJobLogService.DeleteJobLogByIds(uniqueIDs)
	if rows > 0 {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 调度任务日志清空
//
// DELETE /clean
func (s *sysJobLog) Clean(c *gin.Context) {
	err := s.sysJobLogService.CleanJobLog()
	if err != nil {
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}
	c.JSON(200, result.Ok(nil))
}
