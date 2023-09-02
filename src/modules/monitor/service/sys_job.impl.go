package service

import (
	"errors"
	"mask_api_gin/src/framework/constants/common"
	"mask_api_gin/src/framework/cron"
	"mask_api_gin/src/modules/monitor/model"
	"mask_api_gin/src/modules/monitor/repository"
)

// 实例化服务层 SysJobImpl 结构体
var NewSysJobImpl = &SysJobImpl{
	sysJobRepository: repository.NewSysJobImpl,
}

// SysJobImpl 调度任务 服务层处理
type SysJobImpl struct {
	// 调度任务数据信息
	sysJobRepository repository.ISysJob
}

// SelectJobPage 分页查询调度任务集合
func (r *SysJobImpl) SelectJobPage(query map[string]any) map[string]any {
	return r.sysJobRepository.SelectJobPage(query)
}

// SelectJobList 查询调度任务集合
func (r *SysJobImpl) SelectJobList(sysJob model.SysJob) []model.SysJob {
	return r.sysJobRepository.SelectJobList(sysJob)
}

// SelectJobById 通过调度ID查询调度任务信息
func (r *SysJobImpl) SelectJobById(jobId string) model.SysJob {
	if jobId == "" {
		return model.SysJob{}
	}
	jobs := r.sysJobRepository.SelectJobByIds([]string{jobId})
	if len(jobs) > 0 {
		return jobs[0]
	}
	return model.SysJob{}
}

// CheckUniqueJobName 校验调度任务名称和组是否唯一
func (r *SysJobImpl) CheckUniqueJobName(jobName, jobGroup, jobId string) bool {
	uniqueId := r.sysJobRepository.CheckUniqueJob(model.SysJob{
		JobName:  jobName,
		JobGroup: jobGroup,
	})
	if uniqueId == jobId {
		return true
	}
	return uniqueId == ""
}

// InsertJob 新增调度任务信息
func (r *SysJobImpl) InsertJob(sysJob model.SysJob) string {
	insertId := r.sysJobRepository.InsertJob(sysJob)
	if insertId == "" && sysJob.Status == common.STATUS_YES {
		sysJob.JobID = insertId
		r.insertQueueJob(sysJob, true)
	}
	return insertId
}

// UpdateJob 修改调度任务信息
func (r *SysJobImpl) UpdateJob(sysJob model.SysJob) int64 {
	rows := r.sysJobRepository.UpdateJob(sysJob)
	if rows > 0 {
		//状态正常添加队列任务
		if sysJob.Status == common.STATUS_YES {
			r.insertQueueJob(sysJob, true)
		}
		// 状态禁用删除队列任务
		if sysJob.Status == common.STATUS_NO {
			r.deleteQueueJob(sysJob)
		}
	}
	return rows
}

// DeleteJobByIds 批量删除调度任务信息
func (r *SysJobImpl) DeleteJobByIds(jobIds []string) (int64, error) {
	// 检查是否存在
	jobs := r.sysJobRepository.SelectJobByIds(jobIds)
	if len(jobs) <= 0 {
		return 0, errors.New("没有权限访问调度任务数据！")
	}
	if len(jobs) == len(jobIds) {
		// 清除任务
		for _, job := range jobs {
			r.deleteQueueJob(job)
		}
		rows := r.sysJobRepository.DeleteJobByIds(jobIds)
		return rows, nil
	}
	return 0, errors.New("删除调度任务信息失败！")
}

// ChangeStatus 任务调度状态修改
func (r *SysJobImpl) ChangeStatus(sysJob model.SysJob) bool {
	// 更新状态
	newSysJob := model.SysJob{
		JobID:    sysJob.JobID,
		Status:   sysJob.Status,
		UpdateBy: sysJob.UpdateBy,
	}
	rows := r.sysJobRepository.UpdateJob(newSysJob)
	if rows > 0 {
		//状态正常添加队列任务
		if sysJob.Status == common.STATUS_YES {
			r.insertQueueJob(sysJob, true)
		}
		// 状态禁用删除队列任务
		if sysJob.Status == common.STATUS_NO {
			r.deleteQueueJob(sysJob)
		}
		return true
	}
	return false
}

// ResetQueueJob 重置初始调度任务
func (r *SysJobImpl) ResetQueueJob() {
	// 获取注册的队列名称
	queueNames := cron.QueueNames()
	if len(queueNames) == 0 {
		return
	}
	// 查询系统中定义状态为正常启用的任务
	sysJobs := r.sysJobRepository.SelectJobList(model.SysJob{
		Status: common.STATUS_YES,
	})
	for _, sysJob := range sysJobs {
		for _, name := range queueNames {
			if name == sysJob.InvokeTarget {
				r.insertQueueJob(sysJob, true)
			}
		}
	}
}

// RunQueueJob 立即运行一次调度任务
func (r *SysJobImpl) RunQueueJob(sysJob model.SysJob) bool {
	return r.insertQueueJob(sysJob, false)
}

// insertQueueJob 添加调度任务
func (r *SysJobImpl) insertQueueJob(sysJob model.SysJob, repeat bool) bool {
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
func (r *SysJobImpl) deleteQueueJob(sysJob model.SysJob) bool {
	// 获取队列 Processor
	queue := cron.GetQueue(sysJob.InvokeTarget)
	if queue.Name != sysJob.InvokeTarget {
		return false
	}
	return queue.RemoveJob(sysJob.JobID)
}
