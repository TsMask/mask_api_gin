package service

import (
	"fmt"
	"path/filepath"
	"time"

	"mask_api_gin/src/framework/constants"
	"mask_api_gin/src/framework/utils/date"
	"mask_api_gin/src/framework/utils/file"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/repository"
)

// NewSysLogOperate 实例化服务层
var NewSysLogOperate = &SysLogOperate{
	sysLogOperateRepository: repository.NewSysLogOperate,
}

// SysLogOperate 操作日志表 服务层处理
type SysLogOperate struct {
	sysLogOperateRepository *repository.SysLogOperate // 操作日志信息
}

// FindByPage 分页查询列表数据
func (s SysLogOperate) FindByPage(query map[string]string, dataScopeSQL string) ([]model.SysLogOperate, int64) {
	return s.sysLogOperateRepository.SelectByPage(query, dataScopeSQL)
}

// Insert 新增信息
func (s SysLogOperate) Insert(sysLogOperate model.SysLogOperate) int64 {
	if sysLogOperate.CreateBy != "" {
		sysLogOperate.CreateTime = time.Now().UnixMilli()
	}
	return s.sysLogOperateRepository.Insert(sysLogOperate)
}

// Clean 清空操作日志
func (s SysLogOperate) Clean() int64 {
	return s.sysLogOperateRepository.Clean()
}

// ExportData 导出数据表格
func (s SysLogOperate) ExportData(rows []model.SysLogOperate, fileName string) (string, error) {
	// 第一行表头标题
	data := [][]string{
		{"操作序号", "模块标题", "业务类型", "请求URL", "请求方式", "主机地址", "操作地点", "请求参数", "操作消息", "方法名称", "操作人员", "操作时间", "操作状态", "消耗时间"},
	}
	// 从第二行开始的数据
	for _, row := range rows {
		// 业务类型
		businessType := ""
		// 状态
		statusValue := "失败"
		if row.OperateStatus == constants.STATUS_YES {
			statusValue = "成功"
		}
		data = append(data, []string{
			fmt.Sprintf("%d", row.ID),
			row.Title,
			businessType,
			row.OperateUrl,
			row.OperateMethod,
			row.OperateIp,
			row.OperateLocation,
			row.OperateParam,
			row.OperateMsg,
			row.OperateFunction,
			row.CreateBy,
			date.ParseNumberToDate(row.CreateTime).Format(time.RFC3339),
			statusValue,
			fmt.Sprintf("%d", row.OperateCostTime),
		})
	}

	// 导出文件路径
	_, dir := file.UploadResourceDir()
	datePath := date.ParseDatePath()
	filePath := filepath.Join(dir, constants.UPLOAD_EXPORT, datePath, fileName)

	// 导出数据表格
	return filePath, file.CSVWrite(filePath, data, false)
}
