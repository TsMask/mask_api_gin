package model

// SysNotice 通知公告对象 sys_notice
type SysNotice struct {
	// 公告ID
	NoticeID int64 `json:"noticeId"`
	// 公告标题
	NoticeTitle string `json:"noticeTitle"`
	// 公告类型（1通知 2公告）
	NoticeType string `json:"noticeType"`
	// 公告内容
	NoticeContent string `json:"noticeContent"`
	// 公告状态（0关闭 1正常）
	Status string `json:"status"`
	// 删除标志（0代表存在 1代表删除）
	DelFlag string `json:"delFlag"`
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
