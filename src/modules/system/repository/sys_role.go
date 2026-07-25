package repository

import (
	"fmt"

	"mask_api_gin/src/internal/infra/database/db"
	"mask_api_gin/src/modules/system/model"

	"gorm.io/gorm"
)

// NewSysRole 实例化数据层
var NewSysRole = &SysRole{}

// SysRole 角色表 数据层处理
type SysRole struct{}

// sysRoleSortMapping 角色排序字段映射
var sysRoleSortMapping = map[string]string{
	"id":            "id",
	"createTime":    "create_time",
	"updateTime":    "update_time",
	"roleName":      "role_name",
	"roleKey":       "role_key",
	"roleSort":      "role_sort",
	"roleStatus":    "role_status",
	"roleDataScope": "role_data_scope",
}

// SelectByPage 分页查询集合
func (r SysRole) SelectByPage(query map[string]string) ([]model.SysRole, int64, error) {
	tx := db.DB("").Model(&model.SysRole{})
	tx = tx.Where("role_del = '0'")
	// 查询条件拼接
	if v, ok := query["roleName"]; ok && v != "" {
		tx = tx.Where("role_name like ?", v+"%")
	}
	if v, ok := query["roleKey"]; ok && v != "" {
		tx = tx.Where("role_key like ?", v+"%")
	}
	if v, ok := query["roleStatus"]; ok && v != "" {
		tx = tx.Where("role_status = ?", v)
	}
	if v, ok := query["beginTime"]; ok && v != "" {
		if len(v) == 10 {
			v = fmt.Sprintf("%s000", v)
			tx = tx.Where("create_time >= ?", v)
		} else if len(v) == 13 {
			tx = tx.Where("create_time >= ?", v)
		}
	}
	if v, ok := query["endTime"]; ok && v != "" {
		if len(v) == 10 {
			v = fmt.Sprintf("%s999", v)
			tx = tx.Where("create_time <= ?", v)
		} else if len(v) == 13 {
			tx = tx.Where("create_time <= ?", v)
		}
	}
	if v, ok := query["deptId"]; ok && v != "" {
		tx = tx.Where(`id in (
		select distinct role_id from sys_role_dept where sys_role_dept.dept_id IN (
		SELECT id FROM sys_dept WHERE id = ? or find_in_set(? , ancestors )
		))`, v, v)
	}

	// 查询结果
	var total int64 = 0
	rows := []model.SysRole{}

	// 查询数量为0直接返回
	if err := tx.Count(&total).Error; err != nil || total <= 0 {
		return rows, total, nil
	}

	// 查询数据分页
	pageNum, pageSize := db.PageNumSize(query["pageNum"], query["pageSize"])
	tx = tx.Limit(pageSize).Offset(pageSize * pageNum)

	// 排序，默认按 role_sort ASC
	orders := db.SortParse(query["sort"], sysRoleSortMapping)
	tx = db.SortOrDefault(tx, orders, "role_sort asc")

	// 查询数据
	if err := tx.Find(&rows).Error; err != nil {
		return rows, total, err
	}
	return rows, total, nil
}

// Select 查询集合
func (r SysRole) Select(sysRole model.SysRole) ([]model.SysRole, error) {
	tx := db.DB("").Model(&model.SysRole{})
	tx = tx.Where("role_del = '0'")
	// 查询条件拼接
	if sysRole.RoleKey != "" {
		tx = tx.Where("role_key like ?", sysRole.RoleKey+"%")
	}
	if sysRole.RoleName != "" {
		tx = tx.Where("role_name like ?", sysRole.RoleName+"%")
	}
	if sysRole.RoleStatus != "" {
		tx = tx.Where("role_status = ?", sysRole.RoleStatus)
	}

	// 查询数据
	rows := []model.SysRole{}
	if err := tx.Order("role_sort asc").Find(&rows).Error; err != nil {
		return rows, err
	}
	return rows, nil
}

// SelectByIds 通过ID查询信息
func (r SysRole) SelectByIds(roleIds []int64) ([]model.SysRole, error) {
	rows := []model.SysRole{}
	if len(roleIds) <= 0 {
		return rows, nil
	}
	tx := db.DB("").Model(&model.SysRole{})
	// 构建查询条件
	tx = tx.Where("id in ? and role_del = '0'", roleIds)
	// 查询数据
	if err := tx.Find(&rows).Error; err != nil {
		return rows, err
	}
	return rows, nil
}

