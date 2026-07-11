package model

// SysNotice 通知公告表
type SysNotice struct {
	ID            int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`                  // 公告ID
	CreateBy      string `json:"createBy" gorm:"column:create_by"`                              // 创建者
	CreateTime    int64  `json:"createTime" gorm:"column:create_time"`                          // 创建时间
	UpdateBy      string `json:"updateBy" gorm:"column:update_by"`                              // 更新者
	UpdateTime    int64  `json:"updateTime" gorm:"column:update_time"`                          // 更新时间
	Remark        string `json:"remark" gorm:"column:remark"`                                   // 备注
	NoticeTitle   string `json:"noticeTitle" gorm:"column:notice_title" binding:"required"`     // 公告标题
	NoticeType    string `json:"noticeType" gorm:"column:notice_type" binding:"required"`       // 公告类型（1通知 2公告）
	NoticeContent string `json:"noticeContent" gorm:"column:notice_content" binding:"required"` // 公告内容
	NoticeStatus  string `json:"noticeStatus" gorm:"column:notice_status"`                      // 公告状态（0关闭 1正常）
	NoticeDel     string `json:"-" gorm:"column:notice_del"`                                    // 删除标记（0存在 1删除）
}

// TableName 表名称
func (*SysNotice) TableName() string {
	return "sys_notice"
}
