package service

import (
	"encoding/base64"
	"fmt"
	constMenu "mask_api_gin/src/framework/constants/menu"
	constSystem "mask_api_gin/src/framework/constants/system"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/framework/utils/regular"
	"mask_api_gin/src/framework/vo"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/repository"
	"strings"
)

// NewSysMenu 实例化服务层
var NewSysMenu = &SysMenu{
	sysMenuRepository:     repository.NewSysMenu,
	sysRoleMenuRepository: repository.NewSysRoleMenu,
	sysRoleRepository:     repository.NewSysRole,
}

// SysMenu 菜单 服务层处理
type SysMenu struct {
	sysMenuRepository     *repository.SysMenu     // 菜单服务
	sysRoleMenuRepository *repository.SysRoleMenu // 角色与菜单关联服务
	sysRoleRepository     *repository.SysRole     // 角色服务
}

// Find 查询数据
func (s SysMenu) Find(sysMenu model.SysMenu, userId int64) []model.SysMenu {
	return s.sysMenuRepository.Select(sysMenu, userId)
}

// FindById 通过ID查询信息
func (s SysMenu) FindById(menuId int64) model.SysMenu {
	if menuId <= 0 {
		return model.SysMenu{}
	}
	menus := s.sysMenuRepository.SelectByIds([]int64{menuId})
	if len(menus) > 0 {
		return menus[0]
	}
	return model.SysMenu{}
}

// Insert 新增信息
func (s SysMenu) Insert(sysMenu model.SysMenu) int64 {
	return s.sysMenuRepository.Insert(sysMenu)
}

// Update 修改信息
func (s SysMenu) Update(sysMenu model.SysMenu) int64 {
	return s.sysMenuRepository.Update(sysMenu)
}

// DeleteById 删除信息
func (s SysMenu) DeleteById(menuId int64) int64 {
	s.sysRoleMenuRepository.DeleteByMenuIds([]int64{menuId}) // 删除菜单与角色关联
	return s.sysMenuRepository.DeleteById(menuId)
}

// ExistChildrenByMenuIdAndStatus 菜单下同状态存在子节点数量
func (s SysMenu) ExistChildrenByMenuIdAndStatus(menuId int64, status string) int64 {
	return s.sysMenuRepository.ExistChildrenByMenuIdAndStatus(menuId, status)
}

// ExistRoleByMenuId 菜单分配给的角色数量
func (s SysMenu) ExistRoleByMenuId(menuId int64) int64 {
	return s.sysRoleMenuRepository.ExistRoleByMenuId(menuId)
}

// CheckUniqueParentIdByMenuName 检查同级下菜单名称是否唯一
func (s SysMenu) CheckUniqueParentIdByMenuName(parentId int64, menuName string, menuId int64) bool {
	uniqueId := s.sysMenuRepository.CheckUnique(model.SysMenu{
		MenuName: menuName,
		ParentId: parentId,
	})
	if uniqueId == menuId {
		return true
	}
	return uniqueId == 0
}

// CheckUniqueParentIdByMenuPath 检查同级下路由地址是否唯一（针对目录和菜单）
func (s SysMenu) CheckUniqueParentIdByMenuPath(parentId int64, path string, menuId int64) bool {
	uniqueId := s.sysMenuRepository.CheckUnique(model.SysMenu{
		MenuPath: path,
		ParentId: parentId,
	})
	if uniqueId == menuId {
		return true
	}
	return uniqueId == 0
}

// FindPermsByUserId 根据用户ID查询权限标识
func (s SysMenu) FindPermsByUserId(userId int64) []string {
	return s.sysMenuRepository.SelectPermsByUserId(userId)
}

// FindByRoleId 根据角色ID查询菜单树信息 TODO
func (s SysMenu) FindByRoleId(roleId int64) []string {
	roles := s.sysRoleRepository.SelectByIds([]int64{roleId})
	if len(roles) > 0 {
		role := roles[0]
		if role.RoleId == roleId {
			return s.sysMenuRepository.SelectByRoleId(
				role.RoleId,
				role.MenuCheckStrictly == "1",
			)
		}
	}
	return []string{}
}

// BuildTreeMenusByUserId 根据用户ID查询菜单树状嵌套
func (s SysMenu) BuildTreeMenusByUserId(userId int64) []model.SysMenu {
	sysMenus := s.sysMenuRepository.SelectTreeByUserId(userId)
	return s.parseDataToTree(sysMenus)
}

