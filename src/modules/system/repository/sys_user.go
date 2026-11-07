package repository

import (
	"fmt"
	"time"

	"mask_api_gin/src/internal/infra/database"
	"mask_api_gin/src/internal/shared/parse"
	"mask_api_gin/src/modules/system/model"

	"gorm.io/gorm"
)

// NewSysUser 实例化数据层
var NewSysUser = &SysUser{}

// SysUser 用户表 数据层处理
type SysUser struct{}

// sysUserSortMapping 用户排序字段映射
var sysUserSortMapping = map[string]string{
	"id":              "id",
	"create_time":     "create_time",
	"update_time":     "update_time",
	"user_name":       "user_name",
	"user_nick":       "user_nick",
	"user_phone":      "user_phone",
	"user_email":      "user_email",
	"user_status":     "user_status",
	"user_login_ip":   "user_login_ip",
	"user_login_time": "user_login_time",
	"dept_id":         "dept_id",
}

// SelectByPage 分页查询集合
func (r SysUser) SelectByPage(query map[string]string, dataScopeSQL string) ([]model.SysUser, int64, error) {
	tx := database.DB("").Model(&model.SysUser{})
	tx = tx.Where("deleted_at = 0")
	// 查询条件拼接
	if v, ok := query["id"]; ok && v != "" {
		tx = tx.Where("id = ?", v)
	}
	if v, ok := query["user_name"]; ok && v != "" {
		tx = tx.Where("user_name like ?", v+"%")
	}
	if v, ok := query["user_phone"]; ok && v != "" {
		tx = tx.Where("user_phone like ?", v+"%")
	}
	if v, ok := query["user_status"]; ok && v != "" {
		tx = tx.Where("user_status = ?", v)
	}
	if v, ok := query["login_time_begin"]; ok && v != "" {
		if len(v) == 10 {
			v = fmt.Sprintf("%s000", v)
			tx = tx.Where("user_login_time >= ?", v)
		} else if len(v) == 13 {
			tx = tx.Where("user_login_time >= ?", v)
		}
	}
	if v, ok := query["login_time_end"]; ok && v != "" {
		if len(v) == 10 {
			v = fmt.Sprintf("%s999", v)
			tx = tx.Where("user_login_time <= ?", v)
		} else if len(v) == 13 {
			tx = tx.Where("user_login_time <= ?", v)
		}
	}
	if v, ok := query["dept_id"]; ok && v != "" {
		tx = tx.Where(`(dept_id = ? or dept_id in (
		select descendant_id from sys_dept_closure where ancestor_id = ?
		))`, v, v)
	}
	if dataScopeSQL != "" {
		tx = tx.Where(dataScopeSQL)
	}

	// 查询结果
	var total int64 = 0
	rows := []model.SysUser{}

	// 查询数量为0直接返回
	if err := tx.Count(&total).Error; err != nil {
		return rows, total, err
	}
	if total <= 0 {
		return rows, total, nil
	}

	// 查询数据分页
	pageNum, pageSize := database.PageNumSize(query["page_num"], query["page_size"])
	tx = tx.Limit(pageSize).Offset(pageSize * pageNum)

	// 排序，默认按 id DESC
	orders := database.SortParse(query["sort"], sysUserSortMapping)
	tx = database.SortOrDefault(tx, orders, "id desc")

	// 查询数据
	if err := tx.Find(&rows).Error; err != nil {
		return rows, total, err
	}
	return rows, total, nil
}

// Select 查询集合
func (r SysUser) Select(sysUser model.SysUser) ([]model.SysUser, error) {
	tx := database.DB("").Model(&model.SysUser{})
	tx = tx.Where("deleted_at = 0")
	// 查询条件拼接
	if sysUser.UserName != "" {
		tx = tx.Where("user_name like ?", sysUser.UserName+"%")
	}
	if sysUser.UserPhone != "" {
		tx = tx.Where("user_phone like ?", sysUser.UserPhone+"%")
	}
	if sysUser.UserStatus != "" {
		tx = tx.Where("user_status = ?", sysUser.UserStatus)
	}
	if sysUser.ID > 0 {
		tx = tx.Where("id = ?", sysUser.ID)
	}

	// 查询数据
	rows := []model.SysUser{}
	if err := tx.Order("id desc").Find(&rows).Error; err != nil {
		return rows, err
	}
	return rows, nil
}

// SelectByIds 通过ID查询信息
func (r SysUser) SelectByIds(userIds []int64) ([]model.SysUser, error) {
	rows := []model.SysUser{}
	if len(userIds) <= 0 {
		return rows, nil
	}
	tx := database.DB("").Model(&model.SysUser{})
	// 构建查询条件
	tx = tx.Where("id in ? and deleted_at = 0", userIds)
	// 查询数据
	if err := tx.Find(&rows).Error; err != nil {
		return rows, err
	}
	return rows, nil
}

// Insert 新增信息
func (r SysUser) Insert(d *gorm.DB, sysUser model.SysUser) (int64, error) {
	if d == nil {
		d = database.DB("")
	}
	// 执行插入
	if err := d.Create(&sysUser).Error; err != nil {
		return 0, err
	}
	return sysUser.ID, nil
}

