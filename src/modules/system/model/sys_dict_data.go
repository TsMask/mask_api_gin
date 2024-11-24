package model

// SysDictData 字典数据表
type SysDictData struct {
	DataId     string `json:"dataId" gorm:"column:data_id;primary_key"`              // 数据ID
	DictType   string `json:"dictType" gorm:"column:dict_type" binding:"required"`   // 字典类型
	DataLabel  string `json:"dataLabel" gorm:"column:data_label" binding:"required"` // 数据标签
	DataValue  string `json:"dataValue" gorm:"column:data_value" binding:"required"` // 数据键值
	DataSort   int64  `json:"dataSort" gorm:"column:data_sort"`                      // 数据排序
	TagClass   string `json:"tagClass" gorm:"column:tag_class"`                      // 样式属性（样式扩展）
	TagType    string `json:"tagType" gorm:"column:tag_type"`                        // 标签类型（预设颜色）
	StatusFlag string `json:"statusFlag" gorm:"column:status_flag"`                  // 状态（0停用 1正常）
	DelFlag    string `json:"-" gorm:"column:del_flag"`                              // 删除标记（0存在 1删除）
	CreateBy   string `json:"createBy" gorm:"column:create_by"`                      // 创建者
	CreateTime int64  `json:"createTime" gorm:"column:create_time"`                  // 创建时间
	UpdateBy   string `json:"updateBy" gorm:"column:update_by"`                      // 更新者
	UpdateTime int64  `json:"updateTime" gorm:"column:update_time"`                  // 更新时间
	Remark     string `json:"remark" gorm:"column:remark"`                           // 备注
}

// TableName 表名称
func (*SysDictData) TableName() string {
	return "sys_dict_data"
}
