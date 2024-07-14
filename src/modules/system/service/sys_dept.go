package service

import (
	"mask_api_gin/src/framework/vo"
	"mask_api_gin/src/modules/system/model"
)

// ISysDeptService 部门管理 服务层接口
type ISysDeptService interface {
	// Find 查询数据
	Find(sysDept model.SysDept, dataScopeSQL string) []model.SysDept

	// FindById 根据ID查询信息
	FindById(deptId string) model.SysDept

	// Insert 新增信息
	Insert(sysDept model.SysDept) string

	// Update 修改信息
	Update(sysDept model.SysDept) int64

	// DeleteById 删除信息
	DeleteById(deptId string) int64

	// FindDeptIdsByRoleId 根据角色ID查询包含的部门ID
	FindDeptIdsByRoleId(roleId string) []string

	// ExistChildrenByDeptId 部门下存在子节点数量
	ExistChildrenByDeptId(deptId string) int64

	// ExistUserByDeptId 部门下存在用户数量
	ExistUserByDeptId(deptId string) int64

	// CheckUniqueParentIdByDeptName 检查同级下部门名称唯一
	CheckUniqueParentIdByDeptName(parentId, deptName,  deptId string) bool

	// BuildTreeSelect 查询部门树状结构
	BuildTreeSelect(sysDept model.SysDept, dataScopeSQL string) []vo.TreeSelect
}
