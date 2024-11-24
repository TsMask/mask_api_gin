package model

// SysUserPost 用户与岗位关联表
type SysUserPost struct {
	UserId string `json:"userId" gorm:"column:user_id"` // 用户ID
	PostId string `json:"postId" gorm:"column:post_id"` // 岗位ID
}

// TableName 表名称
func (*SysUserPost) TableName() string {
	return "sys_user_post"
}
