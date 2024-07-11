package repository

import (
	"fmt"
	db "mask_api_gin/src/framework/data_source"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/modules/system/model"
	"strings"
	"time"
)

// NewSysDictData 实例化数据层
var NewSysDictData = &SysDictDataRepository{
	selectSql: `select 
	dict_code, dict_sort, dict_label, dict_value, dict_type, tag_class, tag_type, status, create_by, create_time, remark 
	from sys_dict_data`,

	resultMap: map[string]string{
		"dict_code":   "DictCode",
		"dict_sort":   "DictSort",
		"dict_label":  "DictLabel",
		"dict_value":  "DictValue",
		"dict_type":   "DictType",
		"tag_class":   "TagClass",
		"tag_type":    "TagType",
		"status":      "Status",
		"remark":      "Remark",
		"create_by":   "CreateBy",
		"create_time": "CreateTime",
		"update_by":   "UpdateBy",
		"update_time": "UpdateTime",
	},
}

// SysDictDataRepository 字典类型数据表 数据层处理
type SysDictDataRepository struct {
	selectSql string            // 查询视图对象SQL
	resultMap map[string]string // 结果字段与实体映射
}

// convertResultRows 将结果记录转实体结果组
func (r *SysDictDataRepository) convertResultRows(rows []map[string]any) []model.SysDictData {
	arr := make([]model.SysDictData, 0)
	for _, row := range rows {
		sysDictData := model.SysDictData{}
		for key, value := range row {
			if keyMapper, ok := r.resultMap[key]; ok {
				db.SetFieldValue(&sysDictData, keyMapper, value)
			}
		}
		arr = append(arr, sysDictData)
	}
	return arr
}

// SelectByPage 分页查询集合
func (r *SysDictDataRepository) SelectByPage(query map[string]any) map[string]any {
	// 查询条件拼接
	var conditions []string
	var params []any
	if v, ok := query["dictType"]; ok && v != "" {
		conditions = append(conditions, "dict_type = ?")
		params = append(params, v)
	}
	if v, ok := query["dictLabel"]; ok && v != "" {
		conditions = append(conditions, "dict_label like concat(?, '%')")
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
		"rows":  []model.SysDictData{},
	}

	// 查询数量 长度为0直接返回
	totalSql := "select count(1) as 'total' from sys_dict_data"
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
	pageSql := " order by dict_sort asc limit ?,? "
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
func (r *SysDictDataRepository) Select(sysDictData model.SysDictData) []model.SysDictData {
	// 查询条件拼接
	var conditions []string
	var params []any
	if sysDictData.DictLabel != "" {
		conditions = append(conditions, "dict_label like concat(?, '%')")
		params = append(params, sysDictData.DictLabel)
	}
	if sysDictData.DictType != "" {
		conditions = append(conditions, "dict_type = ?")
		params = append(params, sysDictData.DictType)
	}
	if sysDictData.Status != "" {
		conditions = append(conditions, "status = ?")
		params = append(params, sysDictData.Status)
	}

	// 构建查询条件语句
	whereSql := ""
	if len(conditions) > 0 {
		whereSql += " where " + strings.Join(conditions, " and ")
	}

	// 查询数据
	orderSql := " order by dict_sort asc "
	querySql := r.selectSql + whereSql + orderSql
	results, err := db.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysDictData{}
	}

	// 转换实体
	return r.convertResultRows(results)
}

// SelectByCodes 通过Code查询信息
func (r *SysDictDataRepository) SelectByCodes(dictCodes []string) []model.SysDictData {
	placeholder := db.KeyPlaceholderByQuery(len(dictCodes))
	querySql := r.selectSql + " where dict_code in (" + placeholder + ")"
	parameters := db.ConvertIdsSlice(dictCodes)
	results, err := db.RawDB("", querySql, parameters)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysDictData{}
	}
	// 转换实体
	return r.convertResultRows(results)
}

