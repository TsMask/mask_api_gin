package service

import (
	"fmt"
	"path/filepath"
	"time"

	"mask_api_gin/src/framework/constants"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/utils/date"
	"mask_api_gin/src/framework/utils/file"
	"mask_api_gin/src/modules/monitor/model"
	"mask_api_gin/src/modules/monitor/repository"
	systemService "mask_api_gin/src/modules/system/service"
)

// NewSysJobLog 服务层实例化
var NewSysJobLog = &SysJobLog{
	sysJobLogRepository: repository.NewSysJobLog,
	sysDictTypeService:  systemService.NewSysDictType,
}

// SysJobLog 调度任务日志 服务层处理
type SysJobLog struct {
	sysJobLogRepository *repository.SysJobLog      // 调度任务日志数据信息
	sysDictTypeService  *systemService.SysDictType // 字典类型服务
}

// FindByPage 分页查询
func (s SysJobLog) FindByPage(query map[string]string) ([]model.SysJobLog, int64) {
	rows, total, err := s.sysJobLogRepository.SelectByPage(query)
	if err != nil {
		logger.Errorf("query err => %v", err.Error())
	}
	return rows, total
}

// Find 查询
func (s SysJobLog) Find(sysJobLog model.SysJobLog) []model.SysJobLog {
	rows, err := s.sysJobLogRepository.Select(sysJobLog)
	if err != nil {
		logger.Errorf("query err => %v", err.Error())
	}
	return rows
}

// FindById 通过ID查询
func (s SysJobLog) FindById(id int64) model.SysJobLog {
	item, err := s.sysJobLogRepository.SelectById(id)
	if err != nil {
		logger.Errorf("query err => %v", err.Error())
	}
	return item
}

// DeleteByIds 批量删除
func (s SysJobLog) DeleteByIds(ids []int64) int64 {
	rows, err := s.sysJobLogRepository.DeleteByIds(ids)
	if err != nil {
		logger.Errorf("delete err => %v", err.Error())
		return 0
	}
	return rows
}

// Clean 清空
func (s SysJobLog) Clean() int64 {
	rows, err := s.sysJobLogRepository.Clean()
	if err != nil {
		logger.Errorf("clean err => %v", err.Error())
		return 0
	}
	return rows
}

// ExportData 导出数据表格
func (s SysJobLog) ExportData(rows []model.SysJobLog, fileName string) (string, error) {
	// 第一行表头标题
	data := [][]string{
		{"日志序号", "任务名称", "任务组名", "调用目标", "传入参数", "日志信息", "执行状态", "记录时间"},
	}
	// 读取任务组名字典数据
	dictSysJobGroup := s.sysDictTypeService.FindDataByType("sys_job_group")
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
		// 状态
		statusValue := "失败"
		if row.LogStatus == constants.STATUS_YES {
			statusValue = "成功"
		}
		data = append(data, []string{
			fmt.Sprintf("%d", row.ID),
			row.JobName,
			sysJobGroup,
			row.JobInvokeTarget,
			row.JobTargetParams,
			row.LogMsg,
			statusValue,
			date.ParseNumberToDate(row.CreateTime).Format(time.RFC3339),
		})
	}

	// 导出文件路径
	_, dir := file.UploadResourceDir()
	datePath := date.ParseDatePath()
	filePath := filepath.Join(dir, constants.UPLOAD_EXPORT, datePath, fileName)

	// 导出数据表格
	return filePath, file.CSVWrite(filePath, data, false)
}
