package repository

import (
	"mask_api_gin/src/framework/database/db"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/utils/crypto"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/modules/system/model"

	"fmt"
	"time"
)

// NewSysUser 实例化数据层
var NewSysUser = &SysUser{}

// SysUser 用户表 数据层处理
type SysUser struct{}

// SelectByPage 分页查询集合
func (r SysUser) SelectByPage(query map[string]string, dataScopeWhereSQL string) ([]model.SysUser, int64) {
	tx := db.DB("").Model(&model.SysUser{})
	tx = tx.Where("del_flag = '0'")
	// 查询条件拼接
	if v, ok := query["userId"]; ok && v != "" {
		tx = tx.Where("user_id = ?", v)
	}
	if v, ok := query["userName"]; ok && v != "" {
		tx = tx.Where("user_name like concat(?, '%')", v)
	}
	if v, ok := query["phone"]; ok && v != "" {
		tx = tx.Where("phone like concat(?, '%')", v)
	}
	if v, ok := query["statusFlag"]; ok && v != "" {
		tx = tx.Where("status_flag = ?", v)
	}
	if v, ok := query["beginTime"]; ok && v != "" {
		if len(v) == 10 {
			v = fmt.Sprintf("%s000", v)
			tx = tx.Where("login_time >= ?", v)
		} else if len(v) == 13 {
			tx = tx.Where("login_time >= ?", v)
		}
	}
	if v, ok := query["endTime"]; ok && v != "" {
		if len(v) == 10 {
			v = fmt.Sprintf("%s999", v)
			tx = tx.Where("login_time <= ?", v)
		} else if len(v) == 13 {
			tx = tx.Where("login_time <= ?", v)
		}
	}
	if v, ok := query["deptId"]; ok && v != "" {
		tx = tx.Where(`(dept_id = ? or dept_id in ( 
		select t.dept_id from sys_dept t where find_in_set(?, ancestors) 
		))`, v, v)
	}
	if dataScopeWhereSQL != "" {
		tx = tx.Where(dataScopeWhereSQL)
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
	err := tx.Find(&rows).Error
	if err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return rows, total
	}
	return rows, total
}

// Select 查询集合
func (r SysUser) Select(sysUser model.SysUser) []model.SysUser {
	tx := db.DB("").Model(&model.SysUser{})
	tx = tx.Where("del_flag = '0'")
	// 查询条件拼接
	if sysUser.UserName != "" {
		tx = tx.Where("user_name like concat(?, '%')", sysUser.UserName)
	}
	if sysUser.Phone != "" {
		tx = tx.Where("phone like concat(?, '%')", sysUser.Phone)
	}
	if sysUser.StatusFlag != "" {
		tx = tx.Where("status_flag = ?", sysUser.StatusFlag)
	}
	if sysUser.UserId != "" {
		tx = tx.Where("user_id = ?", sysUser.UserId)
	}

	// 查询数据
	rows := []model.SysUser{}
	if err := tx.Order("login_time desc").Find(&rows).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return rows
	}
	return rows
}

// SelectByIds 通过ID查询信息
func (r SysUser) SelectByIds(userIds []string) []model.SysUser {
	rows := []model.SysUser{}
	if len(userIds) <= 0 {
		return rows
	}
	tx := db.DB("").Model(&model.SysUser{})
	// 构建查询条件
	tx = tx.Where("user_id in ? and del_flag = '0'", userIds)
	// 查询数据
	if err := tx.Find(&rows).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return rows
	}
	return rows
}

// Insert 新增信息
func (r SysUser) Insert(sysUser model.SysUser) string {
	sysUser.DelFlag = "0"
	if sysUser.Password != "" {
		sysUser.Password = crypto.BcryptHash(sysUser.Password)
	}
	if sysUser.CreateBy != "" {
		ms := time.Now().UnixMilli()
		sysUser.UpdateBy = sysUser.CreateBy
		sysUser.UpdateTime = ms
		sysUser.CreateTime = ms
	}
	// 执行插入
	if err := db.DB("").Create(&sysUser).Error; err != nil {
		logger.Errorf("insert err => %v", err.Error())
		return ""
	}
	return sysUser.UserId
}

