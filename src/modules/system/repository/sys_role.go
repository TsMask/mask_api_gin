package repository

import (
	"mask_api_gin/src/framework/database/db"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/modules/system/model"

	"fmt"
	"time"
)

// NewSysRole 实例化数据层
var NewSysRole = &SysRole{}

// SysRole 角色表 数据层处理
type SysRole struct{}

// SelectByPage 分页查询集合
func (r SysRole) SelectByPage(query map[string]string, dataScopeSQL string) ([]model.SysRole, int64) {
	tx := db.DB("").Model(&model.SysRole{})
	tx = tx.Where("del_flag = '0'")
	// 查询条件拼接
	if v, ok := query["roleName"]; ok && v != "" {
		tx = tx.Where("role_name like concat(?, '%')", v)
	}
	if v, ok := query["roleKey"]; ok && v != "" {
		tx = tx.Where("role_key like concat(?, '%')", v)
	}
	if v, ok := query["statusFlag"]; ok && v != "" {
		tx = tx.Where("status_flag = ?", v)
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
			v = fmt.Sprintf("%s000", v)
			tx = tx.Where("create_time <= ?", v)
		} else if len(v) == 13 {
			tx = tx.Where("create_time <= ?", v)
		}
	}

	if v, ok := query["deptId"]; ok && v != "" {
		tx = tx.Where(`(dept_id = ? or dept_id in ( 
		select t.dept_id from sys_dept t where find_in_set(?, ancestors)
		))`, v, v)
	}

	// 查询结果
	var total int64 = 0
	rows := []model.SysRole{}

	// 查询数量为0直接返回
	if err := tx.Count(&total).Error; err != nil || total <= 0 {
		return rows, total
	}

	// 查询数据分页
	pageNum, pageSize := db.PageNumSize(query["pageNum"], query["pageSize"])
	tx = tx.Limit(pageSize).Offset(pageSize * pageNum)
	err := tx.Order("role_sort asc").Find(&rows).Error
	if err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return rows, total
	}
	return rows, total
}

// Select 查询集合
func (r SysRole) Select(sysRole model.SysRole, dataScopeWhereSQL string) []model.SysRole {
	tx := db.DB("").Model(&model.SysRole{})
	tx = tx.Where("del_flag = '0'")
	// 查询条件拼接
	if sysRole.RoleKey != "" {
		tx = tx.Where("role_key like concat(?, '%')", sysRole.RoleKey)
	}
	if sysRole.RoleName != "" {
		tx = tx.Where("role_name like concat(?, '%')", sysRole.RoleName)
	}
	if sysRole.StatusFlag != "" {
		tx = tx.Where("status_flag = ?", sysRole.StatusFlag)
	}
	if dataScopeWhereSQL != "" {
		tx = tx.Where(dataScopeWhereSQL)
	}

	// 查询数据
	rows := []model.SysRole{}
	if err := tx.Order("role_sort asc").Find(&rows).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return rows
	}
	return rows
}

// SelectByIds 通过ID查询信息
func (r SysRole) SelectByIds(roleIds []string) []model.SysRole {
	rows := []model.SysRole{}
	if len(roleIds) <= 0 {
		return rows
	}
	tx := db.DB("").Model(&model.SysRole{})
	// 构建查询条件
	tx = tx.Where("role_id in ? and del_flag = '0'", roleIds)
	// 查询数据
	if err := tx.Find(&rows).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return rows
	}
	return rows
}

// Insert 新增信息 返回新增数据ID
func (r SysRole) Insert(sysRole model.SysRole) string {
	sysRole.DelFlag = "0"
	if sysRole.CreateBy != "" {
		ms := time.Now().UnixMilli()
		sysRole.UpdateBy = sysRole.CreateBy
		sysRole.UpdateTime = ms
		sysRole.CreateTime = ms
	}
	// 执行插入
	if err := db.DB("").Create(&sysRole).Error; err != nil {
		logger.Errorf("insert err => %v", err.Error())
		return ""
	}
	return sysRole.RoleId
}

// Update 修改信息 返回受影响行数
func (r SysRole) Update(sysRole model.SysRole) int64 {
	if sysRole.RoleId == "" {
		return 0
	}
	if sysRole.UpdateBy != "" {
		sysRole.UpdateTime = time.Now().UnixMilli()
	}
	tx := db.DB("").Model(&model.SysRole{})
	// 构建查询条件
	tx = tx.Where("role_id = ?", sysRole.RoleId)
	tx = tx.Omit("role_id", "del_flag", "create_by", "create_time")
	// 执行更新
	if err := tx.Updates(sysRole).Error; err != nil {
		logger.Errorf("update err => %v", err.Error())
		return 0
	}
	return tx.RowsAffected
}

// DeleteByIds 批量删除信息 返回受影响行数
func (r SysRole) DeleteByIds(roleIds []string) int64 {
	if len(roleIds) <= 0 {
		return 0
	}
	tx := db.DB("").Model(&model.SysRole{})
	// 构建查询条件
	tx = tx.Where("role_id in ?", roleIds)
	// 执行更新删除标记
	if err := tx.Update("del_flag", "1").Error; err != nil {
		logger.Errorf("update err => %v", err.Error())
		return 0
	}
	return tx.RowsAffected
}

// SelectByUserId 根据用户ID获取角色信息
func (r SysRole) SelectByUserId(userId string) []model.SysRole {
	rows := []model.SysRole{}
	if userId == "" {
		return rows
	}
	tx := db.DB("").Table("sys_user_role ur")
	// 构建查询条件
	tx = tx.Distinct("r.role_id, r.role_name, r.role_key, r.role_sort, r.data_scope, r.menu_check_strictly, r.dept_check_strictly, r.status_flag, r.del_flag, r.create_time, r.remark").
		Joins("left join sys_user u on u.user_id = ur.user_id").
		Joins("left join sys_role r on r.role_id = ur.role_id").
		Where("u.del_flag = '0' AND ur.user_id = ?", userId)

	// 查询数据
	if err := tx.Find(&rows).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return rows
	}
	return rows
}

// CheckUnique 检查信息是否唯一
func (r SysRole) CheckUnique(sysRole model.SysRole) string {
	tx := db.DB("").Model(&model.SysRole{})
	tx = tx.Where("del_flag = '0'")
	// 查询条件拼接
	if sysRole.RoleName != "" {
		tx = tx.Where("role_name = ?", sysRole.RoleName)
	}
	if sysRole.RoleKey != "" {
		tx = tx.Where("role_key = ?", sysRole.RoleKey)
	}

	// 查询数据
	var id string = ""
	if err := tx.Select("role_id").Limit(1).Find(&id).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return id
	}
	return id
}
