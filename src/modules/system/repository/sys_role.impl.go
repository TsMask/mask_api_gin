package repository

import (
	"mask_api_gin/src/framework/datasource"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/utils/repo"
	"mask_api_gin/src/modules/system/model"
	"strings"
)

// SysRoleImpl 角色表 数据层处理
var SysRoleImpl = &sysRoleImpl{
	selectSql: `select distinct 
	r.role_id, r.role_name, r.role_key, r.role_sort, r.data_scope, r.menu_check_strictly, 
	r.dept_check_strictly, r.status, r.del_flag, r.create_time, r.remark 
	from sys_role r
	left join sys_user_role ur on ur.role_id = r.role_id
	left join sys_user u on u.user_id = ur.user_id
	left join sys_dept d on u.dept_id = d.dept_id`,

	resultMap: map[string]string{
		"role_id":             "RoleID",
		"role_name":           "RoleName",
		"role_key":            "RoleKey",
		"role_sort":           "RoleSort",
		"data_scope":          "DataScope",
		"menu_check_strictly": "MenuCheckStrictly",
		"dept_check_strictly": "DeptCheckStrictly",
		"status":              "Status",
		"del_flag":            "DelFlag",
		"create_by":           "CreateBy",
		"create_time":         "CreateTime",
		"update_by":           "UpdateBy",
		"update_time":         "UpdateTime",
		"remark":              "Remark",
	},
}

type sysRoleImpl struct {
	// 查询视图对象SQL
	selectSql string
	// 结果字段与实体映射
	resultMap map[string]string
}

// convertResultRows 将结果记录转实体结果组
func (r *sysRoleImpl) convertResultRows(rows []map[string]interface{}) []model.SysRole {
	arr := make([]model.SysRole, 0)
	for _, row := range rows {
		sysRole := model.SysRole{}
		for key, value := range row {
			if keyMapper, ok := r.resultMap[key]; ok {
				repo.SetFieldValue(&sysRole, keyMapper, value)
			}
		}
		arr = append(arr, sysRole)
	}
	return arr
}

// SelectRolePage 根据条件分页查询角色数据
func (r *sysRoleImpl) SelectRolePage(query map[string]string, dataScopeSQL string) map[string]interface{} {
	return map[string]interface{}{}
}

// SelectRoleList 根据条件查询角色数据
func (r *sysRoleImpl) SelectRoleList(sysRole model.SysRole, dataScopeSQL string) []model.SysRole {
	// 查询条件拼接
	var conditions []string
	var params []interface{}
	if sysRole.RoleID != "" {
		conditions = append(conditions, "r.role_id = ?")
		params = append(params, sysRole.RoleID)
	}
	if sysRole.RoleKey != "" {
		conditions = append(conditions, "r.role_key like concat(?, '%')")
		params = append(params, sysRole.RoleKey)
	}
	if sysRole.RoleName != "" {
		conditions = append(conditions, "r.role_name like concat(?, '%')")
		params = append(params, sysRole.RoleName)
	}
	if sysRole.Status != "" {
		conditions = append(conditions, "r.status = ?")
		params = append(params, sysRole.Status)
	}

	// 构建查询条件语句
	whereSql := " where r.del_flag = '0' "
	if len(conditions) > 0 {
		whereSql += " and " + strings.Join(conditions, " and ")
	}

	// 查询数据
	orderSql := " order by r.role_sort"
	querySql := r.selectSql + whereSql + dataScopeSQL + orderSql
	rows, err := datasource.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysRole{}
	}
	return r.convertResultRows(rows)
}

// SelectRoleListByUserId 根据用户ID获取角色选择框列表
func (r *sysRoleImpl) SelectRoleListByUserId(userId string) []model.SysRole {
	querySql := r.selectSql + " where r.del_flag = '0' and ur.user_id = ?"
	results, err := datasource.RawDB("", querySql, []interface{}{userId})
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysRole{}
	}
	return r.convertResultRows(results)
}

// SelectRoleByIds 通过角色ID查询角色
func (r *sysRoleImpl) SelectRoleByIds(roleIds []string) []model.SysRole {
	return []model.SysRole{}
}

// SelectRolesByUserName 根据用户名查询角色
func (r *sysRoleImpl) SelectRolesByUserName(userName string) []model.SysRole {
	return []model.SysRole{}
}

// CheckUniqueRoleName 校验角色名称是否唯一
func (r *sysRoleImpl) CheckUniqueRoleName(roleName string) string {
	return ""
}

// CheckUniqueRoleKey 校验角色权限是否唯一
func (r *sysRoleImpl) CheckUniqueRoleKey(roleKey string) string {
	return ""
}

// UpdateRole 修改角色信息
func (r *sysRoleImpl) UpdateRole(sysRole model.SysRole) int {
	return 0
}

// InsertRole 新增角色信息
func (r *sysRoleImpl) InsertRole(sysRole model.SysRole) string {
	return ""
}

// DeleteRoleByIds 批量删除角色信息
func (r *sysRoleImpl) DeleteRoleByIds(roleIds []string) int {
	return 0
}
