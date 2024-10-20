package repository

import (
	"fmt"
	"mask_api_gin/src/framework/datasource"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/modules/system/model"
	"strings"
)

// 实例化数据层 SysRoleDeptImpl 结构体
var NewSysRoleDeptImpl = &SysRoleDeptImpl{}

// SysRoleDeptImpl 角色与部门关联表 数据层处理
type SysRoleDeptImpl struct{}

// DeleteRoleDept 批量删除角色部门关联信息
func (r *SysRoleDeptImpl) DeleteRoleDept(roleIds []string) int64 {
	placeholder := datasource.KeyPlaceholderByQuery(len(roleIds))
	sql := "delete from sys_role_dept where role_id in (" + placeholder + ")"
	parameters := datasource.ConvertIdsSlice(roleIds)
	results, err := datasource.ExecDB("", sql, parameters)
	if err != nil {
		logger.Errorf("delete err => %v", err)
		return 0
	}
	return results
}

// DeleteDeptRole 批量删除部门角色关联信息
func (r *SysRoleDeptImpl) DeleteDeptRole(deptIds []string) int64 {
	placeholder := datasource.KeyPlaceholderByQuery(len(deptIds))
	sql := "delete from sys_role_dept where dept_id in (" + placeholder + ")"
	parameters := datasource.ConvertIdsSlice(deptIds)
	results, err := datasource.ExecDB("", sql, parameters)
	if err != nil {
		logger.Errorf("delete err => %v", err)
		return 0
	}
	return results
}

// BatchRoleDept 批量新增角色部门信息
func (r *SysRoleDeptImpl) BatchRoleDept(sysRoleDepts []model.SysRoleDept) int64 {
	keyValues := make([]string, 0)
	for _, item := range sysRoleDepts {
		keyValues = append(keyValues, fmt.Sprintf("(%s,%s)", item.RoleID, item.DeptID))
	}
	sql := "insert into sys_role_dept(role_id, dept_id) values " + strings.Join(keyValues, ",")
	results, err := datasource.ExecDB("", sql, nil)
	if err != nil {
		logger.Errorf("delete err => %v", err)
		return 0
	}
	return results
}
