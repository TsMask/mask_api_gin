package service

import (
	"mask_api_gin/src/framework/vo"
	"mask_api_gin/src/modules/system/model"
)

// ISysDept 部门管理 服务层接口
type ISysDept interface {
	// SelectDeptList 查询部门管理数据
	SelectDeptList(sysDept model.SysDept, dataScopeSQL string) []model.SysDept

	// SelectDeptListByRoleId 根据角色ID查询部门树信息
	SelectDeptListByRoleId(roleId string) []string

	// SelectDeptById 根据部门ID查询信息
	SelectDeptById(deptId string) model.SysDept

	// HasChildByDeptId 是否存在子节点
	HasChildByDeptId(deptId string) int64

	// CheckDeptExistUser 查询部门是否存在用户
	CheckDeptExistUser(deptId string) int64

	// CheckUniqueDeptName 校验同级部门名称是否唯一
	CheckUniqueDeptName(deptName, parentId, deptId string) bool

	// InsertDept 新增部门信息
	InsertDept(sysDept model.SysDept) string

	// UpdateDept 修改部门信息
	UpdateDept(sysDept model.SysDept) int64

	// DeleteDeptById 删除部门管理信息
	DeleteDeptById(deptId string) int64

	// SelectDeptTreeSelect 查询部门树结构信息
	SelectDeptTreeSelect(sysDept model.SysDept, dataScopeSQL string) []vo.TreeSelect
}
