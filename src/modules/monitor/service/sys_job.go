package service

import (
	"fmt"
	"path/filepath"
	"time"

	"mask_api_gin/src/framework/constants"
	"mask_api_gin/src/framework/cron"
	"mask_api_gin/src/framework/utils/date"
	"mask_api_gin/src/framework/utils/file"
	"mask_api_gin/src/modules/monitor/model"
	"mask_api_gin/src/modules/monitor/repository"
	systemService "mask_api_gin/src/modules/system/service"
)

// NewSysJob  服务层实例化
var NewSysJob = &SysJob{
	sysJobRepository:   repository.NewSysJob,
	sysDictTypeService: systemService.NewSysDictType,
}

// SysJob 调度任务 服务层处理
type SysJob struct {
	sysJobRepository   *repository.SysJob         // 调度任务数据信息
	sysDictTypeService *systemService.SysDictType // 字典类型服务
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
func (s SysJob) FindById(id int64) model.SysJob {
	if id <= 0 {
		return model.SysJob{}
	}
	jobs := s.sysJobRepository.SelectByIds([]int64{id})
	if len(jobs) > 0 {
		return jobs[0]
	}
	return model.SysJob{}
}

// Insert 新增
func (s SysJob) Insert(sysJob model.SysJob) int64 {
	if sysJob.CreateBy != "" {
		ms := time.Now().UnixMilli()
		sysJob.UpdateBy = sysJob.CreateBy
		sysJob.UpdateTime = ms
		sysJob.CreateTime = ms
	}
	insertId := s.sysJobRepository.Insert(sysJob)
	if insertId > 0 && sysJob.JobStatus == constants.STATUS_YES {
		sysJob.ID = insertId
		s.insertQueueJob(sysJob, true)
	}
	return insertId
}

// Update 修改
func (s SysJob) Update(sysJob model.SysJob) int64 {
	if sysJob.UpdateBy != "" {
		sysJob.UpdateTime = time.Now().UnixMilli()
	}
	rows := s.sysJobRepository.Update(sysJob)
	if rows > 0 {
		//状态正常添加队列任务
		if sysJob.JobStatus == constants.STATUS_YES {
			s.insertQueueJob(sysJob, true)
		}
		// 状态禁用删除队列任务
		if sysJob.JobStatus == constants.STATUS_NO {
			s.deleteQueueJob(sysJob)
		}
	}
	return rows
}

// DeleteByIds 批量删除
func (s SysJob) DeleteByIds(ids []int64) (int64, error) {
	// 检查是否存在
	jobs := s.sysJobRepository.SelectByIds(ids)
	if len(jobs) <= 0 {
		return 0, fmt.Errorf("没有权限访问调度任务数据！")
	}
	if len(jobs) == len(ids) {
		// 清除任务
		for _, job := range jobs {
			s.deleteQueueJob(job)
		}
		return s.sysJobRepository.DeleteByIds(ids), nil
	}
	return 0, fmt.Errorf("删除调度任务信息失败！")
}

// CheckUniqueByJobName 校验调度任务名称和组是否唯一
func (s SysJob) CheckUniqueByJobName(jobName, jobGroup string, id int64) bool {
	uniqueId := s.sysJobRepository.CheckUnique(model.SysJob{
		JobName:  jobName,
		JobGroup: jobGroup,
	})
	if uniqueId == id {
		return true
	}
	return uniqueId == 0
}

// Run 立即运行一次调度任务
func (s SysJob) Run(sysJob model.SysJob) bool {
	return s.insertQueueJob(sysJob, false)
}

// insertQueueJob 添加调度任务
func (s SysJob) insertQueueJob(sysJob model.SysJob, repeat bool) bool {
	// 获取队列 Processor
	queue := cron.GetQueue(sysJob.JobInvokeTarget)
	if queue == nil || queue.Name != sysJob.JobInvokeTarget {
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
			ID: fmt.Sprintf("%d", sysJob.ID),
		})
		// 执行中或等待中的都返回正常
		return status == cron.Active || status == cron.Waiting
	}

	// 执行重复任务
	queue.RunJob(options, cron.JobOptions{
		ID:   fmt.Sprintf("%d", sysJob.ID),
		Cron: sysJob.JobCron,
	})
	return true
}

// deleteQueueJob 删除调度任务
func (s SysJob) deleteQueueJob(sysJob model.SysJob) bool {
	// 获取队列 Processor
	queue := cron.GetQueue(sysJob.JobInvokeTarget)
	if queue == nil || queue.Name != sysJob.JobInvokeTarget {
		return false
	}
	return queue.RemoveJob(fmt.Sprintf("%d", sysJob.ID))
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
		JobStatus: constants.STATUS_YES,
	})
	for _, sysJob := range sysJobs {
		for _, name := range queueNames {
			if name == sysJob.JobInvokeTarget {
				s.deleteQueueJob(sysJob)
				s.insertQueueJob(sysJob, true)
			}
		}
	}
}

// ExportData 导出数据表格
func (s SysJob) ExportData(rows []model.SysJob, fileName string) (string, error) {
	// 读取任务组名字典数据
	dictSysJobGroup := s.sysDictTypeService.FindDataByType("sys_job_group")
	// 第一行表头标题
	data := [][]string{
		{"任务编号", "任务名称", "任务组名", "调用目标", "传入参数", "执行表达式", "出错策略", "并发执行", "任务状态", "备注说明"},
	}
	// 从第二行开始的数据
	for _, row := range rows {
		// 任务组名
		sysJobGroup := ""
		for _, v := range dictSysJobGroup {
			if row.JobGroup == v.DataValue {
				sysJobGroup = v.DataLabel
				break
			}
		}
		misfirePolicy := "放弃执行"
		switch row.JobMisfirePolicy {
		case "1":
			misfirePolicy = "立即执行"
		case "2":
			misfirePolicy = "执行一次"
		}
		concurrent := "禁止"
		if row.JobConcurrent == "1" {
			concurrent = "允许"
		}
		// 状态
		statusValue := "失败"
		if row.JobStatus == constants.STATUS_YES {
			statusValue = "成功"
		}
		data = append(data, []string{
			fmt.Sprintf("%d", row.ID),
			row.JobName,
			sysJobGroup,
			row.JobInvokeTarget,
			row.JobTargetParams,
			row.JobCron,
			misfirePolicy,
			concurrent,
			statusValue,
			row.Remark,
		})
	}

	// 导出文件路径
	_, dir := file.UploadResourceDir()
	datePath := date.ParseDatePath()
	filePath := filepath.Join(dir, constants.UPLOAD_EXPORT, datePath, fileName)

	// 导出数据表格
	return filePath, file.CSVWrite(filePath, data, false)
}
