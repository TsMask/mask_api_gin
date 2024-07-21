package service

import (
	constSystem "mask_api_gin/src/framework/constants/system"
	"mask_api_gin/src/framework/vo"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/repository"
	"strings"
)

// NewSysDept  实例化服务层
var NewSysDept = &SysDeptService{
	sysDeptRepository:     repository.NewSysDept,
	sysRoleRepository:     repository.NewSysRole,
	sysRoleDeptRepository: repository.NewSysRoleDept,
}

// SysDeptService 部门管理 服务层处理
type SysDeptService struct {
	sysDeptRepository     repository.ISysDeptRepository     // 部门服务
	sysRoleRepository     repository.ISysRoleRepository     // 角色服务
	sysRoleDeptRepository repository.ISysRoleDeptRepository // 角色与部门关联服务
}

// Find 查询数据
func (r *SysDeptService) Find(sysDept model.SysDept, dataScopeSQL string) []model.SysDept {
	return r.sysDeptRepository.Select(sysDept, dataScopeSQL)
}

// FindById 根据ID查询信息
func (r *SysDeptService) FindById(deptId string) model.SysDept {
	return r.sysDeptRepository.SelectById(deptId)
}

// Insert 新增信息
func (r *SysDeptService) Insert(sysDept model.SysDept) string {
	return r.sysDeptRepository.Insert(sysDept)
}

// Update 修改信息
func (r *SysDeptService) Update(sysDept model.SysDept) int64 {
	dept := r.sysDeptRepository.SelectById(sysDept.DeptID)
	parentDept := r.sysDeptRepository.SelectById(sysDept.ParentID)
	// 上级与当前部门祖级列表更新
	if parentDept.DeptID == sysDept.ParentID && dept.DeptID == sysDept.DeptID {
		newAncestors := parentDept.Ancestors + "," + parentDept.DeptID
		oldAncestors := dept.Ancestors
		// 祖级列表不一致时更新
		if newAncestors != oldAncestors {
			dept.Ancestors = newAncestors
			r.updateDeptChildren(dept.DeptID, newAncestors, oldAncestors)
		}
	}
	// 如果该部门是启用状态，则启用该部门的所有上级部门
	if sysDept.Status == constSystem.StatusYes && parentDept.Status == constSystem.StatusNo {
		r.updateDeptStatusNormal(sysDept.Ancestors)
	}
	return r.sysDeptRepository.Update(sysDept)
}

// updateDeptStatusNormal 修改所在部门正常状态
func (r *SysDeptService) updateDeptStatusNormal(ancestors string) int64 {
	if ancestors == "" || ancestors == "0" {
		return 0
	}
	deptIds := strings.Split(ancestors, ",")
	return r.sysDeptRepository.UpdateDeptStatusNormal(deptIds)
}

// updateDeptChildren 修改子元素关系
func (r *SysDeptService) updateDeptChildren(deptId, newAncestors, oldAncestors string) int64 {
	arr := r.sysDeptRepository.SelectChildrenDeptById(deptId)
	if len(arr) == 0 {
		return 0
	}
	// 替换父ID
	for i := range arr {
		item := &arr[i]
		ancestors := strings.Replace(item.Ancestors, oldAncestors, newAncestors, 1)
		item.Ancestors = ancestors
	}
	return r.sysDeptRepository.UpdateDeptChildren(arr)
}

// DeleteById 删除信息
func (r *SysDeptService) DeleteById(deptId string) int64 {
	r.sysRoleDeptRepository.DeleteByDeptIds([]string{deptId}) // 删除角色与部门关联
	return r.sysDeptRepository.DeleteById(deptId)
}

// FindDeptIdsByRoleId 根据角色ID查询包含的部门ID TODO
func (r *SysDeptService) FindDeptIdsByRoleId(roleId string) []string {
	roles := r.sysRoleRepository.SelectByIds([]string{roleId})
	if len(roles) == 0 {
		return []string{}
	}
	role := roles[0]
	if role.RoleID != roleId {
		return []string{}
	}
	return r.sysDeptRepository.SelectDeptIdsByRoleId(
		role.RoleID,
		role.DeptCheckStrictly == "1",
	)
}

// ExistChildrenByDeptId 部门下存在子节点数量
func (r *SysDeptService) ExistChildrenByDeptId(deptId string) int64 {
	return r.sysDeptRepository.ExistChildrenByDeptId(deptId)
}

// ExistUserByDeptId 部门下存在用户数量
func (r *SysDeptService) ExistUserByDeptId(deptId string) int64 {
	return r.sysDeptRepository.ExistUserByDeptId(deptId)
}

// CheckUniqueParentIdByDeptName 检查同级下部门名称唯一
func (r *SysDeptService) CheckUniqueParentIdByDeptName(parentId, deptName, deptId string) bool {
	uniqueId := r.sysDeptRepository.CheckUnique(model.SysDept{
		DeptName: deptName,
		ParentID: parentId,
	})
	if uniqueId == deptId {
		return true
	}
	return uniqueId == ""
}

// BuildTreeSelect 查询部门树状结构 TODO
func (r *SysDeptService) BuildTreeSelect(sysDept model.SysDept, dataScopeSQL string) []vo.TreeSelect {
	arr := r.sysDeptRepository.Select(sysDept, dataScopeSQL)
	treeArr := r.parseDataToTree(arr)
	tree := make([]vo.TreeSelect, 0)
	for _, item := range treeArr {
		tree = append(tree, vo.SysDeptTreeSelect(item))
	}
	return tree
}

// parseDataToTree 将数据解析为树结构，构建前端所需要下拉树结构
func (r *SysDeptService) parseDataToTree(arr []model.SysDept) []model.SysDept {
	// 节点分组
	nodesMap := make(map[string][]model.SysDept)
	// 节点id
	treeIds := make([]string, 0)
	// 树节点
	tree := make([]model.SysDept, 0)

	for _, item := range arr {
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
		iN := r.parseDataToTreeComponent(node, &nodesMap)
		tree[i] = iN
	}

	return tree
}

// parseDataToTreeComponent 递归函数处理子节点
func (r *SysDeptService) parseDataToTreeComponent(node model.SysDept, nodesMap *map[string][]model.SysDept) model.SysDept {
	id := node.DeptID
	children, ok := (*nodesMap)[id]
	if ok {
		node.Children = children
	}
	if len(node.Children) > 0 {
		for i, child := range node.Children {
			icN := r.parseDataToTreeComponent(child, nodesMap)
			node.Children[i] = icN
		}
	}
	return node
}
