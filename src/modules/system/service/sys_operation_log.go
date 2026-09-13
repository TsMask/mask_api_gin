package service

import (
	"fmt"
	"path/filepath"
	"time"

	"mask_api_gin/src/internal/shared/constants"
	"mask_api_gin/src/internal/shared/date"
	"mask_api_gin/src/internal/shared/file"
	"mask_api_gin/src/internal/shared/logger"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/repository"
)

// NewSysOperationLog 实例化服务层
var NewSysOperationLog = &SysOperationLog{
	sysOperationLogRepository: repository.NewSysOperationLog,
}

// SysOperationLog 操作日志表 服务层处理
type SysOperationLog struct {
	sysOperationLogRepository *repository.SysOperationLog // 操作日志信息
}

// FindByPage 分页查询列表数据
func (s SysOperationLog) FindByPage(query map[string]string, dataScopeSQL string) ([]model.SysOperationLog, int64) {
	rows, total, err := s.sysOperationLogRepository.SelectByPage(query, dataScopeSQL)
	if err != nil {
		logger.Errorf("query err => %v", err.Error())
	}
	return rows, total
}

// Insert 新增信息
func (s SysOperationLog) Insert(sysOperationLog model.SysOperationLog) int64 {
	if sysOperationLog.CreateBy != "" {
		sysOperationLog.CreateTime = time.Now().UnixMilli()
	}
	rows, err := s.sysOperationLogRepository.Insert(sysOperationLog)
	if err != nil {
		logger.Errorf("insert err => %v", err.Error())
		return 0
	}
	return rows
}

// Clean 清空操作日志
func (s SysOperationLog) Clean() int64 {
	rows, err := s.sysOperationLogRepository.Clean()
	if err != nil {
		logger.Errorf("clean err => %v", err.Error())
		return 0
	}
	return rows
}

// ExportData 导出数据表格
func (s SysOperationLog) ExportData(rows []model.SysOperationLog, fileName string) (string, error) {
	// 第一行表头标题
	data := [][]string{
		{"操作序号", "模块标题", "业务类型", "请求URL", "请求方式", "主机地址", "操作地点", "请求参数", "操作消息", "方法名称", "操作人员", "操作时间", "操作状态", "消耗时间"},
	}
	// 从第二行开始的数据
	for _, row := range rows {
		// 业务类型
		operationType := ""
		// 状态
		statusValue := "失败"
		if row.OperationStatus == constants.STATUS_YES {
			statusValue = "成功"
		}
		data = append(data, []string{
			fmt.Sprintf("%d", row.ID),
			row.Title,
			operationType,
			row.OperationUrl,
			row.OperationMethod,
			row.OperationIp,
			row.OperationLocation,
			row.OperationParam,
			row.OperationMsg,
			row.OperationFunction,
			row.CreateBy,
			date.ParseNumberToDate(row.CreateTime).Format(time.RFC3339),
			statusValue,
			fmt.Sprintf("%d", row.OperationCostTime),
		})
	}

	// 导出文件路径
	_, dir := file.UploadResourceDir()
	datePath := date.ParseDatePath()
	filePath := filepath.Join(dir, constants.UPLOAD_EXPORT, datePath, fileName)

	// 导出数据表格
	return filePath, file.CSVWrite(filePath, data, false)
}
