package repository

import (
	"mask_api_gin/src/framework/datasource"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/utils/repo"
	"mask_api_gin/src/modules/system/model"
	"strings"
)

// SysPostImpl 岗位表 数据层处理
var SysPostImpl = &sysPostImpl{
	selectSql: `select 
	post_id, post_code, post_name, post_sort, status, create_by, create_time, remark 
	from sys_post`,

	resultMap: map[string]string{
		"post_id":     "PostID",
		"post_code":   "PostCode",
		"post_name":   "PostName",
		"post_sort":   "PostSort",
		"status":      "Status",
		"create_by":   "CreateBy",
		"create_time": "CreateTime",
		"update_by":   "UpdateBy",
		"update_time": "UpdateTime",
		"remark":      "Remark",
	},
}

type sysPostImpl struct {
	// 查询视图对象SQL
	selectSql string
	// 结果字段与实体映射
	resultMap map[string]string
}

// convertResultRows 将结果记录转实体结果组
func (r *sysPostImpl) convertResultRows(rows []map[string]interface{}) []model.SysPost {
	arr := make([]model.SysPost, 0)
	for _, row := range rows {
		sysPost := model.SysPost{}
		for key, value := range row {
			if keyMapper, ok := r.resultMap[key]; ok {
				repo.SetFieldValue(&sysPost, keyMapper, value)
			}
		}
		arr = append(arr, sysPost)
	}
	return arr
}

// SelectPostPage 查询岗位分页数据集合
func (r *sysPostImpl) SelectPostPage(query map[string]string) map[string]interface{} {
	return map[string]interface{}{}
}

// SelectPostList 查询岗位数据集合
func (r *sysPostImpl) SelectPostList(sysPost model.SysPost) []model.SysPost {
	// 查询条件拼接
	var conditions []string
	var params []interface{}
	if sysPost.PostCode != "" {
		conditions = append(conditions, "post_code like concat(?, '%')")
		params = append(params, sysPost.PostCode)
	}
	if sysPost.PostName != "" {
		conditions = append(conditions, "post_name like concat(?, '%')")
		params = append(params, sysPost.PostName)
	}
	if sysPost.Status != "" {
		conditions = append(conditions, "status = ?")
		params = append(params, sysPost.Status)
	}

	// 构建查询条件语句
	whereSql := ""
	if len(conditions) > 0 {
		whereSql += " where " + strings.Join(conditions, " and ")
	}

	// 查询数据
	orderSql := " order by post_sort"
	querySql := r.selectSql + whereSql + orderSql
	rows, err := datasource.RawDB("", querySql, params)
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysPost{}
	}
	return r.convertResultRows(rows)
}

// SelectPostById 通过岗位ID查询岗位信息
func (r *sysPostImpl) SelectPostById(postId string) model.SysPost {
	return model.SysPost{}
}

// SelectPostListByUserId 根据用户ID获取岗位选择框列表
func (r *sysPostImpl) SelectPostListByUserId(userId string) []model.SysPost {
	// 查询数据
	querySql := `select p.post_id, p.post_name, p.post_code 
	from sys_post p 
    left join sys_user_post up on up.post_id = p.post_id 
    left join sys_user u on u.user_id = up.user_id 
    where u.user_id = ? order by p.post_sort`
	rows, err := datasource.RawDB("", querySql, []interface{}{userId})
	if err != nil {
		logger.Errorf("query err => %v", err)
		return []model.SysPost{}
	}
	return r.convertResultRows(rows)
}

// SelectPostsByUserName 查询用户所属岗位组
func (r *sysPostImpl) SelectPostsByUserName(userName string) []model.SysPost {
	return []model.SysPost{}
}

// DeletePostByIds 批量删除岗位信息
func (r *sysPostImpl) DeletePostByIds(postIds []string) int {
	return 0
}

// UpdatePost 修改岗位信息
func (r *sysPostImpl) UpdatePost(sysPost model.SysPost) int {
	return 0
}

// InsertPost 新增岗位信息
func (r *sysPostImpl) InsertPost(sysPost model.SysPost) string {
	return ""
}

// CheckUniquePostName 校验岗位名称
func (r *sysPostImpl) CheckUniquePostName(postName string) string {
	return ""
}

// CheckUniquePostCode 校验岗位编码
func (r *sysPostImpl) CheckUniquePostCode(postCode string) string {
	return ""
}
