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

// NewSysConfig 实例化数据层
var NewSysConfig = &SysConfigRepository{
	selectSql: `select
	config_id, config_name, config_key, config_value, config_type, 
	create_by, create_time, update_by, update_time, remark 
	from sys_config`,

	resultMap: map[string]string{
		"config_id":    "ConfigID",
		"config_name":  "ConfigName",
		"config_key":   "ConfigKey",
		"config_value": "ConfigValue",
		"config_type":  "ConfigType",
		"remark":       "Remark",
		"create_by":    "CreateBy",
		"create_time":  "CreateTime",
		"update_by":    "UpdateBy",
		"update_time":  "UpdateTime",
	},
}

// SysConfigRepository 参数配置表 数据层处理
type SysConfigRepository struct {
	selectSql string            // 查询视图对象SQL
	resultMap map[string]string // 结果字段与实体映射
}

// convertResultRows 将结果记录转实体结果组
func (r *SysConfigRepository) convertResultRows(rows []map[string]any) []model.SysConfig {
	arr := make([]model.SysConfig, 0)
	for _, row := range rows {
		sysConfig := model.SysConfig{}
		for key, value := range row {
			if keyMapper, ok := r.resultMap[key]; ok {
				db.SetFieldValue(&sysConfig, keyMapper, value)
			}
		}
		arr = append(arr, sysConfig)
	}
	return arr
}

// SelectByPage 分页查询集合
func (r *SysConfigRepository) SelectByPage(query map[string]any) map[string]any {
	// 查询条件拼接
	var conditions []string
	var params []any
	if v, ok := query["configName"]; ok && v != "" {
		conditions = append(conditions, "config_name like concat(?, '%')")
		params = append(params, v)
	}
	if v, ok := query["configType"]; ok && v != "" {
		conditions = append(conditions, "config_type = ?")
		params = append(params, v)
	}
	if v, ok := query["configKey"]; ok && v != "" {
		conditions = append(conditions, "config_key like concat(?, '%')")
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
		"rows":  []model.SysConfig{},
	}

	// 查询数量 长度为0直接返回
	totalSql := "select count(1) as 'total' from sys_config"
	totalRows, err := db.RawDB("", totalSql+whereSql, params)
	if err != nil {
		logger.Errorf("total err => %v", err)
		return result
	}

	if total := parse.Number(totalRows[0]["total"]); total > 0 {
		return result
	} else {
		result["total"] = total
	}

	// 分页
	pageNum, pageSize := db.PageNumSize(query["pageNum"], query["pageSize"])
	pageSql := " limit ?,? "
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
func (r *SysConfigRepository) Select(sysConfig model.SysConfig) []model.SysConfig {
	// 查询条件拼接
	var conditions []string
	var params []any
	if sysConfig.ConfigName != "" {
		conditions = append(conditions, "config_name like concat(?, '%')")
		params = append(params, sysConfig.ConfigName)
	}
	if sysConfig.ConfigType != "" {
		conditions = append(conditions, "config_type = ?")
		params = append(params, sysConfig.ConfigType)
	}
	if sysConfig.ConfigKey != "" {
		conditions = append(conditions, "config_key like concat(?, '%')")
		params = append(params, sysConfig.ConfigKey)
	}
	if sysConfig.CreateTime > 0 {
		conditions = append(conditions, "create_time >= ?")
		params = append(params, sysConfig.CreateTime)
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
		return []model.SysConfig{}
	}

	// 转换实体
	return r.convertResultRows(results)
}

// SelectByIds 通过ID查询信息
func (r *SysConfigRepository) SelectByIds(configIds []string) []model.SysConfig {
	placeholder := db.KeyPlaceholderByQuery(len(configIds))
	querySql := r.selectSql + " where config_id in (" + placeholder + ")"
	parameters := db.ConvertIdsSlice(configIds)
	results, err := db.RawDB("", querySql, parameters)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysConfig{}
	}
	// 转换实体
	return r.convertResultRows(results)
}

