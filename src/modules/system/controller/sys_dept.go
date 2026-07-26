package controller

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"

	"mask_api_gin/src/internal/shared/constants"
	"mask_api_gin/src/internal/shared/parse"
	"mask_api_gin/src/internal/web/reqctx"
	"mask_api_gin/src/internal/web/resp"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/service"
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
		ID         int64  `form:"id"`                                       // 部门ID
		ParentId   int64  `form:"parentId"`                                 // 父部门ID
		DeptName   string `form:"deptName"`                                 // 部门名称
		DeptStatus string `form:"deptStatus" binding:"omitempty,oneof=0 1"` // 部门状态（0正常 1停用）
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_PARSER, errMsgs))
		return
	}

	sysDept := model.SysDept{
		ID:         query.ID,
		ParentId:   query.ParentId,
		DeptName:   query.DeptName,
		DeptStatus: query.DeptStatus,
	}
	dataScopeSQL := service.DataScopeSQL(c, "sys_dept", "")
	data := s.sysDeptService.Find(sysDept, dataScopeSQL)
	c.JSON(200, resp.OkData(data))
}

// Info 部门信息
//
// GET /:id
func (s SysDeptController) Info(c *gin.Context) {
	id := parse.Number(c.Param("id"))
	if id <= 0 {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: id is empty"))
		return
	}

	data := s.sysDeptService.FindById(id)
	if data.ID == id {
		c.JSON(200, resp.OkData(data))
		return
	}
	c.JSON(200, resp.Err(nil))
}

// Add 部门新增
//
// POST /
func (s SysDeptController) Add(c *gin.Context) {
	var body model.SysDept
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_PARSER, errMsgs))
		return
	}
	if body.ID > 0 {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: deptId not is empty"))
		return
	}

	// 父级ID不为0是要检查
	if body.ParentId > 0 {
		deptParent := s.sysDeptService.FindById(body.ParentId)
		if deptParent.ID != body.ParentId {
			c.JSON(200, resp.ErrMsg("没有权限访问部门数据！"))
			return
		}
		if deptParent.DeptStatus == constants.STATUS_NO {
			msg := fmt.Sprintf("上级部门【%s】停用，不允许新增", deptParent.DeptName)
			c.JSON(200, resp.ErrMsg(msg))
			return
		}
		if deptParent.DeptDel == constants.STATUS_YES {
			msg := fmt.Sprintf("上级部门【%s】已删除，不允许新增", deptParent.DeptName)
			c.JSON(200, resp.ErrMsg(msg))
			return
		}
		body.Ancestors = fmt.Sprintf("%s,%d", deptParent.Ancestors, body.ParentId)
	} else {
		body.Ancestors = "0"
	}

	// 检查同级下名称唯一
	uniqueName := s.sysDeptService.CheckUniqueParentIdByDeptName(body.ParentId, body.DeptName, 0)
	if !uniqueName {
		msg := fmt.Sprintf("部门新增【%s】失败，部门名称已存在", body.DeptName)
		c.JSON(200, resp.ErrMsg(msg))
		return
	}

	body.CreateBy = reqctx.LoginUserToUserName(c)
	insertId := s.sysDeptService.Insert(body)
	if insertId > 0 {
		c.JSON(200, resp.OkData(insertId))
		return
	}
	c.JSON(200, resp.Err(nil))
}

// Edit 部门修改
//
// PUT /
func (s SysDeptController) Edit(c *gin.Context) {
	var body model.SysDept
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_PARSER, errMsgs))
		return
	}
	if body.ID <= 0 {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: deptId is empty"))
		return
	}

	// 上级部门不能选自己
	if body.ID == body.ParentId {
		msg := fmt.Sprintf("部门修改【%s】失败，上级部门不能是自己", body.DeptName)
		c.JSON(200, resp.ErrMsg(msg))
		return
	}

	// 检查数据是否存在
	deptInfo := s.sysDeptService.FindById(body.ID)
	if deptInfo.ID != body.ID {
		c.JSON(200, resp.ErrMsg("没有权限访问部门数据！"))
		return
	}

	// 父级ID不为0是要检查
	if body.ParentId > 0 {
		deptParent := s.sysDeptService.FindById(body.ParentId)
		if deptParent.ID != body.ParentId {
			c.JSON(200, resp.ErrMsg("没有权限访问部门数据！"))
			return
		}
	}

	// 检查同级下名称唯一
	uniqueDeptName := s.sysDeptService.CheckUniqueParentIdByDeptName(body.ParentId, body.DeptName, body.ID)
	if !uniqueDeptName {
		msg := fmt.Sprintf("部门修改【%s】失败，部门名称已存在", body.DeptName)
		c.JSON(200, resp.ErrMsg(msg))
		return
	}

	// 上级停用需要检查下级是否有在使用
	if body.DeptStatus == constants.STATUS_NO {
		hasChild := s.sysDeptService.ExistChildrenByDeptId(body.ID)
		if hasChild > 0 {
			msg := fmt.Sprintf("该部门包含未停用的子部门数量：%d", hasChild)
			c.JSON(200, resp.ErrMsg(msg))
			return
		}
	}

	deptInfo.DeptName = body.DeptName
	deptInfo.ParentId = body.ParentId
	deptInfo.DeptSort = body.DeptSort
	deptInfo.DeptLeader = body.DeptLeader
	deptInfo.DeptPhone = body.DeptPhone
	deptInfo.DeptEmail = body.DeptEmail
	deptInfo.DeptStatus = body.DeptStatus
	deptInfo.UpdateBy = reqctx.LoginUserToUserName(c)
	rows := s.sysDeptService.Update(deptInfo)
	if rows > 0 {
		c.JSON(200, resp.Ok(nil))
		return
	}
	c.JSON(200, resp.Err(nil))
}

