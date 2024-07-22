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

// NewSysDictType 实例化数据层
var NewSysDictType = &SysDictTypeRepository{
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

// SysDictTypeRepository 字典类型表 数据层处理
type SysDictTypeRepository struct {
	selectSql string            // 查询视图对象SQL
	resultMap map[string]string // 结果字段与实体映射
}

// convertResultRows 将结果记录转实体结果组
func (r *SysDictTypeRepository) convertResultRows(rows []map[string]any) []model.SysDictType {
	arr := make([]model.SysDictType, 0)
	for _, row := range rows {
		sysDictType := model.SysDictType{}
		for key, value := range row {
			if keyMapper, ok := r.resultMap[key]; ok {
				db.SetFieldValue(&sysDictType, keyMapper, value)
			}
		}
		arr = append(arr, sysDictType)
	}
	return arr
}

// SelectByPage 分页查询集合
func (r *SysDictTypeRepository) SelectByPage(query map[string]any) map[string]any {
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
		"total": int64(0),
		"rows":  []model.SysDictType{},
	}

	// 查询数量 长度为0直接返回
	totalSql := "select count(1) as 'total' from sys_dict_type"
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
func (r *SysDictTypeRepository) Select(sysDictType model.SysDictType) []model.SysDictType {
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
	results, err := db.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysDictType{}
	}

	// 转换实体
	return r.convertResultRows(results)
}

// SelectByIds 通过ID查询信息
func (r *SysDictTypeRepository) SelectByIds(dictIds []string) []model.SysDictType {
	placeholder := db.KeyPlaceholderByQuery(len(dictIds))
	querySql := r.selectSql + " where dict_id in (" + placeholder + ")"
	parameters := db.ConvertIdsSlice(dictIds)
	results, err := db.RawDB("", querySql, parameters)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysDictType{}
	}
	// 转换实体
	return r.convertResultRows(results)
}

// Insert 新增信息
func (r *SysDictTypeRepository) Insert(sysDictType model.SysDictType) string {
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
	keys, values, placeholder := db.KeyValuePlaceholderByInsert(params)
	sql := "insert into sys_dict_type (" + keys + ")values(" + placeholder + ")"

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
func (r *SysDictTypeRepository) Update(sysDictType model.SysDictType) int64 {
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
	params["remark"] = sysDictType.Remark
	if sysDictType.UpdateBy != "" {
		params["update_by"] = sysDictType.UpdateBy
		params["update_time"] = time.Now().UnixMilli()
	}

	// 构建执行语句
	keys, values := db.KeyValueByUpdate(params)
	sql := fmt.Sprintf("update sys_dict_type set %s where dict_id = ?", keys)

	// 执行更新
	values = append(values, sysDictType.DictID)
	rows, err := db.ExecDB("", sql, values)
	if err != nil {
		logger.Errorf("update row : %v", err.Error())
		return 0
	}
	return rows
}

// DeleteByIds 批量删除信息
func (r *SysDictTypeRepository) DeleteByIds(dictIds []string) int64 {
	placeholder := db.KeyPlaceholderByQuery(len(dictIds))
	sql := fmt.Sprintf("delete from sys_dict_type where dict_id in (%s)", placeholder)
	parameters := db.ConvertIdsSlice(dictIds)
	results, err := db.ExecDB("", sql, parameters)
	if err != nil {
		logger.Errorf("delete err => %v", err)
		return 0
	}
	return results
}

// CheckUnique 检查信息是否唯一
func (r *SysDictTypeRepository) CheckUnique(sysDictType model.SysDictType) string {
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
		return "-"
	}

	// 查询数据
	querySql := fmt.Sprintf("select dict_id as 'str' from sys_dict_type %s limit 1", whereSql)
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

// SelectByType 通过字典类型查询信息
func (r *SysDictTypeRepository) SelectByType(dictType string) model.SysDictType {
	querySql := r.selectSql + " where dict_type = ?"
	results, err := db.RawDB("", querySql, []any{dictType})
	if err != nil {
		logger.Errorf("query err => %v", err)
		return model.SysDictType{}
	}
	// 转换实体
	if rows := r.convertResultRows(results); len(rows) > 0 {
		return rows[0]
	}
	return model.SysDictType{}
}
