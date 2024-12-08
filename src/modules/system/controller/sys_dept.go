package controller

import (
	"mask_api_gin/src/framework/constants"
	"mask_api_gin/src/framework/response"
	"mask_api_gin/src/framework/utils/ctx"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/service"

	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
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
		DeptId   string `form:"deptId"`   // 部门ID
		ParentId string `form:"parentId"` // 父部门ID
		DeptName string `form:"deptName"` // 部门名称
		Status   string `form:"status"`   // 部门状态（0正常 1停用）
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", response.FormatBindError(err))
		c.JSON(400, response.CodeMsg(40010, errMsgs))
		return
	}

	SysDeptController := model.SysDept{
		DeptId:     query.DeptId,
		ParentId:   query.ParentId,
		DeptName:   query.DeptName,
		StatusFlag: query.Status,
	}
	dataScopeSQL := ctx.LoginUserToDataScopeSQL(c, "d", "")
	data := s.sysDeptService.Find(SysDeptController, dataScopeSQL)
	c.JSON(200, response.OkData(data))
}

// Info 部门信息
//
// GET /:deptId
func (s SysDeptController) Info(c *gin.Context) {
	deptId := c.Param("deptId")
	if deptId == "" {
		c.JSON(400, response.CodeMsg(40010, "bind err: deptId is empty"))
		return
	}

	data := s.sysDeptService.FindById(deptId)
	if data.DeptId == deptId {
		c.JSON(200, response.OkData(data))
		return
	}
	c.JSON(200, response.Err(nil))
}

// Add 部门新增
//
// POST /
func (s SysDeptController) Add(c *gin.Context) {
	var body model.SysDept
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", response.FormatBindError(err))
		c.JSON(400, response.CodeMsg(40010, errMsgs))
		return
	}
	if body.DeptId != "" {
		c.JSON(400, response.CodeMsg(40010, "bind err: deptId not is empty"))
		return
	}

	// 父级ID不为0是要检查
	if body.ParentId != "0" {
		deptParent := s.sysDeptService.FindById(body.ParentId)
		if deptParent.DeptId != body.ParentId {
			c.JSON(200, response.ErrMsg("没有权限访问部门数据！"))
			return
		}
		if deptParent.StatusFlag == constants.STATUS_NO {
			msg := fmt.Sprintf("上级部门【%s】停用，不允许新增", deptParent.DeptName)
			c.JSON(200, response.ErrMsg(msg))
			return
		}
		if deptParent.DelFlag == constants.STATUS_YES {
			msg := fmt.Sprintf("上级部门【%s】已删除，不允许新增", deptParent.DeptName)
			c.JSON(200, response.ErrMsg(msg))
			return
		}
		body.Ancestors = fmt.Sprintf("%s,%s", deptParent.Ancestors, body.ParentId)
	} else {
		body.Ancestors = "0"
	}

	// 检查同级下名称唯一
	uniqueDeptName := s.sysDeptService.CheckUniqueParentIdByDeptName(body.ParentId, body.DeptName, "")
	if !uniqueDeptName {
		msg := fmt.Sprintf("部门新增【%s】失败，部门名称已存在", body.DeptName)
		c.JSON(200, response.ErrMsg(msg))
		return
	}

	body.CreateBy = ctx.LoginUserToUserName(c)
	insertId := s.sysDeptService.Insert(body)
	if insertId != "" {
		c.JSON(200, response.OkData(insertId))
		return
	}
	c.JSON(200, response.Err(nil))
}

// Edit 部门修改
//
// PUT /
func (s SysDeptController) Edit(c *gin.Context) {
	var body model.SysDept
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", response.FormatBindError(err))
		c.JSON(400, response.CodeMsg(40010, errMsgs))
		return
	}
	if body.DeptId == "" {
		c.JSON(400, response.CodeMsg(40010, "bind err: deptId is empty"))
		return
	}

	// 上级部门不能选自己
	if body.DeptId == body.ParentId {
		msg := fmt.Sprintf("部门修改【%s】失败，上级部门不能是自己", body.DeptName)
		c.JSON(200, response.ErrMsg(msg))
		return
	}

	// 检查数据是否存在
	deptInfo := s.sysDeptService.FindById(body.DeptId)
	if deptInfo.DeptId != body.DeptId {
		c.JSON(200, response.ErrMsg("没有权限访问部门数据！"))
		return
	}

	// 父级ID不为0是要检查
	if body.ParentId != "0" {
		deptParent := s.sysDeptService.FindById(body.ParentId)
		if deptParent.DeptId != body.ParentId {
			c.JSON(200, response.ErrMsg("没有权限访问部门数据！"))
			return
		}
	}

	// 检查同级下名称唯一
	uniqueDeptName := s.sysDeptService.CheckUniqueParentIdByDeptName(body.ParentId, body.DeptName, body.DeptId)
	if !uniqueDeptName {
		msg := fmt.Sprintf("部门修改【%s】失败，部门名称已存在", body.DeptName)
		c.JSON(200, response.ErrMsg(msg))
		return
	}

	// 上级停用需要检查下级是否有在使用
	if body.StatusFlag == constants.STATUS_NO {
		hasChild := s.sysDeptService.ExistChildrenByDeptId(body.DeptId)
		if hasChild > 0 {
			msg := fmt.Sprintf("该部门包含未停用的子部门数量：%d", hasChild)
			c.JSON(200, response.ErrMsg(msg))
			return
		}
	}

	deptInfo.DeptName = body.DeptName
	deptInfo.ParentId = body.ParentId
	deptInfo.DeptSort = body.DeptSort
	deptInfo.Leader = body.Leader
	deptInfo.Phone = body.Phone
	deptInfo.Email = body.Email
	deptInfo.StatusFlag = body.StatusFlag
	deptInfo.UpdateBy = ctx.LoginUserToUserName(c)
	rows := s.sysDeptService.Update(deptInfo)
	if rows > 0 {
		c.JSON(200, response.Ok(nil))
		return
	}
	c.JSON(200, response.Err(nil))
}

