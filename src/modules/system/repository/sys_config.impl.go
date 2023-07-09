package repository

import (
	"mask_api_gin/src/framework/datasource"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/service/repo"
	"mask_api_gin/src/framework/utils/date"
	"mask_api_gin/src/modules/system/model"
	"strings"
)

// SysConfigImpl 参数配置表 数据层处理
var SysConfigImpl = new(sysConfigImpl)

// 查询视图对象SQL
var selectSql = "select config_id, config_name, config_key, config_value, config_type, create_by, create_time, update_by, update_time, remark from sys_config"

type sysConfigImpl struct{}

// SelectDictDataPage 分页查询参数配置列表数据
func (r *sysConfigImpl) SelectConfigPage(query map[string]string) map[string]interface{} {
	db := datasource.GetDefaultDB()

	// 查询条件拼接
	var conditions []string
	var params []interface{}
	if v, ok := query["configName"]; ok {
		conditions = append(conditions, "config_name like concat(?, '%')")
		params = append(params, v)
	}
	if v, ok := query["configType"]; ok {
		conditions = append(conditions, "config_type = ?")
		params = append(params, v)
	}
	if v, ok := query["configKey"]; ok {
		conditions = append(conditions, "config_key like concat(?, '%')")
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
	totalSql := "select count(1) as 'total' from sys_config"
	var total int64
	totalRes := db.Raw(totalSql+whereSql, params...).Scan(&total)
	if totalRes.Error != nil {
		logger.Errorf("totalRes err %v", totalRes.Error)
	}
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
	var sysConfig []model.SysConfig
	querySql := selectSql + whereSql + pageSql
	queryRes := db.Raw(querySql, params...).Scan(&sysConfig)
	if queryRes.Error != nil {
		logger.Errorf("queryRes err %v", queryRes.Error)
	}

	rows := repo.ConvertResultRows(sysConfig)
	return map[string]interface{}{
		"total": total,
		"rows":  rows,
	}
}

// SelectConfigList 查询参数配置列表
func (r *sysConfigImpl) SelectConfigList(sysConfig model.SysConfig) []model.SysConfig {
	// 实现具体逻辑
	return []model.SysConfig{}
}

// SelectConfigValueByKey 通过参数键名查询参数键值
func (r *sysConfigImpl) SelectConfigValueByKey(configKey string) string {
	// 实现具体逻辑
	return ""
}

// SelectConfigById 通过配置ID查询参数配置信息
func (r *sysConfigImpl) SelectConfigById(configId string) model.SysConfig {
	// 实现具体逻辑
	return model.SysConfig{}
}

// CheckUniqueConfigKey 校验参数键名是否唯一
func (r *sysConfigImpl) CheckUniqueConfigKey(configKey string) string {
	// 实现具体逻辑
	return ""
}

// InsertConfig 新增参数配置
func (r *sysConfigImpl) InsertConfig(sysConfig model.SysConfig) string {
	// 实现具体逻辑
	return ""
}

// UpdateConfig 修改参数配置
func (r *sysConfigImpl) UpdateConfig(sysConfig model.SysConfig) int {
	// 实现具体逻辑
	return 0
}

// DeleteConfigByIds 批量删除参数配置信息
func (r *sysConfigImpl) DeleteConfigByIds(configIds []string) int {
	// 实现具体逻辑
	return 0
}
