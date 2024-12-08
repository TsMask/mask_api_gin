package service

import (
	"mask_api_gin/src/framework/constants"
	"mask_api_gin/src/framework/vo"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/repository"

	"fmt"
	"strings"
)

// NewSysDept  实例化服务层
var NewSysDept = &SysDept{
	sysDeptRepository:     repository.NewSysDept,
	sysRoleRepository:     repository.NewSysRole,
	sysRoleDeptRepository: repository.NewSysRoleDept,
}

// SysDept 部门管理 服务层处理
type SysDept struct {
	sysDeptRepository     *repository.SysDept     // 部门服务
	sysRoleRepository     *repository.SysRole     // 角色服务
	sysRoleDeptRepository *repository.SysRoleDept // 角色与部门关联服务
}

// Find 查询数据
func (s SysDept) Find(sysDept model.SysDept, dataScopeSQL string) []model.SysDept {
	return s.sysDeptRepository.Select(sysDept, dataScopeSQL)
}

// FindById 根据ID查询信息
func (s SysDept) FindById(deptId string) model.SysDept {
	return s.sysDeptRepository.SelectById(deptId)
}

// Insert 新增信息
func (s SysDept) Insert(sysDept model.SysDept) string {
	return s.sysDeptRepository.Insert(sysDept)
}

// Update 修改信息
func (s SysDept) Update(sysDept model.SysDept) int64 {
	dept := s.sysDeptRepository.SelectById(sysDept.DeptId)
	parentDept := s.sysDeptRepository.SelectById(sysDept.ParentId)
	// 上级与当前部门祖级列表更新
	if parentDept.DeptId == sysDept.ParentId && dept.DeptId == sysDept.DeptId {
		newAncestors := fmt.Sprintf("%s,%s", parentDept.Ancestors, parentDept.DeptId)
		oldAncestors := dept.Ancestors
		// 祖级列表不一致时更新
		if newAncestors != oldAncestors {
			dept.Ancestors = newAncestors
			s.updateDeptChildren(dept.DeptId, newAncestors, oldAncestors)
		}
	}
	// 如果该部门是启用状态，则启用该部门的所有上级部门
	if sysDept.StatusFlag == constants.STATUS_YES && parentDept.StatusFlag == constants.STATUS_NO {
		s.updateDeptStatusNormal(sysDept.Ancestors)
	}
	return s.sysDeptRepository.Update(sysDept)
}

// updateDeptStatusNormal 修改所在部门正常状态
func (s SysDept) updateDeptStatusNormal(ancestors string) int64 {
	if ancestors == "" || ancestors == "0" {
		return 0
	}
	deptIds := strings.Split(ancestors, ",")
	return s.sysDeptRepository.UpdateDeptStatusNormal(deptIds)
}

// updateDeptChildren 修改子元素关系
func (s SysDept) updateDeptChildren(deptId string, newAncestors, oldAncestors string) int64 {
	arr := s.sysDeptRepository.SelectChildrenDeptById(deptId)
	if len(arr) == 0 {
		return 0
	}
	// 替换父ID
	for i := range arr {
		item := &arr[i]
		ancestors := strings.Replace(item.Ancestors, oldAncestors, newAncestors, 1)
		item.Ancestors = ancestors
	}
	return s.sysDeptRepository.UpdateDeptChildren(arr)
}

// DeleteById 删除信息
func (s SysDept) DeleteById(deptId string) int64 {
	s.sysRoleDeptRepository.DeleteByDeptIds([]string{deptId}) // 删除角色与部门关联
	return s.sysDeptRepository.DeleteById(deptId)
}

// FindDeptIdsByRoleId 根据角色ID查询包含的部门ID
func (s SysDept) FindDeptIdsByRoleId(roleId string) []string {
	roles := s.sysRoleRepository.SelectByIds([]string{roleId})
	if len(roles) > 0 {
		role := roles[0]
		if role.RoleId == roleId {
			return s.sysDeptRepository.SelectDeptIdsByRoleId(
				role.RoleId,
				role.DeptCheckStrictly == "1",
			)
		}
	}
	return []string{}
}

// ExistChildrenByDeptId 部门下存在子节点数量
func (s SysDept) ExistChildrenByDeptId(deptId string) int64 {
	return s.sysDeptRepository.ExistChildrenByDeptId(deptId)
}

// ExistUserByDeptId 部门下存在用户数量
func (s SysDept) ExistUserByDeptId(deptId string) int64 {
	return s.sysDeptRepository.ExistUserByDeptId(deptId)
}

// CheckUniqueParentIdByDeptName 检查同级下部门名称唯一
func (s SysDept) CheckUniqueParentIdByDeptName(parentId string, deptName string, deptId string) bool {
	uniqueId := s.sysDeptRepository.CheckUnique(model.SysDept{
		DeptName: deptName,
		ParentId: parentId,
	})
	if uniqueId == deptId {
		return true
	}
	return uniqueId == ""
}

// BuildTreeSelect 查询部门树状结构
func (s SysDept) BuildTreeSelect(sysDept model.SysDept, dataScopeSQL string) []vo.TreeSelect {
	arr := s.sysDeptRepository.Select(sysDept, dataScopeSQL)
	treeArr := s.parseDataToTree(arr)
	tree := make([]vo.TreeSelect, 0)
	for _, item := range treeArr {
		tree = append(tree, vo.SysDeptTreeSelect(item))
	}
	return tree
}

// parseDataToTree 将数据解析为树结构，构建前端所需要下拉树结构
func (s SysDept) parseDataToTree(arr []model.SysDept) []model.SysDept {
	// 节点分组
	nodesMap := make(map[string][]model.SysDept)
	// 节点id
	treeIds := make([]string, 0)
	// 树节点
	tree := make([]model.SysDept, 0)

	for _, item := range arr {
		parentID := item.ParentId
		// 分组
		mapItem, ok := nodesMap[parentID]
		if !ok {
			mapItem = []model.SysDept{}
		}
		mapItem = append(mapItem, item)
		nodesMap[parentID] = mapItem
		// 记录节点ID
		treeIds = append(treeIds, item.DeptId)
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
		iN := s.parseDataToTreeComponent(node, &nodesMap)
		tree[i] = iN
	}

	return tree
}

// parseDataToTreeComponent 递归函数处理子节点
func (s SysDept) parseDataToTreeComponent(node model.SysDept, nodesMap *map[string][]model.SysDept) model.SysDept {
	id := node.DeptId
	children, ok := (*nodesMap)[id]
	if ok {
		node.Children = children
	}
	if len(node.Children) > 0 {
		for i, child := range node.Children {
			icN := s.parseDataToTreeComponent(child, nodesMap)
			node.Children[i] = icN
		}
	}
	return node
}
