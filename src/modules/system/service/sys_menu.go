package service

import (
	"time"

	"mask_api_gin/src/internal/infra/database"
	"mask_api_gin/src/internal/shared/constants"
	"mask_api_gin/src/internal/shared/logger"
	"mask_api_gin/src/internal/shared/regular"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/model/vo"
	"mask_api_gin/src/modules/system/repository"

	"gorm.io/gorm"
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
	rows, err := s.sysMenuRepository.Select(sysMenu, userId)
	if err != nil {
		logger.Errorf("query err => %v", err.Error())
	}
	return rows
}

// FindById 通过ID查询信息
func (s SysMenu) FindById(menuId int64) model.SysMenu {
	if menuId <= 0 {
		return model.SysMenu{}
	}
	menus, err := s.sysMenuRepository.SelectByIds([]int64{menuId})
	if err != nil {
		logger.Errorf("query err => %v", err.Error())
	}
	if len(menus) > 0 {
		return menus[0]
	}
	return model.SysMenu{}
}

// normalizeMenuFields resets fields by menu type (business rule: menu type decides which fields apply)
func (s SysMenu) normalizeMenuFields(sysMenu *model.SysMenu) {
	if sysMenu.MenuIcon == "" {
		sysMenu.MenuIcon = "#"
	}
	switch sysMenu.MenuType {
	case constants.MENU_TYPE_AUTH:
		sysMenu.MenuComponent = ""
		sysMenu.MenuExternal = "0"
		sysMenu.MenuCache = "1"
		sysMenu.MenuVisible = "1"
		sysMenu.MenuPath = ""
		sysMenu.MenuIcon = "#"
	case constants.MENU_TYPE_DIR:
		sysMenu.MenuComponent = ""
		sysMenu.MenuCache = "1"
		sysMenu.MenuPerms = ""
	}
	if sysMenu.MenuType == constants.MENU_TYPE_MENU && regular.ValidHttp(sysMenu.MenuPath) {
		sysMenu.MenuComponent = ""
		sysMenu.MenuPerms = ""
	}
}

// Insert 新增信息
func (s SysMenu) Insert(sysMenu model.SysMenu) int64 {
	if sysMenu.CreateBy != "" {
		ms := time.Now().UnixMilli()
		sysMenu.UpdateBy = sysMenu.CreateBy
		sysMenu.UpdateTime = ms
		sysMenu.CreateTime = ms
	}
	s.normalizeMenuFields(&sysMenu)
	rows, err := s.sysMenuRepository.Insert(sysMenu)
	if err != nil {
		logger.Errorf("insert err => %v", err.Error())
		return 0
	}
	return rows
}

// Update 修改信息
func (s SysMenu) Update(sysMenu model.SysMenu) int64 {
	if sysMenu.UpdateBy != "" {
		sysMenu.UpdateTime = time.Now().UnixMilli()
	}
	s.normalizeMenuFields(&sysMenu)
	rows, err := s.sysMenuRepository.Update(sysMenu)
	if err != nil {
		logger.Errorf("update err => %v", err.Error())
		return 0
	}
	return rows
}

// DeleteById 删除信息
func (s SysMenu) DeleteById(menuId int64) int64 {
	var rows int64
	err := database.Transaction("", func(tx *gorm.DB) error {
		if _, err := s.sysRoleMenuRepository.DeleteByMenuIds(tx, []int64{menuId}); err != nil {
			return err
		}
		r, err := s.sysMenuRepository.DeleteById(tx, menuId)
		rows = r
		return err
	})
	if err != nil {
		return rows
	}
	return rows
}

// ExistChildrenByMenuIdAndStatus 菜单下同状态存在子节点数量
func (s SysMenu) ExistChildrenByMenuIdAndStatus(menuId int64, status string) int64 {
	rows, err := s.sysMenuRepository.ExistChildrenByMenuIdAndStatus(menuId, status)
	if err != nil {
		logger.Errorf("query err => %v", err.Error())
	}
	return rows
}

// ExistRoleByMenuId 菜单分配给的角色数量
func (s SysMenu) ExistRoleByMenuId(menuId int64) int64 {
	rows, err := s.sysRoleMenuRepository.ExistRoleByMenuId(menuId)
	if err != nil {
		logger.Errorf("query err => %v", err.Error())
	}
	return rows
}

