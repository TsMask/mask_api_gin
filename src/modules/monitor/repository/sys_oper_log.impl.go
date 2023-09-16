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

// 实例化数据层 SysOperLogImpl 结构体
var NewSysOperLogImpl = &SysOperLogImpl{
	selectSql: `select 
	oper_id, title, business_type, method, request_method, operator_type, oper_name, dept_name, 
	oper_url, oper_ip, oper_location, oper_param, oper_msg, status, oper_time, cost_time
	from sys_oper_log`,

	resultMap: map[string]string{
		"oper_id":        "OperID",
		"title":          "Title",
		"business_type":  "BusinessType",
		"method":         "Method",
		"request_method": "RequestMethod",
		"operator_type":  "OperatorType",
		"oper_name":      "OperName",
		"dept_name":      "DeptName",
		"oper_url":       "OperURL",
		"oper_ip":        "OperIP",
		"oper_location":  "OperLocation",
		"oper_param":     "OperParam",
		"oper_msg":       "OperMsg",
		"status":         "Status",
		"oper_time":      "OperTime",
		"cost_time":      "CostTime",
	},
}

// SysOperLogImpl 操作日志表 数据层处理
type SysOperLogImpl struct {
	// 查询视图对象SQL
	selectSql string
	// 结果字段与实体映射
	resultMap map[string]string
}

// convertResultRows 将结果记录转实体结果组
func (r *SysOperLogImpl) convertResultRows(rows []map[string]any) []model.SysOperLog {
	arr := make([]model.SysOperLog, 0)
	for _, row := range rows {
		sysOperLog := model.SysOperLog{}
		for key, value := range row {
			if keyMapper, ok := r.resultMap[key]; ok {
				repo.SetFieldValue(&sysOperLog, keyMapper, value)
			}
		}
		arr = append(arr, sysOperLog)
	}
	return arr
}

// SelectOperLogPage 分页查询系统操作日志集合
func (r *SysOperLogImpl) SelectOperLogPage(query map[string]any) map[string]any {
	// 查询条件拼接
	var conditions []string
	var params []any
	if v, ok := query["title"]; ok && v != "" {
		conditions = append(conditions, "title like concat(?, '%')")
		params = append(params, v)
	}
	if v, ok := query["businessType"]; ok && v != "" {
		conditions = append(conditions, "business_type = ?")
		params = append(params, v)
	}
	if v, ok := query["operName"]; ok && v != "" {
		conditions = append(conditions, "oper_name like concat(?, '%')")
		params = append(params, v)
	}
	if v, ok := query["status"]; ok && v != "" {
		conditions = append(conditions, "status = ?")
		params = append(params, v)
	}
	beginTime, ok := query["beginTime"]
	if !ok {
		beginTime, ok = query["params[beginTime]"]
	}
	if ok && beginTime != "" {
		conditions = append(conditions, "oper_time >= ?")
		beginDate := date.ParseStrToDate(beginTime.(string), date.YYYY_MM_DD)
		params = append(params, beginDate.UnixMilli())
	}
	endTime, ok := query["endTime"]
	if !ok {
		endTime, ok = query["params[endTime]"]
	}
	if ok && endTime != "" {
		conditions = append(conditions, "oper_time <= ?")
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
		"rows":  []model.SysOperLog{},
	}

	// 查询数量 长度为0直接返回
	totalSql := "select count(1) as 'total' from sys_oper_log"
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
	pageSql := " order by oper_id desc limit ?,? "
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

// SelectOperLogList 查询系统操作日志集合
func (r *SysOperLogImpl) SelectOperLogList(sysOperLog model.SysOperLog) []model.SysOperLog {
	// 查询条件拼接
	var conditions []string
	var params []any
	if sysOperLog.Title != "" {
		conditions = append(conditions, "title like concat(?, '%')")
		params = append(params, sysOperLog.Title)
	}
	if sysOperLog.BusinessType != "" {
		conditions = append(conditions, "business_type = ?")
		params = append(params, sysOperLog.BusinessType)
	}
	if sysOperLog.OperName != "" {
		conditions = append(conditions, "oper_name like concat(?, '%')")
		params = append(params, sysOperLog.OperName)
	}
	if sysOperLog.Status != "" {
		conditions = append(conditions, "status = ?")
		params = append(params, sysOperLog.Status)
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
		return []model.SysOperLog{}
	}

	// 转换实体
	return r.convertResultRows(results)
}

// SelectOperLogById 查询操作日志详细
func (r *SysOperLogImpl) SelectOperLogById(operId string) model.SysOperLog {
	querySql := r.selectSql + " where oper_id = ?"
	results, err := datasource.RawDB("", querySql, []any{operId})
	if err != nil {
		logger.Errorf("query err => %v", err)
		return model.SysOperLog{}
	}
	// 转换实体
	rows := r.convertResultRows(results)
	if len(rows) > 0 {
		return rows[0]
	}
	return model.SysOperLog{}
}

// InsertOperLog 新增操作日志
func (r *SysOperLogImpl) InsertOperLog(sysOperLog model.SysOperLog) string {
	// 参数拼接
	params := make(map[string]any)
	params["oper_time"] = time.Now().UnixMilli()
	if sysOperLog.Title != "" {
		params["title"] = sysOperLog.Title
	}
	if sysOperLog.BusinessType != "" {
		params["business_type"] = sysOperLog.BusinessType
	}
	if sysOperLog.Method != "" {
		params["method"] = sysOperLog.Method
	}
	if sysOperLog.RequestMethod != "" {
		params["request_method"] = sysOperLog.RequestMethod
	}
	if sysOperLog.OperatorType != "" {
		params["operator_type"] = sysOperLog.OperatorType
	}
	if sysOperLog.DeptName != "" {
		params["dept_name"] = sysOperLog.DeptName
	}
	if sysOperLog.OperName != "" {
		params["oper_name"] = sysOperLog.OperName
	}
	if sysOperLog.OperURL != "" {
		params["oper_url"] = sysOperLog.OperURL
	}
	if sysOperLog.OperIP != "" {
		params["oper_ip"] = sysOperLog.OperIP
	}
	if sysOperLog.OperLocation != "" {
		params["oper_location"] = sysOperLog.OperLocation
	}
	if sysOperLog.OperParam != "" {
		params["oper_param"] = sysOperLog.OperParam
	}
	if sysOperLog.OperMsg != "" {
		params["oper_msg"] = sysOperLog.OperMsg
	}
	if sysOperLog.Status != "" {
		params["status"] = sysOperLog.Status
	}
	if sysOperLog.CostTime > 0 {
		params["cost_time"] = sysOperLog.CostTime
	}

	// 构建执行语句
	keys, placeholder, values := repo.KeyPlaceholderValueByInsert(params)
	sql := "insert into sys_oper_log (" + strings.Join(keys, ",") + ")values(" + placeholder + ")"

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

// DeleteOperLogByIds 批量删除系统操作日志
func (r *SysOperLogImpl) DeleteOperLogByIds(operIds []string) int64 {
	placeholder := repo.KeyPlaceholderByQuery(len(operIds))
	sql := "delete from sys_oper_log where oper_id in (" + placeholder + ")"
	parameters := repo.ConvertIdsSlice(operIds)
	results, err := datasource.ExecDB("", sql, parameters)
	if err != nil {
		logger.Errorf("delete err => %v", err)
		return 0
	}
	return results
}

// CleanOperLog 清空操作日志
func (r *SysOperLogImpl) CleanOperLog() error {
	sql := "truncate table sys_oper_log"
	_, err := datasource.ExecDB("", sql, []any{})
	return err
}
