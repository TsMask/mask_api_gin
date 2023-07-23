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

// SysDeptImpl 部门表 数据层处理
var SysDeptImpl = &sysDeptImpl{
	selectSql: "",
}

type sysDeptImpl struct {
	// 查询视图对象SQL
	selectSql string
}

// SelectDeptList 查询部门管理数据
func (r *sysDeptImpl) SelectDeptList(sysDept model.SysDept, dataScopeSQL string) []model.SysDept {
	return []model.SysDept{}
}

// SelectDeptListByRoleId 根据角色ID查询部门树信息
func (r *sysDeptImpl) SelectDeptListByRoleId(roleId string, deptCheckStrictly bool) []string {
	querySql := `select d.dept_id as 'str' from sys_dept d
    left join sys_role_dept rd on d.dept_id = rd.dept_id
    where rd.role_id = ? `
	var params []interface{}
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
func (r *sysDeptImpl) SelectDeptById(deptId string) model.SysDept {
	return model.SysDept{}
}

// SelectChildrenDeptById 根据ID查询所有子部门
func (r *sysDeptImpl) SelectChildrenDeptById(deptId string) []model.SysDept {
	return []model.SysDept{}
}

// SelectNormalChildrenDeptById 根据ID查询所有子部门（正常状态）
func (r *sysDeptImpl) SelectNormalChildrenDeptById(deptId string) int {
	return 0
}

// HasChildByDeptId 是否存在子节点
func (r *sysDeptImpl) HasChildByDeptId(deptId string) int64 {
	querySql := "select count(1) as 'total' from sys_dept where del_flag = '0' and parent_id = ? limit 1"
	results, err := datasource.RawDB("", querySql, []interface{}{deptId})
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
func (r *sysDeptImpl) CheckDeptExistUser(deptId string) int64 {
	querySql := "select count(1) as 'total' from sys_user where dept_id = ? and del_flag = '0'"
	results, err := datasource.RawDB("", querySql, []interface{}{deptId})
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
func (r *sysDeptImpl) CheckUniqueDept(sysDept model.SysDept) string {
	// 查询条件拼接
	var conditions []string
	var params []interface{}
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
func (r *sysDeptImpl) InsertDept(sysDept model.SysDept) string {
	// 参数拼接
	params := make(map[string]interface{})
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
		params["create_time"] = date.NowTimestamp()
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
func (r *sysDeptImpl) UpdateDept(sysDept model.SysDept) int64 {
	// 参数拼接
	params := make(map[string]interface{})
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
		params["update_time"] = date.NowTimestamp()
	}

	// 构建执行语句
	keys, values := repo.KeyValueByUpdate(params)
	sql := "update sys_dept set " + strings.Join(keys, ",") + " where menu_id = ?"

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
func (r *sysDeptImpl) UpdateDeptStatusNormal(deptIds []string) int {
	return 0
}

// UpdateDeptChildren 修改子元素关系
func (r *sysDeptImpl) UpdateDeptChildren(sysDepts []model.SysDept) int {
	return 0
}

// DeleteDeptById 删除部门管理信息
func (r *sysDeptImpl) DeleteDeptById(deptId string) int64 {
	sql := "update sys_dept set del_flag = '1' where dept_id = ?"
	results, err := datasource.ExecDB("", sql, []interface{}{deptId})
	if err != nil {
		logger.Errorf("delete err => %v", err)
		return 0
	}
	return results
}
