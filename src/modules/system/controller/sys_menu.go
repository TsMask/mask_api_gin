package controller

import (
	"fmt"
	"mask_api_gin/src/framework/config"
	constMenu "mask_api_gin/src/framework/constants/menu"
	constSystem "mask_api_gin/src/framework/constants/system"
	"mask_api_gin/src/framework/utils/ctx"
	"mask_api_gin/src/framework/utils/regular"
	"mask_api_gin/src/framework/vo/result"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/service"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// NewSysMenu 实例化控制层
var NewSysMenu = &SysMenuController{
	sysMenuService: service.NewSysMenu,
}

// SysMenuController 菜单信息
//
// PATH /system/menu
type SysMenuController struct {
	sysMenuService *service.SysMenu // 菜单服务
}

// List 菜单列表
//
// GET /list
func (s *SysMenuController) List(c *gin.Context) {
	query := model.SysMenu{}
	if v, ok := c.GetQuery("menuName"); ok {
		query.MenuName = v
	}
	if v, ok := c.GetQuery("status"); ok {
		query.Status = v
	}

	userId := ctx.LoginUserToUserID(c)
	if config.IsSysAdmin(userId) {
		userId = "*"
	}
	data := s.sysMenuService.Find(query, userId)
	c.JSON(200, result.OkData(data))
}

// Info 菜单信息
//
// GET /:menuId
func (s *SysMenuController) Info(c *gin.Context) {
	menuId := c.Param("menuId")
	if menuId == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}
	data := s.sysMenuService.FindById(menuId)
	if data.MenuID == menuId {
		c.JSON(200, result.OkData(data))
		return
	}
	c.JSON(200, result.Err(nil))
}

