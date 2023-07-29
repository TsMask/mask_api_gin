package processor

import (
	"mask_api_gin/src/framework/cron"
	"mask_api_gin/src/modules/monitor/processor/bar"
	"mask_api_gin/src/modules/monitor/processor/foo"
	"mask_api_gin/src/modules/monitor/processor/simple"
)

// InitCronQueue 初始定时任务队列
func InitCronQueue() {
	cron.AddQueue("simple", simple.Execute)
	cron.AddQueue("foo", foo.Execute)
	cron.AddQueue("bar", bar.Execute)
}