// BuildTreeSelectByUserId 根据用户ID查询菜单树状结构
func (s SysMenu) BuildTreeSelectByUserId(sysMenu model.SysMenu, userId int64) []vo.TreeSelect {
	sysMenus := s.sysMenuRepository.Select(sysMenu, userId)
	menus := s.parseDataToTree(sysMenus)
	tree := make([]vo.TreeSelect, 0)
	for _, v := range menus {
		tree = append(tree, vo.SysMenuTreeSelect(v))
	}
	return tree
}

// parseDataToTree 将数据解析为树结构，构建前端所需要下拉树结构
func (s SysMenu) parseDataToTree(sysMenus []model.SysMenu) []model.SysMenu {
	// 节点分组
	nodesMap := make(map[int64][]model.SysMenu)
	// 节点id
	var treeIds []int64
	// 树节点
	var tree []model.SysMenu

	for _, item := range sysMenus {
		ParentId := item.ParentId
		// 分组
		mapItem, ok := nodesMap[ParentId]
		if !ok {
			mapItem = []model.SysMenu{}
		}
		mapItem = append(mapItem, item)
		nodesMap[ParentId] = mapItem
		// 记录节点ID
		treeIds = append(treeIds, item.MenuId)
	}

	for key, value := range nodesMap {
		// 选择不是节点ID的作为树节点
		found := false
		for _, id := range treeIds {
			if id == key {
				found = true
				break
			}
		}
		if !found {
			tree = append(tree, value...)
		}
	}

	for i, node := range tree {
		iN := s.parseDataToTreeComponent(node, &nodesMap)
		tree[i] = iN
	}

	return tree
}

// parseDataToTreeComponent 递归函数处理子节点
func (s SysMenu) parseDataToTreeComponent(node model.SysMenu, nodesMap *map[int64][]model.SysMenu) model.SysMenu {
	id := node.MenuId
	children, ok := (*nodesMap)[id]
	if ok {
		node.Children = children
	}
	if len(node.Children) > 0 {
		for i, child := range node.Children {
			icN := s.parseDataToTreeComponent(child, nodesMap)
			node.Children[i] = icN
		}
	}
	return node
}

// BuildRouteMenus 构建前端路由所需要的菜单
func (s SysMenu) BuildRouteMenus(sysMenus []model.SysMenu, prefix string) []vo.Router {
	var routers []vo.Router
	for _, item := range sysMenus {
		router := vo.Router{
			Name:      s.getRouteName(item),
			Path:      s.getRouterPath(item),
			Component: s.getComponent(item),
			Meta:      s.getRouteMeta(item),
			Children:  []vo.Router{},
		}

		// 子项菜单 目录类型 非路径链接
		cMenus := item.Children
		if len(cMenus) > 0 && item.MenuType == constMenu.TYPE_DIR && !regular.ValidHttp(item.MenuPath) {
			// 获取重定向地址
			redirectPrefix, redirectPath := s.getRouteRedirect(
				cMenus,
				router.Path,
				prefix,
			)
			router.Redirect = redirectPath
			// 子菜单进入递归
			router.Children = s.BuildRouteMenus(cMenus, redirectPrefix)
		} else if item.ParentId == 0 && item.FrameFlag == constSystem.STATUS_YES && item.MenuType == constMenu.TYPE_MENU {
			// 父菜单 内部跳转 菜单类型
			menuPath := fmt.Sprintf("/%d", item.MenuId)
			childPath := menuPath + s.getRouterPath(item)
			children := vo.Router{
				Name:      s.getRouteName(item),
				Path:      childPath,
				Component: item.Component,
				Meta:      s.getRouteMeta(item),
			}
			router.Meta.HideChildInMenu = true
			router.Children = append(router.Children, children)
			router.Name = fmt.Sprint(item.MenuId)
			router.Path = menuPath
			router.Redirect = childPath
			router.Component = constMenu.COMPONENT_LAYOUT_BASIC
		} else if item.ParentId == 0 && item.FrameFlag == constSystem.STATUS_YES && regular.ValidHttp(item.MenuPath) {
			// 父菜单 内部跳转 路径链接
			menuPath := fmt.Sprintf("/%d", item.MenuId)
			childPath := menuPath + s.getRouterPath(item)
			children := vo.Router{
				Name:      s.getRouteName(item),
				Path:      childPath,
				Component: constMenu.COMPONENT_LAYOUT_LINK,
				Meta:      s.getRouteMeta(item),
			}
			router.Meta.HideChildInMenu = true
			router.Children = append(router.Children, children)
			router.Name = fmt.Sprint(item.MenuId)
			router.Path = menuPath
			router.Redirect = childPath
			router.Component = constMenu.COMPONENT_LAYOUT_BASIC
		}

		routers = append(routers, router)
	}
	return routers
}

