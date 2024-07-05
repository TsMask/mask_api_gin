package model

// SysNotice 通知公告对象 sys_notice
type SysNotice struct {
	NoticeID      string `json:"noticeId"`                         // 公告ID
	NoticeTitle   string `json:"noticeTitle" binding:"required"`   // 公告标题
	NoticeType    string `json:"noticeType" binding:"required"`    // 公告类型（1通知 2公告）
	NoticeContent string `json:"noticeContent" binding:"required"` // 公告内容
	Status        string `json:"status"`                           // 公告状态（0关闭 1正常）
	DelFlag       string `json:"delFlag"`                          // 删除标志（0存在 1删除）
	CreateBy      string `json:"createBy"`                         // 创建者
	CreateTime    int64  `json:"createTime"`                       // 创建时间
	UpdateBy      string `json:"updateBy"`                         // 更新者
	UpdateTime    int64  `json:"updateTime"`                       // 更新时间
	Remark        string `json:"remark"`                           // 备注
}
