package repository

import (
	"fmt"
	"mask_api_gin/src/framework/datasource"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/utils/date"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/framework/utils/repo"
	"mask_api_gin/src/modules/monitor/model"
	"strings"
)

// SysJobImpl 调度任务表 数据层处理
var SysJobImpl = &sysJobImpl{
	selectSql: `select job_id, job_name, job_group, invoke_target, target_params, cron_expression, 
	misfire_policy, concurrent, status, create_by, create_time, remark from sys_job`,

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
		"create_by":       "CreateBy",
		"create_time":     "CreateTime",
		"update_by":       "UpdateBy",
		"update_time":     "UpdateTime",
		"remark":          "Remark",
	},
}

type sysJobImpl struct {
	// 查询视图对象SQL
	selectSql string
	// 结果字段与实体映射
	resultMap map[string]string
}

// convertResultRows 将结果记录转实体结果组
func (r *sysJobImpl) convertResultRows(rows []map[string]interface{}) []model.SysJob {
	arr := make([]model.SysJob, 0)
	for _, row := range rows {
		sysJob := model.SysJob{}
		for key, value := range row {
			if keyMapper, ok := r.resultMap[key]; ok {
				repo.SetFieldValue(&sysJob, keyMapper, value)
			}
		}
		arr = append(arr, sysJob)
	}
	return arr
}

// SelectJobPage 分页查询调度任务集合
func (r *sysJobImpl) SelectJobPage(query map[string]string) map[string]interface{} {
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
	if v, ok := query["invokeTarget"]; ok {
		conditions = append(conditions, "invoke_target like concat(?, '%')")
		params = append(params, v)
	}
	if v, ok := query["status"]; ok {
		conditions = append(conditions, "status = ?")
		params = append(params, v)
	}

	// 构建查询条件语句
	whereSql := ""
	if len(conditions) > 0 {
		whereSql += " where " + strings.Join(conditions, " and ")
	}

	// 查询数量 长度为0直接返回
	totalSql := "select count(1) as 'total' from sys_job"
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
	pageSql := " limit ?,? "
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

// SelectJobList 查询调度任务集合
func (r *sysJobImpl) SelectJobList(sysJob model.SysJob) []model.SysJob {
	// 查询条件拼接
	var conditions []string
	var params []interface{}
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
	results, err := datasource.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysJob{}
	}

	// 转换实体
	return r.convertResultRows(results)
}

// SelectJobByIds 通过调度ID查询调度任务信息
func (r *sysJobImpl) SelectJobByIds(jobIds []string) []model.SysJob {
	placeholder := repo.KeyPlaceholderByQuery(len(jobIds))
	querySql := r.selectSql + " where job_id in (" + placeholder + ")"
	parameters := repo.ConvertIdsSlice(jobIds)
	results, err := datasource.RawDB("", querySql, parameters)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysJob{}
	}
	// 转换实体
	return r.convertResultRows(results)
}

// CheckUniqueJob 校验调度任务是否唯一
func (r *sysJobImpl) CheckUniqueJob(sysJob model.SysJob) string {
	// 查询条件拼接
	var conditions []string
	var params []interface{}
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
		return ""
	}

	// 查询数据
	querySql := "select job_id as 'str' from sys_job " + whereSql + " limit 1"
	results, err := datasource.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err %v", err)
		return ""
	}
	if len(results) > 0 {
		return fmt.Sprintf("%v", results[0]["str"])
	}
	return ""
}

// SelectJobByInvokeTarget 通过调用目标字符串查询调度任务信息
func (r *sysJobImpl) SelectJobByInvokeTarget(invokeTarget string) model.SysJob {
	return model.SysJob{}
}

// InsertJob 新增调度任务信息
func (r *sysJobImpl) InsertJob(sysJob model.SysJob) string {
	// 参数拼接
	params := make(map[string]interface{})
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
	if sysJob.Remark != "" {
		params["remark"] = sysJob.Remark
	}
	if sysJob.CreateBy != "" {
		params["create_by"] = sysJob.CreateBy
		params["create_time"] = date.NowTimestamp()
	}

	// 构建执行语句
	keys, placeholder, values := repo.KeyPlaceholderValueByInsert(params)
	sql := "insert into sys_job (" + strings.Join(keys, ",") + ")values(" + placeholder + ")"

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

// UpdateJob 修改调度任务信息
func (r *sysJobImpl) UpdateJob(sysJob model.SysJob) int64 {
	// 参数拼接
	params := make(map[string]interface{})
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
	if sysJob.Remark != "" {
		params["remark"] = sysJob.Remark
	}
	if sysJob.UpdateBy != "" {
		params["update_by"] = sysJob.UpdateBy
		params["update_time"] = date.NowTimestamp()
	}

	// 构建执行语句
	keys, values := repo.KeyValueByUpdate(params)
	sql := "update sys_job set " + strings.Join(keys, ",") + " where job_id = ?"

	// 执行更新
	values = append(values, sysJob.JobID)
	rows, err := datasource.ExecDB("", sql, values)
	if err != nil {
		logger.Errorf("update row : %v", err.Error())
		return 0
	}
	return rows
}

// DeleteJobByIds 批量删除调度任务信息
func (r *sysJobImpl) DeleteJobByIds(jobIds []string) int64 {
	placeholder := repo.KeyPlaceholderByQuery(len(jobIds))
	sql := "delete from sys_job where job_id in (" + placeholder + ")"
	parameters := repo.ConvertIdsSlice(jobIds)
	results, err := datasource.ExecDB("", sql, parameters)
	if err != nil {
		logger.Errorf("delete err => %v", err)
		return 0
	}
	return results
}
