package controller

import (
	"mask_api_gin/src/framework/constants"
	"mask_api_gin/src/framework/context"
	"mask_api_gin/src/framework/response"
	"mask_api_gin/src/framework/utils/file"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/service"

	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// NewSysRole 实例化控制层
var NewSysRole = &SysRoleController{
	sysRoleService: service.NewSysRole,
	sysUserService: service.NewSysUser,
}

// SysRoleController 角色信息
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
	query := context.QueryMap(c)
	dataScopeSQL := context.LoginUserToDataScopeSQL(c, "d", "")
	rows, total := s.sysRoleService.FindByPage(query, dataScopeSQL)
	c.JSON(200, response.OkData(map[string]any{"rows": rows, "total": total}))
}

// Info 角色信息详情
//
// GET /:roleId
func (s SysRoleController) Info(c *gin.Context) {
	roleId := c.Param("roleId")
	if roleId == "" {
		c.JSON(400, response.CodeMsg(40010, "bind err: roleId is empty"))
		return
	}

	data := s.sysRoleService.FindById(roleId)
	if data.RoleId == roleId {
		c.JSON(200, response.OkData(data))
		return
	}
	c.JSON(200, response.Err(nil))
}

// Add 角色信息新增
//
// POST /
func (s SysRoleController) Add(c *gin.Context) {
	var body model.SysRole
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", response.FormatBindError(err))
		c.JSON(400, response.CodeMsg(40010, errMsgs))
		return
	}
	if body.RoleId != "" {
		c.JSON(400, response.CodeMsg(40010, "bind err: roleId not is empty"))
		return
	}

	// 判断角色名称是否唯一
	uniqueRoleName := s.sysRoleService.CheckUniqueByName(body.RoleName, "")
	if !uniqueRoleName {
		msg := fmt.Sprintf("角色新增【%s】失败，角色名称已存在", body.RoleName)
		c.JSON(200, response.ErrMsg(msg))
		return
	}

	// 判断角色键值是否唯一
	uniqueRoleKey := s.sysRoleService.CheckUniqueByKey(body.RoleKey, "")
	if !uniqueRoleKey {
		msg := fmt.Sprintf("角色新增【%s】失败，角色键值已存在", body.RoleName)
		c.JSON(200, response.ErrMsg(msg))
		return
	}

	body.CreateBy = context.LoginUserToUserName(c)
	insertId := s.sysRoleService.Insert(body)
	if insertId != "" {
		c.JSON(200, response.OkData(insertId))
		return
	}
	c.JSON(200, response.Err(nil))
}

// Edit 角色信息修改
//
// PUT /
func (s SysRoleController) Edit(c *gin.Context) {
	var body model.SysRole
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", response.FormatBindError(err))
		c.JSON(400, response.CodeMsg(40010, errMsgs))
		return
	}
	if body.RoleId == "" {
		c.JSON(400, response.CodeMsg(40010, "bind err: roleId is empty"))
		return
	}

	// 检查是否系统管理员角色
	if body.RoleId == constants.SYS_ROLE_SYSTEM_ID {
		c.JSON(200, response.ErrMsg("不允许操作系统管理员角色"))
		return
	}

	// 检查是否存在
	role := s.sysRoleService.FindById(body.RoleId)
	if role.RoleId != body.RoleId {
		c.JSON(200, response.ErrMsg("没有权限访问角色数据！"))
		return
	}

	// 判断角色名称是否唯一
	uniqueRoleName := s.sysRoleService.CheckUniqueByName(body.RoleName, body.RoleId)
	if !uniqueRoleName {
		msg := fmt.Sprintf("角色修改【%s】失败，角色名称已存在", body.RoleName)
		c.JSON(200, response.ErrMsg(msg))
		return
	}

	// 判断角色键值是否唯一
	uniqueRoleKey := s.sysRoleService.CheckUniqueByKey(body.RoleKey, body.RoleId)
	if !uniqueRoleKey {
		msg := fmt.Sprintf("角色修改【%s】失败，角色键值已存在", body.RoleName)
		c.JSON(200, response.ErrMsg(msg))
		return
	}

	body.UpdateBy = context.LoginUserToUserName(c)
	rows := s.sysRoleService.Update(body)
	if rows > 0 {
		c.JSON(200, response.Ok(nil))
		return
	}
	c.JSON(200, response.Err(nil))
}

