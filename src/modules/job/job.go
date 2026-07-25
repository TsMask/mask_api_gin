package job

import (
	"mask_api_gin/src/internal/infra/cron"
	"mask_api_gin/src/modules/job/testerror"
	"mask_api_gin/src/modules/job/testjob"
	"mask_api_gin/src/modules/job/testprogress"
)

// Setup 模块注册
func Setup() {
	// 初始化定时任务处理
	cron.CreateQueue("test_job", testjob.NewProcessor)
	cron.CreateQueue("test_progress", testprogress.NewProcessor)
	cron.CreateQueue("test_error", testerror.NewProcessor)
}
