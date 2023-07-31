package repository

import (
	"fmt"
	"mask_api_gin/src/framework/datasource"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/framework/utils/repo"
	"mask_api_gin/src/modules/system/model"
	"strings"
)

// 实例化数据层 SysUserPostImpl 结构体
var NewSysUserPostImpl = &SysUserPostImpl{}

// SysUserPostImpl 用户与岗位关联表 数据层处理
type SysUserPostImpl struct{}

// CountUserPostByPostId 通过岗位ID查询岗位使用数量
func (r *SysUserPostImpl) CountUserPostByPostId(postId string) int64 {
	querySql := "select count(1) as total from sys_user_role where role_id = ?"
	results, err := datasource.RawDB("", querySql, []interface{}{postId})
	if err != nil {
		logger.Errorf("query err => %v", err)
		return 0
	}
	if len(results) > 0 {
		return parse.Number(results[0]["total"])
	}
	return 0
}

// DeleteUserPost 批量删除用户和岗位关联
func (r *SysUserPostImpl) DeleteUserPost(userIds []string) int64 {
	placeholder := repo.KeyPlaceholderByQuery(len(userIds))
	sql := "delete from sys_user_post where user_id in  (" + placeholder + ")"
	parameters := repo.ConvertIdsSlice(userIds)
	results, err := datasource.ExecDB("", sql, parameters)
	if err != nil {
		logger.Errorf("delete err => %v", err)
		return 0
	}
	return results
}

// BatchUserPost 批量新增用户岗位信息
func (r *SysUserPostImpl) BatchUserPost(sysUserPosts []model.SysUserPost) int64 {
	keyValues := make([]string, 0)
	for _, item := range sysUserPosts {
		keyValues = append(keyValues, fmt.Sprintf("(%s,%s)", item.UserID, item.PostID))
	}
	sql := "insert into sys_user_post(user_id, post_id) values " + strings.Join(keyValues, ",")
	results, err := datasource.ExecDB("", sql, nil)
	if err != nil {
		logger.Errorf("delete err => %v", err)
		return 0
	}
	return results
}
