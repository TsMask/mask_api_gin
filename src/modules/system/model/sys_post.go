package model

// SysPost 岗位对象 sys_post
type SysPost struct {
	PostID     string `json:"postId"`                      // 岗位ID
	PostCode   string `json:"postCode" binding:"required"` // 岗位编码
	PostName   string `json:"postName" binding:"required"` // 岗位名称
	PostSort   int    `json:"postSort"`                    // 显示顺序
	Status     string `json:"status"`                      // 状态（0停用 1正常）
	CreateBy   string `json:"createBy"`                    // 创建者
	CreateTime int64  `json:"createTime"`                  // 创建时间
	UpdateBy   string `json:"updateBy"`                    // 更新者
	UpdateTime int64  `json:"updateTime"`                  // 更新时间
	Remark     string `json:"remark"`                      // 备注
}
