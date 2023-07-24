package repository

import (
	"fmt"
	"mask_api_gin/src/framework/datasource"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/framework/utils/repo"
	"mask_api_gin/src/modules/system/model"
	"strings"
)

// SysUserRoleImpl 用户与角色关联表 数据层处理
var SysUserRoleImpl = &sysUserRoleImpl{
	selectSql: "",
}

type sysUserRoleImpl struct {
	// 查询视图对象SQL
	selectSql string
}

// CountUserRoleByRoleId 通过角色ID查询角色使用数量
func (r *sysUserRoleImpl) CountUserRoleByRoleId(roleId string) int64 {
	querySql := "select count(1) as total from sys_user_role where role_id = ?"
	results, err := datasource.RawDB("", querySql, []interface{}{roleId})
	if err != nil {
		logger.Errorf("query err => %v", err)
		return 0
	}
	if len(results) > 0 {
		return parse.Number(results[0]["total"])
	}
	return 0
}

// BatchUserRole 批量新增用户角色信息
func (r *sysUserRoleImpl) BatchUserRole(sysUserRoles []model.SysUserRole) int64 {
	keyValues := make([]string, 0)
	for _, item := range sysUserRoles {
		keyValues = append(keyValues, fmt.Sprintf("(%s,%s)", item.UserID, item.RoleID))
	}
	sql := "insert into sys_user_role(user_id, role_id) values " + strings.Join(keyValues, ",")
	results, err := datasource.ExecDB("", sql, nil)
	if err != nil {
		logger.Errorf("delete err => %v", err)
		return 0
	}
	return results
}

// DeleteUserRole 批量删除用户和角色关联
func (r *sysUserRoleImpl) DeleteUserRole(userIds []string) int64 {
	placeholder := repo.KeyPlaceholderByQuery(len(userIds))
	sql := "delete from sys_user_role where user_id in (" + placeholder + ")"
	parameters := repo.ConvertIdsSlice(userIds)
	results, err := datasource.ExecDB("", sql, parameters)
	if err != nil {
		logger.Errorf("delete err => %v", err)
		return 0
	}
	return results
}

// DeleteUserRoleByRoleId 批量取消授权用户角色
func (r *sysUserRoleImpl) DeleteUserRoleByRoleId(roleId string, userIds []string) int64 {
	placeholder := repo.KeyPlaceholderByQuery(len(userIds))
	sql := "delete from sys_user_role where role_id= ? and user_id in (" + placeholder + ")"
	parameters := repo.ConvertIdsSlice(userIds)
	parameters = append([]interface{}{roleId}, parameters...)
	results, err := datasource.ExecDB("", sql, parameters)
	if err != nil {
		logger.Errorf("delete err => %v", err)
		return 0
	}
	return results
}
