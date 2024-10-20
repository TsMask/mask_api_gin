package repository

import (
	"fmt"
	db "mask_api_gin/src/framework/data_source"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/modules/monitor/model"
	"strings"
	"time"
)

// NewSysJob 实例化数据层
var NewSysJob = &SysJob{
	selectSql: `select 
    job_id, job_name, job_group, 
    invoke_target, target_params, cron_expression, 
	misfire_policy, concurrent, status, 
	save_log, create_by, create_time, remark 
	from sys_job`,

	resultMap: map[string]string{
		"job_id":          "JobID",
		"job_name":        "JobName",
		"job_group":       "JobGroup",
		"invoke_target":   "InvokeTarget",
		"target_params":   "TargetParams",
		"cron_expression": "CronExpression",
		"misfire_policy":  "MisfirePolicy",
		"concurrent":      "Concurrent",
		"status":          "Status",
		"save_log":        "SaveLog",
		"create_by":       "CreateBy",
		"create_time":     "CreateTime",
		"update_by":       "UpdateBy",
		"update_time":     "UpdateTime",
		"remark":          "Remark",
	},
}

// SysJob 调度任务 数据层处理
type SysJob struct {
	selectSql string            // 查询视图对象SQL
	resultMap map[string]string // 结果字段与实体映射
}

// convertResultRows 将结果记录转实体结果组
func (r *SysJob) convertResultRows(rows []map[string]any) []model.SysJob {
	arr := make([]model.SysJob, 0)
	for _, row := range rows {
		sysJob := model.SysJob{}
		for key, value := range row {
			if keyMapper, ok := r.resultMap[key]; ok {
				db.SetFieldValue(&sysJob, keyMapper, value)
			}
		}
		arr = append(arr, sysJob)
	}
	return arr
}

// SelectByPage 分页查询集合
func (r *SysJob) SelectByPage(query map[string]any) map[string]any {
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
	if v, ok := query["invokeTarget"]; ok && v != "" {
		conditions = append(conditions, "invoke_target like concat(?, '%')")
		params = append(params, v)
	}
	if v, ok := query["status"]; ok && v != "" {
		conditions = append(conditions, "status = ?")
		params = append(params, v)
	}

	// 构建查询条件语句
	whereSql := ""
	if len(conditions) > 0 {
		whereSql += " where " + strings.Join(conditions, " and ")
	}

	// 查询结果
	result := map[string]any{
		"total": int64(0),
		"rows":  []model.SysJob{},
	}

	// 查询数量 长度为0直接返回
	totalSql := "select count(1) as 'total' from sys_job"
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
	pageSql := " limit ?,? "
	params = append(params, pageNum*pageSize)
	params = append(params, pageSize)

	// 查询数据
	querySql := r.selectSql + whereSql + pageSql
	rows, err := db.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return result
	}

	// 转换实体
	result["rows"] = r.convertResultRows(rows)
	return result
}

// Select 查询集合
func (r *SysJob) Select(sysJob model.SysJob) []model.SysJob {
	// 查询条件拼接
	var conditions []string
	var params []any
	if sysJob.JobName != "" {
		conditions = append(conditions, "job_name like concat(?, '%')")
		params = append(params, sysJob.JobName)
	}
	if sysJob.JobGroup != "" {
		conditions = append(conditions, "job_group = ?")
		params = append(params, sysJob.JobGroup)
	}
	if sysJob.InvokeTarget != "" {
		conditions = append(conditions, "invoke_target like concat(?, '%')")
		params = append(params, sysJob.InvokeTarget)
	}
	if sysJob.Status != "" {
		conditions = append(conditions, "status = ?")
		params = append(params, sysJob.Status)
	}

	// 构建查询条件语句
	whereSql := ""
	if len(conditions) > 0 {
		whereSql += " where " + strings.Join(conditions, " and ")
	}

	// 查询数据
	querySql := r.selectSql + whereSql
	rows, err := db.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysJob{}
	}

	// 转换实体
	return r.convertResultRows(rows)
}

// SelectByIds 通过ID查询信息
func (r *SysJob) SelectByIds(jobIds []string) []model.SysJob {
	placeholder := db.KeyPlaceholderByQuery(len(jobIds))
	querySql := r.selectSql + " where job_id in (" + placeholder + ")"
	parameters := db.ConvertIdsSlice(jobIds)
	rows, err := db.RawDB("", querySql, parameters)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysJob{}
	}
	// 转换实体
	return r.convertResultRows(rows)
}