// Insert 新增信息 返回新增数据ID
func (r SysRole) Insert(d *gorm.DB, sysRole model.SysRole) (int64, error) {
	sysRole.RoleDel = "0"
	if d == nil {
		d = db.DB("")
	}
	// 执行插入
	if err := d.Create(&sysRole).Error; err != nil {
		return 0, err
	}
	return sysRole.ID, nil
}

// Update 修改信息 返回受影响行数
func (r SysRole) Update(d *gorm.DB, sysRole model.SysRole) (int64, error) {
	if sysRole.ID <= 0 {
		return 0, nil
	}
	if d == nil {
		d = db.DB("")
	}
	tx := d.Model(&model.SysRole{})
	// 构建查询条件
	tx = tx.Where("id = ?", sysRole.ID)
	tx = tx.Omit("id", "create_by", "create_time", "role_del")
	// 执行更新
	if err := tx.Updates(sysRole).Error; err != nil {
		return 0, err
	}
	return tx.RowsAffected, nil
}

// DeleteByIds 批量删除信息 返回受影响行数（软删除）
func (r SysRole) DeleteByIds(d *gorm.DB, roleIds []int64) (int64, error) {
	if len(roleIds) <= 0 {
		return 0, nil
	}
	if d == nil {
		d = db.DB("")
	}
	tx := d.Model(&model.SysRole{})
	// 构建查询条件
	tx = tx.Where("id in ?", roleIds)
	// 执行更新删除标记（软删除）
	if err := tx.Update("role_del", "1").Error; err != nil {
		return 0, err
	}
	return tx.RowsAffected, nil
}

// CheckUnique 检查信息是否唯一
func (r SysRole) CheckUnique(sysRole model.SysRole) (int64, error) {
	tx := db.DB("").Model(&model.SysRole{})
	tx = tx.Where("role_del = '0'")
	// 查询条件拼接
	if sysRole.RoleName != "" {
		tx = tx.Where("role_name = ?", sysRole.RoleName)
	}
	if sysRole.RoleKey != "" {
		tx = tx.Where("role_key = ?", sysRole.RoleKey)
	}

	// 查询数据
	var id int64 = 0
	if err := tx.Select("id").Limit(1).Find(&id).Error; err != nil {
		return id, err
	}
	return id, nil
}

// SelectByUserId 根据用户ID获取角色信息
func (r SysRole) SelectByUserId(userId int64) ([]model.SysRole, error) {
	rows := []model.SysRole{}
	if userId <= 0 {
		return rows, nil
	}

	tx := db.DB("").Model(&model.SysRole{}).
		Table("sys_role r").
		Distinct("r.*").
		Joins("LEFT JOIN sys_user_role ur ON ur.role_id = r.id").
		Joins("LEFT JOIN sys_user u ON u.id = ur.user_id").
		Where("r.role_del = '0' AND u.user_del = '0' AND ur.user_id = ?", userId)

	// 查询数据
	if err := tx.Find(&rows).Error; err != nil {
		return rows, err
	}
	return rows, nil
}

// SelectByUserIds 批量根据用户ID获取角色信息，返回 userId → []SysRole 的映射
func (r SysRole) SelectByUserIds(userIds []int64) (map[int64][]model.SysRole, error) {
	result := make(map[int64][]model.SysRole)
	if len(userIds) <= 0 {
		return result, nil
	}

	type userRoleRow struct {
		UserId int64 `gorm:"column:user_id"`
		model.SysRole
	}
	var rows []userRoleRow
	tx := db.DB("").Table("sys_user_role ur").
		Select("ur.user_id, r.*").
		Joins("JOIN sys_role r ON r.id = ur.role_id AND r.role_del = '0'").
		Where("ur.user_id IN ?", userIds)
	if err := tx.Find(&rows).Error; err != nil {
		return result, err
	}
	for _, row := range rows {
		result[row.UserId] = append(result[row.UserId], row.SysRole)
	}
	return result, nil
}
