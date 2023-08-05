package repository

import (
	"fmt"
	"mask_api_gin/src/framework/datasource"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/utils/date"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/framework/utils/repo"
	"mask_api_gin/src/modules/system/model"
	"strings"
	"time"
)

// 实例化数据层 SysConfigImpl 结构体
var NewSysConfigImpl = &SysConfigImpl{
	selectSql: `select
	config_id, config_name, config_key, config_value, config_type, create_by, create_time, update_by, update_time, remark 
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

// SysConfigImpl 参数配置表 数据层处理
type SysConfigImpl struct {
	// 查询视图对象SQL
	selectSql string
	// 结果字段与实体映射
	resultMap map[string]string
}

// convertResultRows 将结果记录转实体结果组
func (r *SysConfigImpl) convertResultRows(rows []map[string]any) []model.SysConfig {
	arr := make([]model.SysConfig, 0)
	for _, row := range rows {
		sysConfig := model.SysConfig{}
		for key, value := range row {
			if keyMapper, ok := r.resultMap[key]; ok {
				repo.SetFieldValue(&sysConfig, keyMapper, value)
			}
		}
		arr = append(arr, sysConfig)
	}
	return arr
}

// SelectDictDataPage 分页查询参数配置列表数据
func (r *SysConfigImpl) SelectConfigPage(query map[string]any) map[string]any {
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

	// 查询数量 长度为0直接返回
	totalSql := "select count(1) as 'total' from sys_config"
	totalRows, err := datasource.RawDB("", totalSql+whereSql, params)
	if err != nil {
		logger.Errorf("total err => %v", err)
	}
	total := parse.Number(totalRows[0]["total"])
	if total == 0 {
		return map[string]any{
			"total": total,
			"rows":  []model.SysConfig{},
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
	return map[string]any{
		"total": total,
		"rows":  rows,
	}
}

// SelectConfigList 查询参数配置列表
func (r *SysConfigImpl) SelectConfigList(sysConfig model.SysConfig) []model.SysConfig {
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
	results, err := datasource.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysConfig{}
	}

	// 转换实体
	return r.convertResultRows(results)
}

// SelectConfigValueByKey 通过参数键名查询参数键值
func (r *SysConfigImpl) SelectConfigValueByKey(configKey string) string {
	querySql := "select config_value as 'str' from sys_config where config_key = ?"
	results, err := datasource.RawDB("", querySql, []any{configKey})
	if err != nil {
		logger.Errorf("query err => %v", err)
		return ""
	}
	if len(results) > 0 {
		return fmt.Sprintf("%v", results[0]["str"])
	}
	return ""
}

// SelectConfigByIds 通过配置ID查询参数配置信息
func (r *SysConfigImpl) SelectConfigByIds(configIds []string) []model.SysConfig {
	placeholder := repo.KeyPlaceholderByQuery(len(configIds))
	querySql := r.selectSql + " where config_id in (" + placeholder + ")"
	parameters := repo.ConvertIdsSlice(configIds)
	results, err := datasource.RawDB("", querySql, parameters)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysConfig{}
	}
	// 转换实体
	return r.convertResultRows(results)
}

// CheckUniqueConfig 校验配置参数是否唯一
func (r *SysConfigImpl) CheckUniqueConfig(sysConfig model.SysConfig) string {
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
	querySql := "select config_id as 'str' from sys_config " + whereSql + " limit 1"
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

// InsertConfig 新增参数配置
func (r *SysConfigImpl) InsertConfig(sysConfig model.SysConfig) string {
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
	keys, placeholder, values := repo.KeyPlaceholderValueByInsert(params)
	sql := "insert into sys_config (" + strings.Join(keys, ",") + ")values(" + placeholder + ")"

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

// UpdateConfig 修改参数配置
func (r *SysConfigImpl) UpdateConfig(sysConfig model.SysConfig) int64 {
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
	if sysConfig.UpdateBy != "" {
		params["update_by"] = sysConfig.UpdateBy
		params["update_time"] = time.Now().UnixMilli()
	}

	// 构建执行语句
	keys, values := repo.KeyValueByUpdate(params)
	sql := "update sys_config set " + strings.Join(keys, ",") + " where config_id = ?"

	// 执行更新
	values = append(values, sysConfig.ConfigID)
	rows, err := datasource.ExecDB("", sql, values)
	if err != nil {
		logger.Errorf("update row : %v", err.Error())
		return 0
	}
	return rows
}

// DeleteConfigByIds 批量删除参数配置信息
func (r *SysConfigImpl) DeleteConfigByIds(configIds []string) int64 {
	placeholder := repo.KeyPlaceholderByQuery(len(configIds))
	sql := "delete from sys_config where config_id in (" + placeholder + ")"
	parameters := repo.ConvertIdsSlice(configIds)
	results, err := datasource.ExecDB("", sql, parameters)
	if err != nil {
		logger.Errorf("delete err => %v", err)
		return 0
	}
	return results
}
