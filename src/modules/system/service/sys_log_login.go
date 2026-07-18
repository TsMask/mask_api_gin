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
		LoginOs:       ilobArr[2],
		LoginBrowser:  ilobArr[3],
		LoginUserName: userName,
		LoginStatus:   status,
		LoginMsg:      msg,
	}
	sysSysLogLogin.CreateTime = time.Now().UnixMilli()
	return s.sysLogLoginRepository.Insert(sysSysLogLogin)
}

// Clean 清空系统登录日志
func (s SysLogLogin) Clean() int64 {
	return s.sysLogLoginRepository.Clean()
}

// ExportData 导出数据表格
func (s SysLogLogin) ExportData(rows []model.SysLogLogin, fileName string) (string, error) {
	// 第一行表头标题
	data := [][]string{
		{"序号", "用户账号", "登录状态", "登录地址", "登录地点", "浏览器", "操作系统", "提示消息", "访问时间"},
	}
	// 从第二行开始的数据
	for _, row := range rows {
		// 状态
		statusValue := "失败"
		if row.LoginStatus == constants.STATUS_YES {
			statusValue = "成功"
		}
		data = append(data, []string{
			fmt.Sprintf("%d", row.ID),
			row.LoginUserName,
			statusValue,
			row.LoginIp,
			row.LoginLocation,
			row.LoginBrowser,
			row.LoginOs,
			row.LoginMsg,
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