// Insert 新增信息
func (r *SysDictDataRepository) Insert(sysDictData model.SysDictData) string {
	// 参数拼接
	params := make(map[string]any)
	if sysDictData.DictSort > 0 {
		params["dict_sort"] = sysDictData.DictSort
	}
	if sysDictData.DictLabel != "" {
		params["dict_label"] = sysDictData.DictLabel
	}
	if sysDictData.DictValue != "" {
		params["dict_value"] = sysDictData.DictValue
	}
	if sysDictData.DictType != "" {
		params["dict_type"] = sysDictData.DictType
	}
	if sysDictData.TagClass != "" {
		params["tag_class"] = sysDictData.TagClass
	}
	if sysDictData.TagType != "" {
		params["tag_type"] = sysDictData.TagType
	}
	if sysDictData.Status != "" {
		params["status"] = sysDictData.Status
	}
	if sysDictData.Remark != "" {
		params["remark"] = sysDictData.Remark
	}
	if sysDictData.CreateBy != "" {
		params["create_by"] = sysDictData.CreateBy
		params["create_time"] = time.Now().UnixMilli()
	}

	// 构建执行语句
	keys, values, placeholder := db.KeyValuePlaceholderByInsert(params)
	sql := fmt.Sprintf("insert into sys_dict_data (%s)values(%s)", keys, placeholder)

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
func (r *SysDictDataRepository) Update(sysDictData model.SysDictData) int64 {
	// 参数拼接
	params := make(map[string]any)
	if sysDictData.DictSort >= 0 {
		params["dict_sort"] = sysDictData.DictSort
	}
	if sysDictData.DictLabel != "" {
		params["dict_label"] = sysDictData.DictLabel
	}
	if sysDictData.DictValue != "" {
		params["dict_value"] = sysDictData.DictValue
	}
	if sysDictData.DictType != "" {
		params["dict_type"] = sysDictData.DictType
	}
	params["tag_class"] = sysDictData.TagClass
	params["tag_type"] = sysDictData.TagType
	if sysDictData.Status != "" {
		params["status"] = sysDictData.Status
	}
	params["remark"] = sysDictData.Remark
	if sysDictData.UpdateBy != "" {
		params["update_by"] = sysDictData.UpdateBy
		params["update_time"] = time.Now().UnixMilli()
	}

	// 构建执行语句
	keys, values := db.KeyValueByUpdate(params)
	sql := fmt.Sprintf("update sys_dict_data set %s where dict_code = ?", keys)

	// 执行更新
	values = append(values, sysDictData.DictCode)
	rows, err := db.ExecDB("", sql, values)
	if err != nil {
		logger.Errorf("update row : %v", err.Error())
		return 0
	}
	return rows
}

// DeleteByCodes 批量删除信息
func (r *SysDictDataRepository) DeleteByCodes(dictCodes []string) int64 {
	placeholder := db.KeyPlaceholderByQuery(len(dictCodes))
	sql := fmt.Sprintf("delete from sys_dict_data where dict_code in (%s)", placeholder)
	parameters := db.ConvertIdsSlice(dictCodes)
	results, err := db.ExecDB("", sql, parameters)
	if err != nil {
		logger.Errorf("delete err => %v", err)
		return 0
	}
	return results
}

// CheckUnique 检查信息是否唯一
func (r *SysDictDataRepository) CheckUnique(sysDictData model.SysDictData) string {
	// 查询条件拼接
	var conditions []string
	var params []any
	if sysDictData.DictType != "" {
		conditions = append(conditions, "dict_type = ?")
		params = append(params, sysDictData.DictType)
	}
	if sysDictData.DictLabel != "" {
		conditions = append(conditions, "dict_label = ?")
		params = append(params, sysDictData.DictLabel)
	}
	if sysDictData.DictValue != "" {
		conditions = append(conditions, "dict_value = ?")
		params = append(params, sysDictData.DictValue)
	}

	// 构建查询条件语句
	whereSql := ""
	if len(conditions) > 0 {
		whereSql += " where " + strings.Join(conditions, " and ")
	} else {
		return "-"
	}

	// 查询数据
	querySql := fmt.Sprintf("select dict_code as 'str' from sys_dict_data %s limit 1", whereSql)
	results, err := db.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err %v", err)
		return "-"
	}
	if len(results) > 0 {
		return fmt.Sprint(results[0]["str"])
	}
	return "-"
}

// ExistDataByDictType 存在数据数量
func (r *SysDictDataRepository) ExistDataByDictType(dictType string) int64 {
	querySql := "select count(1) as 'total' from sys_dict_data where dict_type = ?"
	results, err := db.RawDB("", querySql, []any{dictType})
	if err != nil {
		logger.Errorf("query err => %v", err)
		return 0
	}
	if len(results) > 0 {
		return parse.Number(results[0]["total"])
	}
	return 0
}

// UpdateDataByDictType 更新一组字典类型
func (r *SysDictDataRepository) UpdateDataByDictType(oldDictType string, newDictType string) int64 {
	// 参数拼接
	params := make([]any, 0)
	if oldDictType == "" || newDictType == "" {
		return 0
	}
	params = append(params, newDictType)
	params = append(params, oldDictType)

	// 构建执行语句
	sql := "update sys_dict_data set dict_type = ? where dict_type = ?"

	// 执行更新
	rows, err := db.ExecDB("", sql, params)
	if err != nil {
		logger.Errorf("update row : %v", err.Error())
		return 0
	}
	return rows
}
