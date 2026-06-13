package controller

import (
	"fmt"
	"slices"
	"time"

	"github.com/gin-gonic/gin"

	"mask_api_gin/src/framework/constants"
	"mask_api_gin/src/framework/reqctx"
	"mask_api_gin/src/framework/resp"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/service"
)

// NewSysRole 实例化控制层
var NewSysRole = &SysRoleController{
	sysRoleService: service.NewSysRole,
	sysUserService: service.NewSysUser,
}

// SysRoleController 角色信息 控制层处理
//
// PATH /system/role
type SysRoleController struct {
	sysRoleService *service.SysRole // 角色服务
	sysUserService *service.SysUser // 用户服务
}

// List 角色列表
//
// GET /list
func (s SysRoleController) List(c *gin.Context) {
	query := reqctx.QueryMap(c)
	rows, total := s.sysRoleService.FindByPage(query)
	c.JSON(200, resp.OkData(map[string]any{"rows": rows, "total": total}))
}

// Info 角色信息详情
//
// GET /:roleId
func (s SysRoleController) Info(c *gin.Context) {
	roleId := parse.Number(c.Param("roleId"))
	if roleId <= 0 {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: roleId is empty"))
		return
	}

	data := s.sysRoleService.FindById(roleId)
	if data.RoleId == roleId {
		c.JSON(200, resp.OkData(data))
		return
	}
	c.JSON(200, resp.Err(nil))
}

// Add 角色信息新增
//
// POST /
func (s SysRoleController) Add(c *gin.Context) {
	var body model.SysRole
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_PARSER, errMsgs))
		return
	}
	if body.RoleId > 0 {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: roleId not is empty"))
		return
	}

	// 判断角色名称是否唯一
	uniqueName := s.sysRoleService.CheckUniqueByName(body.RoleName, 0)
	if !uniqueName {
		msg := fmt.Sprintf("角色新增【%s】失败，角色名称已存在", body.RoleName)
		c.JSON(200, resp.ErrMsg(msg))
		return
	}

	// 判断角色键值是否唯一
	uniqueKey := s.sysRoleService.CheckUniqueByKey(body.RoleKey, 0)
	if !uniqueKey {
		msg := fmt.Sprintf("角色新增【%s】失败，角色键值已存在", body.RoleName)
		c.JSON(200, resp.ErrMsg(msg))
		return
	}

	body.CreateBy = reqctx.LoginUserToUserName(c)
	insertId := s.sysRoleService.Insert(body)
	if insertId > 0 {
		c.JSON(200, resp.OkData(insertId))
		return
	}
	c.JSON(200, resp.Err(nil))
}

// Edit 角色信息修改
//
// PUT /
func (s SysRoleController) Edit(c *gin.Context) {
	var body model.SysRole
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_PARSER, errMsgs))
		return
	}
	if body.RoleId <= 0 {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: roleId is empty"))
		return
	}

	// 检查是否分配给当前用户
	info, err := reqctx.LoginUser(c)
	if err != nil {
		c.JSON(401, resp.CodeMsg(resp.CODE_AUTH_INVALID, err.Error()))
		return
	}
	if len(info.User.RoleIds) > 0 {
		if slices.Contains(info.User.RoleIds, body.RoleId) {
			c.JSON(200, resp.ErrMsg("不允许操作分配给自己的角色"))
			return
		}
	}

	// 检查是否系统管理员角色
	if body.RoleId == constants.SYS_ROLE_SYSTEM_ID {
		c.JSON(200, resp.ErrMsg("不允许操作系统管理员角色"))
		return
	}

	// 判断角色名称是否唯一
	uniqueName := s.sysRoleService.CheckUniqueByName(body.RoleName, body.RoleId)
	if !uniqueName {
		msg := fmt.Sprintf("角色修改【%s】失败，角色名称已存在", body.RoleName)
		c.JSON(200, resp.ErrMsg(msg))
		return
	}

	// 判断角色键值是否唯一
	uniqueKey := s.sysRoleService.CheckUniqueByKey(body.RoleKey, body.RoleId)
	if !uniqueKey {
		msg := fmt.Sprintf("角色修改【%s】失败，角色键值已存在", body.RoleName)
		c.JSON(200, resp.ErrMsg(msg))
		return
	}

	// 检查是否存在
	roleInfo := s.sysRoleService.FindById(body.RoleId)
	if roleInfo.RoleId != body.RoleId {
		c.JSON(200, resp.ErrMsg("没有权限访问角色数据！"))
		return
	}

	roleInfo.RoleName = body.RoleName
	roleInfo.RoleKey = body.RoleKey
	roleInfo.RoleSort = body.RoleSort
	roleInfo.DataScope = body.DataScope
	roleInfo.MenuCheckStrictly = body.MenuCheckStrictly
	roleInfo.DeptCheckStrictly = body.DeptCheckStrictly
	roleInfo.StatusFlag = body.StatusFlag
	roleInfo.Remark = body.Remark
	roleInfo.UpdateBy = reqctx.LoginUserToUserName(c)
	roleInfo.MenuIds = body.MenuIds
	roleInfo.DeptIds = body.DeptIds
	rows := s.sysRoleService.Update(roleInfo)
	if rows > 0 {
		c.JSON(200, resp.Ok(nil))
		return
	}
	c.JSON(200, resp.Err(nil))
}

