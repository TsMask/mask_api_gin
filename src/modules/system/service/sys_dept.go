package service

import (
	"fmt"
	"strings"

	"mask_api_gin/src/framework/constants"
	"mask_api_gin/src/framework/database/db"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/model/vo"
	"mask_api_gin/src/modules/system/repository"

	"gorm.io/gorm"
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
func (s SysDept) FindById(deptId int64) model.SysDept {
	return s.sysDeptRepository.SelectById(deptId)
}

// Insert 新增信息
func (s SysDept) Insert(sysDept model.SysDept) int64 {
	return s.sysDeptRepository.Insert(sysDept)
}

// Update 修改信息
func (s SysDept) Update(sysDept model.SysDept) int64 {
	dept := s.sysDeptRepository.SelectById(sysDept.ID)
	parentDept := s.sysDeptRepository.SelectById(sysDept.ParentId)

	var rows int64
	db.Transaction("", func(tx *gorm.DB) error {
		// 上级与当前部门祖级列表更新
		if parentDept.ID == sysDept.ParentId && dept.ID == sysDept.ID {
			newAncestors := fmt.Sprintf("%s,%d", parentDept.Ancestors, parentDept.ID)
			oldAncestors := dept.Ancestors
			// 祖级列表不一致时更新
			if newAncestors != oldAncestors {
				s.updateDeptChildren(tx, dept.ID, newAncestors, oldAncestors)
			}
		}
		// 如果该部门是启用状态，则启用该部门的所有上级部门
		if sysDept.DeptStatus == constants.STATUS_YES && parentDept.DeptStatus == constants.STATUS_NO {
			s.updateDeptStatusNormal(tx, sysDept.Ancestors)
		}
		rows = s.sysDeptRepository.Update(nil, sysDept)
		return nil
	})
	return rows
}

// updateDeptStatusNormal 修改所在部门正常状态
func (s SysDept) updateDeptStatusNormal(d *gorm.DB, ancestors string) int64 {
	if ancestors == "" || ancestors == "0" {
		return 0
	}
	deptIds := make([]int64, 0)
	for _, v := range strings.Split(ancestors, ",") {
		deptIds = append(deptIds, parse.Number(v))
	}
	return s.sysDeptRepository.UpdateDeptStatusNormal(d, deptIds)
}

// updateDeptChildren 修改子元素关系
func (s SysDept) updateDeptChildren(d *gorm.DB, deptId int64, newAncestors, oldAncestors string) int64 {
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
	return s.sysDeptRepository.UpdateDeptChildren(d, arr)
}

// DeleteById 删除信息
func (s SysDept) DeleteById(deptId int64) int64 {
	var rows int64
	db.Transaction("", func(tx *gorm.DB) error {
		s.sysRoleDeptRepository.DeleteByDeptIds(tx, []int64{deptId})
		rows = s.sysDeptRepository.DeleteById(tx, deptId)
		return nil
	})
	return rows
}

// FindDeptIdsByRoleId 根据角色ID查询包含的部门ID
func (s SysDept) FindDeptIdsByRoleId(roleId int64) []int64 {
	roles := s.sysRoleRepository.SelectByIds([]int64{roleId})
	if len(roles) > 0 {
		role := roles[0]
		if role.ID == roleId {
			return s.sysDeptRepository.SelectDeptIdsByRoleId(
				role.ID,
				role.RoleDeptCheckStrictly == "1",
			)
		}
	}
	return []int64{}
}

// ExistChildrenByDeptId 部门下存在子节点数量
func (s SysDept) ExistChildrenByDeptId(deptId int64) int64 {
	return s.sysDeptRepository.ExistChildrenByDeptId(deptId)
}

// ExistUserByDeptId 部门下存在用户数量
func (s SysDept) ExistUserByDeptId(deptId int64) int64 {
	return s.sysDeptRepository.ExistUserByDeptId(deptId)
}

// CheckUniqueParentIdByDeptName 检查同级下部门名称唯一
func (s SysDept) CheckUniqueParentIdByDeptName(parentId int64, deptName string, deptId int64) bool {
	uniqueId := s.sysDeptRepository.CheckUnique(model.SysDept{
		DeptName: deptName,
		ParentId: parentId,
	})
	if uniqueId == deptId {
		return true
	}
	return uniqueId == 0
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
	nodesMap := make(map[int64][]model.SysDept)
	// 节点id
	treeIds := make([]int64, 0)
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
		treeIds = append(treeIds, item.ID)
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
func (s SysDept) parseDataToTreeComponent(node model.SysDept, nodesMap *map[int64][]model.SysDept) model.SysDept {
	if children, ok := (*nodesMap)[node.ID]; ok {
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
