package repository

import (
	"mask_api_gin/src/framework/datasource"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/utils/date"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/framework/utils/repo"
	"mask_api_gin/src/modules/monitor/model"
	"strings"
)

// SysLogininforImpl 系统登录访问表 数据层处理
var SysLogininforImpl = &sysLogininforImpl{
	selectSql: `select info_id, user_name, ipaddr, login_location, 
	browser, os, status, msg, login_time from sys_logininfor`,

	resultMap: map[string]string{
		"info_id":        "InfoID",
		"user_name":      "UserName",
		"status":         "Status",
		"ipaddr":         "IPAddr",
		"login_location": "LoginLocation",
		"browser":        "Browser",
		"os":             "Os",
		"msg":            "Msg",
		"login_time":     "LoginTime",
	},
}

type sysLogininforImpl struct {
	// 查询视图对象SQL
	selectSql string
	// 结果字段与实体映射
	resultMap map[string]string
}

// convertResultRows 将结果记录转实体结果组
func (r *sysLogininforImpl) convertResultRows(rows []map[string]interface{}) []model.SysLogininfor {
	arr := make([]model.SysLogininfor, 0)
	for _, row := range rows {
		sysLogininfor := model.SysLogininfor{}
		for key, value := range row {
			if keyMapper, ok := r.resultMap[key]; ok {
				repo.SetFieldValue(&sysLogininfor, keyMapper, value)
			}
		}
		arr = append(arr, sysLogininfor)
	}
	return arr
}

// SelectLogininforPage 分页查询系统登录日志集合
func (r *sysLogininforImpl) SelectLogininforPage(query map[string]string) map[string]interface{} {
	// 查询条件拼接
	var conditions []string
	var params []interface{}
	if v, ok := query["ipaddr"]; ok {
		conditions = append(conditions, "ipaddr like concat(?, '%')")
		params = append(params, v)
	}
	if v, ok := query["userName"]; ok {
		conditions = append(conditions, "user_name like concat(?, '%')")
		params = append(params, v)
	}
	if v, ok := query["status"]; ok {
		conditions = append(conditions, "status = ?")
		params = append(params, v)
	}
	if v, ok := query["beginTime"]; ok {
		conditions = append(conditions, "login_time >= ?")
		beginDate := date.ParseStrToDate(v, date.YYYY_MM_DD)
		params = append(params, beginDate.UnixNano()/1e6)
	}
	if v, ok := query["endTime"]; ok {
		conditions = append(conditions, "login_time <= ?")
		endDate := date.ParseStrToDate(v, date.YYYY_MM_DD)
		params = append(params, endDate.UnixNano()/1e6)
	}

	// 构建查询条件语句
	whereSql := ""
	if len(conditions) > 0 {
		whereSql += " where " + strings.Join(conditions, " and ")
	}

	// 查询数量 长度为0直接返回
	totalSql := "select count(1) as 'total' from sys_logininfor"
	totalRows, err := datasource.RawDB("", totalSql+whereSql, params)
	if err != nil {
		logger.Errorf("total err => %v", err)
	}
	total := parse.Number(totalRows[0]["total"])
	if total <= 0 {
		return map[string]interface{}{
			"total": 0,
			"rows":  []interface{}{},
		}
	}

	// 分页
	pageNum, pageSize := repo.PageNumSize(query["pageNum"], query["pageSize"])
	pageSql := " order by info_id desc limit ?,? "
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
	return map[string]interface{}{
		"total": total,
		"rows":  rows,
	}
}

// SelectLogininforList 查询系统登录日志集合
func (r *sysLogininforImpl) SelectLogininforList(sysLogininfor model.SysLogininfor) []model.SysLogininfor {
	// 查询条件拼接
	var conditions []string
	var params []interface{}
	if sysLogininfor.IPAddr != "" {
		conditions = append(conditions, "title like concat(?, '%')")
		params = append(params, sysLogininfor.IPAddr)
	}
	if sysLogininfor.UserName != "" {
		conditions = append(conditions, "user_name like concat(?, '%')")
		params = append(params, sysLogininfor.UserName)
	}
	if sysLogininfor.Status != "" {
		conditions = append(conditions, "status = ?")
		params = append(params, sysLogininfor.Status)
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
		return []model.SysLogininfor{}
	}

	// 转换实体
	return r.convertResultRows(results)
}

// InsertLogininfor 新增系统登录日志
func (r *sysLogininforImpl) InsertLogininfor(sysLogininfor model.SysLogininfor) string {
	// 参数拼接
	params := make(map[string]interface{})
	params["login_time"] = date.NowTimestamp()
	if sysLogininfor.UserName != "" {
		params["user_name"] = sysLogininfor.UserName
	}
	if sysLogininfor.Status != "" {
		params["status"] = sysLogininfor.Status
	}
	if sysLogininfor.IPAddr != "" {
		params["ipaddr"] = sysLogininfor.IPAddr
	}
	if sysLogininfor.LoginLocation != "" {
		params["login_location"] = sysLogininfor.LoginLocation
	}
	if sysLogininfor.Browser != "" {
		params["browser"] = sysLogininfor.Browser
	}
	if sysLogininfor.OS != "" {
		params["os"] = sysLogininfor.OS
	}
	if sysLogininfor.Msg != "" {
		params["msg"] = sysLogininfor.Msg
	}

	// 构建执行语句
	keys, placeholder, values := repo.KeyPlaceholderValueByInsert(params)
	sql := "insert into sys_logininfor (" + strings.Join(keys, ",") + ")values(" + placeholder + ")"

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

// DeleteLogininforByIds 批量删除系统登录日志
func (r *sysLogininforImpl) DeleteLogininforByIds(infoIds []string) int64 {
	placeholder := repo.KeyPlaceholderByQuery(len(infoIds))
	sql := "delete from sys_logininfor where info_id in (" + placeholder + ")"
	parameters := repo.ConvertIdsSlice(infoIds)
	results, err := datasource.ExecDB("", sql, parameters)
	if err != nil {
		logger.Errorf("delete err => %v", err)
		return 0
	}
	return results
}

// CleanLogininfor 清空系统登录日志
func (r *sysLogininforImpl) CleanLogininfor() error {
	sql := "truncate table sys_logininfor"
	results, err := datasource.ExecDB("", sql, []interface{}{})
	logger.Errorf("delete results => %v", results)
	return err
}
