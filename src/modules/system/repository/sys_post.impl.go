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

// NewSysPost 实例化数据层
var NewSysPost = &SysPostRepository{
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

// SysPostRepository 岗位表 数据层处理
type SysPostRepository struct {
	selectSql string            // 查询视图对象SQL
	resultMap map[string]string // 结果字段与实体映射
}

// SelectByPage 分页查询集合
func (r *SysPostRepository) SelectByPage(query map[string]any) map[string]any {
	// 查询条件拼接
	var conditions []string
	var params []any
	if v, ok := query["postCode"]; ok && v != "" {
		conditions = append(conditions, "post_code like concat(?, '%')")
		params = append(params, v)
	}
	if v, ok := query["postName"]; ok && v != "" {
		conditions = append(conditions, "post_name like concat(?, '%')")
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
		"rows":  []model.SysPost{},
	}

	// 查询数量 长度为0直接返回
	totalSql := "select count(1) as 'total' from sys_post"
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
	pageSql := " order by post_sort limit ?,? "
	params = append(params, pageNum*pageSize)
	params = append(params, pageSize)

	// 查询数据
	querySql := r.selectSql + whereSql + pageSql
	rows, err := db.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err => %v", err)
	}

	// 转换实体
	result["rows"] = db.ConvertResultRows[model.SysPost](model.SysPost{}, r.resultMap, rows)
	return result
}

// Select 查询集合
func (r *SysPostRepository) Select(sysPost model.SysPost) []model.SysPost {
	// 查询条件拼接
	var conditions []string
	var params []any
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
	rows, err := db.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysPost{}
	}
	return db.ConvertResultRows[model.SysPost](model.SysPost{}, r.resultMap, rows)
}

// SelectByIds 通过ID查询信息
func (r *SysPostRepository) SelectByIds(postIds []string) []model.SysPost {
	placeholder := db.KeyPlaceholderByQuery(len(postIds))
	querySql := r.selectSql + " where post_id in (" + placeholder + ")"
	parameters := db.ConvertIdsSlice(postIds)
	rows, err := db.RawDB("", querySql, parameters)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysPost{}
	}
	// 转换实体
	return db.ConvertResultRows[model.SysPost](model.SysPost{}, r.resultMap, rows)
}

// Update 修改信息
func (r *SysPostRepository) Update(sysPost model.SysPost) int64 {
	// 参数拼接
	params := make(map[string]any)
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
	params["remark"] = sysPost.Remark
	if sysPost.UpdateBy != "" {
		params["update_by"] = sysPost.UpdateBy
		params["update_time"] = time.Now().UnixMilli()
	}

	// 构建执行语句
	keys, values := db.KeyValueByUpdate(params)
	sql := fmt.Sprintf("update sys_post set %s where post_id = ?", keys)

	// 执行更新
	values = append(values, sysPost.PostId)
	rows, err := db.ExecDB("", sql, values)
	if err != nil {
		logger.Errorf("update row : %v", err.Error())
		return 0
	}
	return rows
}

// Insert 新增信息
func (r *SysPostRepository) Insert(sysPost model.SysPost) string {
	// 参数拼接
	params := make(map[string]any)
	if sysPost.PostId != "" {
		params["post_id"] = sysPost.PostId
	}
	if sysPost.PostCode != "" {
		params["post_code"] = sysPost.PostCode
	}
	if sysPost.PostName != "" {
		params["post_name"] = sysPost.PostName
	}
	if sysPost.PostSort > 0 {
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
		params["create_time"] = time.Now().UnixMilli()
	}

	// 构建执行语句
	keys, values, placeholder := db.KeyValuePlaceholderByInsert(params)
	sql := fmt.Sprintf("insert into sys_post (%s)values(%s)", keys, placeholder)

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

// DeleteByIds 批量删除信息
func (r *SysPostRepository) DeleteByIds(postIds []string) int64 {
	placeholder := db.KeyPlaceholderByQuery(len(postIds))
	sql := fmt.Sprintf("delete from sys_post where post_id in (%s)", placeholder)
	parameters := db.ConvertIdsSlice(postIds)
	results, err := db.ExecDB("", sql, parameters)
	if err != nil {
		logger.Errorf("delete err => %v", err)
		return 0
	}
	return results
}

// SelectByUserId 根据用户ID获取岗位选择框列表
func (r *SysPostRepository) SelectByUserId(userId string) []model.SysPost {
	// 查询数据
	querySql := `select distinct 
	p.post_id, p.post_name, p.post_code 
	from sys_post p 
    left join sys_user_post up on up.post_id = p.post_id 
    left join sys_user u on u.user_id = up.user_id 
    where u.user_id = ? order by p.post_id`
	rows, err := db.RawDB("", querySql, []any{userId})
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysPost{}
	}
	return db.ConvertResultRows[model.SysPost](model.SysPost{}, r.resultMap, rows)
}

// CheckUnique 检查信息是否唯一
func (r *SysPostRepository) CheckUnique(sysPost model.SysPost) string {
	// 查询条件拼接
	var conditions []string
	var params []any
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
		return "-"
	}

	// 查询数据
	querySql := fmt.Sprintf("select post_id as 'str' from sys_post %s limit 1", whereSql)
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