// CheckUniqueParentIdByMenuName 检查同级下菜单名称是否唯一
func (s SysMenu) CheckUniqueParentIdByMenuName(parentId int64, menuName string, menuId int64) bool {
	uniqueId, err := s.sysMenuRepository.CheckUnique(model.SysMenu{
		MenuName: menuName,
		ParentId: parentId,
	})
	if err != nil {
		logger.Errorf("query err => %v", err.Error())
	}
	if uniqueId == menuId {
		return true
	}
	return uniqueId == 0
}

// CheckUniqueMenuPath 检查目录和菜单的完整路由地址是否全局唯一
func (s SysMenu) CheckUniqueMenuPath(menuPath string, menuId int64) bool {
	uniqueId, err := s.sysMenuRepository.CheckUniqueMenuPath(menuPath, menuId)
	if err != nil {
		logger.Errorf("query err => %v", err.Error())
	}
	if uniqueId == menuId {
		return true
	}
	return uniqueId == 0
}

// FindPermsByUserId 根据用户ID查询权限标识
func (s SysMenu) FindPermsByUserId(userId int64) []string {
	rows, err := s.sysMenuRepository.SelectPermsByUserId(userId)
	if err != nil {
		logger.Errorf("query err => %v", err.Error())
	}
	return rows
}

// FindByRoleId 根据角色ID查询菜单树信息
func (s SysMenu) FindByRoleId(roleId int64) []int64 {
	roles, err := s.sysRoleRepository.SelectByIds([]int64{roleId})
	if err != nil {
		logger.Errorf("query err => %v", err.Error())
	}
	if len(roles) > 0 {
		role := roles[0]
		if role.ID == roleId {
			menuIds, err := s.sysMenuRepository.SelectByRoleId(
				role.ID,
				role.RoleMenuCheckStrictly == "1",
			)
			if err != nil {
				logger.Errorf("query err => %v", err.Error())
			}
			return menuIds
		}
	}
	return []int64{}
}

// BuildTreeMenusByUserId 根据用户ID查询菜单树状嵌套
func (s SysMenu) BuildTreeMenusByUserId(userId int64) []model.SysMenu {
	arr, err := s.sysMenuRepository.SelectTreeByUserId(userId)
	if err != nil {
		logger.Errorf("query err => %v", err.Error())
	}
	return s.parseDataToTree(arr)
}

// BuildTreeSelectByUserId 根据用户ID查询菜单树状结构
func (s SysMenu) BuildTreeSelectByUserId(sysMenu model.SysMenu, userId int64) []vo.TreeSelect {
	arr, err := s.sysMenuRepository.Select(sysMenu, userId)
	if err != nil {
		logger.Errorf("query err => %v", err.Error())
	}
	menus := s.parseDataToTree(arr)
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
		treeIds = append(treeIds, item.ID)
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
	if children, ok := (*nodesMap)[node.ID]; ok {
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

// BuildAuthMenus 构建当前登录用户可用的语义化菜单树。
func (s SysMenu) BuildAuthMenus(sysMenus []model.SysMenu) []vo.AuthMenu {
	authMenus := make([]vo.AuthMenu, 0, len(sysMenus))
	for _, item := range sysMenus {
		menu := vo.AuthMenu{
			ID:     item.ID,
			Kind:   s.getAuthMenuKind(item),
			Title:  item.MenuName,
			Path:   item.MenuPath,
			Icon:   item.MenuIcon,
			Hidden: item.MenuVisible == constants.STATUS_NO,
			Cache:  item.MenuCache == constants.STATUS_YES,
		}
		if menu.Icon == "#" {
			menu.Icon = ""
		}
		if menu.Kind == "external" {
			menu.Target = "_blank"
		}
		if menu.Kind == "page" {
			menu.Component = item.MenuComponent
		}
		if len(item.Children) > 0 {
			menu.Children = s.BuildAuthMenus(item.Children)
		}
		authMenus = append(authMenus, menu)
	}
	return authMenus
}

func (s SysMenu) getAuthMenuKind(sysMenu model.SysMenu) string {
	if regular.ValidHttp(sysMenu.MenuPath) {
		if sysMenu.MenuExternal == constants.STATUS_YES {
			return "external"
		}
		return "iframe"
	}
	if sysMenu.MenuType == constants.MENU_TYPE_DIR {
		return "directory"
	}
	return "page"
}
