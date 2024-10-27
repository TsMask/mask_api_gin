package model

// SysPost 岗位信息表
type SysPost struct {
	PostId     string `json:"post_id" gorm:"column:post_id;primary_key"`            // 岗位ID
	PostCode   string `json:"post_code" gorm:"column:post_code" binding:"required"` // 岗位编码
	PostName   string `json:"post_name" gorm:"column:post_name" binding:"required"` // 岗位名称
	PostSort   int64  `json:"post_sort" gorm:"column:post_sort"`                    // 显示顺序
	Status     string `json:"status" gorm:"column:status"`                          // 状态（0停用 1正常）
	CreateBy   string `json:"create_by" gorm:"column:create_by"`                    // 创建者
	CreateTime int64  `json:"create_time" gorm:"column:create_time"`                // 创建时间
	UpdateBy   string `json:"update_by" gorm:"column:update_by"`                    // 更新者
	UpdateTime int64  `json:"update_time" gorm:"column:update_time"`                // 更新时间
	Remark     string `json:"remark" gorm:"column:remark"`                          // 备注
}

// TableName 表名称
func (*SysPost) TableName() string {
	return "sys_post"
}
