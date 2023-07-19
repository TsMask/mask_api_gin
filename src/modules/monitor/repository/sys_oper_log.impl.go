package repository

import (
	"mask_api_gin/src/framework/datasource"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/utils/date"
	repoUtils "mask_api_gin/src/framework/utils/repo"
	"mask_api_gin/src/modules/monitor/model"
	"strings"
)

// SysOperLogImpl 操作日志表 数据层处理
var SysOperLogImpl = &sysOperLogImpl{
	selectSql: "",
}

type sysOperLogImpl struct {
	// 查询视图对象SQL
	selectSql string
}

// SelectOperLogPage 分页查询系统操作日志集合
func (r *sysOperLogImpl) SelectOperLogPage(query map[string]string) map[string]interface{} {
	return map[string]interface{}{}
}

// SelectOperLogList 查询系统操作日志集合
func (r *sysOperLogImpl) SelectOperLogList(sysOperLog model.SysOperLog) []model.SysOperLog {
	return []model.SysOperLog{}
}

// InsertOperLog 新增操作日志
func (r *sysOperLogImpl) InsertOperLog(sysOperLog model.SysOperLog) string {
	// 参数拼接
	paramMap := make(map[string]interface{})
	paramMap["oper_time"] = date.NowTimestamp()
	if sysOperLog.Title != "" {
		paramMap["title"] = sysOperLog.Title
	}
	if sysOperLog.BusinessType != "" {
		paramMap["business_type"] = sysOperLog.BusinessType
	}
	if sysOperLog.Method != "" {
		paramMap["method"] = sysOperLog.Method
	}
	if sysOperLog.RequestMethod != "" {
		paramMap["request_method"] = sysOperLog.RequestMethod
	}
	if sysOperLog.OperatorType != "" {
		paramMap["operator_type"] = sysOperLog.OperatorType
	}
	if sysOperLog.DeptName != "" {
		paramMap["dept_name"] = sysOperLog.DeptName
	}
	if sysOperLog.OperName != "" {
		paramMap["oper_name"] = sysOperLog.OperName
	}
	if sysOperLog.OperURL != "" {
		paramMap["oper_url"] = sysOperLog.OperURL
	}
	if sysOperLog.OperID != "" {
		paramMap["oper_ip"] = sysOperLog.OperID
	}
	if sysOperLog.OperLocation != "" {
		paramMap["oper_location"] = sysOperLog.OperLocation
	}
	if sysOperLog.OperParam != "" {
		paramMap["oper_param"] = sysOperLog.OperParam
	}
	if sysOperLog.OperMsg != "" {
		paramMap["oper_msg"] = sysOperLog.OperMsg
	}
	if sysOperLog.Status != "" {
		paramMap["status"] = sysOperLog.Status
	}
	if sysOperLog.CostTime > 0 {
		paramMap["cost_time"] = sysOperLog.CostTime
	}

	// 构建执行语句
	keys, placeholder, values := repoUtils.KeyPlaceholderValueByInsert(paramMap)
	sql := "insert into sys_oper_log (" + strings.Join(keys, ",") + ")values(" + placeholder + ")"

	db := datasource.DefaultDB()
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

// DeleteOperLogByIds 批量删除系统操作日志
func (r *sysOperLogImpl) DeleteOperLogByIds(operIds []string) int64 {
	return 0
}

// SelectOperLogById 查询操作日志详细
func (r *sysOperLogImpl) SelectOperLogById(operId string) model.SysOperLog {
	return model.SysOperLog{}
}

// CleanOperLog 清空操作日志
func (r *sysOperLogImpl) CleanOperLog() error {
	return nil
}
