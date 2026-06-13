package service

import (
	"strconv"
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
	return s.sysLogOperateRepository.Insert(sysLogOperate)
}

// Clean 清空操作日志
func (s SysLogOperate) Clean() int64 {
	return s.sysLogOperateRepository.Clean()
}

// ExportData 导出数据表格
func (s SysLogOperate) ExportData(rows []model.SysLogOperate, fileName string) (string, error) {
	// 第一行表头标题
	headerCells := map[string]string{
		"A1": "操作序号",
		"B1": "模块标题",
		"C1": "业务类型",
		"D1": "请求URL",
		"E1": "请求方式",
		"F1": "主机地址",
		"G1": "操作地点",
		"H1": "请求参数",
		"I1": "操作消息",
		"J1": "方法名称",
		"K1": "操作人员",
		"L1": "操作时间",
		"M1": "操作状态",
		"N1": "消耗时间",
	}
	// 从第二行开始的数据
	dataCells := make([]map[string]any, 0)
	for i, row := range rows {
		idx := strconv.Itoa(i + 2)
		// 业务类型
		businessType := ""
		// 状态
		statusValue := "失败"
		if row.StatusFlag == constants.STATUS_YES {
			statusValue = "成功"
		}
		dataCells = append(dataCells, map[string]any{
			"A" + idx: row.ID,
			"B" + idx: row.Title,
			"C" + idx: businessType,
			"D" + idx: row.OperaUrl,
			"E" + idx: row.OperaUrlMethod,
			"F" + idx: row.OperaIp,
			"G" + idx: row.OperaLocation,
			"H" + idx: row.OperaParam,
			"I" + idx: row.OperaMsg,
			"J" + idx: row.OperaMethod,
			"K" + idx: row.OperaBy,
			"L" + idx: date.ParseNumberToDate(row.OperaTime).Format(time.RFC3339),
			"M" + idx: statusValue,
			"N" + idx: row.CostTime,
		})
	}

	// 导出数据表格
	return file.WriteSheet(headerCells, dataCells, fileName, "")
}
