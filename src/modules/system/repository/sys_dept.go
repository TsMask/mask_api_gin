package repository

import (
	"fmt"

	"mask_api_gin/src/internal/infra/database"
	"mask_api_gin/src/modules/system/model"

	"gorm.io/gorm"
)

// NewSysDept 实例化数据层
var NewSysDept = &SysDept{}

// SysDept 部门表 数据层处理
type SysDept struct{}

// Select 查询集合
func (r SysDept) Select(sysDept model.SysDept, dataScopeSQL string) ([]model.SysDept, error) {
	tx := database.DB("").Model(&model.SysDept{}).
		Select("sys_dept.*, COALESCE(parent_rel.ancestor_id, 0) AS parent_id").
		Joins("LEFT JOIN sys_dept_closure parent_rel ON parent_rel.descendant_id = sys_dept.id AND parent_rel.depth = 1")

	if sysDept.ID > 0 {
		tx = tx.Where("sys_dept.id = ?", sysDept.ID)
	}
	if sysDept.ParentId > 0 {
		tx = tx.Where("parent_rel.ancestor_id = ?", sysDept.ParentId)
	}
	if sysDept.DeptName != "" {
		tx = tx.Where("sys_dept.dept_name like ?", sysDept.DeptName+"%")
	}
	if sysDept.DeptStatus != "" {
		tx = tx.Where("sys_dept.dept_status = ?", sysDept.DeptStatus)
	}
	if dataScopeSQL != "" {
		tx = tx.Where(dataScopeSQL)
	}

	rows := []model.SysDept{}
	if err := tx.Order("sys_dept.dept_sort, sys_dept.id").Find(&rows).Error; err != nil {
		return rows, err
	}
	return rows, nil
}

// SelectExcludeChildren 查询排除指定部门及其子部门的集合
func (r SysDept) SelectExcludeChildren(deptId int64, dataScopeSQL string) ([]model.SysDept, error) {
	tx := database.DB("").Model(&model.SysDept{}).
		Select("sys_dept.*, COALESCE(parent_rel.ancestor_id, 0) AS parent_id").
		Joins("LEFT JOIN sys_dept_closure parent_rel ON parent_rel.descendant_id = sys_dept.id AND parent_rel.depth = 1")

	tx = tx.Where("sys_dept.id <> ?", deptId)
	tx = tx.Where(`NOT EXISTS (
		SELECT 1 FROM sys_dept_closure exclude_rel
		WHERE exclude_rel.ancestor_id = ? AND exclude_rel.descendant_id = sys_dept.id
	)`, deptId)
	if dataScopeSQL != "" {
		tx = tx.Where(dataScopeSQL)
	}

	rows := []model.SysDept{}
	if err := tx.Order("sys_dept.dept_sort, sys_dept.id").Find(&rows).Error; err != nil {
		return rows, err
	}
	return rows, nil
}

// SelectById 通过ID查询信息
func (r SysDept) SelectById(deptId int64) (model.SysDept, error) {
	if deptId <= 0 {
		return model.SysDept{}, nil
	}
	tx := database.DB("").Model(&model.SysDept{}).
		Select("sys_dept.*, COALESCE(parent_rel.ancestor_id, 0) AS parent_id").
		Joins("LEFT JOIN sys_dept_closure parent_rel ON parent_rel.descendant_id = sys_dept.id AND parent_rel.depth = 1")

	tx = tx.Where("sys_dept.id = ?", deptId)

	item := model.SysDept{}
	if err := tx.Limit(1).Find(&item).Error; err != nil {
		return item, err
	}
	return item, nil
}

// SelectMapByIds 批量通过ID查询信息，返回 deptId → SysDept 的映射
func (r SysDept) SelectMapByIds(deptIds []int64) (map[int64]model.SysDept, error) {
	result := make(map[int64]model.SysDept)
	if len(deptIds) <= 0 {
		return result, nil
	}
	tx := database.DB("").Model(&model.SysDept{}).
		Select("sys_dept.*, COALESCE(parent_rel.ancestor_id, 0) AS parent_id").
		Joins("LEFT JOIN sys_dept_closure parent_rel ON parent_rel.descendant_id = sys_dept.id AND parent_rel.depth = 1")

	tx = tx.Where("sys_dept.id IN ?", deptIds)

	rows := []model.SysDept{}
	if err := tx.Find(&rows).Error; err != nil {
		return result, err
	}
	for _, dept := range rows {
		result[dept.ID] = dept
	}
	return result, nil
}

// Insert 新增信息 返回新增数据ID
func (r SysDept) Insert(sysDept model.SysDept) (int64, error) {
	err := database.Transaction("", func(tx *gorm.DB) error {
		if err := tx.Create(&sysDept).Error; err != nil {
			return err
		}
		return insertDeptClosure(tx, sysDept.ID, sysDept.ParentId)
	})
	if err != nil {
		return 0, err
	}
	return sysDept.ID, nil
}

