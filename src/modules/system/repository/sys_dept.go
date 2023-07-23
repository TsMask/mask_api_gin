package repository

import "mask_api_gin/src/modules/system/model"

// ISysDept 部门表 数据层接口
type ISysDept interface {
	// SelectDeptList 查询部门管理数据
	SelectDeptList(sysDept model.SysDept, dataScopeSQL string) []model.SysDept

	// SelectDeptListByRoleId 根据角色ID查询部门树信息
	SelectDeptListByRoleId(roleId string, deptCheckStrictly bool) []string

	// SelectDeptById 根据部门ID查询信息
	SelectDeptById(deptId string) model.SysDept

	// SelectChildrenDeptById 根据ID查询所有子部门
	SelectChildrenDeptById(deptId string) []model.SysDept

	// HasChildByDeptId 是否存在子节点
	HasChildByDeptId(deptId string) int64

	// CheckDeptExistUser 查询部门是否存在用户
	CheckDeptExistUser(deptId string) int64

	// CheckUniqueDept 校验部门是否唯一
	CheckUniqueDept(sysDept model.SysDept) string

	// InsertDept 新增部门信息
	InsertDept(sysDept model.SysDept) string

	// UpdateDept 修改部门信息
	UpdateDept(sysDept model.SysDept) int64

	// UpdateDeptStatusNormal 修改所在部门正常状态
	UpdateDeptStatusNormal(deptIds []string) int64

	// UpdateDeptChildren 修改子元素关系
	UpdateDeptChildren(sysDepts []model.SysDept) int64

	// DeleteDeptById 删除部门管理信息
	DeleteDeptById(deptId string) int64
}
