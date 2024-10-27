package service

import (
	"fmt"
	constSystem "mask_api_gin/src/framework/constants/system"
	"mask_api_gin/src/framework/cron"
	"mask_api_gin/src/modules/monitor/model"
	"mask_api_gin/src/modules/monitor/repository"
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
func (r *SysJob) FindByPage(query map[string]any) ([]model.SysJob, int64) {
	return r.sysJobRepository.SelectByPage(query)
}

// Find 查询
func (r *SysJob) Find(sysJob model.SysJob) []model.SysJob {
	return r.sysJobRepository.Select(sysJob)
}

// FindById 通过ID查询
func (r *SysJob) FindById(jobId string) model.SysJob {
	if jobId == "" {
		return model.SysJob{}
	}
	if jobs := r.sysJobRepository.SelectByIds([]string{jobId}); len(jobs) > 0 {
		return jobs[0]
	}
	return model.SysJob{}
}

// Insert 新增调度任务信息
func (r *SysJob) Insert(sysJob model.SysJob) string {
	insertId := r.sysJobRepository.Insert(sysJob)
	if insertId == "" && sysJob.Status == constSystem.STATUS_YES {
		sysJob.JobID = insertId
		r.insertQueueJob(sysJob, true)
	}
	return insertId
}

// Update 修改
func (r *SysJob) Update(sysJob model.SysJob) int64 {
	rows := r.sysJobRepository.Update(sysJob)
	if rows > 0 {
		//状态正常添加队列任务
		if sysJob.Status == constSystem.STATUS_YES {
			r.insertQueueJob(sysJob, true)
		}
		// 状态禁用删除队列任务
		if sysJob.Status == constSystem.STATUS_NO {
			r.deleteQueueJob(sysJob)
		}
	}
	return rows
}

// DeleteByIds 批量删除
func (r *SysJob) DeleteByIds(jobIds []string) (int64, error) {
	// 检查是否存在
	jobs := r.sysJobRepository.SelectByIds(jobIds)
	if len(jobs) <= 0 {
		return 0, fmt.Errorf("没有权限访问调度任务数据！")
	}
	if len(jobs) == len(jobIds) {
		// 清除任务
		for _, job := range jobs {
			r.deleteQueueJob(job)
		}
		return r.sysJobRepository.DeleteByIds(jobIds), nil
	}
	return 0, fmt.Errorf("删除调度任务信息失败！")
}

// CheckUniqueByJobName 校验调度任务名称和组是否唯一
func (r *SysJob) CheckUniqueByJobName(jobName, jobGroup, jobId string) bool {
	uniqueId := r.sysJobRepository.CheckUniqueJob(model.SysJob{
		JobName:  jobName,
		JobGroup: jobGroup,
	})
	if uniqueId == jobId {
		return true
	}
	return uniqueId == ""
}

// Run 立即运行一次调度任务
func (r *SysJob) Run(sysJob model.SysJob) bool {
	return r.insertQueueJob(sysJob, false)
}

// insertQueueJob 添加调度任务
func (r *SysJob) insertQueueJob(sysJob model.SysJob, repeat bool) bool {
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
			JobId: sysJob.JobID,
		})
		// 执行中或等待中的都返回正常
		return status == cron.Active || status == cron.Waiting
	}

	// 执行重复任务
	queue.RunJob(options, cron.JobOptions{
		JobId: sysJob.JobID,
		Cron:  sysJob.CronExpression,
	})
	return true
}

// deleteQueueJob 删除调度任务
func (r *SysJob) deleteQueueJob(sysJob model.SysJob) bool {
	// 获取队列 Processor
	queue := cron.GetQueue(sysJob.InvokeTarget)
	if queue.Name != sysJob.InvokeTarget {
		return false
	}
	return queue.RemoveJob(sysJob.JobID)
}

// Reset 重置初始调度任务
func (r *SysJob) Reset() {
	// 获取注册的队列名称
	queueNames := cron.QueueNames()
	if len(queueNames) == 0 {
		return
	}
	// 查询系统中定义状态为正常启用的任务
	sysJobs := r.sysJobRepository.Select(model.SysJob{
		Status: constSystem.STATUS_YES,
	})
	for _, sysJob := range sysJobs {
		for _, name := range queueNames {
			if name == sysJob.InvokeTarget {
				r.insertQueueJob(sysJob, true)
			}
		}
	}
}