// Remove 角色信息删除
//
// DELETE /:roleId
func (s SysRoleController) Remove(c *gin.Context) {
	roleIdsStr := c.Param("roleId")
	roleIds := parse.RemoveDuplicatesToArray(roleIdsStr, ",")
	if roleIdsStr == "" || len(roleIds) <= 0 {
		c.JSON(400, response.CodeMsg(40010, "bind err: roleId is empty"))
		return
	}

	// 检查是否系统管理员角色
	for _, id := range roleIds {
		if id == constants.SYS_ROLE_SYSTEM_ID {
			c.JSON(200, response.ErrMsg("不允许操作系统管理员角色"))
			return
		}
	}

	rows, err := s.sysRoleService.DeleteByIds(roleIds)
	if err != nil {
		c.JSON(200, response.ErrMsg(err.Error()))
		return
	}
	msg := fmt.Sprintf("删除成功：%d", rows)
	c.JSON(200, response.OkMsg(msg))
}

// Status 角色状态变更
//
// PUT /status
func (s SysRoleController) Status(c *gin.Context) {
	var body struct {
		RoleID     string `json:"roleId" binding:"required"`               // 角色ID
		StatusFlag string `json:"statusFlag" binding:"required,oneof=0 1"` // 状态
	}
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", response.FormatBindError(err))
		c.JSON(400, response.CodeMsg(40010, errMsgs))
		return
	}

	// 检查是否系统管理员角色
	if body.RoleID == constants.SYS_ROLE_SYSTEM_ID {
		c.JSON(200, response.ErrMsg("不允许操作系统管理员角色"))
		return
	}

	// 检查是否存在
	role := s.sysRoleService.FindById(body.RoleID)
	if role.RoleId != body.RoleID {
		c.JSON(200, response.ErrMsg("没有权限访问角色数据！"))
		return
	}

	// 与旧值相等不变更
	if role.StatusFlag == body.StatusFlag {
		c.JSON(200, response.ErrMsg("变更状态与旧值相等！"))
		return
	}

	// 更新状态不刷新缓存
	role.StatusFlag = body.StatusFlag
	role.UpdateBy = context.LoginUserToUserName(c)
	rows := s.sysRoleService.Update(role)
	if rows > 0 {
		c.JSON(200, response.Ok(nil))
		return
	}
	c.JSON(200, response.Err(nil))
}

// DataScope 角色数据权限修改
//
// PUT /data-scope
func (s SysRoleController) DataScope(c *gin.Context) {
	var body struct {
		RoleId            string   `json:"roleId" binding:"required"`                      // 角色ID
		DeptIds           []string `json:"deptIds"`                                        // 部门组（数据权限）
		DataScope         string   `json:"dataScope" binding:"required,oneof=1 2 3 4 5"`   // 数据范围（1：全部数据权限 2：自定数据权限 3：本部门数据权限 4：本部门及以下数据权限 5：仅本人数据权限）
		DeptCheckStrictly string   `json:"deptCheckStrictly" binding:"required,oneof=0 1"` // 部门树选择项是否关联显示（0：父子不互相关联显示 1：父子互相关联显示）
	}
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", response.FormatBindError(err))
		c.JSON(400, response.CodeMsg(40010, errMsgs))
		return
	}

	// 检查是否系统管理员角色
	if body.RoleId == constants.SYS_ROLE_SYSTEM_ID {
		c.JSON(200, response.ErrMsg("不允许操作系统管理员角色"))
		return
	}

	// 检查是否存在
	roleInfo := s.sysRoleService.FindById(body.RoleId)
	if roleInfo.RoleId != body.RoleId {
		c.JSON(200, response.ErrMsg("没有权限访问角色数据！"))
		return
	}

	// 更新数据权限
	roleInfo.DeptIds = body.DeptIds
	roleInfo.DataScope = body.DataScope
	roleInfo.DeptCheckStrictly = body.DeptCheckStrictly
	roleInfo.UpdateBy = context.LoginUserToUserName(c)
	rows := s.sysRoleService.UpdateAndDataScope(roleInfo)
	if rows > 0 {
		c.JSON(200, response.Ok(nil))
		return
	}
	c.JSON(200, response.Err(nil))
}

