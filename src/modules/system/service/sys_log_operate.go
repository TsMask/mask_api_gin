package service

import (
	"mask_api_gin/src/framework/utils/date"
	"mask_api_gin/src/framework/utils/file"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/repository"
	"strconv"
)

// NewSysLogOperate 实例化服务层
var NewSysLogOperate = &SysLogOperate{
	SysLogOperate: repository.NewSysLogOperate,
}

// SysLogOperate 操作日志表 服务层处理
type SysLogOperate struct {
	SysLogOperate *repository.SysLogOperate // 操作日志信息
}

// FindByPage 分页查询列表数据
func (s SysLogOperate) FindByPage(query map[string]string, dataScopeSQL string) ([]model.SysLogOperate, int64) {
	return s.SysLogOperate.SelectByPage(query, dataScopeSQL)
}

// Insert 新增信息
func (s SysLogOperate) Insert(SysLogOperate model.SysLogOperate) string {
	return s.SysLogOperate.Insert(SysLogOperate)
}

// Clean 清空操作日志
func (s SysLogOperate) Clean() int64 {
	return s.SysLogOperate.Clean()
}

// ExportData 导出数据表格
func (s SysLogOperate) ExportData(rows []model.SysLogOperate, fileName string) (string, error) {
	// 第一行表头标题
	headerCells := map[string]string{
		"A1": "操作序号",
		"B1": "操作模块",
		"C1": "业务类型",
		"D1": "请求方法",
		"E1": "请求方式",
		"F1": "操作类别",
		"G1": "操作人员",
		"H1": "部门名称",
		"I1": "请求地址",
		"J1": "操作地址",
		"K1": "操作地点",
		"L1": "请求参数",
		"M1": "操作消息",
		"N1": "状态",
		"O1": "消耗时间（毫秒）",
		"P1": "操作时间",
	}
	// 从第二行开始的数据
	dataCells := make([]map[string]any, 0)
	for i, row := range rows {
		idx := strconv.Itoa(i + 2)
		// 业务类型
		businessType := ""
		// 操作类别
		operatorType := ""
		// 状态
		statusValue := "失败"
		if row.StatusFlag == "1" {
			statusValue = "成功"
		}
		dataCells = append(dataCells, map[string]any{
			"A" + idx: row.ID,
			"B" + idx: row.Title,
			"C" + idx: businessType,
			"D" + idx: row.OperaMethod,
			"E" + idx: row.OperaUrlMethod,
			"F" + idx: operatorType,
			"G" + idx: row.OperaBy,
			"I" + idx: row.OperaUrl,
			"J" + idx: row.OperaIp,
			"K" + idx: row.OperaLocation,
			"L" + idx: row.OperaParam,
			"M" + idx: row.OperaMsg,
			"N" + idx: statusValue,
			"O" + idx: row.CostTime,
			"P" + idx: date.ParseDateToStr(row.OperaTime, date.YYYY_MM_DD_HH_MM_SS),
		})
	}

	// 导出数据表格
	return file.WriteSheet(headerCells, dataCells, fileName, "")
}
