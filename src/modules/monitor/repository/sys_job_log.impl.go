package repository

import (
	"mask_api_gin/src/framework/datasource"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/utils/date"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/framework/utils/repo"
	"mask_api_gin/src/modules/monitor/model"
	"strings"
)

// SysJobLogImpl 调度任务日志表 数据层处理
var SysJobLogImpl = &sysJobLogImpl{
	selectSql: `select job_log_id, job_name, job_group, invoke_target, 
	target_params, job_msg, status, create_time from sys_job_log`,

	resultMap: map[string]string{
		"job_log_id":    "JobLogID",
		"job_name":      "JobName",
		"job_group":     "JobGroup",
		"invoke_target": "InvokeTarget",
		"target_params": "TargetParams",
		"job_msg":       "JobMsg",
		"status":        "Status",
		"create_time":   "CreateTime",
	},
}

type sysJobLogImpl struct {
	// 查询视图对象SQL
	selectSql string
	// 结果字段与实体映射
	resultMap map[string]string
}

// convertResultRows 将结果记录转实体结果组
func (r *sysJobLogImpl) convertResultRows(rows []map[string]interface{}) []model.SysJobLog {
	arr := make([]model.SysJobLog, 0)
	for _, row := range rows {
		sysJobLog := model.SysJobLog{}
		for key, value := range row {
			if keyMapper, ok := r.resultMap[key]; ok {
				repo.SetFieldValue(&sysJobLog, keyMapper, value)
			}
		}
		arr = append(arr, sysJobLog)
	}
	return arr
}

// 分页查询调度任务日志集合
func (r *sysJobLogImpl) SelectJobLogPage(query map[string]string) map[string]interface{} {
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
	totalRows, err := datasource.RawDB("", totalSql+whereSql, params)
	if err != nil {
		logger.Errorf("total err => %v", err)
	}
	total := parse.Number(totalRows[0]["total"])
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
	querySql := r.selectSql + whereSql + pageSql
	results, err := datasource.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err => %v", err)
	}

	// 转换实体
	rows := r.convertResultRows(results)
	return map[string]interface{}{
		"total": total,
		"rows":  rows,
	}
}

// 查询调度任务日志集合
func (r *sysJobLogImpl) SelectJobLogList(sysJobLog model.SysJobLog) []model.SysJobLog {
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
		whereSql += " where " + strings.Join(conditions, " and ")
	}

	// 查询数据
	querySql := r.selectSql + whereSql
	results, err := datasource.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysJobLog{}
	}

	// 转换实体
	return r.convertResultRows(results)
}

// 通过调度ID查询调度任务日志信息
func (r *sysJobLogImpl) SelectJobLogById(jobLogId string) model.SysJobLog {
	querySql := r.selectSql + " where job_log_id = ?"
	results, err := datasource.RawDB("", querySql, []interface{}{jobLogId})
	if err != nil {
		logger.Errorf("query err => %v", err)
		return model.SysJobLog{}
	}
	// 转换实体
	rows := r.convertResultRows(results)
	if len(rows) > 0 {
		return rows[0]
	}
	return model.SysJobLog{}
}

// 新增调度任务日志信息
func (r *sysJobLogImpl) InsertJobLog(sysJobLog model.SysJobLog) string {
	// 参数拼接
	params := make(map[string]interface{})
	params["create_time"] = date.NowTimestamp()
	if sysJobLog.JobLogID != "" {
		params["job_log_id"] = sysJobLog.JobLogID
	}
	if sysJobLog.JobName != "" {
		params["job_name"] = sysJobLog.JobName
	}
	if sysJobLog.JobGroup != "" {
		params["job_group"] = sysJobLog.JobGroup
	}
	if sysJobLog.InvokeTarget != "" {
		params["invoke_target"] = sysJobLog.InvokeTarget
	}
	if sysJobLog.TargetParams != "" {
		params["target_params"] = sysJobLog.TargetParams
	}
	if sysJobLog.JobMsg != "" {
		params["job_msg"] = sysJobLog.JobMsg
	}
	if sysJobLog.Status != "" {
		params["status"] = sysJobLog.Status
	}

	// 构建执行语句
	keys, placeholder, values := repo.KeyPlaceholderValueByInsert(params)
	sql := "insert into sys_job_log (" + strings.Join(keys, ",") + ")values(" + placeholder + ")"

	db := datasource.DefaultDB()
	// 开启事务
	tx := db.Begin()
	// 执行插入
	err := tx.Exec(sql, values...).Error
	if err != nil {
		logger.Errorf("insert row : %v", err.Error())
		tx.Rollback()
		return ""
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
	placeholder := repo.KeyPlaceholderByQuery(len(jobLogIds))
	sql := "delete from sys_job_log where job_log_id in (" + placeholder + ")"
	parameters := repo.ConvertIdsSlice(jobLogIds)
	results, err := datasource.ExecDB("", sql, parameters)
	if err != nil {
		logger.Errorf("delete err => %v", err)
		return 0
	}
	return results
}

// 清空调度任务日志
func (r *sysJobLogImpl) CleanJobLog() error {
	sql := "truncate table sys_job_log"
	results, err := datasource.ExecDB("", sql, []interface{}{})
	logger.Errorf("delete results => %v", results)
	return err
}
