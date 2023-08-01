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

// 实例化数据层 SysDeptImpl 结构体
var NewSysDeptImpl = &SysDeptImpl{
	selectSql: `select 
	d.dept_id, d.parent_id, d.ancestors, d.dept_name, d.order_num, d.leader, d.phone, d.email, d.status, d.del_flag, d.create_by, d.create_time 
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

// SysDeptImpl 部门表 数据层处理
type SysDeptImpl struct {
	// 查询视图对象SQL
	selectSql string
	// 结果字段与实体映射
	resultMap map[string]string
}

// convertResultRows 将结果记录转实体结果组
func (r *SysDeptImpl) convertResultRows(rows []map[string]any) []model.SysDept {
	arr := make([]model.SysDept, 0)
	for _, row := range rows {
		sysDept := model.SysDept{}
		for key, value := range row {
			if keyMapper, ok := r.resultMap[key]; ok {
				repo.SetFieldValue(&sysDept, keyMapper, value)
			}
		}
		arr = append(arr, sysDept)
	}
	return arr
}

// SelectDeptList 查询部门管理数据
func (r *SysDeptImpl) SelectDeptList(sysDept model.SysDept, dataScopeSQL string) []model.SysDept {
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
	results, err := datasource.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysDept{}
	}

	// 转换实体
	return r.convertResultRows(results)
}

// SelectDeptListByRoleId 根据角色ID查询部门树信息
func (r *SysDeptImpl) SelectDeptListByRoleId(roleId string, deptCheckStrictly bool) []string {
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
	results, err := datasource.RawDB("", querySql, params)
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

// SelectDeptById 根据部门ID查询信息
func (r *SysDeptImpl) SelectDeptById(deptId string) model.SysDept {
	querySql := `select d.dept_id, d.parent_id, d.ancestors,
	d.dept_name, d.order_num, d.leader, d.phone, d.email, d.status,
	(select dept_name from sys_dept where dept_id = d.parent_id) parent_name
	from sys_dept d where d.dept_id = ?`
	results, err := datasource.RawDB("", querySql, []any{deptId})
	if err != nil {
		logger.Errorf("query err => %v", err)
		return model.SysDept{}
	}
	// 转换实体
	rows := r.convertResultRows(results)
	if len(rows) > 0 {
		return rows[0]
	}
	return model.SysDept{}
}

// SelectChildrenDeptById 根据ID查询所有子部门
func (r *SysDeptImpl) SelectChildrenDeptById(deptId string) []model.SysDept {
	querySql := r.selectSql + " where find_in_set(?, d.ancestors)"
	results, err := datasource.RawDB("", querySql, []any{deptId})
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysDept{}
	}

	// 转换实体
	return r.convertResultRows(results)
}

// HasChildByDeptId 是否存在子节点
func (r *SysDeptImpl) HasChildByDeptId(deptId string) int64 {
	querySql := "select count(1) as 'total' from sys_dept where del_flag = '0' and parent_id = ? limit 1"
	results, err := datasource.RawDB("", querySql, []any{deptId})
	if err != nil {
		logger.Errorf("query err => %v", err)
		return 0
	}
	if len(results) > 0 {
		return parse.Number(results[0]["total"])
	}
	return 0
}

// CheckDeptExistUser 查询部门是否存在用户
func (r *SysDeptImpl) CheckDeptExistUser(deptId string) int64 {
	querySql := "select count(1) as 'total' from sys_user where dept_id = ? and del_flag = '0'"
	results, err := datasource.RawDB("", querySql, []any{deptId})
	if err != nil {
		logger.Errorf("query err => %v", err)
		return 0
	}
	if len(results) > 0 {
		return parse.Number(results[0]["total"])
	}
	return 0
}

// CheckUniqueDept 校验部门是否唯一
func (r *SysDeptImpl) CheckUniqueDept(sysDept model.SysDept) string {
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
		return ""
	}

	// 查询数据
	querySql := "select dept_id as 'str' from sys_dept " + whereSql + " and del_flag = '0' limit 1"
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

// InsertDept 新增部门信息
func (r *SysDeptImpl) InsertDept(sysDept model.SysDept) string {
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
	keys, placeholder, values := repo.KeyPlaceholderValueByInsert(params)
	sql := "insert into sys_dept (" + strings.Join(keys, ",") + ")values(" + placeholder + ")"

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

// UpdateDept 修改部门信息
func (r *SysDeptImpl) UpdateDept(sysDept model.SysDept) int64 {
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
	if sysDept.UpdateBy != "" {
		params["update_by"] = sysDept.UpdateBy
		params["update_time"] = time.Now().UnixMilli()
	}

	// 构建执行语句
	keys, values := repo.KeyValueByUpdate(params)
	sql := "update sys_dept set " + strings.Join(keys, ",") + " where dept_id = ?"

	// 执行更新
	values = append(values, sysDept.DeptID)
	rows, err := datasource.ExecDB("", sql, values)
	if err != nil {
		logger.Errorf("update row : %v", err.Error())
		return 0
	}
	return rows
}

// UpdateDeptStatusNormal 修改所在部门正常状态
func (r *SysDeptImpl) UpdateDeptStatusNormal(deptIds []string) int64 {
	placeholder := repo.KeyPlaceholderByQuery(len(deptIds))
	sql := "update sys_dept set status = '1' where dept_id in (" + placeholder + ")"
	parameters := repo.ConvertIdsSlice(deptIds)
	results, err := datasource.ExecDB("", sql, parameters)
	if err != nil {
		logger.Errorf("delete err => %v", err)
		return 0
	}
	return results
}

// UpdateDeptChildren 修改子元素关系
func (r *SysDeptImpl) UpdateDeptChildren(sysDepts []model.SysDept) int64 {
	// 无参数
	if len(sysDepts) == 0 {
		return 0
	}

	// 更新条件拼接
	var conditions []string
	var params []any
	for _, dept := range sysDepts {
		caseSql := fmt.Sprintf("case when %s then %s end", dept.DeptID, dept.Ancestors)
		conditions = append(conditions, caseSql)
		params = append(params, dept.DeptID)
	}

	cases := strings.Join(conditions, " ")
	placeholders := repo.KeyPlaceholderByQuery(len(params))
	sql := "update sys_dept set ancestors = " + cases + " where dept_id in (" + placeholders + ")"
	results, err := datasource.ExecDB("", sql, params)
	if err != nil {
		logger.Errorf("delete err => %v", err)
		return 0
	}
	return results
}

// DeleteDeptById 删除部门管理信息
func (r *SysDeptImpl) DeleteDeptById(deptId string) int64 {
	sql := "update sys_dept set del_flag = '1' where dept_id = ?"
	results, err := datasource.ExecDB("", sql, []any{deptId})
	if err != nil {
		logger.Errorf("delete err => %v", err)
		return 0
	}
	return results
}
