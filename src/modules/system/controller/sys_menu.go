package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"mask_api_gin/src/internal/config"
	"mask_api_gin/src/internal/shared/constants"
	"mask_api_gin/src/internal/shared/parse"
	"mask_api_gin/src/internal/shared/regular"
	"mask_api_gin/src/internal/web/reqctx"
	"mask_api_gin/src/internal/web/resp"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/service"
)

// NewSysMenu 实例化控制层
var NewSysMenu = &SysMenuController{
	sysMenuService: service.NewSysMenu,
}

// SysMenuController 菜单信息 控制层处理
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
	if v, ok := c.GetQuery("menu_name"); ok {
		query.MenuName = v
	}
	if v, ok := c.GetQuery("menu_status"); ok {
		query.MenuStatus = v
	}

	userId := reqctx.LoginUserToUserID(c)
	if config.IsSystemUser(userId) {
		userId = 0
	}
	data := s.sysMenuService.Find(query, userId)
	c.JSON(200, resp.OkData(data))
}

// Info 菜单信息
//
// GET /:id
func (s SysMenuController) Info(c *gin.Context) {
	id := parse.Number(c.Param("id"))
	if id <= 0 {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: id is empty"))
		return
	}

	data := s.sysMenuService.FindById(id)
	if data.ID == id {
		c.JSON(200, resp.OkData(data))
		return
	}
	c.JSON(200, resp.Err(nil))
}

// Add 菜单新增
//
// POST /
func (s SysMenuController) Add(c *gin.Context) {
	var body model.SysMenu
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_PARSER, errMsgs))
		return
	}
	if body.ID > 0 {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: menuId not is empty"))
		return
	}

	// 目录和菜单检查地址唯一
	if constants.MENU_TYPE_DIR == body.MenuType || constants.MENU_TYPE_MENU == body.MenuType {
		uniqueMenuPath := s.sysMenuService.CheckUniqueParentIdByMenuPath(body.ParentId, body.MenuPath, 0)
		if !uniqueMenuPath {
			msg := fmt.Sprintf("菜单新增【%s】失败，菜单路由地址已存在", body.MenuName)
			c.JSON(200, resp.ErrMsg(msg))
			return
		}
	}

	// 检查名称唯一
	uniqueMenuName := s.sysMenuService.CheckUniqueParentIdByMenuName(body.ParentId, body.MenuName, 0)
	if !uniqueMenuName {
		msg := fmt.Sprintf("菜单新增【%s】失败，菜单名称已存在", body.MenuName)
		c.JSON(200, resp.ErrMsg(msg))
		return
	}

	// 外链菜单需要符合网站http(s)开头
	if body.MenuExternal == constants.STATUS_YES && !regular.ValidHttp(body.MenuPath) {
		msg := fmt.Sprintf("菜单新增【%s】失败，非内部地址必须以http(s)://开头", body.MenuName)
		c.JSON(200, resp.ErrMsg(msg))
		return
	}

	body.CreateBy = reqctx.LoginUserToUserName(c)
	insertId := s.sysMenuService.Insert(body)
	if insertId > 0 {
		c.JSON(200, resp.OkData(insertId))
		return
	}
	c.JSON(200, resp.Err(nil))
}

