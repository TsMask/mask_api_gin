package model

// SysPost 岗位信息表（用户 M:N 岗位，sys_user_post）
type SysPost struct {
	ID         int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`        // 岗位ID
	CreateBy   string `json:"createBy" gorm:"column:create_by"`                    // 创建者
	CreateTime int64  `json:"createTime" gorm:"column:create_time"`                // 创建时间
	UpdateBy   string `json:"updateBy" gorm:"column:update_by"`                    // 更新者
	UpdateTime int64  `json:"updateTime" gorm:"column:update_time"`                // 更新时间
	Remark     string `json:"remark" gorm:"column:remark"`                         // 备注
	PostCode   string `json:"postCode" gorm:"column:post_code" binding:"required"` // 岗位编码
	PostName   string `json:"postName" gorm:"column:post_name" binding:"required"` // 岗位名称
	PostSort   int64  `json:"postSort" gorm:"column:post_sort"`                    // 显示顺序
	PostStatus string `json:"postStatus" gorm:"column:post_status"`                // 状态（0停用 1正常）
}

// TableName 表名称
func (*SysPost) TableName() string {
	return "sys_post"
}
