package service

import (
	"mask_api_gin/src/framework/vo"
	"mask_api_gin/src/modules/system/model"
)

// ISysMenuService 菜单 服务层接口
type ISysMenuService interface {
	// Find 查询数据
	Find(sysMenu model.SysMenu, userId string) []model.SysMenu

	// FindById 通过ID查询信息
	FindById(menuId string) model.SysMenu

	// Insert 新增信息
	Insert(sysMenu model.SysMenu) string

	// Update 修改信息
	Update(sysMenu model.SysMenu) int64

	// DeleteById 删除信息
	DeleteById(menuId string) int64

	// ExistChildrenByMenuIdAndStatus 菜单下同状态存在子节点数量
	ExistChildrenByMenuIdAndStatus(menuId, status string) int64

	// ExistRoleByMenuId 菜单分配给的角色数量
	ExistRoleByMenuId(menuId string) int64

	// CheckUniqueParentIdByMenuName 检查同级下菜单名称是否唯一
	CheckUniqueParentIdByMenuName(parentId, menuName, menuId string) bool

	// CheckUniqueParentIdByMenuPath 检查同级下路由地址是否唯一（针对目录和菜单）
	CheckUniqueParentIdByMenuPath(parentId, path, menuId string) bool

	// FindPermsByUserId 根据用户ID查询权限标识
	FindPermsByUserId(userId string) []string

	// FindByRoleId 根据角色ID查询菜单树信息
	FindByRoleId(roleId string) []string

	// BuildTreeMenusByUserId 根据用户ID查询菜单树状嵌套
	BuildTreeMenusByUserId(userId string) []model.SysMenu

	// BuildTreeSelectByUserId 根据用户ID查询菜单树状结构
	BuildTreeSelectByUserId(sysMenu model.SysMenu, userId string) []vo.TreeSelect

	// BuildRouteMenus 构建前端路由所需要的菜单
	BuildRouteMenus(sysMenus []model.SysMenu, prefix string) []vo.Router
}
