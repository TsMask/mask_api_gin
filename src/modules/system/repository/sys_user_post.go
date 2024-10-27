package repository

import (
	"fmt"
	db "mask_api_gin/src/framework/data_source"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/modules/system/model"
	"strings"
)

// NewSysUserPost 实例化数据层
var NewSysUserPost = &SysUserPost{}

// SysUserPost 用户与岗位关联表 数据层处理
type SysUserPost struct{}

// ExistUserByPostId 存在用户使用数量
func (r SysUserPost) ExistUserByPostId(postId string) int64 {
	querySql := "select count(1) as total from sys_user_role where role_id = ?"
	results, err := db.RawDB("", querySql, []any{postId})
	if err != nil {
		logger.Errorf("query err => %v", err)
		return 0
	}
	return parse.Number(results[0]["total"])
}

// DeleteByUserIds 批量删除关联By用户
func (r SysUserPost) DeleteByUserIds(userIds []string) int64 {
	placeholder := db.KeyPlaceholderByQuery(len(userIds))
	sql := fmt.Sprintf("delete from sys_user_post where user_id in  (%s)", placeholder)
	parameters := db.ConvertIdsSlice(userIds)
	results, err := db.ExecDB("", sql, parameters)
	if err != nil {
		logger.Errorf("delete err => %v", err)
		return 0
	}
	return results
}

// BatchInsert 批量新增信息
func (r SysUserPost) BatchInsert(sysUserPosts []model.SysUserPost) int64 {
	up := make([]string, 0)
	for _, item := range sysUserPosts {
		up = append(up, fmt.Sprintf("(%s,%s)", item.UserId, item.PostId))
	}
	sql := fmt.Sprintf("insert into db(user_id, post_id) values %s", strings.Join(up, ","))
	results, err := db.ExecDB("", sql, nil)
	if err != nil {
		logger.Errorf("delete err => %v", err)
		return 0
	}
	return results
}
