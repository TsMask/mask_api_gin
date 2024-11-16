package model

// SysDictType 字典类型表
type SysDictType struct {
	DictId     int64  `json:"dictId" gorm:"column:dict_id;primary_key"`            // 字典ID
	DictName   string `json:"dictName" gorm:"column:dict_name" binding:"required"` // 字典名称
	DictType   string `json:"dictType" gorm:"column:dict_type" binding:"required"` // 字典类型
	StatusFlag string `json:"statusFlag" gorm:"column:status_flag"`                // 状态（0停用 1正常）
	DelFlag    string `json:"-" gorm:"column:del_flag"`                            // 删除标记（0存在 1删除）
	CreateBy   string `json:"createBy" gorm:"column:create_by"`                    // 创建者
	CreateTime int64  `json:"createTime" gorm:"column:create_time"`                // 创建时间
	UpdateBy   string `json:"updateBy" gorm:"column:update_by"`                    // 更新者
	UpdateTime int64  `json:"updateTime" gorm:"column:update_time"`                // 更新时间
	Remark     string `json:"remark" gorm:"column:remark"`                         // 备注
}

// TableName 表名称
func (*SysDictType) TableName() string {
	return "sys_dict_type"
}
