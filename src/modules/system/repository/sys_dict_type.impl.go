package repository

import (
	"mask_api_gin/src/framework/datasource"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/utils/date"
	"mask_api_gin/src/framework/utils/parse"

	"mask_api_gin/src/modules/system/model"
	"strings"
	"time"
)

// 实例化数据层 SysDictTypeImpl 结构体
var NewSysDictTypeImpl = &SysDictTypeImpl{
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

// SysDictTypeImpl 字典类型表 数据层处理
type SysDictTypeImpl struct {
	// 查询视图对象SQL
	selectSql string
	// 结果字段与实体映射
	resultMap map[string]string
}

// convertResultRows 将结果记录转实体结果组
func (r *SysDictTypeImpl) convertResultRows(rows []map[string]any) []model.SysDictType {
	arr := make([]model.SysDictType, 0)
	for _, row := range rows {
		sysDictType := model.SysDictType{}
		for key, value := range row {
			if keyMapper, ok := r.resultMap[key]; ok {
				datasource.SetFieldValue(&sysDictType, keyMapper, value)
			}
		}
		arr = append(arr, sysDictType)
	}
	return arr
}

// SelectDictTypePage 根据条件分页查询字典类型
func (r *SysDictTypeImpl) SelectDictTypePage(query map[string]any) map[string]any {
	// 查询条件拼接
	var conditions []string
	var params []any
	if v, ok := query["dictName"]; ok && v != "" {
		conditions = append(conditions, "dict_name like concat(?, '%')")
		params = append(params, v)
	}
	if v, ok := query["dictType"]; ok && v != "" {
		conditions = append(conditions, "dict_type like concat(?, '%')")
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
		"rows":  []model.SysDictType{},
	}

	// 查询数量 长度为0直接返回
	totalSql := "select count(1) as 'total' from sys_dict_type"
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
	pageNum, pageSize := datasource.PageNumSize(query["pageNum"], query["pageSize"])
	pageSql := " limit ?,? "
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

// SelectDictTypeList 根据条件查询字典类型
func (r *SysDictTypeImpl) SelectDictTypeList(sysDictType model.SysDictType) []model.SysDictType {
	// 查询条件拼接
	var conditions []string
	var params []any
	if sysDictType.DictName != "" {
		conditions = append(conditions, "dict_name like concat(?, '%')")
		params = append(params, sysDictType.DictName)
	}
	if sysDictType.DictType != "" {
		conditions = append(conditions, "dict_type like concat(?, '%')")
		params = append(params, sysDictType.DictType)
	}
	if sysDictType.Status != "" {
		conditions = append(conditions, "status = ?")
		params = append(params, sysDictType.Status)
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
		return []model.SysDictType{}
	}

	// 转换实体
	return r.convertResultRows(results)
}

// SelectDictTypeByIDs 根据字典类型ID查询信息
func (r *SysDictTypeImpl) SelectDictTypeByIDs(dictIDs []string) []model.SysDictType {
	placeholder := datasource.KeyPlaceholderByQuery(len(dictIDs))
	querySql := r.selectSql + " where dict_id in (" + placeholder + ")"
	parameters := datasource.ConvertIdsSlice(dictIDs)
	results, err := datasource.RawDB("", querySql, parameters)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysDictType{}
	}
	// 转换实体
	return r.convertResultRows(results)
}

// SelectDictTypeByType 根据字典类型查询信息
func (r *SysDictTypeImpl) SelectDictTypeByType(dictType string) model.SysDictType {
	querySql := r.selectSql + " where dict_type = ?"
	results, err := datasource.RawDB("", querySql, []any{dictType})
	if err != nil {
		logger.Errorf("query err => %v", err)
		return model.SysDictType{}
	}
	// 转换实体
	rows := r.convertResultRows(results)
	if len(rows) > 0 {
		return rows[0]
	}
	return model.SysDictType{}
}

// CheckUniqueDictType 校验字典是否唯一
func (r *SysDictTypeImpl) CheckUniqueDictType(sysDictType model.SysDictType) string {
	// 查询条件拼接
	var conditions []string
	var params []any
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
		v, ok := results[0]["str"].(string)
		if ok {
			return v
		}
		return ""
	}
	return ""
}

// InsertDictType 新增字典类型信息
func (r *SysDictTypeImpl) InsertDictType(sysDictType model.SysDictType) string {
	// 参数拼接
	params := make(map[string]any)
	if sysDictType.DictName != "" {
		params["dict_name"] = sysDictType.DictName
	}
	if sysDictType.DictType != "" {
		params["dict_type"] = sysDictType.DictType
	}
	if sysDictType.Status != "" {
		params["status"] = sysDictType.Status
	}
	if sysDictType.Remark != "" {
		params["remark"] = sysDictType.Remark
	}
	if sysDictType.CreateBy != "" {
		params["create_by"] = sysDictType.CreateBy
		params["create_time"] = time.Now().UnixMilli()
	}

	// 构建执行语句
	keys, values, placeholder := datasource.KeyValuePlaceholderByInsert(params)
	sql := "insert into sys_dict_type (" + keys + ")values(" + placeholder + ")"

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

// UpdateDictType 修改字典类型信息
func (r *SysDictTypeImpl) UpdateDictType(sysDictType model.SysDictType) int64 {
	// 参数拼接
	params := make(map[string]any)
	if sysDictType.DictName != "" {
		params["dict_name"] = sysDictType.DictName
	}
	if sysDictType.DictType != "" {
		params["dict_type"] = sysDictType.DictType
	}
	if sysDictType.Status != "" {
		params["status"] = sysDictType.Status
	}
	if sysDictType.Remark != "" {
		params["remark"] = sysDictType.Remark
	}
	if sysDictType.UpdateBy != "" {
		params["update_by"] = sysDictType.UpdateBy
		params["update_time"] = time.Now().UnixMilli()
	}

	// 构建执行语句
	keys, values := datasource.KeyValueByUpdate(params)
	sql := "update sys_dict_type set " + keys + " where dict_id = ?"

	// 执行更新
	values = append(values, sysDictType.DictID)
	rows, err := datasource.ExecDB("", sql, values)
	if err != nil {
		logger.Errorf("update row : %v", err.Error())
		return 0
	}
	return rows
}

// DeleteDictTypeByIDs 批量删除字典类型信息
func (r *SysDictTypeImpl) DeleteDictTypeByIDs(dictIDs []string) int64 {
	placeholder := datasource.KeyPlaceholderByQuery(len(dictIDs))
	sql := "delete from sys_dict_type where dict_id in (" + placeholder + ")"
	parameters := datasource.ConvertIdsSlice(dictIDs)
	results, err := datasource.ExecDB("", sql, parameters)
	if err != nil {
		logger.Errorf("delete err => %v", err)
		return 0
	}
	return results
}
