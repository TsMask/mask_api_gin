package service

import (
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/repository"
)

// SysMenuImpl 菜单 数据层处理
var SysMenuImpl = &sysMenuImpl{
	sysUserRepository: repository.SysUserImpl,
}

type sysMenuImpl struct {
	// 用户服务
	sysUserRepository repository.ISysUser
}

// SelectMenuList 查询系统菜单列表
func (r *sysMenuImpl) SelectMenuList(sysMenu model.SysMenu, userId string) []model.SysMenu {
	return []model.SysMenu{}
}

// SelectMenuPermsByRoleId 根据角色ID查询权限
func (r *sysMenuImpl) SelectMenuPermsByRoleId(roleId string) []string {
	return []string{}
}

// SelectMenuPermsByUserId 根据用户ID查询权限
func (r *sysMenuImpl) SelectMenuPermsByUserId(userId string) []string {
	return []string{}
}

// SelectMenuTreeByUserId 根据用户ID查询菜单
func (r *sysMenuImpl) SelectMenuTreeByUserId(userId string) []model.SysMenu {
	return []model.SysMenu{}
}

// SelectMenuListByRoleId 根据角色ID查询菜单树信息
func (r *sysMenuImpl) SelectMenuListByRoleId(roleId string, menuCheckStrictly bool) []string {
	return []string{}
}

// SelectMenuById 根据菜单ID查询信息
func (r *sysMenuImpl) SelectMenuById(menuId string) model.SysMenu {
	return model.SysMenu{}
}

// HasChildByMenuId 是否存在菜单子节点
func (r *sysMenuImpl) HasChildByMenuId(menuId string) int {
	return 0
}

// CheckMenuExistRole 查询菜单是否存在角色
func (r *sysMenuImpl) CheckMenuExistRole(menuId string) int {
	return 0
}

// InsertMenu 新增菜单信息
func (r *sysMenuImpl) InsertMenu(sysMenu model.SysMenu) string {
	return ""
}

// UpdateMenu 修改菜单信息
func (r *sysMenuImpl) UpdateMenu(sysMenu model.SysMenu) int {
	return 0
}

// DeleteMenuById 删除菜单管理信息
func (r *sysMenuImpl) DeleteMenuById(menuId string) int {
	return 0
}

// CheckUniqueMenuName 校验菜单名称是否唯一
func (r *sysMenuImpl) CheckUniqueMenuName(menuName, parentId string) string {
	return ""
}