// Insert 新增信息
func (r *SysConfigRepository) Insert(sysConfig model.SysConfig) string {
	// 参数拼接
	params := make(map[string]any)
	if sysConfig.ConfigName != "" {
		params["config_name"] = sysConfig.ConfigName
	}
	if sysConfig.ConfigKey != "" {
		params["config_key"] = sysConfig.ConfigKey
	}
	if sysConfig.ConfigValue != "" {
		params["config_value"] = sysConfig.ConfigValue
	}
	if sysConfig.ConfigType != "" {
		params["config_type"] = sysConfig.ConfigType
	}
	if sysConfig.Remark != "" {
		params["remark"] = sysConfig.Remark
	}
	if sysConfig.CreateBy != "" {
		params["create_by"] = sysConfig.CreateBy
		params["create_time"] = time.Now().UnixMilli()
	}

	// 构建执行语句
	keys, values, placeholder := db.KeyValuePlaceholderByInsert(params)
	sql := fmt.Sprintf("insert into sys_config (%s)values(%s)", keys, placeholder)

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
func (r *SysConfigRepository) Update(sysConfig model.SysConfig) int64 {
	// 参数拼接
	params := make(map[string]any)
	if sysConfig.ConfigName != "" {
		params["config_name"] = sysConfig.ConfigName
	}
	if sysConfig.ConfigKey != "" {
		params["config_key"] = sysConfig.ConfigKey
	}
	if sysConfig.ConfigValue != "" {
		params["config_value"] = sysConfig.ConfigValue
	}
	if sysConfig.ConfigType != "" {
		params["config_type"] = sysConfig.ConfigType
	}
	params["remark"] = sysConfig.Remark
	if sysConfig.UpdateBy != "" {
		params["update_by"] = sysConfig.UpdateBy
		params["update_time"] = time.Now().UnixMilli()
	}

	// 构建执行语句
	keys, values := db.KeyValueByUpdate(params)
	sql := fmt.Sprintf("update sys_config set %s where config_id = ?", keys)

	// 执行更新
	values = append(values, sysConfig.ConfigID)
	rows, err := db.ExecDB("", sql, values)
	if err != nil {
		logger.Errorf("update row : %v", err.Error())
		return 0
	}
	return rows
}

// DeleteByIds 批量删除信息
func (r *SysConfigRepository) DeleteByIds(configIds []string) int64 {
	placeholder := db.KeyPlaceholderByQuery(len(configIds))
	sql := fmt.Sprintf("delete from sys_config where config_id in (%s)", placeholder)
	parameters := db.ConvertIdsSlice(configIds)
	results, err := db.ExecDB("", sql, parameters)
	if err != nil {
		logger.Errorf("delete err => %v", err)
		return 0
	}
	return results
}

// CheckUnique 检查信息是否唯一
func (r *SysConfigRepository) CheckUnique(sysConfig model.SysConfig) string {
	// 查询条件拼接
	var conditions []string
	var params []any
	if sysConfig.ConfigKey != "" {
		conditions = append(conditions, "config_key = ?")
		params = append(params, sysConfig.ConfigKey)
	}

	// 构建查询条件语句
	whereSql := ""
	if len(conditions) > 0 {
		whereSql += " where " + strings.Join(conditions, " and ")
	} else {
		return ""
	}

	// 查询数据
	querySql := fmt.Sprintf("select config_id as 'str' from sys_config %s limit 1", whereSql)
	results, err := db.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err %v", err)
		return ""
	}
	if len(results) > 0 {
		return fmt.Sprint(results[0]["str"])
	}
	return ""
}

// SelectValueByKey 通过Key查询Value
func (r *SysConfigRepository) SelectValueByKey(configKey string) string {
	querySql := "select config_value as 'str' from sys_config where config_key = ?"
	results, err := db.RawDB("", querySql, []any{configKey})
	if err != nil {
		logger.Errorf("query err => %v", err)
		return ""
	}
	if len(results) > 0 {
		return fmt.Sprint(results[0]["str"])
	}
	return ""
}
