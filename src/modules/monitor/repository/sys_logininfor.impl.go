package repository

import (
	"mask_api_gin/src/framework/datasource"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/utils/date"
	repoUtils "mask_api_gin/src/framework/utils/repo"
	"mask_api_gin/src/modules/monitor/model"
	"strings"
)

// SysLogininforImpl 系统登录访问表 数据层处理
var SysLogininforImpl = &sysLogininforImpl{
	selectSql: "",
}

type sysLogininforImpl struct {
	// 查询视图对象SQL
	selectSql string
}

// SelectLogininforPage 分页查询系统登录日志集合
func (r *sysLogininforImpl) SelectLogininforPage(query map[string]string) map[string]interface{} {
	return map[string]interface{}{}
}

// SelectLogininforList 查询系统登录日志集合
func (r *sysLogininforImpl) SelectLogininforList(sysLogininfor model.SysLogininfor) []model.SysLogininfor {
	return []model.SysLogininfor{}
}

// InsertLogininfor 新增系统登录日志
func (r *sysLogininforImpl) InsertLogininfor(sysLogininfor model.SysLogininfor) string {
	db := datasource.DefaultDB()

	// 参数拼接
	paramMap := make(map[string]interface{})
	paramMap["login_time"] = date.NowTimestamp()
	if sysLogininfor.UserName != "" {
		paramMap["user_name"] = sysLogininfor.UserName
	}
	if sysLogininfor.Status != "" {
		paramMap["status"] = sysLogininfor.Status
	}
	if sysLogininfor.IPAddr != "" {
		paramMap["ipaddr"] = sysLogininfor.IPAddr
	}
	if sysLogininfor.LoginLocation != "" {
		paramMap["login_location"] = sysLogininfor.LoginLocation
	}
	if sysLogininfor.Browser != "" {
		paramMap["browser"] = sysLogininfor.Browser
	}
	if sysLogininfor.OS != "" {
		paramMap["os"] = sysLogininfor.OS
	}
	if sysLogininfor.Msg != "" {
		paramMap["msg"] = sysLogininfor.Msg
	}

	// 构建执行语句
	keys, placeholder, values := repoUtils.KeyPlaceholderValueByInsert(paramMap)
	sql := "insert into sys_logininfor (" + strings.Join(keys, ",") + ")values(" + placeholder + ")"

	// 开启事务
	tx := db.Begin()
	// 执行插入
	err := tx.Exec(sql, values...).Error
	if err != nil {
		logger.Errorf("insert row : %v", err.Error())
		tx.Rollback()
		return err.Error()
	}
	// 获取生成的自增 ID
	var insertedID string
	err = tx.Raw("select last_insert_id()").Row().Scan(&insertedID)
	if err != nil {
		logger.Errorf("insert last id : %v", err.Error())
		tx.Rollback()
		return ""
	}
	// 提交事务
	tx.Commit()
	return insertedID
}

// DeleteLogininforByIds 批量删除系统登录日志
func (r *sysLogininforImpl) DeleteLogininforByIds(infoIds []string) int64 {
	return 0
}

// CleanLogininfor 清空系统登录日志
func (r *sysLogininforImpl) CleanLogininfor() error {
	return nil
}
