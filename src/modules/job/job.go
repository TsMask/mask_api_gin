package job

import (
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"

	"mask_api_gin/src/internal/infra/cron"
	"mask_api_gin/src/internal/shared/constants"
	"mask_api_gin/src/internal/shared/logger"
	"mask_api_gin/src/internal/web/middleware"
	"mask_api_gin/src/modules/job/controller"
	"mask_api_gin/src/modules/job/model"
	"mask_api_gin/src/modules/job/processors/testerror"
	"mask_api_gin/src/modules/job/processors/testjob"
	"mask_api_gin/src/modules/job/processors/testprogress"
	"mask_api_gin/src/modules/job/service"
)

// Setup 模块注册
func Setup(router *gin.Engine) {
	logger.Infof("开始加载 ====> job 模块路由")

	// 启动时加载
	initLoad()

	// 调度任务信息
	sysJobGroup := router.Group("/job")
	sysJob := controller.NewSysJob
	{
		sysJobGroup.GET("/list",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"job:list"}}),
			sysJob.List,
		)
		sysJobGroup.GET("/:id",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"job:query"}}),
			sysJob.Info,
		)
		sysJobGroup.POST("",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"job:add"}}),
			middleware.OperateLog(middleware.OptionNew("调度任务信息", middleware.OPERATE_TYPE_INSERT)),
			sysJob.Add,
		)
		sysJobGroup.PUT("",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"job:edit"}}),
			middleware.OperateLog(middleware.OptionNew("调度任务信息", middleware.OPERATE_TYPE_UPDATE)),
			sysJob.Edit,
		)
		sysJobGroup.DELETE("/:id",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"job:remove"}}),
			middleware.OperateLog(middleware.OptionNew("调度任务信息", middleware.OPERATE_TYPE_DELETE)),
			sysJob.Remove,
		)
		sysJobGroup.PUT("/status",
			middleware.RepeatSubmit(5),
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"job:status"}}),
			middleware.OperateLog(middleware.OptionNew("调度任务信息", middleware.OPERATE_TYPE_UPDATE)),
			sysJob.Status,
		)
		sysJobGroup.PUT("/run/:id",
			middleware.RepeatSubmit(10),
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"job:status"}}),
			middleware.OperateLog(middleware.OptionNew("调度任务信息", middleware.OPERATE_TYPE_UPDATE)),
			sysJob.Run,
		)
		sysJobGroup.PUT("/reset",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"job:status"}}),
			middleware.OperateLog(middleware.OptionNew("调度任务信息", middleware.OPERATE_TYPE_OTHER)),
			sysJob.Reset,
		)
		sysJobGroup.GET("/export",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"job:export"}}),
			middleware.OperateLog(middleware.OptionNew("调度任务信息", middleware.OPERATE_TYPE_EXPORT)),
			sysJob.Export,
		)
	}

	// 调度任务日志信息
	sysJobLogGroup := router.Group("/job/log")
	sysJobLog := controller.NewSysJobLog
	{
		sysJobLogGroup.GET("/list",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"job:list"}}),
			sysJobLog.List,
		)
		sysJobLogGroup.GET("/:id",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"job:query"}}),
			sysJobLog.Info,
		)
		sysJobLogGroup.DELETE("/:id",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"job:remove"}}),
			middleware.OperateLog(middleware.OptionNew("调度任务日志信息", middleware.OPERATE_TYPE_DELETE)),
			sysJobLog.Remove,
		)
		sysJobLogGroup.DELETE("/clean",
			middleware.RepeatSubmit(5),
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"job:remove"}}),
			middleware.OperateLog(middleware.OptionNew("调度任务日志信息", middleware.OPERATE_TYPE_CLEAN)),
			sysJobLog.Clean,
		)
		sysJobLogGroup.GET("/export",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"job:export"}}),
			middleware.OperateLog(middleware.OptionNew("调度任务日志信息", middleware.OPERATE_TYPE_EXPORT)),
			sysJobLog.Export,
		)
	}
}

// initLoad 初始加载
func initLoad() {
	// 启动时，初始化调度任务
	service.NewSysJob.Reset()

	// 初始化定时任务处理
	cron.CreateQueue("test_job", testjob.NewProcessor)
	cron.CreateQueue("test_progress", testprogress.NewProcessor)
	cron.CreateQueue("test_error", testerror.NewProcessor)
}

// RegisterOnJobLog 注入任务执行日志持久化钩子。
//
// 将引擎发出的执行事件映射为调度任务日志并落库；引擎本身不依赖任何业务模型。
// 需在调度任务初始化（Reset）之前调用。
func RegisterOnJobLog() {
	cron.OnJobLog = func(record cron.JobLogRecord) {
		// 还原业务负载
		sysJob, ok := record.Data.(model.SysJob)
		if !ok {
			return
		}

		// 任务配置为不记录日志则跳过
		if sysJob.JobSaveLog == "" || sysJob.JobSaveLog == constants.STATUS_NO {
			return
		}

		// 结果信息 key 的名称
		resultName := "failed"
		if record.Status == constants.STATUS_YES {
			resultName = "completed"
		}

		// 结果信息序列化字符串
		jsonByte, _ := json.Marshal(map[string]any{
			"cron":          sysJob.JobCron,
			"invoke_target": sysJob.JobInvokeTarget,
			"result":        resultName,
			"message":       record.Result,
		})
		jobMsg := string(jsonByte)
		if len(jobMsg) > 500 {
			jobMsg = jobMsg[:500]
		}

		// 计算执行耗时并写入日志
		duration := time.Since(time.UnixMilli(record.Timestamp))
		service.NewSysJobLog.Insert(model.SysJobLog{
			JobName:         sysJob.JobName,
			JobGroup:        sysJob.JobGroup,
			JobInvokeTarget: sysJob.JobInvokeTarget,
			JobTargetParams: sysJob.JobTargetParams,
			LogStatus:       record.Status,
			LogMsg:          jobMsg,
			LogCostTime:     duration.Milliseconds(),
			CreateTime:      time.Now().UnixMilli(),
		})
	}
}
