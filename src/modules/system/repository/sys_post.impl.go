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
)

// SysPostImpl 岗位表 数据层处理
var SysPostImpl = &sysPostImpl{
	selectSql: `select 
	post_id, post_code, post_name, post_sort, status, create_by, create_time, remark 
	from sys_post`,

	resultMap: map[string]string{
		"post_id":     "PostID",
		"post_code":   "PostCode",
		"post_name":   "PostName",
		"post_sort":   "PostSort",
		"status":      "Status",
		"create_by":   "CreateBy",
		"create_time": "CreateTime",
		"update_by":   "UpdateBy",
		"update_time": "UpdateTime",
		"remark":      "Remark",
	},
}

type sysPostImpl struct {
	// 查询视图对象SQL
	selectSql string
	// 结果字段与实体映射
	resultMap map[string]string
}

// convertResultRows 将结果记录转实体结果组
func (r *sysPostImpl) convertResultRows(rows []map[string]interface{}) []model.SysPost {
	arr := make([]model.SysPost, 0)
	for _, row := range rows {
		sysPost := model.SysPost{}
		for key, value := range row {
			if keyMapper, ok := r.resultMap[key]; ok {
				repo.SetFieldValue(&sysPost, keyMapper, value)
			}
		}
		arr = append(arr, sysPost)
	}
	return arr
}

// SelectPostPage 查询岗位分页数据集合
func (r *sysPostImpl) SelectPostPage(query map[string]string) map[string]interface{} {
	// 查询条件拼接
	var conditions []string
	var params []interface{}
	if v, ok := query["postCode"]; ok {
		conditions = append(conditions, "post_code like concat(?, '%')")
		params = append(params, v)
	}
	if v, ok := query["postName"]; ok {
		conditions = append(conditions, "post_name like concat(?, '%')")
		params = append(params, v)
	}
	if v, ok := query["status"]; ok {
		conditions = append(conditions, "status = ?")
		params = append(params, v)
	}

	// 构建查询条件语句
	whereSql := ""
	if len(conditions) > 0 {
		whereSql += " where " + strings.Join(conditions, " and ")
	}

	// 查询数量 长度为0直接返回
	totalSql := "select count(1) as 'total' from sys_post"
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
	return map[string]interface{}{
		"total": total,
		"rows":  rows,
	}
}

// SelectPostList 查询岗位数据集合
func (r *sysPostImpl) SelectPostList(sysPost model.SysPost) []model.SysPost {
	// 查询条件拼接
	var conditions []string
	var params []interface{}
	if sysPost.PostCode != "" {
		conditions = append(conditions, "post_code like concat(?, '%')")
		params = append(params, sysPost.PostCode)
	}
	if sysPost.PostName != "" {
		conditions = append(conditions, "post_name like concat(?, '%')")
		params = append(params, sysPost.PostName)
	}
	if sysPost.Status != "" {
		conditions = append(conditions, "status = ?")
		params = append(params, sysPost.Status)
	}

	// 构建查询条件语句
	whereSql := ""
	if len(conditions) > 0 {
		whereSql += " where " + strings.Join(conditions, " and ")
	}

	// 查询数据
	orderSql := " order by post_sort"
	querySql := r.selectSql + whereSql + orderSql
	rows, err := datasource.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysPost{}
	}
	return r.convertResultRows(rows)
}

// SelectPostByIds 通过岗位ID查询岗位信息
func (r *sysPostImpl) SelectPostByIds(postIds []string) []model.SysPost {
	placeholder := repo.KeyPlaceholderByQuery(len(postIds))
	querySql := r.selectSql + " where post_id in (" + placeholder + ")"
	parameters := repo.ConvertIdsSlice(postIds)
	results, err := datasource.RawDB("", querySql, parameters)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysPost{}
	}
	// 转换实体
	return r.convertResultRows(results)
}

