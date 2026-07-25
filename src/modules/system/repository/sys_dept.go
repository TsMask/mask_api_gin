package repository

import (
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
	tx := database.DB("").Model(&model.SysDept{})
	tx = tx.Where("dept_del = '0'")
	// 查询条件拼接
	if sysDept.ID > 0 {
		tx = tx.Where("id = ?", sysDept.ID)
	}
	if sysDept.ParentId > 0 {
		tx = tx.Where("parent_id = ?", sysDept.ParentId)
	}
	if sysDept.DeptName != "" {
		tx = tx.Where("dept_name like ?", sysDept.DeptName+"%")
	}
	if sysDept.DeptStatus != "" {
		tx = tx.Where("dept_status = ?", sysDept.DeptStatus)
	}
	if dataScopeSQL != "" {
		tx = tx.Where(dataScopeSQL)
	}

	// 查询数据
	rows := []model.SysDept{}
	if err := tx.Order("parent_id, dept_sort asc").Find(&rows).Error; err != nil {
		return rows, err
	}
	return rows, nil
}

// SelectById 通过ID查询信息
func (r SysDept) SelectById(deptId int64) (model.SysDept, error) {
	if deptId <= 0 {
		return model.SysDept{}, nil
	}
	tx := database.DB("").Model(&model.SysDept{})
	// 构建查询条件
	tx = tx.Where("id = ? and dept_del = '0'", deptId)
	// 查询数据
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
	tx := database.DB("").Model(&model.SysDept{})
	tx = tx.Where("id in ? and dept_del = '0'", deptIds)
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
	sysDept.DeptDel = "0"
	// 执行插入
	if err := database.DB("").Create(&sysDept).Error; err != nil {
		return 0, err
	}
	return sysDept.ID, nil
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
	tx = tx.Omit("id", "create_by", "create_time", "dept_del")
	// 执行更新
	if err := tx.Updates(sysDept).Error; err != nil {
		return 0, err
	}
	return tx.RowsAffected, nil
}

// DeleteById 删除信息 返回受影响行数（软删除）
func (r SysDept) DeleteById(d *gorm.DB, deptId int64) (int64, error) {
	if deptId <= 0 {
		return 0, nil
	}
	if d == nil {
		d = database.DB("")
	}
	tx := d.Model(&model.SysDept{})
	// 构建查询条件
	tx = tx.Where("id = ?", deptId)
	// 执行更新删除标记
	if err := tx.Update("dept_del", "1").Error; err != nil {
		return 0, err
	}
	return tx.RowsAffected, nil
}

// CheckUnique 检查信息是否唯一 返回数据ID
func (r SysDept) CheckUnique(sysDept model.SysDept) (int64, error) {
	tx := database.DB("").Model(&model.SysDept{})
	tx = tx.Where("dept_del = '0'")
	// 查询条件拼接
	if sysDept.DeptName != "" {
		tx = tx.Where("dept_name = ?", sysDept.DeptName)
	}
	if sysDept.ParentId > 0 {
		tx = tx.Where("parent_id = ?", sysDept.ParentId)
	}

	// 查询数据
	var id int64 = 0
	if err := tx.Select("id").Limit(1).Find(&id).Error; err != nil {
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
	tx = tx.Where("dept_del = '0' and dept_status = '1'")
	tx = tx.Where("parent_id = ?", deptId)
	// 查询数据
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
	tx = tx.Where("user_del = '0' and dept_id = ?", deptId)
	// 查询数据
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
	tx = tx.Where("dept_del = '0'")
	tx = tx.Where("id in (SELECT DISTINCT dept_id FROM sys_role_dept WHERE role_id = ?)", roleId)
	// 父子互相关联显示，取所有子节点
	if deptCheckStrictly {
		tx = tx.Where(`id not in (
		SELECT d.parent_id FROM sys_dept d
		INNER JOIN sys_role_dept rd ON rd.dept_id = d.id
		AND rd.role_id = ?
		)`, roleId)
	}

	// 查询数据
	rows := []int64{}
	if err := tx.Distinct("id").Find(&rows).Error; err != nil {
		return rows, err
	}
	return rows, nil
}

// SelectChildrenDeptById 根据ID查询所有子部门
func (r SysDept) SelectChildrenDeptById(deptId int64) ([]model.SysDept, error) {
	tx := database.DB("").Model(&model.SysDept{})
	tx = tx.Where("dept_del = '0'")
	tx = tx.Where("find_in_set(?, ancestors)", deptId)
	// 查询数据
	rows := []model.SysDept{}
	if err := tx.Find(&rows).Error; err != nil {
		return rows, err
	}
	return rows, nil
}

// UpdateDeptStatusNormal 修改所在部门正常状态
func (r SysDept) UpdateDeptStatusNormal(d *gorm.DB, deptIds []int64) (int64, error) {
	if len(deptIds) <= 0 {
		return 0, nil
	}
	if d == nil {
		d = database.DB("")
	}
	tx := d.Model(&model.SysDept{})
	// 构建查询条件
	tx = tx.Where("id in ?", deptIds)
	// 执行更新状态标记
	if err := tx.Update("dept_status", "1").Error; err != nil {
		return 0, err
	}
	return tx.RowsAffected, nil
}

// UpdateDeptChildren 修改子元素关系
func (r SysDept) UpdateDeptChildren(d *gorm.DB, arr []model.SysDept) (int64, error) {
	if len(arr) == 0 {
		return 0, nil
	}

	// 逐条更新子部门的 ancestors，使用参数化避免 SQL 注入
	if d == nil {
		d = database.DB("")
	}
	var affected int64
	for _, dept := range arr {
		if dept.ID <= 0 {
			continue
		}
		tx := d.Model(&model.SysDept{})
		tx = tx.Where("id = ?", dept.ID)
		if err := tx.Update("ancestors", dept.Ancestors).Error; err != nil {
			return affected, err
		}
		affected += tx.RowsAffected
	}
	return affected, nil
}
