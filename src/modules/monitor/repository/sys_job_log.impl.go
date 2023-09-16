package repository

import (
	"mask_api_gin/src/framework/datasource"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/utils/date"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/framework/utils/repo"
	"mask_api_gin/src/modules/monitor/model"
	"strings"
	"time"
)

// 实例化数据层 SysJobLogImpl 结构体
var NewSysJobLogImpl = &SysJobLogImpl{
	selectSql: `select job_log_id, job_name, job_group, invoke_target, 
	target_params, job_msg, status, create_time, cost_time from sys_job_log`,

	resultMap: map[string]string{
		"job_log_id":    "JobLogID",
		"job_name":      "JobName",
		"job_group":     "JobGroup",
		"invoke_target": "InvokeTarget",
		"target_params": "TargetParams",
		"job_msg":       "JobMsg",
		"status":        "Status",
		"create_time":   "CreateTime",
		"cost_time":     "CostTime",
	},
}

// SysJobLogImpl 调度任务日志表 数据层处理
type SysJobLogImpl struct {
	// 查询视图对象SQL
	selectSql string
	// 结果字段与实体映射
	resultMap map[string]string
}

// convertResultRows 将结果记录转实体结果组
func (r *SysJobLogImpl) convertResultRows(rows []map[string]any) []model.SysJobLog {
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
func (r *SysJobLogImpl) SelectJobLogPage(query map[string]any) map[string]any {
	// 查询条件拼接
	var conditions []string
	var params []any
	if v, ok := query["jobName"]; ok && v != "" {
		conditions = append(conditions, "job_name like concat(?, '%')")
		params = append(params, v)
	}
	if v, ok := query["jobGroup"]; ok && v != "" {
		conditions = append(conditions, "job_group = ?")
		params = append(params, v)
	}
	if v, ok := query["status"]; ok && v != "" {
		conditions = append(conditions, "status = ?")
		params = append(params, v)
	}
	if v, ok := query["invokeTarget"]; ok && v != "" {
		conditions = append(conditions, "invoke_target like concat(?, '%')")
		params = append(params, v)
	}
	beginTime, ok := query["beginTime"]
	if !ok {
		beginTime, ok = query["params[beginTime]"]
	}
	if ok && beginTime != "" {
		conditions = append(conditions, "create_time >= ?")
		beginDate := date.ParseStrToDate(beginTime.(string), date.YYYY_MM_DD)
		params = append(params, beginDate.UnixMilli())
	}
	endTime, ok := query["endTime"]
	if !ok {
		endTime, ok = query["params[endTime]"]
	}
	if ok && endTime != "" {
		conditions = append(conditions, "create_time <= ?")
		endDate := date.ParseStrToDate(endTime.(string), date.YYYY_MM_DD)
		params = append(params, endDate.UnixMilli())
	}

	// 构建查询条件语句
	whereSql := ""
	if len(conditions) > 0 {
		whereSql += " where " + strings.Join(conditions, " and ")
	}

	// 查询结果
	result := map[string]any{
		"total": 0,
		"rows":  []model.SysJobLog{},
	}

	// 查询数量 长度为0直接返回
	totalSql := "select count(1) as 'total' from sys_job_log"
	totalRows, err := datasource.RawDB("", totalSql+whereSql, params)
	if err != nil {
		logger.Errorf("total err => %v", err)
		return result
	}
	total := parse.Number(totalRows[0]["total"])
	if total == 0 {
		return result
	} else {
		result["total"] = total
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
		return result
	}

	// 转换实体
	result["rows"] = r.convertResultRows(results)
	return result
}

// 查询调度任务日志集合
func (r *SysJobLogImpl) SelectJobLogList(sysJobLog model.SysJobLog) []model.SysJobLog {
	// 查询条件拼接
	var conditions []string
	var params []any
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
func (r *SysJobLogImpl) SelectJobLogById(jobLogId string) model.SysJobLog {
	querySql := r.selectSql + " where job_log_id = ?"
	results, err := datasource.RawDB("", querySql, []any{jobLogId})
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
func (r *SysJobLogImpl) InsertJobLog(sysJobLog model.SysJobLog) string {
	// 参数拼接
	params := make(map[string]any)
	params["create_time"] = time.Now().UnixMilli()
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
	if sysJobLog.CostTime > 0 {
		params["cost_time"] = sysJobLog.CostTime
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
func (r *SysJobLogImpl) DeleteJobLogByIds(jobLogIds []string) int64 {
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
func (r *SysJobLogImpl) CleanJobLog() error {
	sql := "truncate table sys_job_log"
	_, err := datasource.ExecDB("", sql, []any{})
	return err
}
