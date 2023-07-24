package model

// SysDictType 字典类型对象 sys_dict_type
type SysDictType struct {
	// 字典主键
	DictID string `json:"dictId"`
	// 字典名称
	DictName string `json:"dictName" binding:"required"`
	// 字典类型
	DictType string `json:"dictType" binding:"required"`
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
