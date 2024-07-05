package model

// SysDictData 字典数据对象 sys_dict_data
type SysDictData struct {
	DictCode   string `json:"dictCode"`                     // 字典编码
	DictSort   int    `json:"dictSort"`                     // 字典排序
	DictLabel  string `json:"dictLabel" binding:"required"` // 字典标签
	DictValue  string `json:"dictValue" binding:"required"` // 字典键值
	DictType   string `json:"dictType" binding:"required"`  // 字典类型
	TagClass   string `json:"tagClass"`                     // 样式属性（样式扩展）
	TagType    string `json:"tagType"`                      // 标签类型（预设颜色）
	Status     string `json:"status"`                       // 状态（0停用 1正常）
	CreateBy   string `json:"createBy"`                     // 创建者
	CreateTime int64  `json:"createTime"`                   // 创建时间
	UpdateBy   string `json:"updateBy"`                     // 更新者
	UpdateTime int64  `json:"updateTime"`                   // 更新时间
	Remark     string `json:"remark"`                       // 备注
}
