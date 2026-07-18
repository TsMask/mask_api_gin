package repository

import (
	"fmt"

	"mask_api_gin/src/framework/database/db"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/modules/system/model"

	"gorm.io/gorm"
)

// NewSysUser 实例化数据层
var NewSysUser = &SysUser{}

// SysUser 用户表 数据层处理
type SysUser struct{}

// sysUserSortMapping 用户排序字段映射
var sysUserSortMapping = map[string]string{
	"id":            "id",
	"createTime":    "create_time",
	"updateTime":    "update_time",
	"userName":      "user_name",
	"userNick":      "user_nick",
	"userPhone":     "user_phone",
	"userEmail":     "user_email",
	"userStatus":    "user_status",
	"userLoginIp":   "user_login_ip",
	"userLoginTime": "user_login_time",
	"deptId":        "dept_id",
}

// SelectByPage 分页查询集合
func (r SysUser) SelectByPage(query map[string]string, dataScopeSQL string) ([]model.SysUser, int64) {
	tx := db.DB("").Model(&model.SysUser{})
	tx = tx.Where("user_del = '0'")
	// 查询条件拼接
	if v, ok := query["id"]; ok && v != "" {
		tx = tx.Where("id = ?", v)
	}
	if v, ok := query["userName"]; ok && v != "" {
		tx = tx.Where("user_name like ?", v+"%")
	}
	if v, ok := query["userPhone"]; ok && v != "" {
		tx = tx.Where("user_phone like ?", v+"%")
	}
	if v, ok := query["userStatus"]; ok && v != "" {
		tx = tx.Where("user_status = ?", v)
	}
	if v, ok := query["beginTime"]; ok && v != "" {
		if len(v) == 10 {
			v = fmt.Sprintf("%s000", v)
			tx = tx.Where("user_login_time >= ?", v)
		} else if len(v) == 13 {
			tx = tx.Where("user_login_time >= ?", v)
		}
	}
	if v, ok := query["endTime"]; ok && v != "" {
		if len(v) == 10 {
			v = fmt.Sprintf("%s999", v)
			tx = tx.Where("user_login_time <= ?", v)
		} else if len(v) == 13 {
			tx = tx.Where("user_login_time <= ?", v)
		}
	}
	if v, ok := query["deptId"]; ok && v != "" {
		tx = tx.Where(`(dept_id = ? or dept_id in (
		select t.id from sys_dept t where find_in_set(?, ancestors)
		))`, v, v)
	}
	if dataScopeSQL != "" {
		tx = tx.Where(dataScopeSQL)
	}

	// 查询结果
	var total int64 = 0
	rows := []model.SysUser{}

	// 查询数量为0直接返回
	if err := tx.Count(&total).Error; err != nil || total <= 0 {
		return rows, total
	}

	// 查询数据分页
	pageNum, pageSize := db.PageNumSize(query["pageNum"], query["pageSize"])
	tx = tx.Limit(pageSize).Offset(pageSize * pageNum)

	// 排序，默认按 id DESC
	orders := db.SortParse(query["sort"], sysUserSortMapping)
	tx = db.SortOrDefault(tx, orders, "id desc")

	// 查询数据
	if err := tx.Find(&rows).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return rows, total
	}
	return rows, total
}

// Select 查询集合
func (r SysUser) Select(sysUser model.SysUser) []model.SysUser {
	tx := db.DB("").Model(&model.SysUser{})
	tx = tx.Where("user_del = '0'")
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
		logger.Errorf("query find err => %v", err.Error())
		return rows
	}
	return rows
}

// SelectByIds 通过ID查询信息
func (r SysUser) SelectByIds(userIds []int64) []model.SysUser {
	rows := []model.SysUser{}
	if len(userIds) <= 0 {
		return rows
	}
	tx := db.DB("").Model(&model.SysUser{})
	// 构建查询条件
	tx = tx.Where("id in ? and user_del = '0'", userIds)
	// 查询数据
	if err := tx.Find(&rows).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return rows
	}
	return rows
}

