package repository

import (
	"fmt"
	db "mask_api_gin/src/framework/data_source"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/utils/date"
	"mask_api_gin/src/framework/utils/parse"

	"mask_api_gin/src/modules/monitor/model"
	"strings"
	"time"
)

// NewSysJobLogRepository 实例化数据层
var NewSysJobLogRepository = &SysJobLogRepositoryImpl{
	selectSql: `select 
    job_log_id, job_name, job_group, 
    invoke_target, target_params, job_msg, 
    status, create_time, cost_time 
	from sys_job_log`,

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

// SysJobLogRepositoryImpl 调度任务日志表 数据层处理
type SysJobLogRepositoryImpl struct {
	// 查询视图对象SQL
	selectSql string
	// 结果字段与实体映射
	resultMap map[string]string
}

// convertResultRows 将结果记录转实体结果组
func (r *SysJobLogRepositoryImpl) convertResultRows(rows []map[string]any) []model.SysJobLog {
	arr := make([]model.SysJobLog, 0)
	for _, row := range rows {
		sysJobLog := model.SysJobLog{}
		for key, value := range row {
			if keyMapper, ok := r.resultMap[key]; ok {
				db.SetFieldValue(&sysJobLog, keyMapper, value)
			}
		}
		arr = append(arr, sysJobLog)
	}
	return arr
}

// SelectByPage 分页查询集合
func (r *SysJobLogRepositoryImpl) SelectByPage(query map[string]any) map[string]any {
	// 查询条件拼接
	var conditions []string
	var params []any
	if v, ok := query["jobName"]; ok && v != "" {
		conditions = append(conditions, "job_name = ?")
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
		"total": int64(0),
		"rows":  []model.SysJobLog{},
	}

	// 查询数量 长度为0直接返回
	totalSql := "select count(1) as 'total' from sys_job_log"
	totalRows, err := db.RawDB("", totalSql+whereSql, params)
	if err != nil {
		logger.Errorf("total err => %v", err)
		return result
	}
	if total := parse.Number(totalRows[0]["total"]); total > 0 {
		result["total"] = total
	} else {
		return result
	}

	// 分页
	pageNum, pageSize := db.PageNumSize(query["pageNum"], query["pageSize"])
	pageSql := " order by job_log_id desc limit ?,? "
	params = append(params, pageNum*pageSize)
	params = append(params, pageSize)

	// 查询数据
	querySql := r.selectSql + whereSql + pageSql
	results, err := db.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return result
	}

	// 转换实体
	result["rows"] = r.convertResultRows(results)
	return result
}

// Select 查询集合
func (r *SysJobLogRepositoryImpl) Select(sysJobLog model.SysJobLog) []model.SysJobLog {
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
	results, err := db.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysJobLog{}
	}

	// 转换实体
	return r.convertResultRows(results)
}

// SelectById 通过ID查询信息
func (r *SysJobLogRepositoryImpl) SelectById(jobLogId string) model.SysJobLog {
	querySql := r.selectSql + " where job_log_id = ?"
	results, err := db.RawDB("", querySql, []any{jobLogId})
	if err != nil {
		logger.Errorf("query err => %v", err)
		return model.SysJobLog{}
	}
	// 转换实体
	if rows := r.convertResultRows(results); len(rows) > 0 {
		return rows[0]
	}
	return model.SysJobLog{}
}

// Insert 新增信息
func (r *SysJobLogRepositoryImpl) Insert(sysJobLog model.SysJobLog) string {
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
	keys, values, placeholder := db.KeyValuePlaceholderByInsert(params)
	sql := fmt.Sprintf("insert into sys_job_log (%s)values(%s)", keys, placeholder)

	tx := db.DB("").Begin() // 开启事务
	// 执行插入
	if err := tx.Exec(sql, values...).Error; err != nil {
		logger.Errorf("insert row : %v", err.Error())
		tx.Rollback()
		return ""
	}
	// 获取生成的自增 ID
	var insertedID string
	if err := tx.Raw("select last_insert_id()").Row().Scan(&insertedID); err != nil {
		logger.Errorf("insert last id : %v", err.Error())
		tx.Rollback()
		return ""
	}
	tx.Commit() // 提交事务
	return insertedID
}

// DeleteByIds 批量删除信息
func (r *SysJobLogRepositoryImpl) DeleteByIds(jobLogIds []string) int64 {
	placeholder := db.KeyPlaceholderByQuery(len(jobLogIds))
	sql := fmt.Sprintf("delete from sys_job_log where job_log_id in (%s)", placeholder)
	parameters := db.ConvertIdsSlice(jobLogIds)
	results, err := db.ExecDB("", sql, parameters)
	if err != nil {
		logger.Errorf("delete err => %v", err)
		return 0
	}
	return results
}

// Clean 清空集合数据
func (r *SysJobLogRepositoryImpl) Clean() error {
	sql := "truncate table sys_job_log"
	_, err := db.ExecDB("", sql, nil)
	return err
}