// Remove 部门删除
//
// DELETE /:deptId
func (s SysDeptController) Remove(c *gin.Context) {
	deptId := c.Param("deptId")
	if deptId == "" {
		c.JSON(400, response.CodeMsg(40010, "bind err: deptId is empty"))
		return
	}

	// 检查数据是否存在
	dept := s.sysDeptService.FindById(deptId)
	if dept.DeptId != deptId {
		c.JSON(200, response.ErrMsg("没有权限访问部门数据！"))
		return
	}

	// 检查是否存在子部门
	hasChild := s.sysDeptService.ExistChildrenByDeptId(deptId)
	if hasChild > 0 {
		msg := fmt.Sprintf("不允许删除，存在子部门数：%d", hasChild)
		c.JSON(200, response.ErrMsg(msg))
		return
	}

	// 检查是否分配给用户
	existUser := s.sysDeptService.ExistUserByDeptId(deptId)
	if existUser > 0 {
		msg := fmt.Sprintf("不允许删除，部门已分配给用户数：%d", existUser)
		c.JSON(200, response.ErrMsg(msg))
		return
	}

	rows := s.sysDeptService.DeleteById(deptId)
	if rows > 0 {
		msg := fmt.Sprintf("删除成功：%d", rows)
		c.JSON(200, response.OkMsg(msg))
		return
	}
	c.JSON(200, response.Err(nil))
}

// ExcludeChild 部门列表（排除节点）
//
// GET /list/exclude/:deptId
func (s SysDeptController) ExcludeChild(c *gin.Context) {
	deptId := c.Param("deptId")
	if deptId == "" {
		c.JSON(400, response.CodeMsg(40010, "bind err: deptId is empty"))
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
		if !(dept.DeptId == deptId || hasAncestor) {
			filtered = append(filtered, dept)
		}
	}
	c.JSON(200, response.OkData(filtered))
}

// Tree 部门树结构列表
//
// GET /tree
func (s SysDeptController) Tree(c *gin.Context) {
	var query struct {
		DeptId     string `form:"deptId"`     // 部门ID
		ParentId   string `form:"parentId"`   // 父部门ID
		DeptName   string `form:"deptName"`   // 部门名称
		StatusFlag string `form:"statusFlag"` // 部门状态（0正常 1停用）
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", response.FormatBindError(err))
		c.JSON(400, response.CodeMsg(40010, errMsgs))
		return
	}

	sysDept := model.SysDept{
		DeptId:     query.DeptId,
		ParentId:   query.ParentId,
		DeptName:   query.DeptName,
		StatusFlag: query.StatusFlag,
	}
	dataScopeSQL := ctx.LoginUserToDataScopeSQL(c, "d", "")
	data := s.sysDeptService.BuildTreeSelect(sysDept, dataScopeSQL)
	c.JSON(200, response.OkData(data))
}

// TreeRole 部门树结构列表（指定角色）
//
// GET /tree/role/:roleId
func (s SysDeptController) TreeRole(c *gin.Context) {
	roleId := c.Param("roleId")
	if roleId == "" {
		c.JSON(400, response.CodeMsg(40010, "bind err: roleId is empty"))
		return
	}

	dataScopeSQL := ctx.LoginUserToDataScopeSQL(c, "d", "")
	deptTreeSelect := s.sysDeptService.BuildTreeSelect(model.SysDept{}, dataScopeSQL)
	checkedKeys := s.sysDeptService.FindDeptIdsByRoleId(roleId)
	c.JSON(200, response.OkData(map[string]any{
		"depts":       deptTreeSelect,
		"checkedKeys": checkedKeys,
	}))
}
