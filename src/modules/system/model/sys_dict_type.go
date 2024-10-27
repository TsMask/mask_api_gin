package model

// SysDictType 字典类型表
type SysDictType struct {
	DictId     string `json:"dict_id" gorm:"column:dict_id"`                        // 字典主键
	DictName   string `json:"dict_name" gorm:"column:dict_name" binding:"required"` // 字典名称
	DictType   string `json:"dict_type" gorm:"column:dict_type" binding:"required"` // 字典类型
	Status     string `json:"status" gorm:"column:status"`                          // 状态（0停用 1正常）
	CreateBy   string `json:"create_by" gorm:"column:create_by"`                    // 创建者
	CreateTime int64  `json:"create_time" gorm:"column:create_time"`                // 创建时间
	UpdateBy   string `json:"update_by" gorm:"column:update_by"`                    // 更新者
	UpdateTime int64  `json:"update_time" gorm:"column:update_time"`                // 更新时间
	Remark     string `json:"remark" gorm:"column:remark"`                          // 备注
}

// TableName 表名称
func (*SysDictType) TableName() string {
	return "sys_dict_type"
}
