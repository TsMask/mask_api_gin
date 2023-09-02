package cron

import (
	"mask_api_gin/src/framework/logger"
	"testing"
	"time"
)

// 参考文章：
// https://blog.csdn.net/zjbyough/article/details/113853582
// https://mp.weixin.qq.com/s/Ak7RBv1NuS-VBeDNo8_fww
func init() {
	StartCron()
}

// 简单示例 队列任务处理
var NewSimple = &Simple{}

type Simple struct{}

func (s *Simple) Execute(data any) any {
	logger.Infof("执行=> %+v ", data)
	// 实现任务处理逻辑
	return data
}

func TestSimple(t *testing.T) {

	simple := CreateQueue("simple", NewSimple)
	simple.RunJob(map[string]string{
		"ok":   "ok",
		"data": "data",
	}, JobOptions{
		JobId: "101",
	})

	simpleC := CreateQueue("simple", NewSimple)
	simpleC.RunJob(map[string]string{
		"corn": "*/5 * * * * *",
		"id":   "102",
	}, JobOptions{
		JobId: "102",
		Cron:  "*/5 * * * * *",
	})

	// simpleC.RunJob(map[string]string{
	// 	"corn": "*/15 * * * * *",
	// 	"id":   "103",
	// }, JobOptions{
	// 	JobId: "103",
	// 	Cron:  "*/15 * * * * *",
	// })

	// simpleC.RemoveJob("102")

	select {}
}

// Foo 队列任务处理
var NewFooProcessor = &FooProcessor{
	progress: 0,
	count:    0,
}

type FooProcessor struct {
	progress int
	count    int
}

func (s *FooProcessor) Execute(data any) any {
	logger.Infof("执行 %d %d => %+v ", s.count, s.progress, data)
	s.count++

	// 实现任务处理逻辑
	i := 0
	s.progress = i
	for i < 10 {
		// 获取任务进度
		progress := s.progress
		logger.Infof("data: %v => 任务进度：%d", data, progress)
		// 延迟响应
		time.Sleep(time.Second * 2)
		i++
		// 改变任务进度
		s.progress = i
	}
	return data
}

func TestFoo(t *testing.T) {

	foo := CreateQueue("foo", NewFooProcessor)
	foo.RunJob(map[string]string{
		"data": "2",
	}, JobOptions{
		JobId: "2",
	})

	fooC := CreateQueue("foo", NewFooProcessor)
	fooC.RunJob(map[string]string{
		"corn": "*/5 * * * * *",
	}, JobOptions{
		JobId: "3",
		Cron:  "*/5 * * * * *",
	})

	select {}
}

// Bar 队列任务处理
var NewBarProcessor = &BarProcessor{
	progress: 0,
	count:    0,
}

type BarProcessor struct {
	progress int
	count    int
}

func (s *BarProcessor) Execute(data any) any {
	logger.Infof("执行 %d %d => %+v ", s.count, s.progress, data)
	s.count++

	// 实现任务处理逻辑
	i := 0
	s.progress = i
	for i < 5 {
		// 获取任务进度
		progress := s.progress
		logger.Infof("data: %v => 任务进度：%d", data, progress)
		// 延迟响应
		time.Sleep(time.Second * 2)
		// 程序中途执行错误
		if i == 3 {
			// arr := [1]int{1}
			// arr[i] = 3
			// fmt.Println(arr)
			// return "i = 3"
			panic("程序中途执行错误")
		}
		i++
		// 改变任务进度
		s.progress = i
	}

	return data
}

func TestBar(t *testing.T) {

	bar := CreateQueue("bar", NewBarProcessor)
	bar.RunJob(map[string]string{
		"data": "wdf",
	}, JobOptions{
		JobId: "81923",
	})

	barC := CreateQueue("bar", NewBarProcessor)
	barC.RunJob(map[string]string{
		"corn": "*/5 * * * * *",
	}, JobOptions{
		JobId: "789",
		Cron:  "*/5 * * * * *",
	})

	// barDB := CreateQueue("barDB", NewBarProcessor)
	// barDB.RunJob(JobData{
	// 	SysJob: model.SysJob{
	// 		JobID:   "9123",
	// 		JobName: "测试任务",
	// 	},
	// }, JobOptions{
	// 	JobId: "9123",
	// })

	select {}
}
