package service

import (
	"strconv"

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
	return s.sysJobLogRepository.SelectByPage(query)
}

// Find 查询
func (s SysJobLog) Find(sysJobLog model.SysJobLog) []model.SysJobLog {
	return s.sysJobLogRepository.Select(sysJobLog)
}

// FindById 通过ID查询
func (s SysJobLog) FindById(logId int64) model.SysJobLog {
	return s.sysJobLogRepository.SelectById(logId)
}

// RemoveByIds 批量删除
func (s SysJobLog) RemoveByIds(logIds []int64) int64 {
	return s.sysJobLogRepository.DeleteByIds(logIds)
}

// Clean 清空
func (s SysJobLog) Clean() int64 {
	return s.sysJobLogRepository.Clean()
}

// ExportData 导出数据表格
func (s SysJobLog) ExportData(rows []model.SysJobLog, fileName string) (string, error) {
	// 第一行表头标题
	headerCells := map[string]string{
		"A1": "日志序号",
		"B1": "任务名称",
		"C1": "任务组名",
		"D1": "调用目标",
		"E1": "传入参数",
		"F1": "日志信息",
		"G1": "执行状态",
		"H1": "记录时间",
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
		// 状态
		statusValue := "失败"
		if row.StatusFlag == "1" {
			statusValue = "成功"
		}
		dataCells = append(dataCells, map[string]any{
			"A" + idx: row.LogId,
			"B" + idx: row.JobName,
			"C" + idx: sysJobGroup,
			"D" + idx: row.InvokeTarget,
			"E" + idx: row.TargetParams,
			"F" + idx: row.JobMsg,
			"G" + idx: statusValue,
			"H" + idx: date.ParseDateToStr(row.CreateTime, date.YYYY_MM_DD_HH_MM_SS),
		})
	}

	// 导出数据表格
	return file.WriteSheet(headerCells, dataCells, fileName, "")
}