// getRouteName 获取路由名称 路径英文首字母大写
func (s SysMenu) getRouteName(sysMenu model.SysMenu) string {
	routerName := parse.ConvertToCamelCase(sysMenu.MenuPath)
	// 路径链接
	if regular.ValidHttp(sysMenu.MenuPath) {
		routerName = routerName[:5] + "Link"
	}
	// 拼上菜单ID防止name重名
	return fmt.Sprintf("%s_%d", routerName, sysMenu.MenuId)
}

// getRouterPath 获取路由地址
func (s SysMenu) getRouterPath(sysMenu model.SysMenu) string {
	routerPath := sysMenu.MenuPath

	// 显式路径
	if routerPath == "" || strings.HasPrefix(routerPath, "/") {
		return routerPath
	}

	// 路径链接 内部跳转
	if regular.ValidHttp(routerPath) && sysMenu.FrameFlag == constSystem.STATUS_YES {
		routerPath = regular.Replace(`/^http(s)?:\/\/+/`, routerPath, "")
		routerPath = base64.StdEncoding.EncodeToString([]byte(routerPath))
	}

	// 父菜单 内部跳转
	if sysMenu.ParentId == 0 && sysMenu.FrameFlag == constSystem.STATUS_YES {
		routerPath = "/" + routerPath
	}

	return routerPath
}

// getComponent 获取组件信息
func (s SysMenu) getComponent(sysMenu model.SysMenu) string {
	// 内部跳转 路径链接
	if sysMenu.FrameFlag == constSystem.STATUS_YES && regular.ValidHttp(sysMenu.MenuPath) {
		return constMenu.COMPONENT_LAYOUT_LINK
	}

	// 非父菜单 目录类型
	if sysMenu.ParentId != 0 && sysMenu.MenuType == constMenu.TYPE_DIR {
		return constMenu.COMPONENT_LAYOUT_BLANK
	}

	// 组件路径 内部跳转 菜单类型
	if sysMenu.Component != "" && sysMenu.FrameFlag == constSystem.STATUS_YES && sysMenu.MenuType == constMenu.TYPE_MENU {
		// 父菜单套外层布局
		if sysMenu.ParentId == 0 {
			return constMenu.COMPONENT_LAYOUT_BASIC
		}
		return sysMenu.Component
	}

	return constMenu.COMPONENT_LAYOUT_BASIC
}

// getRouteMeta 获取路由元信息
func (s SysMenu) getRouteMeta(sysMenu model.SysMenu) vo.RouterMeta {
	meta := vo.RouterMeta{}
	if sysMenu.Icon == "#" {
		meta.Icon = ""
	} else {
		meta.Icon = sysMenu.Icon
	}
	meta.Title = sysMenu.MenuName
	meta.HideChildInMenu = false
	meta.HideInMenu = sysMenu.VisibleFlag == constSystem.STATUS_NO
	meta.Cache = sysMenu.CacheFlag == constSystem.STATUS_YES
	meta.Target = ""

	// 路径链接 非内部跳转
	if regular.ValidHttp(sysMenu.MenuPath) && sysMenu.FrameFlag == constSystem.STATUS_NO {
		meta.Target = "_blank"
	}

	return meta
}

// getRouteRedirect 获取路由重定向地址（针对目录）
//
// cMenus 子菜单数组
// routerPath 当前菜单路径
// prefix 菜单重定向路径前缀
func (s SysMenu) getRouteRedirect(cMenus []model.SysMenu, routerPath string, prefix string) (string, string) {
	redirectPath := ""

	// 重定向为首个显示并启用的子菜单
	var firstChild *model.SysMenu
	for _, item := range cMenus {
		if item.FrameFlag == constSystem.STATUS_YES && item.VisibleFlag == constSystem.STATUS_YES {
			firstChild = &item
			break
		}
	}

	// 检查内嵌隐藏菜单是否可做重定向
	if firstChild == nil {
		for _, item := range cMenus {
			if item.FrameFlag == constSystem.STATUS_YES && item.VisibleFlag == constSystem.STATUS_NO && strings.Contains(item.MenuPath, constMenu.PATH_INLINE) {
				firstChild = &item
				break
			}
		}
	}

	if firstChild != nil {
		firstChildPath := s.getRouterPath(*firstChild)
		if strings.HasPrefix(firstChildPath, "/") {
			redirectPath = firstChildPath
		} else {
			// 拼接追加路径
			if !strings.HasPrefix(routerPath, "/") {
				prefix += "/"
			}
			prefix = prefix + routerPath
			redirectPath = prefix + "/" + firstChildPath
		}
	}

	return prefix, redirectPath
}
