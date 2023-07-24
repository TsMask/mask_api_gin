package repository

import (
	"fmt"
	"mask_api_gin/src/framework/datasource"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/utils/repo"
	"mask_api_gin/src/modules/system/model"
	"strings"
)

// SysDictTypeImpl 字典类型表 数据层处理
var SysDictTypeImpl = &sysDictTypeImpl{
	selectSql: `select 
	dict_id, dict_name, dict_type, status, create_by, create_time, remark 
	from sys_dict_type`,

	resultMap: map[string]string{
		"dict_id":     "DictID",
		"dict_name":   "DictName",
		"dict_type":   "DictType",
		"remark":      "Remark",
		"status":      "Status",
		"create_by":   "CreateBy",
		"create_time": "CreateTime",
		"update_by":   "UpdateBy",
		"update_time": "UpdateTime",
	},
}

type sysDictTypeImpl struct {
	// 查询视图对象SQL
	selectSql string
	// 结果字段与实体映射
	resultMap map[string]string
}

// convertResultRows 将结果记录转实体结果组
func (r *sysDictTypeImpl) convertResultRows(rows []map[string]interface{}) []model.SysDictType {
	arr := make([]model.SysDictType, 0)
	for _, row := range rows {
		sysDictType := model.SysDictType{}
		for key, value := range row {
			if keyMapper, ok := r.resultMap[key]; ok {
				repo.SetFieldValue(&sysDictType, keyMapper, value)
			}
		}
		arr = append(arr, sysDictType)
	}
	return arr
}

// SelectDictTypePage 根据条件分页查询字典类型
func (r *sysDictTypeImpl) SelectDictTypePage(query map[string]string) map[string]interface{} {
	return map[string]interface{}{}
}

// SelectDictTypeList 根据条件查询字典类型
func (r *sysDictTypeImpl) SelectDictTypeList(sysDictType model.SysDictType) []model.SysDictType {
	return []model.SysDictType{}
}

// SelectDictTypeByIDs 根据字典类型ID查询信息
func (r *sysDictTypeImpl) SelectDictTypeByIDs(dictIDs []string) []model.SysDictType {
	placeholder := repo.KeyPlaceholderByQuery(len(dictIDs))
	querySql := r.selectSql + " where dict_id in (" + placeholder + ")"
	parameters := repo.ConvertIdsSlice(dictIDs)
	results, err := datasource.RawDB("", querySql, parameters)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysDictType{}
	}
	// 转换实体
	return r.convertResultRows(results)
}

// SelectDictTypeByType 根据字典类型查询信息
func (r *sysDictTypeImpl) SelectDictTypeByType(dictType string) model.SysDictType {
	return model.SysDictType{}
}

// CheckUniqueDictType 校验字典是否唯一
func (r *sysDictTypeImpl) CheckUniqueDictType(sysDictType model.SysDictType) string {
	// 查询条件拼接
	var conditions []string
	var params []interface{}
	if sysDictType.DictName != "" {
		conditions = append(conditions, "dict_name = ?")
		params = append(params, sysDictType.DictName)
	}
	if sysDictType.DictType != "" {
		conditions = append(conditions, "dict_type = ?")
		params = append(params, sysDictType.DictType)
	}

	// 构建查询条件语句
	whereSql := ""
	if len(conditions) > 0 {
		whereSql += " where " + strings.Join(conditions, " and ")
	} else {
		return ""
	}

	// 查询数据
	querySql := "select dict_id as 'str' from sys_dict_type " + whereSql + " limit 1"
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

// InsertDictType 新增字典类型信息
func (r *sysDictTypeImpl) InsertDictType(sysDictType model.SysDictType) string {
	return ""
}

// UpdateDictType 修改字典类型信息
func (r *sysDictTypeImpl) UpdateDictType(sysDictType model.SysDictType) int {
	return 0
}

// DeleteDictTypeByIDs 批量删除字典类型信息
func (r *sysDictTypeImpl) DeleteDictTypeByIDs(dictIDs []string) int64 {
	placeholder := repo.KeyPlaceholderByQuery(len(dictIDs))
	sql := "delete from sys_dict_type where dict_id in (" + placeholder + ")"
	parameters := repo.ConvertIdsSlice(dictIDs)
	results, err := datasource.ExecDB("", sql, parameters)
	if err != nil {
		logger.Errorf("delete err => %v", err)
		return 0
	}
	return results
}