// Remove 角色信息删除
//
// DELETE /:roleId
func (s SysRoleController) Remove(c *gin.Context) {
	roleId := c.Param("roleId")
	if roleId == "" {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: roleId is empty"))
		return
	}

	// 处理字符转id数组后去重
	uniqueIDs := parse.RemoveDuplicatesToArray(roleId, ",")
	// 转换成int64数组类型
	ids := make([]int64, 0)
	for _, v := range uniqueIDs {
		ids = append(ids, parse.Number(v))
	}

	// 检查是否系统管理员角色
	for _, id := range ids {
		if id == constants.SYS_ROLE_SYSTEM_ID {
			c.JSON(200, resp.ErrMsg("不允许操作系统管理员角色"))
			return
		}
	}

	rows, err := s.sysRoleService.DeleteByIds(ids)
	if err != nil {
		c.JSON(200, resp.ErrMsg(err.Error()))
		return
	}
	msg := fmt.Sprintf("删除成功：%d", rows)
	c.JSON(200, resp.OkMsg(msg))
}

// Status 角色状态变更
//
// PUT /status
func (s SysRoleController) Status(c *gin.Context) {
	var body struct {
		RoleID     int64  `json:"roleId" binding:"required"`               // 角色ID
		StatusFlag string `json:"statusFlag" binding:"required,oneof=0 1"` // 状态
	}
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_PARSER, errMsgs))
		return
	}

	// 检查是否系统管理员角色
	if body.RoleID == constants.SYS_ROLE_SYSTEM_ID {
		c.JSON(200, resp.ErrMsg("不允许操作系统管理员角色"))
		return
	}

	// 检查是否存在
	role := s.sysRoleService.FindById(body.RoleID)
	if role.RoleId != body.RoleID {
		c.JSON(200, resp.ErrMsg("没有权限访问角色数据！"))
		return
	}

	// 与旧值相等不变更
	if role.StatusFlag == body.StatusFlag {
		c.JSON(200, resp.ErrMsg("变更状态与旧值相等！"))
		return
	}

	// 更新状态不刷新缓存
	role.StatusFlag = body.StatusFlag
	role.UpdateBy = reqctx.LoginUserToUserName(c)
	rows := s.sysRoleService.Update(role)
	if rows > 0 {
		c.JSON(200, resp.Ok(nil))
		return
	}
	c.JSON(200, resp.Err(nil))
}

