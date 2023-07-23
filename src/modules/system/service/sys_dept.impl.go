package service

import (
	"mask_api_gin/src/framework/vo"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/repository"
)

// SysDeptImpl 部门表 数据层处理
var SysDeptImpl = &sysDeptImpl{
	sysDeptRepository: repository.SysDeptImpl,
	sysRoleRepository: repository.SysRoleImpl,
}

type sysDeptImpl struct {
	// 部门服务
	sysDeptRepository repository.ISysDept
	// 角色服务
	sysRoleRepository repository.ISysRole
}

// SelectDeptList 查询部门管理数据
func (r *sysDeptImpl) SelectDeptList(sysDept model.SysDept, dataScopeSQL string) []model.SysDept {
	return r.sysDeptRepository.SelectDeptList(sysDept, dataScopeSQL)
}

// SelectDeptListByRoleId 根据角色ID查询部门树信息 TODO
func (r *sysDeptImpl) SelectDeptListByRoleId(roleId string) []string {
	roles := r.sysRoleRepository.SelectRoleByIds([]string{roleId})
	if len(roles) > 0 {
		role := roles[0]
		if role.RoleID == roleId {
			return r.sysDeptRepository.SelectDeptListByRoleId(
				role.RoleID,
				role.DeptCheckStrictly == "1",
			)
		}
	}
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
func (r *sysDeptImpl) HasChildByDeptId(deptId string) int64 {
	return r.sysDeptRepository.HasChildByDeptId(deptId)
}

// CheckDeptExistUser 查询部门是否存在用户
func (r *sysDeptImpl) CheckDeptExistUser(deptId string) int64 {
	return r.sysDeptRepository.CheckDeptExistUser(deptId)
}

// CheckUniqueDeptName 校验同级部门名称是否唯一
func (r *sysDeptImpl) CheckUniqueDeptName(deptName, parentId, deptId string) bool {
	uniqueId := r.sysDeptRepository.CheckUniqueDept(model.SysDept{
		DeptName: deptName,
		ParentID: parentId,
	})
	if uniqueId == deptId {
		return true
	}
	return uniqueId == ""
}

// InsertDept 新增部门信息
func (r *sysDeptImpl) InsertDept(sysDept model.SysDept) string {
	return r.sysDeptRepository.InsertDept(sysDept)
}

// UpdateDept 修改部门信息
func (r *sysDeptImpl) UpdateDept(sysDept model.SysDept) int64 {
	return r.sysDeptRepository.UpdateDept(sysDept)
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
func (r *sysDeptImpl) DeleteDeptById(deptId string) int64 {
	return r.sysDeptRepository.DeleteDeptById(deptId)
}

// SelectDeptTreeSelect 查询部门树结构信息
func (r *sysDeptImpl) SelectDeptTreeSelect(sysDept model.SysDept, dataScopeSQL string) []vo.TreeSelect {
	sysDepts := r.sysDeptRepository.SelectDeptList(sysDept, dataScopeSQL)
	depts := r.parseDataToTree(sysDepts)
	tree := make([]vo.TreeSelect, 0)
	for _, dept := range depts {
		tree = append(tree, vo.SysDeptTreeSelect(dept))
	}
	return tree
}

// parseDataToTree 将数据解析为树结构，构建前端所需要下拉树结构
func (r *sysDeptImpl) parseDataToTree(sysDepts []model.SysDept) []model.SysDept {
	// 节点分组
	nodesMap := make(map[string][]model.SysDept)
	// 节点id
	treeIds := []string{}
	// 树节点
	tree := []model.SysDept{}

	for _, item := range sysDepts {
		parentID := item.ParentID
		// 分组
		mapItem, ok := nodesMap[parentID]
		if !ok {
			mapItem = []model.SysDept{}
		}
		mapItem = append(mapItem, item)
		nodesMap[parentID] = mapItem
		// 记录节点ID
		treeIds = append(treeIds, item.DeptID)
	}

	for key, value := range nodesMap {
		// 选择不是节点ID的作为树节点
		found := false
		for _, id := range treeIds {
			if id == key {
				found = true
				break
			}
		}
		if !found {
			tree = append(tree, value...)
		}
	}

	for i, node := range tree {
		iN := r.parseDataToTreeComponet(node, &nodesMap)
		tree[i] = iN
	}

	return tree
}

// parseDataToTreeComponet 递归函数处理子节点
func (r *sysDeptImpl) parseDataToTreeComponet(node model.SysDept, nodesMap *map[string][]model.SysDept) model.SysDept {
	id := node.DeptID
	children, ok := (*nodesMap)[id]
	if ok {
		node.Children = children
	}
	if len(node.Children) > 0 {
		for i, child := range node.Children {
			icN := r.parseDataToTreeComponet(child, nodesMap)
			node.Children[i] = icN
		}
	}
	return node
}
