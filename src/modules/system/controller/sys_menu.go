package controller

import (
	"mask_api_gin/src/framework/config"
	"mask_api_gin/src/framework/constants"
	"mask_api_gin/src/framework/response"
	"mask_api_gin/src/framework/utils/ctx"
	"mask_api_gin/src/framework/utils/regular"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/service"

	"fmt"

	"github.com/gin-gonic/gin"
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
func (s SysMenuController) List(c *gin.Context) {
	query := model.SysMenu{}
	if v, ok := c.GetQuery("menuName"); ok {
		query.MenuName = v
	}
	if v, ok := c.GetQuery("statusFlag"); ok {
		query.StatusFlag = v
	}

	userId := ctx.LoginUserToUserID(c)
	if config.IsSystemUser(userId) {
		userId = "0"
	}
	data := s.sysMenuService.Find(query, userId)
	c.JSON(200, response.OkData(data))
}

// Info 菜单信息
//
// GET /:menuId
func (s SysMenuController) Info(c *gin.Context) {
	menuId := c.Param("menuId")
	if menuId == "" {
		c.JSON(400, response.CodeMsg(40010, "bind err: menuId is empty"))
		return
	}

	data := s.sysMenuService.FindById(menuId)
	if data.MenuId == menuId {
		c.JSON(200, response.OkData(data))
		return
	}
	c.JSON(200, response.Err(nil))
}

// Add 菜单新增
//
// POST /
func (s SysMenuController) Add(c *gin.Context) {
	var body model.SysMenu
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", response.FormatBindError(err))
		c.JSON(400, response.CodeMsg(40010, errMsgs))
		return
	}
	if body.MenuId != "" {
		c.JSON(400, response.CodeMsg(40010, "bind err: menuId not is empty"))
		return
	}

	// 目录和菜单检查地址唯一
	if constants.MENU_TYPE_DIR == body.MenuType || constants.MENU_TYPE_MENU == body.MenuType {
		uniqueMenuPath := s.sysMenuService.CheckUniqueParentIdByMenuPath(body.ParentId, body.MenuPath, "")
		if !uniqueMenuPath {
			msg := fmt.Sprintf("菜单新增【%s】失败，菜单路由地址已存在", body.MenuName)
			c.JSON(200, response.ErrMsg(msg))
			return
		}
	}

	// 检查名称唯一
	uniqueMenuName := s.sysMenuService.CheckUniqueParentIdByMenuName(body.ParentId, body.MenuName, "")
	if !uniqueMenuName {
		msg := fmt.Sprintf("菜单新增【%s】失败，菜单名称已存在", body.MenuName)
		c.JSON(200, response.ErrMsg(msg))
		return
	}

	// 外链菜单需要符合网站http(s)开头
	if body.FrameFlag == constants.STATUS_NO && !regular.ValidHttp(body.MenuPath) {
		msg := fmt.Sprintf("菜单新增【%s】失败，非内部地址必须以http(s)://开头", body.MenuName)
		c.JSON(200, response.ErrMsg(msg))
		return
	}

	body.CreateBy = ctx.LoginUserToUserName(c)
	insertId := s.sysMenuService.Insert(body)
	if insertId != "" {
		c.JSON(200, response.OkData(insertId))
		return
	}
	c.JSON(200, response.Err(nil))
}

