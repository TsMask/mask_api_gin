package repository

import (
	"fmt"
	db "mask_api_gin/src/framework/data_source"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/modules/system/model"
	"strings"
)

// NewSysUserRole 实例化数据层
var NewSysUserRole = &SysUserRoleRepository{}

// SysUserRoleRepository 用户与角色关联表 数据层处理
type SysUserRoleRepository struct{}

// ExistUserByRoleId 存在用户使用数量
func (r *SysUserRoleRepository) ExistUserByRoleId(roleId string) int64 {
	querySql := "select count(1) as total from sys_user_role where role_id = ?"
	results, err := db.RawDB("", querySql, []any{roleId})
	if err != nil {
		logger.Errorf("query err => %v", err)
		return 0
	}
	return parse.Number(results[0]["total"])
}

// DeleteByUserIds 批量删除关联By用户
func (r *SysUserRoleRepository) DeleteByUserIds(userIds []string) int64 {
	placeholder := db.KeyPlaceholderByQuery(len(userIds))
	sql := fmt.Sprintf("delete from sys_user_role where user_id in (%s)", placeholder)
	parameters := db.ConvertIdsSlice(userIds)
	results, err := db.ExecDB("", sql, parameters)
	if err != nil {
		logger.Errorf("delete err => %v", err)
		return 0
	}
	return results
}

// DeleteByRoleId 批量删除关联By角色
func (r *SysUserRoleRepository) DeleteByRoleId(roleId string, userIds []string) int64 {
	placeholder := db.KeyPlaceholderByQuery(len(userIds))
	sql := fmt.Sprintf("delete from sys_user_role where role_id= ? and user_id in (%s)", placeholder)
	parameters := db.ConvertIdsSlice(userIds)
	parameters = append([]any{roleId}, parameters...)
	results, err := db.ExecDB("", sql, parameters)
	if err != nil {
		logger.Errorf("delete err => %v", err)
		return 0
	}
	return results
}

// BatchInsert 批量新增信息
func (r *SysUserRoleRepository) BatchInsert(arr []model.SysUserRole) int64 {
	ur := make([]string, 0)
	for _, item := range arr {
		ur = append(ur, fmt.Sprintf("(%s,%s)", item.UserID, item.RoleID))
	}
	sql := fmt.Sprintf("insert into sys_user_role(user_id, role_id) values %s", strings.Join(ur, ","))
	results, err := db.ExecDB("", sql, nil)
	if err != nil {
		logger.Errorf("delete err => %v", err)
		return 0
	}
	return results
}
