package repository

import (
	"mask_api_gin/src/framework/datasource"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/utils/date"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/framework/utils/repo"
	"mask_api_gin/src/modules/system/model"
	"strings"
	"time"
)

// 实例化数据层 SysNoticeImpl 结构体
var NewSysNoticeImpl = &SysNoticeImpl{
	selectSql: `select 
	notice_id, notice_title, notice_type, notice_content, status, del_flag, 
	create_by, create_time, update_by, update_time, remark from sys_notice`,

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

// SysNoticeImpl 通知公告表 数据层处理
type SysNoticeImpl struct {
	// 查询视图对象SQL
	selectSql string
	// 结果字段与实体映射
	resultMap map[string]string
}

// convertResultRows 将结果记录转实体结果组
func (r *SysNoticeImpl) convertResultRows(rows []map[string]any) []model.SysNotice {
	arr := make([]model.SysNotice, 0)
	for _, row := range rows {
		sysNotice := model.SysNotice{}
		for key, value := range row {
			if keyMapper, ok := r.resultMap[key]; ok {
				repo.SetFieldValue(&sysNotice, keyMapper, value)
			}
		}
		arr = append(arr, sysNotice)
	}
	return arr
}

// SelectNoticePage 分页查询公告列表
func (r *SysNoticeImpl) SelectNoticePage(query map[string]any) map[string]any {
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
		"total": 0,
		"rows":  []model.SysNotice{},
	}

	// 查询数量 长度为0直接返回
	totalSql := "select count(1) as 'total' from sys_notice"
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
	result["rows"] = r.convertResultRows(results)
	return result
}

// SelectNoticeList 查询公告列表
func (r *SysNoticeImpl) SelectNoticeList(sysNotice model.SysNotice) []model.SysNotice {
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
	results, err := datasource.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysNotice{}
	}

	// 转换实体
	return r.convertResultRows(results)
}

// SelectNoticeByIds 查询公告信息
func (r *SysNoticeImpl) SelectNoticeByIds(noticeIds []string) []model.SysNotice {
	placeholder := repo.KeyPlaceholderByQuery(len(noticeIds))
	querySql := r.selectSql + " where notice_id in (" + placeholder + ")"
	parameters := repo.ConvertIdsSlice(noticeIds)
	results, err := datasource.RawDB("", querySql, parameters)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysNotice{}
	}
	// 转换实体
	return r.convertResultRows(results)
}

// InsertNotice 新增公告
func (r *SysNoticeImpl) InsertNotice(sysNotice model.SysNotice) string {
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
	keys, placeholder, values := repo.KeyPlaceholderValueByInsert(params)
	sql := "insert into sys_notice (" + strings.Join(keys, ",") + ")values(" + placeholder + ")"

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

// UpdateNotice 修改公告
func (r *SysNoticeImpl) UpdateNotice(sysNotice model.SysNotice) int64 {
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
	if sysNotice.UpdateBy != "" {
		params["update_by"] = sysNotice.UpdateBy
		params["update_time"] = time.Now().UnixMilli()
	}

	// 构建执行语句
	keys, values := repo.KeyValueByUpdate(params)
	sql := "update sys_notice set " + strings.Join(keys, ",") + " where notice_id = ?"

	// 执行更新
	values = append(values, sysNotice.NoticeID)
	rows, err := datasource.ExecDB("", sql, values)
	if err != nil {
		logger.Errorf("update row : %v", err.Error())
		return 0
	}
	return rows
}

// DeleteNoticeByIds 批量删除公告信息
func (r *SysNoticeImpl) DeleteNoticeByIds(noticeIds []string) int64 {
	placeholder := repo.KeyPlaceholderByQuery(len(noticeIds))
	sql := "update sys_notice set del_flag = '1' where notice_id in (" + placeholder + ")"
	parameters := repo.ConvertIdsSlice(noticeIds)
	results, err := datasource.ExecDB("", sql, parameters)
	if err != nil {
		logger.Errorf("update err => %v", err)
		return 0
	}
	return results
}
