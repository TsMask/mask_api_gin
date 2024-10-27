package model

// SysDictData 字典数据表
type SysDictData struct {
	DictCode   string `json:"dict_code" gorm:"column:dict_code"`                      // 字典编码
	DictSort   int64  `json:"dict_sort" gorm:"column:dict_sort"`                      // 字典排序
	DictLabel  string `json:"dict_label" gorm:"column:dict_label" binding:"required"` // 字典标签
	DictValue  string `json:"dict_value" gorm:"column:dict_value" binding:"required"` // 字典键值
	DictType   string `json:"dict_type" gorm:"column:dict_type" binding:"required"`   // 字典类型
	TagClass   string `json:"tag_class" gorm:"column:tag_class"`                      // 样式属性（样式扩展）
	TagType    string `json:"tag_type" gorm:"column:tag_type"`                        // 标签类型（预设颜色）
	Status     string `json:"status" gorm:"column:status"`                            // 状态（0停用 1正常）
	CreateBy   string `json:"create_by" gorm:"column:create_by"`                      // 创建者
	CreateTime int64  `json:"create_time" gorm:"column:create_time"`                  // 创建时间
	UpdateBy   string `json:"update_by" gorm:"column:update_by"`                      // 更新者
	UpdateTime int64  `json:"update_time" gorm:"column:update_time"`                  // 更新时间
	Remark     string `json:"remark" gorm:"column:remark"`                            // 备注
}

// TableName 表名称
func (*SysDictData) TableName() string {
	return "sys_dict_data"
}
