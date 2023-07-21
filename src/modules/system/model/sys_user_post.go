package model

// SysUserPost 用户和岗位关联对象 sys_user_post
type SysUserPost struct {
	UserID string `json:"userId"` // 用户ID
	PostID string `json:"postId"` // 岗位ID
}

// NewSysUserPost 创建用户和岗位关联对象的构造函数
func NewSysUserPost(userID string, postID string) SysUserPost {
	return SysUserPost{
		UserID: userID,
		PostID: postID,
	}
}