func insertDeptClosure(tx *gorm.DB, deptId, parentId int64) error {
	if deptId <= 0 {
		return fmt.Errorf("dept id is empty")
	}

	relations := []model.SysDeptClosure{
		{AncestorID: deptId, DescendantID: deptId, Depth: 0},
	}
	if parentId > 0 {
		parentRelations := []model.SysDeptClosure{}
		if err := tx.Raw(`
			SELECT ancestor_id, descendant_id, depth
			FROM sys_dept_closure
			WHERE descendant_id = ?
			FOR UPDATE
		`, parentId).Scan(&parentRelations).Error; err != nil {
			return err
		}
		if len(parentRelations) == 0 {
			return fmt.Errorf("parent dept closure not found: %d", parentId)
		}
		for _, relation := range parentRelations {
			relations = append(relations, model.SysDeptClosure{
				AncestorID:   relation.AncestorID,
				DescendantID: deptId,
				Depth:        relation.Depth + 1,
			})
		}
	}
	return tx.Create(&relations).Error
}

// Update 修改信息 返回受影响行数
func (r SysDept) Update(d *gorm.DB, sysDept model.SysDept) (int64, error) {
	if sysDept.ID <= 0 {
		return 0, nil
	}
	if d == nil {
		d = database.DB("")
	}
	tx := d.Model(&model.SysDept{})
	// 构建查询条件
	tx = tx.Where("id = ?", sysDept.ID)
	tx = tx.Omit("id", "create_by", "create_time", "parent_id")
	// 执行更新
	if err := tx.Updates(sysDept).Error; err != nil {
		return 0, err
	}
	return tx.RowsAffected, nil
}

// MoveDeptSubtree 移动部门及其子树到新的父级下
func (r SysDept) MoveDeptSubtree(d *gorm.DB, deptId, newParentId int64) error {
	if d == nil {
		return fmt.Errorf("database transaction is nil")
	}
	if deptId <= 0 {
		return fmt.Errorf("dept id is empty")
	}

	var currentParentId int64
	if err := d.Model(&model.SysDeptClosure{}).
		Select("ancestor_id").
		Where("ancestor_id <> descendant_id AND descendant_id = ? AND depth = 1", deptId).
		Scan(&currentParentId).Error; err != nil {
		return err
	}
	if currentParentId == newParentId {
		return nil
	}

	if newParentId > 0 {
		var descendantCount int64
		if err := d.Model(&model.SysDeptClosure{}).
			Where("ancestor_id = ? AND descendant_id = ? AND depth > 0", deptId, newParentId).
			Count(&descendantCount).Error; err != nil {
			return err
		}
		if descendantCount > 0 {
			return fmt.Errorf("new parent dept is a descendant of the moving dept")
		}
	}

	subtree := []model.SysDeptClosure{}
	if err := d.Raw(`
		SELECT ancestor_id, descendant_id, depth
		FROM sys_dept_closure
		WHERE ancestor_id = ?
		FOR UPDATE
	`, deptId).Scan(&subtree).Error; err != nil {
		return err
	}
	if len(subtree) == 0 {
		return fmt.Errorf("dept closure not found: %d", deptId)
	}

	subtreeIds := make([]int64, 0, len(subtree))
	for _, relation := range subtree {
		subtreeIds = append(subtreeIds, relation.DescendantID)
	}
	if err := d.Where("descendant_id IN ? AND ancestor_id NOT IN ?", subtreeIds, subtreeIds).
		Delete(&model.SysDeptClosure{}).Error; err != nil {
		return err
	}

	if newParentId <= 0 {
		return nil
	}

	parentRelations := []model.SysDeptClosure{}
	if err := d.Raw(`
		SELECT ancestor_id, descendant_id, depth
		FROM sys_dept_closure
		WHERE descendant_id = ?
		FOR UPDATE
	`, newParentId).Scan(&parentRelations).Error; err != nil {
		return err
	}
	if len(parentRelations) == 0 {
		return fmt.Errorf("parent dept closure not found: %d", newParentId)
	}

	relations := make([]model.SysDeptClosure, 0, len(parentRelations)*len(subtree))
	for _, parent := range parentRelations {
		for _, descendant := range subtree {
			relations = append(relations, model.SysDeptClosure{
				AncestorID:   parent.AncestorID,
				DescendantID: descendant.DescendantID,
				Depth:        parent.Depth + descendant.Depth + 1,
			})
		}
	}
	return d.CreateInBatches(relations, 500).Error
}

// DeleteById 删除信息 返回受影响行数
func (r SysDept) DeleteById(d *gorm.DB, deptId int64) (int64, error) {
	if deptId <= 0 {
		return 0, nil
	}
	if d == nil {
		d = database.DB("")
	}

	var childCount int64
	if err := d.Model(&model.SysDept{}).
		Joins("JOIN sys_dept_closure c ON c.descendant_id = sys_dept.id").
		Where("c.ancestor_id = ? AND c.depth = 1", deptId).
		Count(&childCount).Error; err != nil {
		return 0, err
	}
	if childCount > 0 {
		return 0, fmt.Errorf("dept contains child departments: %d", childCount)
	}

	if err := d.Where("ancestor_id = ? OR descendant_id = ?", deptId, deptId).
		Delete(&model.SysDeptClosure{}).Error; err != nil {
		return 0, err
	}
	tx := d.Where("id = ?", deptId).Delete(&model.SysDept{})
	if tx.Error != nil {
		return 0, tx.Error
	}
	return tx.RowsAffected, nil
}