// SelectPostListByUserId 根据用户ID获取岗位选择框列表
func (r *sysPostImpl) SelectPostListByUserId(userId string) []model.SysPost {
	// 查询数据
	querySql := `select distinct 
	p.post_id, p.post_name, p.post_code 
	from sys_post p 
    left join sys_user_post up on up.post_id = p.post_id 
    left join sys_user u on u.user_id = up.user_id 
    where u.user_id = ? order by p.post_id`
	rows, err := datasource.RawDB("", querySql, []interface{}{userId})
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysPost{}
	}
	return r.convertResultRows(rows)
}

// DeletePostByIds 批量删除岗位信息
func (r *sysPostImpl) DeletePostByIds(postIds []string) int64 {
	placeholder := repo.KeyPlaceholderByQuery(len(postIds))
	sql := "delete from sys_post where post_id (" + placeholder + ")"
	parameters := repo.ConvertIdsSlice(postIds)
	results, err := datasource.ExecDB("", sql, parameters)
	if err != nil {
		logger.Errorf("delete err => %v", err)
		return 0
	}
	return results
}

// UpdatePost 修改岗位信息
func (r *sysPostImpl) UpdatePost(sysPost model.SysPost) int64 {
	// 参数拼接
	params := make(map[string]interface{})
	if sysPost.PostCode != "" {
		params["post_code"] = sysPost.PostCode
	}
	if sysPost.PostName != "" {
		params["post_name"] = sysPost.PostName
	}
	if sysPost.PostSort >= 0 {
		params["post_sort"] = sysPost.PostSort
	}
	if sysPost.Status != "" {
		params["status"] = sysPost.Status
	}
	if sysPost.Remark != "" {
		params["remark"] = sysPost.Remark
	}
	if sysPost.UpdateBy != "" {
		params["update_by"] = sysPost.UpdateBy
		params["update_time"] = date.NowTimestamp()
	}

	// 构建执行语句
	keys, values := repo.KeyValueByUpdate(params)
	sql := "update sys_post set " + strings.Join(keys, ",") + " where post_id = ?"

	// 执行更新
	values = append(values, sysPost.PostID)
	rows, err := datasource.ExecDB("", sql, values)
	if err != nil {
		logger.Errorf("update row : %v", err.Error())
		return 0
	}
	return rows
}

// InsertPost 新增岗位信息
func (r *sysPostImpl) InsertPost(sysPost model.SysPost) string {
	// 参数拼接
	params := make(map[string]interface{})
	if sysPost.PostID != "" {
		params["post_id"] = sysPost.PostID
	}
	if sysPost.PostCode != "" {
		params["post_code"] = sysPost.PostCode
	}
	if sysPost.PostName != "" {
		params["post_name"] = sysPost.PostName
	}
	if sysPost.PostSort >= 0 {
		params["post_sort"] = sysPost.PostSort
	}
	if sysPost.Status != "" {
		params["status"] = sysPost.Status
	}
	if sysPost.Remark != "" {
		params["remark"] = sysPost.Remark
	}
	if sysPost.CreateBy != "" {
		params["create_by"] = sysPost.CreateBy
		params["create_time"] = date.NowTimestamp()
	}

	// 构建执行语句
	keys, placeholder, values := repo.KeyPlaceholderValueByInsert(params)
	sql := "insert into sys_post (" + strings.Join(keys, ",") + ")values(" + placeholder + ")"

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

// CheckUniquePost 校验岗位唯一
func (r *sysPostImpl) CheckUniquePost(sysPost model.SysPost) string {
	// 查询条件拼接
	var conditions []string
	var params []interface{}
	if sysPost.PostName != "" {
		conditions = append(conditions, "post_name= ?")
		params = append(params, sysPost.PostName)
	}
	if sysPost.PostCode != "" {
		conditions = append(conditions, "post_code = ?")
		params = append(params, sysPost.PostCode)
	}

	// 构建查询条件语句
	whereSql := ""
	if len(conditions) > 0 {
		whereSql += " where " + strings.Join(conditions, " and ")
	} else {
		return ""
	}

	// 查询数据
	querySql := "select post_id as 'str' from sys_post " + whereSql + " limit 1"
	results, err := datasource.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err %v", err)
	}
	if len(results) > 0 {
		return fmt.Sprintf("%v", results[0]["str"])
	}
	return ""
}
