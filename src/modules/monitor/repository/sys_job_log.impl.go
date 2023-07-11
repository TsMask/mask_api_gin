package repository

import (
	"mask_api_gin/src/framework/datasource"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/service/repo"
	"mask_api_gin/src/framework/utils/date"
	"mask_api_gin/src/modules/monitor/model"
	"strings"
)

// SysJobLogImpl 调度任务日志表 数据层处理
var SysJobLogImpl = &sysJobLogImpl{
	selectSql: "select job_log_id, job_name, job_group, invoke_target, target_params, job_msg, status, create_time from sys_job_log",
}

type sysJobLogImpl struct {
	// 查询视图对象SQL
	selectSql string
}

// 分页查询调度任务日志集合
func (r *sysJobLogImpl) SelectJobLogPage(query map[string]string) map[string]interface{} {
	db := datasource.GetDefaultDB()

	// 查询条件拼接
	var conditions []string
	var params []interface{}
	if v, ok := query["jobName"]; ok {
		conditions = append(conditions, "job_name like concat(?, '%')")
		params = append(params, v)
	}
	if v, ok := query["jobGroup"]; ok {
		conditions = append(conditions, "job_group = ?")
		params = append(params, v)
	}
	if v, ok := query["status"]; ok {
		conditions = append(conditions, "status = ?")
		params = append(params, v)
	}
	if v, ok := query["invokeTarget"]; ok {
		conditions = append(conditions, "invoke_target like concat(?, '%')")
		params = append(params, v)
	}
	if v, ok := query["beginTime"]; ok {
		conditions = append(conditions, "create_time >= ?")
		beginDate := date.ParseStrToDate(v, date.YYYY_MM_DD)
		params = append(params, beginDate.UnixNano()/1e6)
	}
	if v, ok := query["endTime"]; ok {
		conditions = append(conditions, "create_time <= ?")
		endDate := date.ParseStrToDate(v, date.YYYY_MM_DD)
		params = append(params, endDate.UnixNano()/1e6)
	}

	// 构建查询条件语句
	whereSql := ""
	if len(conditions) > 0 {
		whereSql += " where " + strings.Join(conditions, " and ")
	}

	// 查询数量 长度为0直接返回
	totalSql := "select count(1) as 'total' from sys_job_log"
	var total int64
	totalRes := db.Raw(totalSql+whereSql, params...).Scan(&total)
	if totalRes.Error != nil {
		logger.Errorf("SelectJobLogPage totalRes err %v", totalRes.Error)
	}
	if total <= 0 {
		return map[string]interface{}{
			"total": 0,
			"rows":  []interface{}{},
		}
	}

	// 分页
	pageNum, pageSize := repo.PageNumSize(query["pageNum"], query["pageSize"])
	pageSql := " order by job_log_id desc limit ?,? "
	params = append(params, pageNum*pageSize)
	params = append(params, pageSize)

	// 查询数据
	var sysJobLog []model.SysJobLog
	querySql := r.selectSql + whereSql + pageSql
	queryRes := db.Raw(querySql, params...).Scan(&sysJobLog)
	if queryRes.Error != nil {
		logger.Errorf("SelectJobLogPage queryRes err %v", queryRes.Error)
	}

	rows := repo.ConvertResultRows(sysJobLog)
	return map[string]interface{}{
		"total": total,
		"rows":  rows,
	}
}

// 查询调度任务日志集合
func (r *sysJobLogImpl) SelectJobLogList(sysJobLog model.SysJobLog) []model.SysJobLog {
	db := datasource.GetDefaultDB()

	// 查询条件拼接
	var conditions []string
	var params []interface{}
	if sysJobLog.JobName != "" {
		conditions = append(conditions, "job_name like concat(?, '%')")
		params = append(params, sysJobLog.JobName)
	}
	if sysJobLog.JobGroup != "" {
		conditions = append(conditions, "job_group = ?")
		params = append(params, sysJobLog.JobGroup)
	}
	if sysJobLog.Status != "" {
		conditions = append(conditions, "status = ?")
		params = append(params, sysJobLog.Status)
	}
	if sysJobLog.InvokeTarget != "" {
		conditions = append(conditions, "invoke_target like concat(?, '%')")
		params = append(params, sysJobLog.InvokeTarget)
	}

	// 构建查询条件语句
	whereSql := ""
	if len(conditions) > 0 {
		whereSql = " where " + strings.Join(conditions, " and ")
	}

	// 查询数据
	var results []model.SysJobLog
	querySql := r.selectSql + whereSql
	queryRes := db.Raw(querySql, params...).Scan(&results)
	if queryRes.Error != nil {
		logger.Errorf("SelectJobLogPage queryRes err %v", queryRes.Error)
	}
	return results
}

// 通过调度ID查询调度任务日志信息
func (r *sysJobLogImpl) SelectJobLogById(jobLogId string) model.SysJobLog {
	db := datasource.GetDefaultDB()

	// 查询数据
	var result model.SysJobLog
	querySql := r.selectSql + " where job_log_id = ?"
	queryRes := db.Raw(querySql, jobLogId).Scan(&result)
	if queryRes.Error != nil {
		logger.Errorf("SelectJobLogById queryRes err %v", queryRes.Error)
	}
	return result
}

// 新增调度任务日志信息
func (r *sysJobLogImpl) InsertJobLog(sysJobLog model.SysJobLog) string {
	db := datasource.GetDefaultDB()

	// 参数拼接
	paramMap := make(map[string]interface{})
	paramMap["create_time"] = date.NowTimestamp()
	if sysJobLog.JobLogID != "" {
		paramMap["job_log_id"] = sysJobLog.JobLogID
	}
	if sysJobLog.JobName != "" {
		paramMap["job_name"] = sysJobLog.JobName
	}
	if sysJobLog.JobGroup != "" {
		paramMap["job_group"] = sysJobLog.JobGroup
	}
	if sysJobLog.InvokeTarget != "" {
		paramMap["invoke_target"] = sysJobLog.InvokeTarget
	}
	if sysJobLog.TargetParams != "" {
		paramMap["target_params"] = sysJobLog.TargetParams
	}
	if sysJobLog.JobMsg != "" {
		paramMap["job_msg"] = sysJobLog.JobMsg
	}
	if sysJobLog.Status != "" {
		paramMap["status"] = sysJobLog.Status
	}

	// 构建执行语句
	keys, placeholder, values := repo.KeyValuePlaceholder(paramMap)
	sql := "insert into sys_job_log (" + strings.Join(keys, ",") + ")values(" + placeholder + ")"

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

// 批量删除调度任务日志信息
func (r *sysJobLogImpl) DeleteJobLogByIds(jobLogIds []string) int64 {
	db := datasource.GetDefaultDB()

	// 构建执行语句
	sql := "delete from sys_job_log where job_log_id in (?)"
	// 执行插入
	result := db.Exec(sql, jobLogIds)
	if err := result.Error; err != nil {
		logger.Errorf("delete rows : %v", err.Error())
		return 0
	}
	return result.RowsAffected
}

// 清空调度任务日志
func (r *sysJobLogImpl) CleanJobLog() error {
	db := datasource.GetDefaultDB()

	return db.Exec("truncate table sys_job_log").Error
}
