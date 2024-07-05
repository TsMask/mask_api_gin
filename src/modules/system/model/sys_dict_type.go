package model

// SysDictType 字典类型对象 sys_dict_type
type SysDictType struct {
	DictID     string `json:"dictId"`                      // 字典主键
	DictName   string `json:"dictName" binding:"required"` // 字典名称
	DictType   string `json:"dictType" binding:"required"` // 字典类型
	Status     string `json:"status"`                      // 状态（0停用 1正常）
	CreateBy   string `json:"createBy"`                    // 创建者
	CreateTime int64  `json:"createTime"`                  // 创建时间
	UpdateBy   string `json:"updateBy"`                    // 更新者
	UpdateTime int64  `json:"updateTime"`                  // 更新时间
	Remark     string `json:"remark"`                      // 备注
}
