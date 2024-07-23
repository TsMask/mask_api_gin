package service

import (
	"encoding/base64"
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
var NewSysMenu = &SysMenuService{
	sysMenuRepository:     repository.NewSysMenu,
	sysRoleMenuRepository: repository.NewSysRoleMenu,
	sysRoleRepository:     repository.NewSysRole,
}

// SysMenuService 菜单 服务层处理
type SysMenuService struct {
	sysMenuRepository     repository.ISysMenuRepository     // 菜单服务
	sysRoleMenuRepository repository.ISysRoleMenuRepository // 角色与菜单关联服务
	sysRoleRepository     repository.ISysRoleRepository     // 角色服务
}

// Find 查询数据
func (r *SysMenuService) Find(sysMenu model.SysMenu, userId string) []model.SysMenu {
	return r.sysMenuRepository.Select(sysMenu, userId)
}

// FindById 通过ID查询信息
func (r *SysMenuService) FindById(menuId string) model.SysMenu {
	if menuId == "" {
		return model.SysMenu{}
	}
	menus := r.sysMenuRepository.SelectByIds([]string{menuId})
	if len(menus) > 0 {
		return menus[0]
	}
	return model.SysMenu{}
}

// Insert 新增信息
func (r *SysMenuService) Insert(sysMenu model.SysMenu) string {
	return r.sysMenuRepository.Insert(sysMenu)
}

// Update 修改信息
func (r *SysMenuService) Update(sysMenu model.SysMenu) int64 {
	return r.sysMenuRepository.Update(sysMenu)
}

// DeleteById 删除信息
func (r *SysMenuService) DeleteById(menuId string) int64 {
	r.sysRoleMenuRepository.DeleteByMenuIds([]string{menuId}) // 删除菜单与角色关联
	return r.sysMenuRepository.DeleteById(menuId)
}

// ExistChildrenByMenuIdAndStatus 菜单下同状态存在子节点数量
func (r *SysMenuService) ExistChildrenByMenuIdAndStatus(menuId, status string) int64 {
	return r.sysMenuRepository.ExistChildrenByMenuIdAndStatus(menuId, status)
}

// ExistRoleByMenuId 菜单分配给的角色数量
func (r *SysMenuService) ExistRoleByMenuId(menuId string) int64 {
	return r.sysRoleMenuRepository.ExistRoleByMenuId(menuId)
}

// CheckUniqueParentIdByMenuName 检查同级下菜单名称是否唯一
func (r *SysMenuService) CheckUniqueParentIdByMenuName(parentId, menuName, menuId string) bool {
	uniqueId := r.sysMenuRepository.CheckUnique(model.SysMenu{
		MenuName: menuName,
		ParentID: parentId,
	})
	if uniqueId == menuId {
		return true
	}
	return uniqueId == ""
}

// CheckUniqueParentIdByMenuPath 检查同级下路由地址是否唯一（针对目录和菜单）
func (r *SysMenuService) CheckUniqueParentIdByMenuPath(parentId, path, menuId string) bool {
	uniqueId := r.sysMenuRepository.CheckUnique(model.SysMenu{
		Path:     path,
		ParentID: parentId,
	})
	if uniqueId == menuId {
		return true
	}
	return uniqueId == ""
}

// FindPermsByUserId 根据用户ID查询权限标识
func (r *SysMenuService) FindPermsByUserId(userId string) []string {
	return r.sysMenuRepository.SelectPermsByUserId(userId)
}

// FindByRoleId 根据角色ID查询菜单树信息 TODO
func (r *SysMenuService) FindByRoleId(roleId string) []string {
	roles := r.sysRoleRepository.SelectByIds([]string{roleId})
	if len(roles) > 0 {
		role := roles[0]
		if role.RoleID == roleId {
			return r.sysMenuRepository.SelectByRoleId(
				role.RoleID,
				role.MenuCheckStrictly == "1",
			)
		}
	}
	return []string{}
}

// BuildTreeMenusByUserId 根据用户ID查询菜单树状嵌套
func (r *SysMenuService) BuildTreeMenusByUserId(userId string) []model.SysMenu {
	sysMenus := r.sysMenuRepository.SelectTreeByUserId(userId)
	return r.parseDataToTree(sysMenus)
}

// BuildTreeSelectByUserId 根据用户ID查询菜单树状结构
func (r *SysMenuService) BuildTreeSelectByUserId(sysMenu model.SysMenu, userId string) []vo.TreeSelect {
	sysMenus := r.sysMenuRepository.Select(sysMenu, userId)
	menus := r.parseDataToTree(sysMenus)
	tree := make([]vo.TreeSelect, 0)
	for _, v := range menus {
		tree = append(tree, vo.SysMenuTreeSelect(v))
	}
	return tree
}

// parseDataToTree 将数据解析为树结构，构建前端所需要下拉树结构
func (r *SysMenuService) parseDataToTree(sysMenus []model.SysMenu) []model.SysMenu {
	// 节点分组
	nodesMap := make(map[string][]model.SysMenu)
	// 节点id
	var treeIds []string
	// 树节点
	var tree []model.SysMenu

	for _, item := range sysMenus {
		parentID := item.ParentID
		// 分组
		mapItem, ok := nodesMap[parentID]
		if !ok {
			mapItem = []model.SysMenu{}
		}
		mapItem = append(mapItem, item)
		nodesMap[parentID] = mapItem
		// 记录节点ID
		treeIds = append(treeIds, item.MenuID)
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
		iN := r.parseDataToTreeComponent(node, &nodesMap)
		tree[i] = iN
	}

	return tree
}

// parseDataToTreeComponent 递归函数处理子节点
func (r *SysMenuService) parseDataToTreeComponent(node model.SysMenu, nodesMap *map[string][]model.SysMenu) model.SysMenu {
	id := node.MenuID
	children, ok := (*nodesMap)[id]
	if ok {
		node.Children = children
	}
	if len(node.Children) > 0 {
		for i, child := range node.Children {
			icN := r.parseDataToTreeComponent(child, nodesMap)
			node.Children[i] = icN
		}
	}
	return node
}

// BuildRouteMenus 构建前端路由所需要的菜单
func (r *SysMenuService) BuildRouteMenus(sysMenus []model.SysMenu, prefix string) []vo.Router {
	var routers []vo.Router
	for _, item := range sysMenus {
		router := vo.Router{}
		router.Name = r.getRouteName(item)
		router.Path = r.getRouterPath(item)
		router.Component = r.getComponent(item)
		router.Meta = r.getRouteMeta(item)

		// 子项菜单 目录类型 非路径链接
		cMenus := item.Children
		if len(cMenus) > 0 && item.MenuType == constMenu.TypeDir && !regular.ValidHttp(item.Path) {
			// 获取重定向地址
			redirectPrefix, redirectPath := r.getRouteRedirect(
				cMenus,
				router.Path,
				prefix,
			)
			router.Redirect = redirectPath
			// 子菜单进入递归
			router.Children = r.BuildRouteMenus(cMenus, redirectPrefix)
		} else if item.ParentID == "0" && item.IsFrame == constSystem.StatusYes && item.MenuType == constMenu.TypeMenu {
			// 父菜单 内部跳转 菜单类型
			menuPath := "/" + item.MenuID
			childPath := menuPath + r.getRouterPath(item)
			children := vo.Router{
				Name:      r.getRouteName(item),
				Path:      childPath,
				Component: item.Component,
				Meta:      r.getRouteMeta(item),
			}
			router.Meta.HideChildInMenu = true
			router.Children = append(router.Children, children)
			router.Name = item.MenuID
			router.Path = menuPath
			router.Redirect = childPath
			router.Component = constMenu.ComponentLayoutBasic
		} else if item.ParentID == "0" && item.IsFrame == constSystem.StatusYes && regular.ValidHttp(item.Path) {
			// 父菜单 内部跳转 路径链接
			menuPath := "/" + item.MenuID
			childPath := menuPath + r.getRouterPath(item)
			children := vo.Router{
				Name:      r.getRouteName(item),
				Path:      childPath,
				Component: constMenu.ComponentLayoutLink,
				Meta:      r.getRouteMeta(item),
			}
			router.Meta.HideChildInMenu = true
			router.Children = append(router.Children, children)
			router.Name = item.MenuID
			router.Path = menuPath
			router.Redirect = childPath
			router.Component = constMenu.ComponentLayoutBasic
		}

		routers = append(routers, router)
	}
	return routers
}

// getRouteName 获取路由名称 路径英文首字母大写
func (r *SysMenuService) getRouteName(sysMenu model.SysMenu) string {
	routerName := parse.ConvertToCamelCase(sysMenu.Path)
	// 路径链接
	if regular.ValidHttp(sysMenu.Path) {
		routerName = routerName[:5] + "Link"
	}
	// 拼上菜单ID防止name重名
	return routerName + "_" + sysMenu.MenuID
}

// getRouterPath 获取路由地址
func (r *SysMenuService) getRouterPath(sysMenu model.SysMenu) string {
	routerPath := sysMenu.Path

	// 显式路径
	if routerPath == "" || strings.HasPrefix(routerPath, "/") {
		return routerPath
	}

	// 路径链接 内部跳转
	if regular.ValidHttp(routerPath) && sysMenu.IsFrame == constSystem.StatusYes {
		routerPath = regular.Replace(`/^http(s)?:\/\/+/`, routerPath, "")
		routerPath = base64.StdEncoding.EncodeToString([]byte(routerPath))
	}

	// 父菜单 内部跳转
	if sysMenu.ParentID == "0" && sysMenu.IsFrame == constSystem.StatusYes {
		routerPath = "/" + routerPath
	}

	return routerPath
}

// getComponent 获取组件信息
func (r *SysMenuService) getComponent(sysMenu model.SysMenu) string {
	// 内部跳转 路径链接
	if sysMenu.IsFrame == constSystem.StatusYes && regular.ValidHttp(sysMenu.Path) {
		return constMenu.ComponentLayoutLink
	}

	// 非父菜单 目录类型
	if sysMenu.ParentID != "0" && sysMenu.MenuType == constMenu.TypeDir {
		return constMenu.ComponentLayoutBlank
	}

	// 组件路径 内部跳转 菜单类型
	if sysMenu.Component != "" && sysMenu.IsFrame == constSystem.StatusYes && sysMenu.MenuType == constMenu.TypeMenu {
		// 父菜单套外层布局
		if sysMenu.ParentID == "0" {
			return constMenu.ComponentLayoutBasic
		}
		return sysMenu.Component
	}

	return constMenu.ComponentLayoutBasic
}

// getRouteMeta 获取路由元信息
func (r *SysMenuService) getRouteMeta(sysMenu model.SysMenu) vo.RouterMeta {
	meta := vo.RouterMeta{}
	if sysMenu.Icon == "#" {
		meta.Icon = ""
	} else {
		meta.Icon = sysMenu.Icon
	}
	meta.Title = sysMenu.MenuName
	meta.HideChildInMenu = false
	meta.HideInMenu = sysMenu.Visible == constSystem.StatusNo
	meta.Cache = sysMenu.IsCache == constSystem.StatusYes
	meta.Target = ""

	// 路径链接 非内部跳转
	if regular.ValidHttp(sysMenu.Path) && sysMenu.IsFrame == constSystem.StatusNo {
		meta.Target = "_blank"
	}

	return meta
}

// getRouteRedirect 获取路由重定向地址（针对目录）
//
// cMenus 子菜单数组
// routerPath 当前菜单路径
// prefix 菜单重定向路径前缀
func (r *SysMenuService) getRouteRedirect(cMenus []model.SysMenu, routerPath string, prefix string) (string, string) {
	redirectPath := ""

	// 重定向为首个显示并启用的子菜单
	var firstChild *model.SysMenu
	for _, item := range cMenus {
		if item.IsFrame == constSystem.StatusYes && item.Visible == constSystem.StatusYes {
			firstChild = &item
			break
		}
	}

	// 检查内嵌隐藏菜单是否可做重定向
	if firstChild == nil {
		for _, item := range cMenus {
			if item.IsFrame == constSystem.StatusYes && item.Visible == constSystem.StatusNo && strings.Contains(item.Path, constMenu.PathInline) {
				firstChild = &item
				break
			}
		}
	}

	if firstChild != nil {
		firstChildPath := r.getRouterPath(*firstChild)
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
