package service

import (
	"encoding/base64"
	"mask_api_gin/src/framework/constants/common"
	"mask_api_gin/src/framework/constants/menu"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/framework/utils/regular"
	"mask_api_gin/src/framework/vo"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/repository"
	"strings"
)

// 实例化服务层 SysMenuImpl 结构体
var NewSysMenuImpl = &SysMenuImpl{
	sysMenuRepository:     repository.NewSysMenuImpl,
	sysRoleMenuRepository: repository.NewSysRoleMenuImpl,
	sysRoleRepository:     repository.NewSysRoleImpl,
}

// SysMenuImpl 菜单 服务层处理
type SysMenuImpl struct {
	// 菜单服务
	sysMenuRepository repository.ISysMenu
	// 角色与菜单关联服务
	sysRoleMenuRepository repository.ISysRoleMenu
	// 角色服务
	sysRoleRepository repository.ISysRole
}

// SelectMenuList 查询系统菜单列表
func (r *SysMenuImpl) SelectMenuList(sysMenu model.SysMenu, userId string) []model.SysMenu {
	return r.sysMenuRepository.SelectMenuList(sysMenu, userId)
}

// SelectMenuPermsByUserId 根据用户ID查询权限
func (r *SysMenuImpl) SelectMenuPermsByUserId(userId string) []string {
	return r.sysMenuRepository.SelectMenuPermsByUserId(userId)
}

// SelectMenuTreeByUserId 根据用户ID查询菜单
func (r *SysMenuImpl) SelectMenuTreeByUserId(userId string) []model.SysMenu {
	sysMenus := r.sysMenuRepository.SelectMenuTreeByUserId(userId)
	return r.parseDataToTree(sysMenus)
}

// SelectMenuTreeSelectByUserId 根据用户ID查询菜单树结构信息
func (r *SysMenuImpl) SelectMenuTreeSelectByUserId(sysMenu model.SysMenu, userId string) []vo.TreeSelect {
	sysMenus := r.sysMenuRepository.SelectMenuList(sysMenu, userId)
	menus := r.parseDataToTree(sysMenus)
	tree := make([]vo.TreeSelect, 0)
	for _, menu := range menus {
		tree = append(tree, vo.SysMenuTreeSelect(menu))
	}
	return tree
}

// SelectMenuListByRoleId 根据角色ID查询菜单树信息 TODO
func (r *SysMenuImpl) SelectMenuListByRoleId(roleId string) []string {
	roles := r.sysRoleRepository.SelectRoleByIds([]string{roleId})
	if len(roles) > 0 {
		role := roles[0]
		if role.RoleID == roleId {
			return r.sysMenuRepository.SelectMenuListByRoleId(
				role.RoleID,
				role.MenuCheckStrictly == "1",
			)
		}
	}
	return []string{}
}

// SelectMenuById 根据菜单ID查询信息
func (r *SysMenuImpl) SelectMenuById(menuId string) model.SysMenu {
	if menuId == "" {
		return model.SysMenu{}
	}
	menus := r.sysMenuRepository.SelectMenuByIds([]string{menuId})
	if len(menus) > 0 {
		return menus[0]
	}
	return model.SysMenu{}
}

// HasChildByMenuIdAndStatus 存在菜单子节点数量与状态
func (r *SysMenuImpl) HasChildByMenuIdAndStatus(menuId, status string) int64 {
	return r.sysMenuRepository.HasChildByMenuIdAndStatus(menuId, status)
}

// CheckMenuExistRole 查询菜单是否存在角色
func (r *SysMenuImpl) CheckMenuExistRole(menuId string) int64 {
	return r.sysRoleMenuRepository.CheckMenuExistRole(menuId)
}

// InsertMenu 新增菜单信息
func (r *SysMenuImpl) InsertMenu(sysMenu model.SysMenu) string {
	return r.sysMenuRepository.InsertMenu(sysMenu)
}

// UpdateMenu 修改菜单信息
func (r *SysMenuImpl) UpdateMenu(sysMenu model.SysMenu) int64 {
	return r.sysMenuRepository.UpdateMenu(sysMenu)
}