// Update 修改信息
func (r SysUser) Update(d *gorm.DB, sysUser model.SysUser) (int64, error) {
	if sysUser.ID <= 0 {
		return 0, nil
	}
	if d == nil {
		d = database.DB("")
	}
	tx := d.Model(&model.SysUser{})
	// 构建查询条件
	tx = tx.Where("id = ?", sysUser.ID)
	tx = tx.Omit("id", "create_by", "create_time", "deleted_at", "deleted_by")
	// 执行更新
	if err := tx.Updates(sysUser).Error; err != nil {
		return 0, err
	}
	return tx.RowsAffected, nil
}

// DeleteByIds 批量删除信息
func (r SysUser) DeleteByIds(d *gorm.DB, userIds []int64, deletedBy string) (int64, error) {
	if len(userIds) <= 0 {
		return 0, nil
	}
	if d == nil {
		d = database.DB("")
	}
	tx := d.Model(&model.SysUser{})
	// 构建查询条件
	tx = tx.Where("id in ?", userIds)
	updates := map[string]interface{}{
		"deleted_at": time.Now().UnixMilli(),
	}
	if deletedBy != "" {
		updates["deleted_by"] = deletedBy
	}
	if err := tx.Updates(updates).Error; err != nil {
		return 0, err
	}
	return tx.RowsAffected, nil
}

// CheckUnique 检查信息是否唯一
func (r SysUser) CheckUnique(sysUser model.SysUser) (int64, error) {
	tx := database.DB("").Model(&model.SysUser{})
	tx = tx.Where("deleted_at = 0")
	// 查询条件拼接
	if sysUser.UserName != "" {
		tx = tx.Where("user_name = ?", sysUser.UserName)
	}
	if sysUser.UserPhone != "" {
		tx = tx.Where("user_phone = ?", sysUser.UserPhone)
	}
	if sysUser.UserEmail != "" {
		tx = tx.Where("user_email = ?", sysUser.UserEmail)
	}

	// 查询数据
	var id int64 = 0
	if err := tx.Select("id").Limit(1).Find(&id).Error; err != nil {
		return id, err
	}
	return id, nil
}

// SelectByUserName 通过登录账号查询信息
func (r SysUser) SelectByUserName(userName string) (model.SysUser, error) {
	item := model.SysUser{}
	if userName == "" {
		return item, nil
	}
	tx := database.DB("").Model(&model.SysUser{})
	// 构建查询条件
	tx = tx.Where("user_name = ? and deleted_at = 0", userName)
	// 查询数据
	if err := tx.Limit(1).Find(&item).Error; err != nil {
		return item, err
	}
	return item, nil
}

// SelectByPhone 通过手机号查询登录用户信息。
func (r SysUser) SelectByPhone(phone string) (model.SysUser, error) {
	item := model.SysUser{}
	if phone == "" {
		return item, nil
	}

	tx := database.DB("").Model(&model.SysUser{})
	tx = tx.Where("user_phone = ? and deleted_at = 0", phone)
	if err := tx.Limit(1).Find(&item).Error; err != nil {
		return item, err
	}
	return item, nil
}

// SelectAuthUsersByPage 分页查询集合By分配用户角色
func (r SysUser) SelectAuthUsersByPage(query map[string]string, dataScopeSQL string) ([]model.SysUser, int64, error) {
	tx := database.DB("").Model(&model.SysUser{})
	tx = tx.Where("deleted_at = 0")
	// 查询条件拼接
	if v, ok := query["user_name"]; ok && v != "" {
		tx = tx.Where("user_name like ?", v+"%")
	}
	if v, ok := query["user_phone"]; ok && v != "" {
		tx = tx.Where("user_phone like ?", v+"%")
	}
	if v, ok := query["user_status"]; ok && v != "" {
		tx = tx.Where("user_status = ?", v)
	}

	// 分配角色的用户（sys_user 主键 id；sys_user_role 外键 user_id/role_id 不变）
	if roleId, ok := query["role_id"]; ok && roleId != "" {
		auth, ok := query["auth"]
		if ok && parse.Boolean(auth) {
			tx = tx.Where(`id in (
				select distinct u.id from sys_user u
				inner join sys_user_role ur on u.id = ur.user_id
				and ur.role_id = ?
			)`, roleId)
		} else {
			tx = tx.Where(`id not in (
				select distinct u.id from sys_user u
				inner join sys_user_role ur on u.id = ur.user_id
				and ur.role_id = ?
			)`, roleId)
		}
	}
	if dataScopeSQL != "" {
		tx = tx.Where(dataScopeSQL)
	}

	// 查询结果
	var total int64 = 0
	rows := []model.SysUser{}

	// 查询数量为0直接返回
	if err := tx.Count(&total).Error; err != nil {
		return rows, total, err
	}
	if total <= 0 {
		return rows, total, nil
	}

	// 查询数据分页
	pageNum, pageSize := database.PageNumSize(query["page_num"], query["page_size"])
	tx = tx.Limit(pageSize).Offset(pageSize * pageNum)

	// 排序，默认按 id DESC
	orders := database.SortParse(query["sort"], sysUserSortMapping)
	tx = database.SortOrDefault(tx, orders, "id desc")

	// 查询数据
	if err := tx.Find(&rows).Error; err != nil {
		return rows, total, err
	}
	return rows, total, nil
}
