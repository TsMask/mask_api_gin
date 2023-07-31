package service

import (
	"errors"
	"fmt"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/repository"
)

// 实例化服务层 SysRoleImpl 结构体
var NewSysRoleImpl = &SysRoleImpl{
	sysRoleRepository:     repository.NewSysRoleImpl,
	sysUserRoleRepository: repository.NewSysUserRoleImpl,
	sysRoleDeptRepository: repository.NewSysRoleDeptImpl,
	sysRoleMenuRepository: repository.NewSysRoleMenuImpl,
}

// SysRoleImpl 角色 服务层处理
type SysRoleImpl struct {
	// 角色服务
	sysRoleRepository repository.ISysRole
	//  用户与角色关联服务
	sysUserRoleRepository repository.ISysUserRole
	// 角色与部门关联服务
	sysRoleDeptRepository repository.ISysRoleDept
	// 角色与菜单关联服务
	sysRoleMenuRepository repository.ISysRoleMenu
}

// SelectRolePage 根据条件分页查询角色数据
func (r *SysRoleImpl) SelectRolePage(query map[string]string, dataScopeSQL string) map[string]interface{} {
	return r.sysRoleRepository.SelectRolePage(query, dataScopeSQL)
}

// SelectRoleList 根据条件查询角色数据
func (r *SysRoleImpl) SelectRoleList(sysRole model.SysRole, dataScopeSQL string) []model.SysRole {
	return r.sysRoleRepository.SelectRoleList(sysRole, dataScopeSQL)
}

// SelectRoleListByUserId 根据用户ID获取角色选择框列表
func (r *SysRoleImpl) SelectRoleListByUserId(userId string) []model.SysRole {
	return r.sysRoleRepository.SelectRoleListByUserId(userId)
}

// SelectRoleById 通过角色ID查询角色
func (r *SysRoleImpl) SelectRoleById(roleId string) model.SysRole {
	if roleId == "" {
		return model.SysRole{}
	}
	posts := r.sysRoleRepository.SelectRoleByIds([]string{roleId})
	if len(posts) > 0 {
		return posts[0]
	}
	return model.SysRole{}
}

// UpdateRole 修改角色信息
func (r *SysRoleImpl) UpdateRole(sysRole model.SysRole) int64 {
	rows := r.sysRoleRepository.UpdateRole(sysRole)
	if rows > 0 && len(sysRole.MenuIds) > 0 {
		// 删除角色与菜单关联
		r.sysRoleMenuRepository.DeleteRoleMenu([]string{sysRole.RoleID})
		r.insertRoleMenu(sysRole.RoleID, sysRole.MenuIds)
	}
	return rows
}

// InsertRole 新增角色信息
func (r *SysRoleImpl) InsertRole(sysRole model.SysRole) string {
	insertId := r.sysRoleRepository.InsertRole(sysRole)
	if insertId != "" && len(sysRole.MenuIds) > 0 {
		r.insertRoleMenu(insertId, sysRole.MenuIds)
	}
	return insertId
}

// insertRoleMenu 新增角色菜单信息
func (r *SysRoleImpl) insertRoleMenu(roleId string, menuIds []string) int64 {
	if roleId == "" || len(menuIds) <= 0 {
		return 0
	}

	sysRoleMenus := []model.SysRoleMenu{}
	for _, menuId := range menuIds {
		if menuId == "" {
			continue
		}
		sysRoleMenus = append(sysRoleMenus, model.NewSysRoleMenu(roleId, menuId))
	}

	return r.sysRoleMenuRepository.BatchRoleMenu(sysRoleMenus)
}

// DeleteRoleByIds 批量删除角色信息
func (r *SysRoleImpl) DeleteRoleByIds(roleIds []string) (int64, error) {
	// 检查是否存在
	roles := r.sysRoleRepository.SelectRoleByIds(roleIds)
	if len(roles) <= 0 {
		return 0, errors.New("没有权限访问角色数据！")
	}
	for _, role := range roles {
		// 检查是否为已删除
		if role.DelFlag == "1" {
			return 0, errors.New(role.RoleID + " 角色信息已经删除！")
		}
		// 检查分配用户
		userCount := r.sysUserRoleRepository.CountUserRoleByRoleId(role.RoleID)
		if userCount > 0 {
			msg := fmt.Sprintf("【%s】已分配给用户,不能删除", role.RoleName)
			return 0, errors.New(msg)
		}
	}
	if len(roles) == len(roleIds) {
		// 删除角色与菜单关联
		r.sysRoleMenuRepository.DeleteRoleMenu(roleIds)
		// 删除角色与部门关联
		r.sysRoleDeptRepository.DeleteRoleDept(roleIds)
		rows := r.sysRoleRepository.DeleteRoleByIds(roleIds)
		return rows, nil
	}
	return 0, errors.New("删除角色信息失败！")
}

// CheckUniqueRoleName 校验角色名称是否唯一
func (r *SysRoleImpl) CheckUniqueRoleName(roleName, roleId string) bool {
	uniqueId := r.sysRoleRepository.CheckUniqueRole(model.SysRole{
		RoleName: roleName,
	})
	if uniqueId == roleId {
		return true
	}
	return uniqueId == ""
}

// CheckUniqueRoleKey 校验角色权限是否唯一
func (r *SysRoleImpl) CheckUniqueRoleKey(roleKey, roleId string) bool {
	uniqueId := r.sysRoleRepository.CheckUniqueRole(model.SysRole{
		RoleKey: roleKey,
	})
	if uniqueId == roleId {
		return true
	}
	return uniqueId == ""
}

// AuthDataScope 修改数据权限信息
func (r *SysRoleImpl) AuthDataScope(sysRole model.SysRole) int64 {
	// 修改角色信息
	rows := r.sysRoleRepository.UpdateRole(sysRole)
	// 删除角色与部门关联
	r.sysRoleDeptRepository.DeleteRoleDept([]string{sysRole.RoleID})
	// 新增角色和部门信息
	if sysRole.DataScope == "2" && len(sysRole.DeptIds) > 0 {
		sysRoleDepts := []model.SysRoleDept{}
		for _, deptId := range sysRole.DeptIds {
			if deptId == "" {
				continue
			}
			sysRoleDepts = append(sysRoleDepts, model.NewSysRoleDept(sysRole.RoleID, deptId))
		}
		rows += r.sysRoleDeptRepository.BatchRoleDept(sysRoleDepts)
	}
	return rows
}

// DeleteAuthUsers 批量取消授权用户角色
func (r *SysRoleImpl) DeleteAuthUsers(roleId string, userIds []string) int64 {
	return r.sysUserRoleRepository.DeleteUserRoleByRoleId(roleId, userIds)
}

// InsertAuthUsers 批量新增授权用户角色
func (r *SysRoleImpl) InsertAuthUsers(roleId string, userIds []string) int64 {
	if roleId == "" || len(userIds) <= 0 {
		return 0
	}

	sysUserRoles := []model.SysUserRole{}
	for _, userId := range userIds {
		if userId == "" {
			continue
		}
		sysUserRoles = append(sysUserRoles, model.NewSysUserRole(userId, roleId))
	}

	return r.sysUserRoleRepository.BatchUserRole(sysUserRoles)
}