// DeleteMenuById 删除菜单管理信息
func (r *SysMenuImpl) DeleteMenuById(menuId string) int64 {
	// 删除菜单与角色关联
	r.sysRoleMenuRepository.DeleteMenuRole([]string{menuId})
	return r.sysMenuRepository.DeleteMenuById(menuId)
}

// CheckUniqueMenuName 校验菜单名称是否唯一
func (r *SysMenuImpl) CheckUniqueMenuName(menuName, parentId, menuId string) bool {
	uniqueId := r.sysMenuRepository.CheckUniqueMenu(model.SysMenu{
		MenuName: menuName,
		ParentID: parentId,
	})
	if uniqueId == menuId {
		return true
	}
	return uniqueId == ""
}

// CheckUniqueMenuPath 校验路由地址是否唯一（针对目录和菜单）
func (r *SysMenuImpl) CheckUniqueMenuPath(path, menuId string) bool {
	uniqueId := r.sysMenuRepository.CheckUniqueMenu(model.SysMenu{
		Path: path,
	})
	if uniqueId == menuId {
		return true
	}
	return uniqueId == ""
}

// BuildRouteMenus 构建前端路由所需要的菜单
func (r *SysMenuImpl) BuildRouteMenus(sysMenus []model.SysMenu, prefix string) []vo.Router {
	routers := []vo.Router{}
	for _, item := range sysMenus {
		router := vo.Router{}
		router.Name = r.getRouteName(item)
		router.Path = r.getRouterPath(item)
		router.Component = r.getComponent(item)
		router.Meta = r.getRouteMeta(item)

		// 子项菜单 目录类型 非路径链接
		cMenus := item.Children
		if len(cMenus) > 0 && item.MenuType == menu.TYPE_DIR && !regular.ValidHttp(item.Path) {
			// 获取重定向地址
			redirectPrefix, redirectPath := r.getRouteRedirect(
				cMenus,
				router.Path,
				prefix,
			)
			router.Redirect = redirectPath
			// 子菜单进入递归
			router.Children = r.BuildRouteMenus(cMenus, redirectPrefix)
		} else if item.ParentID == "0" && item.IsFrame == common.STATUS_YES && item.MenuType == menu.TYPE_MENU {
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
			router.Component = menu.COMPONENT_LAYOUT_BASIC
		} else if item.ParentID == "0" && item.IsFrame == common.STATUS_YES && regular.ValidHttp(item.Path) {
			// 父菜单 内部跳转 路径链接
			menuPath := "/" + item.MenuID
			childPath := menuPath + r.getRouterPath(item)
			children := vo.Router{
				Name:      r.getRouteName(item),
				Path:      childPath,
				Component: menu.COMPONENT_LAYOUT_LINK,
				Meta:      r.getRouteMeta(item),
			}
			router.Meta.HideChildInMenu = true
			router.Children = append(router.Children, children)
			router.Name = item.MenuID
			router.Path = menuPath
			router.Redirect = childPath
			router.Component = menu.COMPONENT_LAYOUT_BASIC
		}

		routers = append(routers, router)
	}
	return routers
}

// getRouteName 获取路由名称 路径英文首字母大写
func (r *SysMenuImpl) getRouteName(sysMenu model.SysMenu) string {
	routerName := parse.ConvertToCamelCase(sysMenu.Path)
	// 路径链接
	if regular.ValidHttp(sysMenu.Path) {
		routerName = routerName[:5] + "Link"
	}
	// 拼上菜单ID防止name重名
	return routerName + "_" + sysMenu.MenuID
}

// getRouterPath 获取路由地址
func (r *SysMenuImpl) getRouterPath(sysMenu model.SysMenu) string {
	routerPath := sysMenu.Path

	// 显式路径
	if routerPath == "" || strings.HasPrefix(routerPath, "/") {
		return routerPath
	}

	// 路径链接 内部跳转
	if regular.ValidHttp(routerPath) && sysMenu.IsFrame == common.STATUS_YES {
		routerPath = regular.Replace(routerPath, `/^http(s)?:\/\/+/`, "")
		routerPath = base64.StdEncoding.EncodeToString([]byte(routerPath))
	}

	// 父菜单 内部跳转
	if sysMenu.ParentID == "0" && sysMenu.IsFrame == common.STATUS_YES {
		routerPath = "/" + routerPath
	}

	return routerPath
}

