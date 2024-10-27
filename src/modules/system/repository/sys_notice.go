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

// NewSysNotice 实例化数据层
var NewSysNotice = &SysNotice{
	selectSql: `select 
	notice_id, notice_title, notice_type, notice_content, status, del_flag, 
	create_by, create_time, update_by, update_time, remark 
	from sys_notice`,

	resultMap: map[string]string{
		"notice_id":      "NoticeID",
		"notice_title":   "NoticeTitle",
		"notice_type":    "NoticeType",
		"notice_content": "NoticeContent",
		"status":         "Status",
		"del_flag":       "DelFlag",
		"create_by":      "CreateBy",
		"create_time":    "CreateTime",
		"update_by":      "UpdateBy",
		"update_time":    "UpdateTime",
		"remark":         "Remark",
	},
}

// SysNotice 通知公告表 数据层处理
type SysNotice struct {
	selectSql string            // 查询视图对象SQL
	resultMap map[string]string // 结果字段与实体映射
}

// SelectByPage 分页查询集合
func (r SysNotice) SelectByPage(query map[string]any) map[string]any {
	// 查询条件拼接
	var conditions []string
	var params []any
	if v, ok := query["noticeTitle"]; ok && v != "" {
		conditions = append(conditions, "notice_title like concat(?, '%')")
		params = append(params, v)
	}
	if v, ok := query["noticeType"]; ok && v != "" {
		conditions = append(conditions, "notice_type = ?")
		params = append(params, v)
	}
	if v, ok := query["createBy"]; ok && v != "" {
		conditions = append(conditions, "create_by like concat(?, '%')")
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
	whereSql := " where del_flag = '0' "
	if len(conditions) > 0 {
		whereSql += " and " + strings.Join(conditions, " and ")
	}

	// 查询结果
	result := map[string]any{
		"total": int64(0),
		"rows":  []model.SysNotice{},
	}

	// 查询数量 长度为0直接返回
	totalSql := "select count(1) as 'total' from sys_notice"
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
	rows, err := db.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err => %v", err)
	}

	// 转换实体
	result["rows"] = db.ConvertResultRows[model.SysNotice](model.SysNotice{}, r.resultMap, rows)
	return result
}

// Select 查询集合
func (r SysNotice) Select(sysNotice model.SysNotice) []model.SysNotice {
	// 查询条件拼接
	var conditions []string
	var params []any
	if sysNotice.NoticeTitle != "" {
		conditions = append(conditions, "notice_title like concat(?, '%')")
		params = append(params, sysNotice.NoticeTitle)
	}
	if sysNotice.NoticeType != "" {
		conditions = append(conditions, "notice_type = ?")
		params = append(params, sysNotice.NoticeType)
	}
	if sysNotice.CreateBy != "" {
		conditions = append(conditions, "create_by like concat(?, '%')")
		params = append(params, sysNotice.CreateBy)
	}
	if sysNotice.Status != "" {
		conditions = append(conditions, "status = ?")
		params = append(params, sysNotice.Status)
	}

	// 构建查询条件语句
	whereSql := " where del_flag = '0' "
	if len(conditions) > 0 {
		whereSql += " and " + strings.Join(conditions, " and ")
	}

	// 查询数据
	querySql := r.selectSql + whereSql
	rows, err := db.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysNotice{}
	}

	// 转换实体
	return db.ConvertResultRows[model.SysNotice](model.SysNotice{}, r.resultMap, rows)
}

// SelectByIds 通过ID查询信息
func (r SysNotice) SelectByIds(noticeIds []string) []model.SysNotice {
	placeholder := db.KeyPlaceholderByQuery(len(noticeIds))
	querySql := r.selectSql + " where notice_id in (" + placeholder + ")"
	parameters := db.ConvertIdsSlice(noticeIds)
	rows, err := db.RawDB("", querySql, parameters)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysNotice{}
	}
	// 转换实体
	return db.ConvertResultRows[model.SysNotice](model.SysNotice{}, r.resultMap, rows)
}

// Insert 新增信息
func (r SysNotice) Insert(sysNotice model.SysNotice) string {
	// 参数拼接
	params := make(map[string]any)
	if sysNotice.NoticeTitle != "" {
		params["notice_title"] = sysNotice.NoticeTitle
	}
	if sysNotice.NoticeType != "" {
		params["notice_type"] = sysNotice.NoticeType
	}
	if sysNotice.NoticeContent != "" {
		params["notice_content"] = sysNotice.NoticeContent
	}
	if sysNotice.Status != "" {
		params["status"] = sysNotice.Status
	}
	if sysNotice.Remark != "" {
		params["remark"] = sysNotice.Remark
	}
	if sysNotice.CreateBy != "" {
		params["create_by"] = sysNotice.CreateBy
		params["create_time"] = time.Now().UnixMilli()
	}

	// 构建执行语句
	keys, values, placeholder := db.KeyValuePlaceholderByInsert(params)
	sql := fmt.Sprintf("insert into sys_notice (%s)values(%s)", keys, placeholder)

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
func (r SysNotice) Update(sysNotice model.SysNotice) int64 {
	// 参数拼接
	params := make(map[string]any)
	if sysNotice.NoticeTitle != "" {
		params["notice_title"] = sysNotice.NoticeTitle
	}
	if sysNotice.NoticeType != "" {
		params["notice_type"] = sysNotice.NoticeType
	}
	if sysNotice.NoticeContent != "" {
		params["notice_content"] = sysNotice.NoticeContent
	}
	if sysNotice.Status != "" {
		params["status"] = sysNotice.Status
	}
	params["remark"] = sysNotice.Remark
	if sysNotice.UpdateBy != "" {
		params["update_by"] = sysNotice.UpdateBy
		params["update_time"] = time.Now().UnixMilli()
	}

	// 构建执行语句
	keys, values := db.KeyValueByUpdate(params)
	sql := fmt.Sprintf("update sys_notice set %s where notice_id = ?", keys)

	// 执行更新
	values = append(values, sysNotice.NoticeId)
	rows, err := db.ExecDB("", sql, values)
	if err != nil {
		logger.Errorf("update row : %v", err.Error())
		return 0
	}
	return rows
}

// DeleteByIds 批量删除信息
func (r SysNotice) DeleteByIds(noticeIds []string) int64 {
	placeholder := db.KeyPlaceholderByQuery(len(noticeIds))
	sql := fmt.Sprintf("update sys_notice set del_flag = '1' where notice_id in (%s)", placeholder)
	parameters := db.ConvertIdsSlice(noticeIds)
	results, err := db.ExecDB("", sql, parameters)
	if err != nil {
		logger.Errorf("update err => %v", err)
		return 0
	}
	return results
}
