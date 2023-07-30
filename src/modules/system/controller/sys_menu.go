package controller

import (
	"fmt"
	"mask_api_gin/src/framework/config"
	"mask_api_gin/src/framework/constants/common"
	"mask_api_gin/src/framework/constants/menu"
	"mask_api_gin/src/framework/utils/ctx"
	"mask_api_gin/src/framework/utils/regular"
	"mask_api_gin/src/framework/vo/result"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/service"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// 菜单信息
//
// PATH /system/menu
var SysMenu = &sysMenu{
	sysMenuService: service.SysMenuImpl,
}

type sysMenu struct {
	// 菜单服务
	sysMenuService service.ISysMenu
}

// 菜单列表
//
// GET /list
func (s *sysMenu) List(c *gin.Context) {
	query := model.SysMenu{}
	if v, ok := c.GetQuery("menuName"); ok {
		query.MenuName = v
	}
	if v, ok := c.GetQuery("status"); ok {
		query.Status = v
	}

	userId := ctx.LoginUserToUserID(c)
	if config.IsAdmin(userId) {
		userId = "*"
	}
	data := s.sysMenuService.SelectMenuList(query, userId)
	c.JSON(200, result.OkData(data))
}

// 菜单信息
//
// GET /:menuId
func (s *sysMenu) Info(c *gin.Context) {
	menuId := c.Param("menuId")
	if menuId == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}
	data := s.sysMenuService.SelectMenuById(menuId)
	if data.MenuID == menuId {
		c.JSON(200, result.OkData(data))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 菜单新增
//
// POST /
func (s *sysMenu) Add(c *gin.Context) {
	var body model.SysMenu
	err := c.ShouldBindBodyWith(&body, binding.JSON)
	if err != nil || body.MenuID != "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 目录和菜单检查地址唯一
	if menu.TYPE_DIR == body.MenuType || menu.TYPE_MENU == body.MenuType {
		uniqueNenuPath := s.sysMenuService.CheckUniqueMenuPath(body.Path, "")
		if !uniqueNenuPath {
			msg := fmt.Sprintf("菜单新增【%s】失败，菜单路由地址已存在", body.MenuName)
			c.JSON(200, result.ErrMsg(msg))
			return
		}
	}

	// 检查名称唯一
	uniqueNenuName := s.sysMenuService.CheckUniqueMenuName(body.MenuName, body.ParentID, "")
	if !uniqueNenuName {
		msg := fmt.Sprintf("菜单新增【%s】失败，菜单名称已存在", body.MenuName)
		c.JSON(200, result.ErrMsg(msg))
		return
	}

	// 外链菜单需要符合网站http(s)开头
	if body.IsFrame == common.STATUS_NO && !regular.ValidHttp(body.Path) {
		msg := fmt.Sprintf("菜单新增【%s】失败，非内部地址必须以http(s)://开头", body.MenuName)
		c.JSON(200, result.ErrMsg(msg))
		return
	}

	body.CreateBy = ctx.LoginUserToUserName(c)
	insertId := s.sysMenuService.InsertMenu(body)
	if insertId != "" {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 菜单修改
//
// PUT /
func (s *sysMenu) Edit(c *gin.Context) {
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
	menuInfo := s.sysMenuService.SelectMenuById(body.MenuID)
	if menuInfo.MenuID != body.MenuID {
		c.JSON(200, result.ErrMsg("没有权限访问菜单数据"))
		return
	}
	// 父级ID不为0是要检查
	if body.ParentID != "0" {
		menuParent := s.sysMenuService.SelectMenuById(body.ParentID)
		if menuParent.MenuID != body.ParentID {
			c.JSON(200, result.ErrMsg("没有权限访问菜单数据"))
			return
		}
	}

	// 目录和菜单检查地址唯一
	if menu.TYPE_DIR == body.MenuType || menu.TYPE_MENU == body.MenuType {
		uniqueNenuPath := s.sysMenuService.CheckUniqueMenuPath(body.Path, body.MenuID)
		if !uniqueNenuPath {
			msg := fmt.Sprintf("菜单修改【%s】失败，菜单路由地址已存在", body.MenuName)
			c.JSON(200, result.ErrMsg(msg))
			return
		}
	}

	// 检查名称唯一
	uniqueNenuName := s.sysMenuService.CheckUniqueMenuName(body.MenuName, body.ParentID, body.MenuID)
	if !uniqueNenuName {
		msg := fmt.Sprintf("菜单修改【%s】失败，菜单名称已存在", body.MenuName)
		c.JSON(200, result.ErrMsg(msg))
		return
	}

	// 外链菜单需要符合网站http(s)开头
	if body.IsFrame == common.STATUS_NO && !regular.ValidHttp(body.Path) {
		msg := fmt.Sprintf("菜单修改【%s】失败，非内部地址必须以http(s)://开头", body.MenuName)
		c.JSON(200, result.ErrMsg(msg))
		return
	}

	body.UpdateBy = ctx.LoginUserToUserName(c)
	rows := s.sysMenuService.UpdateMenu(body)
	if rows > 0 {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 菜单删除
//
// DELETE /:menuId
func (s *sysMenu) Remove(c *gin.Context) {
	menuId := c.Param("menuId")
	if menuId == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 检查数据是否存在
	menu := s.sysMenuService.SelectMenuById(menuId)
	if menu.MenuID != menuId {
		c.JSON(200, result.ErrMsg("没有权限访问菜单数据！"))
		return
	}

	// 检查是否存在子菜单
	hasChild := s.sysMenuService.HasChildByMenuId(menuId)
	if hasChild > 0 {
		msg := fmt.Sprintf("不允许删除，存在子菜单数：%d", hasChild)
		c.JSON(200, result.ErrMsg(msg))
		return
	}

	// 检查是否分配给角色
	existRole := s.sysMenuService.CheckMenuExistRole(menuId)
	if existRole > 0 {
		msg := fmt.Sprintf("不允许删除，菜单已分配给角色数：%d", existRole)
		c.JSON(200, result.ErrMsg(msg))
		return
	}

	rows := s.sysMenuService.DeleteMenuById(menuId)
	if rows > 0 {
		msg := fmt.Sprintf("删除成功：%d", rows)
		c.JSON(200, result.OkMsg(msg))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 菜单树结构列表
//
// GET /treeSelect
func (s *sysMenu) TreeSelect(c *gin.Context) {
	query := model.SysMenu{}
	if v, ok := c.GetQuery("menuName"); ok {
		query.MenuName = v
	}
	if v, ok := c.GetQuery("status"); ok {
		query.Status = v
	}

	userId := ctx.LoginUserToUserID(c)
	if config.IsAdmin(userId) {
		userId = "*"
	}
	data := s.sysMenuService.SelectMenuTreeSelectByUserId(query, userId)
	c.JSON(200, result.OkData(data))

}

// 菜单树结构列表（指定角色）
//
// GET /roleMenuTreeSelect/:roleId
func (s *sysMenu) RoleMenuTreeSelect(c *gin.Context) {
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
	if config.IsAdmin(userId) {
		userId = "*"
	}
	menuTreeSelect := s.sysMenuService.SelectMenuTreeSelectByUserId(query, userId)
	checkedKeys := s.sysMenuService.SelectMenuListByRoleId(roleId)
	c.JSON(200, result.OkData(map[string]interface{}{
		"menus":       menuTreeSelect,
		"checkedKeys": checkedKeys,
	}))
}