// Edit 菜单修改
//
// PUT /
func (s SysMenuController) Edit(c *gin.Context) {
	var body model.SysMenu
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", response.FormatBindError(err))
		c.JSON(400, response.CodeMsg(40010, errMsgs))
		return
	}
	if body.MenuId == "" {
		c.JSON(400, response.CodeMsg(40010, "bind err: menuId is empty"))
		return
	}

	// 上级菜单不能选自己
	if body.MenuId == body.ParentId {
		msg := fmt.Sprintf("菜单修改【%s】失败，上级菜单不能选择自己", body.MenuName)
		c.JSON(200, response.ErrMsg(msg))
		return
	}

	// 检查数据是否存在
	menuInfo := s.sysMenuService.FindById(body.MenuId)
	if menuInfo.MenuId != body.MenuId {
		c.JSON(200, response.ErrMsg("没有权限访问菜单数据"))
		return
	}
	// 父级ID不为0是要检查
	if body.ParentId != "0" {
		menuParent := s.sysMenuService.FindById(body.ParentId)
		if menuParent.MenuId != body.ParentId {
			c.JSON(200, response.ErrMsg("没有权限访问菜单数据"))
			return
		}
		// 禁用菜单时检查父菜单是否使用
		if body.StatusFlag == constants.STATUS_YES && menuParent.StatusFlag == constants.STATUS_NO {
			c.JSON(200, response.ErrMsg("上级菜单未启用！"))
			return
		}
	}

	// 目录和菜单检查地址唯一
	if constants.MENU_TYPE_DIR == body.MenuType || constants.MENU_TYPE_MENU == body.MenuType {
		uniqueMenuPath := s.sysMenuService.CheckUniqueParentIdByMenuPath(body.ParentId, body.MenuPath, body.MenuId)
		if !uniqueMenuPath {
			msg := fmt.Sprintf("菜单修改【%s】失败，菜单路由地址已存在", body.MenuName)
			c.JSON(200, response.ErrMsg(msg))
			return
		}
	}

	// 检查名称唯一
	uniqueMenuName := s.sysMenuService.CheckUniqueParentIdByMenuName(body.ParentId, body.MenuName, body.MenuId)
	if !uniqueMenuName {
		msg := fmt.Sprintf("菜单修改【%s】失败，菜单名称已存在", body.MenuName)
		c.JSON(200, response.ErrMsg(msg))
		return
	}

	// 外链菜单需要符合网站http(s)开头
	if body.FrameFlag == constants.STATUS_NO && !regular.ValidHttp(body.MenuPath) {
		msg := fmt.Sprintf("菜单修改【%s】失败，非内部地址必须以http(s)://开头", body.MenuName)
		c.JSON(200, response.ErrMsg(msg))
		return
	}

	// 禁用菜单时检查子菜单是否使用
	if body.StatusFlag == constants.STATUS_NO {
		hasStatus := s.sysMenuService.ExistChildrenByMenuIdAndStatus(body.MenuId, constants.STATUS_YES)
		if hasStatus > 0 {
			msg := fmt.Sprintf("不允许禁用，存在使用子菜单数：%d", hasStatus)
			c.JSON(200, response.ErrMsg(msg))
			return
		}
	}

	menuInfo.ParentId = body.ParentId
	menuInfo.MenuName = body.MenuName
	menuInfo.MenuType = body.MenuType
	menuInfo.MenuSort = body.MenuSort
	menuInfo.MenuPath = body.MenuPath
	menuInfo.Component = body.Component
	menuInfo.FrameFlag = body.FrameFlag
	menuInfo.CacheFlag = body.CacheFlag
	menuInfo.VisibleFlag = body.VisibleFlag
	menuInfo.StatusFlag = body.StatusFlag
	menuInfo.Perms = body.Perms
	menuInfo.Icon = body.Icon
	menuInfo.Remark = body.Remark
	menuInfo.UpdateBy = ctx.LoginUserToUserName(c)
	rows := s.sysMenuService.Update(menuInfo)
	if rows > 0 {
		c.JSON(200, response.Ok(nil))
		return
	}
	c.JSON(200, response.Err(nil))
}

// Remove 菜单删除
//
// DELETE /:menuId
func (s SysMenuController) Remove(c *gin.Context) {
	menuId := c.Param("menuId")
	if menuId == "" {
		c.JSON(400, response.CodeMsg(40010, "bind err: menuId is empty"))
		return
	}

	// 检查数据是否存在
	menu := s.sysMenuService.FindById(menuId)
	if menu.MenuId != menuId {
		c.JSON(200, response.ErrMsg("没有权限访问菜单数据！"))
		return
	}

	// 检查是否存在子菜单
	hasChild := s.sysMenuService.ExistChildrenByMenuIdAndStatus(menuId, "")
	if hasChild > 0 {
		msg := fmt.Sprintf("不允许删除，存在子菜单数：%d", hasChild)
		c.JSON(200, response.ErrMsg(msg))
		return
	}

	// 检查是否分配给角色
	existRole := s.sysMenuService.ExistRoleByMenuId(menuId)
	if existRole > 0 {
		msg := fmt.Sprintf("不允许删除，菜单已分配给角色数：%d", existRole)
		c.JSON(200, response.ErrMsg(msg))
		return
	}

	rows := s.sysMenuService.DeleteById(menuId)
	if rows > 0 {
		msg := fmt.Sprintf("删除成功：%d", rows)
		c.JSON(200, response.OkMsg(msg))
		return
	}
	c.JSON(200, response.Err(nil))
}

// Tree 菜单树结构列表
//
// GET /tree
func (s SysMenuController) Tree(c *gin.Context) {
	query := model.SysMenu{}
	if v, ok := c.GetQuery("menuName"); ok {
		query.MenuName = v
	}
	if v, ok := c.GetQuery("status"); ok {
		query.StatusFlag = v
	}

	userId := ctx.LoginUserToUserID(c)
	if config.IsSystemUser(userId) {
		userId = "0"
	}
	data := s.sysMenuService.BuildTreeSelectByUserId(query, userId)
	c.JSON(200, response.OkData(data))

}

// TreeRole 菜单树结构列表（指定角色）
//
// GET /tree/role/:roleId
func (s SysMenuController) TreeRole(c *gin.Context) {
	roleId := c.Param("roleId")
	if roleId == "" {
		c.JSON(400, response.CodeMsg(40010, "bind err: roleId is empty"))
		return
	}

	query := model.SysMenu{}
	if v, ok := c.GetQuery("menuName"); ok {
		query.MenuName = v
	}
	if v, ok := c.GetQuery("status"); ok {
		query.StatusFlag = v
	}

	userId := ctx.LoginUserToUserID(c)
	if config.IsSystemUser(userId) {
		userId = "0"
	}
	menuTreeSelect := s.sysMenuService.BuildTreeSelectByUserId(query, userId)
	checkedKeys := s.sysMenuService.FindByRoleId(roleId)
	c.JSON(200, response.OkData(map[string]any{
		"menus":       menuTreeSelect,
		"checkedKeys": checkedKeys,
	}))
}
