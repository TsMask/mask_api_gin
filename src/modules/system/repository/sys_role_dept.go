package repository

import (
	"fmt"
	db "mask_api_gin/src/framework/data_source"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/modules/system/model"
	"strings"
)

// NewSysRoleDept 实例化数据层
var NewSysRoleDept = &SysRoleDept{}

// SysRoleDept 角色与部门关联表 数据层处理
type SysRoleDept struct{}

// DeleteByRoleIds 批量删除信息By角色
func (r SysRoleDept) DeleteByRoleIds(roleIds []string) int64 {
	placeholder := db.KeyPlaceholderByQuery(len(roleIds))
	sql := fmt.Sprintf("delete from sys_role_dept where role_id in (%s)", placeholder)
	parameters := db.ConvertIdsSlice(roleIds)
	results, err := db.ExecDB("", sql, parameters)
	if err != nil {
		logger.Errorf("delete err => %v", err)
		return 0
	}
	return results
}

// DeleteByDeptIds 批量删除信息By部门
func (r SysRoleDept) DeleteByDeptIds(deptIds []string) int64 {
	placeholder := db.KeyPlaceholderByQuery(len(deptIds))
	sql := fmt.Sprintf("delete from sys_role_dept where dept_id in (%s)", placeholder)
	parameters := db.ConvertIdsSlice(deptIds)
	results, err := db.ExecDB("", sql, parameters)
	if err != nil {
		logger.Errorf("delete err => %v", err)
		return 0
	}
	return results
}

// BatchInsert 批量新增信息
func (r SysRoleDept) BatchInsert(arr []model.SysRoleDept) int64 {
	rd := make([]string, 0)
	for _, item := range arr {
		rd = append(rd, fmt.Sprintf("(%s,%s)", item.RoleId, item.DeptId))
	}
	sql := fmt.Sprintf("insert into sys_role_dept(role_id, dept_id) values %s", strings.Join(rd, ","))
	results, err := db.ExecDB("", sql, nil)
	if err != nil {
		logger.Errorf("delete err => %v", err)
		return 0
	}
	return results
}
