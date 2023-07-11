package model

// SysDictData 字典数据对象 sys_dict_data
type SysDictData struct {
	// 字典编码
	DictCode string `json:"dictCode"`
	// 字典排序
	DictSort int `json:"dictSort"`
	// 字典标签
	DictLabel string `json:"dictLabel"`
	// 字典键值
	DictValue string `json:"dictValue"`
	// 字典类型
	DictType string `json:"dictType"`
	// 样式属性（样式扩展）
	TagClass string `json:"tagClass"`
	// 标签类型（预设颜色）
	TagType string `json:"tagType"`
	// 状态（0停用 1正常）
	Status string `json:"status"`
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
