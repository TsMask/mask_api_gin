package repository

import "mask_api_gin/src/modules/system/model"

// ISysDeptRepository 部门表 数据层接口
type ISysDeptRepository interface {
	// Select 查询集合
	Select(sysDept model.SysDept, dataScopeSQL string) []model.SysDept

	// SelectById 通过ID查询信息
	SelectById(deptId string) model.SysDept

	// Insert 新增信息
	Insert(sysDept model.SysDept) string

	// Update 修改信息
	Update(sysDept model.SysDept) int64

	// DeleteById 删除信息
	DeleteById(deptId string) int64

	// CheckUnique 检查信息是否唯一
	CheckUnique(sysDept model.SysDept) string

	// ExistChildrenByDeptId 存在子节点数量
	ExistChildrenByDeptId(deptId string) int64

	// ExistUserByDeptId 存在用户使用数量
	ExistUserByDeptId(deptId string) int64

	// SelectDeptIdsByRoleId 通过角色ID查询包含的部门ID
	SelectDeptIdsByRoleId(roleId string, deptCheckStrictly bool) []string

	// SelectChildrenDeptById 根据ID查询所有子部门
	SelectChildrenDeptById(deptId string) []model.SysDept

	// UpdateDeptStatusNormal 修改所在部门正常状态
	UpdateDeptStatusNormal(deptIds []string) int64

	// UpdateDeptChildren 修改子元素关系
	UpdateDeptChildren(arr []model.SysDept) int64
}
