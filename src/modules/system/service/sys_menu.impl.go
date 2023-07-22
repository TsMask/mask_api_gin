package service

import (
	"encoding/base64"
	"mask_api_gin/src/framework/constants/common"
	menuConstants "mask_api_gin/src/framework/constants/menu"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/framework/utils/regular"
	"mask_api_gin/src/framework/vo"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/repository"
	"strings"
)

// SysMenuImpl 菜单 数据层处理
var SysMenuImpl = &sysMenuImpl{
	sysMenuRepository:     repository.SysMenuImpl,
	sysRoleMenuRepository: repository.SysRoleMenuImpl,
	sysRoleRepository:     repository.SysRoleImpl,
}

type sysMenuImpl struct {
	// 菜单服务
	sysMenuRepository repository.ISysMenu
	// 角色与菜单关联服务
	sysRoleMenuRepository repository.ISysRoleMenu
	// 角色服务
	sysRoleRepository repository.ISysRole
}

// SelectMenuList 查询系统菜单列表
func (r *sysMenuImpl) SelectMenuList(sysMenu model.SysMenu, userId string) []model.SysMenu {
	return r.sysMenuRepository.SelectMenuList(sysMenu, userId)
}

// SelectMenuPermsByUserId 根据用户ID查询权限
func (r *sysMenuImpl) SelectMenuPermsByUserId(userId string) []string {
	return r.sysMenuRepository.SelectMenuPermsByUserId(userId)
}

// SelectMenuTreeByUserId 根据用户ID查询菜单
func (r *sysMenuImpl) SelectMenuTreeByUserId(userId string) []model.SysMenu {
	sysMenus := r.sysMenuRepository.SelectMenuTreeByUserId(userId)
	return parseDataToTree(sysMenus)
}

// SelectMenuTreeSelectByUserId 根据用户ID查询菜单树结构信息
func (r *sysMenuImpl) SelectMenuTreeSelectByUserId(sysMenu model.SysMenu, userId string) []vo.TreeSelect {
	sysMenus := r.sysMenuRepository.SelectMenuList(sysMenu, userId)
	menus := parseDataToTree(sysMenus)
	tree := make([]vo.TreeSelect, 0)
	for _, menu := range menus {
		tree = append(tree, vo.SysMenuTreeSelect(menu))
	}
	return tree
}

// SelectMenuListByRoleId 根据角色ID查询菜单树信息
func (r *sysMenuImpl) SelectMenuListByRoleId(roleId string) []string {
	role := r.sysRoleRepository.SelectRoleById(roleId)
	if role.RoleID != roleId {
		return []string{}
	}
	return r.sysMenuRepository.SelectMenuListByRoleId(
		role.RoleID,
		role.MenuCheckStrictly == "1",
	)
}

// SelectMenuById 根据菜单ID查询信息
func (r *sysMenuImpl) SelectMenuById(menuId string) model.SysMenu {
	if menuId == "" {
		return model.SysMenu{}
	}
	menus := r.sysMenuRepository.SelectMenuByIds([]string{menuId})
	if len(menus) > 0 {
		return menus[0]
	}
	return model.SysMenu{}
}

// HasChildByMenuId 存在菜单子节点数量
func (r *sysMenuImpl) HasChildByMenuId(menuId string) int64 {
	return r.sysMenuRepository.HasChildByMenuId(menuId)
}

// CheckMenuExistRole 查询菜单是否存在角色
func (r *sysMenuImpl) CheckMenuExistRole(menuId string) int64 {
	return r.sysRoleMenuRepository.CheckMenuExistRole(menuId)
}

// InsertMenu 新增菜单信息
func (r *sysMenuImpl) InsertMenu(sysMenu model.SysMenu) string {
	return r.sysMenuRepository.InsertMenu(sysMenu)
}

// UpdateMenu 修改菜单信息
func (r *sysMenuImpl) UpdateMenu(sysMenu model.SysMenu) int64 {
	return r.sysMenuRepository.UpdateMenu(sysMenu)
}

// DeleteMenuById 删除菜单管理信息
func (r *sysMenuImpl) DeleteMenuById(menuId string) int64 {
	return r.sysMenuRepository.DeleteMenuById(menuId)
}