// getComponent 获取组件信息
func (r *SysMenuImpl) getComponent(sysMenu model.SysMenu) string {
	// 内部跳转 路径链接
	if sysMenu.IsFrame == common.STATUS_YES && regular.ValidHttp(sysMenu.Path) {
		return menu.COMPONENT_LAYOUT_LINK
	}

	// 非父菜单 目录类型
	if sysMenu.ParentID != "0" && sysMenu.MenuType == menu.TYPE_DIR {
		return menu.COMPONENT_LAYOUT_BLANK
	}

	// 组件路径 内部跳转 菜单类型
	if sysMenu.Component != "" && sysMenu.IsFrame == common.STATUS_YES && sysMenu.MenuType == menu.TYPE_MENU {
		// 父菜单套外层布局
		if sysMenu.ParentID == "0" {
			return menu.COMPONENT_LAYOUT_BASIC
		}
		return sysMenu.Component
	}

	return menu.COMPONENT_LAYOUT_BASIC
}

// getRouteMeta 获取路由元信息
func (r *SysMenuImpl) getRouteMeta(sysMenu model.SysMenu) vo.RouterMeta {
	meta := vo.RouterMeta{}
	if sysMenu.Icon == "#" {
		meta.Icon = ""
	} else {
		meta.Icon = sysMenu.Icon
	}
	meta.Title = sysMenu.MenuName
	meta.HideChildInMenu = false
	meta.HideInMenu = sysMenu.Visible == common.STATUS_NO
	meta.Cache = sysMenu.IsCache == common.STATUS_YES
	meta.Target = ""

	// 路径链接 非内部跳转
	if regular.ValidHttp(sysMenu.Path) && sysMenu.IsFrame == common.STATUS_NO {
		meta.Target = "_blank"
	}

	return meta
}

// getRouteRedirect 获取路由重定向地址（针对目录）
//
// cMenus 子菜单数组
// routerPath 当前菜单路径
// prefix 菜单重定向路径前缀
func (r *SysMenuImpl) getRouteRedirect(cMenus []model.SysMenu, routerPath string, prefix string) (string, string) {
	redirectPath := ""

	// 重定向为首个显示并启用的子菜单
	var firstChild *model.SysMenu
	for _, item := range cMenus {
		if item.IsFrame == common.STATUS_YES && item.Visible == common.STATUS_YES {
			firstChild = &item
			break
		}
	}

	// 检查内嵌隐藏菜单是否可做重定向
	if firstChild == nil {
		for _, item := range cMenus {
			if item.IsFrame == common.STATUS_YES && item.Visible == common.STATUS_NO && strings.Contains(item.Path, menu.PATH_INLINE) {
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

// parseDataToTree 将数据解析为树结构，构建前端所需要下拉树结构
func (r *SysMenuImpl) parseDataToTree(sysMenus []model.SysMenu) []model.SysMenu {
	// 节点分组
	nodesMap := make(map[string][]model.SysMenu)
	// 节点id
	treeIds := []string{}
	// 树节点
	tree := []model.SysMenu{}

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
		iN := r.parseDataToTreeComponet(node, &nodesMap)
		tree[i] = iN
	}

	return tree
}

// parseDataToTreeComponet 递归函数处理子节点
func (r *SysMenuImpl) parseDataToTreeComponet(node model.SysMenu, nodesMap *map[string][]model.SysMenu) model.SysMenu {
	id := node.MenuID
	children, ok := (*nodesMap)[id]
	if ok {
		node.Children = children
	}
	if len(node.Children) > 0 {
		for i, child := range node.Children {
			icN := r.parseDataToTreeComponet(child, nodesMap)
			node.Children[i] = icN
		}
	}
	return node
}
