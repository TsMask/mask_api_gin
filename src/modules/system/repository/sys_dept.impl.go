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

// NewSysDept 实例化数据层
var NewSysDept = &SysDeptRepository{
	selectSql: `select 
	d.dept_id, d.parent_id, d.ancestors, d.dept_name, d.order_num, 
	d.leader, d.phone, d.email, d.status, 
	d.del_flag, d.create_by, d.create_time 
	from sys_dept d`,

	resultMap: map[string]string{
		"dept_id":     "DeptID",
		"parent_id":   "ParentID",
		"ancestors":   "Ancestors",
		"dept_name":   "DeptName",
		"order_num":   "OrderNum",
		"leader":      "Leader",
		"phone":       "Phone",
		"email":       "Email",
		"status":      "Status",
		"del_flag":    "DelFlag",
		"create_by":   "CreateBy",
		"create_time": "CreateTime",
		"update_by":   "UpdateBy",
		"update_time": "UpdateTime",
		"parent_name": "ParentName",
	},
}

// SysDeptRepository 部门表 数据层处理
type SysDeptRepository struct {
	selectSql string            // 查询视图对象SQL
	resultMap map[string]string // 结果字段与实体映射
}

// Select 查询集合
func (r *SysDeptRepository) Select(sysDept model.SysDept, dataScopeSQL string) []model.SysDept {
	// 查询条件拼接
	var conditions []string
	var params []any
	if sysDept.DeptID != "" {
		conditions = append(conditions, "dept_id = ?")
		params = append(params, sysDept.DeptID)
	}
	if sysDept.ParentID != "" {
		conditions = append(conditions, "parent_id = ?")
		params = append(params, sysDept.ParentID)
	}
	if sysDept.DeptName != "" {
		conditions = append(conditions, "dept_name like concat(?, '%')")
		params = append(params, sysDept.DeptName)
	}
	if sysDept.Status != "" {
		conditions = append(conditions, "status = ?")
		params = append(params, sysDept.Status)
	}

	// 构建查询条件语句
	whereSql := " where d.del_flag = '0' "
	if len(conditions) > 0 {
		whereSql += " and " + strings.Join(conditions, " and ")
	}

	// 查询数据
	orderSql := " order by d.parent_id, d.order_num asc "
	querySql := r.selectSql + whereSql + dataScopeSQL + orderSql
	rows, err := db.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysDept{}
	}

	// 转换实体
	return db.ConvertResultRows[model.SysDept](model.SysDept{}, r.resultMap, rows)
}

// SelectById 通过ID查询信息
func (r *SysDeptRepository) SelectById(deptId string) model.SysDept {
	querySql := `select d.dept_id, d.parent_id, d.ancestors,
	d.dept_name, d.order_num, d.leader, d.phone, d.email, d.status,
	(select dept_name from sys_dept where dept_id = d.parent_id) parent_name
	from sys_dept d where d.dept_id = ?`
	rows, err := db.RawDB("", querySql, []any{deptId})
	if err != nil {
		logger.Errorf("query err => %v", err)
		return model.SysDept{}
	}
	// 转换实体
	if v := db.ConvertResultRows[model.SysDept](model.SysDept{}, r.resultMap, rows); len(v) > 0 {
		return v[0]
	}
	return model.SysDept{}
}