// UserAuthList 角色分配用户列表
//
// GET /user/list
func (s SysRoleController) UserAuthList(c *gin.Context) {
	roleId := c.Query("roleId")
	if roleId == "" {
		c.JSON(400, response.CodeMsg(40010, "bind err: roleId is empty"))
		return
	}

	// 检查是否存在
	role := s.sysRoleService.FindById(roleId)
	if role.RoleId != roleId {
		c.JSON(200, response.ErrMsg("没有权限访问角色数据！"))
		return
	}

	query := context.QueryMap(c)
	dataScopeSQL := context.LoginUserToDataScopeSQL(c, "d", "u")
	rows, total := s.sysUserService.FindAuthUsersPage(query, dataScopeSQL)
	c.JSON(200, response.OkData(map[string]any{"rows": rows, "total": total}))
}

// UserAuthChecked 角色分配选择授权
//
// PUT /user/auth
func (s SysRoleController) UserAuthChecked(c *gin.Context) {
	var body struct {
		RoleId  string `json:"roleId" binding:"required"`  // 角色ID
		UserIds string `json:"userIds" binding:"required"` // 用户ID组
		Auth    bool   `json:"auth"`                       // 选择操作 添加true 取消false
	}
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", response.FormatBindError(err))
		c.JSON(400, response.CodeMsg(40010, errMsgs))
		return
	}
	userIds := parse.RemoveDuplicatesToArray(body.UserIds, ",")
	if len(userIds) <= 0 {
		c.JSON(400, response.CodeMsg(40010, "bind err: userIds is empty"))
		return
	}

	// 检查是否存在
	role := s.sysRoleService.FindById(body.RoleId)
	if role.RoleId != body.RoleId {
		c.JSON(200, response.ErrMsg("没有权限访问角色数据！"))
		return
	}

	var rows int64
	if body.Auth {
		rows = s.sysRoleService.InsertAuthUsers(body.RoleId, userIds)
	} else {
		rows = s.sysRoleService.DeleteAuthUsers(body.RoleId, userIds)
	}
	if rows > 0 {
		c.JSON(200, response.Ok(nil))
		return
	}
	c.JSON(200, response.Err(nil))
}

// Export 导出角色信息
//
// GET /export
func (s SysRoleController) Export(c *gin.Context) {
	// 查询结果，根据查询条件结果，单页最大值限制
	query := context.QueryMap(c)
	dataScopeSQL := context.LoginUserToDataScopeSQL(c, "d", "")
	rows, total := s.sysRoleService.FindByPage(query, dataScopeSQL)
	if total == 0 {
		c.JSON(200, response.CodeMsg(40016, "export data record as empty"))
		return
	}

	// 导出文件名称
	fileName := fmt.Sprintf("role_export_%d_%d.xlsx", len(rows), time.Now().UnixMilli())
	// 第一行表头标题
	headerCells := map[string]string{
		"A1": "角色序号",
		"B1": "角色名称",
		"C1": "角色权限",
		"D1": "角色排序",
		"E1": "数据范围",
		"F1": "角色状态",
	}
	// 从第二行开始的数据
	dataCells := make([]map[string]any, 0)
	for i, row := range rows {
		idx := fmt.Sprintf("%d", i+2)
		// 数据范围
		dataScope := "空"
		if v, ok := constants.ROLE_SCOPE_DATA[row.DataScope]; ok {
			dataScope = v
		}
		// 角色状态
		statusValue := "停用"
		if row.StatusFlag == "1" {
			statusValue = "正常"
		}
		dataCells = append(dataCells, map[string]any{
			"A" + idx: row.RoleId,
			"B" + idx: row.RoleName,
			"C" + idx: row.RoleKey,
			"D" + idx: row.RoleSort,
			"E" + idx: dataScope,
			"F" + idx: statusValue,
		})
	}

	// 导出数据表格
	saveFilePath, err := file.WriteSheet(headerCells, dataCells, fileName, "")
	if err != nil {
		c.JSON(200, response.ErrMsg(err.Error()))
		return
	}

	c.FileAttachment(saveFilePath, fileName)
}