// CheckUniqueMenuName 校验菜单名称是否唯一
func (r *sysMenuImpl) CheckUniqueMenuName(menuName, parentId, menuId string) bool {
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
func (r *sysMenuImpl) CheckUniqueMenuPath(path, menuId string) bool {
	uniqueId := r.sysMenuRepository.CheckUniqueMenu(model.SysMenu{
		Path: path,
	})
	if uniqueId == menuId {
		return true
	}
	return uniqueId == ""
}

// BuildRouteMenus 构建前端路由所需要的菜单
func (r *sysMenuImpl) BuildRouteMenus(sysMenus []model.SysMenu, prefix string) []vo.Router {
	routers := []vo.Router{}
	for _, item := range sysMenus {
		router := vo.Router{}
		router.Name = r.getRouteName(item)
		router.Path = r.getRouterPath(item)
		router.Component = r.getComponent(item)
		router.Meta = r.getRouteMeta(item)

		// 非路径链接 子项菜单目录
		cMenus := item.Children
		if len(cMenus) > 0 && !regular.ValidHttp(item.Path) && item.MenuType == menuConstants.TYPE_DIR {
			// 获取重定向地址
			redirectPrefix, redirectPath := r.getRouteRedirect(
				cMenus,
				router.Path,
				prefix,
			)
			router.Redirect = redirectPath
			// 子菜单进入递归
			router.Children = r.BuildRouteMenus(cMenus, redirectPrefix)
		}
		routers = append(routers, router)
	}
	return routers
}

// 获取路由名称
//
// 路径英文首字母大写
//
// menu 菜单信息
func (r *sysMenuImpl) getRouteName(menu model.SysMenu) string {
	routerName := parse.FirstUpper(menu.Path)
	// 路径链接
	if regular.ValidHttp(menu.Path) {
		return routerName[:5] + "Link" + menu.MenuID
	}
	return routerName
}

// 获取路由地址
//
// menu 菜单信息
func (r *sysMenuImpl) getRouterPath(menu model.SysMenu) string {
	routerPath := menu.Path

	// 显式路径
	if strings.HasPrefix(routerPath, "/") {
		return routerPath
	}

	menuType := menu.MenuType == menuConstants.TYPE_DIR || menu.MenuType == menuConstants.TYPE_MENU

	// 路径链接
	if regular.ValidHttp(routerPath) {
		// 内部跳转 非父菜单 目录类型或菜单类型
		if menu.IsFrame == common.STATUS_YES && menu.ParentID != "0" && menuType {
			routerPath = regular.Replace(routerPath, `/^http(s)?:\/\/+/`, "")
			return base64.StdEncoding.EncodeToString([]byte(routerPath))
		}
		// 非内部跳转
		return routerPath
	}

	// 父菜单 目录类型或菜单类型
	if menu.ParentID == "0" && menuType {
		routerPath = "/" + routerPath
	}

	return routerPath
}

// 获取组件信息
//
// menu 菜单信息
func (r *sysMenuImpl) getComponent(menu model.SysMenu) string {
	menuType := menu.MenuType == menuConstants.TYPE_DIR || menu.MenuType == menuConstants.TYPE_MENU

	// 路径链接 非父菜单 目录类型或菜单类型
	if regular.ValidHttp(menu.Path) && menu.ParentID != "0" && menuType {
		return menuConstants.COMPONENT_LAYOUT_LINK
	}

	// 非父菜单 目录类型
	if menu.ParentID != "0" && menu.MenuType == menuConstants.TYPE_DIR {
		return menuConstants.COMPONENT_LAYOUT_BLANK
	}

	// 菜单类型 内部跳转 有组件路径
	if menu.MenuType == menuConstants.TYPE_MENU && menu.IsFrame == common.STATUS_YES && menu.Component != "" {
		return menu.Component
	}

	return menuConstants.COMPONENT_LAYOUT_BASIC
}

// 获取路由元信息
//
// menu 菜单信息
func (r *sysMenuImpl) getRouteMeta(menu model.SysMenu) vo.RouterMeta {
	meta := vo.RouterMeta{}
	if menu.Icon == "#" {
		meta.Icon = ""
	} else {
		meta.Icon = menu.Icon
	}
	meta.Title = menu.MenuName
	meta.Hide = menu.Visible == common.STATUS_NO
	meta.Cache = menu.IsCache == common.STATUS_YES
	meta.Target = ""

	// 路径链接
	if regular.ValidHttp(menu.Path) {
		menuType := menu.MenuType == menuConstants.TYPE_DIR || menu.MenuType == menuConstants.TYPE_MENU

		// 内部跳转 父菜单 目录类型或菜单类型
		if menu.IsFrame == common.STATUS_YES && menu.ParentID == "0" && menuType {
			meta.Target = "_self"
		}
		// 非内部跳转
		if menu.IsFrame == common.STATUS_NO {
			meta.Target = "_blank"
		}
	}

	return meta
}

// 获取路由重定向地址（针对目录）
//
// cMenus 子菜单数组
// routerPath 当前菜单路径
// prefix 菜单重定向路径前缀
func (r *sysMenuImpl) getRouteRedirect(cMenus []model.SysMenu, routerPath string, prefix string) (string, string) {
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
			if item.IsFrame == common.STATUS_YES && item.Visible == common.STATUS_NO && strings.Contains(item.Path, menuConstants.PATH_INLINE) {
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
func parseDataToTree(sysMenus []model.SysMenu) []model.SysMenu {
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
		iN := componet(node, &nodesMap)
		tree[i] = iN
	}

	return tree
}

// componet 递归函数处理子节点
func componet(node model.SysMenu, nodesMap *map[string][]model.SysMenu) model.SysMenu {
	id := node.MenuID
	children, ok := (*nodesMap)[id]
	if ok {
		node.Children = children
	}
	if len(node.Children) > 0 {
		for i, child := range node.Children {
			icN := componet(child, nodesMap)
			node.Children[i] = icN
		}
	}
	return node
}
