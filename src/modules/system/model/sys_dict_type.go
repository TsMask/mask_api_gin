package model

// SysDictType 字典类型表（dict_type 与 sys_dict_data 字符串关联）
type SysDictType struct {
	ID         int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`         // 字典ID
	CreateBy   string `json:"create_by" gorm:"column:create_by"`                    // 创建者
	CreateTime int64  `json:"create_time" gorm:"column:create_time"`                // 创建时间
	UpdateBy   string `json:"update_by" gorm:"column:update_by"`                    // 更新者
	UpdateTime int64  `json:"update_time" gorm:"column:update_time"`                // 更新时间
	Remark     string `json:"remark" gorm:"column:remark"`                          // 备注
	DictName   string `json:"dict_name" gorm:"column:dict_name" binding:"required"` // 字典名称
	DictType   string `json:"dict_type" gorm:"column:dict_type" binding:"required"` // 字典类型
	DictStatus string `json:"dict_status" gorm:"column:dict_status"`                // 状态（0停用 1正常）
}

// TableName 表名称
func (*SysDictType) TableName() string {
	return "sys_dict_type"
}
