package repository

import (
	"fmt"
	"mask_api_gin/src/framework/database/db"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/modules/system/model"
	"strings"
	"time"
)

// NewSysDept 实例化数据层
var NewSysDept = &SysDept{}

// SysDept 部门表 数据层处理
type SysDept struct{}

// Select 查询集合
func (r SysDept) Select(sysDept model.SysDept, dataScopeSQL string) []model.SysDept {
	tx := db.DB("").Model(&model.SysDept{})
	tx = tx.Where("del_flag = '0'")
	// 查询条件拼接
	if sysDept.DeptId > 0 {
		tx = tx.Where("dept_id = ?", sysDept.DeptId)
	}
	if sysDept.ParentId > 0 {
		tx = tx.Where("parent_id = ?", sysDept.ParentId)
	}
	if sysDept.DeptName != "" {
		tx = tx.Where("dept_name like concat(?, '%')", sysDept.DeptName)
	}
	if sysDept.StatusFlag != "" {
		tx = tx.Where("status_flag = ?", sysDept.StatusFlag)
	}

	// 查询数据
	rows := []model.SysDept{}
	if err := tx.Order("parent_id, dept_sort asc").Find(&rows).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return rows
	}
	return rows
}

// SelectById 通过ID查询信息
func (r SysDept) SelectById(deptId int64) model.SysDept {
	if deptId <= 0 {
		return model.SysDept{}
	}
	tx := db.DB("").Model(&model.SysDept{})
	// 构建查询条件
	tx = tx.Where("dept_id = ? and del_flag = '0'", deptId)
	// 查询数据
	item := model.SysDept{}
	if err := tx.Limit(1).Find(&item).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return item
	}
	return item
}

// Insert 新增信息 返回新增数据ID
func (r SysDept) Insert(sysDept model.SysDept) int64 {
	sysDept.DelFlag = "0"
	if sysDept.CreateBy != "" {
		sysDept.CreateTime = time.Now().UnixMilli()
	}
	// 执行插入
	if err := db.DB("").Create(&sysDept).Error; err != nil {
		logger.Errorf("insert err => %v", err.Error())
		return 0
	}
	return sysDept.DeptId
}

// Update 修改信息 返回受影响行数
func (r SysDept) Update(sysDept model.SysDept) int64 {
	if sysDept.DeptId <= 0 {
		return 0
	}
	if sysDept.UpdateBy != "" {
		sysDept.UpdateTime = time.Now().UnixMilli()
	}
	tx := db.DB("").Model(&model.SysDept{})
	// 构建查询条件
	tx = tx.Where("dept_id = ?", sysDept.DeptId)
	// 执行更新
	if err := tx.Updates(sysDept).Error; err != nil {
		logger.Errorf("update err => %v", err.Error())
		return 0
	}
	return tx.RowsAffected
}

// DeleteById 删除信息 返回受影响行数
func (r SysDept) DeleteById(deptId int64) int64 {
	if deptId <= 0 {
		return 0
	}
	tx := db.DB("").Model(&model.SysDept{})
	// 构建查询条件
	tx = tx.Where("dept_id in ?", deptId)
	// 执行更新删除标记
	if err := tx.Update("del_flag", "1").Error; err != nil {
		logger.Errorf("update err => %v", err.Error())
		return 0
	}
	return tx.RowsAffected
}

// CheckUnique 检查信息是否唯一 返回数据ID
func (r SysDept) CheckUnique(sysDept model.SysDept) int64 {
	tx := db.DB("").Model(&model.SysDept{})
	tx = tx.Where("del_flag = 0")
	// 查询条件拼接
	if sysDept.DeptName != "" {
		tx = tx.Where("dept_name = ?", sysDept.DeptName)
	}
	if sysDept.ParentId > 0 {
		tx = tx.Where("parent_id = ?", sysDept.ParentId)
	}

	// 查询数据
	var id int64 = 0
	if err := tx.Select("dept_id").Limit(1).Find(&id).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return id
	}
	return id
}

// ExistChildrenByDeptId 存在子节点数量
func (r SysDept) ExistChildrenByDeptId(deptId int64) int64 {
	if deptId <= 0 {
		return 0
	}
	tx := db.DB("").Model(&model.SysDept{})
	tx = tx.Where("del_flag = '0' and parent_id = ?", deptId)
	// 查询数据
	var count int64 = 0
	if err := tx.Count(&count).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return count
	}
	return count
}

// ExistUserByDeptId 存在用户使用数量
func (r SysDept) ExistUserByDeptId(deptId int64) int64 {
	if deptId <= 0 {
		return 0
	}
	tx := db.DB("").Model(&model.SysUser{})
	tx = tx.Where("del_flag = '0' and dept_id = ?", deptId)
	// 查询数据
	var count int64 = 0
	if err := tx.Count(&count).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return count
	}
	return count
}

// SelectDeptIdsByRoleId 通过角色ID查询包含的部门ID
func (r SysDept) SelectDeptIdsByRoleId(roleId int64, deptCheckStrictly bool) []string {
	if roleId <= 0 {
		return []string{}
	}

	tx := db.DB("").Model(&model.SysDept{})
	tx = tx.Where("del_flag = '0'")
	tx = tx.Where("dept_id in (SELECT DISTINCT dept_id FROM sys_role_dept WHERE role_id = ?)", roleId)
	// 展开
	if deptCheckStrictly {
		tx = tx.Where(`dept_id not in (
		SELECT d.parent_id FROM sys_dept d 
		INNER JOIN sys_role_dept rd ON rd.dept_id = d.dept_id 
		AND rd.role_id = ?
		)`, roleId)
	}

	// 查询数据
	rows := []string{}
	if err := tx.Distinct("dept_id").Find(&rows).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return rows
	}
	return rows
}

// SelectChildrenDeptById 根据ID查询所有子部门
func (r SysDept) SelectChildrenDeptById(deptId int64) []model.SysDept {
	tx := db.DB("").Model(&model.SysDept{})
	tx = tx.Where("del_flag = 0")
	tx = tx.Where("find_in_set(?, ancestors)", deptId)
	// 查询数据
	rows := []model.SysDept{}
	if err := tx.Find(&rows).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return rows
	}
	return rows
}

// UpdateDeptStatusNormal 修改所在部门正常状态
func (r SysDept) UpdateDeptStatusNormal(deptIds []int64) int64 {
	if len(deptIds) <= 0 {
		return 0
	}
	tx := db.DB("").Model(&model.SysDept{})
	// 构建查询条件
	tx = tx.Where("dept_id in ?", deptIds)
	// 执行更新状态标记
	if err := tx.Update("status_flag", "1").Error; err != nil {
		logger.Errorf("update err => %v", err.Error())
		return 0
	}
	return tx.RowsAffected
}

// UpdateDeptChildren 修改子元素关系
func (r SysDept) UpdateDeptChildren(arr []model.SysDept) int64 {
	if len(arr) == 0 {
		return 0
	}

	// 更新条件拼接
	var conditions []string
	var params []any
	for _, dept := range arr {
		caseSql := fmt.Sprintf("WHEN dept_id = '%d' THEN '%s'", dept.DeptId, dept.Ancestors)
		conditions = append(conditions, caseSql)
		params = append(params, dept.DeptId)
	}

	cases := strings.Join(conditions, " ")
	placeholders := db.KeyPlaceholderByQuery(len(params))
	sql := fmt.Sprintf("update sys_dept set ancestors = CASE %s END where dept_id in (%s)", cases, placeholders)
	results, err := db.ExecDB("", sql, params)
	if err != nil {
		logger.Errorf("delete err => %v", err)
		return 0
	}
	return results
}