// Insert 新增信息
func (r *SysJob) Insert(sysJob model.SysJob) string {
	// 参数拼接
	params := make(map[string]any)
	if sysJob.JobID != "" {
		params["job_id"] = sysJob.JobID
	}
	if sysJob.JobName != "" {
		params["job_name"] = sysJob.JobName
	}
	if sysJob.JobGroup != "" {
		params["job_group"] = sysJob.JobGroup
	}
	if sysJob.InvokeTarget != "" {
		params["invoke_target"] = sysJob.InvokeTarget
	}
	if sysJob.TargetParams != "" {
		params["target_params"] = sysJob.TargetParams
	}
	if sysJob.CronExpression != "" {
		params["cron_expression"] = sysJob.CronExpression
	}
	if sysJob.MisfirePolicy != "" {
		params["misfire_policy"] = sysJob.MisfirePolicy
	}
	if sysJob.Concurrent != "" {
		params["concurrent"] = sysJob.Concurrent
	}
	if sysJob.Status != "" {
		params["status"] = sysJob.Status
	}
	if sysJob.SaveLog != "" {
		params["save_log"] = sysJob.SaveLog
	}
	if sysJob.Remark != "" {
		params["remark"] = sysJob.Remark
	}
	if sysJob.CreateBy != "" {
		params["create_by"] = sysJob.CreateBy
		params["create_time"] = time.Now().UnixMilli()
	}

	// 构建执行语句
	keys, values, placeholder := db.KeyValuePlaceholderByInsert(params)
	sql := fmt.Sprintf("insert into sys_job (%s)values(%s)", keys, placeholder)

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

// Update 修改信息
func (r *SysJob) Update(sysJob model.SysJob) int64 {
	// 参数拼接
	params := make(map[string]any)
	if sysJob.JobName != "" {
		params["job_name"] = sysJob.JobName
	}
	if sysJob.JobGroup != "" {
		params["job_group"] = sysJob.JobGroup
	}
	if sysJob.InvokeTarget != "" {
		params["invoke_target"] = sysJob.InvokeTarget
	}
	params["target_params"] = sysJob.TargetParams
	if sysJob.CronExpression != "" {
		params["cron_expression"] = sysJob.CronExpression
	}
	if sysJob.MisfirePolicy != "" {
		params["misfire_policy"] = sysJob.MisfirePolicy
	}
	if sysJob.Concurrent != "" {
		params["concurrent"] = sysJob.Concurrent
	}
	if sysJob.Status != "" {
		params["status"] = sysJob.Status
	}
	if sysJob.SaveLog != "" {
		params["save_log"] = sysJob.SaveLog
	}
	params["remark"] = sysJob.Remark
	if sysJob.UpdateBy != "" {
		params["update_by"] = sysJob.UpdateBy
		params["update_time"] = time.Now().UnixMilli()
	}

	// 构建执行语句
	keys, values := db.KeyValueByUpdate(params)
	sql := fmt.Sprintf("update sys_job set %s where job_id = ?", keys)

	// 执行更新
	values = append(values, sysJob.JobID)
	rows, err := db.ExecDB("", sql, values)
	if err != nil {
		logger.Errorf("update row : %v", err.Error())
		return 0
	}
	return rows
}

// DeleteByIds 批量删除信息
func (r *SysJob) DeleteByIds(jobIds []string) int64 {
	placeholder := db.KeyPlaceholderByQuery(len(jobIds))
	sql := fmt.Sprintf("delete from sys_job where job_id in (%s)", placeholder)
	parameters := db.ConvertIdsSlice(jobIds)
	results, err := db.ExecDB("", sql, parameters)
	if err != nil {
		logger.Errorf("delete err => %v", err)
		return 0
	}
	return results
}

// CheckUniqueJob 校验信息是否唯一
func (r *SysJob) CheckUniqueJob(sysJob model.SysJob) string {
	// 查询条件拼接
	var conditions []string
	var params []any
	if sysJob.JobName != "" {
		conditions = append(conditions, "job_name = ?")
		params = append(params, sysJob.JobName)
	}
	if sysJob.JobGroup != "" {
		conditions = append(conditions, "job_group = ?")
		params = append(params, sysJob.JobGroup)
	}

	// 构建查询条件语句
	whereSql := ""
	if len(conditions) > 0 {
		whereSql += " where " + strings.Join(conditions, " and ")
	} else {
		return "-"
	}

	// 查询数据
	querySql := fmt.Sprintf("select job_id as 'str' from sys_job %s limit 1", whereSql)
	results, err := db.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err %v", err)
		return "-"
	}
	if len(results) > 0 {
		return fmt.Sprint(results[0]["str"])
	}
	return ""
}
