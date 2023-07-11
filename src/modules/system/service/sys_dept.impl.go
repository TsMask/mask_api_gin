package service

import "mask_api_gin/src/modules/system/model"

// SysDeptImpl 部门表 数据层处理
var SysDeptImpl = &sysDeptImpl{
	selectSql: "",
}

type sysDeptImpl struct {
	// 查询视图对象SQL
	selectSql string
}

// SelectDeptList 查询部门管理数据
func (r *sysDeptImpl) SelectDeptList(sysDept model.SysDept, dataScopeSQL string) []model.SysDept {
	return []model.SysDept{}
}

// SelectDeptListByRoleId 根据角色ID查询部门树信息
func (r *sysDeptImpl) SelectDeptListByRoleId(roleId string, deptCheckStrictly bool) []string {
	return []string{}
}

// SelectDeptById 根据部门ID查询信息
func (r *sysDeptImpl) SelectDeptById(deptId string) model.SysDept {
	return model.SysDept{}
}

// SelectChildrenDeptById 根据ID查询所有子部门
func (r *sysDeptImpl) SelectChildrenDeptById(deptId string) []model.SysDept {
	return []model.SysDept{}
}

// SelectNormalChildrenDeptById 根据ID查询所有子部门（正常状态）
func (r *sysDeptImpl) SelectNormalChildrenDeptById(deptId string) int {
	return 0
}

// HasChildByDeptId 是否存在子节点
func (r *sysDeptImpl) HasChildByDeptId(deptId string) int {
	return 0
}

// CheckDeptExistUser 查询部门是否存在用户
func (r *sysDeptImpl) CheckDeptExistUser(deptId string) int {
	return 0
}

// CheckUniqueDeptName 校验部门名称是否唯一
func (r *sysDeptImpl) CheckUniqueDeptName(deptName string, parentId string) string {
	return ""
}

// InsertDept 新增部门信息
func (r *sysDeptImpl) InsertDept(sysDept model.SysDept) string {
	return ""
}

// UpdateDept 修改部门信息
func (r *sysDeptImpl) UpdateDept(sysDept model.SysDept) int {
	return 0
}

// UpdateDeptStatusNormal 修改所在部门正常状态
func (r *sysDeptImpl) UpdateDeptStatusNormal(deptIds []string) int {
	return 0
}

// UpdateDeptChildren 修改子元素关系
func (r *sysDeptImpl) UpdateDeptChildren(sysDepts []model.SysDept) int {
	return 0
}

// DeleteDeptById 删除部门管理信息
func (r *sysDeptImpl) DeleteDeptById(deptId string) int {
	return 0
}