// Remove 部门删除
//
// DELETE /:id
func (s SysDeptController) Remove(c *gin.Context) {
	id := parse.Number(c.Param("id"))
	if id <= 0 {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: id is empty"))
		return
	}

	// 检查数据是否存在
	dept := s.sysDeptService.FindById(id)
	if dept.ID != id {
		c.JSON(200, resp.ErrMsg("没有权限访问部门数据！"))
		return
	}

	// 检查是否存在子部门
	hasChild := s.sysDeptService.ExistChildrenByDeptId(id)
	if hasChild > 0 {
		msg := fmt.Sprintf("不允许删除，存在子部门数：%d", hasChild)
		c.JSON(200, resp.ErrMsg(msg))
		return
	}

	// 检查是否分配给用户
	existUser := s.sysDeptService.ExistUserByDeptId(id)
	if existUser > 0 {
		msg := fmt.Sprintf("不允许删除，部门已分配给用户数：%d", existUser)
		c.JSON(200, resp.ErrMsg(msg))
		return
	}

	rows := s.sysDeptService.DeleteById(id)
	if rows > 0 {
		msg := fmt.Sprintf("删除成功：%d", rows)
		c.JSON(200, resp.OkMsg(msg))
		return
	}
	c.JSON(200, resp.Err(nil))
}

// ExcludeChild 部门列表（排除节点）
//
// GET /list/exclude/:dept_id
func (s SysDeptController) ExcludeChild(c *gin.Context) {
	deptIdStr := c.Param("dept_id")
	if deptIdStr == "" {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: deptId is empty"))
		return
	}
	deptId := parse.Number(deptIdStr)
	dataScopeSQL := service.DataScopeSQL(c, "sys_dept", "")
	data := s.sysDeptService.Find(model.SysDept{}, dataScopeSQL)

	// 过滤排除节点
	filtered := make([]model.SysDept, 0)
	for _, dept := range data {
		if dept.ID == deptId {
			continue
		}
		// 如果当前部门的ancestors不包含要排除的deptId，则添加到filtered中
		if !strings.Contains(dept.Ancestors, deptIdStr) {
			filtered = append(filtered, dept)
		}
	}
	c.JSON(200, resp.OkData(filtered))
}

// Tree 部门树结构列表
//
// GET /tree
func (s SysDeptController) Tree(c *gin.Context) {
	var query struct {
		DeptId     int64  `form:"deptId"`                                   // 部门ID
		ParentId   int64  `form:"parentId"`                                 // 父部门ID
		DeptName   string `form:"deptName"`                                 // 部门名称
		DeptStatus string `form:"deptStatus" binding:"omitempty,oneof=0 1"` // 部门状态（0正常 1停用）
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_PARSER, errMsgs))
		return
	}

	sysDept := model.SysDept{
		ID:         query.DeptId,
		ParentId:   query.ParentId,
		DeptName:   query.DeptName,
		DeptStatus: query.DeptStatus,
	}
	dataScopeSQL := service.DataScopeSQL(c, "sys_dept", "")
	data := s.sysDeptService.BuildTreeSelect(sysDept, dataScopeSQL)
	c.JSON(200, resp.OkData(data))
}

// TreeRole 部门树结构列表（指定角色）
//
// GET /tree/role/:role_id
func (s SysDeptController) TreeRole(c *gin.Context) {
	roleId := parse.Number(c.Param("role_id"))
	if roleId <= 0 {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "bind err: roleId is empty"))
		return
	}

	dataScopeSQL := service.DataScopeSQL(c, "sys_dept", "")
	deptTreeSelect := s.sysDeptService.BuildTreeSelect(model.SysDept{}, dataScopeSQL)
	checkedKeys := s.sysDeptService.FindDeptIdsByRoleId(roleId)
	c.JSON(200, resp.OkData(map[string]any{
		"depts":        deptTreeSelect,
		"checked_keys": checkedKeys,
	}))
}
