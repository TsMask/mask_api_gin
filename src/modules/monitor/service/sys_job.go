package service

import (
	"fmt"
	"strconv"

	"mask_api_gin/src/framework/constants"
	"mask_api_gin/src/framework/cron"
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
func (s SysJob) FindById(jobId int64) model.SysJob {
	if jobId <= 0 {
		return model.SysJob{}
	}
	jobs := s.sysJobRepository.SelectByIds([]int64{jobId})
	if len(jobs) > 0 {
		return jobs[0]
	}
	return model.SysJob{}
}

// Insert 新增
func (s SysJob) Insert(sysJob model.SysJob) int64 {
	insertId := s.sysJobRepository.Insert(sysJob)
	if insertId > 0 && sysJob.StatusFlag == constants.STATUS_YES {
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
func (s SysJob) DeleteByIds(jobIds []int64) (int64, error) {
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
func (s SysJob) CheckUniqueByJobName(jobName, jobGroup string, jobId int64) bool {
	uniqueId := s.sysJobRepository.CheckUnique(model.SysJob{
		JobName:  jobName,
		JobGroup: jobGroup,
	})
	if uniqueId == jobId {
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
	queue := cron.GetQueue(sysJob.InvokeTarget)
	if queue == nil || queue.Name != sysJob.InvokeTarget {
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
	if queue == nil || queue.Name != sysJob.InvokeTarget {
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

// ExportData 导出数据表格
func (s SysJob) ExportData(rows []model.SysJob, fileName string) (string, error) {
	// 第一行表头标题
	headerCells := map[string]string{
		"A1": "任务编号",
		"B1": "任务名称",
		"C1": "任务组名",
		"D1": "调用目标",
		"E1": "传入参数",
		"F1": "执行表达式",
		"G1": "出错策略",
		"H1": "并发执行",
		"I1": "任务状态",
		"J1": "备注说明",
	}
	// 读取任务组名字典数据
	dictSysJobGroup := s.sysDictTypeService.FindDataByType("sys_job_group")
	// 从第二行开始的数据
	dataCells := make([]map[string]any, 0)
	for i, row := range rows {
		idx := strconv.Itoa(i + 2)
		// 任务组名
		sysJobGroup := ""
		for _, v := range dictSysJobGroup {
			if row.JobGroup == v.DataValue {
				sysJobGroup = v.DataLabel
				break
			}
		}
		misfirePolicy := "放弃执行"
		if row.MisfirePolicy == "1" {
			misfirePolicy = "立即执行"
		} else if row.MisfirePolicy == "2" {
			misfirePolicy = "执行一次"
		}
		concurrent := "禁止"
		if row.Concurrent == "1" {
			concurrent = "允许"
		}
		// 状态
		statusValue := "失败"
		if row.StatusFlag == constants.STATUS_YES {
			statusValue = "成功"
		}
		dataCells = append(dataCells, map[string]any{
			"A" + idx: row.JobId,
			"B" + idx: row.JobName,
			"C" + idx: sysJobGroup,
			"D" + idx: row.InvokeTarget,
			"E" + idx: row.TargetParams,
			"F" + idx: row.CronExpression,
			"G" + idx: misfirePolicy,
			"H" + idx: concurrent,
			"I" + idx: statusValue,
			"J" + idx: row.Remark,
		})
	}

	// 导出数据表格
	return file.WriteSheet(headerCells, dataCells, fileName, "")
}