// Insert 新增信息
func (r SysUser) Insert(d *gorm.DB, sysUser model.SysUser) int64 {
	sysUser.UserDel = "0"
	if d == nil {
		d = db.DB("")
	}
	// 执行插入
	if err := d.Create(&sysUser).Error; err != nil {
		logger.Errorf("insert err => %v", err.Error())
		return 0
	}
	return sysUser.ID
}

// Update 修改信息
func (r SysUser) Update(d *gorm.DB, sysUser model.SysUser) int64 {
	if sysUser.ID <= 0 {
		return 0
	}
	if d == nil {
		d = db.DB("")
	}
	tx := d.Model(&model.SysUser{})
	// 构建查询条件
	tx = tx.Where("id = ?", sysUser.ID)
	tx = tx.Omit("id", "create_by", "create_time", "user_del")
	// 执行更新
	if err := tx.Updates(sysUser).Error; err != nil {
		logger.Errorf("update err => %v", err.Error())
		return 0
	}
	return tx.RowsAffected
}

// DeleteByIds 批量删除信息
func (r SysUser) DeleteByIds(d *gorm.DB, userIds []int64) int64 {
	if len(userIds) <= 0 {
		return 0
	}
	if d == nil {
		d = db.DB("")
	}
	tx := d.Model(&model.SysUser{})
	// 构建查询条件
	tx = tx.Where("id in ?", userIds)
	// 执行更新删除标记（软删除）
	if err := tx.Update("user_del", "1").Error; err != nil {
		logger.Errorf("update err => %v", err.Error())
		return 0
	}
	return tx.RowsAffected
}

// CheckUnique 检查信息是否唯一
func (r SysUser) CheckUnique(sysUser model.SysUser) int64 {
	tx := db.DB("").Model(&model.SysUser{})
	tx = tx.Where("user_del = '0'")
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
		logger.Errorf("query find err => %v", err.Error())
		return id
	}
	return id
}

// SelectByUserName 通过登录账号查询信息
func (r SysUser) SelectByUserName(userName string) model.SysUser {
	item := model.SysUser{}
	if userName == "" {
		return item
	}
	tx := db.DB("").Model(&model.SysUser{})
	// 构建查询条件
	tx = tx.Where("user_name = ? and user_del = '0'", userName)
	// 查询数据
	if err := tx.Limit(1).Find(&item).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return item
	}
	return item
}

// SelectAuthUsersByPage 分页查询集合By分配用户角色
func (r SysUser) SelectAuthUsersByPage(query map[string]string, dataScopeSQL string) ([]model.SysUser, int64) {
	tx := db.DB("").Model(&model.SysUser{})
	tx = tx.Where("user_del = '0'")
	// 查询条件拼接
	if v, ok := query["userName"]; ok && v != "" {
		tx = tx.Where("user_name like ?", v+"%")
	}
	if v, ok := query["userPhone"]; ok && v != "" {
		tx = tx.Where("user_phone like ?", v+"%")
	}
	if v, ok := query["userStatus"]; ok && v != "" {
		tx = tx.Where("user_status = ?", v)
	}

	// 分配角色的用户（sys_user 主键 id；sys_user_role 外键 user_id/role_id 不变）
	if roleId, ok := query["roleId"]; ok && roleId != "" {
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
	if err := tx.Count(&total).Error; err != nil || total <= 0 {
		return rows, total
	}

	// 查询数据分页
	pageNum, pageSize := db.PageNumSize(query["pageNum"], query["pageSize"])
	tx = tx.Limit(pageSize).Offset(pageSize * pageNum)

	// 排序，默认按 id DESC
	orders := db.SortParse(query["sort"], sysUserSortMapping)
	tx = db.SortOrDefault(tx, orders, "id desc")

	// 查询数据
	if err := tx.Find(&rows).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return rows, total
	}
	return rows, total
}