// Insert 新增信息
func (r *SysDeptRepository) Insert(sysDept model.SysDept) string {
	// 参数拼接
	params := make(map[string]any)
	if sysDept.DeptID != "" {
		params["dept_id"] = sysDept.DeptID
	}
	if sysDept.ParentID != "" {
		params["parent_id"] = sysDept.ParentID
	}
	if sysDept.DeptName != "" {
		params["dept_name"] = sysDept.DeptName
	}
	if sysDept.Ancestors != "" {
		params["ancestors"] = sysDept.Ancestors
	}
	if sysDept.OrderNum > 0 {
		params["order_num"] = sysDept.OrderNum
	}
	if sysDept.Leader != "" {
		params["leader"] = sysDept.Leader
	}
	if sysDept.Phone != "" {
		params["phone"] = sysDept.Phone
	}
	if sysDept.Email != "" {
		params["email"] = sysDept.Email
	}
	if sysDept.Status != "" {
		params["status"] = sysDept.Status
	}
	if sysDept.CreateBy != "" {
		params["create_by"] = sysDept.CreateBy
		params["create_time"] = time.Now().UnixMilli()
	}

	// 构建执行语句
	keys, values, placeholder := db.KeyValuePlaceholderByInsert(params)
	sql := fmt.Sprintf("insert into sys_dept (%s)values(%s)", keys, placeholder)

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
func (r *SysDeptRepository) Update(sysDept model.SysDept) int64 {
	// 参数拼接
	params := make(map[string]any)
	if sysDept.ParentID != "" {
		params["parent_id"] = sysDept.ParentID
	}
	if sysDept.DeptName != "" {
		params["dept_name"] = sysDept.DeptName
	}
	if sysDept.Ancestors != "" {
		params["ancestors"] = sysDept.Ancestors
	}
	if sysDept.OrderNum >= 0 {
		params["order_num"] = sysDept.OrderNum
	}
	params["leader"] = sysDept.Leader
	params["phone"] = sysDept.Phone
	params["email"] = sysDept.Email
	if sysDept.Status != "" {
		params["status"] = sysDept.Status
	}
	if sysDept.UpdateBy != "" {
		params["update_by"] = sysDept.UpdateBy
		params["update_time"] = time.Now().UnixMilli()
	}

	// 构建执行语句
	keys, values := db.KeyValueByUpdate(params)
	sql := fmt.Sprintf("update sys_dept set %s where dept_id = ?", keys)

	// 执行更新
	values = append(values, sysDept.DeptID)
	rows, err := db.ExecDB("", sql, values)
	if err != nil {
		logger.Errorf("update row : %v", err.Error())
		return 0
	}
	return rows
}

// DeleteById 删除信息
func (r *SysDeptRepository) DeleteById(deptId string) int64 {
	sql := "update sys_dept set status = '0', del_flag = '1' where dept_id = ?"
	results, err := db.ExecDB("", sql, []any{deptId})
	if err != nil {
		logger.Errorf("delete err => %v", err)
		return 0
	}
	return results
}

// CheckUnique 检查信息是否唯一
func (r *SysDeptRepository) CheckUnique(sysDept model.SysDept) string {
	// 查询条件拼接
	var conditions []string
	var params []any
	if sysDept.DeptName != "" {
		conditions = append(conditions, "dept_name = ?")
		params = append(params, sysDept.DeptName)
	}
	if sysDept.ParentID != "" {
		conditions = append(conditions, "parent_id = ?")
		params = append(params, sysDept.ParentID)
	}

	// 构建查询条件语句
	whereSql := ""
	if len(conditions) > 0 {
		whereSql += " where " + strings.Join(conditions, " and ")
	}
	if whereSql == "" {
		return "-"
	}

	// 查询数据
	querySql := fmt.Sprintf("select dept_id as 'str' from sys_dept %s and del_flag = '0' limit 1", whereSql)
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

// ExistChildrenByDeptId 存在子节点数量
func (r *SysDeptRepository) ExistChildrenByDeptId(deptId string) int64 {
	querySql := "select count(1) as 'total' from sys_dept where status = '1' and parent_id = ? limit 1"
	results, err := db.RawDB("", querySql, []any{deptId})
	if err != nil {
		logger.Errorf("query err => %v", err)
		return 0
	}
	if len(results) > 0 {
		return parse.Number(results[0]["total"])
	}
	return 0
}

// ExistUserByDeptId 存在用户使用数量
func (r *SysDeptRepository) ExistUserByDeptId(deptId string) int64 {
	querySql := "select count(1) as 'total' from sys_user where dept_id = ? and del_flag = '0'"
	results, err := db.RawDB("", querySql, []any{deptId})
	if err != nil {
		logger.Errorf("query err => %v", err)
		return 0
	}
	if len(results) > 0 {
		return parse.Number(results[0]["total"])
	}
	return 0
}

// SelectDeptIdsByRoleId 通过角色ID查询包含的部门ID
func (r *SysDeptRepository) SelectDeptIdsByRoleId(roleId string, deptCheckStrictly bool) []string {
	querySql := `select d.dept_id as 'str' from sys_dept d
    left join sys_role_dept rd on d.dept_id = rd.dept_id
    where rd.role_id = ? `
	var params []any
	params = append(params, roleId)
	// 展开
	if deptCheckStrictly {
		querySql += ` and d.dept_id not in 
		(select d.parent_id from sys_dept d
		inner join sys_role_dept rd on d.dept_id = rd.dept_id 
		and rd.role_id = ?) `
		params = append(params, roleId)
	}
	querySql += "order by d.parent_id, d.order_num"

	// 查询结果
	results, err := db.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []string{}
	}

	if len(results) > 0 {
		ids := make([]string, 0)
		for _, v := range results {
			ids = append(ids, fmt.Sprintf("%v", v["str"]))
		}
		return ids
	}
	return []string{}
}

// SelectChildrenDeptById 根据ID查询所有子部门
func (r *SysDeptRepository) SelectChildrenDeptById(deptId string) []model.SysDept {
	querySql := r.selectSql + " where find_in_set(?, d.ancestors)"
	rows, err := db.RawDB("", querySql, []any{deptId})
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysDept{}
	}

	// 转换实体
	return db.ConvertResultRows[model.SysDept](model.SysDept{}, r.resultMap, rows)
}

// UpdateDeptStatusNormal 修改所在部门正常状态
func (r *SysDeptRepository) UpdateDeptStatusNormal(deptIds []string) int64 {
	placeholder := db.KeyPlaceholderByQuery(len(deptIds))
	sql := fmt.Sprintf("update sys_dept set status = '1' where dept_id in (%s)", placeholder)
	parameters := db.ConvertIdsSlice(deptIds)
	results, err := db.ExecDB("", sql, parameters)
	if err != nil {
		logger.Errorf("update err => %v", err)
		return 0
	}
	return results
}

// UpdateDeptChildren 修改子元素关系
func (r *SysDeptRepository) UpdateDeptChildren(arr []model.SysDept) int64 {
	// 无参数
	if len(arr) == 0 {
		return 0
	}

	// 更新条件拼接
	var conditions []string
	var params []any
	for _, dept := range arr {
		caseSql := fmt.Sprintf("WHEN dept_id = '%s' THEN '%s'", dept.DeptID, dept.Ancestors)
		conditions = append(conditions, caseSql)
		params = append(params, dept.DeptID)
	}

	cases := strings.Join(conditions, " ")
	placeholders := db.KeyPlaceholderByQuery(len(params))
	sql := fmt.Sprintf("update sys_dept set ancestors = CASE %s END where dept_id in (%s)", cases, placeholders)
	results, err := db.ExecDB("", sql, params)
	if err != nil {
		logger.Errorf("delete err => %v", err)
		return 0
	}
	return results
}
