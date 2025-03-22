package service

import (
	"strconv"

	"mask_api_gin/src/framework/constants"
	"mask_api_gin/src/framework/utils/date"
	"mask_api_gin/src/framework/utils/file"
	"mask_api_gin/src/modules/oauth2/model"
	"mask_api_gin/src/modules/oauth2/repository"
)

// NewOauth2LogLogin 实例化服务层
var NewOauth2LogLogin = &Oauth2LogLogin{
	oauth2LogLoginRepository: repository.NewOauth2LogLogin,
}

// Oauth2LogLogin 用户授权第三方应用登录日志 服务层处理
type Oauth2LogLogin struct {
	oauth2LogLoginRepository *repository.Oauth2LogLoginRepository // 用户授权第三方应用登录日志信息
}

// FindByPage 分页查询列表数据
func (s Oauth2LogLogin) FindByPage(query map[string]string) ([]model.Oauth2LogLogin, int64) {
	return s.oauth2LogLoginRepository.SelectByPage(query)
}

// Insert 新增信息
func (s Oauth2LogLogin) Insert(clientId, status, msg string, ilobArr [4]string) int64 {
	sysOauth2LogLogin := model.Oauth2LogLogin{
		LoginIp:       ilobArr[0],
		LoginLocation: ilobArr[1],
		OS:            ilobArr[2],
		Browser:       ilobArr[3],
		ClientId:      clientId,
		StatusFlag:    status,
		Msg:           msg,
	}
	return s.oauth2LogLoginRepository.Insert(sysOauth2LogLogin)
}

// Clean 清空用户授权第三方应用登录日志
func (s Oauth2LogLogin) Clean() int64 {
	return s.oauth2LogLoginRepository.Clean()
}

// ExportData 导出数据表格
func (s Oauth2LogLogin) ExportData(rows []model.Oauth2LogLogin, fileName string) (string, error) {
	// 第一行表头标题
	headerCells := map[string]string{
		"A1": "序号",
		"B1": "应用的唯一标识",
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
			"B" + idx: row.ClientId,
			"C" + idx: statusValue,
			"D" + idx: row.LoginIp,
			"E" + idx: row.LoginLocation,
			"F" + idx: row.Browser,
			"G" + idx: row.OS,
			"H" + idx: row.Msg,
			"I" + idx: date.ParseDateToStr(row.LoginTime, date.YYYY_MM_DD_HH_MM_SS),
		})
	}

	// 导出数据表格
	return file.WriteSheet(headerCells, dataCells, fileName, "")
}
