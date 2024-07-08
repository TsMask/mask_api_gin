package repository

import (
	"fmt"
	db "mask_api_gin/src/framework/data_source"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/utils/date"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/modules/system/model"
	"strings"
	"time"
)

// NewSysLogOperate 实例化数据层
var NewSysLogOperate = &SysLogOperateRepository{
	selectSql: `select 
	opera_id, title, business_type, method, request_method, operator_type, 
	opera_name, dept_name, opera_url, opera_ip, opera_location, opera_param, 
	opera_msg, status, opera_time, cost_time
	from sys_log_operate`,

	resultMap: map[string]string{
		"opera_id":       "OperaID",
		"title":          "Title",
		"business_type":  "BusinessType",
		"method":         "Method",
		"request_method": "RequestMethod",
		"operator_type":  "OperatorType",
		"opera_name":     "OperaName",
		"dept_name":      "DeptName",
		"opera_url":      "OperaURL",
		"opera_ip":       "OperaIP",
		"opera_location": "OperaLocation",
		"opera_param":    "OperaParam",
		"opera_msg":      "OperaMsg",
		"status":         "Status",
		"opera_time":     "OperaTime",
		"cost_time":      "CostTime",
	},
}

// SysLogOperateRepository 操作日志表 数据层处理
type SysLogOperateRepository struct {
	selectSql string            // 查询视图对象SQL
	resultMap map[string]string // 结果字段与实体映射
}

// convertResultRows 将结果记录转实体结果组
func (r *SysLogOperateRepository) convertResultRows(rows []map[string]any) []model.SysLogOperate {
	arr := make([]model.SysLogOperate, 0)
	for _, row := range rows {
		SysLogOperate := model.SysLogOperate{}
		for key, value := range row {
			if keyMapper, ok := r.resultMap[key]; ok {
				db.SetFieldValue(&SysLogOperate, keyMapper, value)
			}
		}
		arr = append(arr, SysLogOperate)
	}
	return arr
}

// SelectByPage 分页查询集合
func (r *SysLogOperateRepository) SelectByPage(query map[string]any) map[string]any {
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
	if v, ok := query["operaName"]; ok && v != "" {
		conditions = append(conditions, "opera_name like concat(?, '%')")
		params = append(params, v)
	}
	if v, ok := query["operaIp"]; ok && v != "" {
		conditions = append(conditions, "opera_ip like concat(?, '%')")
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
		conditions = append(conditions, "opera_time >= ?")
		beginDate := date.ParseStrToDate(beginTime.(string), date.YYYY_MM_DD)
		params = append(params, beginDate.UnixMilli())
	}
	endTime, ok := query["endTime"]
	if !ok {
		endTime, ok = query["params[endTime]"]
	}
	if ok && endTime != "" {
		conditions = append(conditions, "opera_time <= ?")
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
		"rows":  []model.SysLogOperate{},
	}

	// 查询数量 长度为0直接返回
	totalSql := "select count(1) as 'total' from sys_log_operate"
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
	pageSql := " order by opera_id desc limit ?,? "
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
func (r *SysLogOperateRepository) Select(SysLogOperate model.SysLogOperate) []model.SysLogOperate {
	// 查询条件拼接
	var conditions []string
	var params []any
	if SysLogOperate.Title != "" {
		conditions = append(conditions, "title like concat(?, '%')")
		params = append(params, SysLogOperate.Title)
	}
	if SysLogOperate.BusinessType != "" {
		conditions = append(conditions, "business_type = ?")
		params = append(params, SysLogOperate.BusinessType)
	}
	if SysLogOperate.OperaName != "" {
		conditions = append(conditions, "opera_name like concat(?, '%')")
		params = append(params, SysLogOperate.OperaName)
	}
	if SysLogOperate.Status != "" {
		conditions = append(conditions, "status = ?")
		params = append(params, SysLogOperate.Status)
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
		return []model.SysLogOperate{}
	}

	// 转换实体
	return r.convertResultRows(results)
}

// SelectById 通过ID查询信息
func (r *SysLogOperateRepository) SelectById(operaId string) model.SysLogOperate {
	querySql := r.selectSql + " where opera_id = ?"
	results, err := db.RawDB("", querySql, []any{operaId})
	if err != nil {
		logger.Errorf("query err => %v", err)
		return model.SysLogOperate{}
	}
	// 转换实体
	rows := r.convertResultRows(results)
	if len(rows) > 0 {
		return rows[0]
	}
	return model.SysLogOperate{}
}

// Insert 新增信息
func (r *SysLogOperateRepository) Insert(SysLogOperate model.SysLogOperate) string {
	// 参数拼接
	params := make(map[string]any)
	params["opera_time"] = time.Now().UnixMilli()
	if SysLogOperate.Title != "" {
		params["title"] = SysLogOperate.Title
	}
	if SysLogOperate.BusinessType != "" {
		params["business_type"] = SysLogOperate.BusinessType
	}
	if SysLogOperate.Method != "" {
		params["method"] = SysLogOperate.Method
	}
	if SysLogOperate.RequestMethod != "" {
		params["request_method"] = SysLogOperate.RequestMethod
	}
	if SysLogOperate.OperatorType != "" {
		params["operator_type"] = SysLogOperate.OperatorType
	}
	if SysLogOperate.DeptName != "" {
		params["dept_name"] = SysLogOperate.DeptName
	}
	if SysLogOperate.OperaName != "" {
		params["opera_name"] = SysLogOperate.OperaName
	}
	if SysLogOperate.OperaURL != "" {
		params["opera_url"] = SysLogOperate.OperaURL
	}
	if SysLogOperate.OperaIP != "" {
		params["opera_ip"] = SysLogOperate.OperaIP
	}
	if SysLogOperate.OperaLocation != "" {
		params["opera_location"] = SysLogOperate.OperaLocation
	}
	if SysLogOperate.OperaParam != "" {
		params["opera_param"] = SysLogOperate.OperaParam
	}
	if SysLogOperate.OperaMsg != "" {
		params["opera_msg"] = SysLogOperate.OperaMsg
	}
	if SysLogOperate.Status != "" {
		params["status"] = parse.Number(SysLogOperate.Status)
	}
	if SysLogOperate.CostTime > 0 {
		params["cost_time"] = SysLogOperate.CostTime
	}

	// 构建执行语句
	keys, values, placeholder := db.KeyValuePlaceholderByInsert(params)
	sql := fmt.Sprintf("insert into sys_log_operate (%s)values(%s)", keys, placeholder)

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
func (r *SysLogOperateRepository) DeleteByIds(operaIds []string) int64 {
	placeholder := db.KeyPlaceholderByQuery(len(operaIds))
	sql := fmt.Sprintf("delete from sys_log_operate where opera_id in (%s)", placeholder)
	parameters := db.ConvertIdsSlice(operaIds)
	results, err := db.ExecDB("", sql, parameters)
	if err != nil {
		logger.Errorf("delete err => %v", err)
		return 0
	}
	return results
}

// Clean 清空信息
func (r *SysLogOperateRepository) Clean() error {
	sql := "truncate table sys_log_operate"
	_, err := db.ExecDB("", sql, []any{})
	return err
}
