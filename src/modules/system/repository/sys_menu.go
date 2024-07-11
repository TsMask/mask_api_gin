package repository

import "mask_api_gin/src/modules/system/model"

// ISysMenuRepository 菜单表 数据层接口
type ISysMenuRepository interface {
	// Select 查询集合
	Select(sysMenu model.SysMenu, userId string) []model.SysMenu

	// SelectByIds 通过ID查询信息
	SelectByIds(menuIds []string) []model.SysMenu

	// Insert 新增信息
	Insert(sysMenu model.SysMenu) string

	// Update 修改信息
	Update(sysMenu model.SysMenu) int64

	// DeleteById 删除信息
	DeleteById(menuId string) int64

	// ExistChildrenByMenuIdAndStatus 菜单下同状态存在子节点数量
	ExistChildrenByMenuIdAndStatus(menuId, status string) int64

	// CheckUnique 检查信息是否唯一
	CheckUnique(sysMenu model.SysMenu) string

	// SelectPermsByUserId 根据用户ID查询权限标识
	SelectPermsByUserId(userId string) []string
	
	// SelectByRoleId 根据角色ID查询菜单树信息
	SelectByRoleId(roleId string, menuCheckStrictly bool) []string

	// SelectTreeByUserId 根据用户ID查询菜单
	SelectTreeByUserId(userId string) []model.SysMenu
}