// Update 修改信息
func (r SysUser) Update(sysUser model.SysUser) int64 {
	if sysUser.UserId == "" {
		return 0
	}
	if sysUser.Password != "" {
		sysUser.Password = crypto.BcryptHash(sysUser.Password)
	}
	if sysUser.UpdateBy != "" {
		sysUser.UpdateTime = time.Now().UnixMilli()
	}
	tx := db.DB("").Model(&model.SysUser{})
	// 构建查询条件
	tx = tx.Where("user_id = ?", sysUser.UserId)
	tx = tx.Omit("user_id", "del_flag", "create_by", "create_time")
	// 执行更新
	if err := tx.Updates(sysUser).Error; err != nil {
		logger.Errorf("update err => %v", err.Error())
		return 0
	}
	return tx.RowsAffected
}

// DeleteByIds 批量删除信息
func (r SysUser) DeleteByIds(userIds []string) int64 {
	if len(userIds) <= 0 {
		return 0
	}
	tx := db.DB("").Model(&model.SysUser{})
	// 构建查询条件
	tx = tx.Where("user_id in ?", userIds)
	// 执行更新删除标记
	if err := tx.Update("del_flag", "1").Error; err != nil {
		logger.Errorf("update err => %v", err.Error())
		return 0
	}
	return tx.RowsAffected
}

// CheckUnique 检查信息是否唯一
func (r SysUser) CheckUnique(sysUser model.SysUser) string {
	tx := db.DB("").Model(&model.SysUser{})
	tx = tx.Where("del_flag = 0")
	// 查询条件拼接
	if sysUser.UserName != "" {
		tx = tx.Where("user_name = ?", sysUser.UserName)
	}
	if sysUser.Phone != "" {
		tx = tx.Where("phone = ?", sysUser.Phone)
	}
	if sysUser.Email != "" {
		tx = tx.Where("email = ?", sysUser.Email)
	}

	// 查询数据
	var id string = ""
	if err := tx.Select("user_id").Limit(1).Find(&id).Error; err != nil {
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
	tx = tx.Where("user_name = ? and del_flag = '0'", userName)
	// 查询数据
	if err := tx.Limit(1).Find(&item).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return item
	}
	return item
}

// SelectAuthUsersByPage 分页查询集合By分配用户角色
func (r SysUser) SelectAuthUsersByPage(query map[string]string, dataScopeWhereSQL string) ([]model.SysUser, int64) {
	tx := db.DB("").Model(&model.SysUser{})
	tx = tx.Where("del_flag = '0'")
	// 查询条件拼接
	if v, ok := query["userName"]; ok && v != "" {
		tx = tx.Where("user_name like concat(?, '%')", v)
	}
	if v, ok := query["phone"]; ok && v != "" {
		tx = tx.Where("phone like concat(?, '%')", v)
	}
	if v, ok := query["statusFlag"]; ok && v != "" {
		tx = tx.Where("status_flag = ?", v)
	}
	if dataScopeWhereSQL != "" {
		tx = tx.Where(dataScopeWhereSQL)
	}

	// 分配角色的用户
	if roleId, ok := query["roleId"]; ok && roleId != "" {
		auth, ok := query["auth"]
		if ok && parse.Boolean(auth) {
			tx = tx.Where(`user_id in (
				select distinct u.user_id from sys_user u 
				inner join sys_user_role ur on u.user_id = ur.user_id 
				and ur.role_id = ?
			)`, roleId)
		} else {
			tx = tx.Where(`user_id not in (
				select distinct u.user_id from sys_user u 
				inner join sys_user_role ur on u.user_id = ur.user_id 
				and ur.role_id = ?
			)`, roleId)
		}
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
	err := tx.Find(&rows).Error
	if err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return rows, total
	}
	return rows, total
}
