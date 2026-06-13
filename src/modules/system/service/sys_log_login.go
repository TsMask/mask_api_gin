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

// NewSysLogLogin 实例化服务层
var NewSysLogLogin = &SysLogLogin{
	sysLogLoginRepository: repository.NewSysLogLogin,
}

// SysLogLogin 系统登录日志 服务层处理
type SysLogLogin struct {
	sysLogLoginRepository *repository.SysLogLogin // 系统登录日志信息
}

// FindByPage 分页查询列表数据
func (s SysLogLogin) FindByPage(query map[string]string, dataScopeSQL string) ([]model.SysLogLogin, int64) {
	return s.sysLogLoginRepository.SelectByPage(query, dataScopeSQL)
}

// Insert 新增信息
func (s SysLogLogin) Insert(userName, status, msg string, ilobArr [4]string) int64 {
	sysSysLogLogin := model.SysLogLogin{
		LoginIp:       ilobArr[0],
		LoginLocation: ilobArr[1],
		OS:            ilobArr[2],
		Browser:       ilobArr[3],
		UserName:      userName,
		StatusFlag:    status,
		Msg:           msg,
	}
	return s.sysLogLoginRepository.Insert(sysSysLogLogin)
}

// Clean 清空系统登录日志
func (s SysLogLogin) Clean() int64 {
	return s.sysLogLoginRepository.Clean()
}

// ExportData 导出数据表格
func (s SysLogLogin) ExportData(rows []model.SysLogLogin, fileName string) (string, error) {
	// 第一行表头标题
	headerCells := map[string]string{
		"A1": "序号",
		"B1": "用户账号",
		"C1": "登录状态",
		"D1": "登录地址",
		"E1": "登录地点",
		"F1": "浏览器",
		"G1": "操作系统",
		"H1": "提示消息",
		"I1": "访问时间",
	}
	// 从第二行开始的数据
	dataCells := make([]map[string]any, 0)
	for i, row := range rows {
		idx := strconv.Itoa(i + 2)
		// 状态
		statusValue := "失败"
		if row.StatusFlag == constants.STATUS_YES {
			statusValue = "成功"
		}
		dataCells = append(dataCells, map[string]any{
			"A" + idx: row.ID,
			"B" + idx: row.UserName,
			"C" + idx: statusValue,
			"D" + idx: row.LoginIp,
			"E" + idx: row.LoginLocation,
			"F" + idx: row.Browser,
			"G" + idx: row.OS,
			"H" + idx: row.Msg,
			"I" + idx: date.ParseNumberToDate(row.LoginTime).Format(time.RFC3339),
		})
	}

	// 导出数据表格
	return file.WriteSheet(headerCells, dataCells, fileName, "")
}
