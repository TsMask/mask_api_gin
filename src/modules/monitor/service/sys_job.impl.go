package service

import (
	"errors"
	"mask_api_gin/src/framework/constants/common"
	"mask_api_gin/src/modules/monitor/model"
	"mask_api_gin/src/modules/monitor/repository"
)

// SysJobImpl 调度任务 业务层处理
var SysJobImpl = &sysJobImpl{
	sysJobRepository: repository.SysJobImpl,
}

type sysJobImpl struct {
	// 调度任务日志信息
	sysJobRepository repository.ISysJob
}

// SelectJobPage 分页查询调度任务集合
func (r *sysJobImpl) SelectJobPage(query map[string]string) map[string]interface{} {
	return r.sysJobRepository.SelectJobPage(query)
}

// SelectJobList 查询调度任务集合
func (r *sysJobImpl) SelectJobList(sysJob model.SysJob) []model.SysJob {
	return r.sysJobRepository.SelectJobList(sysJob)
}

// SelectJobById 通过调度ID查询调度任务信息
func (r *sysJobImpl) SelectJobById(jobId string) model.SysJob {
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
func (r *sysJobImpl) CheckUniqueJobName(jobName, jobGroup, jobId string) bool {
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
func (r *sysJobImpl) InsertJob(sysJob model.SysJob) string {
	insertId := r.sysJobRepository.InsertJob(sysJob)
	if insertId == "" && sysJob.Status == common.STATUS_YES {
		sysJob.JobID = insertId
		r.insertQueueJob(sysJob, true)
	}
	return insertId
}

// UpdateJob 修改调度任务信息
func (r *sysJobImpl) UpdateJob(sysJob model.SysJob) int64 {
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
func (r *sysJobImpl) DeleteJobByIds(jobIds []string) (int64, error) {
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
func (r *sysJobImpl) ChangeStatus(sysJob model.SysJob) bool {
	//状态正常添加队列任务
	if sysJob.Status == common.STATUS_YES {
		r.insertQueueJob(sysJob, true)
	}
	// 状态禁用删除队列任务
	if sysJob.Status == common.STATUS_NO {
		r.deleteQueueJob(sysJob)
	}
	// 更新状态
	newSysJob := model.SysJob{
		JobID:    sysJob.JobID,
		Status:   sysJob.Status,
		UpdateBy: sysJob.UpdateBy,
	}
	rows := r.sysJobRepository.UpdateJob(newSysJob)
	return rows > 0
}

// ResetQueueJob 重置初始调度任务
func (r *sysJobImpl) ResetQueueJob() {
}

// RunQueueJob 立即运行一次调度任务
func (r *sysJobImpl) RunQueueJob(sysJob model.SysJob) bool {
	return false
}

// insertQueueJob 添加调度任务
func (r *sysJobImpl) insertQueueJob(sysJob model.SysJob, repeat bool) bool {
	return false
}

// deleteQueueJob 删除调度任务
func (r *sysJobImpl) deleteQueueJob(sysJob model.SysJob) error {
	return nil
}
