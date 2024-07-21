package controller

import (
	"fmt"
	constRoleDataScope "mask_api_gin/src/framework/constants/role_data_scope"
	constSystem "mask_api_gin/src/framework/constants/system"
	"mask_api_gin/src/framework/utils/ctx"
	"mask_api_gin/src/framework/utils/file"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/framework/vo/result"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/service"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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
	sysRoleService service.ISysRoleService // 角色服务
	sysUserService service.ISysUserService // 用户服务
}

// List 角色列表
//
// GET /list
func (s *SysRoleController) List(c *gin.Context) {
	query := ctx.QueryMap(c)
	dataScopeSQL := ctx.LoginUserToDataScopeSQL(c, "d", "")
	data := s.sysRoleService.FindByPage(query, dataScopeSQL)
	c.JSON(200, result.Ok(data))
}

// Info 角色信息详情
//
// GET /:roleId
func (s *SysRoleController) Info(c *gin.Context) {
	roleId := c.Param("roleId")
	if roleId == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}
	data := s.sysRoleService.FindById(roleId)
	if data.RoleID == roleId {
		c.JSON(200, result.OkData(data))
		return
	}
	c.JSON(200, result.Err(nil))
}

// Add 角色信息新增
//
// POST /
func (s *SysRoleController) Add(c *gin.Context) {
	var body model.SysRole
	err := c.ShouldBindBodyWith(&body, binding.JSON)
	if err != nil || body.RoleID != "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 判断角色名称是否唯一
	uniqueRoleName := s.sysRoleService.CheckUniqueByName(body.RoleName, "")
	if !uniqueRoleName {
		msg := fmt.Sprintf("角色新增【%s】失败，角色名称已存在", body.RoleName)
		c.JSON(200, result.ErrMsg(msg))
		return
	}

	// 判断角色键值是否唯一
	uniqueRoleKey := s.sysRoleService.CheckUniqueByKey(body.RoleKey, "")
	if !uniqueRoleKey {
		msg := fmt.Sprintf("角色新增【%s】失败，角色键值已存在", body.RoleName)
		c.JSON(200, result.ErrMsg(msg))
		return
	}

	body.CreateBy = ctx.LoginUserToUserName(c)
	insertId := s.sysRoleService.Insert(body)
	if insertId != "" {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// Edit 角色信息修改
//
// PUT /
func (s *SysRoleController) Edit(c *gin.Context) {
	var body model.SysRole
	err := c.ShouldBindBodyWith(&body, binding.JSON)
	if err != nil || body.RoleID == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 检查是否系统管理员角色
	if body.RoleID == constSystem.RoleId {
		c.JSON(200, result.ErrMsg("不允许操作系统管理员角色"))
		return
	}

	// 检查是否存在
	role := s.sysRoleService.FindById(body.RoleID)
	if role.RoleID != body.RoleID {
		c.JSON(200, result.ErrMsg("没有权限访问角色数据！"))
		return
	}

	// 判断角色名称是否唯一
	uniqueRoleName := s.sysRoleService.CheckUniqueByName(body.RoleName, body.RoleID)
	if !uniqueRoleName {
		msg := fmt.Sprintf("角色修改【%s】失败，角色名称已存在", body.RoleName)
		c.JSON(200, result.ErrMsg(msg))
		return
	}

	// 判断角色键值是否唯一
	uniqueRoleKey := s.sysRoleService.CheckUniqueByKey(body.RoleKey, body.RoleID)
	if !uniqueRoleKey {
		msg := fmt.Sprintf("角色修改【%s】失败，角色键值已存在", body.RoleName)
		c.JSON(200, result.ErrMsg(msg))
		return
	}

	body.UpdateBy = ctx.LoginUserToUserName(c)
	rows := s.sysRoleService.Update(body)
	if rows > 0 {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// Remove 角色信息删除
//
// DELETE /:roleIds
func (s *SysRoleController) Remove(c *gin.Context) {
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
	// 检查是否系统管理员角色
	for _, id := range uniqueIDs {
		if id == constSystem.RoleId {
			c.JSON(200, result.ErrMsg("不允许操作系统管理员角色"))
			return
		}
	}
	rows, err := s.sysRoleService.DeleteByIds(uniqueIDs)
	if err != nil {
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}
	msg := fmt.Sprintf("删除成功：%d", rows)
	c.JSON(200, result.OkMsg(msg))
}

// Status 角色状态变更
//
// PUT /changeStatus
func (s *SysRoleController) Status(c *gin.Context) {
	var body struct {
		RoleID string `json:"roleId" binding:"required"` // 角色ID
		Status string `json:"status" binding:"required"` // 状态
	}
	err := c.ShouldBindBodyWith(&body, binding.JSON)
	if err != nil {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 检查是否系统管理员角色
	if body.RoleID == constSystem.RoleId {
		c.JSON(200, result.ErrMsg("不允许操作系统管理员角色"))
		return
	}

	// 检查是否存在
	role := s.sysRoleService.FindById(body.RoleID)
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
	role.Status = body.Status
	role.UpdateBy = ctx.LoginUserToUserName(c)
	rows := s.sysRoleService.Update(role)
	if rows > 0 {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// DataScope 角色数据权限修改
//
// PUT /dataScope
func (s *SysRoleController) DataScope(c *gin.Context) {
	var body struct {
		RoleID            string   `json:"roleId"`            // 角色ID
		DeptIds           []string `json:"deptIds"`           // 部门组（数据权限）
		DataScope         string   `json:"dataScope"`         // 数据范围（1：全部数据权限 2：自定数据权限 3：本部门数据权限 4：本部门及以下数据权限 5：仅本人数据权限）
		DeptCheckStrictly string   `json:"deptCheckStrictly"` // 部门树选择项是否关联显示（0：父子不互相关联显示 1：父子互相关联显示）
	}
	err := c.ShouldBindBodyWith(&body, binding.JSON)
	if err != nil {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 检查是否系统管理员角色
	if body.RoleID == constSystem.RoleId {
		c.JSON(200, result.ErrMsg("不允许操作系统管理员角色"))
		return
	}

	// 检查是否存在
	role := s.sysRoleService.FindById(body.RoleID)
	if role.RoleID != body.RoleID {
		c.JSON(200, result.ErrMsg("没有权限访问角色数据！"))
		return
	}

	// 更新数据权限
	userName := ctx.LoginUserToUserName(c)
	SysRoleController := model.SysRole{
		RoleID:            body.RoleID,
		DeptIds:           body.DeptIds,
		DataScope:         body.DataScope,
		DeptCheckStrictly: body.DeptCheckStrictly,
		UpdateBy:          userName,
	}
	rows := s.sysRoleService.UpdateAndDataScope(SysRoleController)
	if rows > 0 {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// AuthUserAllocatedList 角色分配用户列表
//
// GET /authUser/allocatedList
func (s *SysRoleController) AuthUserAllocatedList(c *gin.Context) {
	query := ctx.QueryMap(c)
	roleId, ok := query["roleId"]
	if !ok || roleId == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 检查是否存在
	role := s.sysRoleService.FindById(roleId.(string))
	if role.RoleID != roleId {
		c.JSON(200, result.ErrMsg("没有权限访问角色数据！"))
		return
	}

	dataScopeSQL := ctx.LoginUserToDataScopeSQL(c, "d", "u")
	data := s.sysUserService.FindAllocatedPage(query, dataScopeSQL)
	c.JSON(200, result.Ok(data))
}

// AuthUserChecked 角色分配选择授权
//
// PUT /authUser/checked
func (s *SysRoleController) AuthUserChecked(c *gin.Context) {
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
	role := s.sysRoleService.FindById(body.RoleID)
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

// Export 导出角色信息
//
// POST /export
func (s *SysRoleController) Export(c *gin.Context) {
	// 查询结果，根据查询条件结果，单页最大值限制
	query := ctx.BodyJSONMap(c)
	dataScopeSQL := ctx.LoginUserToDataScopeSQL(c, "d", "")
	data := s.sysRoleService.FindByPage(query, dataScopeSQL)
	if data["total"].(int64) == 0 {
		c.JSON(200, result.ErrMsg("导出数据记录为空"))
		return
	}
	rows := data["rows"].([]model.SysRole)

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
		idx := strconv.Itoa(i + 2)
		// 数据范围
		dataScope := "空"
		if v, ok := constRoleDataScope.RoleDataScope[row.DataScope]; ok {
			dataScope = v
		}
		// 角色状态
		statusValue := "停用"
		if row.Status == "1" {
			statusValue = "正常"
		}
		dataCells = append(dataCells, map[string]any{
			"A" + idx: row.RoleID,
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
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}

	c.FileAttachment(saveFilePath, fileName)
}
