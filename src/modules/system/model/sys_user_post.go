package model

// SysUserPost 用户与岗位关联表
type SysUserPost struct {
	UserId int64 `json:"user_id" gorm:"column:user_id;primaryKey"` // 用户ID
	PostId int64 `json:"post_id" gorm:"column:post_id;primaryKey"` // 岗位ID
}

// TableName 表名称
func (*SysUserPost) TableName() string {
	return "sys_user_post"
}
