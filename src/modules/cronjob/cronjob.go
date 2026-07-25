package cronjob

import (
	"mask_api_gin/src/internal/infra/cron"
	"mask_api_gin/src/modules/cronjob/testerror"
	"mask_api_gin/src/modules/cronjob/testjob"
	"mask_api_gin/src/modules/cronjob/testprogress"
)

// Setup 模块注册
func Setup() {
	// 初始化定时任务处理
	cron.CreateQueue("test_job", testjob.NewProcessor)
	cron.CreateQueue("test_progress", testprogress.NewProcessor)
	cron.CreateQueue("test_error", testerror.NewProcessor)
}
