package controller

import (
	"fmt"
	"mask_api_gin/src/framework/constants/admin"
	"mask_api_gin/src/framework/constants/common"
	"mask_api_gin/src/framework/constants/roledatascope"
	"mask_api_gin/src/framework/utils/ctx"
	"mask_api_gin/src/framework/utils/date"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/framework/vo/result"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/service"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/xuri/excelize/v2"
)

// 角色信息
//
// PATH /system/role
var SysRole = &sysRole{
	sysRoleService: service.SysRoleImpl,
	sysUserService: service.SysUserImpl,
}

type sysRole struct {
	// 角色服务
	sysRoleService service.ISysRole
	// 用户服务
	sysUserService service.ISysUser
}

// 角色列表
//
// GET /list
func (s *sysRole) List(c *gin.Context) {
	querys := ctx.QueryMapString(c)
	dataScopeSQL := ctx.LoginUserToDataScopeSQL(c, "d", "")
	data := s.sysRoleService.SelectRolePage(querys, dataScopeSQL)
	c.JSON(200, result.Ok(data))
}

// 角色信息详情
//
// GET /:roleId
func (s *sysRole) Info(c *gin.Context) {
	roleId := c.Param("roleId")
	if roleId == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}
	data := s.sysRoleService.SelectRoleById(roleId)
	if data.RoleID == roleId {
		c.JSON(200, result.OkData(data))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 角色信息新增
//
// POST /
func (s *sysRole) Add(c *gin.Context) {
	var body model.SysRole
	err := c.ShouldBindBodyWith(&body, binding.JSON)
	if err != nil || body.RoleID != "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 判断角色名称是否唯一
	uniqueRoleName := s.sysRoleService.CheckUniqueRoleName(body.RoleName, "")
	if !uniqueRoleName {
		msg := fmt.Sprintf("角色新增【%s】失败，角色名称已存在", body.RoleName)
		c.JSON(200, result.ErrMsg(msg))
		return
	}

	// 判断角色键值是否唯一
	uniqueRoleKey := s.sysRoleService.CheckUniqueRoleKey(body.RoleKey, "")
	if !uniqueRoleKey {
		msg := fmt.Sprintf("角色新增【%s】失败，角色键值已存在", body.RoleName)
		c.JSON(200, result.ErrMsg(msg))
		return
	}

	body.CreateBy = ctx.LoginUserToUserName(c)
	insertId := s.sysRoleService.InsertRole(body)
	if insertId != "" {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 角色信息修改
//
// PUT /
func (s *sysRole) Edit(c *gin.Context) {
	var body model.SysRole
	err := c.ShouldBindBodyWith(&body, binding.JSON)
	if err != nil || body.RoleID == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 检查是否管理员角色
	if body.RoleID == admin.ROLE_ID {
		c.JSON(200, result.ErrMsg("不允许操作管理员角色"))
		return
	}

	// 检查是否存在
	role := s.sysRoleService.SelectRoleById(body.RoleID)
	if role.RoleID != body.RoleID {
		c.JSON(200, result.ErrMsg("没有权限访问角色数据！"))
		return
	}

	// 判断角色名称是否唯一
	uniqueRoleName := s.sysRoleService.CheckUniqueRoleName(body.RoleName, body.RoleID)
	if !uniqueRoleName {
		msg := fmt.Sprintf("角色修改【%s】失败，角色名称已存在", body.RoleName)
		c.JSON(200, result.ErrMsg(msg))
		return
	}

	// 判断角色键值是否唯一
	uniqueRoleKey := s.sysRoleService.CheckUniqueRoleKey(body.RoleKey, body.RoleID)
	if !uniqueRoleKey {
		msg := fmt.Sprintf("角色修改【%s】失败，角色键值已存在", body.RoleName)
		c.JSON(200, result.ErrMsg(msg))
		return
	}

	body.UpdateBy = ctx.LoginUserToUserName(c)
	rows := s.sysRoleService.UpdateRole(body)
	if rows > 0 {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 角色信息删除
//
// DELETE /:roleIds
func (s *sysRole) Remove(c *gin.Context) {
	roleIds := c.Param("roleIds")
	if roleIds == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}
	// 处理字符转id数组后去重
	ids := strings.Split(roleIds, ",")
	uniqueIDs := parse.RemoveDuplicates(ids)
	if len(uniqueIDs) <= 0 {
		c.JSON(200, result.Err(nil))
		return
	}
	// 检查是否管理员角色
	for _, id := range uniqueIDs {
		if id == admin.ROLE_ID {
			c.JSON(200, result.ErrMsg("不允许操作管理员角色"))
			return
		}
	}
	rows, err := s.sysRoleService.DeleteRoleByIds(uniqueIDs)
	if err != nil {
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}
	msg := fmt.Sprintf("删除成功：%d", rows)
	c.JSON(200, result.OkMsg(msg))
}

// 角色状态变更
//
// PUT /changeStatus
func (s *sysRole) Status(c *gin.Context) {
	var body struct {
		// 角色ID
		RoleID string `json:"roleId" binding:"required"`
		// 状态
		Status string `json:"status" binding:"required"`
	}
	err := c.ShouldBindBodyWith(&body, binding.JSON)
	if err != nil {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 检查是否管理员角色
	if body.RoleID == admin.ROLE_ID {
		c.JSON(200, result.ErrMsg("不允许操作管理员角色"))
		return
	}

	// 检查是否存在
	role := s.sysRoleService.SelectRoleById(body.RoleID)
	if role.RoleID != body.RoleID {
		c.JSON(200, result.ErrMsg("没有权限访问角色数据！"))
		return
	}

	// 与旧值相等不变更
	if role.Status == body.Status {
		c.JSON(200, result.ErrMsg("变更状态与旧值相等！"))
		return
	}

	// 更新状态不刷新缓存
	userName := ctx.LoginUserToUserName(c)
	sysRole := model.SysRole{
		RoleID:   body.RoleID,
		Status:   body.Status,
		UpdateBy: userName,
	}
	rows := s.sysRoleService.UpdateRole(sysRole)
	if rows > 0 {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 角色数据权限修改
//
// PUT /dataScope
func (s *sysRole) DataScope(c *gin.Context) {
	var body struct {
		// 角色ID
		RoleID string `json:"roleId"`
		// 部门组（数据权限）
		DeptIds []string `json:"deptIds"`
		// 数据范围（1：全部数据权限 2：自定数据权限 3：本部门数据权限 4：本部门及以下数据权限 5：仅本人数据权限）
		DataScope string `json:"dataScope"`
		// 部门树选择项是否关联显示（0：父子不互相关联显示 1：父子互相关联显示）
		DeptCheckStrictly string `json:"deptCheckStrictly"`
	}
	err := c.ShouldBindBodyWith(&body, binding.JSON)
	if err != nil {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 检查是否管理员角色
	if body.RoleID == admin.ROLE_ID {
		c.JSON(200, result.ErrMsg("不允许操作管理员角色"))
		return
	}

	// 检查是否存在
	role := s.sysRoleService.SelectRoleById(body.RoleID)
	if role.RoleID != body.RoleID {
		c.JSON(200, result.ErrMsg("没有权限访问角色数据！"))
		return
	}

	// 更新数据权限
	userName := ctx.LoginUserToUserName(c)
	sysRole := model.SysRole{
		RoleID:            body.RoleID,
		DeptIds:           body.DeptIds,
		DataScope:         body.DataScope,
		DeptCheckStrictly: body.DeptCheckStrictly,
		UpdateBy:          userName,
	}
	rows := s.sysRoleService.AuthDataScope(sysRole)
	if rows > 0 {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 角色分配用户列表
//
// GET /authUser/allocatedList
func (s *sysRole) AuthUserAllocatedList(c *gin.Context) {
	querys := ctx.QueryMapString(c)
	roleId, ok := querys["roleId"]
	if !ok {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 检查是否存在
	role := s.sysRoleService.SelectRoleById(roleId)
	if role.RoleID != roleId {
		c.JSON(200, result.ErrMsg("没有权限访问角色数据！"))
		return
	}

	dataScopeSQL := ctx.LoginUserToDataScopeSQL(c, "d", "u")
	data := s.sysUserService.SelectAllocatedPage(querys, dataScopeSQL)
	c.JSON(200, result.Ok(data))
}

// 角色分配选择授权
//
// PUT /authUser/checked
func (s *sysRole) AuthUserChecked(c *gin.Context) {
	var body struct {
		// 角色ID
		RoleID string `json:"roleId" binding:"required"`
		// 用户ID组
		UserIDs string `json:"userIds" binding:"required"`
		// 选择操作 添加true 取消false
		Checked bool `json:"checked"`
	}
	err := c.ShouldBindBodyWith(&body, binding.JSON)
	if err != nil {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 处理字符转id数组后去重
	ids := strings.Split(body.UserIDs, ",")
	uniqueIDs := parse.RemoveDuplicates(ids)
	if len(uniqueIDs) <= 0 {
		c.JSON(200, result.Err(nil))
		return
	}

	// 检查是否存在
	role := s.sysRoleService.SelectRoleById(body.RoleID)
	if role.RoleID != body.RoleID {
		c.JSON(200, result.ErrMsg("没有权限访问角色数据！"))
		return
	}

	var rows int64
	if body.Checked {
		rows = s.sysRoleService.InsertAuthUsers(body.RoleID, uniqueIDs)
	} else {
		rows = s.sysRoleService.DeleteAuthUsers(body.RoleID, uniqueIDs)
	}
	if rows > 0 {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 导出角色信息
//
// POST /export
func (s *sysRole) Export(c *gin.Context) {
	// 查询结果，根据查询条件结果，单页最大值限制
	querys := ctx.QueryMapString(c)
	dataScopeSQL := ctx.LoginUserToDataScopeSQL(c, "d", "")
	data := s.sysRoleService.SelectRolePage(querys, dataScopeSQL)

	// 导出数据组装
	fileName := fmt.Sprintf("role_export_%d_%d.xlsx", data["total"], date.NowTimestamp())
	file := excelize.NewFile()
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// 创建一个工作表
	sheet := "Sheet1"
	index, err := file.NewSheet(sheet)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 设置工作簿的默认工作表
	file.SetActiveSheet(index)
	// 设置名为 Sheet1 工作表上 A 到 H 列的宽度为 20
	file.SetColWidth("Sheet1", "A", "H", 20)
	// 设置单元格的值
	file.SetCellValue(sheet, "A1", "角色序号")
	file.SetCellValue(sheet, "B1", "角色名称")
	file.SetCellValue(sheet, "C1", "角色权限")
	file.SetCellValue(sheet, "D1", "角色排序")
	file.SetCellValue(sheet, "E1", "数据范围")
	file.SetCellValue(sheet, "F1", "角色状态")

	for i, row := range data["rows"].([]model.SysRole) {
		idx := i + 2
		file.SetCellValue(sheet, "A"+strconv.Itoa(idx), row.RoleID)
		file.SetCellValue(sheet, "B"+strconv.Itoa(idx), row.RoleName)
		file.SetCellValue(sheet, "C"+strconv.Itoa(idx), row.RoleKey)
		file.SetCellValue(sheet, "D"+strconv.Itoa(idx), row.RoleSort)
		if v, ok := roledatascope.RoleDataScope[row.DataScope]; ok {
			file.SetCellValue(sheet, "E"+strconv.Itoa(idx), v)
		} else {
			file.SetCellValue(sheet, "E"+strconv.Itoa(idx), "")
		}
		if row.Status == common.STATUS_NO {
			file.SetCellValue(sheet, "F"+strconv.Itoa(idx), "停用")
		} else {
			file.SetCellValue(sheet, "F"+strconv.Itoa(idx), "正常")
		}
	}

	// 根据指定路径保存文件
	if err := file.SaveAs(fileName); err != nil {
		fmt.Println(err)
	}
	// 导出数据表格
	c.FileAttachment(fileName, fileName)
}