// CheckUnique 检查同级部门名称唯一
func (r SysDept) CheckUnique(sysDept model.SysDept) (int64, error) {
	tx := database.DB("").Model(&model.SysDept{}).
		Select("sys_dept.*, COALESCE(parent_rel.ancestor_id, 0) AS parent_id").
		Joins("LEFT JOIN sys_dept_closure parent_rel ON parent_rel.descendant_id = sys_dept.id AND parent_rel.depth = 1")

	if sysDept.DeptName != "" {
		tx = tx.Where("sys_dept.dept_name = ?", sysDept.DeptName)
	}
	if sysDept.ParentId > 0 {
		tx = tx.Where("parent_rel.ancestor_id = ? AND parent_rel.depth = 1", sysDept.ParentId)
	} else {
		tx = tx.Where("parent_rel.ancestor_id IS NULL")
	}

	var id int64 = 0
	if err := tx.Select("sys_dept.id").Limit(1).Find(&id).Error; err != nil {
		return id, err
	}
	return id, nil
}

// ExistChildrenByDeptId 存在子节点数量
func (r SysDept) ExistChildrenByDeptId(deptId int64) (int64, error) {
	if deptId <= 0 {
		return 0, nil
	}
	tx := database.DB("").Model(&model.SysDept{})
	tx = tx.Joins("JOIN sys_dept_closure c ON c.descendant_id = sys_dept.id")
	tx = tx.Where("c.ancestor_id = ? AND c.depth = 1", deptId)

	var count int64 = 0
	if err := tx.Count(&count).Error; err != nil {
		return count, err
	}
	return count, nil
}

// ExistUserByDeptId 存在用户使用数量
func (r SysDept) ExistUserByDeptId(deptId int64) (int64, error) {
	if deptId <= 0 {
		return 0, nil
	}
	tx := database.DB("").Model(&model.SysUser{})
	tx = tx.Where("deleted_at = 0 AND dept_id = ?", deptId)

	var count int64 = 0
	if err := tx.Count(&count).Error; err != nil {
		return count, err
	}
	return count, nil
}

// SelectDeptIdsByRoleId 通过角色ID查询包含的部门ID
func (r SysDept) SelectDeptIdsByRoleId(roleId int64, deptCheckStrictly bool) ([]int64, error) {
	if roleId <= 0 {
		return []int64{}, nil
	}

	tx := database.DB("").Model(&model.SysDept{})
	tx = tx.Where("sys_dept.id IN (SELECT DISTINCT dept_id FROM sys_role_dept WHERE role_id = ?)", roleId)
	if deptCheckStrictly {
		tx = tx.Where(`sys_dept.id NOT IN (
			SELECT c.ancestor_id
			FROM sys_dept_closure c
			INNER JOIN sys_role_dept rd ON rd.dept_id = c.descendant_id AND rd.role_id = ?
			WHERE c.depth = 1
		)`, roleId)
	}

	rows := []int64{}
	if err := tx.Distinct("sys_dept.id").Find(&rows).Error; err != nil {
		return rows, err
	}
	return rows, nil
}

// SelectChildrenDeptById 根据ID查询所有子部门
func (r SysDept) SelectChildrenDeptById(deptId int64) ([]model.SysDept, error) {
	tx := database.DB("").Model(&model.SysDept{}).
		Select("sys_dept.*, COALESCE(parent_rel.ancestor_id, 0) AS parent_id").
		Joins("LEFT JOIN sys_dept_closure parent_rel ON parent_rel.descendant_id = sys_dept.id AND parent_rel.depth = 1")

	tx = tx.Joins("JOIN sys_dept_closure child_rel ON child_rel.descendant_id = sys_dept.id").
		Where("child_rel.ancestor_id = ? AND child_rel.depth > 0", deptId)

	rows := []model.SysDept{}
	if err := tx.Order("sys_dept.dept_sort, sys_dept.id").Find(&rows).Error; err != nil {
		return rows, err
	}
	return rows, nil
}

// UpdateDeptStatusNormal 启用指定部门的全部祖先
func (r SysDept) UpdateDeptStatusNormal(d *gorm.DB, deptId int64) (int64, error) {
	if deptId <= 0 {
		return 0, nil
	}
	if d == nil {
		d = database.DB("")
	}
	result := d.Exec(`
		UPDATE sys_dept d
		INNER JOIN sys_dept_closure c ON c.descendant_id = d.id
		SET d.dept_status = '1'
		WHERE c.ancestor_id <> c.descendant_id
		  AND c.ancestor_id = ?
	`, deptId)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
