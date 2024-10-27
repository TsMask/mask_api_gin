package controller

import (
	"fmt"
	constSystem "mask_api_gin/src/framework/constants/system"
	"mask_api_gin/src/framework/utils/ctx"
	"mask_api_gin/src/framework/vo/result"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/service"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// NewSysDept 实例化控制层
var NewSysDept = &SysDeptController{
	sysDeptService: service.NewSysDept,
}

// SysDeptController 部门信息 控制层处理
//
// PATH /system/dept
type SysDeptController struct {
	sysDeptService *service.SysDept // 部门服务
}

// List 部门列表
//
// GET /list
func (s SysDeptController) List(c *gin.Context) {
	var query struct {
		DeptID   string `form:"deptId"`   // 部门ID
		ParentID string `form:"parentId"` // 父部门ID
		DeptName string `form:"deptName"` // 部门名称
		Status   string `form:"status"`   // 部门状态（0正常 1停用）
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	SysDeptController := model.SysDept{
		DeptID:   query.DeptID,
		ParentID: query.ParentID,
		DeptName: query.DeptName,
		Status:   query.Status,
	}
	dataScopeSQL := ctx.LoginUserToDataScopeSQL(c, "d", "")
	data := s.sysDeptService.Find(SysDeptController, dataScopeSQL)
	c.JSON(200, result.OkData(data))
}

// Info 部门信息
//
// GET /:deptId
func (s SysDeptController) Info(c *gin.Context) {
	deptId := c.Param("deptId")
	if deptId == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}
	data := s.sysDeptService.FindById(deptId)
	if data.DeptID == deptId {
		c.JSON(200, result.OkData(data))
		return
	}
	c.JSON(200, result.Err(nil))
}

// Add 部门新增
//
// POST /
func (s SysDeptController) Add(c *gin.Context) {
	var body model.SysDept
	err := c.ShouldBindBodyWith(&body, binding.JSON)
	if err != nil || body.DeptID != "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 父级ID不为0是要检查
	if body.ParentID != "0" {
		deptParent := s.sysDeptService.FindById(body.ParentID)
		if deptParent.DeptID != body.ParentID {
			c.JSON(200, result.ErrMsg("没有权限访问部门数据！"))
			return
		}
		if deptParent.Status == constSystem.STATUS_NO {
			msg := fmt.Sprintf("上级部门【%s】停用，不允许新增", deptParent.DeptName)
			c.JSON(200, result.ErrMsg(msg))
			return
		}
		if deptParent.DelFlag == constSystem.STATUS_YES {
			msg := fmt.Sprintf("上级部门【%s】已删除，不允许新增", deptParent.DeptName)
			c.JSON(200, result.ErrMsg(msg))
			return
		}
		body.Ancestors = deptParent.Ancestors + "," + body.ParentID
	} else {
		body.Ancestors = "0"
	}

	// 检查同级下名称唯一
	uniqueDeptName := s.sysDeptService.CheckUniqueParentIdByDeptName(body.ParentID, body.DeptName, "")
	if !uniqueDeptName {
		msg := fmt.Sprintf("部门新增【%s】失败，部门名称已存在", body.DeptName)
		c.JSON(200, result.ErrMsg(msg))
		return
	}

	body.CreateBy = ctx.LoginUserToUserName(c)
	insertId := s.sysDeptService.Insert(body)
	if insertId != "" {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// Edit 部门修改
//
// PUT /
func (s SysDeptController) Edit(c *gin.Context) {
	var body model.SysDept
	err := c.ShouldBindBodyWith(&body, binding.JSON)
	if err != nil || body.DeptID == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 上级部门不能选自己
	if body.DeptID == body.ParentID {
		msg := fmt.Sprintf("部门修改【%s】失败，上级部门不能是自己", body.DeptName)
		c.JSON(200, result.ErrMsg(msg))
		return
	}

	// 检查数据是否存在
	deptInfo := s.sysDeptService.FindById(body.DeptID)
	if deptInfo.DeptID != body.DeptID {
		c.JSON(200, result.ErrMsg("没有权限访问部门数据！"))
		return
	}

	// 父级ID不为0是要检查
	if body.ParentID != "0" {
		deptParent := s.sysDeptService.FindById(body.ParentID)
		if deptParent.DeptID != body.ParentID {
			c.JSON(200, result.ErrMsg("没有权限访问部门数据！"))
			return
		}
	}

	// 检查同级下名称唯一
	uniqueDeptName := s.sysDeptService.CheckUniqueParentIdByDeptName(body.ParentID, body.DeptName, body.DeptID)
	if !uniqueDeptName {
		msg := fmt.Sprintf("部门修改【%s】失败，部门名称已存在", body.DeptName)
		c.JSON(200, result.ErrMsg(msg))
		return
	}

	// 上级停用需要检查下级是否有在使用
	if body.Status == constSystem.STATUS_NO {
		hasChild := s.sysDeptService.ExistChildrenByDeptId(body.DeptID)
		if hasChild > 0 {
			msg := fmt.Sprintf("该部门包含未停用的子部门数量：%d", hasChild)
			c.JSON(200, result.ErrMsg(msg))
			return
		}
	}

	body.UpdateBy = ctx.LoginUserToUserName(c)
	rows := s.sysDeptService.Update(body)
	if rows > 0 {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// Remove 部门删除
//
// DELETE /:deptId
func (s SysDeptController) Remove(c *gin.Context) {
	deptId := c.Param("deptId")
	if deptId == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 检查数据是否存在
	dept := s.sysDeptService.FindById(deptId)
	if dept.DeptID != deptId {
		c.JSON(200, result.ErrMsg("没有权限访问部门数据！"))
		return
	}

	// 检查是否存在子部门
	hasChild := s.sysDeptService.ExistUserByDeptId(deptId)
	if hasChild > 0 {
		msg := fmt.Sprintf("不允许删除，存在子部门数：%d", hasChild)
		c.JSON(200, result.ErrMsg(msg))
		return
	}

	// 检查是否分配给用户
	existUser := s.sysDeptService.ExistUserByDeptId(deptId)
	if existUser > 0 {
		msg := fmt.Sprintf("不允许删除，部门已分配给用户数：%d", existUser)
		c.JSON(200, result.ErrMsg(msg))
		return
	}

	rows := s.sysDeptService.DeleteById(deptId)
	if rows > 0 {
		msg := fmt.Sprintf("删除成功：%d", rows)
		c.JSON(200, result.OkMsg(msg))
		return
	}
	c.JSON(200, result.Err(nil))
}

// ExcludeChild 部门列表（排除节点）
//
// GET /list/exclude/:deptId
func (s SysDeptController) ExcludeChild(c *gin.Context) {
	deptId := c.Param("deptId")
	if deptId == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	dataScopeSQL := ctx.LoginUserToDataScopeSQL(c, "d", "")
	data := s.sysDeptService.Find(model.SysDept{}, dataScopeSQL)

	// 过滤排除节点
	filtered := make([]model.SysDept, 0)
	for _, dept := range data {
		hasAncestor := false
		ancestorList := strings.Split(dept.Ancestors, ",")
		for _, ancestor := range ancestorList {
			if ancestor == deptId {
				hasAncestor = true
				break
			}
		}
		if !(dept.DeptID == deptId || hasAncestor) {
			filtered = append(filtered, dept)
		}
	}
	c.JSON(200, result.OkData(filtered))
}

// TreeSelect 部门树结构列表
//
// GET /treeSelect
func (s SysDeptController) TreeSelect(c *gin.Context) {
	var query struct {
		DeptID   string `form:"deptId"`   // 部门ID
		ParentID string `form:"parentId"` // 父部门ID
		DeptName string `form:"deptName"` // 部门名称
		Status   string `form:"status"`   // 部门状态（0正常 1停用）
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	sysDept := model.SysDept{
		DeptID:   query.DeptID,
		ParentID: query.ParentID,
		DeptName: query.DeptName,
		Status:   query.Status,
	}
	dataScopeSQL := ctx.LoginUserToDataScopeSQL(c, "d", "")
	data := s.sysDeptService.BuildTreeSelect(sysDept, dataScopeSQL)
	c.JSON(200, result.OkData(data))
}

// RoleDeptTreeSelect 部门树结构列表（指定角色）
//
// GET /roleDeptTreeSelect/:roleId
func (s SysDeptController) RoleDeptTreeSelect(c *gin.Context) {
	roleId := c.Param("roleId")
	if roleId == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	dataScopeSQL := ctx.LoginUserToDataScopeSQL(c, "d", "")
	deptTreeSelect := s.sysDeptService.BuildTreeSelect(model.SysDept{}, dataScopeSQL)
	checkedKeys := s.sysDeptService.FindDeptIdsByRoleId(roleId)
	c.JSON(200, result.OkData(map[string]any{
		"depts":       deptTreeSelect,
		"checkedKeys": checkedKeys,
	}))
}
