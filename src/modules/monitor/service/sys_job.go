package service

import (
	"mask_api_gin/src/framework/constants"
	"mask_api_gin/src/framework/cron"
	"mask_api_gin/src/modules/monitor/model"
	"mask_api_gin/src/modules/monitor/repository"

	"fmt"
)

// NewSysJob  服务层实例化
var NewSysJob = &SysJob{
	sysJobRepository: repository.NewSysJob,
}

// SysJob 调度任务 服务层处理
type SysJob struct {
	sysJobRepository *repository.SysJob // 调度任务数据信息
}

// FindByPage 分页查询
func (s SysJob) FindByPage(query map[string]string) ([]model.SysJob, int64) {
	return s.sysJobRepository.SelectByPage(query)
}

// Find 查询
func (s SysJob) Find(sysJob model.SysJob) []model.SysJob {
	return s.sysJobRepository.Select(sysJob)
}

// FindById 通过ID查询
func (s SysJob) FindById(jobId string) model.SysJob {
	if jobId == "" {
		return model.SysJob{}
	}
	if jobs := s.sysJobRepository.SelectByIds([]string{jobId}); len(jobs) > 0 {
		return jobs[0]
	}
	return model.SysJob{}
}

// Insert 新增调度任务信息
func (s SysJob) Insert(sysJob model.SysJob) string {
	insertId := s.sysJobRepository.Insert(sysJob)
	if insertId != "" && sysJob.StatusFlag == constants.STATUS_YES {
		sysJob.JobId = insertId
		s.insertQueueJob(sysJob, true)
	}
	return insertId
}

// Update 修改
func (s SysJob) Update(sysJob model.SysJob) int64 {
	rows := s.sysJobRepository.Update(sysJob)
	if rows > 0 {
		//状态正常添加队列任务
		if sysJob.StatusFlag == constants.STATUS_YES {
			s.insertQueueJob(sysJob, true)
		}
		// 状态禁用删除队列任务
		if sysJob.StatusFlag == constants.STATUS_NO {
			s.deleteQueueJob(sysJob)
		}
	}
	return rows
}

// DeleteByIds 批量删除
func (s SysJob) DeleteByIds(jobIds []string) (int64, error) {
	// 检查是否存在
	jobs := s.sysJobRepository.SelectByIds(jobIds)
	if len(jobs) <= 0 {
		return 0, fmt.Errorf("没有权限访问调度任务数据！")
	}
	if len(jobs) == len(jobIds) {
		// 清除任务
		for _, job := range jobs {
			s.deleteQueueJob(job)
		}
		return s.sysJobRepository.DeleteByIds(jobIds), nil
	}
	return 0, fmt.Errorf("删除调度任务信息失败！")
}

// CheckUniqueByJobName 校验调度任务名称和组是否唯一
func (s SysJob) CheckUniqueByJobName(jobName, jobGroup, jobId string) bool {
	uniqueId := s.sysJobRepository.CheckUniqueJob(model.SysJob{
		JobName:  jobName,
		JobGroup: jobGroup,
	})
	if uniqueId == jobId {
		return true
	}
	return uniqueId == ""
}

// Run 立即运行一次调度任务
func (s SysJob) Run(sysJob model.SysJob) bool {
	return s.insertQueueJob(sysJob, false)
}

// insertQueueJob 添加调度任务
func (s SysJob) insertQueueJob(sysJob model.SysJob, repeat bool) bool {
	// 获取队列 Processor
	queue := cron.GetQueue(sysJob.InvokeTarget)
	if queue.Name != sysJob.InvokeTarget {
		return false
	}

	// 给执行任务数据参数
	options := cron.JobData{
		Repeat: repeat,
		SysJob: sysJob,
	}

	// 不是重复任务的情况，立即执行一次
	if !repeat {
		// 执行单次任务
		status := queue.RunJob(options, cron.JobOptions{
			JobId: sysJob.JobId,
		})
		// 执行中或等待中的都返回正常
		return status == cron.Active || status == cron.Waiting
	}

	// 执行重复任务
	queue.RunJob(options, cron.JobOptions{
		JobId: sysJob.JobId,
		Cron:  sysJob.CronExpression,
	})
	return true
}

// deleteQueueJob 删除调度任务
func (s SysJob) deleteQueueJob(sysJob model.SysJob) bool {
	// 获取队列 Processor
	queue := cron.GetQueue(sysJob.InvokeTarget)
	if queue.Name != sysJob.InvokeTarget {
		return false
	}
	return queue.RemoveJob(sysJob.JobId)
}

// Reset 重置初始调度任务
func (s SysJob) Reset() {
	// 获取注册的队列名称
	queueNames := cron.QueueNames()
	if len(queueNames) == 0 {
		return
	}
	// 查询系统中定义状态为正常启用的任务
	sysJobs := s.sysJobRepository.Select(model.SysJob{
		StatusFlag: constants.STATUS_YES,
	})
	for _, sysJob := range sysJobs {
		for _, name := range queueNames {
			if name == sysJob.InvokeTarget {
				s.deleteQueueJob(sysJob)
				s.insertQueueJob(sysJob, true)
			}
		}
	}
}