// Edit 菜单修改
//
// PUT /
func (s SysMenuController) Edit(c *gin.Context) {
	var body model.SysMenu
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_PARSER, errMsgs))
		return
	}
	if body.ID <= 0 {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: menuId is empty"))
		return
	}

	// 上级菜单不能选自己
	if body.ID == body.ParentId {
		msg := fmt.Sprintf("菜单修改【%s】失败，上级菜单不能选择自己", body.MenuName)
		c.JSON(200, resp.ErrMsg(msg))
		return
	}

	// 检查数据是否存在
	menuInfo := s.sysMenuService.FindById(body.ID)
	if menuInfo.ID != body.ID {
		c.JSON(200, resp.ErrMsg("没有权限访问菜单数据"))
		return
	}
	// 父级ID不为0是要检查
	if body.ParentId > 0 {
		menuParent := s.sysMenuService.FindById(body.ParentId)
		if menuParent.ID != body.ParentId {
			c.JSON(200, resp.ErrMsg("没有权限访问菜单数据"))
			return
		}
		// 禁用菜单时检查父菜单是否使用
		if body.MenuStatus == constants.STATUS_YES && menuParent.MenuStatus == constants.STATUS_NO {
			c.JSON(200, resp.ErrMsg("上级菜单未启用！"))
			return
		}
	}

	// 目录和菜单检查地址唯一
	if constants.MENU_TYPE_DIR == body.MenuType || constants.MENU_TYPE_MENU == body.MenuType {
		uniqueMenuPath := s.sysMenuService.CheckUniqueParentIdByMenuPath(body.ParentId, body.MenuPath, body.ID)
		if !uniqueMenuPath {
			msg := fmt.Sprintf("菜单修改【%s】失败，菜单路由地址已存在", body.MenuName)
			c.JSON(200, resp.ErrMsg(msg))
			return
		}
	}

	// 检查名称唯一
	uniqueMenuName := s.sysMenuService.CheckUniqueParentIdByMenuName(body.ParentId, body.MenuName, body.ID)
	if !uniqueMenuName {
		msg := fmt.Sprintf("菜单修改【%s】失败，菜单名称已存在", body.MenuName)
		c.JSON(200, resp.ErrMsg(msg))
		return
	}

	// 外链菜单需要符合网站http(s)开头
	if body.MenuExternal == constants.STATUS_YES && !regular.ValidHttp(body.MenuPath) {
		msg := fmt.Sprintf("菜单修改【%s】失败，非内部地址必须以http(s)://开头", body.MenuName)
		c.JSON(200, resp.ErrMsg(msg))
		return
	}

	// 禁用菜单时检查子菜单是否使用
	if body.MenuStatus == constants.STATUS_NO {
		hasStatus := s.sysMenuService.ExistChildrenByMenuIdAndStatus(body.ID, constants.STATUS_YES)
		if hasStatus > 0 {
			msg := fmt.Sprintf("不允许禁用，存在使用子菜单数：%d", hasStatus)
			c.JSON(200, resp.ErrMsg(msg))
			return
		}
	}

	menuInfo.ParentId = body.ParentId
	menuInfo.MenuName = body.MenuName
	menuInfo.MenuType = body.MenuType
	menuInfo.MenuSort = body.MenuSort
	menuInfo.MenuPath = body.MenuPath
	menuInfo.MenuComponent = body.MenuComponent
	menuInfo.MenuExternal = body.MenuExternal
	menuInfo.MenuCache = body.MenuCache
	menuInfo.MenuVisible = body.MenuVisible
	menuInfo.MenuStatus = body.MenuStatus
	menuInfo.MenuPerms = body.MenuPerms
	menuInfo.MenuIcon = body.MenuIcon
	menuInfo.Remark = body.Remark
	menuInfo.UpdateBy = reqctx.LoginUserToUserName(c)
	rows := s.sysMenuService.Update(menuInfo)
	if rows > 0 {
		c.JSON(200, resp.Ok(nil))
		return
	}
	c.JSON(200, resp.Err(nil))
}

// Remove 菜单删除
//
// DELETE /:menu_id
func (s SysMenuController) Remove(c *gin.Context) {
	id := parse.Number(c.Param("id"))
	if id <= 0 {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: id is empty"))
		return
	}

	// 检查数据是否存在
	menu := s.sysMenuService.FindById(id)
	if menu.ID != id {
		c.JSON(200, resp.ErrMsg("没有权限访问菜单数据！"))
		return
	}

	// 检查是否存在子菜单
	hasChild := s.sysMenuService.ExistChildrenByMenuIdAndStatus(id, "")
	if hasChild > 0 {
		msg := fmt.Sprintf("不允许删除，存在子菜单数：%d", hasChild)
		c.JSON(200, resp.ErrMsg(msg))
		return
	}

	// 检查是否分配给角色
	existRole := s.sysMenuService.ExistRoleByMenuId(id)
	if existRole > 0 {
		msg := fmt.Sprintf("不允许删除，菜单已分配给角色数：%d", existRole)
		c.JSON(200, resp.ErrMsg(msg))
		return
	}

	rows := s.sysMenuService.DeleteById(id)
	if rows > 0 {
		msg := fmt.Sprintf("删除成功：%d", rows)
		c.JSON(200, resp.OkMsg(msg))
		return
	}
	c.JSON(200, resp.Err(nil))
}

// Tree 菜单树结构列表
//
// GET /tree
func (s SysMenuController) Tree(c *gin.Context) {
	query := model.SysMenu{}
	if v, ok := c.GetQuery("menu_name"); ok {
		query.MenuName = v
	}
	if v, ok := c.GetQuery("menu_status"); ok {
		query.MenuStatus = v
	}

	userId := reqctx.LoginUserToUserID(c)
	if config.IsSystemUser(userId) {
		userId = 0
	}
	trees := s.sysMenuService.BuildTreeSelectByUserId(query, userId)
	c.JSON(200, resp.OkData(trees))
}

// TreeRole 菜单树结构列表（指定角色）
//
// GET /tree/role/:role_id
func (s SysMenuController) TreeRole(c *gin.Context) {
	roleId := parse.Number(c.Param("role_id"))
	if roleId <= 0 {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: roleId is empty"))
		return
	}

	sysMenu := model.SysMenu{}
	if v, ok := c.GetQuery("menu_name"); ok {
		sysMenu.MenuName = v
	}
	if v, ok := c.GetQuery("menu_status"); ok {
		sysMenu.MenuStatus = v
	}

	userId := reqctx.LoginUserToUserID(c)
	if config.IsSystemUser(userId) {
		userId = 0
	}
	menuTreeSelect := s.sysMenuService.BuildTreeSelectByUserId(sysMenu, userId)
	checkedKeys := s.sysMenuService.FindByRoleId(roleId)
	c.JSON(200, resp.OkData(map[string]any{
		"menus":       menuTreeSelect,
		"checkedKeys": checkedKeys,
	}))
}
