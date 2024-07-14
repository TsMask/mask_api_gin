package service

import (
	"fmt"
	constRoleDataScope "mask_api_gin/src/framework/constants/role_data_scope"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/repository"
)

// NewSysRole 实例化服务层
var NewSysRole = &SysRoleService{
	sysRoleRepository:     repository.NewSysRole,
	sysUserRoleRepository: repository.NewSysUserRole,
	sysRoleDeptRepository: repository.NewSysRoleDept,
	sysRoleMenuRepository: repository.NewSysRoleMenu,
}

// SysRoleService 角色 服务层处理
type SysRoleService struct {
	sysRoleRepository     repository.ISysRoleRepository     // 角色服务
	sysUserRoleRepository repository.ISysUserRoleRepository // 用户与角色关联服务
	sysRoleDeptRepository repository.ISysRoleDeptRepository // 角色与部门关联服务
	sysRoleMenuRepository repository.ISysRoleMenuRepository // 角色与菜单关联服务
}

// FindByPage 分页查询列表数据
func (r *SysRoleService) FindByPage(query map[string]any, dataScopeSQL string) map[string]any {
	return r.sysRoleRepository.SelectByPage(query, dataScopeSQL)
}

// Find 查询列表数据
func (r *SysRoleService) Find(sysRole model.SysRole, dataScopeSQL string) []model.SysRole {
	return r.sysRoleRepository.Select(sysRole, dataScopeSQL)
}

// FindById 通过ID查询信息
func (r *SysRoleService) FindById(roleId string) model.SysRole {
	if roleId == "" {
		return model.SysRole{}
	}
	posts := r.sysRoleRepository.SelectByIds([]string{roleId})
	if len(posts) > 0 {
		return posts[0]
	}
	return model.SysRole{}
}

// Insert 新增信息
func (r *SysRoleService) Insert(sysRole model.SysRole) string {
	insertId := r.sysRoleRepository.Insert(sysRole)
	if insertId != "" && len(sysRole.MenuIds) > 0 {
		r.insertRoleMenu(insertId, sysRole.MenuIds)
	}
	return insertId
}

// insertRoleMenu 新增角色菜单信息
func (r *SysRoleService) insertRoleMenu(roleId string, menuIds []string) int64 {
	if roleId == "" || len(menuIds) <= 0 {
		return 0
	}
	sysRoleMenus := make([]model.SysRoleMenu, 0)
	for _, menuId := range menuIds {
		if menuId == "" {
			continue
		}
		sysRoleMenus = append(sysRoleMenus, model.SysRoleMenu{
			RoleID: roleId, MenuID: menuId,
		})
	}
	return r.sysRoleMenuRepository.BatchInsert(sysRoleMenus)
}

// Update 修改信息
func (r *SysRoleService) Update(sysRole model.SysRole) int64 {
	rows := r.sysRoleRepository.Update(sysRole)
	if rows > 0 && len(sysRole.MenuIds) > 0 {
		// 删除角色与菜单关联
		r.sysRoleMenuRepository.DeleteByRoleIds([]string{sysRole.RoleID})
		r.insertRoleMenu(sysRole.RoleID, sysRole.MenuIds)
	}
	return rows
}

// DeleteByIds 批量删除信息
func (r *SysRoleService) DeleteByIds(roleIds []string) (int64, error) {
	// 检查是否存在
	roles := r.sysRoleRepository.SelectByIds(roleIds)
	if len(roles) <= 0 {
		return 0, fmt.Errorf("没有权限访问角色数据！")
	}
	for _, role := range roles {
		// 检查是否为已删除
		if role.DelFlag == "1" {
			return 0, fmt.Errorf("%s 角色信息已经删除！", role.RoleID)
		}
		// 检查分配用户
		if useCount := r.sysUserRoleRepository.ExistUserByRoleId(role.RoleID); useCount > 0 {
			return 0, fmt.Errorf("【%s】已分配给用户,不能删除", role.RoleName)
		}
	}
	if len(roles) == len(roleIds) {
		r.sysRoleMenuRepository.DeleteByRoleIds(roleIds) // 删除角色与菜单关联
		r.sysRoleDeptRepository.DeleteByRoleIds(roleIds) // 删除角色与部门关联
		return r.sysRoleRepository.DeleteByIds(roleIds), nil
	}
	return 0, fmt.Errorf("删除角色信息失败！")
}

// FindByUserId 根据用户ID获取角色选择框列表
func (r *SysRoleService) FindByUserId(userId string) []model.SysRole {
	return r.sysRoleRepository.SelectByUserId(userId)
}

// CheckUniqueByName 检查角色名称是否唯一
func (r *SysRoleService) CheckUniqueByName(roleName, roleId string) bool {
	uniqueId := r.sysRoleRepository.CheckUnique(model.SysRole{
		RoleName: roleName,
	})
	if uniqueId == roleId {
		return true
	}
	return uniqueId == ""
}

// CheckUniqueByKey 检查角色权限是否唯一
func (r *SysRoleService) CheckUniqueByKey(roleKey, roleId string) bool {
	uniqueId := r.sysRoleRepository.CheckUnique(model.SysRole{
		RoleKey: roleKey,
	})
	if uniqueId == roleId {
		return true
	}
	return uniqueId == ""
}

// UpdateAndDataScope 修改信息同时更新数据权限信息
func (r *SysRoleService) UpdateAndDataScope(sysRole model.SysRole) int64 {
	// 修改角色信息
	rows := r.sysRoleRepository.Update(sysRole)
	if rows > 0 {
		// 删除角色与部门关联
		r.sysRoleDeptRepository.DeleteByRoleIds([]string{sysRole.RoleID})
		// 新增角色和部门信息
		if sysRole.DataScope == constRoleDataScope.Custom && len(sysRole.DeptIds) > 0 {
			arr := make([]model.SysRoleDept, 0)
			for _, deptId := range sysRole.DeptIds {
				if deptId == "" {
					continue
				}
				arr = append(arr, model.SysRoleDept{
					RoleID: sysRole.RoleID, DeptID: deptId,
				})
			}
			r.sysRoleDeptRepository.BatchInsert(arr)
		}
	}
	return rows
}

// InsertAuthUsers 批量新增授权用户角色
func (r *SysRoleService) InsertAuthUsers(roleId string, userIds []string) int64 {
	if roleId == "" || len(userIds) <= 0 {
		return 0
	}
	sysUserRoles := make([]model.SysUserRole, 0)
	for _, userId := range userIds {
		if userId == "" {
			continue
		}
		sysUserRoles = append(sysUserRoles, model.SysUserRole{
			RoleID: roleId, UserID: userId,
		})
	}
	return r.sysUserRoleRepository.BatchInsert(sysUserRoles)
}

// DeleteAuthUsers 批量取消授权用户角色
func (r *SysRoleService) DeleteAuthUsers(roleId string, userIds []string) int64 {
	return r.sysUserRoleRepository.DeleteByRoleId(roleId, userIds)
}
