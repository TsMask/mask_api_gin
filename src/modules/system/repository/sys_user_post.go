package repository

import (
	"mask_api_gin/src/framework/database/db"
	"mask_api_gin/src/modules/system/model"

	"gorm.io/gorm"
)

// NewSysUserPost 实例化数据层
var NewSysUserPost = &SysUserPost{}

// SysUserPost 用户与岗位关联表 数据层处理
type SysUserPost struct{}

// ExistUserByPostId 存在用户使用数量
func (r SysUserPost) ExistUserByPostId(postId int64) (int64, error) {
	if postId <= 0 {
		return 0, nil
	}
	tx := db.DB("").Model(&model.SysUserPost{})
	tx = tx.Where("post_id = ?", postId)
	// 查询数据
	var count int64 = 0
	if err := tx.Count(&count).Error; err != nil {
		return count, err
	}
	return count, nil
}

// DeleteByUserIds 批量删除关联By用户
func (r SysUserPost) DeleteByUserIds(d *gorm.DB, userIds []int64) (int64, error) {
	if len(userIds) <= 0 {
		return 0, nil
	}
	if d == nil {
		d = db.DB("")
	}
	tx := d.Where("user_id in ?", userIds)
	// 执行删除
	if err := tx.Delete(&model.SysUserPost{}).Error; err != nil {
		return 0, err
	}
	return tx.RowsAffected, nil
}

// BatchInsert 批量新增信息
func (r SysUserPost) BatchInsert(d *gorm.DB, sysUserPosts []model.SysUserPost) (int64, error) {
	if len(sysUserPosts) <= 0 {
		return 0, nil
	}
	if d == nil {
		d = db.DB("")
	}
	tx := d.CreateInBatches(sysUserPosts, 500)
	if err := tx.Error; err != nil {
		return 0, err
	}
	return tx.RowsAffected, nil
}
