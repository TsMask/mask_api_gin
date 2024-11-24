package repository

import (
	"mask_api_gin/src/framework/database/db"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/modules/system/model"
)

// NewSysUserPost 实例化数据层
var NewSysUserPost = &SysUserPost{}

// SysUserPost 用户与岗位关联表 数据层处理
type SysUserPost struct{}

// ExistUserByPostId 存在用户使用数量
func (r SysUserPost) ExistUserByPostId(postId string) int64 {
	if postId == "" {
		return 0
	}
	tx := db.DB("").Model(&model.SysUserPost{})
	tx = tx.Where("post_id = ?", postId)
	// 查询数据
	var count int64 = 0
	if err := tx.Count(&count).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return count
	}
	return count
}

// DeleteByUserIds 批量删除关联By用户
func (r SysUserPost) DeleteByUserIds(userIds []string) int64 {
	if len(userIds) <= 0 {
		return 0
	}
	tx := db.DB("").Where("user_id in ?", userIds)
	// 执行删除
	if err := tx.Delete(&model.SysUserPost{}).Error; err != nil {
		logger.Errorf("delete err => %v", err.Error())
		return 0
	}
	return tx.RowsAffected
}

// BatchInsert 批量新增信息
func (r SysUserPost) BatchInsert(sysUserPosts []model.SysUserPost) int64 {
	if len(sysUserPosts) <= 0 {
		return 0
	}
	// 执行批量删除
	tx := db.DB("").CreateInBatches(sysUserPosts, 500)
	if err := tx.Error; err != nil {
		logger.Errorf("delete batch err => %v", err.Error())
		return 0
	}
	return tx.RowsAffected
}
