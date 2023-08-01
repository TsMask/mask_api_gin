package repository

import (
	"fmt"
	"mask_api_gin/src/framework/datasource"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/framework/utils/repo"
	"mask_api_gin/src/modules/system/model"
	"strings"
	"time"
)

// 实例化数据层 SysDictDataImpl 结构体
var NewSysDictDataImpl = &SysDictDataImpl{
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

// SysDictDataImpl 字典类型数据表 数据层处理
type SysDictDataImpl struct {
	// 查询视图对象SQL
	selectSql string
	// 结果字段与实体映射
	resultMap map[string]string
}

// convertResultRows 将结果记录转实体结果组
func (r *SysDictDataImpl) convertResultRows(rows []map[string]any) []model.SysDictData {
	arr := make([]model.SysDictData, 0)
	for _, row := range rows {
		sysDictData := model.SysDictData{}
		for key, value := range row {
			if keyMapper, ok := r.resultMap[key]; ok {
				repo.SetFieldValue(&sysDictData, keyMapper, value)
			}
		}
		arr = append(arr, sysDictData)
	}
	return arr
}

// SelectDictDataPage 根据条件分页查询字典数据
func (r *SysDictDataImpl) SelectDictDataPage(query map[string]any) map[string]any {
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

	// 查询数量 长度为0直接返回
	totalSql := "select count(1) as 'total' from sys_dict_data"
	totalRows, err := datasource.RawDB("", totalSql+whereSql, params)
	if err != nil {
		logger.Errorf("total err => %v", err)
	}
	total := parse.Number(totalRows[0]["total"])
	if total == 0 {
		return map[string]any{
			"total": total,
			"rows":  []model.SysDictData{},
		}
	}

	// 分页
	pageNum, pageSize := repo.PageNumSize(query["pageNum"], query["pageSize"])
	pageSql := " order by dict_sort asc limit ?,? "
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

// SelectDictDataList 根据条件查询字典数据
func (r *SysDictDataImpl) SelectDictDataList(sysDictData model.SysDictData) []model.SysDictData {
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
	results, err := datasource.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysDictData{}
	}

	// 转换实体
	return r.convertResultRows(results)
}

// SelectDictDataByCodes 根据字典数据编码查询信息
func (r *SysDictDataImpl) SelectDictDataByCodes(dictCodes []string) []model.SysDictData {
	placeholder := repo.KeyPlaceholderByQuery(len(dictCodes))
	querySql := r.selectSql + " where dict_code in (" + placeholder + ")"
	parameters := repo.ConvertIdsSlice(dictCodes)
	results, err := datasource.RawDB("", querySql, parameters)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysDictData{}
	}
	// 转换实体
	return r.convertResultRows(results)
}

// CountDictDataByType 查询字典数据
func (r *SysDictDataImpl) CountDictDataByType(dictType string) int64 {
	querySql := "select count(1) as 'total' from sys_dict_data where dict_type = ?"
	results, err := datasource.RawDB("", querySql, []any{dictType})
	if err != nil {
		logger.Errorf("query err => %v", err)
		return 0
	}
	if len(results) > 0 {
		return parse.Number(results[0]["total"])
	}
	return 0
}

// CheckUniqueDictData 校验字典数据是否唯一
func (r *SysDictDataImpl) CheckUniqueDictData(sysDictData model.SysDictData) string {
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
		return ""
	}

	// 查询数据
	querySql := "select dict_code as 'str' from sys_dict_data " + whereSql + " limit 1"
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

// DeleteDictDataByCodes 批量删除字典数据信息
func (r *SysDictDataImpl) DeleteDictDataByCodes(dictCodes []string) int64 {
	placeholder := repo.KeyPlaceholderByQuery(len(dictCodes))
	sql := "delete from sys_dict_data where dict_code in (" + placeholder + ")"
	parameters := repo.ConvertIdsSlice(dictCodes)
	results, err := datasource.ExecDB("", sql, parameters)
	if err != nil {
		logger.Errorf("delete err => %v", err)
		return 0
	}
	return results
}

// InsertDictData 新增字典数据信息
func (r *SysDictDataImpl) InsertDictData(sysDictData model.SysDictData) string {
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
	keys, placeholder, values := repo.KeyPlaceholderValueByInsert(params)
	sql := "insert into sys_dict_data (" + strings.Join(keys, ",") + ")values(" + placeholder + ")"

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

// UpdateDictData 修改字典数据信息
func (r *SysDictDataImpl) UpdateDictData(sysDictData model.SysDictData) int64 {
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
	if sysDictData.UpdateBy != "" {
		params["update_by"] = sysDictData.UpdateBy
		params["update_time"] = time.Now().UnixMilli()
	}

	// 构建执行语句
	keys, values := repo.KeyValueByUpdate(params)
	sql := "update sys_dict_data set " + strings.Join(keys, ",") + " where dict_code = ?"

	// 执行更新
	values = append(values, sysDictData.DictCode)
	rows, err := datasource.ExecDB("", sql, values)
	if err != nil {
		logger.Errorf("update row : %v", err.Error())
		return 0
	}
	return rows
}

// UpdateDictDataType 同步修改字典类型
func (r *SysDictDataImpl) UpdateDictDataType(oldDictType string, newDictType string) int64 {
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
	rows, err := datasource.ExecDB("", sql, params)
	if err != nil {
		logger.Errorf("update row : %v", err.Error())
		return 0
	}
	return rows
}