// DataScope 角色数据权限修改
//
// PUT /data-scope
func (s SysRoleController) DataScope(c *gin.Context) {
	var body struct {
		RoleId            int64   `json:"roleId" binding:"required"`                      // 角色ID
		DeptIds           []int64 `json:"deptIds"`                                        // 部门组（数据权限）
		DataScope         string  `json:"dataScope" binding:"required,oneof=1 2 3 4 5"`   // 数据范围（1：全部数据权限 2：自定数据权限 3：本部门数据权限 4：本部门及以下数据权限 5：仅本人数据权限）
		DeptCheckStrictly string  `json:"deptCheckStrictly" binding:"required,oneof=0 1"` // 部门树选择项是否关联显示（0：父子不互相关联显示 1：父子互相关联显示）
	}
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_PARSER, errMsgs))
		return
	}

	// 检查是否系统管理员角色
	if body.RoleId == constants.SYS_ROLE_SYSTEM_ID {
		c.JSON(200, resp.ErrMsg("不允许操作系统管理员角色"))
		return
	}

	// 检查是否存在
	roleInfo := s.sysRoleService.FindById(body.RoleId)
	if roleInfo.RoleId != body.RoleId {
		c.JSON(200, resp.ErrMsg("没有权限访问角色数据！"))
		return
	}

	// 更新数据权限
	roleInfo.DeptIds = body.DeptIds
	roleInfo.DataScope = body.DataScope
	roleInfo.DeptCheckStrictly = body.DeptCheckStrictly
	roleInfo.UpdateBy = reqctx.LoginUserToUserName(c)
	rows := s.sysRoleService.UpdateAndDataScope(roleInfo)
	if rows > 0 {
		c.JSON(200, resp.Ok(nil))
		return
	}
	c.JSON(200, resp.Err(nil))
}

// UserAuthList 角色分配用户列表
//
// GET /user/list
func (s SysRoleController) UserAuthList(c *gin.Context) {
	roleId := parse.Number(c.Query("roleId"))
	if roleId <= 0 {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: roleId is empty"))
		return
	}

	// 检查是否存在
	role := s.sysRoleService.FindById(roleId)
	if role.RoleId != roleId {
		c.JSON(200, resp.ErrMsg("没有权限访问角色数据！"))
		return
	}

	query := reqctx.QueryMap(c)
	dataScopeSQL := reqctx.LoginUserToDataScopeSQL(c, "sys_user", "sys_user")
	rows, total := s.sysUserService.FindAuthUsersPage(query, dataScopeSQL)
	c.JSON(200, resp.OkData(map[string]any{"rows": rows, "total": total}))
}

// UserAuthChecked 角色分配选择授权
//
// PUT /user/auth
func (s SysRoleController) UserAuthChecked(c *gin.Context) {
	var body struct {
		RoleId  int64   `json:"roleId" binding:"required"`  // 角色ID
		UserIds []int64 `json:"userIds" binding:"required"` // 用户ID组
		Auth    bool    `json:"auth" binding:"omitempty"`   // 选择操作 添加true 取消false
	}
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_PARSER, errMsgs))
		return
	}
	if len(body.UserIds) <= 0 {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: userIds is empty"))
		return
	}

	// 检查是否存在
	role := s.sysRoleService.FindById(body.RoleId)
	if role.RoleId != body.RoleId {
		c.JSON(200, resp.ErrMsg("没有权限访问角色数据！"))
		return
	}

	var rows int64
	if body.Auth {
		rows = s.sysRoleService.InsertAuthUsers(body.RoleId, body.UserIds)
	} else {
		rows = s.sysRoleService.DeleteAuthUsers(body.RoleId, body.UserIds)
	}
	if rows > 0 {
		c.JSON(200, resp.Ok(nil))
		return
	}
	c.JSON(200, resp.Err(nil))
}

// Export 导出角色信息
//
// GET /export
func (s SysRoleController) Export(c *gin.Context) {
	// 查询结果，根据查询条件结果，单页最大值限制
	query := reqctx.QueryMap(c)
	rows, total := s.sysRoleService.FindByPage(query)
	if total == 0 {
		c.JSON(200, resp.ErrMsg("export data record as empty"))
		return
	}

	// 导出文件名称
	fileName := fmt.Sprintf("role_export_%d_%d.csv", len(rows), time.Now().UnixMilli())
	// 导出数据表格
	saveFilePath, err := s.sysRoleService.ExportData(rows, fileName)
	if err != nil {
		c.JSON(200, resp.ErrMsg(err.Error()))
		return
	}

	c.FileAttachment(saveFilePath, fileName)
}
