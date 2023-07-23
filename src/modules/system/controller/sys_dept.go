package controller

import (
	"fmt"
	"mask_api_gin/src/framework/constants/common"
	"mask_api_gin/src/framework/utils/ctx"
	"mask_api_gin/src/framework/vo/result"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/service"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// 部门信息
//
// PATH /system/dept
var SysDept = &sysDept{
	sysDeptService: service.SysDeptImpl,
}

type sysDept struct {
	sysDeptService service.ISysDept
}

// 部门列表
//
// GET /list
func (s *sysDept) List(c *gin.Context) {
	var querys struct {
		// 部门ID
		DeptID string `json:"deptId"`
		// 父部门ID
		ParentID string `json:"parentId" `
		// 部门名称
		DeptName string `json:"deptName" `
		// 部门状态（0正常 1停用）
		Status string `json:"status"`
	}
	err := c.ShouldBindQuery(&querys)
	if err != nil {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	sysDept := model.SysDept{
		DeptID:   querys.DeptID,
		ParentID: querys.ParentID,
		DeptName: querys.DeptName,
		Status:   querys.Status,
	}
	dataScopeSQL := ctx.LoginUserToDataScopeSQL(c, "d", "")
	data := s.sysDeptService.SelectDeptList(sysDept, dataScopeSQL)
	c.JSON(200, result.OkData(data))
}

// 部门信息
//
// GET /:deptId
func (s *sysDept) Info(c *gin.Context) {
	deptId := c.Param("deptId")
	if deptId == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}
	data := s.sysDeptService.SelectDeptById(deptId)
	if data.DeptID == deptId {
		c.JSON(200, result.OkData(data))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 部门新增
//
// POST /
func (s *sysDept) Add(c *gin.Context) {
	var body model.SysDept
	err := c.ShouldBindBodyWith(&body, binding.JSON)
	if err != nil || body.DeptID != "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 父级ID不为0是要检查
	if body.ParentID != "0" {
		deptParent := s.sysDeptService.SelectDeptById(body.ParentID)
		if deptParent.DeptID != body.ParentID {
			c.JSON(200, result.OkMsg("没有权限访问部门数据！"))
			return
		}
		if deptParent.Status == common.STATUS_NO {
			msg := fmt.Sprintf("上级部门【%s】停用，不允许新增", deptParent.DeptName)
			c.JSON(200, result.OkMsg(msg))
			return
		}
		if deptParent.DelFlag == common.STATUS_YES {
			msg := fmt.Sprintf("上级部门【%s】已删除，不允许新增", deptParent.DeptName)
			c.JSON(200, result.OkMsg(msg))
			return
		}
		body.Ancestors = deptParent.Ancestors + "," + body.ParentID
	} else {
		body.Ancestors = "0"
	}

	// 检查同级下名称唯一
	uniqueDeptName := s.sysDeptService.CheckUniqueDeptName(body.DeptName, body.ParentID, "")
	if !uniqueDeptName {
		msg := fmt.Sprintf("部门新增【%s】失败，部门名称已存在", body.DeptName)
		c.JSON(200, result.OkMsg(msg))
		return
	}

	body.CreateBy = ctx.LoginUserToUserName(c)
	insertId := s.sysDeptService.InsertDept(body)
	if insertId != "" {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 部门修改
//
// PUT /
func (s *sysDept) Edit(c *gin.Context) {
	var body model.SysDept
	err := c.ShouldBindBodyWith(&body, binding.JSON)
	if err != nil || body.DeptID == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 上级部门不能选自己
	if body.DeptID == body.ParentID {
		msg := fmt.Sprintf("部门修改【%s】失败，上级部门不能是自己", body.DeptName)
		c.JSON(200, result.OkMsg(msg))
		return
	}

	// 检查数据是否存在
	deptInfo := s.sysDeptService.SelectDeptById(body.DeptID)
	if deptInfo.DeptID != body.DeptID {
		c.JSON(200, result.OkMsg("没有权限访问部门数据！"))
		return
	}
	// 父级ID不为0是要检查
	if body.ParentID != "0" {
		deptParent := s.sysDeptService.SelectDeptById(body.ParentID)
		if deptParent.DeptID != body.ParentID {
			c.JSON(200, result.OkMsg("没有权限访问部门数据！"))
			return
		}
	}

	// 检查同级下名称唯一
	uniqueDeptName := s.sysDeptService.CheckUniqueDeptName(body.DeptName, body.ParentID, body.DeptID)
	if !uniqueDeptName {
		msg := fmt.Sprintf("部门修改【%s】失败，部门名称已存在", body.DeptName)
		c.JSON(200, result.OkMsg(msg))
		return
	}

	// 上级停用需要检查下级是否有在使用
	if body.Status == common.STATUS_NO {
		hasChild := s.sysDeptService.HasChildByDeptId(body.DeptID)
		if hasChild > 0 {
			msg := fmt.Sprintf("该部门包含未停用的子部门数量：%d", hasChild)
			c.JSON(200, result.OkMsg(msg))
			return
		}
	}

	body.UpdateBy = ctx.LoginUserToUserName(c)
	rows := s.sysDeptService.UpdateDept(body)
	if rows > 0 {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 部门删除
//
// DELETE /:deptId
func (s *sysDept) Remove(c *gin.Context) {
	deptId := c.Param("deptId")
	if deptId == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 检查数据是否存在
	dept := s.sysDeptService.SelectDeptById(deptId)
	if dept.DeptID != deptId {
		c.JSON(200, result.ErrMsg("没有权限访问部门数据！"))
		return
	}

	// 检查是否存在子部门
	hasChild := s.sysDeptService.HasChildByDeptId(deptId)
	if hasChild > 0 {
		msg := fmt.Sprintf("不允许删除，存在子部门数：%d", hasChild)
		c.JSON(200, result.ErrMsg(msg))
		return
	}

	// 检查是否分配给用户
	existUser := s.sysDeptService.CheckDeptExistUser(deptId)
	if existUser > 0 {
		msg := fmt.Sprintf("不允许删除，部门已分配给用户数：%d", existUser)
		c.JSON(200, result.ErrMsg(msg))
		return
	}

	rows := s.sysDeptService.DeleteDeptById(deptId)
	if rows > 0 {
		msg := fmt.Sprintf("删除成功：%d", rows)
		c.JSON(200, result.OkMsg(msg))
	}
	c.JSON(200, result.Err(nil))
}

// 部门列表（排除节点）
//
// GET /list/exclude/:deptId
func (s *sysDept) ExcludeChild(c *gin.Context) {
	deptId := c.Param("deptId")
	if deptId == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	dataScopeSQL := ctx.LoginUserToDataScopeSQL(c, "d", "")
	data := s.sysDeptService.SelectDeptList(model.SysDept{}, dataScopeSQL)

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

// 部门树结构列表
//
// GET /treeSelect
func (s *sysDept) TreeSelect(c *gin.Context) {
	var querys struct {
		// 部门ID
		DeptID string `json:"deptId"`
		// 父部门ID
		ParentID string `json:"parentId" `
		// 部门名称
		DeptName string `json:"deptName" `
		// 部门状态（0正常 1停用）
		Status string `json:"status"`
	}
	err := c.ShouldBindQuery(&querys)
	if err != nil {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	sysDept := model.SysDept{
		DeptID:   querys.DeptID,
		ParentID: querys.ParentID,
		DeptName: querys.DeptName,
		Status:   querys.Status,
	}
	dataScopeSQL := ctx.LoginUserToDataScopeSQL(c, "d", "")
	data := s.sysDeptService.SelectDeptTreeSelect(sysDept, dataScopeSQL)
	c.JSON(200, result.OkData(data))
}

// 部门树结构列表（指定角色）
//
// GET /roleDeptTreeSelect/:roleId
func (s *sysDept) RoleDeptTreeSelect(c *gin.Context) {
	roleId := c.Param("roleId")
	if roleId == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	dataScopeSQL := ctx.LoginUserToDataScopeSQL(c, "d", "")
	deptTreeSelect := s.sysDeptService.SelectDeptList(model.SysDept{}, dataScopeSQL)
	checkedKeys := s.sysDeptService.SelectDeptListByRoleId(roleId)
	c.JSON(200, result.OkData(map[string]interface{}{
		"depts":       deptTreeSelect,
		"checkedKeys": checkedKeys,
	}))
}
