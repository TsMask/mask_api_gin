package service

import (
	"mask_api_gin/src/framework/constants/common"
	"mask_api_gin/src/framework/vo"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/repository"
	"strings"
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
	return r.sysDeptRepository.SelectDeptById(deptId)
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
	newParentDept := r.sysDeptRepository.SelectDeptById(sysDept.ParentID)
	oldDept := r.sysDeptRepository.SelectDeptById(sysDept.DeptID)
	// 修改子元素关系
	if newParentDept.DeptID == sysDept.ParentID && oldDept.DeptID == sysDept.DeptID {
		newAncestors := newParentDept.Ancestors + "," + newParentDept.DeptID
		oldAncestors := oldDept.Ancestors
		sysDept.Ancestors = newAncestors
		r.updateDeptChildren(sysDept.DeptID, newAncestors, oldAncestors)
	}
	// 如果该部门是启用状态，则启用该部门的所有上级部门
	if sysDept.Status == common.STATUS_YES && sysDept.Ancestors != "0" {
		r.updateDeptStatusNormal(sysDept.Ancestors)
	}
	return r.sysDeptRepository.UpdateDept(sysDept)
}

// updateDeptStatusNormal 修改所在部门正常状态
func (r *sysDeptImpl) updateDeptStatusNormal(ancestors string) int64 {
	if ancestors == "" || ancestors == "0" {
		return 0
	}
	deptIds := strings.Split(ancestors, ",")
	return r.sysDeptRepository.UpdateDeptStatusNormal(deptIds)
}

// updateDeptChildren 修改子元素关系
func (r *sysDeptImpl) updateDeptChildren(deptId, newAncestors, oldAncestors string) int64 {
	childrens := r.sysDeptRepository.SelectChildrenDeptById(deptId)
	if len(childrens) == 0 {
		return 0
	}
	// 替换父ID
	for _, child := range childrens {
		ancestors := strings.Replace(child.Ancestors, oldAncestors, newAncestors, -1)
		child.Ancestors = ancestors
	}
	return r.sysDeptRepository.UpdateDeptChildren(childrens)
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