// Add 菜单新增
//
// POST /
func (s *SysMenuController) Add(c *gin.Context) {
	var body model.SysMenu
	err := c.ShouldBindBodyWith(&body, binding.JSON)
	if err != nil || body.MenuID != "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 目录和菜单检查地址唯一
	if constMenu.TYPE_DIR == body.MenuType || constMenu.TYPE_MENU == body.MenuType {
		uniqueMenuPath := s.sysMenuService.CheckUniqueParentIdByMenuPath(body.ParentID, body.Path, "")
		if !uniqueMenuPath {
			msg := fmt.Sprintf("菜单新增【%s】失败，菜单路由地址已存在", body.MenuName)
			c.JSON(200, result.ErrMsg(msg))
			return
		}
	}

	// 检查名称唯一
	uniqueMenuName := s.sysMenuService.CheckUniqueParentIdByMenuName(body.ParentID, body.MenuName, "")
	if !uniqueMenuName {
		msg := fmt.Sprintf("菜单新增【%s】失败，菜单名称已存在", body.MenuName)
		c.JSON(200, result.ErrMsg(msg))
		return
	}

	// 外链菜单需要符合网站http(s)开头
	if body.IsFrame == constSystem.STATUS_NO && !regular.ValidHttp(body.Path) {
		msg := fmt.Sprintf("菜单新增【%s】失败，非内部地址必须以http(s)://开头", body.MenuName)
		c.JSON(200, result.ErrMsg(msg))
		return
	}

	body.CreateBy = ctx.LoginUserToUserName(c)
	insertId := s.sysMenuService.Insert(body)
	if insertId != "" {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// Edit 菜单修改
//
// PUT /
func (s *SysMenuController) Edit(c *gin.Context) {
	var body model.SysMenu
	err := c.ShouldBindBodyWith(&body, binding.JSON)
	if err != nil || body.MenuID == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 上级菜单不能选自己
	if body.MenuID == body.ParentID {
		msg := fmt.Sprintf("菜单修改【%s】失败，上级菜单不能选择自己", body.MenuName)
		c.JSON(200, result.ErrMsg(msg))
		return
	}

	// 检查数据是否存在
	menuInfo := s.sysMenuService.FindById(body.MenuID)
	if menuInfo.MenuID != body.MenuID {
		c.JSON(200, result.ErrMsg("没有权限访问菜单数据"))
		return
	}
	// 父级ID不为0是要检查
	if body.ParentID != "0" {
		menuParent := s.sysMenuService.FindById(body.ParentID)
		if menuParent.MenuID != body.ParentID {
			c.JSON(200, result.ErrMsg("没有权限访问菜单数据"))
			return
		}
		// 禁用菜单时检查父菜单是否使用
		if body.Status == constSystem.STATUS_YES && menuParent.Status == constSystem.STATUS_NO {
			c.JSON(200, result.ErrMsg("上级菜单未启用！"))
			return
		}
	}

	// 目录和菜单检查地址唯一
	if constMenu.TYPE_DIR == body.MenuType || constMenu.TYPE_MENU == body.MenuType {
		uniqueMenuPath := s.sysMenuService.CheckUniqueParentIdByMenuPath(body.ParentID, body.Path, body.MenuID)
		if !uniqueMenuPath {
			msg := fmt.Sprintf("菜单修改【%s】失败，菜单路由地址已存在", body.MenuName)
			c.JSON(200, result.ErrMsg(msg))
			return
		}
	}

	// 检查名称唯一
	uniqueMenuName := s.sysMenuService.CheckUniqueParentIdByMenuName(body.ParentID, body.MenuName, body.MenuID)
	if !uniqueMenuName {
		msg := fmt.Sprintf("菜单修改【%s】失败，菜单名称已存在", body.MenuName)
		c.JSON(200, result.ErrMsg(msg))
		return
	}

	// 外链菜单需要符合网站http(s)开头
	if body.IsFrame == constSystem.STATUS_NO && !regular.ValidHttp(body.Path) {
		msg := fmt.Sprintf("菜单修改【%s】失败，非内部地址必须以http(s)://开头", body.MenuName)
		c.JSON(200, result.ErrMsg(msg))
		return
	}

	// 禁用菜单时检查子菜单是否使用
	if body.Status == constSystem.STATUS_NO {
		hasStatus := s.sysMenuService.ExistChildrenByMenuIdAndStatus(body.MenuID, constSystem.STATUS_YES)
		if hasStatus > 0 {
			msg := fmt.Sprintf("不允许禁用，存在使用子菜单数：%d", hasStatus)
			c.JSON(200, result.ErrMsg(msg))
			return
		}
	}

	body.UpdateBy = ctx.LoginUserToUserName(c)
	rows := s.sysMenuService.Update(body)
	if rows > 0 {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// Remove 菜单删除
//
// DELETE /:menuId
func (s *SysMenuController) Remove(c *gin.Context) {
	menuId := c.Param("menuId")
	if menuId == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 检查数据是否存在
	menu := s.sysMenuService.FindById(menuId)
	if menu.MenuID != menuId {
		c.JSON(200, result.ErrMsg("没有权限访问菜单数据！"))
		return
	}

	// 检查是否存在子菜单
	hasChild := s.sysMenuService.ExistChildrenByMenuIdAndStatus(menuId, "")
	if hasChild > 0 {
		msg := fmt.Sprintf("不允许删除，存在子菜单数：%d", hasChild)
		c.JSON(200, result.ErrMsg(msg))
		return
	}

	// 检查是否分配给角色
	existRole := s.sysMenuService.ExistRoleByMenuId(menuId)
	if existRole > 0 {
		msg := fmt.Sprintf("不允许删除，菜单已分配给角色数：%d", existRole)
		c.JSON(200, result.ErrMsg(msg))
		return
	}

	rows := s.sysMenuService.DeleteById(menuId)
	if rows > 0 {
		msg := fmt.Sprintf("删除成功：%d", rows)
		c.JSON(200, result.OkMsg(msg))
		return
	}
	c.JSON(200, result.Err(nil))
}

// TreeSelect 菜单树结构列表
//
// GET /treeSelect
func (s *SysMenuController) TreeSelect(c *gin.Context) {
	query := model.SysMenu{}
	if v, ok := c.GetQuery("menuName"); ok {
		query.MenuName = v
	}
	if v, ok := c.GetQuery("status"); ok {
		query.Status = v
	}

	userId := ctx.LoginUserToUserID(c)
	if config.IsSysAdmin(userId) {
		userId = "*"
	}
	data := s.sysMenuService.BuildTreeSelectByUserId(query, userId)
	c.JSON(200, result.OkData(data))

}

// RoleMenuTreeSelect 菜单树结构列表（指定角色）
//
// GET /roleMenuTreeSelect/:roleId
func (s *SysMenuController) RoleMenuTreeSelect(c *gin.Context) {
	roleId := c.Param("roleId")
	if roleId == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	query := model.SysMenu{}
	if v, ok := c.GetQuery("menuName"); ok {
		query.MenuName = v
	}
	if v, ok := c.GetQuery("status"); ok {
		query.Status = v
	}

	userId := ctx.LoginUserToUserID(c)
	if config.IsSysAdmin(userId) {
		userId = "*"
	}
	menuTreeSelect := s.sysMenuService.BuildTreeSelectByUserId(query, userId)
	checkedKeys := s.sysMenuService.FindByRoleId(roleId)
	c.JSON(200, result.OkData(map[string]any{
		"menus":       menuTreeSelect,
		"checkedKeys": checkedKeys,
	}))
}
