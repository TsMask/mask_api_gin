package model

// SysDictData 字典数据表（按 dict_type 字符串归属到字典类型）
type SysDictData struct {
	ID         int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`           // 数据ID
	CreateBy   string `json:"create_by" gorm:"column:create_by"`                      // 创建者
	CreateTime int64  `json:"create_time" gorm:"column:create_time"`                  // 创建时间
	UpdateBy   string `json:"update_by" gorm:"column:update_by"`                      // 更新者
	UpdateTime int64  `json:"update_time" gorm:"column:update_time"`                  // 更新时间
	Remark     string `json:"remark" gorm:"column:remark"`                            // 备注
	DictType   string `json:"dict_type" gorm:"column:dict_type" binding:"required"`   // 字典类型
	DataLabel  string `json:"data_label" gorm:"column:data_label" binding:"required"` // 数据标签
	DataValue  string `json:"data_value" gorm:"column:data_value" binding:"required"` // 数据键值
	DataSort   int64  `json:"data_sort" gorm:"column:data_sort"`                      // 数据排序
	DataStatus string `json:"data_status" gorm:"column:data_status"`                  // 状态（0停用 1正常）
	TagClass   string `json:"tag_class" gorm:"column:tag_class"`                      // 样式属性（样式扩展）
	TagType    string `json:"tag_type" gorm:"column:tag_type"`                        // 标签类型（预设颜色）
}

// TableName 表名称
func (*SysDictData) TableName() string {
	return "sys_dict_data"
}
