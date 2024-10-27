package model

// SysNotice 通知公告表
type SysNotice struct {
	NoticeId      string `json:"notice_id" gorm:"column:notice_id;primary_key"`                  // 公告ID
	NoticeTitle   string `json:"notice_title" gorm:"column:notice_title" binding:"required"`     // 公告标题
	NoticeType    string `json:"notice_type" gorm:"column:notice_type" binding:"required"`       // 公告类型（1通知 2公告）
	NoticeContent string `json:"notice_content" gorm:"column:notice_content" binding:"required"` // 公告内容
	Status        string `json:"status" gorm:"column:status"`                                    // 公告状态（0关闭 1正常）
	DelFlag       string `json:"del_flag" gorm:"column:del_flag"`                                // 删除标志（0存在 1删除）
	CreateBy      string `json:"create_by" gorm:"column:create_by"`                              // 创建者
	CreateTime    int64  `json:"create_time" gorm:"column:create_time"`                          // 创建时间
	UpdateBy      string `json:"update_by" gorm:"column:update_by"`                              // 更新者
	UpdateTime    int64  `json:"update_time" gorm:"column:update_time"`                          // 更新时间
	Remark        string `json:"remark" gorm:"column:remark"`                                    // 备注
}

// TableName 表名称
func (*SysNotice) TableName() string {
	return "sys_notice"
}
