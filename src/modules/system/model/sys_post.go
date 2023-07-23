package model

// SysPost 岗位对象 sys_post
type SysPost struct {
	// 岗位ID
	PostID string `json:"postId"`
	// 岗位编码
	PostCode string `json:"postCode" binding:"required"`
	// 岗位名称
	PostName string `json:"postName" binding:"required"`
	// 显示顺序
	PostSort int `json:"postSort"`
	// 状态（0停用 1正常）
	Status string `json:"status"`
	// 创建者
	CreateBy string `json:"createBy"`
	// 创建时间
	CreateTime int64 `json:"createTime"`
	// 更新者
	UpdateBy string `json:"updateBy"`
	// 更新时间
	UpdateTime int64 `json:"updateTime"`
	// 备注
	Remark string `json:"remark"`
}
